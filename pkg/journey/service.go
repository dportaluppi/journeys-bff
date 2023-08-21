package journey

import (
	"context"
	"errors"
	"fmt"
	"github.com/dportaluppi/journeys-bff/pkg/product"
	"github.com/dportaluppi/journeys-bff/pkg/segment"
)

type Service interface {
	CreateJourney(ctx context.Context, j *JourneyWriteModel) (*JourneyWriteModel, error)
	GetJourneys(ctx context.Context, accountId string) ([]JourneyReadModel, error)
	GetJourneyByID(ctx context.Context, accountId, journeyID string) (*JourneyReadModel, error)
}

type service struct {
	journeyRepo Repository
	productRepo product.Searcher
	segmentRepo segment.Getter
}

func NewService(journeyRepo Repository, productRepo product.Searcher, segmentRepo segment.Getter) *service {
	return &service{journeyRepo: journeyRepo, productRepo: productRepo, segmentRepo: segmentRepo}
}

func (s *service) CreateJourney(ctx context.Context, j *JourneyWriteModel) (*JourneyWriteModel, error) {
	// TODO: add business logic here
	// validate journey: audiences, products, storefront, botId, startAt, endAt
	// review selection strategy for all products and all users.

	// generate experiences ids
	for i := range j.Experiences {
		j.Experiences[i].ExperienceID = s.journeyRepo.UniqueID()
	}

	// generate stepId
	// set stepId to experience who has type `order_reminder_message`

	// update rule notification

	return s.journeyRepo.Create(ctx, j)
}

func (s *service) GetJourneys(ctx context.Context, accountId string) ([]JourneyReadModel, error) {
	journeys, err := s.journeyRepo.GetJourneys(ctx, accountId)
	if err != nil {
		return nil, fmt.Errorf("error fetching journeys: %w", err)
	}

	readModels := make([]JourneyReadModel, 0, len(journeys))
	for _, j := range journeys {
		jrm, err := s.mapWriteModelToReadModel(ctx, &j)
		if err != nil {
			return nil, fmt.Errorf("error mapping write model to read model: %w", err)
		}
		readModels = append(readModels, *jrm)
	}

	return readModels, nil
}

func (s *service) GetJourneyByID(ctx context.Context, accountId, journeyID string) (*JourneyReadModel, error) {
	j, err := s.journeyRepo.GetJourneyByID(ctx, accountId, journeyID)
	if err != nil {
		return nil, fmt.Errorf("error fetching journey by ID: %w", err)
	}

	return s.mapWriteModelToReadModel(ctx, j)
}

func (s *service) fetchProducts(ctx context.Context, j *JourneyWriteModel) ([]product.Product, error) {
	allSKUs := append(j.Recommendations.AllowList, j.Recommendations.ExcludeList...)
	products := make([]product.Product, 0, len(allSKUs))

	for _, sku := range allSKUs {
		aProduct, err := s.productRepo.GetBySKU(ctx, j.Storefront, sku)
		if err != nil {
			if errors.Is(err, product.ErrProductNotFound) {
				aProduct = product.Product{
					SKU: sku,
				}
			} else {
				return nil, err
			}
		}
		products = append(products, aProduct)
	}

	return products, nil
}

func (s *service) fetchSegments(ctx context.Context, j *JourneyWriteModel) ([]segment.Segment, error) {
	segments := make([]segment.Segment, 0, len(j.Audiences))

	for _, segmentID := range j.Audiences {
		seg, err := s.segmentRepo.GetByID(ctx, &segment.Filter{
			Provider: j.AccountID, // TODO: What should it be? the storefrontName, the accountID or the customerID? Is it an attribute that we should store in the journey?
			ID:       segmentID,
		})
		if err != nil {
			return nil, err
		}
		segments = append(segments, seg)
	}

	return segments, nil
}

func (s *service) mapWriteModelToReadModel(ctx context.Context, j *JourneyWriteModel) (*JourneyReadModel, error) {
	products, err := s.fetchProducts(ctx, j)
	if err != nil {
		return nil, fmt.Errorf("error fetching products for journey: %w", err)
	}

	segments, err := s.fetchSegments(ctx, j)
	if err != nil {
		return nil, fmt.Errorf("error fetching segments for journey: %w", err)
	}

	jrm := &JourneyReadModel{
		ID:          j.ID,
		AccountID:   j.AccountID,
		Type:        j.Type,
		Experiences: j.Experiences,
		BotID:       j.BotID,
		Storefront:  j.Storefront,
		StartAt:     j.StartAt,
		EndAt:       j.EndAt,
		Audiences:   segments,
		Recommendations: RecommendationsReadModel{
			AllowList:   products[:len(j.Recommendations.AllowList)],
			ExcludeList: products[len(j.Recommendations.AllowList):],
		},
	}
	return jrm, nil
}

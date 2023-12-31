package segment

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net/http"
)

var ErrSegmentNotFound = errors.New("segment not found")

type HTTPAudienceRepo struct {
	baseURL string
	token   string
}

func NewHTTPRepo(baseURL, token string) *HTTPAudienceRepo {
	return &HTTPAudienceRepo{baseURL: baseURL, token: token}
}

func (repo *HTTPAudienceRepo) GetSegments(ctx context.Context, filter *Filter, pageSize int, pageNumber int) (*SegmentsResponse, error) {
	var repoResponse RepoSegmentsResponse

	url := repo.baseURL + fmt.Sprintf("?provider=%s&page=%d&pageSize=%d", filter.Provider, pageNumber, pageSize)
	if filter.Name != "" {
		url += "&name=" + filter.Name
	}

	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+repo.token)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed with status: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bodyBytes, &repoResponse)
	if err != nil {
		return nil, err
	}

	domainAudiences := make([]Segment, len(repoResponse.Data))
	totalCustomers := 0
	for i, repoAud := range repoResponse.Data {
		domainAudiences[i] = toDomain(repoAud)
		totalCustomers += domainAudiences[i].LastGeneratedSize
	}

	response := &SegmentsResponse{
		Data: Data{
			Audiences: domainAudiences,
			Customers: CustomersInfo{Total: totalCustomers}, // TODO: get total customers from new endpoint
		},
		Metadata: repoResponse.Metadata,
	}

	return response, nil
}

func (repo *HTTPAudienceRepo) GetByID(ctx context.Context, filter *Filter) (Segment, error) {
	url := fmt.Sprintf("%s/%s?provider=%s", repo.baseURL, filter.ID, filter.Provider)

	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return Segment{}, err
	}

	req.Header.Add("Authorization", "Bearer "+repo.token)
	resp, err := client.Do(req)
	if err != nil {
		return Segment{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return Segment{}, ErrSegmentNotFound
		}

		return Segment{}, fmt.Errorf("failed with status: %d", resp.StatusCode)
	}

	var repoSegment RepoSegment
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return Segment{}, err
	}

	err = json.Unmarshal(bodyBytes, &repoSegment)
	if err != nil {
		return Segment{}, err
	}

	return toDomain(repoSegment), nil
}

func toDomain(segmentAud RepoSegment) Segment {
	customerAttributesFilterJSON, err := json.Marshal(segmentAud.Filters.CustomerAttributesFilter)
	if err != nil {
		// TODO: log error
	}

	eventsFilterJSON, err := json.Marshal(segmentAud.Filters.EventsFilter)
	if err != nil {
		// TODO: log error
	}

	return Segment{
		ID:   segmentAud.ID,
		Name: segmentAud.Name,
		Filters: SegmentFilters{
			CustomerAttributesFilter: string(customerAttributesFilterJSON),
			EventsFilter:             string(eventsFilterJSON),
		},
		Tags:              segmentAud.Tags,
		LastGeneratedSize: segmentAud.LastGeneratedSize,
		UpdatedAt:         segmentAud.UpdatedAt,
	}
}

// --------- Repository model -------------

type RepoSegment struct {
	ID                string      `json:"id"`
	Provider          string      `json:"provider"`
	Name              string      `json:"name"`
	Filters           RepoFilters `json:"filters"`
	Tags              []string    `json:"tags"`
	LastGeneratedSize int         `json:"last_generated_size"`
	UpdatedAt         string      `json:"last_generated_at"`
	Folder            string      `json:"folder"`
}

type RepoFilters struct {
	CustomerAttributesFilter map[string]interface{} `json:"customerAttributesFilter"`
	EventsFilter             map[string]interface{} `json:"eventsFilter"`
}

type RepoSegmentsResponse struct {
	Data     []RepoSegment `json:"data"`
	Metadata Metadata      `json:"metadata"`
}

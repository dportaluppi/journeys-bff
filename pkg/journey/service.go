package journey

import (
	"context"
)

type Service interface {
	CreateJourney(ctx context.Context, j *Journey) (*Journey, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) CreateJourney(ctx context.Context, j *Journey) (*Journey, error) {
	// TODO: add business logic here
	// validate journey
	// validate account
	// validate segments or audiences: if id is `*` then we don't need to validate
	// validate products: if id is `*` then we don't need to validate.

	// generate experiences ids
	for i := range j.Experiences {
		j.Experiences[i].ExperienceID = s.repo.UniqueID()
	}

	// generate stepId
	// set stepId to experience who has type `order_reminder_message`

	// update rule notification

	return s.repo.Create(ctx, j)
}

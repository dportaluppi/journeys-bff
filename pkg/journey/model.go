package journey

import (
	"github.com/dportaluppi/journeys-bff/pkg/product"
	"github.com/dportaluppi/journeys-bff/pkg/segment"
	"time"
)

type Experience struct {
	ExperienceID string                 `json:"experienceId"`
	Type         string                 `json:"type"`
	Parameters   map[string]interface{} `json:"parameters,omitempty"`
	EnabledAt    time.Time              `json:"enabledAt"`
}

type RecommendationsWriteModel struct {
	AllowList   []string `json:"allowList"`
	ExcludeList []string `json:"excludeList"`
}

type JourneyWriteModel struct {
	ID              string                    `json:"journeyId"`
	AccountID       string                    `json:"accountId"`
	Type            string                    `json:"type"`
	Experiences     []Experience              `json:"experiences"`
	BotID           string                    `json:"botId"`
	Storefront      string                    `json:"storefront"`
	Audiences       []string                  `json:"audiences"`
	Recommendations RecommendationsWriteModel `json:"recommendations"`
	StartAt         time.Time                 `json:"startAt"`
	EndAt           time.Time                 `json:"endAt"`
}

type RecommendationsReadModel struct {
	AllowList   []product.Product `json:"allowList"`
	ExcludeList []product.Product `json:"excludeList"`
}

type JourneyReadModel struct {
	ID              string                   `json:"journeyId"`
	AccountID       string                   `json:"accountId"`
	Type            string                   `json:"type"`
	Experiences     []Experience             `json:"experiences"`
	BotID           string                   `json:"botId"`
	Storefront      string                   `json:"storefront"`
	Audiences       []segment.Segment        `json:"audiences"`
	Recommendations RecommendationsReadModel `json:"recommendations"`
	StartAt         time.Time                `json:"startAt"`
	EndAt           time.Time                `json:"endAt"`
}

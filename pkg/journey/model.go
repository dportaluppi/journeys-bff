package journey

import (
	"time"
)

type Experience struct {
	ExperienceID string                 `json:"experienceId"`
	Type         string                 `json:"type"`
	Parameters   map[string]interface{} `json:"parameters,omitempty"`
	EnabledAt    time.Time              `json:"enabledAt"`
}

type Recommendations struct {
	AllowList []string `json:"allowList"`
	DenyList  []string `json:"denyList"`
}

type Journey struct {
	ID              string          `json:"journeyId"`
	AccountID       string          `json:"accountId"`
	Type            string          `json:"type"`
	Experiences     []Experience    `json:"experiences"`
	BotID           string          `json:"botId"`
	Storefront      string          `json:"storefront"`
	Audiences       []string        `json:"audiences"`
	Recommendations Recommendations `json:"recommendations"`
	StartAt         time.Time       `json:"startAt"`
	EndAt           time.Time       `json:"endAt"`
}

package audience

import "context"

type Audience struct {
	ID                string          `json:"id"`
	Name              string          `json:"name"`
	Filters           AudienceFilters `json:"filters"`
	Tags              []string        `json:"tags"`
	LastGeneratedSize int             `json:"lastGeneratedSize"`
	UpdatedAt         string          `json:"updatedAt"`
}

type AudienceFilters struct {
	CustomerAttributesFilter string `json:"customerAttributesFilter"`
	EventsFilter             string `json:"eventsFilter"`
}

type CustomersInfo struct {
	Total int `json:"total"`
}

type Metadata struct {
	Page       int `json:"page"`
	PageSize   int `json:"pageSize"`
	TotalPages int `json:"totalPages"`
	Total      int `json:"total"`
}

type Data struct {
	Audiences []Audience    `json:"audiences"`
	Customers CustomersInfo `json:"customers"`
}

type AudiencesResponse struct {
	Data     Data     `json:"data"`
	Metadata Metadata `json:"metadata"`
}

type Filter struct {
	Provider string
	Query    string
}

type Getter interface {
	GetAudiences(ctx context.Context, filter *Filter, pageSize int, pageNumber int) (*AudiencesResponse, error)
}

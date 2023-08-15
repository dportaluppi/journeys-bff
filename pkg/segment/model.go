package segment

import "context"

type Segment struct {
	ID                string         `json:"id"`
	Name              string         `json:"name"`
	Filters           SegmentFilters `json:"filters"`
	Tags              []string       `json:"tags"`
	LastGeneratedSize int            `json:"lastGeneratedSize"`
	UpdatedAt         string         `json:"lastGeneratedAt"`
}

type SegmentFilters struct {
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
	Audiences []Segment     `json:"audiences"`
	Customers CustomersInfo `json:"customers"`
}

type SegmentsResponse struct {
	Data     Data     `json:"data"`
	Metadata Metadata `json:"metadata"`
}

type Filter struct {
	Provider string
	Query    string
}

type Getter interface {
	GetSegments(ctx context.Context, filter *Filter, pageSize int, pageNumber int) (*SegmentsResponse, error)
}

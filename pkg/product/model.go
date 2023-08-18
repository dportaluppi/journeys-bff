package product

import "context"

type Product struct {
	ImageURL    []string `json:"imageURL"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	SKU         string   `json:"sku"`
	Category    string   `json:"category"`
}

type SearchResult struct {
	Data     []Product `json:"data"`
	Metadata Metadata  `json:"metadata"`
}

type Metadata struct {
	PageSize     int `json:"pageSize"`
	PageNumber   int `json:"pageNumber"`
	NextPage     int `json:"nextPage"`
	PreviousPage int `json:"previousPage"`
	Total        int `json:"total"`
	TotalPages   int `json:"totalPages"`
}
type Filter struct {
	Storefront string
	Name       string
	SKU        string
	Category   string
}
type Searcher interface {
	Search(ctx context.Context, filter *Filter, pageSize int, pageNumber int) (SearchResult, error)
	GetBySKU(ctx context.Context, storefront, sku string) (Product, error)
}

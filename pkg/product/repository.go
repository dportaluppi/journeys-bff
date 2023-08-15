package product

import (
	"context"
	"github.com/machinebox/graphql"
)

type repository struct {
	graphQLClient *graphql.Client
}

func NewProductRepository(endpoint string) *repository {
	client := graphql.NewClient(endpoint)
	return &repository{graphQLClient: client}
}

func (r *repository) Search(ctx context.Context, filter *Filter, pageSize int, pageNumber int) (SearchResult, error) {
	req := graphql.NewRequest(`
		query GetProducts($storefrontName: String!, $pagination: PaginationInput, $filter: FilterProductInput) {
			getProducts(storefrontName: $storefrontName, pagination: $pagination, filter: $filter) {
				products {
					imageURL
					name
					description
					sku
					category
				}
				pagination {
					pageNumber
					pageSize
					nextPage
					previousPage
					total
					totalPages
				}
			}
		}
	`)

	req.Var("storefrontName", filter.Storefront)
	req.Var("pagination", map[string]int{
		"pageNumber": pageNumber,
		"pageSize":   pageSize,
	})
	req.Var("filter", map[string]interface{}{
		"isActive": true,
		"name":     filter.Name,
	})

	var resp struct {
		GetProducts struct {
			Products   []Product `json:"products"`
			Pagination Metadata  `json:"pagination"`
		} `json:"getProducts"`
	}

	if err := r.graphQLClient.Run(ctx, req, &resp); err != nil {
		return SearchResult{}, err
	}

	searchResult := SearchResult{
		Data:     resp.GetProducts.Products,
		Metadata: resp.GetProducts.Pagination,
	}

	return searchResult, nil
}
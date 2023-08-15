package product

import (
	"context"
)

type productSearcher struct {
	repo Repository
}

func NewSearcher(repo Repository) Searcher {
	return &productSearcher{repo: repo}
}

func (s *productSearcher) Search(ctx context.Context, filter *Filter, pageSize int, pageNumber int) (SearchResult, error) {
	// TODO: add transformation logic here
	return s.repo.Search(ctx, filter, pageSize, pageNumber)
}

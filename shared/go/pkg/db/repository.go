package db

import (
	"context"

	"github.com/FreibergVlad/url-shortener/shared/go/pkg/db/filters"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/db/sort"
)

type PaginatedQuery struct {
	Filter   filters.Filter
	Sort     *sort.Sort
	PageSize int32
	PageNum  int32
}

type PaginatedResult[T any] struct {
	Data  []T
	Total int32
}

type Repository[T any] interface {
	Create(ctx context.Context, resource *T) error
	FindOne(ctx context.Context, filter filters.Filter) (*T, error)
	FindMany(ctx context.Context, query PaginatedQuery) (*PaginatedResult[*T], error)
	ReplaceOne(ctx context.Context, filter filters.Filter, resource *T) error
	UpdateOne(ctx context.Context, filter filters.Filter, update map[string]any) (*T, error)
}

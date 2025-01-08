package mocks

import (
	"context"

	"github.com/FreibergVlad/url-shortener/shared/go/pkg/db"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/db/filters"
	"github.com/stretchr/testify/mock"
)

type MockedRepository[T any] struct {
	mock.Mock
}

func (r *MockedRepository[T]) Create(ctx context.Context, resource *T) error {
	args := r.Called(ctx, resource)
	return args.Error(0)
}

func (r *MockedRepository[T]) FindOne(ctx context.Context, filter filters.Filter) (*T, error) {
	args := r.Called(ctx, filter)
	return args.Get(0).(*T), args.Error(1)
}

func (r *MockedRepository[T]) FindMany(ctx context.Context, query db.PaginatedQuery) (*db.PaginatedResult[*T], error) {
	args := r.Called(ctx, query)
	return args.Get(0).(*db.PaginatedResult[*T]), args.Error(1)
}

func (r *MockedRepository[T]) ReplaceOne(ctx context.Context, filter filters.Filter, resource *T) error {
	args := r.Called(ctx, filter, resource)
	return args.Error(0)
}

func (r *MockedRepository[T]) UpdateOne(ctx context.Context, filter filters.Filter, update map[string]any) (*T, error) {
	args := r.Called(ctx, filter, update)
	return args.Get(0).(*T), args.Error(1)
}

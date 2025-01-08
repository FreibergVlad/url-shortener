package mocks

import (
	"context"

	"github.com/FreibergVlad/url-shortener/shared/go/pkg/cache"
	"github.com/stretchr/testify/mock"
)

type MockedCacheProvider[T any] struct {
	mock.Mock
}

func (c *MockedCacheProvider[T]) Set(ctx context.Context, item *cache.Item[T]) error {
	args := c.Called(ctx, item)
	return args.Error(0)
}

func (c *MockedCacheProvider[T]) Get(ctx context.Context, key string) (*T, error) {
	args := c.Called(ctx, key)
	return args.Get(0).(*T), args.Error(1)
}

func (c *MockedCacheProvider[T]) Delete(ctx context.Context, key string) error {
	args := c.Called(ctx, key)
	return args.Error(0)
}

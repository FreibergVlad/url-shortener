package cache

import (
	"context"
	"errors"
	"time"
)

var ErrCacheMiss = errors.New("cache miss")

type Item[T any] struct {
	Key   string
	Value *T
	TTL   time.Duration
}

type Cache[T any] interface {
	Set(ctx context.Context, item *Item[T]) error
	Get(ctx context.Context, key string) (*T, error)
	Delete(ctx context.Context, key string) error
}

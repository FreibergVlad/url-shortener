package redis

import (
	"context"
	"errors"

	"github.com/FreibergVlad/url-shortener/shared/go/pkg/cache"
	redisCacheImpl "github.com/go-redis/cache/v9"
)

type Cache[T any] struct {
	redis *redisCacheImpl.Cache
}

func New[T any](opts *redisCacheImpl.Options) *Cache[T] {
	return &Cache[T]{redis: redisCacheImpl.New(opts)}
}

func (c *Cache[T]) Set(ctx context.Context, item *cache.Item[T]) error {
	return c.redis.Set(&redisCacheImpl.Item{
		Ctx:   ctx,
		Key:   item.Key,
		Value: item.Value,
		TTL:   item.TTL,
	})
}

func (c *Cache[T]) Get(ctx context.Context, key string) (*T, error) {
	var value T
	err := c.redis.Get(ctx, key, &value)
	if errors.Is(err, redisCacheImpl.ErrCacheMiss) {
		return nil, cache.ErrCacheMiss
	}
	if err != nil {
		return nil, err
	}
	return &value, err
}

func (c *Cache[T]) Delete(ctx context.Context, key string) error {
	return c.redis.Delete(ctx, key)
}

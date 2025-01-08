package redis

import (
	"context"

	"github.com/FreibergVlad/url-shortener/shared/go/pkg/cache"
	redisCacheImpl "github.com/go-redis/cache/v9"
)

type redisCache[T any] struct {
	redis *redisCacheImpl.Cache
}

func New[T any](opts *redisCacheImpl.Options) *redisCache[T] {
	return &redisCache[T]{redis: redisCacheImpl.New(opts)}
}

func (c *redisCache[T]) Set(ctx context.Context, item *cache.Item[T]) error {
	return c.redis.Set(&redisCacheImpl.Item{
		Ctx:   ctx,
		Key:   item.Key,
		Value: item.Value,
		TTL:   item.TTL,
	})
}

func (c *redisCache[T]) Get(ctx context.Context, key string) (*T, error) {
	var value T
	err := c.redis.Get(ctx, key, &value)
	if err == redisCacheImpl.ErrCacheMiss {
		return nil, cache.ErrCacheMiss
	}
	if err != nil {
		return nil, err
	}
	return &value, err
}

func (c *redisCache[T]) Delete(ctx context.Context, key string) error {
	return c.redis.Delete(ctx, key)
}

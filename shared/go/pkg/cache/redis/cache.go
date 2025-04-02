package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/FreibergVlad/url-shortener/shared/go/pkg/cache"
	redisImpl "github.com/redis/go-redis/v9"
)

type Cache[T any] struct {
	redis *redisImpl.Client
}

func New[T any](client *redisImpl.Client) *Cache[T] {
	return &Cache[T]{redis: client}
}

func (c *Cache[T]) Set(ctx context.Context, item *cache.Item[T]) error {
	value, err := json.Marshal(item.Value)
	if err != nil {
		return fmt.Errorf("redis: failed to marshal value: %w", err)
	}
	return c.redis.Set(ctx, item.Key, value, item.TTL).Err()
}

func (c *Cache[T]) Get(ctx context.Context, key string) (*T, error) {
	rawValue, err := c.redis.Get(ctx, key).Result()
	if errors.Is(err, redisImpl.Nil) {
		return nil, cache.ErrCacheMiss
	}
	if err != nil {
		return nil, err
	}

	var value T
	if err = json.Unmarshal([]byte(rawValue), &value); err != nil {
		return nil, fmt.Errorf("redis: failed to unmarshal value: %w", err)
	}
	return &value, err
}

func (c *Cache[T]) Delete(ctx context.Context, keys ...string) error {
	return c.redis.Del(ctx, keys...).Err()
}

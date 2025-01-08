package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/FreibergVlad/url-shortener/shared/go/pkg/cache"
	"github.com/FreibergVlad/url-shortener/short-url-management-service/internal/db"
	"github.com/rs/zerolog/log"
)

const ShortURLCacheTTL = time.Hour

type ShortURLCache struct {
	cacheProvider cache.Cache
	ttl           time.Duration
}

func New(cacheProvider cache.Cache, ttl time.Duration) *ShortURLCache {
	return &ShortURLCache{cacheProvider: cacheProvider, ttl: ttl}
}

func (c *ShortURLCache) Set(ctx context.Context, shortURL *db.ShortURL) {
	item := cache.Item{
		Key:   oneShortURLCacheKey(shortURL.Domain, shortURL.Alias),
		Value: shortURL,
		TTL:   ShortURLCacheTTL,
	}
	if err := c.cacheProvider.Set(ctx, &item); err != nil {
		log.Error().Err(err).Msg("Error caching short URL")
	}
}

func (c *ShortURLCache) Get(ctx context.Context, domain, alias string) (*db.ShortURL, error) {
	var shortURL *db.ShortURL
	err := c.cacheProvider.Get(ctx, oneShortURLCacheKey(domain, alias), &shortURL)
	if err != nil {
		if err != cache.ErrCacheMiss {
			log.Error().Err(err).Msg("Error getting short URL from cache")
		}
		return nil, err
	}
	return shortURL, nil
}

func (c *ShortURLCache) Delete(ctx context.Context, shortURL *db.ShortURL) error {
	err := c.cacheProvider.Delete(ctx, oneShortURLCacheKey(shortURL.Domain, shortURL.Alias))
	if err != nil {
		log.Error().Err(err).Msg("Error deleting short URL from cache")
	}
	return err
}

func oneShortURLCacheKey(domain, alias string) string {
	return fmt.Sprintf("short_url:%s:%s", domain, alias)
}

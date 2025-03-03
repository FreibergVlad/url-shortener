package shorturls

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/FreibergVlad/url-shortener/shared/go/pkg/cache"
	"github.com/FreibergVlad/url-shortener/short-url-management-service/internal/db/schema"
	"github.com/rs/zerolog/log"
)

type CachedRepository struct {
	repository Repository
	cache      cache.Cache[schema.ShortURL]
	ttl        time.Duration
}

func NewCachedRepository(
	repository Repository,
	cache cache.Cache[schema.ShortURL],
	ttl time.Duration,
) *CachedRepository {
	return &CachedRepository{repository: repository, cache: cache, ttl: ttl}
}

func (r *CachedRepository) Create(ctx context.Context, shorturl *schema.ShortURL) error {
	err := r.repository.Create(ctx, shorturl)
	if err != nil {
		return err
	}
	r.cacheShortURL(ctx, shorturl)
	return nil
}

func (r *CachedRepository) GetByGlobalKey(ctx context.Context, key schema.ShortURLGlobalKey) (*schema.ShortURL, error) {
	shorturl := r.getShortURLFromCache(ctx, key)
	if shorturl != nil {
		return shorturl, nil
	}
	shorturl, err := r.repository.GetByGlobalKey(ctx, key)
	if err != nil {
		return nil, err
	}
	r.cacheShortURL(ctx, shorturl)
	return shorturl, nil
}

func (r *CachedRepository) GetByOrganizationIDAndGlobalKey(
	ctx context.Context, organizationID string, key schema.ShortURLGlobalKey,
) (*schema.ShortURL, error) {
	shorturl := r.getShortURLFromCache(ctx, key)
	if shorturl != nil {
		if shorturl.OrganizationID != organizationID {
			return nil, ErrNotFound
		}
		return shorturl, nil
	}
	shorturl, err := r.repository.GetByOrganizationIDAndGlobalKey(ctx, organizationID, key)
	if err != nil {
		return nil, err
	}
	r.cacheShortURL(ctx, shorturl)
	return shorturl, nil
}

func (r *CachedRepository) GetByOrganizationKeyOrGlobalKey(
	ctx context.Context, orgKey schema.ShortURLOrganizationKey, globalKey schema.ShortURLGlobalKey,
) (*schema.ShortURL, error) {
	return r.repository.GetByOrganizationKeyOrGlobalKey(ctx, orgKey, globalKey)
}

func (r *CachedRepository) ListByOrganizationID(
	ctx context.Context, organizationID string, pageNum, pageSize int,
) (*PaginatedResult, error) {
	return r.repository.ListByOrganizationID(ctx, organizationID, pageNum, pageSize)
}

func (r *CachedRepository) ReplaceByOrganizationIDAndGlobalKey(
	ctx context.Context, organizationID string, key schema.ShortURLGlobalKey, replaceWith *schema.ShortURL,
) error {
	err := r.repository.ReplaceByOrganizationIDAndGlobalKey(ctx, organizationID, key, replaceWith)
	if err != nil {
		return err
	}
	r.cacheShortURL(ctx, replaceWith)
	return nil
}

func (r *CachedRepository) DeleteByOrganizationIDAndGlobalKey(
	ctx context.Context, organizationID string, key schema.ShortURLGlobalKey,
) (*schema.ShortURL, error) {
	shorturl, err := r.repository.DeleteByOrganizationIDAndGlobalKey(ctx, organizationID, key)
	if err != nil {
		return nil, err
	}
	cacheKey := r.shortURLCacheKey(key)
	if err = r.cache.Delete(ctx, cacheKey); err != nil {
		log.Err(err).Msgf("error deleting short url from cache %s", cacheKey)
	}
	return shorturl, nil
}

func (r *CachedRepository) getShortURLFromCache(ctx context.Context, key schema.ShortURLGlobalKey) *schema.ShortURL {
	cacheKey := r.shortURLCacheKey(key)
	shorturl, err := r.cache.Get(ctx, cacheKey)
	if err == nil {
		return shorturl
	}
	if !errors.Is(err, cache.ErrCacheMiss) {
		log.Err(err).Msgf("error getting short url from cache %s", cacheKey)
	}
	return nil
}

func (r *CachedRepository) cacheShortURL(ctx context.Context, shorturl *schema.ShortURL) {
	item := cache.Item[schema.ShortURL]{
		Key:   r.shortURLCacheKey(shorturl.GlobalKey()),
		Value: shorturl,
		TTL:   r.ttl,
	}
	if err := r.cache.Set(ctx, &item); err != nil {
		log.Err(err).Msgf("error caching short url %s", item.Key)
	}
}

func (r *CachedRepository) shortURLCacheKey(key schema.ShortURLGlobalKey) string {
	return fmt.Sprintf("shorturls:domain:%s:alias:%s", key.Domain, key.Alias)
}

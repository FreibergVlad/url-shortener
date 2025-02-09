package users

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/FreibergVlad/url-shortener/auth-service/internal/db/schema"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/cache"
	"github.com/rs/zerolog/log"
)

const (
	cacheKeyTypeID    = "id"
	cacheKeyTypeEmail = "email"
)

type CachedRepository struct {
	repository Repository
	cache      cache.Cache[schema.User]
	ttl        time.Duration
}

func NewCachedRepository(
	repository Repository,
	cache cache.Cache[schema.User],
	ttl time.Duration,
) *CachedRepository {
	return &CachedRepository{
		repository: repository,
		cache:      cache,
		ttl:        ttl,
	}
}

func (r *CachedRepository) Create(ctx context.Context, user *schema.User) error {
	err := r.repository.Create(ctx, user)
	if err != nil {
		return err
	}

	r.cacheUser(ctx, user)

	return nil
}

func (r *CachedRepository) GetByID(ctx context.Context, userID string) (*schema.User, error) {
	user := r.tryCacheGet(ctx, r.cacheKey(cacheKeyTypeID, userID))
	if user != nil {
		return user, nil
	}

	user, err := r.repository.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	r.cacheUser(ctx, user)

	return user, nil
}

func (r *CachedRepository) GetByEmail(ctx context.Context, email string) (*schema.User, error) {
	user := r.tryCacheGet(ctx, r.cacheKey(cacheKeyTypeEmail, email))
	if user != nil {
		return user, nil
	}

	user, err := r.repository.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	r.cacheUser(ctx, user)

	return user, nil
}

func (r *CachedRepository) cacheUser(ctx context.Context, user *schema.User) {
	for _, key := range []string{r.cacheKey(cacheKeyTypeID, user.ID), r.cacheKey(cacheKeyTypeEmail, user.Email)} {
		err := r.cache.Set(ctx, &cache.Item[schema.User]{
			Key:   key,
			Value: user,
			TTL:   r.ttl,
		})
		if err != nil {
			log.Error().Err(err).Msgf("error setting %s:%s to cache", key, user)
		}
	}
}

func (r *CachedRepository) tryCacheGet(ctx context.Context, key string) *schema.User {
	resource, err := r.cache.Get(ctx, key)
	if err == nil {
		return resource
	}
	if !errors.Is(err, cache.ErrCacheMiss) {
		log.Error().Err(err).Msgf("error getting %s from cache", key)
	}
	return nil
}

func (r *CachedRepository) cacheKey(keyType, keyValue string) string {
	return fmt.Sprintf("users:%s:%s", keyType, keyValue)
}

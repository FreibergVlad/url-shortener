package users

import (
	"context"
	"fmt"
	"time"

	"github.com/FreibergVlad/url-shortener/auth-service/internal/db/schema"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/cache"
	"github.com/rs/zerolog/log"
)

const (
	cacheKeyTypeId    = "id"
	cacheKeyTypeEmail = "email"
)

type cachedUserRepository struct {
	repository Repository
	cache      cache.Cache[schema.User]
	ttl        time.Duration
}

func NewCachedUserRepository(
	repository Repository,
	cache cache.Cache[schema.User],
	ttl time.Duration,
) *cachedUserRepository {
	return &cachedUserRepository{
		repository: repository,
		cache:      cache,
		ttl:        ttl,
	}
}

func (r *cachedUserRepository) Create(ctx context.Context, user *schema.User) error {
	err := r.repository.Create(ctx, user)
	if err != nil {
		return err
	}

	r.cacheUser(ctx, user)

	return nil
}

func (r *cachedUserRepository) GetById(ctx context.Context, id string) (*schema.User, error) {
	user := r.tryCacheGet(ctx, r.cacheKey(cacheKeyTypeId, id))
	if user != nil {
		return user, nil
	}

	user, err := r.repository.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	r.cacheUser(ctx, user)

	return user, nil
}

func (r *cachedUserRepository) GetByEmail(ctx context.Context, email string) (*schema.User, error) {
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

func (r *cachedUserRepository) cacheUser(ctx context.Context, user *schema.User) {
	for _, key := range []string{r.cacheKey(cacheKeyTypeId, user.ID), r.cacheKey(cacheKeyTypeEmail, user.Email)} {
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

func (r *cachedUserRepository) tryCacheGet(ctx context.Context, key string) *schema.User {
	resource, err := r.cache.Get(ctx, key)
	if err == nil {
		return resource
	}
	if err != cache.ErrCacheMiss {
		log.Error().Err(err).Msgf("error getting %s from cache", key)
	}
	return nil
}

func (r *cachedUserRepository) cacheKey(keyType, keyValue string) string {
	return fmt.Sprintf("users:%s:%s", keyType, keyValue)
}

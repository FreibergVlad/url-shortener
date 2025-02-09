package organizations

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/FreibergVlad/url-shortener/auth-service/internal/db/schema"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/cache"
	"github.com/rs/zerolog/log"
)

type CachedRepository struct {
	repository                  Repository
	organizationMembershipCache cache.Cache[schema.OrganizationMemberships]
	ttl                         time.Duration
}

func NewCachedRepository(
	repository Repository,
	organizationMembershipCache cache.Cache[schema.OrganizationMemberships],
	ttl time.Duration,
) *CachedRepository {
	return &CachedRepository{
		repository:                  repository,
		organizationMembershipCache: organizationMembershipCache,
		ttl:                         ttl,
	}
}

func (r *CachedRepository) Create(ctx context.Context, organization *schema.Organization) error {
	if err := r.repository.Create(ctx, organization); err != nil {
		return err
	}

	r.uncacheOrganizationMemberships(ctx, organization.CreatedBy)

	return nil
}

func (r *CachedRepository) CreateOrganizationMembership(
	ctx context.Context, params *CreateOrganizationMembershipParams,
) error {
	if err := r.repository.CreateOrganizationMembership(ctx, params); err != nil {
		return err
	}

	r.uncacheOrganizationMemberships(ctx, params.UserID)

	return nil
}

func (r *CachedRepository) ListOrganizationMembershipsByUserID(
	ctx context.Context, userID string,
) (schema.OrganizationMemberships, error) {
	memberships := r.getOrganizationMembershipsFromCache(ctx, userID)
	if memberships != nil {
		return memberships, nil
	}

	memberships, err := r.repository.ListOrganizationMembershipsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	r.cacheOrganizationMemberships(ctx, memberships, userID)

	return memberships, err
}

func (r *CachedRepository) getOrganizationMembershipsFromCache(
	ctx context.Context, userID string,
) schema.OrganizationMemberships {
	key := organizationMembershipsCacheKey(userID)

	memberships, err := r.organizationMembershipCache.Get(ctx, key)
	if err == nil {
		return *memberships
	}

	if !errors.Is(err, cache.ErrCacheMiss) {
		log.Error().Err(err).Msgf("error getting %s from cache", key)
	}

	return nil
}

func (r *CachedRepository) cacheOrganizationMemberships(
	ctx context.Context, memberships schema.OrganizationMemberships, userID string,
) {
	key := organizationMembershipsCacheKey(userID)
	item := cache.Item[schema.OrganizationMemberships]{
		Key:   key,
		Value: &memberships,
		TTL:   r.ttl,
	}
	if err := r.organizationMembershipCache.Set(ctx, &item); err != nil {
		log.Error().Err(err).Msgf("error setting %s:%v+ to cache", key, memberships)
	}
}

func (r *CachedRepository) uncacheOrganizationMemberships(ctx context.Context, userID string) {
	key := organizationMembershipsCacheKey(userID)
	if err := r.organizationMembershipCache.Delete(ctx, key); err != nil {
		log.Error().Err(err).Msgf("error deleting %s from cache", key)
	}
}

func organizationMembershipsCacheKey(userID string) string {
	return fmt.Sprintf("organization-memberships:user-id:%s", userID)
}

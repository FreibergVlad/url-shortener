package organizations

import (
	"context"
	"fmt"
	"time"

	"github.com/FreibergVlad/url-shortener/auth-service/internal/db/schema"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/cache"
	"github.com/rs/zerolog/log"
)

type cachedOrganizationRepository struct {
	repository                  Repository
	organizationMembershipCache cache.Cache[schema.OrganizationMemberships]
	ttl                         time.Duration
}

func NewCachedUserRepository(
	repository Repository,
	organizationMembershipCache cache.Cache[schema.OrganizationMemberships],
	ttl time.Duration,
) *cachedOrganizationRepository {
	return &cachedOrganizationRepository{
		repository:                  repository,
		organizationMembershipCache: organizationMembershipCache,
		ttl:                         ttl,
	}
}

func (r *cachedOrganizationRepository) Create(ctx context.Context, organization *schema.Organization) error {
	err := r.repository.Create(ctx, organization)
	if err != nil {
		return err
	}
	r.organizationMembershipCache.Delete(ctx, organizationMembershipsCacheKey(organization.CreatedBy))
	return nil
}

func (r *cachedOrganizationRepository) CreateOrganizationMembership(ctx context.Context, params *CreateOrganizationMembershipParams) error {
	err := r.repository.CreateOrganizationMembership(ctx, params)
	if err != nil {
		return err
	}
	r.organizationMembershipCache.Delete(ctx, organizationMembershipsCacheKey(params.UserID))
	return nil
}

func (r *cachedOrganizationRepository) ListOrganizationMembershipsByUserId(ctx context.Context, userId string) (schema.OrganizationMemberships, error) {
	memberships := r.getOrganizationMembershipsFromCache(ctx, userId)
	if memberships != nil {
		return memberships, nil
	}

	memberships, err := r.repository.ListOrganizationMembershipsByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}

	r.cacheOrganizationMemberships(ctx, memberships, userId)

	return memberships, err
}

func (r cachedOrganizationRepository) getOrganizationMembershipsFromCache(ctx context.Context, userId string) schema.OrganizationMemberships {
	key := organizationMembershipsCacheKey(userId)
	memberships, err := r.organizationMembershipCache.Get(ctx, key)
	if err == nil {
		return *memberships
	}
	if err != cache.ErrCacheMiss {
		log.Error().Err(err).Msgf("error getting %s from cache", key)
	}
	return nil
}

func (r *cachedOrganizationRepository) cacheOrganizationMemberships(ctx context.Context, memberships schema.OrganizationMemberships, userId string) {
	key := organizationMembershipsCacheKey(userId)
	item := cache.Item[schema.OrganizationMemberships]{
		Key:   key,
		Value: &memberships,
		TTL:   r.ttl,
	}
	if err := r.organizationMembershipCache.Set(ctx, &item); err != nil {
		log.Error().Err(err).Msgf("error setting %s:%v+ to cache", key, memberships)
	}
}

func organizationMembershipsCacheKey(userId string) string {
	return fmt.Sprintf("organization-memberships:user-id:%s", userId)
}

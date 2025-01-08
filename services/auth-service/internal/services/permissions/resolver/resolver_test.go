package resolver_test

import (
	"context"
	"testing"

	"github.com/FreibergVlad/url-shortener/auth-service/internal/db/schema"
	"github.com/FreibergVlad/url-shortener/auth-service/internal/services/permissions/resolver"
	"github.com/FreibergVlad/url-shortener/auth-service/internal/services/roles"
	testUtils "github.com/FreibergVlad/url-shortener/auth-service/internal/testing"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/errors"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestHasPermissionsWhenErrorGettingUser(t *testing.T) {
	t.Parallel()

	userRepo := &testUtils.MockedUserRepository{}
	organizationRepo := &testUtils.MockedOrganizationRepository{}
	resolver := resolver.New(userRepo, organizationRepo)
	ctx, userId := context.Background(), gofakeit.UUID()

	userRepo.On("GetById", ctx, userId).Return(&schema.User{}, errors.ErrResourceNotFound)

	hasPermissons, err := resolver.HasPermissions(ctx, []string{}, userId, nil)

	assert.False(t, hasPermissons)
	assert.ErrorIs(t, err, errors.ErrResourceNotFound)

	userRepo.AssertExpectations(t)
	organizationRepo.AssertExpectations(t)
}

func TestHasPermissionsWhenErrorGettingOrganizationMemberships(t *testing.T) {
	t.Parallel()

	userRepo := &testUtils.MockedUserRepository{}
	organizationRepo := &testUtils.MockedOrganizationRepository{}
	resolver := resolver.New(userRepo, organizationRepo)
	ctx, organizationId := context.Background(), gofakeit.UUID()
	user := &schema.User{ID: gofakeit.UUID(), RoleID: roles.RoleIdProvisional}

	userRepo.On("GetById", ctx, user.ID).Return(user, nil)
	organizationRepo.On("ListOrganizationMembershipsByUserId", ctx, user.ID).Return(schema.OrganizationMemberships{}, errors.ErrResourceNotFound)

	hasPermissons, err := resolver.HasPermissions(ctx, []string{"short-url:read"}, user.ID, &organizationId)

	assert.False(t, hasPermissons)
	assert.ErrorIs(t, err, errors.ErrResourceNotFound)
}

func TestHasPermissionsWithGlobalRole(t *testing.T) {
	t.Parallel()

	userRepo := &testUtils.MockedUserRepository{}
	organizationRepo := &testUtils.MockedOrganizationRepository{}
	resolver := resolver.New(userRepo, organizationRepo)
	ctx := context.Background()
	user := &schema.User{ID: gofakeit.UUID(), RoleID: roles.RoleIdProvisional}

	userRepo.On("GetById", ctx, user.ID).Return(user, nil)

	hasPermissons, err := resolver.HasPermissions(ctx, []string{"me:read"}, user.ID, nil)

	assert.True(t, hasPermissons)
	assert.NoError(t, err)

	hasPermissons, err = resolver.HasPermissions(ctx, []string{"short-url:list"}, user.ID, nil)

	assert.False(t, hasPermissons)
	assert.NoError(t, err)
}

func TestHasPermissionsWithOrganizationalRole(t *testing.T) {
	t.Parallel()

	userRepo := &testUtils.MockedUserRepository{}
	organizationRepo := &testUtils.MockedOrganizationRepository{}
	resolver := resolver.New(userRepo, organizationRepo)
	ctx, organizationId := context.Background(), gofakeit.UUID()
	user := &schema.User{ID: gofakeit.UUID(), RoleID: roles.RoleIdProvisional}
	membership := schema.OrganizationMembership{
		Organization: schema.ShortOrganization{ID: organizationId},
		RoleID:       roles.RoleIdMember,
	}

	userRepo.On("GetById", ctx, user.ID).Return(user, nil)
	organizationRepo.On("ListOrganizationMembershipsByUserId", ctx, user.ID).Return(schema.OrganizationMemberships{&membership}, nil)

	hasPermissons, err := resolver.HasPermissions(ctx, []string{"short-url:read"}, user.ID, &organizationId)

	assert.True(t, hasPermissons)
	assert.NoError(t, err)

	organizationId = gofakeit.UUID()
	hasPermissons, err = resolver.HasPermissions(ctx, []string{"short-url:read"}, user.ID, &organizationId)

	assert.False(t, hasPermissons)
	assert.NoError(t, err)
}

package resolver_test

import (
	"context"
	"testing"

	"github.com/FreibergVlad/url-shortener/auth-service/internal/db/repositories"
	"github.com/FreibergVlad/url-shortener/auth-service/internal/db/schema"
	"github.com/FreibergVlad/url-shortener/auth-service/internal/services/permissions/resolver"
	"github.com/FreibergVlad/url-shortener/auth-service/internal/services/roles"
	testUtils "github.com/FreibergVlad/url-shortener/auth-service/internal/testing"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
)

func TestHasPermissionsWhenErrorGettingUser(t *testing.T) {
	t.Parallel()

	userRepo := &testUtils.MockedUserRepository{}
	organizationRepo := &testUtils.MockedOrganizationRepository{}
	resolver := resolver.New(userRepo, organizationRepo)
	ctx, userID := context.Background(), gofakeit.UUID()

	userRepo.On("GetByID", ctx, userID).Return(&schema.User{}, repositories.ErrNotFound)

	hasPermissons, err := resolver.HasPermissions(ctx, []string{}, userID, nil)

	require.False(t, hasPermissons)
	require.ErrorIs(t, err, repositories.ErrNotFound)

	userRepo.AssertExpectations(t)
	organizationRepo.AssertExpectations(t)
}

func TestHasPermissionsWhenErrorGettingOrganizationMemberships(t *testing.T) {
	t.Parallel()

	userRepo := &testUtils.MockedUserRepository{}
	organizationRepo := &testUtils.MockedOrganizationRepository{}
	resolver := resolver.New(userRepo, organizationRepo)
	ctx, organizationID := context.Background(), gofakeit.UUID()
	user := &schema.User{ID: gofakeit.UUID(), RoleID: roles.RoleIDProvisional}

	userRepo.On("GetByID", ctx, user.ID).Return(user, nil)
	organizationRepo.
		On("ListOrganizationMembershipsByUserID", ctx, user.ID).
		Return(schema.OrganizationMemberships{}, repositories.ErrNotFound)

	hasPermissons, err := resolver.HasPermissions(ctx, []string{"short-url:read"}, user.ID, &organizationID)

	require.False(t, hasPermissons)
	require.ErrorIs(t, err, repositories.ErrNotFound)
}

func TestHasPermissionsWithGlobalRole(t *testing.T) {
	t.Parallel()

	userRepo := &testUtils.MockedUserRepository{}
	organizationRepo := &testUtils.MockedOrganizationRepository{}
	resolver := resolver.New(userRepo, organizationRepo)
	ctx := context.Background()
	user := &schema.User{ID: gofakeit.UUID(), RoleID: roles.RoleIDProvisional}

	userRepo.On("GetByID", ctx, user.ID).Return(user, nil)

	hasPermissons, err := resolver.HasPermissions(ctx, []string{"me:read"}, user.ID, nil)

	require.True(t, hasPermissons)
	require.NoError(t, err)

	hasPermissons, err = resolver.HasPermissions(ctx, []string{"short-url:list"}, user.ID, nil)

	require.False(t, hasPermissons)
	require.NoError(t, err)
}

func TestHasPermissionsWithOrganizationalRole(t *testing.T) {
	t.Parallel()

	userRepo := &testUtils.MockedUserRepository{}
	organizationRepo := &testUtils.MockedOrganizationRepository{}
	resolver := resolver.New(userRepo, organizationRepo)
	ctx, organizationID := context.Background(), gofakeit.UUID()
	user := &schema.User{ID: gofakeit.UUID(), RoleID: roles.RoleIDProvisional}
	membership := schema.OrganizationMembership{
		Organization: schema.ShortOrganization{ID: organizationID},
		RoleID:       roles.RoleIDMember,
	}

	userRepo.On("GetByID", ctx, user.ID).Return(user, nil)
	organizationRepo.
		On("ListOrganizationMembershipsByUserID", ctx, user.ID).
		Return(schema.OrganizationMemberships{&membership}, nil)

	hasPermissons, err := resolver.HasPermissions(ctx, []string{"short-url:read"}, user.ID, &organizationID)

	require.True(t, hasPermissons)
	require.NoError(t, err)

	organizationID = gofakeit.UUID()
	hasPermissons, err = resolver.HasPermissions(ctx, []string{"short-url:read"}, user.ID, &organizationID)

	require.False(t, hasPermissons)
	require.NoError(t, err)
}

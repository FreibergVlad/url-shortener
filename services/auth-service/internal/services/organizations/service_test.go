package organizations_test

import (
	"context"
	"testing"
	"time"

	"github.com/FreibergVlad/url-shortener/auth-service/internal/db/schema"
	"github.com/FreibergVlad/url-shortener/auth-service/internal/services/organizations"
	"github.com/FreibergVlad/url-shortener/auth-service/internal/services/roles"
	testUtils "github.com/FreibergVlad/url-shortener/auth-service/internal/testing"
	organizationServiceMessages "github.com/FreibergVlad/url-shortener/proto/pkg/organizations/messages/v1"
	grpcUtils "github.com/FreibergVlad/url-shortener/shared/go/pkg/api/grpc/utils"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/clock"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestCreateOrganization(t *testing.T) {
	t.Parallel()

	organizationRepo := &testUtils.MockedOrganizationRepository{}
	clock := clock.NewFixedClock(time.Now())
	userID := gofakeit.UUID()
	organizationService := organizations.New(organizationRepo, clock)
	ctx := grpcUtils.IncomingContextWithUserID(context.Background(), userID)
	request := testUtils.CreateTestOrganizationRequest()

	organizationRepo.On("Create", ctx, mock.Anything).Return(nil)

	response, err := organizationService.CreateOrganization(ctx, request)

	require.NoError(t, err)

	assert.Equal(t, request.Name, response.Organization.Name)
	assert.Equal(t, request.Slug, response.Organization.Slug)
	assert.Equal(t, timestamppb.New(clock.Now()), response.Organization.CreatedAt)
	assert.Equal(t, userID, response.Organization.CreatedBy)

	organizationRepo.AssertExpectations(t)
}

func TestCreateOrganizationWhenDatabaseError(t *testing.T) {
	t.Parallel()

	organizationRepo := &testUtils.MockedOrganizationRepository{}
	organizationService := organizations.New(organizationRepo, clock.NewFixedClock(time.Now()))
	ctx := grpcUtils.IncomingContextWithUserID(context.Background(), gofakeit.UUID())
	request := testUtils.CreateTestOrganizationRequest()
	wantErr := gofakeit.ErrorDatabase()

	organizationRepo.On("Create", ctx, mock.Anything).Return(wantErr)

	response, err := organizationService.CreateOrganization(ctx, request)

	require.ErrorIs(t, err, wantErr)

	assert.Nil(t, response)

	organizationRepo.AssertExpectations(t)
}

func TestListMeOrganizationMemberships(t *testing.T) {
	t.Parallel()

	organizationRepo := &testUtils.MockedOrganizationRepository{}
	clock := clock.NewFixedClock(time.Now())
	userID := gofakeit.UUID()
	organizationService := organizations.New(organizationRepo, clock)
	ctx := grpcUtils.IncomingContextWithUserID(context.Background(), userID)
	request := &organizationServiceMessages.ListMeOrganizationMembershipsRequest{}
	membership := schema.OrganizationMembership{
		Organization: schema.ShortOrganization{
			ID:   gofakeit.UUID(),
			Slug: gofakeit.Company(),
		},
		RoleID:    roles.RoleIDMember,
		CreatedAt: clock.Now(),
	}

	organizationRepo.
		On("ListOrganizationMembershipsByUserID", ctx, userID).
		Return(schema.OrganizationMemberships{&membership}, nil)

	response, err := organizationService.ListMeOrganizationMemberships(ctx, request)

	require.NoError(t, err)

	assert.Len(t, response.Data, 1)

	responseMembership := response.Data[0]
	assert.Equal(t, membership.Organization.ID, responseMembership.Organization.Id)
	assert.Equal(t, membership.Organization.Slug, responseMembership.Organization.Slug)
	assert.Equal(t, timestamppb.New(clock.Now()), responseMembership.CreatedAt)
	assert.Equal(t, membership.RoleID, responseMembership.Role.Id)
	assert.Equal(t, roles.RoleMember.Name, responseMembership.Role.Name)
	assert.Equal(t, roles.RoleMember.Description, responseMembership.Role.Description)

	organizationRepo.AssertExpectations(t)
}

func TestListMeOrganizationMembershipsWhenDatabaseError(t *testing.T) {
	t.Parallel()

	organizationRepo := &testUtils.MockedOrganizationRepository{}
	clock := clock.NewFixedClock(time.Now())
	userID := gofakeit.UUID()
	organizationService := organizations.New(organizationRepo, clock)
	ctx := grpcUtils.IncomingContextWithUserID(context.Background(), userID)
	request := &organizationServiceMessages.ListMeOrganizationMembershipsRequest{}
	wantErr := gofakeit.ErrorDatabase()

	organizationRepo.
		On("ListOrganizationMembershipsByUserID", ctx, userID).
		Return(schema.OrganizationMemberships{}, wantErr)

	response, err := organizationService.ListMeOrganizationMemberships(ctx, request)

	require.ErrorIs(t, err, wantErr)

	assert.Nil(t, response)

	organizationRepo.AssertExpectations(t)
}

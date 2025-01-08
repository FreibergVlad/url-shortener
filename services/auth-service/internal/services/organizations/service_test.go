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
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/errors"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestCreateOrganization(t *testing.T) {
	t.Parallel()

	organizationRepo := &testUtils.MockedOrganizationRepository{}
	clock := clock.NewFixedClock(time.Now())
	userId := gofakeit.UUID()
	organizationService := organizations.New(organizationRepo, clock)
	ctx, request := grpcUtils.IncomingContextWithUserID(context.Background(), userId), testUtils.CreateTestOrganizationRequest()

	organizationRepo.On("Create", ctx, mock.Anything).Return(nil)

	response, err := organizationService.CreateOrganization(ctx, request)

	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, request.Name, response.Organization.Name)
	assert.Equal(t, request.Slug, response.Organization.Slug)
	assert.Equal(t, timestamppb.New(clock.Now()), response.Organization.CreatedAt)
	assert.Equal(t, userId, response.Organization.CreatedBy)

	organizationRepo.AssertExpectations(t)
}

func TestCreateOrganizationWhenDatabaseError(t *testing.T) {
	t.Parallel()

	organizationRepo := &testUtils.MockedOrganizationRepository{}
	organizationService := organizations.New(organizationRepo, clock.NewFixedClock(time.Now()))
	ctx, request := grpcUtils.IncomingContextWithUserID(context.Background(), gofakeit.UUID()), testUtils.CreateTestOrganizationRequest()

	organizationRepo.On("Create", ctx, mock.Anything).Return(errors.ErrDuplicateResource)

	response, err := organizationService.CreateOrganization(ctx, request)

	assert.Nil(t, response)
	assert.ErrorIs(t, err, errors.ErrDuplicateResource)

	organizationRepo.AssertExpectations(t)
}

func TestListMeOrganizationMemberships(t *testing.T) {
	t.Parallel()

	organizationRepo := &testUtils.MockedOrganizationRepository{}
	clock := clock.NewFixedClock(time.Now())
	userId := gofakeit.UUID()
	organizationService := organizations.New(organizationRepo, clock)
	ctx := grpcUtils.IncomingContextWithUserID(context.Background(), userId)
	request := &organizationServiceMessages.ListMeOrganizationMembershipsRequest{}
	membership := schema.OrganizationMembership{
		Organization: schema.ShortOrganization{
			ID:   gofakeit.UUID(),
			Slug: gofakeit.Company(),
		},
		RoleID:    roles.RoleIdMember,
		CreatedAt: clock.Now(),
	}

	organizationRepo.
		On("ListOrganizationMembershipsByUserId", ctx, userId).
		Return(schema.OrganizationMemberships{&membership}, nil)

	response, err := organizationService.ListMeOrganizationMemberships(ctx, request)

	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, 1, len(response.Data))

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
	userId := gofakeit.UUID()
	organizationService := organizations.New(organizationRepo, clock)
	ctx := grpcUtils.IncomingContextWithUserID(context.Background(), userId)
	request := &organizationServiceMessages.ListMeOrganizationMembershipsRequest{}

	organizationRepo.
		On("ListOrganizationMembershipsByUserId", ctx, userId).
		Return(schema.OrganizationMemberships{}, errors.ErrResourceNotFound)

	response, err := organizationService.ListMeOrganizationMemberships(ctx, request)

	assert.Nil(t, response)
	assert.ErrorIs(t, err, errors.ErrResourceNotFound)

	organizationRepo.AssertExpectations(t)
}

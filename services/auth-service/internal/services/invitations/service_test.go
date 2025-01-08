package invitations_test

import (
	"context"
	"testing"
	"time"

	"github.com/FreibergVlad/url-shortener/auth-service/internal/db/repositories/organizations"
	"github.com/FreibergVlad/url-shortener/auth-service/internal/db/schema"
	"github.com/FreibergVlad/url-shortener/auth-service/internal/services/invitations"
	"github.com/FreibergVlad/url-shortener/auth-service/internal/services/roles"
	testUtils "github.com/FreibergVlad/url-shortener/auth-service/internal/testing"
	invitationServiceMessages "github.com/FreibergVlad/url-shortener/proto/pkg/invitations/messages/v1"
	grpcUtils "github.com/FreibergVlad/url-shortener/shared/go/pkg/api/grpc/utils"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/clock"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/errors"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type mockedInvitationRepository struct {
	mock.Mock
}

func (r *mockedInvitationRepository) Create(ctx context.Context, invitation *schema.Invitation) error {
	args := r.Called(ctx, invitation)
	return args.Error(0)
}

func (r *mockedInvitationRepository) GetById(ctx context.Context, id string) (*schema.Invitation, error) {
	args := r.Called(ctx, id)
	return args.Get(0).(*schema.Invitation), args.Error(1)
}

func (r *mockedInvitationRepository) UpdateStatusById(ctx context.Context, id, status string) error {
	args := r.Called(ctx, id, status)
	return args.Error(0)
}

func TestCreateInvitation(t *testing.T) {
	t.Parallel()

	userRepo, organizationRepo, invitationRepo := testUtils.MockedUserRepository{}, testUtils.MockedOrganizationRepository{}, mockedInvitationRepository{}
	clock := clock.NewFixedClock(time.Now())
	invitationService := invitations.New(&userRepo, &invitationRepo, &organizationRepo, clock)
	userId := gofakeit.UUID()
	ctx := grpcUtils.IncomingContextWithUserID(context.Background(), userId)
	request := invitationServiceMessages.CreateInvitationRequest{
		Email:          gofakeit.Email(),
		OrganizationId: gofakeit.UUID(),
		RoleId:         roles.RoleIdMember,
	}

	userRepo.On("GetByEmail", ctx, request.Email).Return(&schema.User{}, errors.ErrResourceNotFound)
	invitationRepo.On("Create", ctx, mock.Anything).Return(nil)

	response, err := invitationService.CreateInvitation(ctx, &request)

	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, request.OrganizationId, response.Invitation.OrganizationId)
	assert.Equal(t, invitationServiceMessages.InvitationStatus_INVITATION_STATUS_ACTIVE, response.Invitation.Status)
	assert.Equal(t, request.Email, response.Invitation.Email)
	assert.Equal(t, request.RoleId, response.Invitation.RoleId)
	assert.Equal(t, timestamppb.New(clock.Now()), response.Invitation.CreatedAt)
	assert.Equal(t, timestamppb.New(clock.Now().Add(invitations.InviteLifetime)), response.Invitation.ExpiresAt)
	assert.Equal(t, userId, response.Invitation.CreatedBy)

	userRepo.AssertExpectations(t)
	invitationRepo.AssertExpectations(t)
	organizationRepo.AssertExpectations(t)
}

func TestCreateInvitationWhenErrorGettingInvitedUser(t *testing.T) {
	t.Parallel()

	userRepo, organizationRepo, invitationRepo := testUtils.MockedUserRepository{}, testUtils.MockedOrganizationRepository{}, mockedInvitationRepository{}
	clock := clock.NewFixedClock(time.Now())
	invitationService := invitations.New(&userRepo, &invitationRepo, &organizationRepo, clock)
	ctx := grpcUtils.IncomingContextWithUserID(context.Background(), gofakeit.UUID())
	request := invitationServiceMessages.CreateInvitationRequest{
		Email:          gofakeit.Email(),
		OrganizationId: gofakeit.UUID(),
		RoleId:         roles.RoleIdMember,
	}

	userRepo.On("GetByEmail", ctx, request.Email).Return(&schema.User{}, errors.NewInternalError("err"))

	response, err := invitationService.CreateInvitation(ctx, &request)

	assert.ErrorIs(t, err, errors.NewInternalError("err"))
	assert.Nil(t, response)

	userRepo.AssertExpectations(t)
	invitationRepo.AssertExpectations(t)
	organizationRepo.AssertExpectations(t)
}

func TestCreateInvitationWhenErrorGettingInvitedUserMemberships(t *testing.T) {
	t.Parallel()

	userRepo, organizationRepo, invitationRepo := testUtils.MockedUserRepository{}, testUtils.MockedOrganizationRepository{}, mockedInvitationRepository{}
	clock := clock.NewFixedClock(time.Now())
	invitationService := invitations.New(&userRepo, &invitationRepo, &organizationRepo, clock)
	ctx := grpcUtils.IncomingContextWithUserID(context.Background(), gofakeit.UUID())
	request := invitationServiceMessages.CreateInvitationRequest{
		Email:          gofakeit.Email(),
		OrganizationId: gofakeit.UUID(),
		RoleId:         roles.RoleIdMember,
	}
	invitedUser := schema.User{
		ID:    gofakeit.UUID(),
		Email: request.Email,
	}

	userRepo.On("GetByEmail", ctx, request.Email).Return(&invitedUser, nil)
	organizationRepo.
		On("ListOrganizationMembershipsByUserId", ctx, invitedUser.ID).
		Return(schema.OrganizationMemberships{}, errors.ErrResourceNotFound)

	response, err := invitationService.CreateInvitation(ctx, &request)

	assert.ErrorIs(t, err, errors.ErrResourceNotFound)
	assert.Nil(t, response)

	userRepo.AssertExpectations(t)
	invitationRepo.AssertExpectations(t)
	organizationRepo.AssertExpectations(t)
}

func TestCreateInvitationWhenUserAlreadyBelongsToOrganization(t *testing.T) {
	t.Parallel()

	userRepo, organizationRepo, invitationRepo := testUtils.MockedUserRepository{}, testUtils.MockedOrganizationRepository{}, mockedInvitationRepository{}
	clock := clock.NewFixedClock(time.Now())
	invitationService := invitations.New(&userRepo, &invitationRepo, &organizationRepo, clock)
	ctx := grpcUtils.IncomingContextWithUserID(context.Background(), gofakeit.UUID())
	request := invitationServiceMessages.CreateInvitationRequest{
		Email:          gofakeit.Email(),
		OrganizationId: gofakeit.UUID(),
		RoleId:         roles.RoleIdMember,
	}
	invitedUser := schema.User{
		ID:    gofakeit.UUID(),
		Email: request.Email,
	}
	membership := schema.OrganizationMembership{
		Organization: schema.ShortOrganization{
			Slug: gofakeit.Company(),
			ID:   request.OrganizationId,
		},
		RoleID: roles.RoleIdMember,
	}

	userRepo.On("GetByEmail", ctx, request.Email).Return(&invitedUser, nil)
	organizationRepo.
		On("ListOrganizationMembershipsByUserId", ctx, invitedUser.ID).
		Return(schema.OrganizationMemberships{&membership}, nil)

	response, err := invitationService.CreateInvitation(ctx, &request)

	assert.ErrorContains(t, err, "already belongs to organization")
	assert.Nil(t, response)

	userRepo.AssertExpectations(t)
	invitationRepo.AssertExpectations(t)
	organizationRepo.AssertExpectations(t)
}

func TestCreateInvitationWhenErrorCreatingInvitation(t *testing.T) {
	t.Parallel()

	userRepo, organizationRepo, invitationRepo := testUtils.MockedUserRepository{}, testUtils.MockedOrganizationRepository{}, mockedInvitationRepository{}
	clock := clock.NewFixedClock(time.Now())
	invitationService := invitations.New(&userRepo, &invitationRepo, &organizationRepo, clock)
	ctx := grpcUtils.IncomingContextWithUserID(context.Background(), gofakeit.UUID())
	request := invitationServiceMessages.CreateInvitationRequest{
		Email:          gofakeit.Email(),
		OrganizationId: gofakeit.UUID(),
		RoleId:         roles.RoleIdMember,
	}

	userRepo.On("GetByEmail", ctx, request.Email).Return(&schema.User{}, errors.ErrResourceNotFound)
	invitationRepo.On("Create", ctx, mock.Anything).Return(errors.NewInternalError("err"))

	response, err := invitationService.CreateInvitation(ctx, &request)

	assert.ErrorIs(t, err, errors.NewInternalError("err"))
	assert.Nil(t, response)

	userRepo.AssertExpectations(t)
	invitationRepo.AssertExpectations(t)
	organizationRepo.AssertExpectations(t)
}

func TestAcceptInvitation(t *testing.T) {
	t.Parallel()

	userRepo, organizationRepo, invitationRepo := testUtils.MockedUserRepository{}, testUtils.MockedOrganizationRepository{}, mockedInvitationRepository{}
	clock := clock.NewFixedClock(time.Now())
	invitationService := invitations.New(&userRepo, &invitationRepo, &organizationRepo, clock)
	user, ctx := fakeUserWithContext()
	invitation := fakeActiveInvitation(user.Email)
	request := invitationServiceMessages.AcceptInvitationRequest{Id: invitation.ID}

	userRepo.On("GetById", ctx, user.ID).Return(user, nil)
	invitationRepo.On("GetById", ctx, invitation.ID).Return(invitation, nil)
	organizationRepo.
		On(
			"CreateOrganizationMembership",
			ctx,
			&organizations.CreateOrganizationMembershipParams{
				OrganizationID: invitation.OrganizationID,
				UserID:         user.ID,
				RoleID:         invitation.RoleID,
				CreatedAt:      clock.Now(),
			},
		).
		Return(nil)
	invitationRepo.
		On(
			"UpdateStatusById",
			ctx,
			invitation.ID,
			invitationServiceMessages.InvitationStatus_INVITATION_STATUS_REDEEMED.String(),
		).
		Return(nil)

	response, err := invitationService.AcceptInvitation(ctx, &request)

	assert.NoError(t, err)
	assert.NotNil(t, response)

	userRepo.AssertExpectations(t)
	invitationRepo.AssertExpectations(t)
	organizationRepo.AssertExpectations(t)
}

func TestAcceptInvitationWhenErrorGettingInvitedUser(t *testing.T) {
	t.Parallel()

	userRepo, organizationRepo, invitationRepo := testUtils.MockedUserRepository{}, testUtils.MockedOrganizationRepository{}, mockedInvitationRepository{}
	clock := clock.NewFixedClock(time.Now())
	invitationService := invitations.New(&userRepo, &invitationRepo, &organizationRepo, clock)
	user, ctx := fakeUserWithContext()
	request := invitationServiceMessages.AcceptInvitationRequest{Id: gofakeit.UUID()}

	userRepo.On("GetById", ctx, user.ID).Return(&schema.User{}, errors.ErrResourceNotFound)

	response, err := invitationService.AcceptInvitation(ctx, &request)

	assert.Nil(t, response)
	assert.ErrorIs(t, err, errors.ErrResourceNotFound)

	userRepo.AssertExpectations(t)
	invitationRepo.AssertExpectations(t)
	organizationRepo.AssertExpectations(t)
}

func TestAcceptInvitationWhenInvitationValidationFailed(t *testing.T) {
	t.Parallel()

	userRepo, organizationRepo, invitationRepo := testUtils.MockedUserRepository{}, testUtils.MockedOrganizationRepository{}, mockedInvitationRepository{}
	clock := clock.NewFixedClock(time.Now())
	invitationService := invitations.New(&userRepo, &invitationRepo, &organizationRepo, clock)
	user, ctx := fakeUserWithContext()

	nonExistentInvitationId := gofakeit.UUID()
	userRepo.On("GetById", ctx, user.ID).Return(user, nil)
	invitationRepo.On("GetById", ctx, nonExistentInvitationId).Return(&schema.Invitation{}, errors.ErrResourceNotFound)

	request := invitationServiceMessages.AcceptInvitationRequest{Id: nonExistentInvitationId}
	response, err := invitationService.AcceptInvitation(ctx, &request)

	assert.Nil(t, response)
	assert.ErrorIs(t, err, errors.ErrResourceNotFound)

	invitationWithIncorrectEmail := fakeActiveInvitation(gofakeit.Email())
	userRepo.On("GetById", ctx, user.ID).Return(user, nil)
	invitationRepo.On("GetById", ctx, invitationWithIncorrectEmail.ID).Return(invitationWithIncorrectEmail, nil)

	request = invitationServiceMessages.AcceptInvitationRequest{Id: invitationWithIncorrectEmail.ID}
	response, err = invitationService.AcceptInvitation(ctx, &request)

	assert.Nil(t, response)
	assert.ErrorContains(t, err, "permission denied")

	expiredInvitation := schema.Invitation{
		ID:             gofakeit.UUID(),
		Status:         invitationServiceMessages.InvitationStatus_INVITATION_STATUS_EXPIRED.String(),
		OrganizationID: gofakeit.UUID(),
		Email:          user.Email,
		RoleID:         roles.RoleIdMember,
	}
	userRepo.On("GetById", ctx, user.ID).Return(user, nil)
	invitationRepo.On("GetById", ctx, expiredInvitation.ID).Return(&expiredInvitation, nil)

	request = invitationServiceMessages.AcceptInvitationRequest{Id: expiredInvitation.ID}
	response, err = invitationService.AcceptInvitation(ctx, &request)

	assert.Nil(t, response)
	assert.ErrorContains(t, err, "invite expired")

	userRepo.AssertExpectations(t)
	invitationRepo.AssertExpectations(t)
	organizationRepo.AssertExpectations(t)
}

func TestAcceptInvitationWhenErrorCreatingUserMembership(t *testing.T) {
	t.Parallel()

	userRepo, organizationRepo, invitationRepo := testUtils.MockedUserRepository{}, testUtils.MockedOrganizationRepository{}, mockedInvitationRepository{}
	clock := clock.NewFixedClock(time.Now())
	invitationService := invitations.New(&userRepo, &invitationRepo, &organizationRepo, clock)
	user, ctx := fakeUserWithContext()
	invitation := fakeActiveInvitation(user.Email)
	request := invitationServiceMessages.AcceptInvitationRequest{Id: invitation.ID}

	userRepo.On("GetById", ctx, user.ID).Return(user, nil)
	invitationRepo.On("GetById", ctx, invitation.ID).Return(invitation, nil)
	organizationRepo.
		On(
			"CreateOrganizationMembership",
			ctx,
			&organizations.CreateOrganizationMembershipParams{
				OrganizationID: invitation.OrganizationID,
				UserID:         user.ID,
				RoleID:         invitation.RoleID,
				CreatedAt:      clock.Now(),
			},
		).
		Return(errors.NewInternalError("err"))

	response, err := invitationService.AcceptInvitation(ctx, &request)

	assert.Nil(t, response)
	assert.ErrorContains(t, err, "err")

	userRepo.AssertExpectations(t)
	invitationRepo.AssertExpectations(t)
	organizationRepo.AssertExpectations(t)
}

func TestAcceptInvitationWhenErrorUpdatingInvitationStatus(t *testing.T) {
	t.Parallel()

	userRepo, organizationRepo, invitationRepo := testUtils.MockedUserRepository{}, testUtils.MockedOrganizationRepository{}, mockedInvitationRepository{}
	clock := clock.NewFixedClock(time.Now())
	invitationService := invitations.New(&userRepo, &invitationRepo, &organizationRepo, clock)
	user, ctx := fakeUserWithContext()
	invitation := fakeActiveInvitation(user.Email)
	request := invitationServiceMessages.AcceptInvitationRequest{Id: invitation.ID}

	userRepo.On("GetById", ctx, user.ID).Return(user, nil)
	invitationRepo.On("GetById", ctx, invitation.ID).Return(invitation, nil)
	organizationRepo.
		On(
			"CreateOrganizationMembership",
			ctx,
			&organizations.CreateOrganizationMembershipParams{
				OrganizationID: invitation.OrganizationID,
				UserID:         user.ID,
				RoleID:         invitation.RoleID,
				CreatedAt:      clock.Now(),
			},
		).
		Return(nil)
	invitationRepo.
		On(
			"UpdateStatusById",
			ctx,
			invitation.ID,
			invitationServiceMessages.InvitationStatus_INVITATION_STATUS_REDEEMED.String(),
		).
		Return(errors.NewInternalError("err"))

	response, err := invitationService.AcceptInvitation(ctx, &request)

	assert.Nil(t, response)
	assert.ErrorContains(t, err, "err")

	userRepo.AssertExpectations(t)
	invitationRepo.AssertExpectations(t)
	organizationRepo.AssertExpectations(t)
}

func fakeUserWithContext() (*schema.User, context.Context) {
	user := schema.User{ID: gofakeit.UUID(), Email: gofakeit.Email()}
	ctx := grpcUtils.IncomingContextWithUserID(context.Background(), user.ID)
	return &user, ctx
}

func fakeActiveInvitation(email string) *schema.Invitation {
	return &schema.Invitation{
		ID:             gofakeit.UUID(),
		Status:         invitationServiceMessages.InvitationStatus_INVITATION_STATUS_ACTIVE.String(),
		OrganizationID: gofakeit.UUID(),
		Email:          email,
		RoleID:         roles.RoleIdMember,
	}
}

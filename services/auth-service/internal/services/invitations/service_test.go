package invitations_test

import (
	"context"
	"testing"
	"time"

	"github.com/FreibergVlad/url-shortener/auth-service/internal/db/repositories"
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
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type mockedInvitationRepository struct {
	mock.Mock
}

func (r *mockedInvitationRepository) Create(ctx context.Context, invitation *schema.Invitation) error {
	args := r.Called(ctx, invitation)
	return args.Error(0)
}

func (r *mockedInvitationRepository) GetByID(ctx context.Context, id string) (*schema.Invitation, error) {
	args := r.Called(ctx, id)
	return args.Get(0).(*schema.Invitation), args.Error(1)
}

func (r *mockedInvitationRepository) UpdateStatusByID(ctx context.Context, id, status string) error {
	args := r.Called(ctx, id, status)
	return args.Error(0)
}

func TestCreateInvitation(t *testing.T) {
	t.Parallel()

	userRepo := testUtils.MockedUserRepository{}
	organizationRepo := testUtils.MockedOrganizationRepository{}
	invitationRepo := mockedInvitationRepository{}

	clock := clock.NewFixedClock(time.Now())
	invitationService := invitations.New(&userRepo, &invitationRepo, &organizationRepo, clock)
	userID := gofakeit.UUID()
	ctx := grpcUtils.IncomingContextWithUserID(context.Background(), userID)
	request := invitationServiceMessages.CreateInvitationRequest{
		Email:          gofakeit.Email(),
		OrganizationId: gofakeit.UUID(),
		RoleId:         roles.RoleIDMember,
	}

	userRepo.On("GetByEmail", ctx, request.Email).Return(&schema.User{}, repositories.ErrNotFound)
	invitationRepo.On("Create", ctx, mock.Anything).Return(nil)

	response, err := invitationService.CreateInvitation(ctx, &request)

	require.NoError(t, err)

	assert.Equal(t, request.OrganizationId, response.Invitation.OrganizationId)
	assert.Equal(t, invitationServiceMessages.InvitationStatus_INVITATION_STATUS_ACTIVE, response.Invitation.Status)
	assert.Equal(t, request.Email, response.Invitation.Email)
	assert.Equal(t, request.RoleId, response.Invitation.RoleId)
	assert.Equal(t, timestamppb.New(clock.Now()), response.Invitation.CreatedAt)
	assert.Equal(t, timestamppb.New(clock.Now().Add(invitations.InviteLifetime)), response.Invitation.ExpiresAt)
	assert.Equal(t, userID, response.Invitation.CreatedBy)

	userRepo.AssertExpectations(t)
	invitationRepo.AssertExpectations(t)
	organizationRepo.AssertExpectations(t)
}

func TestCreateInvitationWhenErrorGettingInvitedUser(t *testing.T) {
	t.Parallel()

	userRepo := testUtils.MockedUserRepository{}
	organizationRepo := testUtils.MockedOrganizationRepository{}
	invitationRepo := mockedInvitationRepository{}

	clock := clock.NewFixedClock(time.Now())
	invitationService := invitations.New(&userRepo, &invitationRepo, &organizationRepo, clock)
	ctx := grpcUtils.IncomingContextWithUserID(context.Background(), gofakeit.UUID())
	request := invitationServiceMessages.CreateInvitationRequest{
		Email:          gofakeit.Email(),
		OrganizationId: gofakeit.UUID(),
		RoleId:         roles.RoleIDMember,
	}
	wantErr := gofakeit.ErrorDatabase()

	userRepo.On("GetByEmail", ctx, request.Email).Return(&schema.User{}, wantErr)

	response, err := invitationService.CreateInvitation(ctx, &request)

	require.ErrorIs(t, err, wantErr)

	assert.Nil(t, response)

	userRepo.AssertExpectations(t)
	invitationRepo.AssertExpectations(t)
	organizationRepo.AssertExpectations(t)
}

func TestCreateInvitationWhenErrorGettingInvitedUserMemberships(t *testing.T) {
	t.Parallel()

	userRepo := testUtils.MockedUserRepository{}
	organizationRepo := testUtils.MockedOrganizationRepository{}
	invitationRepo := mockedInvitationRepository{}

	clock := clock.NewFixedClock(time.Now())
	invitationService := invitations.New(&userRepo, &invitationRepo, &organizationRepo, clock)
	ctx := grpcUtils.IncomingContextWithUserID(context.Background(), gofakeit.UUID())
	request := invitationServiceMessages.CreateInvitationRequest{
		Email:          gofakeit.Email(),
		OrganizationId: gofakeit.UUID(),
		RoleId:         roles.RoleIDMember,
	}
	invitedUser := schema.User{
		ID:    gofakeit.UUID(),
		Email: request.Email,
	}
	wantErr := gofakeit.ErrorDatabase()

	userRepo.On("GetByEmail", ctx, request.Email).Return(&invitedUser, nil)
	organizationRepo.
		On("ListOrganizationMembershipsByUserID", ctx, invitedUser.ID).
		Return(schema.OrganizationMemberships{}, wantErr)

	response, err := invitationService.CreateInvitation(ctx, &request)

	require.ErrorIs(t, err, wantErr)

	assert.Nil(t, response)

	userRepo.AssertExpectations(t)
	invitationRepo.AssertExpectations(t)
	organizationRepo.AssertExpectations(t)
}

func TestCreateInvitationWhenUserAlreadyBelongsToOrganization(t *testing.T) {
	t.Parallel()

	userRepo := testUtils.MockedUserRepository{}
	organizationRepo := testUtils.MockedOrganizationRepository{}
	invitationRepo := mockedInvitationRepository{}

	clock := clock.NewFixedClock(time.Now())
	invitationService := invitations.New(&userRepo, &invitationRepo, &organizationRepo, clock)
	ctx := grpcUtils.IncomingContextWithUserID(context.Background(), gofakeit.UUID())
	request := invitationServiceMessages.CreateInvitationRequest{
		Email:          gofakeit.Email(),
		OrganizationId: gofakeit.UUID(),
		RoleId:         roles.RoleIDMember,
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
		RoleID: roles.RoleIDMember,
	}

	userRepo.On("GetByEmail", ctx, request.Email).Return(&invitedUser, nil)
	organizationRepo.
		On("ListOrganizationMembershipsByUserID", ctx, invitedUser.ID).
		Return(schema.OrganizationMemberships{&membership}, nil)

	response, err := invitationService.CreateInvitation(ctx, &request)

	require.ErrorIs(t, err, errors.ErrOrganizationMembershipAlreadyExists)

	assert.Nil(t, response)

	userRepo.AssertExpectations(t)
	invitationRepo.AssertExpectations(t)
	organizationRepo.AssertExpectations(t)
}

func TestCreateInvitationWhenErrorCreatingInvitation(t *testing.T) {
	t.Parallel()

	unknownDBErr := gofakeit.ErrorDatabase()

	tests := []struct {
		name    string
		gotErr  error
		wantErr error
	}{
		{name: "unknown error creating invitation", gotErr: unknownDBErr, wantErr: unknownDBErr},
		{
			name:    "invitation already exists",
			gotErr:  repositories.ErrAlreadyExists,
			wantErr: errors.ErrOrganizationInvitationAlreadyExists,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			userRepo := testUtils.MockedUserRepository{}
			organizationRepo := testUtils.MockedOrganizationRepository{}
			invitationRepo := mockedInvitationRepository{}

			clock := clock.NewFixedClock(time.Now())
			invitationService := invitations.New(&userRepo, &invitationRepo, &organizationRepo, clock)
			ctx := grpcUtils.IncomingContextWithUserID(context.Background(), gofakeit.UUID())
			request := invitationServiceMessages.CreateInvitationRequest{
				Email:          gofakeit.Email(),
				OrganizationId: gofakeit.UUID(),
				RoleId:         roles.RoleIDMember,
			}

			userRepo.On("GetByEmail", ctx, request.Email).Return(&schema.User{}, repositories.ErrNotFound)
			invitationRepo.On("Create", ctx, mock.Anything).Return(test.gotErr)

			response, err := invitationService.CreateInvitation(ctx, &request)

			require.ErrorIs(t, err, test.wantErr)

			assert.Nil(t, response)

			userRepo.AssertExpectations(t)
			invitationRepo.AssertExpectations(t)
			organizationRepo.AssertExpectations(t)
		})
	}
}

func TestAcceptInvitation(t *testing.T) {
	t.Parallel()

	userRepo := testUtils.MockedUserRepository{}
	organizationRepo := testUtils.MockedOrganizationRepository{}
	invitationRepo := mockedInvitationRepository{}

	clock := clock.NewFixedClock(time.Now())
	invitationService := invitations.New(&userRepo, &invitationRepo, &organizationRepo, clock)
	user, ctx := fakeUserWithContext()
	invitation := fakeActiveInvitation(user.Email)
	request := invitationServiceMessages.AcceptInvitationRequest{Id: invitation.ID}

	userRepo.On("GetByID", ctx, user.ID).Return(user, nil)
	invitationRepo.On("GetByID", ctx, invitation.ID).Return(invitation, nil)
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
			"UpdateStatusByID",
			ctx,
			invitation.ID,
			invitationServiceMessages.InvitationStatus_INVITATION_STATUS_REDEEMED.String(),
		).
		Return(nil)

	response, err := invitationService.AcceptInvitation(ctx, &request)

	require.NoError(t, err)

	assert.NotNil(t, response)

	userRepo.AssertExpectations(t)
	invitationRepo.AssertExpectations(t)
	organizationRepo.AssertExpectations(t)
}

func TestAcceptInvitationWhenErrorGettingInvitedUser(t *testing.T) {
	t.Parallel()

	userRepo := testUtils.MockedUserRepository{}
	organizationRepo := testUtils.MockedOrganizationRepository{}
	invitationRepo := mockedInvitationRepository{}

	clock := clock.NewFixedClock(time.Now())
	invitationService := invitations.New(&userRepo, &invitationRepo, &organizationRepo, clock)
	user, ctx := fakeUserWithContext()
	request := invitationServiceMessages.AcceptInvitationRequest{Id: gofakeit.UUID()}
	wantErr := gofakeit.ErrorDatabase()

	userRepo.On("GetByID", ctx, user.ID).Return(&schema.User{}, wantErr)

	response, err := invitationService.AcceptInvitation(ctx, &request)

	require.ErrorIs(t, err, wantErr)

	assert.Nil(t, response)

	userRepo.AssertExpectations(t)
	invitationRepo.AssertExpectations(t)
	organizationRepo.AssertExpectations(t)
}

func TestAcceptInvitationWhenErrorGettingInvitation(t *testing.T) {
	t.Parallel()

	unknownDBErr := gofakeit.ErrorDatabase()

	tests := []struct {
		name    string
		gotErr  error
		wantErr error
	}{
		{name: "unknown error getting invitation", gotErr: unknownDBErr, wantErr: unknownDBErr},
		{
			name:    "invitation not found",
			gotErr:  repositories.ErrNotFound,
			wantErr: errors.ErrOrganizationInvitationNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			userRepo := testUtils.MockedUserRepository{}
			organizationRepo := testUtils.MockedOrganizationRepository{}
			invitationRepo := mockedInvitationRepository{}

			clock := clock.NewFixedClock(time.Now())
			invitationService := invitations.New(&userRepo, &invitationRepo, &organizationRepo, clock)
			user, ctx := fakeUserWithContext()

			invitationID := gofakeit.UUID()
			userRepo.On("GetByID", ctx, user.ID).Return(user, nil)
			invitationRepo.On("GetByID", ctx, invitationID).Return(&schema.Invitation{}, test.gotErr)

			request := invitationServiceMessages.AcceptInvitationRequest{Id: invitationID}
			response, err := invitationService.AcceptInvitation(ctx, &request)

			require.ErrorIs(t, err, test.wantErr)
			assert.Nil(t, response)

			userRepo.AssertExpectations(t)
			invitationRepo.AssertExpectations(t)
			organizationRepo.AssertExpectations(t)
		})
	}
}

func TestAcceptInvitationWhenInvitationBelongsToAnotherUser(t *testing.T) {
	t.Parallel()

	userRepo := testUtils.MockedUserRepository{}
	organizationRepo := testUtils.MockedOrganizationRepository{}
	invitationRepo := mockedInvitationRepository{}

	clock := clock.NewFixedClock(time.Now())
	invitationService := invitations.New(&userRepo, &invitationRepo, &organizationRepo, clock)
	user, ctx := fakeUserWithContext()

	invitation := fakeActiveInvitation(gofakeit.Email())
	userRepo.On("GetByID", ctx, user.ID).Return(user, nil)
	invitationRepo.On("GetByID", ctx, invitation.ID).Return(invitation, nil)

	request := invitationServiceMessages.AcceptInvitationRequest{Id: invitation.ID}
	response, err := invitationService.AcceptInvitation(ctx, &request)

	require.ErrorIs(t, err, errors.ErrOrganizationInvitationNotFound)
	assert.Nil(t, response)

	userRepo.AssertExpectations(t)
	invitationRepo.AssertExpectations(t)
	organizationRepo.AssertExpectations(t)
}

func TestAcceptInvitationWhenInvitationExpired(t *testing.T) {
	t.Parallel()

	userRepo := testUtils.MockedUserRepository{}
	organizationRepo := testUtils.MockedOrganizationRepository{}
	invitationRepo := mockedInvitationRepository{}

	clock := clock.NewFixedClock(time.Now())
	invitationService := invitations.New(&userRepo, &invitationRepo, &organizationRepo, clock)
	user, ctx := fakeUserWithContext()

	invitation := schema.Invitation{
		ID:             gofakeit.UUID(),
		Status:         invitationServiceMessages.InvitationStatus_INVITATION_STATUS_EXPIRED.String(),
		OrganizationID: gofakeit.UUID(),
		Email:          user.Email,
		RoleID:         roles.RoleIDMember,
	}
	userRepo.On("GetByID", ctx, user.ID).Return(user, nil)
	invitationRepo.On("GetByID", ctx, invitation.ID).Return(&invitation, nil)

	request := invitationServiceMessages.AcceptInvitationRequest{Id: invitation.ID}
	response, err := invitationService.AcceptInvitation(ctx, &request)

	require.ErrorIs(t, err, errors.ErrOrganizationInvitationNotFound)
	assert.Nil(t, response)

	userRepo.AssertExpectations(t)
	invitationRepo.AssertExpectations(t)
	organizationRepo.AssertExpectations(t)
}

func TestAcceptInvitationWhenErrorCreatingUserMembership(t *testing.T) {
	t.Parallel()

	unknownDBErr := gofakeit.ErrorDatabase()

	tests := []struct {
		name    string
		gotErr  error
		wantErr error
	}{
		{name: "unknown error creating membership", gotErr: unknownDBErr, wantErr: unknownDBErr},
		{
			name:    "membership already exists",
			gotErr:  repositories.ErrAlreadyExists,
			wantErr: errors.ErrOrganizationMembershipAlreadyExists,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			userRepo := testUtils.MockedUserRepository{}
			organizationRepo := testUtils.MockedOrganizationRepository{}
			invitationRepo := mockedInvitationRepository{}

			clock := clock.NewFixedClock(time.Now())
			invitationService := invitations.New(&userRepo, &invitationRepo, &organizationRepo, clock)
			user, ctx := fakeUserWithContext()
			invitation := fakeActiveInvitation(user.Email)
			request := invitationServiceMessages.AcceptInvitationRequest{Id: invitation.ID}

			userRepo.On("GetByID", ctx, user.ID).Return(user, nil)
			invitationRepo.On("GetByID", ctx, invitation.ID).Return(invitation, nil)
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
				Return(test.gotErr)

			response, err := invitationService.AcceptInvitation(ctx, &request)

			require.ErrorIs(t, err, test.wantErr)

			assert.Nil(t, response)

			userRepo.AssertExpectations(t)
			invitationRepo.AssertExpectations(t)
			organizationRepo.AssertExpectations(t)
		})
	}
}

func TestAcceptInvitationWhenErrorUpdatingInvitationStatus(t *testing.T) {
	t.Parallel()

	userRepo := testUtils.MockedUserRepository{}
	organizationRepo := testUtils.MockedOrganizationRepository{}
	invitationRepo := mockedInvitationRepository{}

	clock := clock.NewFixedClock(time.Now())
	invitationService := invitations.New(&userRepo, &invitationRepo, &organizationRepo, clock)
	user, ctx := fakeUserWithContext()
	invitation := fakeActiveInvitation(user.Email)
	request := invitationServiceMessages.AcceptInvitationRequest{Id: invitation.ID}
	wantErr := gofakeit.ErrorDatabase()

	userRepo.On("GetByID", ctx, user.ID).Return(user, nil)
	invitationRepo.On("GetByID", ctx, invitation.ID).Return(invitation, nil)
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
			"UpdateStatusByID",
			ctx,
			invitation.ID,
			invitationServiceMessages.InvitationStatus_INVITATION_STATUS_REDEEMED.String(),
		).
		Return(wantErr)

	response, err := invitationService.AcceptInvitation(ctx, &request)

	require.ErrorIs(t, err, wantErr)

	assert.Nil(t, response)

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
		RoleID:         roles.RoleIDMember,
	}
}

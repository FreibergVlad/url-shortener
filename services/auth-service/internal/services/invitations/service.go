package invitations

import (
	"context"
	"errors"
	"time"

	"github.com/FreibergVlad/url-shortener/auth-service/internal/db/repositories/invitations"
	"github.com/FreibergVlad/url-shortener/auth-service/internal/db/repositories/organizations"
	"github.com/FreibergVlad/url-shortener/auth-service/internal/db/repositories/users"
	"github.com/FreibergVlad/url-shortener/auth-service/internal/db/schema"
	protoMessages "github.com/FreibergVlad/url-shortener/proto/pkg/invitations/messages/v1"
	protoService "github.com/FreibergVlad/url-shortener/proto/pkg/invitations/service/v1"
	grpcUtils "github.com/FreibergVlad/url-shortener/shared/go/pkg/api/grpc/utils"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/clock"
	serviceErrors "github.com/FreibergVlad/url-shortener/shared/go/pkg/errors"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const InviteLifetime = time.Hour * 24 * 7

var (
	invitationStatusActive   = protoMessages.InvitationStatus_INVITATION_STATUS_ACTIVE
	invitationStatusRedemeed = protoMessages.InvitationStatus_INVITATION_STATUS_REDEEMED
)

type InvitationService struct {
	protoService.UnimplementedInvitationServiceServer
	userRepository         users.Repository
	invitationRepository   invitations.Repository
	organizationRepository organizations.Repository
	clock                  clock.Clock
}

func New(
	userRepository users.Repository,
	invitationRepository invitations.Repository,
	organizationRepository organizations.Repository,
	clock clock.Clock,
) *InvitationService {
	return &InvitationService{
		userRepository:         userRepository,
		invitationRepository:   invitationRepository,
		organizationRepository: organizationRepository,
		clock:                  clock,
	}
}

func (s *InvitationService) CreateInvitation(
	ctx context.Context,
	req *protoMessages.CreateInvitationRequest,
) (*protoMessages.CreateInvitationResponse, error) {
	invitedUser, err := s.userRepository.GetByEmail(ctx, req.Email)
	if err == nil {
		invitedUserMemberships, err := s.organizationRepository.ListOrganizationMembershipsByUserID(ctx, invitedUser.ID)
		if err != nil {
			return nil, err
		}

		if invitedUserMemberships.HasOrganization(req.OrganizationId) {
			return nil, serviceErrors.NewValidationError(
				"user %s already belongs to organization %s", req.Email, req.OrganizationId,
			)
		}
	} else if !errors.Is(err, serviceErrors.ErrResourceNotFound) {
		return nil, err
	}

	invite := schema.Invitation{
		ID:             uuid.NewString(),
		Status:         invitationStatusActive.String(),
		OrganizationID: req.OrganizationId,
		Email:          req.Email,
		RoleID:         req.RoleId,
		CreatedAt:      s.clock.Now(),
		ExpiresAt:      s.clock.Now().Add(InviteLifetime),
		CreatedBy:      grpcUtils.UserIDFromIncomingContext(ctx),
	}

	if err = s.invitationRepository.Create(ctx, &invite); err != nil {
		return nil, err
	}

	return &protoMessages.CreateInvitationResponse{
		Invitation: &protoMessages.Invitation{
			Id:             invite.ID,
			OrganizationId: invite.OrganizationID,
			Status:         invitationStatusActive,
			Email:          invite.Email,
			RoleId:         invite.RoleID,
			CreatedAt:      timestamppb.New(invite.CreatedAt),
			ExpiresAt:      timestamppb.New(invite.ExpiresAt),
			CreatedBy:      invite.CreatedBy,
		},
	}, nil
}

func (s *InvitationService) AcceptInvitation(
	ctx context.Context,
	req *protoMessages.AcceptInvitationRequest,
) (*protoMessages.AcceptInvitationResponse, error) {
	userID := grpcUtils.UserIDFromIncomingContext(ctx)
	user, err := s.userRepository.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	invitation, err := s.invitationRepository.GetByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	if invitation.Email != user.Email {
		return nil, serviceErrors.NewPermissionDeniedError("permission denied")
	}
	if invitation.Status != invitationStatusActive.String() {
		return nil, serviceErrors.NewValidationError("invite expired")
	}

	membershipParams := organizations.CreateOrganizationMembershipParams{
		OrganizationID: invitation.OrganizationID,
		UserID:         userID,
		RoleID:         invitation.RoleID,
		CreatedAt:      s.clock.Now(),
	}
	if err = s.organizationRepository.CreateOrganizationMembership(ctx, &membershipParams); err != nil {
		return nil, err
	}

	err = s.invitationRepository.UpdateStatusByID(ctx, invitation.ID, invitationStatusRedemeed.String())
	if err != nil {
		return nil, err
	}

	return &protoMessages.AcceptInvitationResponse{}, nil
}

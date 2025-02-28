package organizations

import (
	"context"
	"errors"
	"fmt"

	"github.com/FreibergVlad/url-shortener/auth-service/internal/db/repositories"
	"github.com/FreibergVlad/url-shortener/auth-service/internal/db/repositories/organizations"
	"github.com/FreibergVlad/url-shortener/auth-service/internal/db/schema"
	"github.com/FreibergVlad/url-shortener/auth-service/internal/services/roles"
	protoMessages "github.com/FreibergVlad/url-shortener/proto/pkg/organizations/messages/v1"
	protoService "github.com/FreibergVlad/url-shortener/proto/pkg/organizations/service/v1"
	grpcUtils "github.com/FreibergVlad/url-shortener/shared/go/pkg/api/grpc/utils"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/clock"
	serviceErrors "github.com/FreibergVlad/url-shortener/shared/go/pkg/errors"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type OrganizationService struct {
	protoService.UnimplementedOrganizationServiceServer
	organizationRepository organizations.Repository
	clock                  clock.Clock
}

func New(organizationRepository organizations.Repository, clock clock.Clock) *OrganizationService {
	return &OrganizationService{
		organizationRepository: organizationRepository,
		clock:                  clock,
	}
}

func (s *OrganizationService) CreateOrganization(
	ctx context.Context,
	req *protoMessages.CreateOrganizationRequest,
) (*protoMessages.CreateOrganizationResponse, error) {
	userID := grpcUtils.UserIDFromIncomingContext(ctx)
	organization := schema.Organization{
		ID:        uuid.NewString(),
		Name:      req.Name,
		Slug:      req.Slug,
		CreatedAt: s.clock.Now(),
		CreatedBy: userID,
	}

	err := s.organizationRepository.Create(ctx, &organization)
	if err != nil {
		if errors.Is(err, repositories.ErrAlreadyExists) {
			return nil, serviceErrors.ErrOrganizationAlreadyExists
		}
		return nil, fmt.Errorf("error creating organization: %w", err)
	}

	return &protoMessages.CreateOrganizationResponse{Organization: organizationToProto(&organization)}, nil
}

func (s *OrganizationService) ListMeOrganizationMemberships(
	ctx context.Context,
	_ *protoMessages.ListMeOrganizationMembershipsRequest,
) (*protoMessages.ListMeOrganizationMembershipsResponse, error) {
	userID := grpcUtils.UserIDFromIncomingContext(ctx)
	memberships, err := s.organizationRepository.ListOrganizationMembershipsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error listing organization memberships %w", err)
	}

	response := make([]*protoMessages.OrganizationMembership, len(memberships))
	for i, membership := range memberships {
		response[i] = &protoMessages.OrganizationMembership{
			Organization: &protoMessages.ShortOrganization{
				Id:   membership.Organization.ID,
				Slug: membership.Organization.Slug,
			},
			Role:      roles.GetRoleProto(membership.RoleID),
			CreatedAt: timestamppb.New(membership.CreatedAt),
		}
	}

	return &protoMessages.ListMeOrganizationMembershipsResponse{Data: response}, nil
}

func organizationToProto(organization *schema.Organization) *protoMessages.Organization {
	return &protoMessages.Organization{
		Id:        organization.ID,
		Name:      organization.Name,
		Slug:      organization.Slug,
		CreatedAt: timestamppb.New(organization.CreatedAt),
		CreatedBy: organization.CreatedBy,
	}
}

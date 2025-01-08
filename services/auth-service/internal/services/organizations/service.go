package organizations

import (
	"context"

	"github.com/FreibergVlad/url-shortener/auth-service/internal/db/repositories/organizations"
	"github.com/FreibergVlad/url-shortener/auth-service/internal/db/schema"
	"github.com/FreibergVlad/url-shortener/auth-service/internal/services/roles"
	protoMessages "github.com/FreibergVlad/url-shortener/proto/pkg/organizations/messages/v1"
	protoService "github.com/FreibergVlad/url-shortener/proto/pkg/organizations/service/v1"
	grpcUtils "github.com/FreibergVlad/url-shortener/shared/go/pkg/api/grpc/utils"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/clock"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type organizationService struct {
	protoService.UnimplementedOrganizationServiceServer
	organizationRepository organizations.Repository
	clock                  clock.Clock
}

func New(
	organizationRepository organizations.Repository,
	clock clock.Clock,
) *organizationService {
	return &organizationService{
		organizationRepository: organizationRepository,
		clock:                  clock,
	}
}

func (s *organizationService) CreateOrganization(
	ctx context.Context,
	req *protoMessages.CreateOrganizationRequest,
) (*protoMessages.CreateOrganizationResponse, error) {
	userId := grpcUtils.UserIDFromIncomingContext(ctx)
	organization := schema.Organization{
		ID:        uuid.NewString(),
		Name:      req.Name,
		Slug:      req.Slug,
		CreatedAt: s.clock.Now(),
		CreatedBy: userId,
	}

	err := s.organizationRepository.Create(ctx, &organization)
	if err != nil {
		return nil, err
	}

	return &protoMessages.CreateOrganizationResponse{Organization: organizationToResponse(&organization)}, nil
}

func (s *organizationService) ListMeOrganizationMemberships(
	ctx context.Context,
	req *protoMessages.ListMeOrganizationMembershipsRequest,
) (*protoMessages.ListMeOrganizationMembershipsResponse, error) {
	userId := grpcUtils.UserIDFromIncomingContext(ctx)
	memberships, err := s.organizationRepository.ListOrganizationMembershipsByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}

	response := make([]*protoMessages.OrganizationMembership, len(memberships))
	for i, m := range memberships {
		response[i] = &protoMessages.OrganizationMembership{
			Organization: &protoMessages.ShortOrganization{
				Id:   m.Organization.ID,
				Slug: m.Organization.Slug,
			},
			Role:      roles.GetRoleProto(m.RoleID),
			CreatedAt: timestamppb.New(m.CreatedAt),
		}
	}

	return &protoMessages.ListMeOrganizationMembershipsResponse{Data: response}, nil
}

func organizationToResponse(organization *schema.Organization) *protoMessages.Organization {
	return &protoMessages.Organization{
		Id:        organization.ID,
		Name:      organization.Name,
		Slug:      organization.Slug,
		CreatedAt: timestamppb.New(organization.CreatedAt),
		CreatedBy: organization.CreatedBy,
	}
}

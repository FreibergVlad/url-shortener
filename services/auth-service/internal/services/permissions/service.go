package permissions

import (
	"context"

	protoMessages "github.com/FreibergVlad/url-shortener/proto/pkg/permissions/messages/v1"
	protoService "github.com/FreibergVlad/url-shortener/proto/pkg/permissions/service/v1"
	grpcUtils "github.com/FreibergVlad/url-shortener/shared/go/pkg/api/grpc/utils"
)

type PermissionResolver interface {
	HasPermissions(ctx context.Context, scopes []string, userId string, organizationId *string) (bool, error)
}

type permissionService struct {
	protoService.UnimplementedPermissionServiceServer
	permissionResolver PermissionResolver
}

func New(permissionResolver PermissionResolver) *permissionService {
	return &permissionService{permissionResolver: permissionResolver}
}

func (s *permissionService) HasPermissions(ctx context.Context, req *protoMessages.HasPermissionsRequest) (*protoMessages.HasPermissionsResponse, error) {
	userId := grpcUtils.UserIDFromIncomingContext(ctx)
	var organizationId *string
	if req.OrganizationContext != nil {
		organizationId = &req.OrganizationContext.OrganizationId
	}
	hasPermissions, err := s.permissionResolver.HasPermissions(ctx, req.Scopes, userId, organizationId)
	if err != nil {
		return nil, err
	}
	return &protoMessages.HasPermissionsResponse{HasPermissions: hasPermissions}, nil
}

package permissions

import (
	"context"
	"fmt"

	protoMessages "github.com/FreibergVlad/url-shortener/proto/pkg/permissions/messages/v1"
	protoService "github.com/FreibergVlad/url-shortener/proto/pkg/permissions/service/v1"
	grpcUtils "github.com/FreibergVlad/url-shortener/shared/go/pkg/api/grpc/utils"
)

type PermissionResolver interface {
	HasPermissions(ctx context.Context, scopes []string, userID string, organizationID *string) (bool, error)
}

type PermissionService struct {
	protoService.UnimplementedPermissionServiceServer
	permissionResolver PermissionResolver
}

func New(permissionResolver PermissionResolver) *PermissionService {
	return &PermissionService{permissionResolver: permissionResolver}
}

func (s *PermissionService) HasPermissions(
	ctx context.Context, req *protoMessages.HasPermissionsRequest,
) (*protoMessages.HasPermissionsResponse, error) {
	userID := grpcUtils.UserIDFromIncomingContext(ctx)

	var organizationID *string
	if req.OrganizationContext != nil {
		organizationID = &req.OrganizationContext.OrganizationId
	}

	hasPermissions, err := s.permissionResolver.HasPermissions(ctx, req.Scopes, userID, organizationID)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve permissions: %w", err)
	}

	return &protoMessages.HasPermissionsResponse{HasPermissions: hasPermissions}, nil
}

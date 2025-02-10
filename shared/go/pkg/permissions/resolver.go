package permissions

import (
	"context"

	permissionMessagesProto "github.com/FreibergVlad/url-shortener/proto/pkg/permissions/messages/v1"
	permissionServiceProto "github.com/FreibergVlad/url-shortener/proto/pkg/permissions/service/v1"
	grpcUtils "github.com/FreibergVlad/url-shortener/shared/go/pkg/api/grpc/utils"
)

type PermissionResolver struct {
	permissionServiceClient permissionServiceProto.PermissionServiceClient
}

func NewPermissionResolver(permissionServiceClient permissionServiceProto.PermissionServiceClient) *PermissionResolver {
	return &PermissionResolver{permissionServiceClient: permissionServiceClient}
}

func (s *PermissionResolver) HasPermissions(
	ctx context.Context, scopes []string, userID string, organizationID *string,
) (bool, error) {
	var organizationContext *permissionMessagesProto.OrganizationContext
	if organizationID != nil {
		organizationContext = &permissionMessagesProto.OrganizationContext{OrganizationId: *organizationID}
	}

	req := permissionMessagesProto.HasPermissionsRequest{
		Scopes:              scopes,
		OrganizationContext: organizationContext,
	}

	resp, err := s.permissionServiceClient.HasPermissions(grpcUtils.OutgoingContextWithUserID(ctx, userID), &req)
	if err != nil {
		return false, err
	}

	return resp.HasPermissions, nil
}

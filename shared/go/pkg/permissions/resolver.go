package permissions

import (
	"context"

	permissionMessagesProto "github.com/FreibergVlad/url-shortener/proto/pkg/permissions/messages/v1"
	permissionServiceProto "github.com/FreibergVlad/url-shortener/proto/pkg/permissions/service/v1"
	grpcUtils "github.com/FreibergVlad/url-shortener/shared/go/pkg/api/grpc/utils"
)

type permissionResolver struct {
	permissionServiceClient permissionServiceProto.PermissionServiceClient
}

func NewPermissionResolver(permissionServiceClient permissionServiceProto.PermissionServiceClient) *permissionResolver {
	return &permissionResolver{permissionServiceClient: permissionServiceClient}
}

func (s *permissionResolver) HasPermissions(ctx context.Context, scopes []string, userId string, organization_id *string) (bool, error) {
	var organizationContext *permissionMessagesProto.OrganizationContext
	if organization_id != nil {
		organizationContext = &permissionMessagesProto.OrganizationContext{OrganizationId: *organization_id}
	}
	req := permissionMessagesProto.HasPermissionsRequest{
		Scopes:              scopes,
		OrganizationContext: organizationContext,
	}
	resp, err := s.permissionServiceClient.HasPermissions(grpcUtils.OutgoingContextWithUserID(ctx, userId), &req)
	if err != nil {
		return false, err
	}
	return resp.HasPermissions, nil
}

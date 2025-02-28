package mocks

import (
	"context"

	permissionServiceMessages "github.com/FreibergVlad/url-shortener/proto/pkg/permissions/messages/v1"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type MockedPermissionServiceClient struct {
	mock.Mock
}

func NewMockedPermissionServiceClient(allow bool) *MockedPermissionServiceClient {
	permissionServiceClient := &MockedPermissionServiceClient{}
	permissionServiceClient.
		On("HasPermissions", mock.Anything, mock.Anything, mock.Anything).
		Return(&permissionServiceMessages.HasPermissionsResponse{HasPermissions: allow}, nil)
	return permissionServiceClient
}

func (c *MockedPermissionServiceClient) HasPermissions(
	ctx context.Context,
	in *permissionServiceMessages.HasPermissionsRequest,
	opts ...grpc.CallOption,
) (*permissionServiceMessages.HasPermissionsResponse, error) {
	args := c.Called(ctx, in, opts)
	return args.Get(0).(*permissionServiceMessages.HasPermissionsResponse), args.Error(1)
}

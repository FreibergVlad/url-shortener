package permissions_test

import (
	"context"
	"testing"

	"github.com/FreibergVlad/url-shortener/auth-service/internal/services/permissions"
	permissionServiceMessages "github.com/FreibergVlad/url-shortener/proto/pkg/permissions/messages/v1"
	grpcUtils "github.com/FreibergVlad/url-shortener/shared/go/pkg/api/grpc/utils"
	serviceErrors "github.com/FreibergVlad/url-shortener/shared/go/pkg/errors"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockedPermissionResolver struct {
	mock.Mock
}

func (r *mockedPermissionResolver) HasPermissions(
	ctx context.Context, scopes []string, userID string, organizationID *string,
) (bool, error) {
	args := r.Called(ctx, scopes, userID, organizationID)
	return args.Bool(0), args.Error(1)
}

func TestHasPermissions(t *testing.T) {
	t.Parallel()

	permissionResolver := mockedPermissionResolver{}
	permissionService := permissions.New(&permissionResolver)
	scopes, userID := []string{"me:read"}, gofakeit.UUID()
	ctx := grpcUtils.IncomingContextWithUserID(context.Background(), userID)
	request := permissionServiceMessages.HasPermissionsRequest{
		Scopes:              scopes,
		OrganizationContext: &permissionServiceMessages.OrganizationContext{OrganizationId: gofakeit.UUID()},
	}

	permissionResolver.
		On("HasPermissions", ctx, scopes, userID, &request.OrganizationContext.OrganizationId).
		Return(true, nil)

	response, err := permissionService.HasPermissions(ctx, &request)

	require.NoError(t, err)

	assert.True(t, response.HasPermissions)

	permissionResolver.AssertExpectations(t)
}

func TestHasPermissionsWhenError(t *testing.T) {
	t.Parallel()

	permissionResolver := mockedPermissionResolver{}
	permissionService := permissions.New(&permissionResolver)
	scopes, userID := []string{"me:read"}, gofakeit.UUID()
	ctx := grpcUtils.IncomingContextWithUserID(context.Background(), userID)
	request := permissionServiceMessages.HasPermissionsRequest{
		Scopes:              scopes,
		OrganizationContext: &permissionServiceMessages.OrganizationContext{OrganizationId: gofakeit.UUID()},
	}

	mockErr := serviceErrors.NewPermissionDeniedError("mock error")
	permissionResolver.
		On("HasPermissions", ctx, scopes, userID, &request.OrganizationContext.OrganizationId).
		Return(false, mockErr)

	response, err := permissionService.HasPermissions(ctx, &request)

	assert.Nil(t, response)

	require.ErrorIs(t, err, mockErr)

	permissionResolver.AssertExpectations(t)
}

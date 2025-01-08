package permissions_test

import (
	"context"
	"errors"
	"testing"

	"github.com/FreibergVlad/url-shortener/auth-service/internal/services/permissions"
	permissionServiceMessages "github.com/FreibergVlad/url-shortener/proto/pkg/permissions/messages/v1"
	grpcUtils "github.com/FreibergVlad/url-shortener/shared/go/pkg/api/grpc/utils"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockedPermissionResolver struct {
	mock.Mock
}

func (r *mockedPermissionResolver) HasPermissions(ctx context.Context, scopes []string, userId string, organizationId *string) (bool, error) {
	args := r.Called(ctx, scopes, userId, organizationId)
	return args.Bool(0), args.Error(1)
}

func TestHasPermissions(t *testing.T) {
	t.Parallel()

	permissionResolver := mockedPermissionResolver{}
	permissionService := permissions.New(&permissionResolver)
	scopes, userId := []string{"me:read"}, gofakeit.UUID()
	ctx := grpcUtils.IncomingContextWithUserID(context.Background(), userId)
	request := permissionServiceMessages.HasPermissionsRequest{
		Scopes:              scopes,
		OrganizationContext: &permissionServiceMessages.OrganizationContext{OrganizationId: gofakeit.UUID()},
	}

	permissionResolver.
		On("HasPermissions", ctx, scopes, userId, &request.OrganizationContext.OrganizationId).
		Return(true, nil)

	response, err := permissionService.HasPermissions(ctx, &request)

	if !assert.NoError(t, err) {
		return
	}
	assert.True(t, response.HasPermissions)

	permissionResolver.AssertExpectations(t)
}

func TestHasPermissionsWhenError(t *testing.T) {
	t.Parallel()

	permissionResolver := mockedPermissionResolver{}
	permissionService := permissions.New(&permissionResolver)
	scopes, userId := []string{"me:read"}, gofakeit.UUID()
	ctx := grpcUtils.IncomingContextWithUserID(context.Background(), userId)
	request := permissionServiceMessages.HasPermissionsRequest{
		Scopes:              scopes,
		OrganizationContext: &permissionServiceMessages.OrganizationContext{OrganizationId: gofakeit.UUID()},
	}

	mockErr := errors.New("mock error")
	permissionResolver.
		On("HasPermissions", ctx, scopes, userId, &request.OrganizationContext.OrganizationId).
		Return(false, mockErr)

	response, err := permissionService.HasPermissions(ctx, &request)

	assert.Nil(t, response)
	assert.ErrorIs(t, err, mockErr)

	permissionResolver.AssertExpectations(t)
}

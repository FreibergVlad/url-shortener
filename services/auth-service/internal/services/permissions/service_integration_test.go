package permissions_test

import (
	"context"
	"testing"

	testUtils "github.com/FreibergVlad/url-shortener/auth-service/internal/testing"
	permissionServiceMessages "github.com/FreibergVlad/url-shortener/proto/pkg/permissions/messages/v1"
	grpcUtils "github.com/FreibergVlad/url-shortener/shared/go/pkg/api/grpc/utils"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/testing/integration"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestHasPermissionsWhenUnauthenticated_Integration(t *testing.T) {
	t.Parallel()
	integration.MaybeSkipIntegrationTest(t)

	server, teardown := testUtils.BootstrapServer()
	defer teardown()

	request := &permissionServiceMessages.HasPermissionsRequest{}

	response, err := server.PermissionServiceClient.HasPermissions(context.Background(), request)

	assert.Nil(t, response)
	assert.ErrorContains(t, err, "Unauthenticated")
}

func TestHasPermissionsWhenUserDoesntExist_Integration(t *testing.T) {
	t.Parallel()
	integration.MaybeSkipIntegrationTest(t)

	server, teardown := testUtils.BootstrapServer()
	defer teardown()

	request := &permissionServiceMessages.HasPermissionsRequest{}
	ctx := grpcUtils.OutgoingContextWithUserID(context.Background(), gofakeit.UUID())

	response, err := server.PermissionServiceClient.HasPermissions(ctx, request)

	assert.Nil(t, response)
	assert.ErrorContains(t, err, "NotFound")
}

func TestHasPermissionsWhenGlobalRoleMatches_Integration(t *testing.T) {
	t.Parallel()
	integration.MaybeSkipIntegrationTest(t)

	server, teardown := testUtils.BootstrapServer()
	defer teardown()

	user := testUtils.CreateTestUser(server, t)
	scopes := []string{"me:read"}
	request := &permissionServiceMessages.HasPermissionsRequest{Scopes: scopes}
	ctx := grpcUtils.OutgoingContextWithUserID(context.Background(), user.Id)

	response, err := server.PermissionServiceClient.HasPermissions(ctx, request)

	if !assert.NoError(t, err) {
		return
	}
	assert.True(t, response.HasPermissions)
}

func TestHasPermissionsWhenHasNoPermissions_Integration(t *testing.T) {
	t.Parallel()
	integration.MaybeSkipIntegrationTest(t)

	server, teardown := testUtils.BootstrapServer()
	defer teardown()

	user := testUtils.CreateTestUser(server, t)
	scopes := []string{"short-url:read"}
	request := &permissionServiceMessages.HasPermissionsRequest{Scopes: scopes}
	ctx := grpcUtils.OutgoingContextWithUserID(context.Background(), user.Id)

	response, err := server.PermissionServiceClient.HasPermissions(ctx, request)

	if !assert.NoError(t, err) {
		return
	}
	assert.False(t, response.HasPermissions)
}

func TestHasPermissionsWhenOrganizationRoleMatches_Integration(t *testing.T) {
	t.Parallel()
	integration.MaybeSkipIntegrationTest(t)

	server, teardown := testUtils.BootstrapServer()
	defer teardown()

	user := testUtils.CreateTestUser(server, t)
	organization := testUtils.CreateTestOrganization(server, user.Id, t)
	scopes := []string{"short-url:read"}
	request := &permissionServiceMessages.HasPermissionsRequest{
		Scopes:              scopes,
		OrganizationContext: &permissionServiceMessages.OrganizationContext{OrganizationId: organization.Id},
	}
	ctx := grpcUtils.OutgoingContextWithUserID(context.Background(), user.Id)

	response, err := server.PermissionServiceClient.HasPermissions(ctx, request)

	if !assert.NoError(t, err) {
		return
	}
	assert.True(t, response.HasPermissions)
}

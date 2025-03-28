package users_test

import (
	"context"
	"testing"

	"github.com/FreibergVlad/url-shortener/auth-service/internal/services/roles"
	testUtils "github.com/FreibergVlad/url-shortener/auth-service/internal/testing"
	userServiceMessages "github.com/FreibergVlad/url-shortener/proto/pkg/users/messages/v1"
	grpcUtils "github.com/FreibergVlad/url-shortener/shared/go/pkg/api/grpc/utils"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/errors"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/testing/asserts"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/testing/integration"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateUser_Integration(t *testing.T) {
	t.Parallel()
	integration.MaybeSkipIntegrationTest(t)

	server, teardown := testUtils.BootstrapServer()

	t.Cleanup(teardown)

	request := testUtils.CreateTestUserRequest()

	response, err := server.UserServiceClient.CreateUser(context.Background(), request)

	require.NoError(t, err)

	assert.Equal(t, request.Email, response.User.Email)
	assert.Equal(t, request.FullName, response.User.FullName)
	assert.Equal(t, roles.RoleIDProvisional, response.User.Role.Id)
	assert.Equal(t, roles.RoleProvisional.Description, response.User.Role.Description)
}

func TestCreateUserWhenDuplicateEmail_Integration(t *testing.T) {
	t.Parallel()
	integration.MaybeSkipIntegrationTest(t)

	server, teardown := testUtils.BootstrapServer()

	t.Cleanup(teardown)

	request := testUtils.CreateTestUserRequest()

	_, err := server.UserServiceClient.CreateUser(context.Background(), request)

	require.NoError(t, err)

	response, err := server.UserServiceClient.CreateUser(context.Background(), request)
	assert.Nil(t, response)
	assert.ErrorIs(t, err, errors.ErrUserAlreadyExists)
}

func TestCreateUserWhenInvalidRequest_Integration(t *testing.T) {
	t.Parallel()
	integration.MaybeSkipIntegrationTest(t)

	server, teardown := testUtils.BootstrapServer()

	t.Cleanup(teardown)

	response, err := server.UserServiceClient.CreateUser(context.Background(), &userServiceMessages.CreateUserRequest{})

	assert.Nil(t, response)

	fieldViolations := map[string][]string{"email": {"value is required"}, "password": {"value is required"}}
	asserts.AssertValidationErrorContainsFieldViolations(t, err, fieldViolations)
}

func TestCreateUserWhenInvalidEmail_Integration(t *testing.T) {
	t.Parallel()
	integration.MaybeSkipIntegrationTest(t)

	server, teardown := testUtils.BootstrapServer()

	t.Cleanup(teardown)

	request := userServiceMessages.CreateUserRequest{
		Email:    "not-an-email",
		Password: gofakeit.Password(true, true, true, true, true, 10),
	}

	response, err := server.UserServiceClient.CreateUser(context.Background(), &request)

	assert.Nil(t, response)
	asserts.AssertValidationErrorContainsFieldViolation(t, err, "email", "must be a valid email address")
}

func TestGetMe_Integration(t *testing.T) {
	t.Parallel()
	integration.MaybeSkipIntegrationTest(t)

	server, teardown := testUtils.BootstrapServer()

	t.Cleanup(teardown)

	user := testUtils.CreateTestUser(t, server)

	getMeRequest := userServiceMessages.GetMeRequest{}
	getMeResponse, err := server.UserServiceClient.GetMe(
		grpcUtils.OutgoingContextWithUserID(context.Background(), user.Id),
		&getMeRequest,
	)

	require.NoError(t, err)

	assert.Equal(t, user.Id, getMeResponse.User.Id)
	assert.Equal(t, user.Email, getMeResponse.User.Email)
	assert.Equal(t, user.FullName, getMeResponse.User.FullName)
	assert.Equal(t, user.Role.Id, getMeResponse.User.Role.Id)
	assert.Equal(t, user.Role.Description, getMeResponse.User.Role.Description)
}

func TestGetMeWhenUnauthenticated_Integration(t *testing.T) {
	t.Parallel()
	integration.MaybeSkipIntegrationTest(t)

	server, teardown := testUtils.BootstrapServer()

	t.Cleanup(teardown)

	request := userServiceMessages.GetMeRequest{}

	response, err := server.UserServiceClient.GetMe(context.Background(), &request)

	assert.Nil(t, response)
	assert.ErrorIs(t, err, errors.ErrUnauthenticated)
}

package testing

import (
	"context"
	"regexp"
	"strings"
	"testing"

	"github.com/FreibergVlad/url-shortener/auth-service/internal/services/roles"
	organizationServiceMessages "github.com/FreibergVlad/url-shortener/proto/pkg/organizations/messages/v1"
	userServiceMessages "github.com/FreibergVlad/url-shortener/proto/pkg/users/messages/v1"
	grpcUtils "github.com/FreibergVlad/url-shortener/shared/go/pkg/api/grpc/utils"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const fakePasswordLength = 10

func CreateTestUserRequest() *userServiceMessages.CreateUserRequest {
	return &userServiceMessages.CreateUserRequest{
		Email:    gofakeit.Email(),
		FullName: gofakeit.Name(),
		Password: gofakeit.Password(true, true, true, true, true, fakePasswordLength),
	}
}

func CreateTestOrganizationRequest() *organizationServiceMessages.CreateOrganizationRequest {
	name := gofakeit.Company()

	return &organizationServiceMessages.CreateOrganizationRequest{
		Name: name,
		Slug: slugify(name),
	}
}

func CreateTestOrganization(
	t *testing.T, server *Server, userID string,
) *organizationServiceMessages.Organization {
	t.Helper()

	request := CreateTestOrganizationRequest()
	ctx := grpcUtils.OutgoingContextWithUserID(context.Background(), userID)

	response, err := server.OrganizationServiceClient.CreateOrganization(ctx, request)

	require.NoError(t, err)

	assert.Equal(t, request.Name, response.Organization.Name)
	assert.Equal(t, request.Slug, response.Organization.Slug)
	assert.Equal(t, userID, response.Organization.CreatedBy)

	return response.Organization
}

func CreateTestUser(t *testing.T, server *Server) *userServiceMessages.User {
	t.Helper()

	request := CreateTestUserRequest()
	response, err := server.UserServiceClient.CreateUser(context.Background(), request)

	require.NoError(t, err)

	assert.Equal(t, request.Email, response.User.Email)
	assert.Equal(t, request.FullName, response.User.FullName)
	assert.Equal(t, roles.RoleIDProvisional, response.User.Role.Id)
	assert.Equal(t, roles.RoleProvisional.Description, response.User.Role.Description)

	return response.User
}

func slugify(input string) string {
	input = strings.ToLower(input)
	input = regexp.MustCompile(`[^a-z0-9\s-]`).ReplaceAllString(input, "")
	input = strings.Join(strings.Fields(input), " ")
	return strings.ReplaceAll(input, " ", "-")
}

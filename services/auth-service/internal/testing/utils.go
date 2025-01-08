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
)

func CreateTestUserRequest() *userServiceMessages.CreateUserRequest {
	return &userServiceMessages.CreateUserRequest{
		Email:    gofakeit.Email(),
		Password: gofakeit.Password(true, true, true, true, true, 10),
	}
}

func CreateTestOrganizationRequest() *organizationServiceMessages.CreateOrganizationRequest {
	name := gofakeit.Company()

	return &organizationServiceMessages.CreateOrganizationRequest{
		Name: name,
		Slug: slugify(name),
	}
}

func CreateTestOrganization(server *TestingServer, userId string, t *testing.T) *organizationServiceMessages.Organization {
	request := CreateTestOrganizationRequest()
	ctx := grpcUtils.OutgoingContextWithUserID(context.Background(), userId)

	response, err := server.OrganizationServiceClient.CreateOrganization(ctx, request)

	assert.NoError(t, err)

	assert.Equal(t, request.Name, response.Organization.Name)
	assert.Equal(t, request.Slug, response.Organization.Slug)
	assert.Equal(t, userId, response.Organization.CreatedBy)

	return response.Organization
}

func CreateTestUser(server *TestingServer, t *testing.T) *userServiceMessages.User {
	request := CreateTestUserRequest()
	response, err := server.UserServiceClient.CreateUser(context.Background(), request)

	assert.NoError(t, err)

	assert.Equal(t, request.Email, response.User.Email)
	assert.Equal(t, roles.RoleIdProvisional, response.User.Role.Id)
	assert.Equal(t, roles.RoleProvisional.Description, response.User.Role.Description)

	return response.User
}

func slugify(input string) string {
	input = strings.ToLower(input)
	input = regexp.MustCompile(`[^a-z0-9\s-]`).ReplaceAllString(input, "")
	input = strings.Join(strings.Fields(input), " ")
	return strings.ReplaceAll(input, " ", "-")
}

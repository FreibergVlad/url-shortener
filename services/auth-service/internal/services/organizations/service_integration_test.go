package organizations_test

import (
	"context"
	"strings"
	"testing"

	"github.com/FreibergVlad/url-shortener/auth-service/internal/services/roles"
	testUtils "github.com/FreibergVlad/url-shortener/auth-service/internal/testing"
	organizationServiceMessages "github.com/FreibergVlad/url-shortener/proto/pkg/organizations/messages/v1"
	grpcUtils "github.com/FreibergVlad/url-shortener/shared/go/pkg/api/grpc/utils"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/testing/integration"
	"github.com/stretchr/testify/assert"
)

func TestCreateOrganizationWhenDuplicateSlug_Integration(t *testing.T) {
	t.Parallel()
	integration.MaybeSkipIntegrationTest(t)

	server, teardown := testUtils.BootstrapServer()
	defer teardown()

	user := testUtils.CreateTestUser(server, t)
	request := testUtils.CreateTestOrganizationRequest()

	_, err := server.OrganizationServiceClient.CreateOrganization(
		grpcUtils.OutgoingContextWithUserID(context.Background(), user.Id),
		request,
	)

	assert.NoError(t, err)

	response, err := server.OrganizationServiceClient.CreateOrganization(
		grpcUtils.OutgoingContextWithUserID(context.Background(), user.Id),
		request,
	)

	assert.Nil(t, response)
	assert.ErrorContains(t, err, "InvalidArgument")
}

func TestCreateOrganizationWhenUnauthenticated_Integration(t *testing.T) {
	t.Parallel()
	integration.MaybeSkipIntegrationTest(t)

	server, teardown := testUtils.BootstrapServer()
	defer teardown()

	request := testUtils.CreateTestOrganizationRequest()

	response, err := server.OrganizationServiceClient.CreateOrganization(context.Background(), request)

	assert.Nil(t, response)
	assert.ErrorContains(t, err, "Unauthenticated")
}

func TestCreateOrganizationWhenInvalidRequest_Integration(t *testing.T) {
	t.Parallel()
	integration.MaybeSkipIntegrationTest(t)

	server, teardown := testUtils.BootstrapServer()
	defer teardown()

	user := testUtils.CreateTestUser(server, t)
	tests := []struct {
		test string
		name string
		slug string
	}{
		{"empty name and slug", "", ""},
		{"too short name and slug", "a", "a"},
		{"too long name and slug", strings.Repeat("a", 51), strings.Repeat("a", 51)},
		{"invalid slug", "Test Name", "TES4 NAME!"},
	}

	for _, input := range tests {
		t.Run(input.name, func(t *testing.T) {
			request := &organizationServiceMessages.CreateOrganizationRequest{
				Name: input.name,
				Slug: input.slug,
			}
			response, err := server.OrganizationServiceClient.CreateOrganization(
				grpcUtils.OutgoingContextWithUserID(context.Background(), user.Id),
				request,
			)

			assert.Nil(t, response)
			assert.ErrorContains(t, err, "InvalidArgument")
		})
	}
}

func TestCreateAndListOrganizationMemberships_Integration(t *testing.T) {
	t.Parallel()
	integration.MaybeSkipIntegrationTest(t)

	server, teardown := testUtils.BootstrapServer()
	defer teardown()

	user := testUtils.CreateTestUser(server, t)
	organization := testUtils.CreateTestOrganization(server, user.Id, t)

	response, err := server.OrganizationServiceClient.ListMeOrganizationMemberships(
		grpcUtils.OutgoingContextWithUserID(context.Background(), user.Id),
		&organizationServiceMessages.ListMeOrganizationMembershipsRequest{},
	)

	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, 1, len(response.Data))

	assert.Equal(t, organization.Id, response.Data[0].Organization.Id)
	assert.Equal(t, organization.Slug, response.Data[0].Organization.Slug)
	assert.Equal(t, roles.RoleIdOwner, response.Data[0].Role.Id)
}

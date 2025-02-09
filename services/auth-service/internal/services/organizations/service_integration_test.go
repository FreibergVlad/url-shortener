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
	"github.com/stretchr/testify/require"
)

func TestCreateOrganizationWhenDuplicateSlug_Integration(t *testing.T) {
	t.Parallel()
	integration.MaybeSkipIntegrationTest(t)

	server, teardown := testUtils.BootstrapServer()

	t.Cleanup(teardown)

	user := testUtils.CreateTestUser(t, server)
	request := testUtils.CreateTestOrganizationRequest()

	_, err := server.OrganizationServiceClient.CreateOrganization(
		grpcUtils.OutgoingContextWithUserID(context.Background(), user.Id),
		request,
	)

	require.NoError(t, err)

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

	t.Cleanup(teardown)

	request := testUtils.CreateTestOrganizationRequest()

	response, err := server.OrganizationServiceClient.CreateOrganization(context.Background(), request)

	require.ErrorContains(t, err, "Unauthenticated")

	assert.Nil(t, response)
}

func TestCreateOrganizationWhenInvalidRequest_Integration(t *testing.T) {
	t.Parallel()
	integration.MaybeSkipIntegrationTest(t)

	server, teardown := testUtils.BootstrapServer()

	t.Cleanup(teardown)

	user := testUtils.CreateTestUser(t, server)
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
			t.Parallel()

			request := &organizationServiceMessages.CreateOrganizationRequest{
				Name: input.name,
				Slug: input.slug,
			}
			response, err := server.OrganizationServiceClient.CreateOrganization(
				grpcUtils.OutgoingContextWithUserID(context.Background(), user.Id),
				request,
			)

			require.ErrorContains(t, err, "InvalidArgument")

			assert.Nil(t, response)
		})
	}
}

func TestCreateAndListOrganizationMemberships_Integration(t *testing.T) {
	t.Parallel()
	integration.MaybeSkipIntegrationTest(t)

	server, teardown := testUtils.BootstrapServer()

	t.Cleanup(teardown)

	user := testUtils.CreateTestUser(t, server)
	organization := testUtils.CreateTestOrganization(t, server, user.Id)

	response, err := server.OrganizationServiceClient.ListMeOrganizationMemberships(
		grpcUtils.OutgoingContextWithUserID(context.Background(), user.Id),
		&organizationServiceMessages.ListMeOrganizationMembershipsRequest{},
	)

	require.NoError(t, err)

	assert.Len(t, response.Data, 1)

	assert.Equal(t, organization.Id, response.Data[0].Organization.Id)
	assert.Equal(t, organization.Slug, response.Data[0].Organization.Slug)
	assert.Equal(t, roles.RoleIDOwner, response.Data[0].Role.Id)
}

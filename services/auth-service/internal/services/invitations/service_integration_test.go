package invitations_test

import (
	"context"
	"testing"

	"github.com/FreibergVlad/url-shortener/auth-service/internal/services/roles"
	testUtils "github.com/FreibergVlad/url-shortener/auth-service/internal/testing"
	invitationServiceMessages "github.com/FreibergVlad/url-shortener/proto/pkg/invitations/messages/v1"
	organizationServiceMessages "github.com/FreibergVlad/url-shortener/proto/pkg/organizations/messages/v1"
	grpcUtils "github.com/FreibergVlad/url-shortener/shared/go/pkg/api/grpc/utils"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/testing/integration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateAndAcceptInvitation(t *testing.T) {
	t.Parallel()
	integration.MaybeSkipIntegrationTest(t)

	server, teardown := testUtils.BootstrapServer()
	defer teardown()

	organizationOwner := testUtils.CreateTestUser(t, server)
	organization := testUtils.CreateTestOrganization(t, server, organizationOwner.Id)
	invitedUser := testUtils.CreateTestUser(t, server)

	createResponse, err := server.InvitationServiceClient.CreateInvitation(
		grpcUtils.OutgoingContextWithUserID(context.Background(), organizationOwner.Id),
		&invitationServiceMessages.CreateInvitationRequest{
			OrganizationId: organization.Id,
			Email:          invitedUser.Email,
			RoleId:         roles.RoleIDAdmin,
		},
	)

	require.NoError(t, err)
	require.NotNil(t, createResponse)

	assert.Equal(t, organization.Id, createResponse.Invitation.OrganizationId)
	assert.Equal(t, invitationServiceMessages.InvitationStatus_INVITATION_STATUS_ACTIVE, createResponse.Invitation.Status)
	assert.Equal(t, invitedUser.Email, createResponse.Invitation.Email)
	assert.Equal(t, roles.RoleIDAdmin, createResponse.Invitation.RoleId)
	assert.Equal(t, organizationOwner.Id, createResponse.Invitation.CreatedBy)

	acceptResponse, err := server.InvitationServiceClient.AcceptInvitation(
		grpcUtils.OutgoingContextWithUserID(context.Background(), invitedUser.Id),
		&invitationServiceMessages.AcceptInvitationRequest{Id: createResponse.Invitation.Id},
	)

	require.NoError(t, err)
	require.NotNil(t, acceptResponse)

	membershipsResponse, err := server.OrganizationServiceClient.ListMeOrganizationMemberships(
		grpcUtils.OutgoingContextWithUserID(context.Background(), invitedUser.Id),
		&organizationServiceMessages.ListMeOrganizationMembershipsRequest{},
	)

	require.NoError(t, err)

	assert.Len(t, membershipsResponse.Data, 1)
	assert.Equal(t, roles.RoleIDAdmin, membershipsResponse.Data[0].Role.Id)
	assert.Equal(t, organization.Id, membershipsResponse.Data[0].Organization.Id)
}

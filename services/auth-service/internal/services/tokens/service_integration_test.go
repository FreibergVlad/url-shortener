package tokens_test

import (
	"context"
	"testing"

	testUtils "github.com/FreibergVlad/url-shortener/auth-service/internal/testing"
	tokenServiceMessages "github.com/FreibergVlad/url-shortener/proto/pkg/tokens/messages/v1"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/testing/integration"
	"github.com/stretchr/testify/assert"
)

func TestIssueAndRefreshAuthenticationTokenFlow(t *testing.T) {
	t.Parallel()
	integration.MaybeSkipIntegrationTest(t)

	server, teardown := testUtils.BootstrapServer()
	defer teardown()

	createUserRequest := testUtils.CreateTestUserRequest()
	_, err := server.UserServiceClient.CreateUser(context.Background(), createUserRequest)
	assert.NoError(t, err)

	issueTokenRequest := tokenServiceMessages.IssueAuthenticationTokenRequest{
		Email:    createUserRequest.Email,
		Password: createUserRequest.Password,
	}
	issueTokenResponse, err := server.TokenServiceClient.IssueAuthenticationToken(context.Background(), &issueTokenRequest)
	assert.NoError(t, err)
	assert.NotEmpty(t, issueTokenResponse.Token)
	assert.NotEmpty(t, issueTokenResponse.RefreshToken)

	refreshTokenRequest := tokenServiceMessages.RefreshAuthenticationTokenRequest{RefreshToken: issueTokenResponse.RefreshToken}
	refreshTokenResponse, err := server.TokenServiceClient.RefreshAuthenticationToken(context.Background(), &refreshTokenRequest)
	assert.NoError(t, err)
	assert.NotEmpty(t, refreshTokenResponse.Token)
}

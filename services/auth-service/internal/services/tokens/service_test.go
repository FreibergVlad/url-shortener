package tokens_test

import (
	"context"
	"testing"
	"time"

	"github.com/FreibergVlad/url-shortener/auth-service/internal/config"
	"github.com/FreibergVlad/url-shortener/auth-service/internal/db/schema"
	"github.com/FreibergVlad/url-shortener/auth-service/internal/services/tokens"
	testUtils "github.com/FreibergVlad/url-shortener/auth-service/internal/testing"
	tokenServiceMessages "github.com/FreibergVlad/url-shortener/proto/pkg/tokens/messages/v1"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/clock"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/errors"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/jwt"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/must"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	fakePassword     = "fake-pass"
	fakePasswordHash = "$2a$10$duILFNZWb9k1VAnQoON1Zu9NpA70BylznDP6lpU6yJIfHPMB/Kfza" //nolint:gosec
)

func TestIssueAuthenticationToken(t *testing.T) {
	t.Parallel()

	config := config.IdentityServiceConfig{}
	userRepo := testUtils.MockedUserRepository{}
	tokenService := tokens.New(config, &userRepo, clock.NewFixedClock(time.Now()))
	user := schema.User{
		ID:           gofakeit.UUID(),
		Email:        gofakeit.Email(),
		PasswordHash: fakePasswordHash,
	}
	request := tokenServiceMessages.IssueAuthenticationTokenRequest{
		Email:    user.Email,
		Password: fakePassword,
	}
	ctx := context.Background()

	userRepo.On("GetByEmail", ctx, user.Email).Return(&user, nil)

	response, err := tokenService.IssueAuthenticationToken(context.Background(), &request)

	require.NoError(t, err)

	assert.NotEmpty(t, response.Token)
	assert.NotEmpty(t, response.RefreshToken)

	userRepo.AssertExpectations(t)
}

func TestIssueAuthenticationTokenWhenDatabaseError(t *testing.T) {
	t.Parallel()

	config := config.IdentityServiceConfig{}
	userRepo := testUtils.MockedUserRepository{}
	tokenService := tokens.New(config, &userRepo, clock.NewFixedClock(time.Now()))
	request := tokenServiceMessages.IssueAuthenticationTokenRequest{
		Email:    gofakeit.Email(),
		Password: fakePassword,
	}
	ctx := context.Background()

	userRepo.On("GetByEmail", ctx, request.Email).Return(&schema.User{}, errors.ErrResourceNotFound)

	response, err := tokenService.IssueAuthenticationToken(context.Background(), &request)

	require.ErrorContains(t, err, "invalid credentials")

	assert.Nil(t, response)

	userRepo.AssertExpectations(t)
}

func TestIssueAuthenticationTokenWhenInvalidPassword(t *testing.T) {
	t.Parallel()

	config := config.IdentityServiceConfig{}
	userRepo := testUtils.MockedUserRepository{}
	tokenService := tokens.New(config, &userRepo, clock.NewFixedClock(time.Now()))
	user := schema.User{
		ID:           gofakeit.UUID(),
		Email:        gofakeit.Email(),
		PasswordHash: fakePasswordHash,
	}
	request := tokenServiceMessages.IssueAuthenticationTokenRequest{
		Email:    user.Email,
		Password: "invalid",
	}
	ctx := context.Background()

	userRepo.On("GetByEmail", ctx, user.Email).Return(&user, nil)

	response, err := tokenService.IssueAuthenticationToken(ctx, &request)

	require.ErrorContains(t, err, "invalid credentials")

	assert.Nil(t, response)

	userRepo.AssertExpectations(t)
}

func TestRefreshAuthenticationToken(t *testing.T) {
	t.Parallel()

	config := config.IdentityServiceConfig{}
	userRepo := testUtils.MockedUserRepository{}
	tokenService := tokens.New(config, &userRepo, clock.NewFixedClock(time.Now()))
	refreshToken := must.Return(jwt.IssueForUserID(gofakeit.UUID(), "", time.Now(), 10))
	request := tokenServiceMessages.RefreshAuthenticationTokenRequest{RefreshToken: refreshToken}

	response, err := tokenService.RefreshAuthenticationToken(context.Background(), &request)

	require.NoError(t, err)

	assert.NotEmpty(t, response.Token)
}

func TestRefreshAuthenticationTokenWhenInvalidToken(t *testing.T) {
	t.Parallel()

	config := config.IdentityServiceConfig{}
	userRepo := testUtils.MockedUserRepository{}
	tokenService := tokens.New(config, &userRepo, clock.NewFixedClock(time.Now()))
	request := tokenServiceMessages.RefreshAuthenticationTokenRequest{RefreshToken: ""}

	response, err := tokenService.RefreshAuthenticationToken(context.Background(), &request)

	require.ErrorContains(t, err, "token is malformed")

	assert.Nil(t, response)
}

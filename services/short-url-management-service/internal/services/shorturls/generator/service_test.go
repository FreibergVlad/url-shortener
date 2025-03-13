package generator_test

import (
	"context"
	"testing"
	"time"

	shorturlgeneratorproto "github.com/FreibergVlad/url-shortener/proto/pkg/shorturls/generator/messages/v1"
	grpcutils "github.com/FreibergVlad/url-shortener/shared/go/pkg/api/grpc/utils"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/clock"
	serviceerrors "github.com/FreibergVlad/url-shortener/shared/go/pkg/errors"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/testing/asserts"
	"github.com/FreibergVlad/url-shortener/short-url-management-service/internal/clients/domains"
	"github.com/FreibergVlad/url-shortener/short-url-management-service/internal/config"
	shorturlrepository "github.com/FreibergVlad/url-shortener/short-url-management-service/internal/db/repositories/shorturls"
	"github.com/FreibergVlad/url-shortener/short-url-management-service/internal/db/schema"
	shorturlgeneratorservice "github.com/FreibergVlad/url-shortener/short-url-management-service/internal/services/shorturls/generator"
	testutils "github.com/FreibergVlad/url-shortener/short-url-management-service/internal/testing"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockedAliasGenerator struct {
	mock.Mock
}

func (g *MockedAliasGenerator) Generate(shorturl *schema.ShortURL) (string, error) {
	args := g.Called(shorturl)
	return args.String(0), args.Error(1)
}

func (g *MockedAliasGenerator) GenerateWithSalt(shorturl *schema.ShortURL) (string, error) {
	args := g.Called(shorturl)
	return args.String(0), args.Error(1)
}

type MockedUserInfoGetter struct {
	mock.Mock
}

func (g *MockedUserInfoGetter) GetUserInfo(ctx context.Context, userID string) (schema.User, error) {
	args := g.Called(ctx, userID)
	return args.Get(0).(schema.User), args.Error(1)
}

type MockedDomainChecker struct {
	mock.Mock
}

func (c *MockedDomainChecker) HasDomain(ctx context.Context, userID, organizationID, fqdn string) error {
	return c.Called(ctx, userID, organizationID, fqdn).Error(0)
}

func TestCreateShortURLWhenErrorValidatingURL(t *testing.T) {
	t.Parallel()

	service := shorturlgeneratorservice.New(
		&testutils.MockedShortURLRepository{},
		&MockedDomainChecker{},
		&MockedUserInfoGetter{},
		&MockedAliasGenerator{},
		clock.NewFixedClock(time.Now()),
		fakeConfig(),
	)
	request := &shorturlgeneratorproto.CreateShortURLRequest{LongUrl: "not a URL"}

	response, err := service.CreateShortURL(context.Background(), request)

	asserts.AssertValidationErrorContainsFieldViolation(t, err, "long_url", "invalid URL")
	assert.Nil(t, response)
}

func TestCreateShortURLWhenErrorGettingUserInfo(t *testing.T) {
	t.Parallel()

	userID := gofakeit.UUID()
	ctx := grpcutils.IncomingContextWithUserID(context.Background(), userID)

	userInfoGetter := &MockedUserInfoGetter{}
	userInfoGetterErr := gofakeit.ErrorGRPC()
	userInfoGetter.On("GetUserInfo", ctx, userID).Return(schema.User{}, userInfoGetterErr)

	service := shorturlgeneratorservice.New(
		&testutils.MockedShortURLRepository{},
		&MockedDomainChecker{},
		userInfoGetter,
		&MockedAliasGenerator{},
		clock.NewFixedClock(time.Now()),
		fakeConfig(),
	)
	request := &shorturlgeneratorproto.CreateShortURLRequest{LongUrl: gofakeit.URL()}

	response, err := service.CreateShortURL(ctx, request)

	require.ErrorIs(t, err, userInfoGetterErr)
	require.Nil(t, response)

	userInfoGetter.AssertExpectations(t)
}

func TestCreateShortURLWhenErrorCheckingDomain(t *testing.T) {
	t.Parallel()

	grpcError := gofakeit.ErrorGRPC()
	tests := []struct {
		name             string
		domainCheckerErr error
		serviceErr       error
	}{
		{
			"domain not allowed",
			domains.ErrDomainNotAllowed,
			shorturlgeneratorservice.ErrDomainNotAllowed,
		},
		{"internal error", grpcError, grpcError},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			userID := gofakeit.UUID()
			ctx := grpcutils.IncomingContextWithUserID(context.Background(), userID)

			request := &shorturlgeneratorproto.CreateShortURLRequest{
				OrganizationId: gofakeit.UUID(),
				Domain:         gofakeit.DomainName(),
				LongUrl:        gofakeit.URL(),
			}

			userInfoGetter := &MockedUserInfoGetter{}
			userInfoGetter.On("GetUserInfo", ctx, userID).Return(schema.User{ID: userID}, nil)

			domainChecker := &MockedDomainChecker{}
			domainChecker.
				On("HasDomain", ctx, userID, request.OrganizationId, request.Domain).
				Return(test.domainCheckerErr)

			service := shorturlgeneratorservice.New(
				&testutils.MockedShortURLRepository{},
				domainChecker,
				userInfoGetter,
				&MockedAliasGenerator{},
				clock.NewFixedClock(time.Now()),
				fakeConfig(),
			)

			response, err := service.CreateShortURL(ctx, request)

			require.ErrorIs(t, err, test.serviceErr)
			require.Nil(t, response)

			userInfoGetter.AssertExpectations(t)
			domainChecker.AssertExpectations(t)
		})
	}
}

func TestCreateShortURLWithCustomAlias(t *testing.T) {
	t.Parallel()

	fakeDatabaseErr := gofakeit.ErrorDatabase()
	tests := []struct {
		name       string
		repoErr    error
		serviceErr error
	}{
		{"success", nil, nil},
		{"short URL already exists", shorturlrepository.ErrAlreadyExists, serviceerrors.ErrShortURLAlreadyExists},
		{"internal DB error", fakeDatabaseErr, fakeDatabaseErr},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			wantShorturl := &schema.ShortURL{}
			repository := &testutils.MockedShortURLRepository{}

			repository.On("Create", ctx, wantShorturl).Return(test.repoErr)

			service := shorturlgeneratorservice.New(
				repository,
				&MockedDomainChecker{},
				&MockedUserInfoGetter{},
				&MockedAliasGenerator{},
				clock.NewFixedClock(time.Now()),
				fakeConfig(),
			)

			gotShorturl, err := service.СreateWithCustomAlias(ctx, wantShorturl)

			if test.serviceErr != nil {
				require.Nil(t, gotShorturl)
				require.ErrorIs(t, err, test.serviceErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, wantShorturl, gotShorturl)
			}

			repository.AssertExpectations(t)
		})
	}
}

func TestCreateShortURLWithGeneratedAlias(t *testing.T) {
	t.Parallel()

	fakeDatabaseErr := gofakeit.ErrorDatabase()
	fakeGeneratorErr := gofakeit.ErrorRuntime()

	tests := []struct {
		name         string
		generatorErr error
		repoErr      error
		serviceErr   error
	}{
		{"success", nil, nil, nil},
		{"error generating alias", fakeGeneratorErr, nil, fakeGeneratorErr},
		{"internal DB error", nil, fakeDatabaseErr, fakeDatabaseErr},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			shorturl := &schema.ShortURL{}
			alias := "fake-alias"

			generator := &MockedAliasGenerator{}
			generator.On("Generate", shorturl).Return(alias, test.generatorErr)

			repository := &testutils.MockedShortURLRepository{}
			if test.generatorErr == nil {
				repository.On("Create", ctx, shorturl).Return(test.repoErr)
			}

			service := shorturlgeneratorservice.New(
				repository,
				&MockedDomainChecker{},
				&MockedUserInfoGetter{},
				generator,
				clock.NewFixedClock(time.Now()),
				fakeConfig(),
			)

			gotShorturl, err := service.CreateWithGeneratedAlias(ctx, shorturl)

			if test.serviceErr != nil {
				require.Nil(t, gotShorturl)
				require.ErrorIs(t, err, test.serviceErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, shorturl, gotShorturl)
			}

			repository.AssertExpectations(t)
			generator.AssertExpectations(t)
		})
	}
}

func TestCreateShortURLWithGeneratedAliasWhenCollisionDetected(t *testing.T) {
	t.Parallel()

	fakeDatabaseErr := gofakeit.ErrorDatabase()
	fakeGeneratorErr := gofakeit.ErrorRuntime()

	alias, saltedAlias := "fake-alias", "fake-salted-alias"

	tests := []struct {
		name                string
		repoErrOnGet        error
		repoErrOnCreate     error
		conflictingShorturl *schema.ShortURL
		generatorErr        error
		serviceErr          error
	}{
		{"internal DB error while detecting collision", fakeDatabaseErr, nil, nil, nil, fakeDatabaseErr},
		{
			"short URL already exists within organization",
			nil,
			nil,
			&schema.ShortURL{Alias: alias},
			nil,
			serviceerrors.ErrShortURLAlreadyExists,
		},
		{
			"collision detected: error generating salted alias",
			nil,
			nil,
			&schema.ShortURL{OrganizationID: "another"},
			fakeGeneratorErr,
			fakeGeneratorErr,
		},
		{
			"collision detected: internal error creating short URL with salted alias",
			nil,
			fakeDatabaseErr,
			&schema.ShortURL{OrganizationID: "another"},
			nil,
			fakeDatabaseErr,
		},
		{
			"collision detected: exceede all retries while creating short URL with salted alias",
			nil,
			shorturlrepository.ErrAlreadyExists,
			&schema.ShortURL{OrganizationID: "another"},
			nil,
			shorturlgeneratorservice.ErrShortURLCollisionRetriesExceeded,
		},
		{
			"collision detected: short URL successfully created with salted alias",
			nil,
			nil,
			&schema.ShortURL{OrganizationID: "another"},
			nil,
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			shorturl := &schema.ShortURL{Alias: alias}

			generator := &MockedAliasGenerator{}
			generator.On("Generate", shorturl).Return(alias, nil)
			generator.On("GenerateWithSalt", shorturl).Return(saltedAlias, test.generatorErr)

			repository := &testutils.MockedShortURLRepository{}
			repository.On("Create", ctx, shorturl).Return(shorturlrepository.ErrAlreadyExists).Once()
			repository.On("Create", ctx, shorturl).Return(test.repoErrOnCreate).Once()
			repository.
				On("GetByOrganizationKeyOrGlobalKey", ctx, shorturl.OrganizationKey(), shorturl.GlobalKey()).
				Return(test.conflictingShorturl, test.repoErrOnGet)

			service := shorturlgeneratorservice.New(
				repository,
				&MockedDomainChecker{},
				&MockedUserInfoGetter{},
				generator,
				clock.NewFixedClock(time.Now()),
				fakeConfig(),
			)

			gotShorturl, err := service.CreateWithGeneratedAlias(ctx, shorturl)

			if test.serviceErr != nil {
				require.Nil(t, gotShorturl)
				require.ErrorIs(t, err, test.serviceErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, shorturl, gotShorturl)
			}
		})
	}
}

func fakeConfig() config.ShortURLGeneratorConfig {
	return config.ShortURLGeneratorConfig{MaxRetriesOnCollision: 1}
}

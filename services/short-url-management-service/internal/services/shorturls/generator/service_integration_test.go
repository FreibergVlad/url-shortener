package generator_test

import (
	"context"
	"strings"
	"testing"
	"time"

	shortUrlGeneratorMessages "github.com/FreibergVlad/url-shortener/proto/pkg/shorturls/generator/messages/v1"
	shortUrlManagementMessages "github.com/FreibergVlad/url-shortener/proto/pkg/shorturls/management/messages/v1"
	grpcUtils "github.com/FreibergVlad/url-shortener/shared/go/pkg/api/grpc/utils"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/errors"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/testing/asserts"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/testing/integration"
	testUtils "github.com/FreibergVlad/url-shortener/short-url-management-service/internal/testing"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var ShortURLExpirationDate = time.Now().Add(time.Hour * 24)

func TestCreateShortUrlWhenInvalidRequest_Integration(t *testing.T) {
	t.Parallel()
	integration.MaybeSkipIntegrationTest(t)

	server, teardown := testUtils.BootstrapServer()

	t.Cleanup(teardown)

	tests := []struct {
		name            string
		organizationID  string
		domain          string
		alias           *wrapperspb.StringValue
		longURL         string
		expiresAt       time.Time
		tags            []string
		description     string
		fieldViolations map[string][]string
	}{
		{
			name:            "invalid organization ID",
			organizationID:  "fake-id",
			domain:          gofakeit.DomainName(),
			longURL:         gofakeit.URL(),
			expiresAt:       ShortURLExpirationDate,
			tags:            []string{},
			description:     "fake-description",
			fieldViolations: map[string][]string{"organization_id": {"value must be a valid UUID"}},
		},
		{
			name:            "invalid domain",
			organizationID:  gofakeit.UUID(),
			domain:          gofakeit.DomainName(),
			longURL:         gofakeit.URL(),
			expiresAt:       ShortURLExpirationDate,
			tags:            []string{},
			description:     "fake-description",
			fieldViolations: map[string][]string{"domain": {"domain is not allowed to use"}},
		},
		{
			name:            "invalid alias (not base62)",
			organizationID:  gofakeit.UUID(),
			domain:          gofakeit.DomainName(),
			alias:           wrapperspb.String("!nv&alid-alias"),
			longURL:         gofakeit.URL(),
			expiresAt:       ShortURLExpirationDate,
			tags:            []string{},
			description:     "fake-description",
			fieldViolations: map[string][]string{"alias": {"value does not match regex pattern"}},
		},
		{
			name:            "invalid alias (too short)",
			organizationID:  gofakeit.UUID(),
			domain:          gofakeit.DomainName(),
			alias:           wrapperspb.String("a"),
			longURL:         gofakeit.URL(),
			expiresAt:       ShortURLExpirationDate,
			tags:            []string{},
			description:     "fake-description",
			fieldViolations: map[string][]string{"alias": {"value length must be at least 4"}},
		},
		{
			name:            "invalid alias (too long)",
			organizationID:  gofakeit.UUID(),
			domain:          gofakeit.DomainName(),
			alias:           wrapperspb.String(strings.Repeat("a", 31)),
			longURL:         gofakeit.URL(),
			expiresAt:       ShortURLExpirationDate,
			tags:            []string{},
			description:     "fake-description",
			fieldViolations: map[string][]string{"alias": {"value length must be at most 30"}},
		},
		{
			name:            "invalid long URL",
			organizationID:  gofakeit.UUID(),
			domain:          gofakeit.DomainName(),
			longURL:         "fake-url",
			expiresAt:       ShortURLExpirationDate,
			tags:            []string{},
			description:     "fake-description",
			fieldViolations: map[string][]string{"long_url": {"invalid URI"}},
		},
		{
			name:            "invalid expiration date (date in past)",
			organizationID:  gofakeit.UUID(),
			domain:          gofakeit.DomainName(),
			longURL:         gofakeit.URL(),
			expiresAt:       time.Now().Add(time.Hour * -24),
			tags:            []string{},
			description:     "fake-description",
			fieldViolations: map[string][]string{"expires_at": {"value must be greater than now"}},
		},
		{
			name:            "invalid tags (duplicates)",
			organizationID:  gofakeit.UUID(),
			domain:          gofakeit.DomainName(),
			longURL:         gofakeit.URL(),
			expiresAt:       ShortURLExpirationDate,
			tags:            []string{"tag1", "tag1"},
			description:     "fake-description",
			fieldViolations: map[string][]string{"tags": {"repeated value must contain unique items"}},
		},
		{
			name:            "invalid tags (too short tag)",
			organizationID:  gofakeit.UUID(),
			domain:          gofakeit.DomainName(),
			longURL:         gofakeit.URL(),
			expiresAt:       ShortURLExpirationDate,
			tags:            []string{""},
			description:     "fake-description",
			fieldViolations: map[string][]string{"tags": {"value length must be at least 1"}},
		},
		{
			name:            "invalid tags (too long tag)",
			organizationID:  gofakeit.UUID(),
			domain:          gofakeit.DomainName(),
			longURL:         gofakeit.URL(),
			expiresAt:       ShortURLExpirationDate,
			tags:            []string{strings.Repeat("a", 51)},
			description:     "fake-description",
			fieldViolations: map[string][]string{"tags": {"value length must be at most 50"}},
		},
		{
			name:            "invalid tags (too long tag array)",
			organizationID:  gofakeit.UUID(),
			domain:          gofakeit.DomainName(),
			longURL:         gofakeit.URL(),
			expiresAt:       ShortURLExpirationDate,
			tags:            []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11"},
			description:     "fake-description",
			fieldViolations: map[string][]string{"tags": {"value must contain no more than 10"}},
		},
		{
			name:            "invalid description (too long)",
			organizationID:  gofakeit.UUID(),
			domain:          gofakeit.DomainName(),
			longURL:         gofakeit.URL(),
			expiresAt:       ShortURLExpirationDate,
			tags:            []string{},
			description:     strings.Repeat("a", 1001),
			fieldViolations: map[string][]string{"description": {"value length must be at most 1000"}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			request := shortUrlGeneratorMessages.CreateShortURLRequest{
				OrganizationId: test.organizationID,
				Domain:         test.domain,
				Alias:          test.alias,
				LongUrl:        test.longURL,
				ExpiresAt:      timestamppb.New(test.expiresAt),
				Tags:           test.tags,
				Description:    test.description,
			}

			response, err := server.ShortURLGeneratorServiceClient.CreateShortURL(
				grpcUtils.OutgoingContextWithUserID(context.Background(), gofakeit.UUID()),
				&request,
			)

			assert.Nil(t, response)
			asserts.AssertValidationErrorContainsFieldViolations(t, err, test.fieldViolations)
		})
	}
}

func TestCreateShortUrlWithGeneratedAlias_Integration(t *testing.T) {
	t.Parallel()
	integration.MaybeSkipIntegrationTest(t)

	server, teardown := testUtils.BootstrapServer()

	t.Cleanup(teardown)

	testUtils.AssertCreateShortURL(t, server, testUtils.CreateShortURLRequest())
}

func TestCreateShortUrlWithCustomAlias_Integration(t *testing.T) {
	t.Parallel()
	integration.MaybeSkipIntegrationTest(t)

	server, teardown := testUtils.BootstrapServer()

	t.Cleanup(teardown)

	request := testUtils.CreateShortURLRequest()
	request.Alias = wrapperspb.String(testUtils.RandomBase62StringInRange(4, 30))

	testUtils.AssertCreateShortURL(t, server, request)
}

func TestCreateShortUrlWithCustomAliasWhenDuplicate_Integration(t *testing.T) {
	t.Parallel()
	integration.MaybeSkipIntegrationTest(t)

	server, teardown := testUtils.BootstrapServer()

	t.Cleanup(teardown)

	request := testUtils.CreateShortURLRequest()

	shorturl := testUtils.AssertCreateShortURL(t, server, request)

	request.Alias = wrapperspb.String(shorturl.Alias)

	response, err := server.ShortURLGeneratorServiceClient.CreateShortURL(
		grpcUtils.OutgoingContextWithUserID(context.Background(), testUtils.FakeUser.ID),
		request,
	)

	assert.Nil(t, response)
	assert.ErrorIs(t, err, errors.ErrShortURLAlreadyExists)
}

func TestCreateShortUrlWithGeneratedAliasWhenDuplicate_Integration(t *testing.T) {
	t.Parallel()
	integration.MaybeSkipIntegrationTest(t)

	server, teardown := testUtils.BootstrapServer()

	t.Cleanup(teardown)

	request := testUtils.CreateShortURLRequest()

	testUtils.AssertCreateShortURL(t, server, request)

	response, err := server.ShortURLGeneratorServiceClient.CreateShortURL(
		grpcUtils.OutgoingContextWithUserID(context.Background(), testUtils.FakeUser.ID),
		request,
	)

	assert.Nil(t, response)
	assert.ErrorIs(t, err, errors.ErrShortURLAlreadyExists)
}

func TestCreateShortUrlWithGeneratedAliasWhenCollisionDetected_Integration(t *testing.T) {
	t.Parallel()
	integration.MaybeSkipIntegrationTest(t)

	server, teardown := testUtils.BootstrapServer()

	t.Cleanup(teardown)

	request := testUtils.CreateShortURLRequest()
	shorturl := testUtils.AssertCreateShortURL(t, server, request)

	_, err := server.ShortURLManagementServiceClient.UpdateShortURLByOrganizationID(
		grpcUtils.OutgoingContextWithUserID(context.Background(), gofakeit.UUID()),
		&shortUrlManagementMessages.UpdateShortURLByOrganizationIDRequest{
			OrganizationId: shorturl.OrganizationId,
			Domain:         shorturl.Domain,
			Alias:          shorturl.Alias,
			ShortUrl:       &shortUrlManagementMessages.UpdateShortURLParams{LongUrl: gofakeit.URL()},
			UpdateMask:     &fieldmaskpb.FieldMask{Paths: []string{"long_url"}},
		},
	)

	require.NoError(t, err)

	testUtils.AssertCreateShortURL(t, server, request)
}

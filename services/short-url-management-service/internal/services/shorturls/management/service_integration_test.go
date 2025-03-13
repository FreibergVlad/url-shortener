package management_test

import (
	"context"
	"strings"
	"testing"
	"time"

	shortURLManagementMessages "github.com/FreibergVlad/url-shortener/proto/pkg/shorturls/management/messages/v1"
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

func TestGetShortURLByOrganizationID_Integration(t *testing.T) {
	t.Parallel()
	integration.MaybeSkipIntegrationTest(t)

	server, teardown := testUtils.BootstrapServer()

	t.Cleanup(teardown)

	shorturl := testUtils.AssertCreateShortURL(t, server, testUtils.CreateShortURLRequest())

	response, err := server.ShortURLManagementServiceClient.GetShortURLByOrganizationID(
		grpcUtils.OutgoingContextWithUserID(context.Background(), shorturl.CreatedBy.Id),
		&shortURLManagementMessages.GetShortURLByOrganizationIDRequest{
			OrganizationId: shorturl.OrganizationId,
			Domain:         shorturl.Domain,
			Alias:          shorturl.Alias,
		},
	)

	require.NoError(t, err)
	testUtils.AssertShortURLEqual(t, shorturl, response.ShortUrl)
}

func TestGetShortURLByOrganizationIDWhenShortURLNotExists_Integration(t *testing.T) {
	t.Parallel()
	integration.MaybeSkipIntegrationTest(t)

	server, teardown := testUtils.BootstrapServer()

	t.Cleanup(teardown)

	response, err := server.ShortURLManagementServiceClient.GetShortURLByOrganizationID(
		grpcUtils.OutgoingContextWithUserID(context.Background(), gofakeit.UUID()),
		&shortURLManagementMessages.GetShortURLByOrganizationIDRequest{
			OrganizationId: gofakeit.UUID(),
			Domain:         gofakeit.DomainName(),
			Alias:          "fake-alias",
		},
	)

	assert.Nil(t, response)
	assert.ErrorIs(t, err, errors.ErrShortURLNotFound)
}

func TestGetShortURL_Integration(t *testing.T) {
	t.Parallel()
	integration.MaybeSkipIntegrationTest(t)

	server, teardown := testUtils.BootstrapServer()

	t.Cleanup(teardown)

	shorturl := testUtils.AssertCreateShortURL(t, server, testUtils.CreateShortURLRequest())

	response, err := server.ShortURLManagementServiceClient.GetShortURL(
		grpcUtils.OutgoingContextWithUserID(context.Background(), shorturl.CreatedBy.Id),
		&shortURLManagementMessages.GetShortURLRequest{Domain: shorturl.Domain, Alias: shorturl.Alias},
	)

	require.NoError(t, err)

	testUtils.AssertShortURLEqual(t, shorturl, response.ShortUrl)
}

func TestGetShortURLWhenShortURLNotExists_Integration(t *testing.T) {
	t.Parallel()
	integration.MaybeSkipIntegrationTest(t)

	server, teardown := testUtils.BootstrapServer()

	t.Cleanup(teardown)

	response, err := server.ShortURLManagementServiceClient.GetShortURL(
		grpcUtils.OutgoingContextWithUserID(context.Background(), gofakeit.UUID()),
		&shortURLManagementMessages.GetShortURLRequest{Domain: gofakeit.DomainName(), Alias: "fake-alias"},
	)

	assert.Nil(t, response)
	assert.ErrorIs(t, err, errors.ErrShortURLNotFound)
}

func TestListShortURLByOrganizationID_Integration(t *testing.T) {
	t.Parallel()
	integration.MaybeSkipIntegrationTest(t)

	server, teardown := testUtils.BootstrapServer()

	t.Cleanup(teardown)

	shorturls := map[string]*shortURLManagementMessages.ShortURL{}
	organizationID := gofakeit.UUID()
	shorturlsCount, pageCount, pageSize := 3, int32(2), int32(2)

	for range shorturlsCount {
		request := testUtils.CreateShortURLRequest()
		request.OrganizationId = organizationID
		shorturl := testUtils.AssertCreateShortURL(t, server, request)
		shorturls[shorturl.Alias] = shorturl
	}

	actualShorturls := map[string]*shortURLManagementMessages.ShortURL{}
	for currPage := range pageCount {
		response, err := server.ShortURLManagementServiceClient.ListShortURLByOrganizationID(
			grpcUtils.OutgoingContextWithUserID(context.Background(), gofakeit.UUID()),
			&shortURLManagementMessages.ListShortURLByOrganizationIDRequest{
				OrganizationId: organizationID,
				PageSize:       pageSize,
				PageNum:        currPage + 1,
			},
		)

		require.NoError(t, err)

		assert.Equal(t, int32(shorturlsCount), response.Total)
		for _, actualShorturl := range response.Data {
			actualShorturls[actualShorturl.Alias] = actualShorturl
		}
	}

	require.Equal(t, len(shorturls), len(actualShorturls))

	for alias, shorturl := range actualShorturls {
		testUtils.AssertShortURLEqual(t, shorturls[alias], shorturl)
	}
}

func TestUpdateShortURLByOrganizationID_Integration(t *testing.T) {
	t.Parallel()
	integration.MaybeSkipIntegrationTest(t)

	server, teardown := testUtils.BootstrapServer()

	t.Cleanup(teardown)

	shorturl := testUtils.AssertCreateShortURL(t, server, testUtils.CreateShortURLRequest())
	request := &shortURLManagementMessages.UpdateShortURLByOrganizationIDRequest{
		OrganizationId: shorturl.OrganizationId,
		Domain:         shorturl.Domain,
		Alias:          shorturl.Alias,
		ShortUrl: &shortURLManagementMessages.UpdateShortURLParams{
			Alias:       wrapperspb.String(testUtils.RandomShortURLAlias()),
			LongUrl:     gofakeit.URL(),
			Description: "updated description",
			ExpiresAt:   timestamppb.New(time.Now().Add(time.Hour)),
			Tags:        []string{"new-tag-1"},
		},
		UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"alias", "long_url", "description", "expires_at", "tags"}},
	}

	response, err := server.ShortURLManagementServiceClient.UpdateShortURLByOrganizationID(
		grpcUtils.OutgoingContextWithUserID(context.Background(), shorturl.CreatedBy.Id),
		request,
	)

	require.NoError(t, err)

	assert.Equal(t, request.ShortUrl.Alias.GetValue(), response.ShortUrl.Alias)
	assert.Equal(t, request.ShortUrl.LongUrl, response.ShortUrl.LongUrl.Assembled)
	assert.Equal(t, request.ShortUrl.Description, response.ShortUrl.Description)
	testUtils.AssertTimestampNearlyEqual(t, request.ShortUrl.ExpiresAt, response.ShortUrl.ExpiresAt)
	assert.Equal(t, request.ShortUrl.Tags, response.ShortUrl.Tags)
}

func TestUpdateShortURLByOrganizationIDWhenShortURLNotFound_Integration(t *testing.T) {
	t.Parallel()
	integration.MaybeSkipIntegrationTest(t)

	server, teardown := testUtils.BootstrapServer()

	request := &shortURLManagementMessages.UpdateShortURLByOrganizationIDRequest{
		OrganizationId: gofakeit.UUID(),
		Domain:         gofakeit.DomainName(),
		Alias:          testUtils.RandomShortURLAlias(),
		ShortUrl:       &shortURLManagementMessages.UpdateShortURLParams{Description: gofakeit.ProductDescription()},
		UpdateMask:     &fieldmaskpb.FieldMask{Paths: []string{"description"}},
	}

	response, err := server.ShortURLManagementServiceClient.UpdateShortURLByOrganizationID(
		grpcUtils.OutgoingContextWithUserID(context.Background(), gofakeit.UUID()),
		request,
	)

	require.ErrorIs(t, err, errors.ErrShortURLNotFound)

	assert.Nil(t, response)

	t.Cleanup(teardown)
}

func TestUpdateShortURLByOrganizationIDWhenInvalidRequest_Integration(t *testing.T) {
	t.Parallel()
	integration.MaybeSkipIntegrationTest(t)

	server, teardown := testUtils.BootstrapServer()

	t.Cleanup(teardown)

	shortURLForUpdate := testUtils.AssertCreateShortURL(t, server, testUtils.CreateShortURLRequest())

	tests := []struct {
		name            string
		alias           string
		fieldMask       []string
		params          *shortURLManagementMessages.UpdateShortURLParams
		fieldViolations map[string][]string
	}{
		{
			name:            "invalid field in field mask",
			fieldMask:       []string{"not-exist"},
			params:          &shortURLManagementMessages.UpdateShortURLParams{},
			fieldViolations: map[string][]string{"update_mask": {"must contain only valid fields"}},
		},
		{
			name:            "invalid alias (empty)",
			fieldMask:       []string{"alias"},
			params:          &shortURLManagementMessages.UpdateShortURLParams{},
			fieldViolations: map[string][]string{"short_url.alias": {"alias can't be empty"}},
		},
		{
			name:            "invalid alias (too short)",
			fieldMask:       []string{"alias"},
			params:          &shortURLManagementMessages.UpdateShortURLParams{Alias: wrapperspb.String("a")},
			fieldViolations: map[string][]string{"short_url.alias": {"value length must be at least 4"}},
		},
		{
			name:      "invalid alias (too long)",
			fieldMask: []string{"alias"},
			params: &shortURLManagementMessages.UpdateShortURLParams{
				Alias: wrapperspb.String(strings.Repeat("a", 31)),
			},
			fieldViolations: map[string][]string{"short_url.alias": {"value length must be at most 30"}},
		},
		{
			name:      "invalid alias (not base62)",
			fieldMask: []string{"alias"},
			params: &shortURLManagementMessages.UpdateShortURLParams{
				Alias: wrapperspb.String("n{}t_b@se_62"),
			},
			fieldViolations: map[string][]string{"short_url.alias": {"value does not match regex pattern"}},
		},
		{
			name:            "invalid long url (not a valid URI)",
			fieldMask:       []string{"long_url"},
			params:          &shortURLManagementMessages.UpdateShortURLParams{LongUrl: "invalid"},
			fieldViolations: map[string][]string{"short_url.long_url": {"invalid URI"}},
		},
		{
			name:            "invalid description (too long)",
			fieldMask:       []string{"description"},
			params:          &shortURLManagementMessages.UpdateShortURLParams{Description: strings.Repeat("a", 1001)},
			fieldViolations: map[string][]string{"short_url.description": {"value length must be at most 1000"}},
		},
		{
			name:      "invalid timestamp: value in past",
			fieldMask: []string{"expires_at"},
			params: &shortURLManagementMessages.UpdateShortURLParams{
				ExpiresAt: timestamppb.New(time.Now().Add(-time.Hour)),
			},
			fieldViolations: map[string][]string{"short_url.expires_at": {"value must be greater than now"}},
		},
		{
			name:            "invalid tags (duplicate values)",
			fieldMask:       []string{"tags"},
			params:          &shortURLManagementMessages.UpdateShortURLParams{Tags: []string{"a", "a"}},
			fieldViolations: map[string][]string{"short_url.tags": {"repeated value must contain unique items"}},
		},
		{
			name:            "invalid tags (too short)",
			fieldMask:       []string{"tags"},
			params:          &shortURLManagementMessages.UpdateShortURLParams{Tags: []string{""}},
			fieldViolations: map[string][]string{"short_url.tags": {"value length must be at least 1"}},
		},
		{
			name:            "invalid tags (too long)",
			fieldMask:       []string{"tags"},
			params:          &shortURLManagementMessages.UpdateShortURLParams{Tags: []string{strings.Repeat("a", 51)}},
			fieldViolations: map[string][]string{"short_url.tags": {"value length must be at most 50"}},
		},
		{
			name:      "invalid tags (too much tags)",
			fieldMask: []string{"tags"},
			params: &shortURLManagementMessages.UpdateShortURLParams{
				Tags: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11"},
			},
			fieldViolations: map[string][]string{"short_url.tags": {"value must contain no more than 10"}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			alias := test.alias
			if alias == "" {
				alias = shortURLForUpdate.Alias
			}

			request := &shortURLManagementMessages.UpdateShortURLByOrganizationIDRequest{
				OrganizationId: shortURLForUpdate.OrganizationId,
				Domain:         shortURLForUpdate.Domain,
				Alias:          alias,
				ShortUrl:       test.params,
				UpdateMask:     &fieldmaskpb.FieldMask{Paths: test.fieldMask},
			}

			response, err := server.ShortURLManagementServiceClient.UpdateShortURLByOrganizationID(
				grpcUtils.OutgoingContextWithUserID(context.Background(), gofakeit.UUID()),
				request,
			)

			asserts.AssertValidationErrorContainsFieldViolations(t, err, test.fieldViolations)

			assert.Nil(t, response)
		})
	}
}

func TestUpdateShortURLByOrganizationIDWhenDuplicatedLongURL_Integration(t *testing.T) {
	t.Parallel()
	integration.MaybeSkipIntegrationTest(t)

	server, teardown := testUtils.BootstrapServer()

	t.Cleanup(teardown)

	createShortURLForUpdateRequst := testUtils.CreateShortURLRequest()
	shortURLForUpdate := testUtils.AssertCreateShortURL(t, server, createShortURLForUpdateRequst)

	createConflictingShorturlRequest := testUtils.CreateShortURLRequest()
	createConflictingShorturlRequest.OrganizationId = shortURLForUpdate.OrganizationId
	conflictingShorturl := testUtils.AssertCreateShortURL(t, server, createConflictingShorturlRequest)

	request := &shortURLManagementMessages.UpdateShortURLByOrganizationIDRequest{
		OrganizationId: shortURLForUpdate.OrganizationId,
		Domain:         shortURLForUpdate.Domain,
		Alias:          shortURLForUpdate.Alias,
		ShortUrl: &shortURLManagementMessages.UpdateShortURLParams{
			LongUrl: conflictingShorturl.LongUrl.Assembled,
		},
		UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"long_url"}},
	}

	response, err := server.ShortURLManagementServiceClient.UpdateShortURLByOrganizationID(
		grpcUtils.OutgoingContextWithUserID(context.Background(), gofakeit.UUID()),
		request,
	)

	require.ErrorIs(t, err, errors.ErrShortURLAlreadyExists)
	assert.Nil(t, response)
}

func TestDeleteShortURLByOrganizationID_Integration(t *testing.T) {
	t.Parallel()
	integration.MaybeSkipIntegrationTest(t)

	server, teardown := testUtils.BootstrapServer()

	t.Cleanup(teardown)

	shorturl := testUtils.AssertCreateShortURL(t, server, testUtils.CreateShortURLRequest())
	deleteShortURLRequest := &shortURLManagementMessages.DeleteShortURLByOrganizationIDRequest{
		OrganizationId: shorturl.OrganizationId,
		Domain:         shorturl.Domain,
		Alias:          shorturl.Alias,
	}

	response, err := server.ShortURLManagementServiceClient.DeleteShortURLByOrganizationID(
		grpcUtils.OutgoingContextWithUserID(context.Background(), shorturl.CreatedBy.Id),
		deleteShortURLRequest,
	)

	require.NoError(t, err)
	require.NotNil(t, response)

	response, err = server.ShortURLManagementServiceClient.DeleteShortURLByOrganizationID(
		grpcUtils.OutgoingContextWithUserID(context.Background(), shorturl.CreatedBy.Id),
		deleteShortURLRequest,
	)

	require.ErrorIs(t, err, errors.ErrShortURLNotFound)
	assert.Nil(t, response)
}

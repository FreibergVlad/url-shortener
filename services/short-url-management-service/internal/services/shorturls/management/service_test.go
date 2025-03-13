package management_test

import (
	"context"
	"testing"
	"time"

	shorturlproto "github.com/FreibergVlad/url-shortener/proto/pkg/shorturls/management/messages/v1"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/clock"
	serviceerrors "github.com/FreibergVlad/url-shortener/shared/go/pkg/errors"
	shorturlrepository "github.com/FreibergVlad/url-shortener/short-url-management-service/internal/db/repositories/shorturls"
	"github.com/FreibergVlad/url-shortener/short-url-management-service/internal/db/schema"
	shorturlservice "github.com/FreibergVlad/url-shortener/short-url-management-service/internal/services/shorturls"
	shorturlmanagementservice "github.com/FreibergVlad/url-shortener/short-url-management-service/internal/services/shorturls/management"
	testutils "github.com/FreibergVlad/url-shortener/short-url-management-service/internal/testing"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestGetShortURLByOrganizationID(t *testing.T) {
	t.Parallel()

	fakeDatabaseErr := gofakeit.ErrorDatabase()
	tests := []struct {
		name       string
		repoErr    error
		serviceErr error
	}{
		{"success", nil, nil},
		{"short URL not found", shorturlrepository.ErrNotFound, serviceerrors.ErrShortURLNotFound},
		{"unknown DB error", fakeDatabaseErr, fakeDatabaseErr},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			shorturl := testutils.RandomShortURL()
			request := &shorturlproto.GetShortURLByOrganizationIDRequest{
				OrganizationId: shorturl.OrganizationID,
				Domain:         shorturl.Domain,
				Alias:          shorturl.Alias,
			}
			repository := &testutils.MockedShortURLRepository{}
			service := shorturlmanagementservice.New(repository, clock.NewFixedClock(time.Now()))
			ctx := context.Background()

			repository.
				On("GetByOrganizationIDAndGlobalKey", ctx, request.OrganizationId, shorturl.GlobalKey()).
				Return(shorturl, test.repoErr)

			response, err := service.GetShortURLByOrganizationID(ctx, request)

			if test.serviceErr != nil {
				require.ErrorIs(t, err, test.serviceErr)
				require.Nil(t, response)
			} else {
				require.NoError(t, err)
				testutils.AssertShortURLEqual(t, shorturlservice.ShortURLToResponse(shorturl), response.ShortUrl)
			}

			repository.AssertExpectations(t)
		})
	}
}

func TestGetShortURL(t *testing.T) {
	t.Parallel()

	fakeDatabaseErr := gofakeit.ErrorDatabase()
	tests := []struct {
		name       string
		repoErr    error
		serviceErr error
	}{
		{"success", nil, nil},
		{"short URL not found", shorturlrepository.ErrNotFound, serviceerrors.ErrShortURLNotFound},
		{"unknown DB error", fakeDatabaseErr, fakeDatabaseErr},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			shorturl := testutils.RandomShortURL()
			request := &shorturlproto.GetShortURLRequest{Domain: shorturl.Domain, Alias: shorturl.Alias}
			repository := &testutils.MockedShortURLRepository{}
			service := shorturlmanagementservice.New(repository, clock.NewFixedClock(time.Now()))
			ctx := context.Background()

			repository.
				On("GetByGlobalKey", ctx, shorturl.GlobalKey()).
				Return(shorturl, test.repoErr)

			response, err := service.GetShortURL(ctx, request)

			if test.serviceErr != nil {
				require.ErrorIs(t, err, test.serviceErr)
				require.Nil(t, response)
			} else {
				require.NoError(t, err)
				testutils.AssertShortURLEqual(t, shorturlservice.ShortURLToResponse(shorturl), response.ShortUrl)
			}

			repository.AssertExpectations(t)
		})
	}
}

func TestListShortURLByOrganizationID(t *testing.T) {
	t.Parallel()

	fakeDatabaseErr := gofakeit.ErrorDatabase()
	tests := []struct {
		name       string
		repoErr    error
		serviceErr error
	}{{"success", nil, nil}, {"unknown DB error", fakeDatabaseErr, fakeDatabaseErr}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			total := int32(5)
			shorturls := make([]*schema.ShortURL, total)
			for i := range total {
				shorturls[i] = testutils.RandomShortURL()
			}

			request := &shorturlproto.ListShortURLByOrganizationIDRequest{
				OrganizationId: gofakeit.UUID(),
				PageSize:       5,
				PageNum:        1,
			}
			repository := &testutils.MockedShortURLRepository{}
			service := shorturlmanagementservice.New(repository, clock.NewFixedClock(time.Now()))
			ctx := context.Background()

			repository.
				On("ListByOrganizationID", ctx, request.OrganizationId, int(request.PageNum), int(request.PageSize)).
				Return(&shorturlrepository.PaginatedResult{Data: shorturls, Total: total}, test.repoErr)

			response, err := service.ListShortURLByOrganizationID(ctx, request)

			if test.serviceErr != nil {
				require.ErrorIs(t, err, test.serviceErr)
				require.Nil(t, response)
			} else {
				require.NoError(t, err)
				require.Equal(t, total, response.Total)
				for i, shorturl := range shorturls {
					testutils.AssertShortURLEqual(t, shorturlservice.ShortURLToResponse(shorturl), response.Data[i])
				}
			}

			repository.AssertExpectations(t)
		})
	}
}

func TestUpdateShortURLByOrganizationID(t *testing.T) {
	t.Parallel()

	fakeDatabaseErr := gofakeit.ErrorDatabase()
	tests := []struct {
		name                   string
		getShortURLRepoErr     error
		replaceShortURLRepoErr error
		serviceErr             error
	}{
		{"success", nil, nil, nil},
		{"unknown DB error when getting short URL", fakeDatabaseErr, nil, fakeDatabaseErr},
		{"short URL not found", shorturlrepository.ErrNotFound, nil, serviceerrors.ErrShortURLNotFound},
		{"unknown DB error when replacing short URL", nil, fakeDatabaseErr, fakeDatabaseErr},
		{
			"unique constraint violation when replacing short URL",
			nil,
			shorturlrepository.ErrAlreadyExists,
			serviceerrors.ErrShortURLAlreadyExists,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			shorturl := testutils.RandomShortURL()
			request := &shorturlproto.UpdateShortURLByOrganizationIDRequest{
				OrganizationId: gofakeit.UUID(),
				Domain:         shorturl.Domain,
				Alias:          shorturl.Alias,
				ShortUrl: &shorturlproto.UpdateShortURLParams{
					Alias: wrapperspb.String(testutils.RandomShortURLAlias()),
				},
				UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"alias"}},
			}
			repository := &testutils.MockedShortURLRepository{}
			service := shorturlmanagementservice.New(repository, clock.NewFixedClock(time.Now()))
			ctx := context.Background()

			repository.
				On("GetByOrganizationIDAndGlobalKey", ctx, request.OrganizationId, shorturl.GlobalKey()).
				Return(shorturl, test.getShortURLRepoErr)

			if test.getShortURLRepoErr == nil {
				repository.
					On(
						"ReplaceByOrganizationIDAndGlobalKey",
						ctx,
						request.OrganizationId,
						shorturl.GlobalKey(),
						shorturl,
					).
					Return(test.replaceShortURLRepoErr)
			}

			response, err := service.UpdateShortURLByOrganizationID(ctx, request)

			if test.serviceErr != nil {
				require.ErrorIs(t, err, test.serviceErr)
				require.Nil(t, response)
			} else {
				require.NoError(t, err)
				require.Equal(t, request.ShortUrl.Alias.Value, shorturl.Alias)
				testutils.AssertShortURLEqual(t, shorturlservice.ShortURLToResponse(shorturl), response.ShortUrl)
			}

			repository.AssertExpectations(t)
		})
	}
}

func TestDeleteShortURLByOrganizationID(t *testing.T) {
	t.Parallel()

	fakeDatabaseErr := gofakeit.ErrorDatabase()
	tests := []struct {
		name       string
		repoErr    error
		serviceErr error
	}{
		{"success", nil, nil},
		{"short URL not found", shorturlrepository.ErrNotFound, serviceerrors.ErrShortURLNotFound},
		{"unknown DB error", fakeDatabaseErr, fakeDatabaseErr},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			shorturl := testutils.RandomShortURL()
			request := &shorturlproto.DeleteShortURLByOrganizationIDRequest{
				OrganizationId: shorturl.OrganizationID,
				Domain:         shorturl.Domain,
				Alias:          shorturl.Alias,
			}
			repository := &testutils.MockedShortURLRepository{}
			service := shorturlmanagementservice.New(repository, clock.NewFixedClock(time.Now()))
			ctx := context.Background()

			repository.
				On("DeleteByOrganizationIDAndGlobalKey", ctx, request.OrganizationId, shorturl.GlobalKey()).
				Return(shorturl, test.repoErr)

			response, err := service.DeleteShortURLByOrganizationID(ctx, request)

			if test.serviceErr != nil {
				require.ErrorIs(t, err, test.serviceErr)
				require.Nil(t, response)
			} else {
				require.NoError(t, err)
				require.Equal(t, &shorturlproto.DeleteShortURLByOrganizationIDResponse{}, response)
			}

			repository.AssertExpectations(t)
		})
	}
}

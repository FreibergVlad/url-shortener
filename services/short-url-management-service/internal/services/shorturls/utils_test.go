package shorturls_test

import (
	"fmt"
	"testing"
	"time"

	shortURLGeneratorServiceMessages "github.com/FreibergVlad/url-shortener/proto/pkg/shorturls/generator/messages/v1"
	shortURLManagementServiceMessages "github.com/FreibergVlad/url-shortener/proto/pkg/shorturls/management/messages/v1"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/clock"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/testing/asserts"
	"github.com/FreibergVlad/url-shortener/short-url-management-service/internal/db/schema"
	"github.com/FreibergVlad/url-shortener/short-url-management-service/internal/services/shorturls"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestShortURLFromRequest(t *testing.T) {
	t.Parallel()

	longURL := "https://fake-long-url.com/path?q1=1&q2=2"
	createdAt, expiresAt := time.Now().UTC(), time.Now().UTC()

	input := &shortURLGeneratorServiceMessages.CreateShortURLRequest{
		OrganizationId: gofakeit.UUID(),
		Domain:         gofakeit.DomainName(),
		LongUrl:        longURL,
		Tags:           []string{"tag1", "tag2"},
		ExpiresAt:      timestamppb.New(expiresAt),
		Description:    gofakeit.ProductDescription(),
	}
	expected := &schema.ShortURL{
		OrganizationID: input.OrganizationId,
		LongURL: schema.LongURL{
			Hash:      "2c4b0e9bd811502f90d68b4a0364c91c84ef18e2473160455b5268db0cfbd4ba",
			Assembled: longURL,
			Scheme:    "https",
			Host:      "fake-long-url.com",
			Path:      "/path",
			Query:     "q1=1&q2=2",
		},
		Scheme:      "https",
		Domain:      input.Domain,
		CreatedAt:   createdAt,
		ExpiresAt:   &expiresAt,
		CreatedBy:   schema.User{ID: gofakeit.UUID(), Email: gofakeit.Email()},
		Description: input.Description,
		Tags:        input.Tags,
	}

	actual := shorturls.ShortURLFromRequest(input, clock.NewFixedClock(createdAt), "https", expected.CreatedBy)

	assert.Equal(t, expected, actual)
}

func TestShortToResponse(t *testing.T) {
	t.Parallel()

	input := &schema.ShortURL{
		OrganizationID: gofakeit.UUID(),
		LongURL: schema.LongURL{
			Assembled: gofakeit.URL(),
			Hash:      gofakeit.UUID(),
			Scheme:    gofakeit.RandomString([]string{"http", "https"}),
			Host:      gofakeit.DomainName(),
		},
		Scheme:      "https",
		Domain:      gofakeit.DomainName(),
		Alias:       "alias",
		CreatedAt:   time.Now().UTC(),
		CreatedBy:   schema.User{ID: gofakeit.UUID(), Email: gofakeit.Email()},
		Description: gofakeit.ProductDescription(),
		Tags:        []string{"tag1", "tag2"},
	}
	expected := &shortURLManagementServiceMessages.ShortURL{
		OrganizationId: input.OrganizationID,
		ShortUrl:       fmt.Sprintf("https://%s/%s", input.Domain, input.Alias),
		LongUrl: &shortURLManagementServiceMessages.LongURL{
			Assembled: input.LongURL.Assembled,
			Hash:      input.LongURL.Hash,
			Scheme:    input.LongURL.Scheme,
			Host:      input.LongURL.Host,
		},
		Domain:    input.Domain,
		Alias:     input.Alias,
		CreatedAt: timestamppb.New(input.CreatedAt),
		CreatedBy: &shortURLManagementServiceMessages.User{
			Id:    input.CreatedBy.ID,
			Email: input.CreatedBy.Email,
		},
		Tags:        input.Tags,
		Description: input.Description,
	}

	actual := shorturls.ShortURLToResponse(input)

	assert.Equal(t, expected, actual)
}

func TestUpdateShortURLFromRequest(t *testing.T) {
	t.Parallel()

	expiresAt := time.Now().UTC()

	tests := []struct {
		name             string
		originalShortURL *schema.ShortURL
		expectedShortURL *schema.ShortURL
		fieldMask        []string
		updateParams     *shortURLManagementServiceMessages.UpdateShortURLParams
		fieldViolations  map[string][]string
	}{
		{
			name:             "update alias",
			originalShortURL: &schema.ShortURL{Alias: "old-alias"},
			expectedShortURL: &schema.ShortURL{Alias: "new-alias"},
			fieldMask:        []string{"alias"},
			updateParams: &shortURLManagementServiceMessages.UpdateShortURLParams{
				Alias: wrapperspb.String("new-alias"),
			},
		},
		{
			name:             "update alias when empty",
			originalShortURL: &schema.ShortURL{},
			fieldMask:        []string{"alias"},
			updateParams:     &shortURLManagementServiceMessages.UpdateShortURLParams{Alias: wrapperspb.String("")},
			fieldViolations:  map[string][]string{"short_url.alias": {"alias can't be empty"}},
		},
		{
			name:             "update long URL",
			originalShortURL: &schema.ShortURL{},
			expectedShortURL: &schema.ShortURL{
				LongURL: schema.LongURL{
					Hash:      "fb100732afd5cb53ef0d8b649bab0e6f750d72bc931b38198d83d42ec0242267",
					Assembled: "https://fake-url.com/path?k1=v1&k2=v2",
					Scheme:    "https",
					Host:      "fake-url.com",
					Path:      "/path",
					Query:     "k1=v1&k2=v2",
				},
			},
			fieldMask: []string{"long_url"},
			updateParams: &shortURLManagementServiceMessages.UpdateShortURLParams{
				LongUrl: "https://fake-url.com/path?k1=v1&k2=v2",
			},
		},
		{
			name:             "update long URL when empty",
			originalShortURL: &schema.ShortURL{},
			fieldMask:        []string{"long_url"},
			updateParams:     &shortURLManagementServiceMessages.UpdateShortURLParams{LongUrl: ""},
			fieldViolations:  map[string][]string{"short_url.long_url": {"invalid URL"}},
		},
		{
			name:             "update description",
			originalShortURL: &schema.ShortURL{Description: "old description"},
			expectedShortURL: &schema.ShortURL{Description: "new description"},
			fieldMask:        []string{"description"},
			updateParams: &shortURLManagementServiceMessages.UpdateShortURLParams{
				Description: "new description",
			},
		},
		{
			name:             "set expiration date",
			originalShortURL: &schema.ShortURL{},
			expectedShortURL: &schema.ShortURL{ExpiresAt: &expiresAt},
			fieldMask:        []string{"expires_at"},
			updateParams: &shortURLManagementServiceMessages.UpdateShortURLParams{
				ExpiresAt: timestamppb.New(expiresAt),
			},
		},
		{
			name:             "reset expiration date",
			originalShortURL: &schema.ShortURL{ExpiresAt: &expiresAt},
			expectedShortURL: &schema.ShortURL{},
			fieldMask:        []string{"expires_at"},
			updateParams:     &shortURLManagementServiceMessages.UpdateShortURLParams{},
		},
		{
			name:             "update tags",
			originalShortURL: &schema.ShortURL{Tags: []string{"old-tag"}},
			expectedShortURL: &schema.ShortURL{Tags: []string{"new-tag"}},
			fieldMask:        []string{"tags"},
			updateParams: &shortURLManagementServiceMessages.UpdateShortURLParams{
				Tags: []string{"new-tag"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			request := &shortURLManagementServiceMessages.UpdateShortURLByOrganizationIDRequest{
				UpdateMask: &fieldmaskpb.FieldMask{Paths: test.fieldMask},
				ShortUrl:   test.updateParams,
			}

			err := shorturls.UpdateShortURLFromRequest(test.originalShortURL, request)

			if test.fieldViolations != nil {
				asserts.AssertValidationErrorContainsFieldViolations(t, err, test.fieldViolations)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.expectedShortURL, test.originalShortURL)
			}
		})
	}
}

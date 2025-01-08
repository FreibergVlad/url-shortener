package services_test

import (
	"fmt"
	"testing"
	"time"

	urlGeneratorProtobufMessages "github.com/FreibergVlad/url-shortener/proto/pkg/shorturl/generator/messages/v1"
	urlManagementProtobufMessages "github.com/FreibergVlad/url-shortener/proto/pkg/shorturl/management/messages/v1"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/clock"
	"github.com/FreibergVlad/url-shortener/short-url-management-service/internal/db"
	services "github.com/FreibergVlad/url-shortener/short-url-management-service/internal/services/shorturl"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestShortURLFromRequest(t *testing.T) {
	t.Parallel()

	organizationId, createdBy := "fake-org-id", "fake-user"
	domain, longUrl := "test.com", "https://fake-long-url.com/path?q1=1&q2=2"
	expiresAtTimestamp := timestamppb.New(time.Now())
	createdAt, expiredAt := time.Now(), expiresAtTimestamp.AsTime()
	description := "fake-description"
	tags := []string{"tag1", "tag2"}

	input := urlGeneratorProtobufMessages.CreateShortURLRequest{
		OrganizationId: organizationId,
		Domain:         domain,
		LongUrl:        longUrl,
		Tags:           tags,
		ExpiresAt:      expiresAtTimestamp,
		Description:    wrapperspb.String(description),
	}
	expected := &db.ShortURL{
		OrganizationID: organizationId,
		LongURL: db.LongURL{
			Hash:      "2c4b0e9bd811502f90d68b4a0364c91c84ef18e2473160455b5268db0cfbd4ba",
			Assembled: longUrl,
			Scheme:    "https",
			Host:      "fake-long-url.com",
			Path:      "/path",
			Query:     "q1=1&q2=2",
		},
		Scheme:      "https",
		Domain:      domain,
		Status:      urlManagementProtobufMessages.ShortURLStatus_SHORT_URL_STATUS_ACTIVE.String(),
		CreatedAt:   createdAt,
		ExpiresAt:   &expiredAt,
		CreatedBy:   createdBy,
		Description: &description,
		Tags:        tags,
	}

	actual := services.ShortURLFromRequest(&input, clock.NewFixedClock(createdAt), createdBy, "https")
	assert.Equal(t, expected, actual)
}

func TestShortToResponse(t *testing.T) {
	t.Parallel()

	organizationId, createdBy := "fake-org-id", "fake-user"
	domain, alias := "test.com", "alias"
	shortUrl := fmt.Sprintf("https://%s/%s", domain, alias)
	createdAt := time.Now()
	createdAtTimestamp := timestamppb.New(createdAt)
	description, tags := "fake-description", []string{"tag1", "tag2"}

	input := &db.ShortURL{
		OrganizationID: organizationId,
		LongURL: db.LongURL{
			Assembled: "https://fake-long-url.com",
			Hash:      "fake-hash",
			Scheme:    "https",
			Host:      "fake-long-url.com",
			Path:      "",
			Query:     "",
		},
		Scheme:      "https",
		Domain:      domain,
		Alias:       alias,
		Status:      urlManagementProtobufMessages.ShortURLStatus_SHORT_URL_STATUS_ACTIVE.String(),
		CreatedAt:   createdAt,
		CreatedBy:   createdBy,
		Description: &description,
		Tags:        tags,
	}
	expected := &urlManagementProtobufMessages.ShortURL{
		OrganizationId: organizationId,
		ShortUrl:       shortUrl,
		LongUrl: &urlManagementProtobufMessages.LongURL{
			Assembled: "https://fake-long-url.com",
			Hash:      "fake-hash",
			Scheme:    "https",
			Host:      "fake-long-url.com",
			Path:      "",
			Query:     "",
		},
		Domain:      domain,
		Alias:       alias,
		Status:      urlManagementProtobufMessages.ShortURLStatus_SHORT_URL_STATUS_ACTIVE,
		CreatedAt:   createdAtTimestamp,
		CreatedBy:   createdBy,
		Tags:        tags,
		Description: wrapperspb.String(description),
	}

	actual := services.ShortURLToResponse(input)
	assert.Equal(t, actual, expected)
}

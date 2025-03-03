package shorturls

import (
	"crypto/sha256"
	"encoding/hex"
	"net/url"
	"slices"
	"time"

	shortURLGeneratorProtoMessages "github.com/FreibergVlad/url-shortener/proto/pkg/shorturls/generator/messages/v1"
	shortURLManagementProtoMessages "github.com/FreibergVlad/url-shortener/proto/pkg/shorturls/management/messages/v1"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/clock"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/errors"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/must"
	"github.com/FreibergVlad/url-shortener/short-url-management-service/internal/db/schema"
	"github.com/FreibergVlad/url-shortener/short-url-management-service/internal/validators"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ShortURLFromRequest(
	req *shortURLGeneratorProtoMessages.CreateShortURLRequest, clock clock.Clock, scheme string, user schema.User,
) *schema.ShortURL {
	shorturl := &schema.ShortURL{
		OrganizationID: req.OrganizationId,
		LongURL:        ParseLongURL(req.LongUrl),
		Scheme:         scheme,
		Domain:         req.Domain,
		CreatedAt:      clock.Now(),
		CreatedBy:      user,
		Description:    req.Description,
		Tags:           req.Tags,
	}
	if req.ExpiresAt != nil {
		expiresAt := req.ExpiresAt.AsTime()
		shorturl.ExpiresAt = &expiresAt
	}
	return shorturl
}

func ShortURLToResponse(shorturl *schema.ShortURL) *shortURLManagementProtoMessages.ShortURL {
	return &shortURLManagementProtoMessages.ShortURL{
		OrganizationId: shorturl.OrganizationID,
		ShortUrl:       assembleShortURL(shorturl),
		LongUrl: &shortURLManagementProtoMessages.LongURL{
			Hash:      shorturl.LongURL.Hash,
			Assembled: shorturl.LongURL.Assembled,
			Scheme:    shorturl.LongURL.Scheme,
			Host:      shorturl.LongURL.Host,
			Path:      shorturl.LongURL.Path,
			Query:     shorturl.LongURL.Query,
		},
		Domain:    shorturl.Domain,
		Alias:     shorturl.Alias,
		ExpiresAt: timestamppbFromTime(shorturl.ExpiresAt),
		CreatedAt: timestamppbFromTime(&shorturl.CreatedAt),
		CreatedBy: &shortURLManagementProtoMessages.User{
			Id:    shorturl.CreatedBy.ID,
			Email: shorturl.CreatedBy.Email,
		},
		Description: shorturl.Description,
		Tags:        shorturl.Tags,
	}
}

func ParseLongURL(longURL string) schema.LongURL {
	parsedLongURL := must.Return(url.ParseRequestURI(longURL))
	return schema.LongURL{
		Hash:      longURLHash(longURL),
		Assembled: longURL,
		Scheme:    parsedLongURL.Scheme,
		Host:      parsedLongURL.Host,
		Path:      parsedLongURL.EscapedPath(),
		Query:     parsedLongURL.RawQuery,
	}
}

func UpdateShortURLFromRequest(
	shorturl *schema.ShortURL, req *shortURLManagementProtoMessages.UpdateShortURLByOrganizationIDRequest,
) error {
	if slices.Contains(req.UpdateMask.Paths, "description") {
		shorturl.Description = req.ShortUrl.Description
	}
	if slices.Contains(req.UpdateMask.Paths, "tags") {
		shorturl.Tags = req.ShortUrl.Tags
	}
	if slices.Contains(req.UpdateMask.Paths, "expires_at") {
		if req.ShortUrl.ExpiresAt == nil {
			shorturl.ExpiresAt = nil
		} else {
			expiresAt := req.ShortUrl.ExpiresAt.AsTime()
			shorturl.ExpiresAt = &expiresAt
		}
	}
	if slices.Contains(req.UpdateMask.Paths, "long_url") {
		if err := validators.ValidateURL(req.ShortUrl.LongUrl); err != nil {
			return errors.NewValidationError(map[string][]string{"short_url.long_url": {err.Error()}})
		}
		shorturl.LongURL = ParseLongURL(req.ShortUrl.LongUrl)
	}
	if slices.Contains(req.UpdateMask.Paths, "alias") {
		if req.ShortUrl.Alias.GetValue() == "" {
			return errors.NewValidationError(map[string][]string{"short_url.alias": {"alias can't be empty"}})
		}
		shorturl.Alias = req.ShortUrl.Alias.GetValue()
	}
	return nil
}

func assembleShortURL(shorturl *schema.ShortURL) string {
	url := url.URL{Scheme: shorturl.Scheme, Host: shorturl.Domain, Path: shorturl.Alias}
	return url.String()
}

func longURLHash(longURL string) string {
	longURLHashRaw := sha256.Sum256([]byte(longURL))
	return hex.EncodeToString(longURLHashRaw[:])
}

func timestamppbFromTime(t *time.Time) *timestamppb.Timestamp {
	var timestamp *timestamppb.Timestamp
	if t != nil {
		timestamp = timestamppb.New(*t)
	}
	return timestamp
}

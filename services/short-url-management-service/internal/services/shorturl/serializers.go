package services

import (
	"crypto/sha256"
	"encoding/hex"
	"net/url"
	"time"

	shortUrlGeneratorProtoMessages "github.com/FreibergVlad/url-shortener/proto/pkg/shorturl/generator/messages/v1"
	shortUrlManagementProtoMessages "github.com/FreibergVlad/url-shortener/proto/pkg/shorturl/management/messages/v1"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/clock"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/must"
	"github.com/FreibergVlad/url-shortener/short-url-management-service/internal/db"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func ShortURLFromRequest(req *shortUrlGeneratorProtoMessages.CreateShortURLRequest, clock clock.Clock, scheme, userId string) *db.ShortURL {
	parsedLongUrl := must.Return(url.ParseRequestURI(req.LongUrl))
	shortUrl := &db.ShortURL{
		OrganizationID: req.OrganizationId,
		LongURL: db.LongURL{
			Hash:      longURLHash(req.LongUrl),
			Assembled: req.LongUrl,
			Scheme:    parsedLongUrl.Scheme,
			Host:      parsedLongUrl.Host,
			Path:      parsedLongUrl.EscapedPath(),
			Query:     parsedLongUrl.RawQuery,
		},
		Scheme:    scheme,
		Domain:    req.Domain,
		Status:    shortUrlManagementProtoMessages.ShortURLStatus_SHORT_URL_STATUS_ACTIVE.String(),
		CreatedAt: clock.Now(),
		CreatedBy: userId,
		Tags:      req.Tags,
	}
	if req.ExpiresAt != nil {
		expiresAt := req.ExpiresAt.AsTime()
		shortUrl.ExpiresAt = &expiresAt
	}
	if req.Description != nil {
		shortUrl.Description = &req.Description.Value
	}
	return shortUrl
}

func ShortURLToResponse(url *db.ShortURL) *shortUrlManagementProtoMessages.ShortURL {
	var description *wrapperspb.StringValue
	if url.Description != nil {
		description = &wrapperspb.StringValue{Value: *url.Description}
	}

	status := shortUrlManagementProtoMessages.ShortURLStatus_SHORT_URL_STATUS_UNSPECIFIED
	if _, ok := shortUrlManagementProtoMessages.ShortURLStatus_value[url.Status]; ok {
		status = shortUrlManagementProtoMessages.ShortURLStatus(shortUrlManagementProtoMessages.ShortURLStatus_value[url.Status])
	}

	return &shortUrlManagementProtoMessages.ShortURL{
		OrganizationId: url.OrganizationID,
		ShortUrl:       assembleShortURL(url),
		LongUrl: &shortUrlManagementProtoMessages.LongURL{
			Hash:      url.LongURL.Hash,
			Assembled: url.LongURL.Assembled,
			Scheme:    url.LongURL.Scheme,
			Host:      url.LongURL.Host,
			Path:      url.LongURL.Path,
			Query:     url.LongURL.Query,
		},
		Domain:      url.Domain,
		Alias:       url.Alias,
		Status:      status,
		ExpiresAt:   timestamppbFromTime(url.ExpiresAt),
		CreatedAt:   timestamppbFromTime(&url.CreatedAt),
		DeletedAt:   timestamppbFromTime(url.DeletedAt),
		CreatedBy:   url.CreatedBy,
		Description: description,
		Tags:        url.Tags,
	}
}

func assembleShortURL(shortURL *db.ShortURL) string {
	url := url.URL{Scheme: shortURL.Scheme, Host: shortURL.Domain, Path: shortURL.Alias}
	return url.String()
}

func longURLHash(longUrl string) string {
	longUrlHashRaw := sha256.Sum256([]byte(longUrl))
	return hex.EncodeToString(longUrlHashRaw[:])
}

func timestamppbFromTime(t *time.Time) *timestamppb.Timestamp {
	var timestamp *timestamppb.Timestamp
	if t != nil {
		timestamp = timestamppb.New(*t)
	}
	return timestamp
}

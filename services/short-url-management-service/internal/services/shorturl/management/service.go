package management

import (
	"context"
	"slices"

	protoMessages "github.com/FreibergVlad/url-shortener/proto/pkg/shorturl/management/messages/v1"
	protoService "github.com/FreibergVlad/url-shortener/proto/pkg/shorturl/management/service/v1"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/clock"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/db"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/db/sort"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/errors"
	schema "github.com/FreibergVlad/url-shortener/short-url-management-service/internal/db"

	"github.com/FreibergVlad/url-shortener/short-url-management-service/internal/cache"
	services "github.com/FreibergVlad/url-shortener/short-url-management-service/internal/services/shorturl"
)

type shortUrlManagementService struct {
	protoService.UnimplementedShortURLManagementServiceServer
	shortUrlRepository db.Repository[schema.ShortURL]
	cache              *cache.ShortURLCache
	clock              clock.Clock
}

func New(shortUrlRepository db.Repository[schema.ShortURL], cache *cache.ShortURLCache, clock clock.Clock) *shortUrlManagementService {
	return &shortUrlManagementService{shortUrlRepository: shortUrlRepository, cache: cache, clock: clock}
}

func (s *shortUrlManagementService) GetShortURLByOrganizationID(ctx context.Context, req *protoMessages.GetShortURLByOrganizationIDRequest) (*protoMessages.GetShortURLByOrganizationIDResponse, error) {
	shortUrl, err := s.cache.Get(ctx, req.Domain, req.Alias)
	if err == nil {
		if shortUrl.OrganizationID != req.OrganizationId {
			return nil, errors.ErrResourceNotFound
		}
		return &protoMessages.GetShortURLByOrganizationIDResponse{ShortUrl: services.ShortURLToResponse(shortUrl)}, nil
	}

	filter := services.OneNonDeletedShortURLByOrgIDFilter(req.OrganizationId, req.Domain, req.Alias)
	shortUrl, err = s.shortUrlRepository.FindOne(ctx, filter)
	if err != nil {
		return nil, err
	}

	go s.cache.Set(context.WithoutCancel(ctx), shortUrl)

	return &protoMessages.GetShortURLByOrganizationIDResponse{ShortUrl: services.ShortURLToResponse(shortUrl)}, nil
}

func (s *shortUrlManagementService) GetShortURL(ctx context.Context, req *protoMessages.GetShortURLRequest) (*protoMessages.GetShortURLResponse, error) {
	shortUrl, err := s.cache.Get(ctx, req.Domain, req.Alias)
	if err == nil {
		return &protoMessages.GetShortURLResponse{ShortUrl: services.ShortURLToResponse(shortUrl)}, nil
	}

	filter := services.OneNonDeletedShortURLFilter(req.Domain, req.Alias)
	shortUrl, err = s.shortUrlRepository.FindOne(ctx, filter)
	if err != nil {
		return nil, err
	}

	go s.cache.Set(context.WithoutCancel(ctx), shortUrl)

	return &protoMessages.GetShortURLResponse{ShortUrl: services.ShortURLToResponse(shortUrl)}, nil
}

func (s *shortUrlManagementService) ListShortURLByOrganizationID(ctx context.Context, req *protoMessages.ListShortURLByOrganizationIDRequest) (*protoMessages.ListShortURLByOrganizationIDResponse, error) {
	query := db.PaginatedQuery{
		Filter:   services.ManyNonDeletedShortURLByOrgIDFilter(req.OrganizationId),
		Sort:     &sort.Sort{Field: "created_at", Order: -1},
		PageSize: req.PageSize,
		PageNum:  req.PageNum,
	}

	result, err := s.shortUrlRepository.FindMany(ctx, query)
	if err != nil {
		return nil, err
	}

	shortUrls := make([]*protoMessages.ShortURL, len(result.Data))
	for i, shortUrl := range result.Data {
		shortUrls[i] = services.ShortURLToResponse(shortUrl)
	}

	return &protoMessages.ListShortURLByOrganizationIDResponse{Data: shortUrls, Total: result.Total}, nil
}

func (s *shortUrlManagementService) UpdateShortURLByOrganizationID(ctx context.Context, req *protoMessages.UpdateShortURLByOrganizationIDRequest) (*protoMessages.UpdateShortURLByOrganizationIDResponse, error) {
	fieldsForUpdate := make(map[string]any)
	if slices.Contains(req.UpdateMask.Paths, "description") {
		if req.Description == nil {
			fieldsForUpdate["description"] = nil
		} else {
			fieldsForUpdate["description"] = req.Description.Value
		}
	}
	if slices.Contains(req.UpdateMask.Paths, "tags") {
		fieldsForUpdate["tags"] = req.Tags
	}
	if slices.Contains(req.UpdateMask.Paths, "expires_at") {
		if req.ExpiresAt == nil {
			fieldsForUpdate["expires_at"] = nil
		} else {
			fieldsForUpdate["expires_at"] = req.ExpiresAt.AsTime()
		}
	}

	filter := services.OneNonDeletedShortURLByOrgIDFilter(req.OrganizationId, req.Domain, req.Alias)
	shortUrl, err := s.shortUrlRepository.UpdateOne(ctx, filter, fieldsForUpdate)
	if err != nil {
		return nil, err
	}

	go s.cache.Set(context.WithoutCancel(ctx), shortUrl)

	return &protoMessages.UpdateShortURLByOrganizationIDResponse{ShortUrl: services.ShortURLToResponse(shortUrl)}, nil
}

func (s *shortUrlManagementService) DeleteShortURLByOrganizationID(ctx context.Context, req *protoMessages.DeleteShortURLByOrganizationIDRequest) (*protoMessages.DeleteShortURLByOrganizationIDResponse, error) {
	fieldsForUpdate := map[string]any{
		"status":     protoMessages.ShortURLStatus_SHORT_URL_STATUS_DELETED.String(),
		"deleted_at": s.clock.Now(),
	}
	filter := services.OneNonDeletedShortURLByOrgIDFilter(req.OrganizationId, req.Domain, req.Alias)
	shortUrl, err := s.shortUrlRepository.UpdateOne(ctx, filter, fieldsForUpdate)
	if err != nil {
		return nil, err
	}

	go s.cache.Delete(context.WithoutCancel(ctx), shortUrl)

	return &protoMessages.DeleteShortURLByOrganizationIDResponse{}, nil
}

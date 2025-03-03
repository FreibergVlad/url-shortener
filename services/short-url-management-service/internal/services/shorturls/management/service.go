package management

import (
	"context"
	"errors"
	"fmt"

	protoMessages "github.com/FreibergVlad/url-shortener/proto/pkg/shorturls/management/messages/v1"
	protoService "github.com/FreibergVlad/url-shortener/proto/pkg/shorturls/management/service/v1"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/clock"
	serviceErrors "github.com/FreibergVlad/url-shortener/shared/go/pkg/errors"
	shortURLRepository "github.com/FreibergVlad/url-shortener/short-url-management-service/internal/db/repositories/shorturls"
	"github.com/FreibergVlad/url-shortener/short-url-management-service/internal/db/schema"
	shortURLUtils "github.com/FreibergVlad/url-shortener/short-url-management-service/internal/services/shorturls"
)

type ShortURLManagementService struct {
	protoService.UnimplementedShortURLManagementServiceServer
	shortURLRepository shortURLRepository.Repository
	clock              clock.Clock
}

func New(shortURLRepository shortURLRepository.Repository, clock clock.Clock) *ShortURLManagementService {
	return &ShortURLManagementService{shortURLRepository: shortURLRepository, clock: clock}
}

func (s *ShortURLManagementService) GetShortURLByOrganizationID(
	ctx context.Context, req *protoMessages.GetShortURLByOrganizationIDRequest,
) (*protoMessages.GetShortURLByOrganizationIDResponse, error) {
	key := schema.ShortURLGlobalKey{Domain: req.Domain, Alias: req.Alias}
	shorturl, err := s.shortURLRepository.GetByOrganizationIDAndGlobalKey(ctx, req.OrganizationId, key)
	if err != nil {
		if errors.Is(err, shortURLRepository.ErrNotFound) {
			return nil, serviceErrors.ErrShortURLNotFound
		}
		return nil, fmt.Errorf("error getting short url: %w", err)
	}

	return &protoMessages.GetShortURLByOrganizationIDResponse{ShortUrl: shortURLUtils.ShortURLToResponse(shorturl)}, nil
}

func (s *ShortURLManagementService) GetShortURL(
	ctx context.Context, req *protoMessages.GetShortURLRequest,
) (*protoMessages.GetShortURLResponse, error) {
	key := schema.ShortURLGlobalKey{Domain: req.Domain, Alias: req.Alias}
	shorturl, err := s.shortURLRepository.GetByGlobalKey(ctx, key)
	if err != nil {
		if errors.Is(err, shortURLRepository.ErrNotFound) {
			return nil, serviceErrors.ErrShortURLNotFound
		}
		return nil, fmt.Errorf("error getting short url: %w", err)
	}

	return &protoMessages.GetShortURLResponse{ShortUrl: shortURLUtils.ShortURLToResponse(shorturl)}, nil
}

func (s *ShortURLManagementService) ListShortURLByOrganizationID(
	ctx context.Context, req *protoMessages.ListShortURLByOrganizationIDRequest,
) (*protoMessages.ListShortURLByOrganizationIDResponse, error) {
	result, err := s.shortURLRepository.ListByOrganizationID(
		ctx, req.OrganizationId, int(req.PageNum), int(req.PageSize),
	)
	if err != nil {
		return nil, fmt.Errorf("error listing short url: %w", err)
	}

	shorturls := make([]*protoMessages.ShortURL, len(result.Data))
	for i, shorturl := range result.Data {
		shorturls[i] = shortURLUtils.ShortURLToResponse(shorturl)
	}

	return &protoMessages.ListShortURLByOrganizationIDResponse{Data: shorturls, Total: result.Total}, nil
}

func (s *ShortURLManagementService) UpdateShortURLByOrganizationID(
	ctx context.Context, req *protoMessages.UpdateShortURLByOrganizationIDRequest,
) (*protoMessages.UpdateShortURLByOrganizationIDResponse, error) {
	key := schema.ShortURLGlobalKey{Domain: req.Domain, Alias: req.Alias}
	shorturl, err := s.shortURLRepository.GetByOrganizationIDAndGlobalKey(ctx, req.OrganizationId, key)
	if err != nil {
		if errors.Is(err, shortURLRepository.ErrNotFound) {
			return nil, serviceErrors.ErrShortURLNotFound
		}
		return nil, fmt.Errorf("error getting short url: %w", err)
	}

	if err = shortURLUtils.UpdateShortURLFromRequest(shorturl, req); err != nil {
		return nil, fmt.Errorf("error validating update request: %w", err)
	}

	err = s.shortURLRepository.ReplaceByOrganizationIDAndGlobalKey(ctx, req.OrganizationId, key, shorturl)
	if err != nil {
		if errors.Is(err, shortURLRepository.ErrAlreadyExists) {
			return nil, serviceErrors.ErrShortURLAlreadyExists
		}
		return nil, fmt.Errorf("error updating short url: %w", err)
	}

	return &protoMessages.UpdateShortURLByOrganizationIDResponse{
		ShortUrl: shortURLUtils.ShortURLToResponse(shorturl),
	}, nil
}

func (s *ShortURLManagementService) DeleteShortURLByOrganizationID(
	ctx context.Context, req *protoMessages.DeleteShortURLByOrganizationIDRequest,
) (*protoMessages.DeleteShortURLByOrganizationIDResponse, error) {
	key := schema.ShortURLGlobalKey{Domain: req.Domain, Alias: req.Alias}
	_, err := s.shortURLRepository.DeleteByOrganizationIDAndGlobalKey(ctx, req.OrganizationId, key)
	if err != nil {
		if errors.Is(err, shortURLRepository.ErrNotFound) {
			return nil, serviceErrors.ErrShortURLNotFound
		}
		return nil, fmt.Errorf("error deleting short url: %w", err)
	}

	return &protoMessages.DeleteShortURLByOrganizationIDResponse{}, nil
}

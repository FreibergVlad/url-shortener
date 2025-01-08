package generator

import (
	"context"
	"crypto/sha256"
	"io"
	"slices"

	domainServiceMessages "github.com/FreibergVlad/url-shortener/proto/pkg/domain/messages/v1"
	domainService "github.com/FreibergVlad/url-shortener/proto/pkg/domain/service/v1"
	protoMessages "github.com/FreibergVlad/url-shortener/proto/pkg/shorturl/generator/messages/v1"
	protoService "github.com/FreibergVlad/url-shortener/proto/pkg/shorturl/generator/service/v1"
	grpcUtils "github.com/FreibergVlad/url-shortener/shared/go/pkg/api/grpc/utils"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/clock"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/db"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/db/filters"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/errors"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/validators"
	"github.com/FreibergVlad/url-shortener/short-url-management-service/internal/cache"
	"github.com/FreibergVlad/url-shortener/short-url-management-service/internal/config"
	schema "github.com/FreibergVlad/url-shortener/short-url-management-service/internal/db"
	services "github.com/FreibergVlad/url-shortener/short-url-management-service/internal/services/shorturl"
)

type ShortURLAliasEncoder interface {
	Encode(in []byte, len int) (string, error)
}

type shortUrlGeneratorService struct {
	protoService.UnimplementedShortURLGeneratorServiceServer
	shortUrlRepository  db.Repository[schema.ShortURL]
	domainServiceClient domainService.DomainServiceClient
	cache               *cache.ShortURLCache
	encoder             ShortURLAliasEncoder
	saltProvider        io.Reader
	clock               clock.Clock
	config              config.ShortURLGeneratorConfig
}

func New(
	shortUrlRepository db.Repository[schema.ShortURL],
	domainServiceClient domainService.DomainServiceClient,
	cache *cache.ShortURLCache,
	encoder ShortURLAliasEncoder,
	saltProvider io.Reader,
	clock clock.Clock,
	config config.ShortURLGeneratorConfig,
) *shortUrlGeneratorService {
	return &shortUrlGeneratorService{
		shortUrlRepository:  shortUrlRepository,
		domainServiceClient: domainServiceClient,
		cache:               cache,
		encoder:             encoder,
		saltProvider:        saltProvider,
		clock:               clock,
		config:              config,
	}
}

func (s *shortUrlGeneratorService) CreateShortURL(ctx context.Context, req *protoMessages.CreateShortURLRequest) (*protoMessages.CreateShortURLResponse, error) {
	if err := validators.ValidateURL(req.LongUrl); err != nil {
		return nil, errors.NewValidationError("invalid request: %s", err.Error())
	}

	userId := grpcUtils.UserIDFromIncomingContext(ctx)
	draftShortURL := services.ShortURLFromRequest(req, s.clock, s.config.ShortURLScheme, userId)

	if err := s.validateDomain(ctx, draftShortURL); err != nil {
		return nil, err
	}

	shortUrl, err := s.createWithGeneratedAlias(ctx, draftShortURL)
	if err != nil {
		return nil, err
	}

	go s.cache.Set(context.WithoutCancel(ctx), shortUrl)

	return &protoMessages.CreateShortURLResponse{ShortUrl: services.ShortURLToResponse(shortUrl)}, nil
}

func (s *shortUrlGeneratorService) createWithGeneratedAlias(ctx context.Context, shortUrl *schema.ShortURL) (*schema.ShortURL, error) {
	alias, err := s.generateShortURLAlias(shortUrl, nil)
	if err != nil {
		return nil, err
	}

	shortUrl.Alias = alias
	if err = s.shortUrlRepository.Create(ctx, shortUrl); err == nil {
		return shortUrl, nil
	}
	if err != errors.ErrDuplicateResource {
		return nil, err
	}

	filter := filters.Or{Filters: []filters.Filter{
		services.OneShortURLFilter(shortUrl.Domain, shortUrl.Alias),
		filters.And{Filters: []filters.Filter{
			filters.Equals[string]{Field: "organization_id", Value: shortUrl.OrganizationID},
			filters.Equals[string]{Field: "domain", Value: shortUrl.Domain},
			filters.Equals[string]{Field: "long_url.hash", Value: shortUrl.LongURL.Hash},
		}},
	}}
	conflictingUrl, err := s.shortUrlRepository.FindOne(ctx, filter)
	if err != nil {
		return nil, err
	}

	if conflictingUrl.OrganizationID == shortUrl.OrganizationID &&
		shortUrl.Domain == conflictingUrl.Domain &&
		conflictingUrl.LongURL.Hash == shortUrl.LongURL.Hash {
		return s.restoreShortURL(ctx, shortUrl, conflictingUrl)
	}

	return s.handleShortURLCollision(ctx, shortUrl)
}

func (s *shortUrlGeneratorService) restoreShortURL(ctx context.Context, newShortUrl, oldShortUrl *schema.ShortURL) (*schema.ShortURL, error) {
	replaceFilter := services.OneShortURLFilter(oldShortUrl.Domain, oldShortUrl.Alias)
	if err := s.shortUrlRepository.ReplaceOne(ctx, replaceFilter, newShortUrl); err != nil {
		return nil, err
	}
	return newShortUrl, nil
}

func (s *shortUrlGeneratorService) handleShortURLCollision(ctx context.Context, shortUrl *schema.ShortURL) (*schema.ShortURL, error) {
	salt := make([]byte, 4)
	for range s.config.MaxRetriesOnCollision {
		if _, err := s.saltProvider.Read(salt); err != nil {
			return nil, err
		}
		alias, err := s.generateShortURLAlias(shortUrl, salt)
		if err != nil {
			return nil, err
		}

		shortUrl.Alias = alias
		if err = s.shortUrlRepository.Create(ctx, shortUrl); err != nil {
			if err == errors.ErrDuplicateResource {
				continue
			}
			return nil, err
		}
		return shortUrl, nil
	}
	return nil, errors.ErrShortURLCollisionRetriesExceeded
}

func (s *shortUrlGeneratorService) generateShortURLAlias(shortUrl *schema.ShortURL, salt []byte) (string, error) {
	ukey := append([]byte(shortUrl.OrganizationID), []byte(shortUrl.Domain)...)
	ukey = append(ukey, []byte(shortUrl.LongURL.Assembled)...)
	if salt != nil {
		ukey = append(ukey, salt...)
	}
	ukeyHash := sha256.Sum256(ukey)

	alias, err := s.encoder.Encode(ukeyHash[:], s.config.ShortURLAliasLength)
	if err != nil {
		return "", err
	}
	return alias, nil
}

func (s *shortUrlGeneratorService) validateDomain(ctx context.Context, shortUrl *schema.ShortURL) error {
	req := &domainServiceMessages.ListOrganizationDomainRequest{OrganizationId: shortUrl.OrganizationID}
	response, err := s.domainServiceClient.ListOrganizationDomain(grpcUtils.OutgoingContextWithUserID(ctx, shortUrl.CreatedBy), req)
	if err != nil {
		return err
	}
	hasDomain := slices.ContainsFunc(response.Data, func(domain *domainServiceMessages.Domain) bool {
		return domain.Fqdn == shortUrl.Domain
	})
	if hasDomain {
		return nil
	}
	return errors.NewValidationError("invalid request: organization %s is not allowed to use domain %s", shortUrl.OrganizationID, shortUrl.Domain)
}

package generator

import (
	"context"
	"errors"
	"fmt"

	shorturlgeneratorproto "github.com/FreibergVlad/url-shortener/proto/pkg/shorturls/generator/messages/v1"
	shorturlgeneratorserviceproto "github.com/FreibergVlad/url-shortener/proto/pkg/shorturls/generator/service/v1"
	grpcutils "github.com/FreibergVlad/url-shortener/shared/go/pkg/api/grpc/utils"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/clock"
	serviceerrors "github.com/FreibergVlad/url-shortener/shared/go/pkg/errors"
	"github.com/FreibergVlad/url-shortener/short-url-management-service/internal/clients/domains"
	"github.com/FreibergVlad/url-shortener/short-url-management-service/internal/config"
	shorturlrepository "github.com/FreibergVlad/url-shortener/short-url-management-service/internal/db/repositories/shorturls"
	"github.com/FreibergVlad/url-shortener/short-url-management-service/internal/db/schema"
	shorturlutils "github.com/FreibergVlad/url-shortener/short-url-management-service/internal/services/shorturls"
	"github.com/FreibergVlad/url-shortener/short-url-management-service/internal/validators"
)

type ShortURLAliasGenerator interface {
	Generate(shorturl *schema.ShortURL) (string, error)
	GenerateWithSalt(shorturl *schema.ShortURL) (string, error)
}

type DomainChecker interface {
	HasDomain(ctx context.Context, userID, organizationID, fqdn string) error
}

type UserInfoGetter interface {
	GetUserInfo(ctx context.Context, userID string) (schema.User, error)
}

var (
	ErrShortURLCollisionRetriesExceeded = errors.New("maximum retries exceeded while handling collisions of short URL")
	ErrDomainNotAllowed                 = serviceerrors.NewValidationError(
		map[string][]string{"domain": {"domain is not allowed to use"}},
	)
)

type ShortURLGeneratorService struct {
	shorturlgeneratorserviceproto.UnimplementedShortURLGeneratorServiceServer
	shortURLRepository shorturlrepository.Repository
	domainChecker      DomainChecker
	userInfoGetter     UserInfoGetter
	generator          ShortURLAliasGenerator
	clock              clock.Clock
	config             config.ShortURLGeneratorConfig
}

func New(
	shortURLRepository shorturlrepository.Repository,
	domainChecker DomainChecker,
	userInfoGetter UserInfoGetter,
	generator ShortURLAliasGenerator,
	clock clock.Clock,
	config config.ShortURLGeneratorConfig,
) *ShortURLGeneratorService {
	return &ShortURLGeneratorService{
		shortURLRepository: shortURLRepository,
		domainChecker:      domainChecker,
		userInfoGetter:     userInfoGetter,
		generator:          generator,
		clock:              clock,
		config:             config,
	}
}

func (s *ShortURLGeneratorService) CreateShortURL(
	ctx context.Context, req *shorturlgeneratorproto.CreateShortURLRequest,
) (*shorturlgeneratorproto.CreateShortURLResponse, error) {
	if err := validators.ValidateURL(req.LongUrl); err != nil {
		return nil, serviceerrors.NewValidationError(map[string][]string{"long_url": {err.Error()}})
	}

	user, err := s.userInfoGetter.GetUserInfo(ctx, grpcutils.UserIDFromIncomingContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("error getting user info: %w", err)
	}

	draftShortURL := shorturlutils.ShortURLFromRequest(req, s.clock, s.config.ShortURLScheme, user)

	if err = s.domainChecker.HasDomain(ctx, user.ID, draftShortURL.OrganizationID, draftShortURL.Domain); err != nil {
		if errors.Is(err, domains.ErrDomainNotAllowed) {
			return nil, ErrDomainNotAllowed
		}
		return nil, fmt.Errorf("error validating short url domain: %w", err)
	}

	var shorturl *schema.ShortURL
	if customAlias := req.GetAlias().GetValue(); customAlias != "" {
		draftShortURL.Alias = customAlias
		shorturl, err = s.СreateWithCustomAlias(ctx, draftShortURL)
	} else {
		shorturl, err = s.CreateWithGeneratedAlias(ctx, draftShortURL)
	}

	if err != nil {
		return nil, err
	}

	return &shorturlgeneratorproto.CreateShortURLResponse{ShortUrl: shorturlutils.ShortURLToResponse(shorturl)}, nil
}

func (s *ShortURLGeneratorService) СreateWithCustomAlias(
	ctx context.Context, shorturl *schema.ShortURL,
) (*schema.ShortURL, error) {
	if err := s.shortURLRepository.Create(ctx, shorturl); err != nil {
		if errors.Is(err, shorturlrepository.ErrAlreadyExists) {
			return nil, serviceerrors.ErrShortURLAlreadyExists
		}
		return nil, fmt.Errorf("error creating short url with custom alias: %w", err)
	}
	return shorturl, nil
}

func (s *ShortURLGeneratorService) CreateWithGeneratedAlias(
	ctx context.Context, shorturl *schema.ShortURL,
) (*schema.ShortURL, error) {
	alias, err := s.generator.Generate(shorturl)
	if err != nil {
		return nil, fmt.Errorf("error generating short url alias: %w", err)
	}

	shorturl.Alias = alias
	if err = s.shortURLRepository.Create(ctx, shorturl); err == nil {
		return shorturl, nil
	}
	if !errors.Is(err, shorturlrepository.ErrAlreadyExists) {
		return nil, fmt.Errorf("error creating short url with generated alias: %w", err)
	}

	conflictingShortURL, err := s.shortURLRepository.GetByOrganizationKeyOrGlobalKey(
		ctx,
		shorturl.OrganizationKey(),
		shorturl.GlobalKey(),
	)
	if err != nil {
		return nil, fmt.Errorf("error getting short url to detect collision: %w", err)
	}

	if conflictingShortURL.OrganizationKey() == shorturl.OrganizationKey() {
		return nil, serviceerrors.ErrShortURLAlreadyExists
	}

	return s.handleShortURLCollision(ctx, shorturl)
}

func (s *ShortURLGeneratorService) handleShortURLCollision(
	ctx context.Context, shorturl *schema.ShortURL,
) (*schema.ShortURL, error) {
	for range s.config.MaxRetriesOnCollision {
		alias, err := s.generator.GenerateWithSalt(shorturl)
		if err != nil {
			return nil, fmt.Errorf("error generating short url alias salt: %w", err)
		}

		shorturl.Alias = alias
		if err = s.shortURLRepository.Create(ctx, shorturl); err != nil {
			if errors.Is(err, shorturlrepository.ErrAlreadyExists) {
				continue
			}
			return nil, fmt.Errorf("error creating short url during collision handling: %w", err)
		}
		return shorturl, nil
	}
	return nil, ErrShortURLCollisionRetriesExceeded
}

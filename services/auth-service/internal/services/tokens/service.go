package tokens

import (
	"context"
	"errors"
	"fmt"

	"github.com/FreibergVlad/url-shortener/auth-service/internal/config"
	"github.com/FreibergVlad/url-shortener/auth-service/internal/db/repositories"
	"github.com/FreibergVlad/url-shortener/auth-service/internal/db/repositories/users"
	protoMessages "github.com/FreibergVlad/url-shortener/proto/pkg/tokens/messages/v1"
	protoService "github.com/FreibergVlad/url-shortener/proto/pkg/tokens/service/v1"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/clock"
	serviceErrors "github.com/FreibergVlad/url-shortener/shared/go/pkg/errors"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/jwt"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/must"
	"golang.org/x/crypto/bcrypt"
)

type TokenService struct {
	protoService.UnimplementedTokenServiceServer
	config         config.IdentityServiceConfig
	userRepository users.Repository
	clock          clock.Clock
}

func New(
	config config.IdentityServiceConfig,
	userRepository users.Repository,
	clock clock.Clock,
) *TokenService {
	return &TokenService{config: config, userRepository: userRepository, clock: clock}
}

func (s *TokenService) IssueAuthenticationToken(
	ctx context.Context, req *protoMessages.IssueAuthenticationTokenRequest,
) (*protoMessages.IssueAuthenticationTokenResponse, error) {
	user, err := s.userRepository.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			return nil, serviceErrors.ErrInvalidCredentials
		}
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, serviceErrors.ErrInvalidCredentials
	}

	token := must.Return(jwt.IssueForUserID(user.ID, s.config.JWT.Secret, s.clock.Now(), s.config.JWT.LifetimeSeconds))
	refreshToken := must.Return(
		jwt.IssueForUserID(user.ID, s.config.JWT.RefreshSecret, s.clock.Now(), s.config.JWT.RefreshLifetimeSeconds),
	)

	return &protoMessages.IssueAuthenticationTokenResponse{Token: token, RefreshToken: refreshToken}, nil
}

func (s *TokenService) RefreshAuthenticationToken(
	_ context.Context, req *protoMessages.RefreshAuthenticationTokenRequest,
) (*protoMessages.RefreshAuthenticationTokenResponse, error) {
	userID, err := jwt.VerifyAndParseUserID(req.RefreshToken, s.config.JWT.RefreshSecret)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, serviceErrors.ErrTokenExpired
		}
		return nil, serviceErrors.ErrInvalidCredentials
	}

	token := must.Return(jwt.IssueForUserID(userID, s.config.JWT.Secret, s.clock.Now(), s.config.JWT.LifetimeSeconds))

	return &protoMessages.RefreshAuthenticationTokenResponse{Token: token}, nil
}

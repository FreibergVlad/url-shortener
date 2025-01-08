package tokens

import (
	"context"

	"github.com/FreibergVlad/url-shortener/auth-service/internal/config"
	"github.com/FreibergVlad/url-shortener/auth-service/internal/db/repositories/users"
	protoMessages "github.com/FreibergVlad/url-shortener/proto/pkg/tokens/messages/v1"
	protoService "github.com/FreibergVlad/url-shortener/proto/pkg/tokens/service/v1"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/clock"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/errors"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/jwt"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/must"
	"golang.org/x/crypto/bcrypt"
)

type tokenService struct {
	protoService.UnimplementedTokenServiceServer
	config         config.IdentityServiceConfig
	userRepository users.Repository
	clock          clock.Clock
}

func New(
	config config.IdentityServiceConfig,
	userRepository users.Repository,
	clock clock.Clock,
) *tokenService {
	return &tokenService{config: config, userRepository: userRepository, clock: clock}
}

func (s *tokenService) IssueAuthenticationToken(ctx context.Context, req *protoMessages.IssueAuthenticationTokenRequest) (*protoMessages.IssueAuthenticationTokenResponse, error) {
	user, err := s.userRepository.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.NewPermissionDeniedError("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.NewPermissionDeniedError("invalid credentials")
	}

	token := must.Return(jwt.IssueForUserId(user.ID, s.config.JWT.Secret, s.clock.Now(), s.config.JWT.LifetimeSeconds))
	refreshToken := must.Return(jwt.IssueForUserId(user.ID, s.config.JWT.RefreshSecret, s.clock.Now(), s.config.JWT.RefreshLifetimeSeconds))

	return &protoMessages.IssueAuthenticationTokenResponse{Token: token, RefreshToken: refreshToken}, nil
}

func (s *tokenService) RefreshAuthenticationToken(ctx context.Context, req *protoMessages.RefreshAuthenticationTokenRequest) (*protoMessages.RefreshAuthenticationTokenResponse, error) {
	userId, err := jwt.VerifyAndParseUserId(req.RefreshToken, s.config.JWT.RefreshSecret)
	if err != nil {
		return nil, errors.NewPermissionDeniedError("failed to refresh token: %s", err.Error())
	}
	token := must.Return(jwt.IssueForUserId(userId, s.config.JWT.Secret, s.clock.Now(), s.config.JWT.LifetimeSeconds))
	return &protoMessages.RefreshAuthenticationTokenResponse{Token: token}, nil
}

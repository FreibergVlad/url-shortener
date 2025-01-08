package server

import (
	"context"
	"net"
	"time"

	"github.com/FreibergVlad/url-shortener/auth-service/internal/config"
	"github.com/FreibergVlad/url-shortener/auth-service/internal/db/repositories/invitations"
	"github.com/FreibergVlad/url-shortener/auth-service/internal/db/repositories/organizations"
	"github.com/FreibergVlad/url-shortener/auth-service/internal/db/repositories/users"
	"github.com/FreibergVlad/url-shortener/auth-service/internal/db/schema"
	invitationService "github.com/FreibergVlad/url-shortener/auth-service/internal/services/invitations"
	organizationService "github.com/FreibergVlad/url-shortener/auth-service/internal/services/organizations"
	permissionService "github.com/FreibergVlad/url-shortener/auth-service/internal/services/permissions"
	"github.com/FreibergVlad/url-shortener/auth-service/internal/services/permissions/resolver"
	"github.com/FreibergVlad/url-shortener/auth-service/internal/services/roles"
	tokenService "github.com/FreibergVlad/url-shortener/auth-service/internal/services/tokens"
	userService "github.com/FreibergVlad/url-shortener/auth-service/internal/services/users"
	invitationServiceProto "github.com/FreibergVlad/url-shortener/proto/pkg/invitations/service/v1"
	organizationServiceProto "github.com/FreibergVlad/url-shortener/proto/pkg/organizations/service/v1"
	permissionServiceProto "github.com/FreibergVlad/url-shortener/proto/pkg/permissions/service/v1"
	tokenServiceProto "github.com/FreibergVlad/url-shortener/proto/pkg/tokens/service/v1"
	userServiceProto "github.com/FreibergVlad/url-shortener/proto/pkg/users/service/v1"
	grpcAuthorizationMiddleware "github.com/FreibergVlad/url-shortener/shared/go/pkg/api/grpc/middlewares/authorization"
	grpcLoggingMiddleware "github.com/FreibergVlad/url-shortener/shared/go/pkg/api/grpc/middlewares/logging"
	grpcRecoverMiddleware "github.com/FreibergVlad/url-shortener/shared/go/pkg/api/grpc/middlewares/recover"
	grpcValidationMiddleware "github.com/FreibergVlad/url-shortener/shared/go/pkg/api/grpc/middlewares/validation"
	redisCache "github.com/FreibergVlad/url-shortener/shared/go/pkg/cache/redis"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/clock"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/errors"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/must"
	grpcServer "github.com/FreibergVlad/url-shortener/shared/go/pkg/server/grpc"
	"github.com/go-redis/cache/v9"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func ensureAdminUserExists(userRepository users.Repository, config config.IdentityServiceConfig) error {
	user := schema.User{
		ID:           uuid.NewString(),
		Email:        config.Admin.Email,
		PasswordHash: string(must.Return(bcrypt.GenerateFromPassword([]byte(config.Admin.Password), 10))),
		RoleID:       roles.RoleIdSuperAdmin,
		CreatedAt:    time.Now(),
	}
	err := userRepository.Create(context.Background(), &user)
	if err == errors.ErrDuplicateResource {
		return nil
	}
	return err
}

func Bootstrap(config config.IdentityServiceConfig, listener net.Listener) (*grpcServer.GRPCServerWithGracefulShutdown, func()) {
	postgresConnPool := must.Return(pgxpool.New(context.TODO(), config.Postgres.DSN))

	redisOptions := must.Return(redis.ParseURL(config.Redis.DSN))
	redisClient := redis.NewClient(redisOptions)

	redisUserCache := redisCache.New[schema.User](&cache.Options{Redis: redisClient})
	redisOrganizationMembershipsCache := redisCache.New[schema.OrganizationMemberships](&cache.Options{Redis: redisClient})

	userRepository := users.NewUserRepository(postgresConnPool)
	cachedUserRepository := users.NewCachedUserRepository(userRepository, redisUserCache, time.Hour)
	organizationRepository := organizations.NewOrganizationRepository(postgresConnPool)
	cachedOrganizationRepository := organizations.NewCachedUserRepository(
		organizationRepository,
		redisOrganizationMembershipsCache,
		time.Hour,
	)
	invitationRepository := invitations.NewInvitationRepository(postgresConnPool)

	clock := clock.NewSystemClock()

	tokenService := tokenService.New(config, cachedUserRepository, clock)
	permissionResolver := resolver.New(cachedUserRepository, cachedOrganizationRepository)
	permissionService := permissionService.New(permissionResolver)
	userService := userService.New(cachedUserRepository, clock)
	organizationService := organizationService.New(cachedOrganizationRepository, clock)
	invitationService := invitationService.New(
		cachedUserRepository,
		invitationRepository,
		cachedOrganizationRepository,
		clock,
	)

	must.Do(ensureAdminUserExists(cachedUserRepository, config))

	server := grpcServer.NewServer(
		listener,
		grpc.ChainUnaryInterceptor(
			grpcRecoverMiddleware.New(),
			grpcLoggingMiddleware.New(),
			grpcValidationMiddleware.New(),
			grpcAuthorizationMiddleware.New(
				map[string]protoreflect.ServiceDescriptor{
					"TokenService":        tokenServiceProto.File_tokens_service_v1_service_proto.Services().ByName("TokenService"),
					"PermissionService":   permissionServiceProto.File_permissions_service_v1_service_proto.Services().ByName("PermissionService"),
					"UserService":         userServiceProto.File_users_service_v1_service_proto.Services().ByName("UserService"),
					"OrganizationService": organizationServiceProto.File_organizations_service_v1_service_proto.Services().ByName("OrganizationService"),
					"InvitationService":   invitationServiceProto.File_invitations_service_v1_service_proto.Services().ByName("InvitationService"),
				},
				permissionResolver,
			),
		),
	)

	tokenServiceProto.RegisterTokenServiceServer(server, tokenService)
	permissionServiceProto.RegisterPermissionServiceServer(server, permissionService)
	userServiceProto.RegisterUserServiceServer(server, userService)
	organizationServiceProto.RegisterOrganizationServiceServer(server, organizationService)
	invitationServiceProto.RegisterInvitationServiceServer(server, invitationService)

	teardown := func() {
		postgresConnPool.Close()
		must.Do(redisClient.Close())
	}

	return server, teardown
}

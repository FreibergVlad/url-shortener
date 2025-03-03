package server

import (
	"context"
	"crypto/rand"
	"net"
	"time"

	domainServiceProto "github.com/FreibergVlad/url-shortener/proto/pkg/domains/service/v1"
	permissionServiceProto "github.com/FreibergVlad/url-shortener/proto/pkg/permissions/service/v1"
	shortURLGeneratorServiceProto "github.com/FreibergVlad/url-shortener/proto/pkg/shorturls/generator/service/v1"
	shortURLManagementServiceProto "github.com/FreibergVlad/url-shortener/proto/pkg/shorturls/management/service/v1"
	grpcAuthorizationMiddleware "github.com/FreibergVlad/url-shortener/shared/go/pkg/api/grpc/middlewares/authorization"
	grpcLoggingMiddleware "github.com/FreibergVlad/url-shortener/shared/go/pkg/api/grpc/middlewares/logging"
	grpcRecoverMiddleware "github.com/FreibergVlad/url-shortener/shared/go/pkg/api/grpc/middlewares/recoverer"
	grpcValidationMiddleware "github.com/FreibergVlad/url-shortener/shared/go/pkg/api/grpc/middlewares/validation"
	redisCache "github.com/FreibergVlad/url-shortener/shared/go/pkg/cache/redis"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/clock"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/must"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/permissions"
	grpcServer "github.com/FreibergVlad/url-shortener/shared/go/pkg/server/grpc"
	"github.com/FreibergVlad/url-shortener/short-url-management-service/internal/clients/domains"
	"github.com/FreibergVlad/url-shortener/short-url-management-service/internal/clients/users"
	"github.com/FreibergVlad/url-shortener/short-url-management-service/internal/config"
	"github.com/FreibergVlad/url-shortener/short-url-management-service/internal/db/repositories/shorturls"
	"github.com/FreibergVlad/url-shortener/short-url-management-service/internal/db/schema"
	shortURLGeneratorService "github.com/FreibergVlad/url-shortener/short-url-management-service/internal/services/shorturls/generator"
	"github.com/FreibergVlad/url-shortener/short-url-management-service/internal/services/shorturls/generator/encoders/base62"
	sha256AliasGenerator "github.com/FreibergVlad/url-shortener/short-url-management-service/internal/services/shorturls/generator/generators/sha256"
	shortURLManagementService "github.com/FreibergVlad/url-shortener/short-url-management-service/internal/services/shorturls/management"
	"github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	mongoDriver "go.mongodb.org/mongo-driver/v2/mongo"
	mongoDriverOptions "go.mongodb.org/mongo-driver/v2/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func BootstrapShortURLGeneratorService(
	config config.ShortURLGeneratorConfig,
	listener net.Listener,
	grpcPermissionServiceClient permissionServiceProto.PermissionServiceClient,
	grpcUserServiceClient users.UserInfoGetter,
	grpcDomainServiceClient domainServiceProto.DomainServiceClient,
) (*grpcServer.ServerWithGracefulShutdown, func()) {
	logLevel := must.Return(zerolog.ParseLevel(config.LogLevel))
	zerolog.SetGlobalLevel(logLevel)

	mongoClientOptions := mongoDriverOptions.Client().ApplyURI(config.MongoDB.DSN)
	client := must.Return(mongoDriver.Connect(mongoClientOptions))

	redisOptions := must.Return(redis.ParseURL(config.Redis.DSN))
	redisClient := redis.NewClient(redisOptions)
	shortURLCache := redisCache.New[schema.ShortURL](&cache.Options{Redis: redisClient})

	shortURLCollection := client.Database(config.MongoDB.DBName).Collection("short_urls")
	shortURLRepository := shorturls.NewMongoRepository(shortURLCollection)
	cachedShortURLRepository := shorturls.NewCachedRepository(shortURLRepository, shortURLCache, time.Hour)

	aliasGenerator := sha256AliasGenerator.New(base62.New(), rand.Reader, config.ShortURLAliasLength)

	shortURLGeneratorService := shortURLGeneratorService.New(
		cachedShortURLRepository,
		domains.NewServiceClient(grpcDomainServiceClient),
		users.NewServiceClient(grpcUserServiceClient),
		aliasGenerator,
		clock.NewSystemClock(),
		config,
	)
	shortURLGeneratorServiceDesc := shortURLGeneratorServiceProto.
		File_shorturls_generator_service_v1_service_proto.
		Services().
		ByName(protoreflect.Name("ShortURLGeneratorService"))

	teardown := func() {
		must.Do(client.Disconnect(context.TODO()))
		must.Do(redisClient.Close())
	}

	server := grpcServer.NewServer(
		listener,
		grpc.ChainUnaryInterceptor(
			grpcRecoverMiddleware.New(),
			grpcLoggingMiddleware.New(),
			grpcValidationMiddleware.New(),
			grpcAuthorizationMiddleware.New(
				map[string]protoreflect.ServiceDescriptor{"ShortURLGeneratorService": shortURLGeneratorServiceDesc},
				permissions.NewPermissionResolver(grpcPermissionServiceClient),
			),
		),
	)
	shortURLGeneratorServiceProto.RegisterShortURLGeneratorServiceServer(server, shortURLGeneratorService)

	return server, teardown
}

func BootstrapShortURLManagementService(
	config config.ShortURLManagerConfig,
	listener net.Listener,
	grpcPermissionServiceClient permissionServiceProto.PermissionServiceClient,
) (*grpcServer.ServerWithGracefulShutdown, func()) {
	logLevel := must.Return(zerolog.ParseLevel(config.LogLevel))
	zerolog.SetGlobalLevel(logLevel)

	mongoClientOptions := mongoDriverOptions.Client().ApplyURI(config.MongoDB.DSN)
	client := must.Return(mongoDriver.Connect(mongoClientOptions))

	redisOptions := must.Return(redis.ParseURL(config.Redis.DSN))
	redisClient := redis.NewClient(redisOptions)
	shortURLCache := redisCache.New[schema.ShortURL](&cache.Options{Redis: redisClient})

	shortURLCollection := client.Database(config.MongoDB.DBName).Collection("short_urls")
	shortURLRepository := shorturls.NewMongoRepository(shortURLCollection)
	cachedShortURLRepository := shorturls.NewCachedRepository(shortURLRepository, shortURLCache, time.Hour)
	clock := clock.NewSystemClock()
	shortURLManagementService := shortURLManagementService.New(cachedShortURLRepository, clock)
	shortURLManagementServiceDesc := shortURLManagementServiceProto.
		File_shorturls_management_service_v1_service_proto.
		Services().
		ByName(protoreflect.Name("ShortURLManagementService"))

	teardown := func() {
		must.Do(client.Disconnect(context.TODO()))
		must.Do(redisClient.Close())
	}

	server := grpcServer.NewServer(
		listener,
		grpc.ChainUnaryInterceptor(
			grpcRecoverMiddleware.New(),
			grpcLoggingMiddleware.New(),
			grpcValidationMiddleware.New(),
			grpcAuthorizationMiddleware.New(
				map[string]protoreflect.ServiceDescriptor{"ShortURLManagementService": shortURLManagementServiceDesc},
				permissions.NewPermissionResolver(grpcPermissionServiceClient),
			),
		),
	)
	shortURLManagementServiceProto.RegisterShortURLManagementServiceServer(server, shortURLManagementService)

	return server, teardown
}

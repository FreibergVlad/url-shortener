package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	permissionServiceProto "github.com/FreibergVlad/url-shortener/proto/pkg/auth/permissions/service/v1"
	domainServiceProto "github.com/FreibergVlad/url-shortener/proto/pkg/domain/service/v1"
	shortUrlGeneratorServiceProto "github.com/FreibergVlad/url-shortener/proto/pkg/shorturl/generator/service/v1"
	grpcAuthorizationMiddleware "github.com/FreibergVlad/url-shortener/shared/go/pkg/api/grpc/middlewares/authorization"
	grpcLoggingMiddleware "github.com/FreibergVlad/url-shortener/shared/go/pkg/api/grpc/middlewares/logging"
	grpcRecoverMiddleware "github.com/FreibergVlad/url-shortener/shared/go/pkg/api/grpc/middlewares/recover"
	grpcValidationMiddleware "github.com/FreibergVlad/url-shortener/shared/go/pkg/api/grpc/middlewares/validation"
	redisCache "github.com/FreibergVlad/url-shortener/shared/go/pkg/cache/redis"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/clock"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/db/mongo"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/must"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/permissions"
	grpcServer "github.com/FreibergVlad/url-shortener/shared/go/pkg/server/grpc"
	shortURLCache "github.com/FreibergVlad/url-shortener/short-url-management-service/internal/cache"
	"github.com/FreibergVlad/url-shortener/short-url-management-service/internal/config"
	schema "github.com/FreibergVlad/url-shortener/short-url-management-service/internal/db"
	shortUrlGeneratorService "github.com/FreibergVlad/url-shortener/short-url-management-service/internal/services/shorturl/generator"
	"github.com/FreibergVlad/url-shortener/short-url-management-service/internal/services/shorturl/generator/encoders/base62"
	"github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	mongoDriver "go.mongodb.org/mongo-driver/v2/mongo"
	mongoDriverOptions "go.mongodb.org/mongo-driver/v2/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func main() {
	config := config.NewShortURLGeneratorConfig()

	logLevel := must.Return(zerolog.ParseLevel(config.LogLevel))
	zerolog.SetGlobalLevel(logLevel)

	mongoClientOptions := mongoDriverOptions.Client().ApplyURI(config.MongoDB.DSN)
	client := must.Return(mongoDriver.Connect(mongoClientOptions))

	defer func() {
		must.Do(client.Disconnect(context.TODO()))
	}()

	redisOptions := must.Return(redis.ParseURL(config.Redis.DSN))
	redisClient := redis.NewClient(redisOptions)
	redisCache := redisCache.New(&cache.Options{Redis: redisClient})
	shortURLCache := shortURLCache.New(redisCache, time.Hour)

	domainServiceConn := must.Return(grpc.NewClient(config.DomainServiceDSN, grpc.WithTransportCredentials(insecure.NewCredentials())))
	domainServiceClient := domainServiceProto.NewDomainServiceClient(domainServiceConn)

	permissionServiceConn := must.Return(grpc.NewClient(config.AuthServiceDSN, grpc.WithTransportCredentials(insecure.NewCredentials())))
	permissionServiceClient := permissionServiceProto.NewPermissionServiceClient(permissionServiceConn)

	shortUrlCollection := client.Database(config.MongoDB.DBName).Collection("short_urls")
	shortUrlRepository := mongo.NewRepository[schema.ShortURL](shortUrlCollection)
	clock := clock.NewSystemClock()
	shortUrlGeneratorService := shortUrlGeneratorService.New(
		shortUrlRepository,
		domainServiceClient,
		shortURLCache,
		base62.New(),
		rand.Reader,
		clock,
		config,
	)

	server := must.Return(grpcServer.NewServer(
		fmt.Sprintf(":%d", config.Port),
		grpc.ChainUnaryInterceptor(
			grpcRecoverMiddleware.New(),
			grpcLoggingMiddleware.New(),
			grpcValidationMiddleware.New(),
			grpcAuthorizationMiddleware.New(
				map[string]protoreflect.ServiceDescriptor{
					"ShortURLGeneratorService": shortUrlGeneratorServiceProto.File_shorturl_generator_service_v1_service_proto.Services().ByName(protoreflect.Name("ShortURLGeneratorService")),
				},
				permissions.NewPermissionResolver(permissionServiceClient),
			),
		),
	))
	shortUrlGeneratorServiceProto.RegisterShortURLGeneratorServiceServer(server, shortUrlGeneratorService)

	server.Run()
}

package server

import (
	"net"

	"github.com/FreibergVlad/url-shortener/domain-service/internal/config"
	domainService "github.com/FreibergVlad/url-shortener/domain-service/internal/services/domain"
	domainServiceProto "github.com/FreibergVlad/url-shortener/proto/pkg/domains/service/v1"
	permissionServiceProto "github.com/FreibergVlad/url-shortener/proto/pkg/permissions/service/v1"
	grpcAuthorizationMiddleware "github.com/FreibergVlad/url-shortener/shared/go/pkg/api/grpc/middlewares/authorization"
	grpcLoggingMiddleware "github.com/FreibergVlad/url-shortener/shared/go/pkg/api/grpc/middlewares/logging"
	grpcRecoverMiddleware "github.com/FreibergVlad/url-shortener/shared/go/pkg/api/grpc/middlewares/recover"
	grpcValidationMiddleware "github.com/FreibergVlad/url-shortener/shared/go/pkg/api/grpc/middlewares/validation"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/must"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/permissions"
	grpcServer "github.com/FreibergVlad/url-shortener/shared/go/pkg/server/grpc"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func Bootstrap(
	listener net.Listener,
	permissionServiceClient permissionServiceProto.PermissionServiceClient,
	config config.Config,
) *grpcServer.GRPCServerWithGracefulShutdown {
	logLevel := must.Return(zerolog.ParseLevel(config.LogLevel))
	zerolog.SetGlobalLevel(logLevel)

	domainService := domainService.New(config)

	server := grpcServer.NewServer(
		listener,
		grpc.ChainUnaryInterceptor(
			grpcRecoverMiddleware.New(),
			grpcLoggingMiddleware.New(),
			grpcValidationMiddleware.New(),
			grpcAuthorizationMiddleware.New(
				map[string]protoreflect.ServiceDescriptor{
					"DomainService": domainServiceProto.File_domains_service_v1_service_proto.Services().ByName(protoreflect.Name("DomainService")),
				},
				permissions.NewPermissionResolver(permissionServiceClient),
			),
		),
	)
	domainServiceProto.RegisterDomainServiceServer(server, domainService)

	return server
}

package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/FreibergVlad/url-shortener/api-gateway/internal/config"
	"github.com/FreibergVlad/url-shortener/api-gateway/internal/middlewares/authentication"
	"github.com/FreibergVlad/url-shortener/api-gateway/internal/middlewares/cors"
	domainServiceProto "github.com/FreibergVlad/url-shortener/proto/pkg/domains/service/v1"
	invitationServiceProto "github.com/FreibergVlad/url-shortener/proto/pkg/invitations/service/v1"
	organizationServiceProto "github.com/FreibergVlad/url-shortener/proto/pkg/organizations/service/v1"
	shortUrlGeneratorServiceProto "github.com/FreibergVlad/url-shortener/proto/pkg/shorturls/generator/service/v1"
	shortUrlManagementServiceProto "github.com/FreibergVlad/url-shortener/proto/pkg/shorturls/management/service/v1"
	tokenServiceProto "github.com/FreibergVlad/url-shortener/proto/pkg/tokens/service/v1"
	userServiceProto "github.com/FreibergVlad/url-shortener/proto/pkg/users/service/v1"
	grpcUtils "github.com/FreibergVlad/url-shortener/shared/go/pkg/api/grpc/utils"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/must"
	httpWithGracefulShutdown "github.com/FreibergVlad/url-shortener/shared/go/pkg/server/http"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type ctxUserIDKey struct{}

func main() {
	config := config.New()

	logLevel := must.Return(zerolog.ParseLevel(config.LogLevel))
	zerolog.SetGlobalLevel(logLevel)

	// extract user ID from context and pass it to all GRPC requests downstream
	mux := runtime.NewServeMux(runtime.WithMetadata(func(ctx context.Context, _ *http.Request) metadata.MD {
		if userID, ok := ctx.Value(ctxUserIDKey{}).(string); ok {
			return grpcUtils.UserIDMetadata(userID)
		}
		return metadata.MD{}
	}))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	registerGRPCEndpoints(ctx, mux, config)
	httpHandler := registerHTTPMiddlewares(mux, config)

	server := httpWithGracefulShutdown.NewServer(&http.Server{
		Addr:              fmt.Sprintf(":%d", config.Port),
		Handler:           httpHandler,
		ReadHeaderTimeout: httpWithGracefulShutdown.ReadHeaderTimeout,
	})
	server.Run()
}

func registerHTTPMiddlewares(mux *runtime.ServeMux, config config.Config) http.Handler {
	httpHandler := authentication.New(mux, config.JWTSecret, mux, ctxUserIDKey{})
	httpHandler = cors.New(httpHandler)

	return httpHandler
}

func registerGRPCEndpoints(ctx context.Context, mux *runtime.ServeMux, config config.Config) {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	must.Do(userServiceProto.RegisterUserServiceHandlerFromEndpoint(ctx, mux, config.AuthServiceDSN, opts))
	must.Do(tokenServiceProto.RegisterTokenServiceHandlerFromEndpoint(ctx, mux, config.AuthServiceDSN, opts))
	must.Do(
		organizationServiceProto.RegisterOrganizationServiceHandlerFromEndpoint(ctx, mux, config.AuthServiceDSN, opts),
	)
	must.Do(invitationServiceProto.RegisterInvitationServiceHandlerFromEndpoint(ctx, mux, config.AuthServiceDSN, opts))
	must.Do(domainServiceProto.RegisterDomainServiceHandlerFromEndpoint(ctx, mux, config.DomainServiceDSN, opts))
	must.Do(
		shortUrlManagementServiceProto.RegisterShortURLManagementServiceHandlerFromEndpoint(
			ctx, mux, config.ShortURLManagementServiceDSN, opts,
		),
	)
	must.Do(
		shortUrlGeneratorServiceProto.RegisterShortURLGeneratorServiceHandlerFromEndpoint(
			ctx, mux, config.ShortURLGeneratorServiceDSN, opts,
		),
	)
}

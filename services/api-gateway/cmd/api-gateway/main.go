package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/FreibergVlad/url-shortener/api-gateway/internal/config"
	"github.com/FreibergVlad/url-shortener/api-gateway/internal/middlewares/authentication"
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

type ctxUserIdKey struct{}

func main() {
	config := config.New()

	logLevel := must.Return(zerolog.ParseLevel(config.LogLevel))
	zerolog.SetGlobalLevel(logLevel)

	// extract user ID from context and pass it to all GRPC requests downstream
	mux := runtime.NewServeMux(runtime.WithMetadata(func(ctx context.Context, r *http.Request) metadata.MD {
		if userId, ok := ctx.Value(ctxUserIdKey{}).(string); ok {
			return grpcUtils.UserIdMetadata(userId)
		}
		return metadata.MD{}
	}))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	must.Do(userServiceProto.RegisterUserServiceHandlerFromEndpoint(ctx, mux, config.AuthServiceDSN, opts))
	must.Do(tokenServiceProto.RegisterTokenServiceHandlerFromEndpoint(ctx, mux, config.AuthServiceDSN, opts))
	must.Do(organizationServiceProto.RegisterOrganizationServiceHandlerFromEndpoint(ctx, mux, config.AuthServiceDSN, opts))
	must.Do(invitationServiceProto.RegisterInvitationServiceHandlerFromEndpoint(ctx, mux, config.AuthServiceDSN, opts))
	must.Do(domainServiceProto.RegisterDomainServiceHandlerFromEndpoint(ctx, mux, config.DomainServiceDSN, opts))
	must.Do(shortUrlManagementServiceProto.RegisterShortURLManagementServiceHandlerFromEndpoint(ctx, mux, config.ShortUrlManagementServiceDSN, opts))
	must.Do(shortUrlGeneratorServiceProto.RegisterShortURLGeneratorServiceHandlerFromEndpoint(ctx, mux, config.ShortUrlGeneratorServiceDSN, opts))

	handler := authentication.New(mux, config.JWTSecret, mux, ctxUserIdKey{})

	server := httpWithGracefulShutdown.NewServer(&http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: handler,
	})
	server.Run()
}

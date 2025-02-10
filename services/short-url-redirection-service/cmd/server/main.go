package main

import (
	"fmt"
	"net/http"

	protoService "github.com/FreibergVlad/url-shortener/proto/pkg/shorturls/management/service/v1"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/must"
	httpWithGracefulShutdown "github.com/FreibergVlad/url-shortener/shared/go/pkg/server/http"
	"github.com/FreibergVlad/url-shortener/url-redirection-service/internal/config"
	"github.com/FreibergVlad/url-shortener/url-redirection-service/internal/router"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	config := config.New()
	grpcOpts := grpc.WithTransportCredentials(insecure.NewCredentials())

	shortURLManagementServiceConn := must.Return(grpc.NewClient(config.ShortURLManagementServiceDSN, grpcOpts))
	defer shortURLManagementServiceConn.Close()

	shortURLManagementServiceClient := protoService.NewShortURLManagementServiceClient(shortURLManagementServiceConn)
	router := router.New(shortURLManagementServiceClient, config)

	server := httpWithGracefulShutdown.NewServer(&http.Server{
		Addr:              fmt.Sprintf(":%d", config.Port),
		Handler:           router,
		ReadHeaderTimeout: httpWithGracefulShutdown.ReadHeaderTimeout,
	})

	server.Run()
}

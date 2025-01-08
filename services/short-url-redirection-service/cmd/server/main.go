package main

import (
	"fmt"
	"net/http"

	protoService "github.com/FreibergVlad/url-shortener/proto/pkg/shorturl/management/service/v1"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/must"
	httpWithGracefulShutdown "github.com/FreibergVlad/url-shortener/shared/go/pkg/server/http"
	"github.com/FreibergVlad/url-shortener/url-redirection-service/internal/config"
	"github.com/FreibergVlad/url-shortener/url-redirection-service/internal/router"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	config := config.New()

	umsConn := must.Return(grpc.NewClient(config.ShortUrlManagementServiceDSN, grpc.WithTransportCredentials(insecure.NewCredentials())))
	defer umsConn.Close()

	sumsClient := protoService.NewShortURLManagementServiceClient(umsConn)
	router := router.New(sumsClient, config)
	server := httpWithGracefulShutdown.NewServer(&http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: router,
	})

	server.Run()
}

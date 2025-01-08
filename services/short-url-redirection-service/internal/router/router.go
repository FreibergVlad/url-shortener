package router

import (
	"net/http"

	protoService "github.com/FreibergVlad/url-shortener/proto/pkg/shorturl/management/service/v1"
	"github.com/FreibergVlad/url-shortener/url-redirection-service/internal/config"
	"github.com/FreibergVlad/url-shortener/url-redirection-service/internal/handlers"
)

type router struct {
	*http.ServeMux
}

func New(shortUrlManagementServiceClient protoService.ShortURLManagementServiceClient, config config.Config) *router {
	router := &router{ServeMux: http.NewServeMux()}
	redirectionHandler := handlers.NewRedirectionHandler(shortUrlManagementServiceClient, config)

	router.HandleFunc("GET /{alias}", redirectionHandler.HandleRedirect)

	return router
}

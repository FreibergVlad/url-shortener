package router

import (
	"net/http"

	protoService "github.com/FreibergVlad/url-shortener/proto/pkg/shorturls/management/service/v1"
	"github.com/FreibergVlad/url-shortener/url-redirection-service/internal/config"
	"github.com/FreibergVlad/url-shortener/url-redirection-service/internal/handlers"
)

type Router struct {
	*http.ServeMux
}

func New(shortURLManagementServiceClient protoService.ShortURLManagementServiceClient, config config.Config) *Router {
	router := &Router{ServeMux: http.NewServeMux()}
	redirectionHandler := handlers.NewRedirectionHandler(shortURLManagementServiceClient, config)

	router.HandleFunc("GET /{alias}", redirectionHandler.HandleRedirect)

	return router
}

package handlers

import (
	"context"
	"net/http"

	protoMessages "github.com/FreibergVlad/url-shortener/proto/pkg/shorturls/management/messages/v1"
	protoService "github.com/FreibergVlad/url-shortener/proto/pkg/shorturls/management/service/v1"
	"github.com/FreibergVlad/url-shortener/url-redirection-service/internal/config"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RedirectionHandler struct {
	shortURLManagementServiceClient protoService.ShortURLManagementServiceClient
	config                          config.Config
}

func NewRedirectionHandler(
	shortURLServiceManagementClient protoService.ShortURLManagementServiceClient,
	config config.Config,
) *RedirectionHandler {
	return &RedirectionHandler{shortURLManagementServiceClient: shortURLServiceManagementClient, config: config}
}

func (h *RedirectionHandler) HandleRedirect(httpRespWriter http.ResponseWriter, httpReq *http.Request) {
	alias := httpReq.PathValue("alias")
	if alias == "" {
		http.NotFound(httpRespWriter, httpReq)
		return
	}

	redirectURL, err := h.getRedirectURL(httpReq.Context(), alias)
	if err != nil {
		h.handleURLLookupErr(httpRespWriter, httpReq, err)
		return
	}

	http.Redirect(httpRespWriter, httpReq, redirectURL, http.StatusMovedPermanently)
}

func (h *RedirectionHandler) handleURLLookupErr(httpRespWriter http.ResponseWriter, httpReq *http.Request, err error) {
	status, ok := status.FromError(err)
	if !ok || status.Code() != codes.NotFound {
		log.Error().Err(err).Msg("Error while looking for a short URL")
		http.Error(httpRespWriter, "Internal Server Error", http.StatusInternalServerError)
	} else {
		http.NotFound(httpRespWriter, httpReq)
	}
}

func (h *RedirectionHandler) getRedirectURL(ctx context.Context, alias string) (string, error) {
	request := protoMessages.GetShortURLRequest{Alias: alias, Domain: h.config.Domain}

	response, err := h.shortURLManagementServiceClient.GetShortURL(ctx, &request)
	if err != nil {
		return "", err
	}

	return response.ShortUrl.LongUrl.Assembled, nil
}

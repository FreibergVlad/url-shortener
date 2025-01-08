package handlers

import (
	"net/http"

	protoMessages "github.com/FreibergVlad/url-shortener/proto/pkg/shorturl/management/messages/v1"
	protoService "github.com/FreibergVlad/url-shortener/proto/pkg/shorturl/management/service/v1"
	"github.com/FreibergVlad/url-shortener/url-redirection-service/internal/config"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RedirectionHandler struct {
	shortUrlManagementServiceClient protoService.ShortURLManagementServiceClient
	config                          config.Config
}

func NewRedirectionHandler(
	shortUrlServiceManagementClient protoService.ShortURLManagementServiceClient,
	config config.Config,
) *RedirectionHandler {
	return &RedirectionHandler{shortUrlManagementServiceClient: shortUrlServiceManagementClient, config: config}
}

func (h *RedirectionHandler) HandleRedirect(w http.ResponseWriter, r *http.Request) {
	alias := r.PathValue("alias")
	req := protoMessages.GetShortURLRequest{Alias: alias, Domain: h.config.Domain}
	resp, err := h.shortUrlManagementServiceClient.GetShortURL(r.Context(), &req)
	if err != nil {
		h.handleShortURLLookupError(w, r, err)
		return
	}
	if resp.ShortUrl.Status != protoMessages.ShortURLStatus_SHORT_URL_STATUS_ACTIVE {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, resp.ShortUrl.LongUrl.Assembled, http.StatusMovedPermanently)
}

func (h *RedirectionHandler) handleShortURLLookupError(w http.ResponseWriter, r *http.Request, err error) {
	status, ok := status.FromError(err)
	if !ok || status.Code() != codes.NotFound {
		log.Error().Err(err).Msg("Error while looking for a short URL")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	} else {
		http.NotFound(w, r)
	}
}

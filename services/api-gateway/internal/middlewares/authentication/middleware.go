package authentication

import (
	"context"
	"net/http"
	"strings"

	"github.com/FreibergVlad/url-shortener/shared/go/pkg/errors"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/jwt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

func New(
	next http.Handler,
	jwtSecret string,
	mux *runtime.ServeMux,
	ctxUserIDKey any,
) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		authHeader := request.Header.Get("Authorization")
		if authHeader == "" {
			next.ServeHTTP(responseWriter, request)
			return
		}

		userID, err := jwt.VerifyAndParseUserID(strings.TrimPrefix(authHeader, "Bearer "), jwtSecret)
		if err != nil || userID == "" {
			runtime.HTTPError(
				request.Context(),
				mux,
				&runtime.JSONPb{},
				responseWriter,
				request,
				errors.NewUnauthenticatedError("unauthenticated"),
			)
			return
		}

		context := context.WithValue(request.Context(), ctxUserIDKey, userID)

		next.ServeHTTP(responseWriter, request.WithContext(context))
	})
}

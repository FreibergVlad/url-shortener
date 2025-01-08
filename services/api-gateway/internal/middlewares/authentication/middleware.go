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
	ctxUserIdKey any,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			next.ServeHTTP(w, r)
			return
		}

		userId, err := jwt.VerifyAndParseUserId(strings.TrimPrefix(authHeader, "Bearer "), jwtSecret)
		if err != nil || userId == "" {
			runtime.HTTPError(r.Context(), mux, &runtime.JSONPb{}, w, r, errors.NewUnauthenticatedError("unauthenticated"))
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxUserIdKey, userId)))
	})
}

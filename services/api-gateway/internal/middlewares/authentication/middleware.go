package authentication

import (
	"context"
	"errors"
	"net/http"
	"strings"

	serviceErrors "github.com/FreibergVlad/url-shortener/shared/go/pkg/errors"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/jwt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/encoding/protojson"
)

var defaultMarshaler = &runtime.HTTPBodyMarshaler{
	Marshaler: &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			EmitUnpopulated: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	},
}

func New(next http.Handler, jwtSecret string, mux *runtime.ServeMux, ctxUserIDKey any) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		token := strings.TrimPrefix(request.Header.Get("Authorization"), "Bearer ")
		if token == "" {
			// pass request downstream, if authorization is required for this API call,
			// it will be rejected by downstream service
			next.ServeHTTP(responseWriter, request)
			return
		}

		userID, err := verifyTokenAndParseUserID(token, jwtSecret)
		if err != nil {
			runtime.HTTPError(request.Context(), mux, defaultMarshaler, responseWriter, request, err)
			return
		}

		context := context.WithValue(request.Context(), ctxUserIDKey, userID)
		next.ServeHTTP(responseWriter, request.WithContext(context))
	})
}

func verifyTokenAndParseUserID(token string, jwtSecret string) (string, error) {
	userID, err := jwt.VerifyAndParseUserID(token, jwtSecret)
	if err == nil && userID != "" {
		return userID, nil
	}
	if errors.Is(err, jwt.ErrTokenExpired) {
		return "", serviceErrors.ErrTokenExpired
	}
	return "", serviceErrors.ErrUnauthenticated
}

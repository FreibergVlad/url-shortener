package cors

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

const PreflightMaxAge = time.Hour * 24

var (
	AllowedMethods = strings.Join([]string{"GET", "PUT", "POST", "PATCH", "DELETE", "OPTIONS"}, ",")
	AllowedHeaders = strings.Join([]string{"Authorization", "Origin", "Content-Length", "Content-Type"}, ",")
)

func New(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		origin := request.Header.Get("Origin")
		if origin == "" {
			next.ServeHTTP(responseWriter, request)
			return
		}

		headers := responseWriter.Header()

		headers.Set("Access-Control-Allow-Origin", origin)
		headers.Set("Access-Control-Allow-Credentials", "false")
		headers.Add("Vary", "Origin")

		if request.Method == http.MethodOptions {
			headers.Add("Vary", "Access-Control-Request-Method")
			headers.Add("Vary", "Access-Control-Request-Headers")
			headers.Set("Access-Control-Allow-Methods", AllowedMethods)
			headers.Set("Access-Control-Allow-Headers", AllowedHeaders)
			headers.Set("Access-Control-Max-Age", strconv.FormatInt(int64(PreflightMaxAge.Seconds()), 10))

			responseWriter.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(responseWriter, request)
	})
}

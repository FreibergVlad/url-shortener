package recoverer

import (
	"context"
	"runtime/debug"

	"github.com/FreibergVlad/url-shortener/shared/go/pkg/errors"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

func New() grpc.UnaryServerInterceptor {
	return middleware
}

//nolint:nonamedreturns
func middleware(
	ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (resp any, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Error().Msgf("panic occurred: %s", string(debug.Stack()))
			err = errors.ErrInternalError
		}
	}()
	return handler(ctx, req)
}

package logging

import (
	"context"
	"time"

	"github.com/FreibergVlad/url-shortener/shared/go/pkg/errors"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func New() grpc.UnaryServerInterceptor {
	return middleware
}

func middleware(
	ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (any, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	latency := time.Since(start)

	event, code, errMsg := log.Info(), codes.OK, "" //nolint:zerologlint
	if err != nil {
		grpcStatus := ensureGRPCStatus(err)

		code, errMsg = grpcStatus.Code(), grpcStatus.Message()
		event = log.Error() //nolint:zerologlint

		if code == codes.Internal {
			err = errors.ErrInternalError
		}
	}

	event.
		Str("code", code.String()).
		Str("method", info.FullMethod).
		Dur("latency", latency).
		Msg(errMsg)

	return resp, err
}

func ensureGRPCStatus(err error) *status.Status {
	errStatus, ok := status.FromError(err)
	if !ok {
		errStatus = status.New(codes.Internal, err.Error())
	}
	return errStatus
}

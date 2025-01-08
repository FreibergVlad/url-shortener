package logging

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func New() grpc.UnaryServerInterceptor {
	return middleware
}

func middleware(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	start := time.Now()
	resp, err = handler(ctx, req)
	latency := time.Since(start)

	event, code, errMsg := log.Info(), codes.OK, ""
	if err != nil {
		errStatus, ok := status.FromError(err)
		if !ok {
			errStatus = status.New(codes.Internal, err.Error())
		}
		err, code, errMsg = errStatus.Err(), errStatus.Code(), errStatus.Message()
		event = log.Error()
	}

	event.
		Str("code", code.String()).
		Str("method", info.FullMethod).
		Dur("latency", latency).
		Msg(errMsg)

	return resp, err
}

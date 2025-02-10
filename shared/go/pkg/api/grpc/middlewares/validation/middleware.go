package validation

import (
	"context"

	"github.com/FreibergVlad/url-shortener/shared/go/pkg/errors"
	protoValidator "github.com/bufbuild/protovalidate-go"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

func New() grpc.UnaryServerInterceptor {
	return middleware
}

func middleware(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	protoReq, ok := req.(proto.Message)
	if !ok {
		return nil, errors.NewValidationError("invalid request: not a google.golang.org/protobuf/proto.Message")
	}
	if err := protoValidator.Validate(protoReq); err != nil {
		return nil, errors.NewValidationError("invalid request: %s", err.Error())
	}
	return handler(ctx, req)
}

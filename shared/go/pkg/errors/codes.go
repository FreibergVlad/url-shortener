package errors

import grpcCodes "google.golang.org/grpc/codes"

type code int

const (
	CodeInternal code = iota
	CodeUnauthenticated
	CodePermissionDenied
	CodeInvalidArgument
	CodeAlreadyExists
	CodeNotFound
)

var grpcCodesMapping = map[code]grpcCodes.Code{
	CodeInternal:         grpcCodes.Internal,
	CodeUnauthenticated:  grpcCodes.Unauthenticated,
	CodePermissionDenied: grpcCodes.PermissionDenied,
	CodeInvalidArgument:  grpcCodes.InvalidArgument,
	CodeAlreadyExists:    grpcCodes.AlreadyExists,
	CodeNotFound:         grpcCodes.NotFound,
}

func (c code) grpsCode() grpcCodes.Code {
	code, ok := grpcCodesMapping[c]
	if !ok {
		return grpcCodes.Unknown
	}
	return code
}

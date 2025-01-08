package errors

import (
	"fmt"

	"google.golang.org/grpc/status"
)

var (
	ErrDuplicateResource                = NewValidationError("resource already exists")
	ErrResourceNotFound                 = NewNotFoundError("resource doesn't exist")
	ErrShortURLCollisionRetriesExceeded = NewInternalError("maximum retries exceeded while handling collisions of short URL")
)

type ServiceError interface {
	GRPCError() error
	Error() string
}

type serviceError struct {
	code code
	msg  string
}

func NewError(code code, msg string, a ...any) serviceError {
	return serviceError{code: code, msg: fmt.Sprintf(msg, a...)}
}

func NewPermissionDeniedError(msg string, a ...any) serviceError {
	return NewError(CodePermissionDenied, msg, a...)
}

func NewUnauthenticatedError(msg string, a ...any) serviceError {
	return NewError(CodeUnauthenticated, msg, a...)
}

func NewValidationError(msg string, a ...any) serviceError {
	return NewError(CodeInvalidArgument, msg, a...)
}

func NewInternalError(msg string, a ...any) serviceError {
	return NewError(CodeInternal, msg, a...)
}

func NewNotFoundError(msg string, a ...any) serviceError {
	return NewError(CodeNotFound, msg, a...)
}

func (err serviceError) GRPCStatus() *status.Status {
	return status.New(err.code.grpsCode(), err.msg)
}

func (err serviceError) Error() string {
	return err.msg
}

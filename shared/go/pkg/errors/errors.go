package errors

import (
	"fmt"

	"google.golang.org/grpc/status"
)

var (
	ErrDuplicateResource = NewValidationError("resource already exists")
	ErrResourceNotFound  = NewNotFoundError("resource doesn't exist")
)

type ServiceError struct {
	code code
	msg  string
}

func NewError(code code, msg string, a ...any) ServiceError {
	return ServiceError{code: code, msg: fmt.Sprintf(msg, a...)}
}

func NewPermissionDeniedError(msg string, a ...any) ServiceError {
	return NewError(CodePermissionDenied, msg, a...)
}

func NewUnauthenticatedError(msg string, a ...any) ServiceError {
	return NewError(CodeUnauthenticated, msg, a...)
}

func NewValidationError(msg string, a ...any) ServiceError {
	return NewError(CodeInvalidArgument, msg, a...)
}

func NewInternalError(msg string, a ...any) ServiceError {
	return NewError(CodeInternal, msg, a...)
}

func NewNotFoundError(msg string, a ...any) ServiceError {
	return NewError(CodeNotFound, msg, a...)
}

func (err ServiceError) GRPCStatus() *status.Status {
	return status.New(err.code.grpsCode(), err.msg)
}

func (err ServiceError) Error() string {
	return err.msg
}

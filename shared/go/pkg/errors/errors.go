package errors

import (
	"fmt"

	"github.com/FreibergVlad/url-shortener/shared/go/pkg/must"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/protoadapt"
)

var (
	ErrUnauthenticated         = newUnauthenticatedError()
	ErrInvalidCredentials      = newInvalidCredentialsError()
	ErrTokenExpired            = newTokenExpiredError()
	ErrInsufficientPermissions = newInsufficientPermissionsError()

	ErrUserNotFound                   = newNotFoundError(ResourceTypeUser)
	ErrOrganizationInvitationNotFound = newNotFoundError(ResourceTypeOrganizationInvitation)
	ErrShortURLNotFound               = newNotFoundError(ResourceTypeShortURL)

	ErrUserAlreadyExists                   = newDuplicatedResourceError(ResourceTypeUser)
	ErrOrganizationAlreadyExists           = newDuplicatedResourceError(ResourceTypeOrganization)
	ErrOrganizationMembershipAlreadyExists = newDuplicatedResourceError(ResourceTypeOrganizationMembership)
	ErrOrganizationInvitationAlreadyExists = newDuplicatedResourceError(ResourceTypeOrganizationInvitation)
	ErrShortURLAlreadyExists               = newDuplicatedResourceError(ResourceTypeShortURL)

	ErrInternalError = newInternalError()
)

type ErrReason string

const (
	ErrReasonInternal                ErrReason = "INTERNAL_SERVER_ERROR"
	ErrReasonBadRequest              ErrReason = "BAD_REQUEST"
	ErrReasonUnauthenticated         ErrReason = "UNAUTHENTICATED"
	ErrReasonInvalidCredentials      ErrReason = "INVALID_CREDENTIALS" //nolint:gosec
	ErrReasonInsufficientPermissions ErrReason = "INSUFFICIENT_PERMISSIONS"
	ErrReasonTokenExpired            ErrReason = "TOKEN_EXPIRED"
	ErrReasonNotFound                ErrReason = "NOT_FOUND"
	ErrReasonAlreadyExists           ErrReason = "ALREADY_EXISTS"
)

type ResourceType string

const (
	ResourceTypeUser                   ResourceType = "user"
	ResourceTypeOrganization           ResourceType = "organization"
	ResourceTypeOrganizationMembership ResourceType = "organization-membership"
	ResourceTypeOrganizationInvitation ResourceType = "organization-invitaton"
	ResourceTypeShortURL               ResourceType = "short-url"
)

func NewValidationError(fieldViolations map[string][]string) error {
	errorInfo := &errdetails.BadRequest{}
	for field, violations := range fieldViolations {
		for _, violation := range violations {
			fieldViolation := &errdetails.BadRequest_FieldViolation{Field: field, Description: violation}
			errorInfo.FieldViolations = append(errorInfo.FieldViolations, fieldViolation)
		}
	}

	return newError(codes.InvalidArgument, "Invalid request parameters", ErrReasonBadRequest, errorInfo)
}

func newDuplicatedResourceError(resourceType ResourceType) error {
	msg := fmt.Sprintf("Resource '%s' already exists", resourceType)
	return newError(codes.AlreadyExists, msg, ErrReasonAlreadyExists)
}

func newNotFoundError(resourceType ResourceType) error {
	msg := fmt.Sprintf("Resource '%s' not found", resourceType)
	return newError(codes.NotFound, msg, ErrReasonNotFound)
}

func newUnauthenticatedError() error {
	return newError(codes.Unauthenticated, "Unauthenticated", ErrReasonUnauthenticated)
}

func newInvalidCredentialsError() error {
	return newError(codes.Unauthenticated, "Invalid credentials", ErrReasonInvalidCredentials)
}

func newTokenExpiredError() error {
	return newError(codes.Unauthenticated, "Token expired", ErrReasonTokenExpired)
}

func newInsufficientPermissionsError() error {
	msg := "Not enough permissions to access resource"
	return newError(codes.PermissionDenied, msg, ErrReasonInsufficientPermissions)
}

func newInternalError() error {
	return newError(codes.Internal, "Internal server error", ErrReasonInternal)
}

func newError(code codes.Code, msg string, reason ErrReason, extraDetails ...protoadapt.MessageV1) error {
	errorInfo := &errdetails.ErrorInfo{Reason: string(reason)}
	details := append([]protoadapt.MessageV1{errorInfo}, extraDetails...)

	status := must.Return(status.New(code, msg).WithDetails(details...))

	return status.Err()
}

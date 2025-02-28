package validation

import (
	"context"
	"errors"
	"strings"

	"buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	serviceErrors "github.com/FreibergVlad/url-shortener/shared/go/pkg/errors"
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
		fieldValidations := map[string][]string{"request": {"not a google.golang.org/protobuf/proto.Message"}}
		return nil, serviceErrors.NewValidationError(fieldValidations)
	}
	if err := protoValidator.Validate(protoReq); err != nil {
		return nil, unwrapValidationErr(err)
	}
	return handler(ctx, req)
}

func unwrapValidationErr(err error) error {
	fieldViolations := map[string][]string{}
	var validationErr *protoValidator.ValidationError
	if ok := errors.As(err, &validationErr); ok {
		for _, violation := range validationErr.ToProto().Violations {
			field := fieldNameFromViolation(violation)
			violations, ok := fieldViolations[field]
			if ok {
				fieldViolations[field] = append(violations, *violation.Message)
			} else {
				fieldViolations[field] = []string{*violation.Message}
			}
		}
	}
	return serviceErrors.NewValidationError(fieldViolations)
}

func fieldNameFromViolation(violation *validate.Violation) string {
	fieldName := strings.Builder{}
	for i, e := range violation.Field.Elements {
		_, _ = fieldName.WriteString(*e.FieldName)
		if i != len(violation.Field.Elements)-1 {
			fieldName.WriteString(".")
		}
	}
	return fieldName.String()
}

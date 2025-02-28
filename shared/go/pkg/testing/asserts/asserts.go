package asserts

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func AssertValidationErrorContainsFieldViolation(t *testing.T, err error, field, msg string) {
	t.Helper()

	require.Error(t, err, "err is nil")

	status, ok := status.FromError(err) //nolint:varnamelen

	require.True(t, ok, "error is not compatible with google.golang.org/grpc/status")
	require.Equal(t, codes.InvalidArgument, status.Code(), "status code is not 'codes.InvalidArgument'")

	var errInfo *errdetails.BadRequest
	for _, detail := range status.Details() {
		if errInfo, ok = detail.(*errdetails.BadRequest); ok {
			break
		}
	}

	require.NotNil(t, errInfo, "status doesn't contain 'errdetails.BadRequest' detail")

	for _, violation := range errInfo.FieldViolations {
		if violation.Field == field && strings.Contains(violation.Description, msg) {
			return
		}
	}

	failMsg := fmt.Sprintf("violation of field [%s] want: [%s], got: [%s]", field, msg, errInfo.FieldViolations)
	require.Fail(t, failMsg)
}

func AssertValidationErrorContainsFieldViolations(t *testing.T, err error, fieldViolations map[string][]string) {
	t.Helper()

	for field, messages := range fieldViolations {
		for _, msg := range messages {
			AssertValidationErrorContainsFieldViolation(t, err, field, msg)
		}
	}
}

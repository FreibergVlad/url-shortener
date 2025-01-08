package jwt_test

import (
	"testing"
	"time"

	"github.com/FreibergVlad/url-shortener/shared/go/pkg/jwt"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/must"
	jwtLib "github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestIssueForUserId(t *testing.T) {
	t.Parallel()

	issuedAt := must.Return(time.Parse(time.RFC3339Nano, "2006-01-02T15:04:05.999999999Z"))
	expectedToken := "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJmYWtlLXVzZXItaWQiLCJleHAiOjExMzYyMTQyNzUsImlhdCI6MTEzNjIxNDI0NX0.I93GY0_CDHZdqpQt8kv9k-QNqJKklDDMC9PhBA9E-_Pv42pBaIf6RJpqlK3aRwnoEMTsWaKlpoji0dvVwS5hEQ"

	actualToken, err := jwt.IssueForUserId("fake-user-id", "fake-secret", issuedAt, 30)

	assert.NoError(t, err)
	assert.Equal(t, expectedToken, actualToken)
}

func TestVerifyAndParseUserId(t *testing.T) {
	t.Parallel()

	tests := []struct {
		issueSecret    string
		verifySecret   string
		issuedAt       time.Time
		expectedUserId string
		errorContains  string
	}{
		{issueSecret: "fake-secret", verifySecret: "fake-secret", issuedAt: time.Now(), expectedUserId: "fake-user-id"},
		{issueSecret: "fake-secret", verifySecret: "incorrect-secret", issuedAt: time.Now(), errorContains: "signature is invalid"},
		{issueSecret: "fake-secret", verifySecret: "fake-secret", issuedAt: time.Now().Add(-time.Hour), errorContains: "token is expired"},
	}

	for _, test := range tests {
		token := must.Return(jwt.IssueForUserId(test.expectedUserId, test.issueSecret, test.issuedAt, 60))
		actualUserId, err := jwt.VerifyAndParseUserId(token, test.verifySecret)

		if test.errorContains != "" {
			assert.ErrorContains(t, err, test.errorContains)
		} else {
			assert.NoError(t, err)
		}
		assert.Equal(t, test.expectedUserId, actualUserId)
	}
}

func TestVerifyAndParseUserIdWithInvalidAlgorithm(t *testing.T) {
	t.Parallel()

	token := must.Return(jwtLib.New(jwtLib.SigningMethodNone).SignedString(jwtLib.UnsafeAllowNoneSignatureType))

	actualUserId, err := jwt.VerifyAndParseUserId(token, string(jwtLib.UnsafeAllowNoneSignatureType))

	assert.ErrorContains(t, err, "signing method none is invalid")
	assert.Equal(t, "", actualUserId)
}

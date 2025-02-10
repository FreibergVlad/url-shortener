package jwt_test

import (
	"testing"
	"time"

	"github.com/FreibergVlad/url-shortener/shared/go/pkg/jwt"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/must"
	jwtLib "github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIssueForUserId(t *testing.T) {
	t.Parallel()

	issuedAt := must.Return(time.Parse(time.RFC3339Nano, "2006-01-02T15:04:05.999999999Z"))
	expectedToken := "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJmYWtlLXVzZXItaWQiLCJleHAiOjExMzYyMTQyNzUsImlhdCI6MTEzNjIxNDI0NX0.I93GY0_CDHZdqpQt8kv9k-QNqJKklDDMC9PhBA9E-_Pv42pBaIf6RJpqlK3aRwnoEMTsWaKlpoji0dvVwS5hEQ" //nolint: gosec,lll

	actualToken, err := jwt.IssueForUserID("fake-user-id", "fake-secret", issuedAt, 30)

	require.NoError(t, err)

	assert.Equal(t, expectedToken, actualToken)
}

func TestVerifyAndParseUserId(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		issueSecret    string
		verifySecret   string
		issuedAt       time.Time
		expectedUserID string
		errorContains  string
	}{
		{
			name:           "success",
			issueSecret:    "fake-secret",
			verifySecret:   "fake-secret",
			issuedAt:       time.Now(),
			expectedUserID: "fake-user-id",
		},
		{
			name:          "incorrect secret used to verify token",
			issueSecret:   "fake-secret",
			verifySecret:  "incorrect-secret",
			issuedAt:      time.Now(),
			errorContains: "signature is invalid",
		},
		{
			name:          "token is expired",
			issueSecret:   "fake-secret",
			verifySecret:  "fake-secret",
			issuedAt:      time.Now().Add(-time.Hour),
			errorContains: "token is expired",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			token := must.Return(jwt.IssueForUserID(test.expectedUserID, test.issueSecret, test.issuedAt, 60))
			actualUserID, err := jwt.VerifyAndParseUserID(token, test.verifySecret)

			if test.errorContains != "" {
				require.ErrorContains(t, err, test.errorContains)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, test.expectedUserID, actualUserID)
		})
	}
}

func TestVerifyAndParseUserIdWithInvalidAlgorithm(t *testing.T) {
	t.Parallel()

	token := must.Return(jwtLib.New(jwtLib.SigningMethodNone).SignedString(jwtLib.UnsafeAllowNoneSignatureType))

	actualUserID, err := jwt.VerifyAndParseUserID(token, string(jwtLib.UnsafeAllowNoneSignatureType))

	require.ErrorContains(t, err, "signing method none is invalid")

	assert.Equal(t, "", actualUserID)
}

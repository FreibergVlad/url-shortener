package validators_test

import (
	"strings"
	"testing"

	"github.com/FreibergVlad/url-shortener/short-url-management-service/internal/validators"
	"github.com/stretchr/testify/assert"
)

func TestValidateURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		wantErr  bool
		errorMsg string
	}{
		// valid cases
		{"Valid HTTP URL", "http://example.com/", false, ""},
		{"Valid HTTPS URL", "https://example.com/", false, ""},
		{"Valid Subdomain URL", "https://sub.example.co.uk", false, ""},
		{"Valid URL with Port", "https://example.com:8080", false, ""},
		{"Valid URL with Query Params", "https://example.com?param1=val1&param2=val2", false, ""},

		// invalid cases
		{"Too Long Input", strings.Repeat("a", validators.MaxURLLen+1), true, "URL length should be"},
		{"Empty URL", "", true, "empty url"},
		{"Malformed URL", "invalid-url", true, "invalid URI"},
		{"Unsupported Scheme", "ftp://example.com", true, "only http and https"},
		{
			"Too Long Hostname",
			"https://" + strings.Repeat("a", validators.MaxHostnameLen+1) + ".com",
			true,
			"hostname length should be",
		},
		{"No TLD", "http://example", true, "is not valid hostname"},
		{"Leading Hyphen", "http://-example.com", true, "is not valid hostname"},
		{"Double Dot in Host", "http://example..com", true, "is not valid hostname"},
		{"Too Short TLD", "http://example.c", true, "is not valid hostname"},
		{"Special Characters", "https://ex$amp&le!.com", true, "is not valid hostname"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := validators.ValidateURL(test.input)

			if test.wantErr {
				assert.ErrorContains(t, err, test.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

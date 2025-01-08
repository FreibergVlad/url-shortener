package validators_test

import (
	"strings"
	"testing"

	"github.com/FreibergVlad/url-shortener/shared/go/pkg/validators"
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
		{"Too Long Input", strings.Repeat("a", validators.MaxURLLen+1), true, "invalid URL length"},
		{"Empty URL", "", true, "invalid URL format"},
		{"Malformed URL", "invalid-url", true, "invalid URL format"},
		{"Unsupported Scheme", "ftp://example.com", true, "invalid URL scheme"},
		{"Too Long Hostname", "https://" + strings.Repeat("a", validators.MaxHostnameLen+1) + ".com", true, "invalid hostname length"},
		{"No TLD", "http://example", true, "invalid host format"},
		{"Leading Hyphen", "http://-example.com", true, "invalid host format"},
		{"Double Dot in Host", "http://example..com", true, "invalid host format"},
		{"Too Short TLD", "http://example.c", true, "invalid host format"},
		{"Special Characters", "https://ex$amp&le!.com", true, "invalid host format"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validators.ValidateURL(test.input)

			if test.wantErr {
				assert.ErrorContains(t, err, test.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

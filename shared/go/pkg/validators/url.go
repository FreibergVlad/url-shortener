package validators

import (
	"fmt"
	"net/url"
	"regexp"
)

var validDomainRegex = regexp.MustCompile(`^([a-zA-Z0-9][a-zA-Z0-9-]{0,61}[a-zA-Z0-9]\.)+[a-zA-Z]{2,}$`)

const (
	MaxURLLen      = 2000
	MaxHostnameLen = 253
)

func ValidateURL(input string) error {
	if len(input) > MaxURLLen {
		return fmt.Errorf("invalid URL length: length should be <= %d", MaxURLLen)
	}

	parsedURL, err := url.ParseRequestURI(input)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("invalid URL scheme: '%s' (only http and https are allowed)", parsedURL.Scheme)
	}

	hostname := parsedURL.Hostname()

	if len(hostname) > MaxHostnameLen {
		return fmt.Errorf("invalid hostname length: length should be <= %d", MaxHostnameLen)
	}

	if !validDomainRegex.MatchString(hostname) {
		return fmt.Errorf("invalid host format")
	}

	return nil
}

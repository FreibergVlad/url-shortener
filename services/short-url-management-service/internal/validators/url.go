package validators

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
)

var (
	validDomainRegex = regexp.MustCompile(`^([a-zA-Z0-9][a-zA-Z0-9-]{0,61}[a-zA-Z0-9]\.)+[a-zA-Z]{2,}$`)
	ErrInvalidURL    = errors.New("invalid URL")
)

const (
	MaxURLLen      = 2000
	MaxHostnameLen = 253
)

func ValidateURL(input string) error {
	if len(input) > MaxURLLen {
		return fmt.Errorf("%w: URL length should be <= %d", ErrInvalidURL, MaxURLLen)
	}

	parsedURL, err := url.ParseRequestURI(input)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrInvalidURL, err.Error())
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("%w: only http and https are allowed, '%s' given", ErrInvalidURL, parsedURL.Scheme)
	}

	hostname := parsedURL.Hostname()

	if len(hostname) > MaxHostnameLen {
		return fmt.Errorf("%w: hostname length should be <= %d", ErrInvalidURL, MaxHostnameLen)
	}

	if !validDomainRegex.MatchString(hostname) {
		return fmt.Errorf("%w: %s is not valid hostname", ErrInvalidURL, hostname)
	}

	return nil
}

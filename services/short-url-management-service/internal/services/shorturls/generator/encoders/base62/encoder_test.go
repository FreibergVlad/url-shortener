package base62_test

import (
	"testing"

	"github.com/FreibergVlad/url-shortener/short-url-management-service/internal/services/shorturls/generator/encoders/base62"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncode_InvalidLength(t *testing.T) {
	t.Parallel()

	encoder, length := base62.New(), 6

	_, err := encoder.Encode([]byte{}, length)
	require.ErrorContains(t, err, "encoded string <")

	_, err = encoder.Encode([]byte("td"), length)
	require.ErrorContains(t, err, "encoded string <")
}

func TestEncode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in       []byte
		expected string
	}{
		{in: []byte("test"), expected: "289lyu"},
		{in: []byte("random-string"), expected: "345jPN"},
		{in: []byte("?!?*_*"), expected: "Ji2KiC"},
	}

	for _, test := range tests {
		encoder, length := base62.New(), 6

		actual, err := encoder.Encode(test.in, length)

		require.NoError(t, err)
		assert.Equal(t, test.expected, actual)
		assert.Len(t, actual, length)
	}
}

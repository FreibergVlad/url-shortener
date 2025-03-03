package testing

import (
	"context"
	"crypto/rand"
	"math/big"
	"testing"
	"time"

	shortURLGeneratorMessages "github.com/FreibergVlad/url-shortener/proto/pkg/shorturls/generator/messages/v1"
	shortURLManagementMessages "github.com/FreibergVlad/url-shortener/proto/pkg/shorturls/management/messages/v1"
	grpcUtils "github.com/FreibergVlad/url-shortener/shared/go/pkg/api/grpc/utils"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/must"
	"github.com/FreibergVlad/url-shortener/short-url-management-service/internal/db/schema"
	"github.com/FreibergVlad/url-shortener/short-url-management-service/internal/services/shorturls/generator/encoders/base62"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var FakeUser = schema.User{
	ID:    gofakeit.UUID(),
	Email: gofakeit.Email(),
}

const shortURLExpirationInterval = time.Hour * 24

func CreateShortURLRequest() *shortURLGeneratorMessages.CreateShortURLRequest {
	return &shortURLGeneratorMessages.CreateShortURLRequest{
		OrganizationId: gofakeit.UUID(),
		Domain:         ShortURLDomain,
		LongUrl:        gofakeit.URL(),
		ExpiresAt:      timestamppb.New(time.Now().Add(shortURLExpirationInterval)),
		Tags:           []string{"tag1", "tag2"},
		Description:    gofakeit.ProductDescription(),
	}
}

func AssertCreateShortURL(
	t *testing.T, server *Server, request *shortURLGeneratorMessages.CreateShortURLRequest,
) *shortURLManagementMessages.ShortURL {
	t.Helper()

	response, err := server.ShortURLGeneratorServiceClient.CreateShortURL(
		grpcUtils.OutgoingContextWithUserID(context.Background(), FakeUser.ID),
		request,
	)

	require.NoError(t, err)

	assert.Equal(t, request.OrganizationId, response.ShortUrl.OrganizationId)
	assert.Equal(t, request.LongUrl, response.ShortUrl.LongUrl.Assembled)
	assert.NotEmpty(t, response.ShortUrl.ShortUrl)
	assert.Equal(t, request.Domain, response.ShortUrl.Domain)
	assert.NotEmpty(t, response.ShortUrl.Alias)
	AssertTimestampNearlyEqual(t, request.ExpiresAt, response.ShortUrl.ExpiresAt)
	assert.NotEmpty(t, response.ShortUrl.CreatedAt)
	assert.Equal(t, request.Description, response.ShortUrl.Description)
	assert.Equal(t, request.Tags, response.ShortUrl.Tags)

	user := &shortURLManagementMessages.User{Id: FakeUser.ID, Email: FakeUser.Email}
	assert.Equal(t, user, response.ShortUrl.CreatedBy)

	return response.ShortUrl
}

func AssertShortURLEqual(
	t *testing.T, want *shortURLManagementMessages.ShortURL, got *shortURLManagementMessages.ShortURL,
) {
	t.Helper()

	assert.Equal(t, want.OrganizationId, got.OrganizationId)
	assert.Equal(t, want.LongUrl, got.LongUrl)
	assert.Equal(t, want.ShortUrl, got.ShortUrl)
	assert.Equal(t, want.Domain, got.Domain)
	assert.Equal(t, want.Alias, got.Alias)
	AssertTimestampNearlyEqual(t, want.ExpiresAt, got.ExpiresAt)
	AssertTimestampNearlyEqual(t, want.CreatedAt, got.CreatedAt)
	assert.Equal(t, want.CreatedBy, got.CreatedBy)
	assert.Equal(t, want.Description, got.Description)
	assert.Equal(t, want.Tags, got.Tags)
}

func AssertTimestampNearlyEqual(t *testing.T, want *timestamppb.Timestamp, got *timestamppb.Timestamp) {
	t.Helper()

	assert.WithinDuration(t, want.AsTime(), got.AsTime(), time.Second)
}

func RandomBase62String(length int) string {
	result := make([]byte, length)
	for i := range length {
		index := must.Return(rand.Int(rand.Reader, big.NewInt(int64(base62.Base))))
		result[i] = base62.Alphabet[index.Int64()]
	}

	return string(result)
}

func RandomBase62StringInRange(minLen, maxLen int) string {
	length := must.Return(rand.Int(rand.Reader, big.NewInt(int64(maxLen-minLen+1))))
	finalLength := int(length.Int64()) + minLen

	return RandomBase62String(finalLength)
}

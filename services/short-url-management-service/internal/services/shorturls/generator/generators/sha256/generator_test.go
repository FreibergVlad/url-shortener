package sha256_test

import (
	"testing"

	"github.com/FreibergVlad/url-shortener/short-url-management-service/internal/db/schema"
	sha256AliasGenerator "github.com/FreibergVlad/url-shortener/short-url-management-service/internal/services/shorturls/generator/generators/sha256"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	fakeOrganizationID = "fake-org-id"
	fakeDomain         = "fake-domain.com"
	fakeLongURL        = "https://fake-url.com"
	fakeAliasLength    = 5
)

var (
	fakeSalt            = []byte{0x60, 0x61, 0x62, 0x63, 0x64}
	fakeShortURLKeyHash = []byte{
		0x84, 0xe8, 0xad, 0xed, 0xc3, 0x70, 0xda, 0x76, 0x9a, 0xf3, 0xad, 0x8a, 0xe, 0x30, 0xf0, 0xa4, 0x1e, 0x19, 0xb5,
		0xb8, 0x80, 0xa3, 0xc8, 0x8c, 0x43, 0x60, 0xd2, 0x5c, 0x4, 0x5a, 0x5c, 0xbb,
	}
	fakeSaltedShortURLKeyHash = []byte{
		0xb0, 0xda, 0xc6, 0xeb, 0x5, 0x30, 0x56, 0x46, 0x65, 0x2, 0x10, 0x1e, 0x19, 0xb5, 0x8, 0x43, 0x56, 0xff, 0xd6,
		0xe6, 0x9d, 0xc, 0x23, 0xde, 0x7, 0x1a, 0xa5, 0xa0, 0xac, 0x72, 0x62, 0x99,
	}
	fakeShortURLAlias = "fake"
)

type mockedSaltProvider struct {
	mock.Mock
}

func (sp *mockedSaltProvider) Read(p []byte) (int, error) {
	args := sp.Called(p)
	return args.Int(0), args.Error(1)
}

type mockedAliasEncoder struct {
	mock.Mock
}

func (e *mockedAliasEncoder) Encode(in []byte, length int) (string, error) {
	args := e.Called(in, length)
	return args.String(0), args.Error(1)
}

func TestGenerate(t *testing.T) {
	t.Parallel()

	encoder := &mockedAliasEncoder{}
	saltProvider := &mockedSaltProvider{}
	generator := sha256AliasGenerator.New(encoder, saltProvider, fakeAliasLength)
	shorturl := &schema.ShortURL{
		OrganizationID: fakeOrganizationID,
		Domain:         fakeDomain,
		LongURL:        schema.LongURL{Assembled: fakeLongURL},
	}

	encoder.On("Encode", fakeShortURLKeyHash, fakeAliasLength).Return(fakeShortURLAlias, nil)

	alias, err := generator.Generate(shorturl)

	require.NoError(t, err)
	require.Equal(t, fakeShortURLAlias, alias)

	encoder.AssertExpectations(t)
	saltProvider.AssertExpectations(t)
}

func TestGenerateWithSalt(t *testing.T) {
	t.Parallel()

	encoder := &mockedAliasEncoder{}
	saltProvider := &mockedSaltProvider{}
	generator := sha256AliasGenerator.New(encoder, saltProvider, fakeAliasLength)
	shorturl := &schema.ShortURL{
		OrganizationID: fakeOrganizationID,
		Domain:         fakeDomain,
		LongURL:        schema.LongURL{Assembled: fakeLongURL},
	}

	encoder.On("Encode", fakeSaltedShortURLKeyHash, fakeAliasLength).Return(fakeShortURLAlias, nil)
	saltProvider.
		On("Read", mock.Anything).
		Run(func(args mock.Arguments) {
			arg := args.Get(0).([]byte)
			copy(arg, fakeSalt)
		}).
		Return(sha256AliasGenerator.SaltLengthBytes, nil)

	alias, err := generator.GenerateWithSalt(shorturl)

	require.NoError(t, err)
	require.Equal(t, fakeShortURLAlias, alias)

	encoder.AssertExpectations(t)
	saltProvider.AssertExpectations(t)
}

func TestGenerateWithSaltWhenErrorReadingSalt(t *testing.T) {
	t.Parallel()

	encoder := &mockedAliasEncoder{}
	saltProvider := &mockedSaltProvider{}
	generator := sha256AliasGenerator.New(encoder, saltProvider, fakeAliasLength)
	shorturl := &schema.ShortURL{
		OrganizationID: fakeOrganizationID,
		Domain:         fakeDomain,
		LongURL:        schema.LongURL{Assembled: fakeLongURL},
	}
	wantErr := gofakeit.Error()

	saltProvider.On("Read", mock.Anything).Return(0, wantErr)

	alias, err := generator.GenerateWithSalt(shorturl)

	require.ErrorIs(t, err, wantErr)
	require.Equal(t, "", alias)

	encoder.AssertExpectations(t)
	saltProvider.AssertExpectations(t)
}

package sha256

import (
	"crypto/sha256"
	"io"
	"slices"

	"github.com/FreibergVlad/url-shortener/short-url-management-service/internal/db/schema"
)

const SaltLengthBytes = 4

type ShortURLAliasEncoder interface {
	Encode(in []byte, length int) (string, error)
}

type AliasGenerator struct {
	encoder      ShortURLAliasEncoder
	saltProvider io.Reader
	aliasLegnth  int
}

func New(encoder ShortURLAliasEncoder, saltProvider io.Reader, aliasLength int) *AliasGenerator {
	return &AliasGenerator{
		encoder:      encoder,
		saltProvider: saltProvider,
		aliasLegnth:  aliasLength,
	}
}

func (g *AliasGenerator) Generate(shorturl *schema.ShortURL) (string, error) {
	return g.encode(g.generate(shorturl, nil))
}

func (g *AliasGenerator) GenerateWithSalt(shorturl *schema.ShortURL) (string, error) {
	salt := make([]byte, SaltLengthBytes)
	if _, err := g.saltProvider.Read(salt); err != nil {
		return "", err
	}

	return g.encode(g.generate(shorturl, salt))
}

func (g *AliasGenerator) generate(shorturl *schema.ShortURL, salt []byte) []byte {
	ukey := slices.Concat([]byte(shorturl.OrganizationID), []byte(shorturl.Domain), []byte(shorturl.LongURL.Assembled))
	if salt != nil {
		ukey = append(ukey, salt...)
	}

	hash := sha256.Sum256(ukey)
	return hash[:]
}

func (g *AliasGenerator) encode(alias []byte) (string, error) {
	return g.encoder.Encode(alias, g.aliasLegnth)
}

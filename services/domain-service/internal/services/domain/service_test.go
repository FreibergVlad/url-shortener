package domain_test

import (
	"context"
	"testing"

	"github.com/FreibergVlad/url-shortener/domain-service/internal/config"
	"github.com/FreibergVlad/url-shortener/domain-service/internal/services/domain"
	protoMessages "github.com/FreibergVlad/url-shortener/proto/pkg/domains/messages/v1"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListOrganizationDomain(t *testing.T) {
	t.Parallel()

	fqdn := gofakeit.DomainName()
	config := config.Config{PublicDomains: []string{fqdn}}
	domainService := domain.New(config)
	request := protoMessages.ListOrganizationDomainRequest{OrganizationId: gofakeit.UUID()}

	response, err := domainService.ListOrganizationDomain(context.TODO(), &request)

	require.NoError(t, err)

	assert.Equal(t, int64(1), response.Total)
	assert.Equal(t, []*protoMessages.Domain{{Fqdn: fqdn}}, response.Data)
}

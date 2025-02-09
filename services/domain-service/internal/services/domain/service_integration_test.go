package domain_test

import (
	"context"
	"testing"

	"github.com/FreibergVlad/url-shortener/domain-service/internal/config"
	testUtils "github.com/FreibergVlad/url-shortener/domain-service/internal/testing"
	domainServiceMessages "github.com/FreibergVlad/url-shortener/proto/pkg/domains/messages/v1"
	grpcUtils "github.com/FreibergVlad/url-shortener/shared/go/pkg/api/grpc/utils"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/testing/integration"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListOrganizationDomain_Integration(t *testing.T) {
	t.Parallel()
	integration.MaybeSkipIntegrationTest(t)

	server, teardown := testUtils.BootstrapServer()

	t.Cleanup(teardown)

	request := domainServiceMessages.ListOrganizationDomainRequest{OrganizationId: gofakeit.UUID()}

	response, err := server.DomainServiceClient.ListOrganizationDomain(
		grpcUtils.OutgoingContextWithUserID(context.Background(), gofakeit.UUID()),
		&request,
	)

	require.NoError(t, err)

	assert.Equal(t, publicDomains(), response.Data)
}

func TestListOrganizationDomainWhenUnauthenticated_Integration(t *testing.T) {
	t.Parallel()
	integration.MaybeSkipIntegrationTest(t)

	server, teardown := testUtils.BootstrapServer()

	t.Cleanup(teardown)

	request := domainServiceMessages.ListOrganizationDomainRequest{OrganizationId: gofakeit.UUID()}

	response, err := server.DomainServiceClient.ListOrganizationDomain(context.Background(), &request)

	assert.Nil(t, response)
	assert.ErrorContains(t, err, "unauthenticated")
}

func publicDomains() []*domainServiceMessages.Domain {
	config := config.New()
	domains := make([]*domainServiceMessages.Domain, len(config.PublicDomains))
	for i, domain := range config.PublicDomains {
		domains[i] = &domainServiceMessages.Domain{Fqdn: domain}
	}
	return domains
}

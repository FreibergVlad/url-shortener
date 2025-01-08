package domain_test

import (
	"context"
	"testing"

	"github.com/FreibergVlad/url-shortener/domain-service/internal/config"
	"github.com/FreibergVlad/url-shortener/domain-service/internal/server"
	domainServiceMessages "github.com/FreibergVlad/url-shortener/proto/pkg/domains/messages/v1"
	domainServiceProto "github.com/FreibergVlad/url-shortener/proto/pkg/domains/service/v1"
	permissionServiceProto "github.com/FreibergVlad/url-shortener/proto/pkg/permissions/service/v1"
	grpcUtils "github.com/FreibergVlad/url-shortener/shared/go/pkg/api/grpc/utils"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/testing/integration"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/testing/integration/mocks"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/test/bufconn"
)

func TestListOrganizationDomain_Integration(t *testing.T) {
	t.Parallel()
	integration.MaybeSkipIntegrationTest(t)

	config := config.New()
	client := bootstrapTestingGrpcServer(config, mocks.MockedPermissionServiceClient(true))
	request := domainServiceMessages.ListOrganizationDomainRequest{OrganizationId: gofakeit.UUID()}

	response, err := client.ListOrganizationDomain(grpcUtils.OutgoingContextWithUserID(context.Background(), gofakeit.UUID()), &request)

	assert.NoError(t, err)
	assert.Equal(t, publicDomains(config), response.Data)
}

func TestListOrganizationDomainWhenUnauthorized_Integration(t *testing.T) {
	t.Parallel()
	integration.MaybeSkipIntegrationTest(t)

	client := bootstrapTestingGrpcServer(config.New(), mocks.MockedPermissionServiceClient(false))
	request := domainServiceMessages.ListOrganizationDomainRequest{OrganizationId: gofakeit.UUID()}

	response, err := client.ListOrganizationDomain(grpcUtils.OutgoingContextWithUserID(context.Background(), gofakeit.UUID()), &request)

	assert.Nil(t, response)
	assert.ErrorContains(t, err, "is not allowed to perform")
}

func TestListOrganizationDomainWhenUnauthenticated_Integration(t *testing.T) {
	t.Parallel()
	integration.MaybeSkipIntegrationTest(t)

	client := bootstrapTestingGrpcServer(config.New(), mocks.MockedPermissionServiceClient(true))
	request := domainServiceMessages.ListOrganizationDomainRequest{OrganizationId: gofakeit.UUID()}

	response, err := client.ListOrganizationDomain(context.Background(), &request)

	assert.Nil(t, response)
	assert.ErrorContains(t, err, "unauthenticated")
}

func bootstrapTestingGrpcServer(config config.Config, permissionServiceClient permissionServiceProto.PermissionServiceClient) domainServiceProto.DomainServiceClient {
	listener := bufconn.Listen(1024 * 1024)
	server := server.Bootstrap(listener, permissionServiceClient, config)

	go server.Run()

	return domainServiceProto.NewDomainServiceClient(integration.BufnetGrpcClient(listener))
}

func publicDomains(config config.Config) []*domainServiceMessages.Domain {
	var domains []*domainServiceMessages.Domain
	for _, fqdn := range config.PublicDomains {
		domains = append(domains, &domainServiceMessages.Domain{Fqdn: fqdn})
	}
	return domains
}

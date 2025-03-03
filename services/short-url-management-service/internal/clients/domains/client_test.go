package domains_test

import (
	"context"
	"testing"

	domainServiceMessages "github.com/FreibergVlad/url-shortener/proto/pkg/domains/messages/v1"
	"github.com/FreibergVlad/url-shortener/short-url-management-service/internal/clients/domains"
	testUtils "github.com/FreibergVlad/url-shortener/short-url-management-service/internal/testing"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHasDomain(t *testing.T) {
	t.Parallel()

	fqdn := gofakeit.DomainName()
	grpcClient := testUtils.NewMockedDomainServiceClient([]string{fqdn})
	client := domains.NewServiceClient(grpcClient)

	err := client.HasDomain(context.Background(), gofakeit.UUID(), gofakeit.UUID(), fqdn)

	require.NoError(t, err)
}

func TestHasDomainWhenGrpcError(t *testing.T) {
	t.Parallel()

	wantErr := gofakeit.ErrorGRPC()
	grpcClient := &testUtils.MockedDomainServiceClient{}
	grpcClient.
		On("ListOrganizationDomain", mock.Anything, mock.Anything, mock.Anything).
		Return(&domainServiceMessages.ListOrganizationDomainResponse{}, wantErr)
	client := domains.NewServiceClient(grpcClient)

	err := client.HasDomain(context.Background(), gofakeit.UUID(), gofakeit.UUID(), gofakeit.DomainName())

	require.ErrorIs(t, err, wantErr)
}

func TestHasDomainWhenInvalidDomain(t *testing.T) {
	t.Parallel()

	client := domains.NewServiceClient(testUtils.NewMockedDomainServiceClient([]string{}))

	err := client.HasDomain(context.Background(), gofakeit.UUID(), gofakeit.UUID(), gofakeit.DomainName())

	require.ErrorIs(t, err, domains.ErrDomainNotAllowed)
}

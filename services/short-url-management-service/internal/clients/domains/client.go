package domains

import (
	"context"
	"errors"
	"fmt"
	"slices"

	domainServiceMessages "github.com/FreibergVlad/url-shortener/proto/pkg/domains/messages/v1"
	domainService "github.com/FreibergVlad/url-shortener/proto/pkg/domains/service/v1"
	grpcUtils "github.com/FreibergVlad/url-shortener/shared/go/pkg/api/grpc/utils"
)

var ErrDomainNotAllowed = errors.New("domain is not allowed to use")

type ServiceClient struct {
	grpcClient domainService.DomainServiceClient
}

func NewServiceClient(grpcClient domainService.DomainServiceClient) *ServiceClient {
	return &ServiceClient{
		grpcClient: grpcClient,
	}
}

func (c *ServiceClient) HasDomain(ctx context.Context, userID, organizationID, fqdn string) error {
	response, err := c.grpcClient.ListOrganizationDomain(
		grpcUtils.OutgoingContextWithUserID(ctx, userID),
		&domainServiceMessages.ListOrganizationDomainRequest{OrganizationId: organizationID},
	)
	if err != nil {
		return fmt.Errorf("error listing organization domains: %w", err)
	}

	hasDomainFun := func(d *domainServiceMessages.Domain) bool { return d.Fqdn == fqdn }
	if slices.ContainsFunc(response.Data, hasDomainFun) {
		return nil
	}

	return ErrDomainNotAllowed
}

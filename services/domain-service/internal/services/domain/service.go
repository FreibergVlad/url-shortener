package domain

import (
	"context"

	"github.com/FreibergVlad/url-shortener/domain-service/internal/config"
	protoMessages "github.com/FreibergVlad/url-shortener/proto/pkg/domains/messages/v1"
	protoService "github.com/FreibergVlad/url-shortener/proto/pkg/domains/service/v1"
)

type domainService struct {
	protoService.UnimplementedDomainServiceServer
	config config.Config
}

func New(config config.Config) *domainService {
	return &domainService{config: config}
}

func (s *domainService) ListOrganizationDomain(
	ctx context.Context,
	req *protoMessages.ListOrganizationDomainRequest,
) (*protoMessages.ListOrganizationDomainResponse, error) {
	var domains []*protoMessages.Domain
	for _, publicDomain := range s.config.PublicDomains {
		domains = append(domains, &protoMessages.Domain{Fqdn: publicDomain})
	}
	return &protoMessages.ListOrganizationDomainResponse{
		Data:  domains,
		Total: int32(len(s.config.PublicDomains)),
	}, nil
}

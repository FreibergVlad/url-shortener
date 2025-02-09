package domain

import (
	"context"

	"github.com/FreibergVlad/url-shortener/domain-service/internal/config"
	protoMessages "github.com/FreibergVlad/url-shortener/proto/pkg/domains/messages/v1"
	protoService "github.com/FreibergVlad/url-shortener/proto/pkg/domains/service/v1"
)

type Service struct {
	protoService.UnimplementedDomainServiceServer
	config config.Config
}

func New(config config.Config) *Service {
	return &Service{config: config}
}

func (s *Service) ListOrganizationDomain(
	_ context.Context,
	_ *protoMessages.ListOrganizationDomainRequest,
) (*protoMessages.ListOrganizationDomainResponse, error) {
	domains := make([]*protoMessages.Domain, len(s.config.PublicDomains))
	for i, publicDomain := range s.config.PublicDomains {
		domains[i] = &protoMessages.Domain{Fqdn: publicDomain}
	}

	return &protoMessages.ListOrganizationDomainResponse{
		Data:  domains,
		Total: int64(len(s.config.PublicDomains)),
	}, nil
}

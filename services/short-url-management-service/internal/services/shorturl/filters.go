package services

import (
	protoMessages "github.com/FreibergVlad/url-shortener/proto/pkg/shorturl/management/messages/v1"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/db/filters"
)

func OneShortURLFilter(domain, alias string) filters.Filter {
	return filters.And{
		Filters: []filters.Filter{
			filters.Equals[string]{Field: "domain", Value: domain},
			filters.Equals[string]{Field: "alias", Value: alias},
		},
	}
}

func OneNonDeletedShortURLFilter(domain, alias string) filters.Filter {
	return filters.And{
		Filters: []filters.Filter{
			OneShortURLFilter(domain, alias),
			filters.Nor{Filters: []filters.Filter{filters.Equals[string]{Field: "status", Value: protoMessages.ShortURLStatus_SHORT_URL_STATUS_DELETED.String()}}},
		},
	}
}

func OneNonDeletedShortURLByOrgIDFilter(orgID, domain, alias string) filters.Filter {
	return filters.And{
		Filters: []filters.Filter{
			filters.Equals[string]{Field: "organization_id", Value: orgID},
			filters.Equals[string]{Field: "domain", Value: domain},
			filters.Equals[string]{Field: "alias", Value: alias},
			filters.Nor{Filters: []filters.Filter{filters.Equals[string]{Field: "status", Value: protoMessages.ShortURLStatus_SHORT_URL_STATUS_DELETED.String()}}},
		},
	}
}

func ManyNonDeletedShortURLByOrgIDFilter(orgID string) filters.Filter {
	return filters.And{
		Filters: []filters.Filter{
			filters.Equals[string]{Field: "organization_id", Value: orgID},
			filters.Nor{Filters: []filters.Filter{filters.Equals[string]{Field: "status", Value: protoMessages.ShortURLStatus_SHORT_URL_STATUS_DELETED.String()}}},
		},
	}
}

package resolver

import (
	"context"

	"github.com/FreibergVlad/url-shortener/auth-service/internal/db/repositories/organizations"
	"github.com/FreibergVlad/url-shortener/auth-service/internal/db/repositories/users"
	"github.com/FreibergVlad/url-shortener/auth-service/internal/services/roles"
)

type permissionResolver struct {
	userRepository         users.Repository
	organizationRepository organizations.Repository
}

func New(
	userRepository users.Repository,
	organizationRepository organizations.Repository,
) *permissionResolver {
	return &permissionResolver{
		userRepository:         userRepository,
		organizationRepository: organizationRepository,
	}
}

func (r *permissionResolver) HasPermissions(ctx context.Context, scopes []string, userId string, organizationId *string) (bool, error) {
	user, err := r.userRepository.GetById(ctx, userId)
	if err != nil {
		return false, err
	}

	globalRole, ok := roles.GetRole(user.RoleID)
	if ok && globalRole.HasScopes(scopes) {
		return true, nil
	}

	if organizationId != nil {
		memberships, err := r.organizationRepository.ListOrganizationMembershipsByUserId(ctx, userId)
		if err != nil {
			return false, err
		}
		if membership := memberships.OrganizationMembership(*organizationId); membership != nil {
			orgRole, ok := roles.GetRole(membership.RoleID)
			if ok && orgRole.HasScopes(scopes) {
				return true, nil
			}
		}
	}

	return false, nil
}

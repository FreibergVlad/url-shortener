package resolver

import (
	"context"
	"fmt"

	"github.com/FreibergVlad/url-shortener/auth-service/internal/db/repositories/organizations"
	"github.com/FreibergVlad/url-shortener/auth-service/internal/db/repositories/users"
	"github.com/FreibergVlad/url-shortener/auth-service/internal/services/roles"
)

type PermissionResolver struct {
	userRepository         users.Repository
	organizationRepository organizations.Repository
}

func New(
	userRepository users.Repository,
	organizationRepository organizations.Repository,
) *PermissionResolver {
	return &PermissionResolver{
		userRepository:         userRepository,
		organizationRepository: organizationRepository,
	}
}

func (r *PermissionResolver) HasPermissions(
	ctx context.Context, scopes []string, userID string, organizationID *string,
) (bool, error) {
	user, err := r.userRepository.GetByID(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("error getting user: %w", err)
	}

	globalRole, ok := roles.GetRole(user.RoleID)
	if ok && globalRole.HasScopes(scopes) {
		return true, nil
	}

	if organizationID != nil {
		memberships, err := r.organizationRepository.ListOrganizationMembershipsByUserID(ctx, userID)
		if err != nil {
			return false, fmt.Errorf("failed to get user memberships %w", err)
		}
		if membership := memberships.OrganizationMembership(*organizationID); membership != nil {
			orgRole, ok := roles.GetRole(membership.RoleID)
			if ok && orgRole.HasScopes(scopes) {
				return true, nil
			}
		}
	}

	return false, nil
}

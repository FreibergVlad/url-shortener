package roles

import permissionServiceMessages "github.com/FreibergVlad/url-shortener/proto/pkg/permissions/messages/v1"

type (
	RoleId   = string
	RoleType = string
)

var (
	RoleTypeGlobal         = "global"
	RoleTypeOrganizational = "organizational"
)

var (
	RoleIdSuperAdmin  RoleId = "super_admin"
	RoleIdOwner       RoleId = "owner"
	RoleIdAdmin       RoleId = "admin"
	RoleIdMember      RoleId = "member"
	RoleIdProvisional RoleId = "provisional"
)

var (
	RoleSuperAdmin = role{
		ID:          RoleIdSuperAdmin,
		Name:        "Super Admin",
		Type:        RoleTypeGlobal,
		Description: "Full access across all organizations and short URLs.",
		Scopes: set(
			"organization:create",
			"organization-membership:list",
			"organization-invite:create",
			"domain:list",
			"short-url:read",
			"short-url:list",
			"short-url:create",
			"short-url:update",
			"short-url:delete",
			"me:read",
		),
	}
	RoleProvisional = role{
		ID:          RoleIdProvisional,
		Name:        "Provisional",
		Type:        RoleTypeGlobal,
		Description: "Automatically assigned role to the newly registered user.",
		Scopes: set(
			"organization:create",
			"organization-membership:list",
			"me:read",
		),
	}
	RoleOwner = role{
		ID:          RoleIdOwner,
		Name:        "Owner",
		Type:        RoleTypeOrganizational,
		Description: "Full access within a specific organization.",
		Scopes: set(
			"organization-invite:create",
			"domain:list",
			"short-url:read",
			"short-url:list",
			"short-url:create",
			"short-url:update",
			"short-url:delete",
		),
	}
	RoleAdmin = role{
		ID:          RoleIdAdmin,
		Name:        "Admin",
		Type:        RoleTypeOrganizational,
		Description: "Full short URL access within a specific organization.",
		Scopes: set(
			"organization-invite:create",
			"domain:list",
			"short-url:read",
			"short-url:list",
			"short-url:create",
			"short-url:update",
			"short-url:delete",
		),
	}
	RoleMember = role{
		ID:          RoleIdMember,
		Name:        "Member",
		Type:        RoleTypeOrganizational,
		Description: "Read-only access within a specific organization.",
		Scopes: set(
			"short-url:read",
			"short-url:list",
		),
	}
)

var roles map[RoleId]role = map[RoleId]role{
	RoleIdSuperAdmin:  RoleSuperAdmin,
	RoleIdOwner:       RoleOwner,
	RoleIdAdmin:       RoleAdmin,
	RoleIdMember:      RoleMember,
	RoleIdProvisional: RoleProvisional,
}

type role struct {
	ID          RoleId
	Name        string
	Type        RoleType
	Description string
	Scopes      map[string]struct{}
}

func (r role) HasScopes(scopes []string) bool {
	for _, scope := range scopes {
		if _, ok := r.Scopes[scope]; !ok {
			return false
		}
	}
	return true
}

func GetRole(id RoleId) (role, bool) {
	if role, ok := roles[id]; ok {
		return role, true
	}
	return role{}, false
}

func GetRoleProto(id RoleId) *permissionServiceMessages.Role {
	role, ok := GetRole(id)
	if !ok {
		return nil
	}
	return &permissionServiceMessages.Role{
		Id:          id,
		Name:        role.Name,
		Description: role.Description,
	}
}

func set(els ...string) map[string]struct{} {
	set := make(map[string]struct{})
	for _, el := range els {
		set[el] = struct{}{}
	}
	return set
}

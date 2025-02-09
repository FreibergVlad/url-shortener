package roles

import permissionServiceMessages "github.com/FreibergVlad/url-shortener/proto/pkg/permissions/messages/v1"

type (
	RoleID   = string
	RoleType = string
)

var (
	RoleTypeGlobal         = "global"
	RoleTypeOrganizational = "organizational"
)

var (
	RoleIDSuperAdmin  RoleID = "super_admin"
	RoleIDOwner       RoleID = "owner"
	RoleIDAdmin       RoleID = "admin"
	RoleIDMember      RoleID = "member"
	RoleIDProvisional RoleID = "provisional"
)

var (
	RoleSuperAdmin = Role{
		ID:          RoleIDSuperAdmin,
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
	RoleProvisional = Role{
		ID:          RoleIDProvisional,
		Name:        "Provisional",
		Type:        RoleTypeGlobal,
		Description: "Automatically assigned role to the newly registered user.",
		Scopes: set(
			"organization:create",
			"organization-membership:list",
			"me:read",
		),
	}
	RoleOwner = Role{
		ID:          RoleIDOwner,
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
	RoleAdmin = Role{
		ID:          RoleIDAdmin,
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
	RoleMember = Role{
		ID:          RoleIDMember,
		Name:        "Member",
		Type:        RoleTypeOrganizational,
		Description: "Read-only access within a specific organization.",
		Scopes: set(
			"short-url:read",
			"short-url:list",
		),
	}
)

var roles = map[RoleID]Role{
	RoleIDSuperAdmin:  RoleSuperAdmin,
	RoleIDOwner:       RoleOwner,
	RoleIDAdmin:       RoleAdmin,
	RoleIDMember:      RoleMember,
	RoleIDProvisional: RoleProvisional,
}

type Role struct {
	ID          RoleID
	Name        string
	Type        RoleType
	Description string
	Scopes      map[string]struct{}
}

func (r Role) HasScopes(scopes []string) bool {
	for _, scope := range scopes {
		if _, ok := r.Scopes[scope]; !ok {
			return false
		}
	}
	return true
}

func GetRole(id RoleID) (Role, bool) {
	if role, ok := roles[id]; ok {
		return role, true
	}
	return Role{}, false
}

func GetRoleProto(id RoleID) *permissionServiceMessages.Role {
	role, ok := GetRole(id)
	if !ok {
		return nil
	}
	return &permissionServiceMessages.Role{
		Id:          role.ID,
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

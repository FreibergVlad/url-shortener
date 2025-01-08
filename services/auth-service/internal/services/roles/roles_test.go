package roles_test

import (
	"testing"

	"github.com/FreibergVlad/url-shortener/auth-service/internal/services/roles"
	"github.com/stretchr/testify/assert"
)

func TestHasScopes(t *testing.T) {
	t.Parallel()

	role := roles.RoleSuperAdmin
	assert.True(t, role.HasScopes([]string{"short-url:read", "short-url:update"}))
	assert.False(t, role.HasScopes([]string{"short-url:read", "invalid-scope"}))
}

func TestGetRole(t *testing.T) {
	t.Parallel()

	actualRole, ok := roles.GetRole(roles.RoleIdSuperAdmin)
	assert.True(t, ok)
	assert.Equal(t, roles.RoleSuperAdmin, actualRole)

	_, ok = roles.GetRole(roles.RoleId("invalid-role-id"))
	assert.False(t, ok)
}

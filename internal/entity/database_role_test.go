package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckRoleName(t *testing.T) {
	databaseRole := &DatabaseRole{Name: UserRO}
	assert.True(t, databaseRole.IsUserRO())

	databaseRole = &DatabaseRole{Name: Developer}
	assert.True(t, databaseRole.IsDeveloper())

	databaseRole = &DatabaseRole{Name: DevOps}
	assert.True(t, databaseRole.IsDevOps())

	databaseRole = &DatabaseRole{Name: Application}
	assert.True(t, databaseRole.IsApplication())
}

func TestValidateRoleName(t *testing.T) {
	assert.True(t, ValidateRoleName("user_ro"))
	assert.True(t, ValidateRoleName("developer"))
	assert.True(t, ValidateRoleName("devops"))
	assert.True(t, ValidateRoleName("application"))
	assert.False(t, ValidateRoleName("invalid"))
}

func TestCheckRoleApplication(t *testing.T) {
	assert.True(t, CheckRoleApplication("application"))
	assert.False(t, CheckRoleApplication("developer"))
	assert.False(t, CheckRoleApplication("invalid"))
}

func TestCheckRole(t *testing.T) {
	assert.True(t, CheckRole(string(UserRO), UserRO))
	assert.True(t, CheckRole(string(Developer), Developer))
	assert.True(t, CheckRole(string(DevOps), DevOps))
	assert.True(t, CheckRole(string(Application), Application))
	assert.False(t, CheckRole("invalid", UserRO))
	assert.False(t, CheckRole(string(UserRO), Developer))
}

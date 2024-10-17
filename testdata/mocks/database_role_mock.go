package mocks

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/zgsolucoes/zg-data-guard/internal/entity"
)

type DatabaseRoleStorageMock struct {
	mock.Mock
}

func (m *DatabaseRoleStorageMock) FindAll() ([]*entity.DatabaseRole, error) {
	args := m.Called()
	return args.Get(0).([]*entity.DatabaseRole), args.Error(1)
}

func (m *DatabaseRoleStorageMock) FindByID(id string) (*entity.DatabaseRole, error) {
	args := m.Called(id)
	return args.Get(0).(*entity.DatabaseRole), args.Error(1)
}

func BuildRolesList() []*entity.DatabaseRole {
	role := BuildDeveloperRole()
	role2 := BuildReadOnlyRole()
	roles := []*entity.DatabaseRole{role, role2}
	return roles
}

func BuildDeveloperRole() *entity.DatabaseRole {
	return &entity.DatabaseRole{
		ID:          uuid.MustParse("1eb93da6-e739-4396-902f-19f79aa74e39"),
		Name:        "developer",
		DisplayName: "Developer",
		Description: "Role for developers",
		ReadOnly:    false,
	}
}

func BuildReadOnlyRole() *entity.DatabaseRole {
	return &entity.DatabaseRole{
		ID:          uuid.MustParse("cd7f93a4-a2ff-41db-9ad2-6dd67dd285c7"),
		Name:        "user_ro",
		DisplayName: "Read Only User",
		Description: "Role for read only users",
		ReadOnly:    true,
	}
}

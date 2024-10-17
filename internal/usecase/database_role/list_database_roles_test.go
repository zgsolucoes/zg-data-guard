package role

import (
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/zgsolucoes/zg-data-guard/internal/entity"
	"github.com/zgsolucoes/zg-data-guard/testdata/mocks"
)

func TestGivenAnErrorInDb_WhenExecuteList_ThenShouldReturnError(t *testing.T) {
	roleStorage := new(mocks.DatabaseRoleStorageMock)
	roleStorage.On("FindAll").Return([]*entity.DatabaseRole{}, sql.ErrConnDone).Once()

	uc := NewListDatabaseRolesUseCase(roleStorage)
	obtainedRoles, err := uc.Execute()

	assert.Error(t, err, "error expected when some error in db")
	assert.EqualError(t, err, sql.ErrConnDone.Error())
	assert.Nil(t, obtainedRoles)
	assert.Equal(t, len(obtainedRoles), 0, "0 database roles expected")
	roleStorage.AssertNumberOfCalls(t, "FindAll", 1)
}

func TestGivenSomeRoles_WhenExecuteList_ThenShouldListAll(t *testing.T) {
	role := &entity.DatabaseRole{
		ID:              uuid.MustParse("e5dcca5c-163b-4296-8b36-f50db233589f"),
		Name:            "developer",
		DisplayName:     "Developer",
		Description:     "Developer role",
		ReadOnly:        false,
		CreatedAt:       time.Now(),
		CreatedByUserID: "1",
	}
	role2 := &entity.DatabaseRole{
		ID:              uuid.MustParse("918f0198-7c6e-4b75-a4c6-f84ad6c2d090"),
		Name:            "user_ro",
		DisplayName:     "User Read Only",
		Description:     "Read Only role",
		ReadOnly:        true,
		CreatedAt:       time.Now(),
		CreatedByUserID: "1",
	}
	role3 := &entity.DatabaseRole{
		ID:              uuid.MustParse("93afa9f8-e43c-48c8-90db-78539d068557"),
		Name:            "devops",
		DisplayName:     "DevOps",
		Description:     "DevOps role",
		ReadOnly:        false,
		CreatedAt:       time.Now(),
		CreatedByUserID: "1",
	}
	roles := []*entity.DatabaseRole{role, role2, role3}
	roleStorage := new(mocks.DatabaseRoleStorageMock)
	roleStorage.On("FindAll").Return(roles, nil).Once()

	uc := NewListDatabaseRolesUseCase(roleStorage)
	obtainedRoles, err := uc.Execute()

	assert.NoError(t, err, "no error expected ")
	assert.Equal(t, len(obtainedRoles), len(roles), "3 database roles expected")
	roleStorage.AssertNumberOfCalls(t, "FindAll", 1)
}

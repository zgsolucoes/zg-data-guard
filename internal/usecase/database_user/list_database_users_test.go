package dbuser

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	"github.com/zgsolucoes/zg-data-guard/testdata/mocks"
)

func TestGivenAnErrorInDb_WhenExecuteList_ThenShouldReturnError(t *testing.T) {
	dbUserStorage := new(mocks.DatabaseUserStorageMock)
	dbUserStorage.On("FindAllDTOs", []string{}).Return([]*dto.DatabaseUserOutputDTO{}, sql.ErrConnDone).Once()

	uc := NewListDatabaseUsersUseCase(dbUserStorage)
	obtainedUsers, err := uc.Execute(false)

	assert.Error(t, err, "error expected when some error in db")
	assert.EqualError(t, err, sql.ErrConnDone.Error())
	assert.Nil(t, obtainedUsers)
	assert.Equal(t, len(obtainedUsers), 0, "0 database users expected")
	dbUserStorage.AssertNumberOfCalls(t, "FindAllDTOs", 1)
}

func TestGivenSomeUsers_WhenExecuteList_ThenShouldListAll(t *testing.T) {
	users := mocks.BuildDbUserDTOList()
	dbUserStorage := new(mocks.DatabaseUserStorageMock)
	dbUserStorage.On("FindAllDTOsEnabled").Return(users, nil).Once()

	uc := NewListDatabaseUsersUseCase(dbUserStorage)
	obtainedUsers, err := uc.Execute(true)

	assert.NoError(t, err, "no error expected ")
	assert.Equal(t, len(obtainedUsers), len(users), "2 database users expected")
	dbUserStorage.AssertNumberOfCalls(t, "FindAllDTOsEnabled", 1)
}

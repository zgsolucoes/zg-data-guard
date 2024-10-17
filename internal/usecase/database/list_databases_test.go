package database

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	"github.com/zgsolucoes/zg-data-guard/testdata/mocks"
)

func TestGivenAnErrorInDb_WhenExecuteListDatabases_ThenShouldReturnError(t *testing.T) {
	databaseStorage := new(mocks.DatabaseStorageMock)
	databaseStorage.On("FindAllDTOs", mocks.EcosystemId, mocks.DatabaseInstanceId).Return([]*dto.DatabaseOutputDTO{}, sql.ErrConnDone).Once()

	uc := NewListDatabasesUseCase(databaseStorage)
	databasesObtained, err := uc.Execute(mocks.EcosystemId, mocks.DatabaseInstanceId)

	assert.Error(t, err, "error expected when some error in db")
	assert.EqualError(t, err, sql.ErrConnDone.Error())
	assert.Nil(t, databasesObtained)
	assert.Equal(t, len(databasesObtained), 0, "0 database expected")
	databaseStorage.AssertNumberOfCalls(t, "FindAllDTOs", 1)
}

func TestGivenSomeDbs_WhenExecuteListDatabases_ThenShouldListAllDatabases(t *testing.T) {
	databaseDtos := mocks.BuildDatabaseDTOList()
	databaseStorage := new(mocks.DatabaseStorageMock)
	databaseStorage.On("FindAllDTOs", mocks.EcosystemId, mocks.DatabaseInstanceId).Return(databaseDtos, nil).Once()

	uc := NewListDatabasesUseCase(databaseStorage)
	databasesObtained, err := uc.Execute(mocks.EcosystemId, mocks.DatabaseInstanceId)

	assert.NoError(t, err, "no error expected")
	assert.Equal(t, len(databaseDtos), len(databasesObtained), "3 databases expected")
	databaseStorage.AssertNumberOfCalls(t, "FindAllDTOs", 1)
}

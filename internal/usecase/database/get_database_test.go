package database

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	"github.com/zgsolucoes/zg-data-guard/testdata/mocks"
)

func TestGivenANonexistentId_WhenExecuteGetDatabase_ThenShouldReturnError(t *testing.T) {
	databaseStorage := new(mocks.DatabaseStorageMock)
	uc := NewGetDatabaseUseCase(databaseStorage)

	databaseStorage.On("FindDTOByID", mocks.DatabaseID).Return(&dto.DatabaseOutputDTO{}, sql.ErrNoRows)
	_, err := uc.Execute(mocks.DatabaseID)
	assert.Error(t, err)
	assert.EqualError(t, err, ErrDatabaseNotFound.Error())
	databaseStorage.AssertNumberOfCalls(t, "FindDTOByID", 1)
}

func TestGivenAnErrorInDb_WhenExecuteGetDatabase_ThenShouldReturnError(t *testing.T) {
	databaseStorage := new(mocks.DatabaseStorageMock)
	uc := NewGetDatabaseUseCase(databaseStorage)

	databaseStorage.On("FindDTOByID", mocks.DatabaseID).Return(&dto.DatabaseOutputDTO{}, sql.ErrTxDone)
	_, err := uc.Execute(mocks.DatabaseID)
	assert.Error(t, err)
	assert.EqualError(t, err, sql.ErrTxDone.Error())
	databaseStorage.AssertNumberOfCalls(t, "FindDTOByID", 1)
}

func TestGivenAValidId_WhenExecuteGet_ThenShouldReturnDatabase(t *testing.T) {
	databaseStorage := new(mocks.DatabaseStorageMock)
	uc := NewGetDatabaseUseCase(databaseStorage)

	database := mocks.BuildDatabaseDTOExample()
	databaseStorage.On("FindDTOByID", mocks.DatabaseID).Return(database, nil)
	dbDTO, err := uc.Execute(mocks.DatabaseID)
	assert.NoError(t, err)
	assert.NotNil(t, dbDTO)
	assert.Equal(t, database.ID, dbDTO.ID)
	assert.Equal(t, database.Name, dbDTO.Name)
	assert.Equal(t, database.CurrentSize, dbDTO.CurrentSize)
	assert.Equal(t, database.DatabaseInstanceID, dbDTO.DatabaseInstanceID)
	assert.Equal(t, database.DatabaseInstanceName, dbDTO.DatabaseInstanceName)
	assert.Equal(t, database.EcosystemID, dbDTO.EcosystemID)
	assert.Equal(t, database.EcosystemName, dbDTO.EcosystemName)
	assert.Equal(t, database.DatabaseTechnologyID, dbDTO.DatabaseTechnologyID)
	assert.Equal(t, database.DatabaseTechnologyName, dbDTO.DatabaseTechnologyName)
	assert.Equal(t, database.DatabaseTechnologyVersion, dbDTO.DatabaseTechnologyVersion)
	assert.Equal(t, database.Enabled, dbDTO.Enabled)
	assert.Equal(t, database.Description, dbDTO.Description)
	assert.Equal(t, database.CreatedByUserID, dbDTO.CreatedByUserID)
	assert.Equal(t, database.CreatedByUser, dbDTO.CreatedByUser)
	assert.Equal(t, database.CreatedAt, dbDTO.CreatedAt)
	assert.Equal(t, database.UpdatedAt, dbDTO.UpdatedAt)
	assert.Equal(t, database.LastDatabaseSync, dbDTO.LastDatabaseSync)
	assert.Equal(t, database.DisabledAt, dbDTO.DisabledAt)
	databaseStorage.AssertNumberOfCalls(t, "FindDTOByID", 1)
}

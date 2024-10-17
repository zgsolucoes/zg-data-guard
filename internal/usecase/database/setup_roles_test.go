package database

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	"github.com/zgsolucoes/zg-data-guard/internal/entity"
	"github.com/zgsolucoes/zg-data-guard/testdata/mocks"
)

func TestGivenAnErrorInDb_WhenExecuteSetupRoles_ThenShouldReturnError(t *testing.T) {
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	databaseStorage := new(mocks.DatabaseStorageMock)
	databaseStorage.On("FindAllEnabled", "").Return([]*entity.Database{}, sql.ErrConnDone).Once()

	uc := NewSetupRolesInDatabasesUseCase(dbInstanceStorage, databaseStorage)
	outputs, err := uc.Execute(dto.SetupRolesInputDTO{}, mocks.UserID)

	assert.Error(t, err, "error expected when some error in db")
	assert.EqualError(t, err, sql.ErrConnDone.Error())
	assert.Nil(t, outputs)
	databaseStorage.AssertNumberOfCalls(t, "FindAllEnabled", 1)
}

func TestGivenAnErrorInDbWhileFetchingDatabases_WhenExecuteSetupRoles_ThenShouldReturnError(t *testing.T) {
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	databaseStorage := new(mocks.DatabaseStorageMock)
	databaseStorage.On("FindAll", "1", mock.AnythingOfType("[]string")).Return([]*entity.Database{}, sql.ErrConnDone).Once()

	uc := NewSetupRolesInDatabasesUseCase(dbInstanceStorage, databaseStorage)
	outputs, err := uc.Execute(dto.SetupRolesInputDTO{DatabaseInstanceID: "1"}, mocks.UserID)

	assert.Error(t, err, "error expected when some error in db")
	assert.EqualError(t, err, sql.ErrConnDone.Error())
	assert.Nil(t, outputs)
	databaseStorage.AssertNumberOfCalls(t, "FindAll", 1)
}

func TestGivenInputWithIdsNotExistent_WhenExecuteSetupRoles_ThenShouldReturnError(t *testing.T) {
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	databaseStorage := new(mocks.DatabaseStorageMock)
	databaseStorage.On("FindAll", "", []string{"1", "2"}).Return([]*entity.Database{}, sql.ErrNoRows).Once()

	uc := NewSetupRolesInDatabasesUseCase(dbInstanceStorage, databaseStorage)
	outputs, err := uc.Execute(dto.SetupRolesInputDTO{DatabasesIDs: []string{"1", "2"}}, mocks.UserID)

	assert.Error(t, err, "error expected when no databases found")
	assert.EqualError(t, err, ErrNoDatabasesFound.Error())
	assert.Nil(t, outputs)
	databaseStorage.AssertNumberOfCalls(t, "FindAll", 1)
}

func TestGivenAnErrorInDbWhileFetchingInstanceDTO_WhenExecuteSetupRoles_ThenShouldReturnSuccessFalse(t *testing.T) {
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	databaseStorage := new(mocks.DatabaseStorageMock)
	databasesToProcess := mocks.BuildDatabaseListSameInstanceAndOnlyEnabled()
	databaseStorage.On("FindAllEnabled", "").Return(databasesToProcess, nil).Once()
	dbInstanceStorage.On("FindDTOByID", mocks.DatabaseInstanceId).Return(&dto.DatabaseInstanceOutputDTO{}, sql.ErrConnDone).Once()

	uc := NewSetupRolesInDatabasesUseCase(dbInstanceStorage, databaseStorage)
	outputs, err := uc.Execute(dto.SetupRolesInputDTO{}, mocks.UserID)

	assert.NoError(t, err, "no error expected")
	assert.NotNil(t, outputs)
	assert.Len(t, outputs, len(databasesToProcess))
	for _, output := range outputs {
		assert.False(t, output.Success)
		assert.Equal(t, sql.ErrConnDone.Error(), output.Message)
	}
	dbInstanceStorage.AssertNumberOfCalls(t, "FindDTOByID", 1)
	databaseStorage.AssertNumberOfCalls(t, "FindAllEnabled", 1)
}

func TestGivenADisabledInstance_WhenExecuteSetupRoles_ThenShouldReturnSuccessFalse(t *testing.T) {
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	databaseStorage := new(mocks.DatabaseStorageMock)
	databasesToProcess := mocks.BuildDatabaseListSameInstanceAndOnlyEnabled()
	databaseStorage.On("FindAllEnabled", "").Return(databasesToProcess, nil).Once()
	databaseInstance := mocks.BuildAzInstanceDTO()
	databaseInstance.Enabled = false
	dbInstanceStorage.On("FindDTOByID", databaseInstance.ID).Return(databaseInstance, nil).Once()

	uc := NewSetupRolesInDatabasesUseCase(dbInstanceStorage, databaseStorage)
	outputs, err := uc.Execute(dto.SetupRolesInputDTO{}, mocks.UserID)

	assert.NoError(t, err, "no error expected")
	assert.NotNil(t, outputs, "outputs expected all with success false")
	assert.Len(t, outputs, len(databasesToProcess))
	for _, output := range outputs {
		assert.False(t, output.Success)
		assert.Equal(t, "this database belongs to a disabled instance.", output.Message)
	}
	dbInstanceStorage.AssertNumberOfCalls(t, "FindDTOByID", 1)
	databaseStorage.AssertNumberOfCalls(t, "FindAllEnabled", 1)
}

func TestGivenAnErrorInDbWhileUpdateDatabase_WhenExecuteSetupRoles_ThenShouldReturnSuccessFalse(t *testing.T) {
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	databaseStorage := new(mocks.DatabaseStorageMock)
	databasesToProcess := mocks.BuildDatabaseListSameInstanceAndOnlyEnabled()
	databasesQty := len(databasesToProcess)
	databaseStorage.On("FindAllEnabled", "").Return(databasesToProcess, nil).Once()
	databaseStorage.On("Update", mock.Anything).Return(sql.ErrConnDone).Times(databasesQty)
	databaseInstance := mocks.BuildAzInstanceDTO()
	dbInstanceStorage.On("FindDTOByID", databaseInstance.ID).Return(databaseInstance, nil).Once()

	uc := NewSetupRolesInDatabasesUseCase(dbInstanceStorage, databaseStorage)
	outputs, err := uc.Execute(dto.SetupRolesInputDTO{}, mocks.UserID)

	assert.NoError(t, err, "no error expected")
	assert.NotNil(t, outputs, "outputs expected all with success false")
	assert.Len(t, outputs, databasesQty)
	for _, output := range outputs {
		assert.False(t, output.Success)
		assert.Equal(t, sql.ErrConnDone.Error(), output.Message)
	}
	dbInstanceStorage.AssertNumberOfCalls(t, "FindDTOByID", 1)
	databaseStorage.AssertNumberOfCalls(t, "FindAllEnabled", 1)
	databaseStorage.AssertNumberOfCalls(t, "Update", databasesQty)
}

func TestGivenInstanceWithConnectorNotImplemented_WhenExecuteSetupRoles_ThenShouldReturnSuccessFalse(t *testing.T) {
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	databaseStorage := new(mocks.DatabaseStorageMock)
	dbInstanceMysql := mocks.BuildConnectorNotImplementedInstance()
	databaseMysql := &entity.Database{
		Name:               "mysql-db",
		DatabaseInstanceID: dbInstanceMysql.ID,
		Enabled:            true,
	}
	databasesToProcess := []*entity.Database{databaseMysql}
	databasesQty := len(databasesToProcess)
	databaseStorage.On("FindAllEnabled", "").Return(databasesToProcess, nil).Once()
	dbInstanceStorage.On("FindDTOByID", dbInstanceMysql.ID).Return(dbInstanceMysql, nil).Once()

	uc := NewSetupRolesInDatabasesUseCase(dbInstanceStorage, databaseStorage)
	outputs, err := uc.Execute(dto.SetupRolesInputDTO{}, mocks.UserID)
	notImplementedConnectorOutput := outputs[0]

	assert.NoError(t, err, "no error expected")
	assert.NotNil(t, outputs)
	assert.Len(t, outputs, databasesQty)
	assert.False(t, notImplementedConnectorOutput.Success)
	assert.Equal(t, dbInstanceMysql.ID, notImplementedConnectorOutput.DatabaseInstanceID)
	assert.Equal(t, dbInstanceMysql.Name, notImplementedConnectorOutput.Instance)
	assert.Equal(t, databaseMysql.Name, notImplementedConnectorOutput.DatabaseName)
	assert.Equal(t, "the database technology 'mysql' don't have a connector implemented", notImplementedConnectorOutput.Message)
	databaseStorage.AssertNumberOfCalls(t, "FindAllEnabled", 1)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindDTOByID", 1)
}

func TestGivenSomeDatabases_WhenExecuteSetupRolesByIds_ThenShouldSetupOnlyValidOnes(t *testing.T) {
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	databaseStorage := new(mocks.DatabaseStorageMock)
	databasesToProcess := mocks.BuildDatabaseListMixedScenarios()
	var databaseIds []string
	for _, database := range databasesToProcess {
		databaseIds = append(databaseIds, database.ID.String())
	}
	databasesQty := len(databasesToProcess)
	databaseStorage.On("FindAll", "", databaseIds).Return(databasesToProcess, nil).Once()
	databaseStorage.On("Update", mock.Anything).Return(nil).Times(2)
	dbInstanceStorage.On("FindDTOByID", mocks.QAInstanceId).Return(mocks.BuildQAInstanceEnabled(), nil).Once()
	dbInstanceStorage.On("FindDTOByID", mocks.DatabaseInstanceId).Return(mocks.BuildAzInstanceDTO(), nil).Once()
	dbInstanceStorage.On("FindDTOByID", mocks.DummyErrorInstanceId).Return(mocks.BuildDummyErrorInstance(), nil).Once()

	uc := NewSetupRolesInDatabasesUseCase(dbInstanceStorage, databaseStorage)
	outputs, err := uc.Execute(dto.SetupRolesInputDTO{DatabasesIDs: databaseIds}, mocks.UserID)

	successCount := 0
	failureCount := 0
	for _, output := range outputs {
		if output.Success {
			assert.Contains(t, []string{"112c05c4-8e33-49db-b41d-9c3b8d9e2676", "431d710e-b519-4490-8574-ea9fb84a8d33"}, output.DatabaseID)
			assert.Contains(t, []string{mocks.QAInstanceId, mocks.DatabaseInstanceId}, output.DatabaseInstanceID)
			successCount++
		} else {
			assert.Contains(t, []string{"0f0e9ec3-f5c1-48fd-b16b-702a0e9f1804", "b55a4c72-159c-449f-bad4-0080944eb4da"}, output.DatabaseID)
			assert.Contains(t, []string{mocks.DatabaseInstanceId, mocks.DummyErrorInstanceId}, output.DatabaseInstanceID)
			failureCount++
		}
	}
	assert.NoError(t, err, "no error expected")
	assert.NotNil(t, outputs)
	assert.Len(t, outputs, databasesQty)
	assert.Equal(t, 2, successCount)
	assert.Equal(t, 2, failureCount)
	databaseStorage.AssertNumberOfCalls(t, "FindAll", 1)
	databaseStorage.AssertNumberOfCalls(t, "Update", 2)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindDTOByID", 3)
}

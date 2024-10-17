package database

import (
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/zgsolucoes/zg-data-guard/internal/database/connector"
	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	"github.com/zgsolucoes/zg-data-guard/internal/entity"
	"github.com/zgsolucoes/zg-data-guard/testdata/mocks"
)

func TestGivenAnErrorInDb_WhenExecuteSyncDatabases_ThenShouldReturnError(t *testing.T) {
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	databaseStorage := new(mocks.DatabaseStorageMock)
	dbInstanceStorage.On("FindAllDTOsEnabled", "", "").Return([]*dto.DatabaseInstanceOutputDTO{}, sql.ErrConnDone).Once()

	uc := NewSyncDatabasesUseCase(dbInstanceStorage, databaseStorage)
	outputs, err := uc.Execute(dto.SyncDatabasesInputDTO{}, mocks.UserID)

	assert.Error(t, err, "error expected when some error in db")
	assert.EqualError(t, err, sql.ErrConnDone.Error())
	assert.Nil(t, outputs)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindAllDTOsEnabled", 1)
}

func TestGivenAnErrorInDbWhileFetchingInstances_WhenExecuteSyncDatabases_ThenShouldReturnError(t *testing.T) {
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	databaseStorage := new(mocks.DatabaseStorageMock)
	dbInstanceStorage.On("FindAllDTOs", "", "", []string{"1", "2"}).Return([]*dto.DatabaseInstanceOutputDTO{}, sql.ErrConnDone).Once()

	uc := NewSyncDatabasesUseCase(dbInstanceStorage, databaseStorage)
	outputs, err := uc.Execute(dto.SyncDatabasesInputDTO{DatabaseInstancesIDs: []string{"1", "2"}}, mocks.UserID)

	assert.Error(t, err, "error expected when some error in db")
	assert.EqualError(t, err, sql.ErrConnDone.Error())
	assert.Nil(t, outputs)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindAllDTOs", 1)
}

func TestGivenAnErrorInDbWhileFetchingDatabases_WhenExecuteSyncDatabases_ThenShouldReturnSuccessFalse(t *testing.T) {
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	databaseStorage := new(mocks.DatabaseStorageMock)
	dbInstances := []*dto.DatabaseInstanceOutputDTO{mocks.BuildAzInstanceDTO()}
	dbInstanceStorage.On("FindAllDTOsEnabled", "", "").Return(dbInstances, nil).Once()
	databaseStorage.On("FindAll", dbInstances[0].ID, []string{}).Return([]*entity.Database{}, sql.ErrConnDone).Once()

	uc := NewSyncDatabasesUseCase(dbInstanceStorage, databaseStorage)
	outputs, err := uc.Execute(dto.SyncDatabasesInputDTO{}, mocks.UserID)
	unsuccessfulOutput := outputs[0]

	assert.NoError(t, err, "no error expected")
	assert.NotNil(t, outputs)
	assert.Len(t, outputs, 1)
	assert.False(t, unsuccessfulOutput.Success)
	assert.Equal(t, dbInstances[0].ID, unsuccessfulOutput.DatabaseInstanceID)
	assert.Equal(t, dbInstances[0].Name, unsuccessfulOutput.Instance)
	assert.Equal(t, sql.ErrConnDone.Error(), unsuccessfulOutput.Message)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindAllDTOsEnabled", 1)
	databaseStorage.AssertNumberOfCalls(t, "FindAll", 1)
}

func TestGivenAnErrorInDbWhileSaveDatabase_WhenExecuteSyncDatabases_ThenShouldReturnSuccessFalse(t *testing.T) {
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	databaseStorage := new(mocks.DatabaseStorageMock)
	dbInstances := []*dto.DatabaseInstanceOutputDTO{mocks.BuildAzInstanceDTO()}
	dbInstanceStorage.On("FindAllDTOsEnabled", "", "").Return(dbInstances, nil).Once()
	databases := mocks.BuildDatabaseList()
	databaseStorage.On("FindAll", dbInstances[0].ID, []string{}).Return(databases, nil).Once()
	databaseStorage.On("Save", mock.Anything).Return(sql.ErrConnDone).Once()
	databaseStorage.On("Update", mock.Anything).Return(nil).Twice()

	uc := NewSyncDatabasesUseCase(dbInstanceStorage, databaseStorage)
	outputs, err := uc.Execute(dto.SyncDatabasesInputDTO{}, mocks.UserID)
	unsuccessfulOutput := outputs[0]

	assert.NoError(t, err, "no error expected")
	assert.NotNil(t, outputs)
	assert.Len(t, outputs, 1)
	assert.False(t, unsuccessfulOutput.Success)
	assert.Equal(t, dbInstances[0].ID, unsuccessfulOutput.DatabaseInstanceID)
	assert.Equal(t, dbInstances[0].Name, unsuccessfulOutput.Instance)
	assert.Equal(t, sql.ErrConnDone.Error(), unsuccessfulOutput.Message)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindAllDTOsEnabled", 1)
	databaseStorage.AssertNumberOfCalls(t, "FindAll", 1)
	databaseStorage.AssertNumberOfCalls(t, "Save", 1)
	databaseStorage.AssertNumberOfCalls(t, "Update", 2)
}

func TestGivenAnErrorInDbWhileUpdateDatabase_WhenExecuteSyncDatabases_ThenShouldReturnSuccessFalse(t *testing.T) {
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	databaseStorage := new(mocks.DatabaseStorageMock)
	dbInstances := []*dto.DatabaseInstanceOutputDTO{mocks.BuildAzInstanceDTO()}
	dbInstanceStorage.On("FindAllDTOsEnabled", "", "").Return(dbInstances, nil).Once()
	databases := mocks.BuildDatabaseList()
	databaseStorage.On("FindAll", dbInstances[0].ID, []string{}).Return(databases, nil).Once()
	databaseStorage.On("Update", mock.Anything).Return(sql.ErrConnDone).Once()

	uc := NewSyncDatabasesUseCase(dbInstanceStorage, databaseStorage)
	outputs, err := uc.Execute(dto.SyncDatabasesInputDTO{}, mocks.UserID)
	unsuccessfulOutput := outputs[0]

	assert.NoError(t, err, "no error expected")
	assert.NotNil(t, outputs)
	assert.Len(t, outputs, 1)
	assert.False(t, unsuccessfulOutput.Success)
	assert.Equal(t, dbInstances[0].ID, unsuccessfulOutput.DatabaseInstanceID)
	assert.Equal(t, dbInstances[0].Name, unsuccessfulOutput.Instance)
	assert.Equal(t, sql.ErrConnDone.Error(), unsuccessfulOutput.Message)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindAllDTOsEnabled", 1)
	databaseStorage.AssertNumberOfCalls(t, "FindAll", 1)
	databaseStorage.AssertNumberOfCalls(t, "Update", 1)
}

func TestGivenAnErrorInDbWhileFetchingInstanceToUpdate_WhenExecuteSyncDatabases_ThenShouldReturnSuccessFalse(t *testing.T) {
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	databaseStorage := new(mocks.DatabaseStorageMock)
	dbInstances := []*dto.DatabaseInstanceOutputDTO{mocks.BuildAzInstanceDTO()}
	dbInstanceStorage.On("FindAllDTOsEnabled", "", "").Return(dbInstances, nil).Once()
	dbInstanceStorage.On("FindByID", dbInstances[0].ID).Return(&entity.DatabaseInstance{}, sql.ErrConnDone).Once()
	databases := mocks.BuildDatabaseList()
	databaseStorage.On("FindAll", dbInstances[0].ID, []string{}).Return(databases, nil).Once()
	databaseStorage.On("Save", mock.Anything).Return(nil).Once()
	databaseStorage.On("Update", mock.Anything).Return(nil).Times(3)

	uc := NewSyncDatabasesUseCase(dbInstanceStorage, databaseStorage)
	outputs, err := uc.Execute(dto.SyncDatabasesInputDTO{}, mocks.UserID)
	unsuccessfulOutput := outputs[0]

	assert.NoError(t, err, "no error expected")
	assert.NotNil(t, outputs)
	assert.Len(t, outputs, 1)
	assert.False(t, unsuccessfulOutput.Success)
	assert.Equal(t, dbInstances[0].ID, unsuccessfulOutput.DatabaseInstanceID)
	assert.Equal(t, dbInstances[0].Name, unsuccessfulOutput.Instance)
	assert.Equal(t, "databases synchronized but the is an error fetching the database instance 1a3af483-89e5-4820-a579-13ed6d90b0cc. Cause: sql: connection is already closed", unsuccessfulOutput.Message)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindAllDTOsEnabled", 1)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindByID", 1)
	databaseStorage.AssertNumberOfCalls(t, "FindAll", 1)
	databaseStorage.AssertNumberOfCalls(t, "Save", 1)
	databaseStorage.AssertNumberOfCalls(t, "Update", 3)
}

func TestGivenAnErrorInDbWhileUpdateInstance_WhenExecuteSyncDatabases_ThenShouldReturnSuccessFalse(t *testing.T) {
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	databaseStorage := new(mocks.DatabaseStorageMock)
	dbInstances := []*dto.DatabaseInstanceOutputDTO{mocks.BuildAzInstanceDTO()}
	dbInstanceStorage.On("FindAllDTOsEnabled", "", "").Return(dbInstances, nil).Once()
	dbInstanceStorage.On("FindByID", dbInstances[0].ID).Return(&entity.DatabaseInstance{}, nil).Once()
	dbInstanceStorage.On("Update", mock.Anything).Return(sql.ErrConnDone).Once()
	databases := mocks.BuildDatabaseList()
	databaseStorage.On("FindAll", dbInstances[0].ID, []string{}).Return(databases, nil).Once()
	databaseStorage.On("Save", mock.Anything).Return(nil).Once()
	databaseStorage.On("Update", mock.Anything).Return(nil).Times(3)

	uc := NewSyncDatabasesUseCase(dbInstanceStorage, databaseStorage)
	outputs, err := uc.Execute(dto.SyncDatabasesInputDTO{}, mocks.UserID)
	unsuccessfulOutput := outputs[0]

	assert.NoError(t, err, "no error expected")
	assert.NotNil(t, outputs)
	assert.Len(t, outputs, 1)
	assert.False(t, unsuccessfulOutput.Success)
	assert.Equal(t, dbInstances[0].ID, unsuccessfulOutput.DatabaseInstanceID)
	assert.Equal(t, dbInstances[0].Name, unsuccessfulOutput.Instance)
	assert.Equal(t, "databases synchronized but the is an error updating last sync date for database instance 1a3af483-89e5-4820-a579-13ed6d90b0cc. Cause: sql: connection is already closed", unsuccessfulOutput.Message)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindAllDTOsEnabled", 1)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindByID", 1)
	databaseStorage.AssertNumberOfCalls(t, "FindAll", 1)
	databaseStorage.AssertNumberOfCalls(t, "Save", 1)
	databaseStorage.AssertNumberOfCalls(t, "Update", 3)
}

func TestGivenInputWithIdsNotExistent_WhenExecuteSyncDatabases_ThenShouldReturnError(t *testing.T) {
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	databaseStorage := new(mocks.DatabaseStorageMock)
	dbInstanceStorage.On("FindAllDTOs", "", "", []string{"1", "2"}).Return([]*dto.DatabaseInstanceOutputDTO{}, sql.ErrNoRows).Once()

	uc := NewSyncDatabasesUseCase(dbInstanceStorage, databaseStorage)
	outputs, err := uc.Execute(dto.SyncDatabasesInputDTO{DatabaseInstancesIDs: []string{"1", "2"}}, mocks.UserID)

	assert.Error(t, err, "error expected when no database instances found")
	assert.EqualError(t, err, ErrNoDatabaseInstancesFound.Error())
	assert.Nil(t, outputs)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindAllDTOs", 1)
}

func TestGivenInstanceWithConnectorNotImplemented_WhenExecuteSyncDatabases_ThenShouldReturnSuccessFalse(t *testing.T) {
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	dbInstanceMysql := mocks.BuildConnectorNotImplementedInstance()
	dbInstances := []*dto.DatabaseInstanceOutputDTO{dbInstanceMysql}
	dbInstanceStorage.On("FindAllDTOsEnabled", "", "").Return(dbInstances, nil).Once()

	uc := NewSyncDatabasesUseCase(dbInstanceStorage, nil)
	outputs, err := uc.Execute(dto.SyncDatabasesInputDTO{}, mocks.UserID)
	notImplementedConnectorOutput := outputs[0]

	assert.NoError(t, err, "no error expected")
	assert.NotNil(t, outputs)
	assert.Len(t, outputs, 1)
	assert.False(t, notImplementedConnectorOutput.Success)
	assert.Equal(t, dbInstanceMysql.ID, notImplementedConnectorOutput.DatabaseInstanceID)
	assert.Equal(t, dbInstanceMysql.Name, notImplementedConnectorOutput.Instance)
	assert.Equal(t, "the database technology 'mysql' don't have a connector implemented", notImplementedConnectorOutput.Message)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindAllDTOsEnabled", 1)
}

func TestGivenSomeDbInstances_WhenExecuteSyncDatabasesByIds_ThenShouldSyncDatabases(t *testing.T) {
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	dbInstances := mocks.BuildInstancesList()
	var dbInstancesIds []string
	for _, dbInstance := range dbInstances {
		dbInstancesIds = append(dbInstancesIds, dbInstance.ID)
	}
	dbInstanceAzure := &entity.DatabaseInstance{
		ID:   uuid.MustParse(dbInstances[1].ID),
		Name: connector.DummyTest + " - Azure",
	}
	dbInstanceStorage.On("FindAllDTOs", "", "", dbInstancesIds).Return(dbInstances, nil).Once()
	dbInstanceStorage.On("FindByID", dbInstanceAzure.ID.String()).Return(dbInstanceAzure, nil).Once()
	dbInstanceStorage.On("Update", mock.Anything).Return(nil).Once()
	databaseStorage := new(mocks.DatabaseStorageMock)
	databases := mocks.BuildDatabaseList()
	databaseStorage.On("FindAll", dbInstanceAzure.ID.String(), []string{}).Return(databases, nil).Once()
	databaseStorage.On("Save", mock.Anything).Return(nil).Once()
	databaseStorage.On("Update", mock.Anything).Return(nil).Times(3)

	uc := NewSyncDatabasesUseCase(dbInstanceStorage, databaseStorage)
	outputs, err := uc.Execute(dto.SyncDatabasesInputDTO{DatabaseInstancesIDs: dbInstancesIds}, mocks.UserID)

	for _, output := range outputs {
		if output.Success {
			assert.Equal(t, dbInstanceAzure.ID.String(), output.DatabaseInstanceID)
			assert.Equal(t, dbInstanceAzure.Name, output.Instance)
			assert.Equal(t, 3, output.TotalDatabases)
		} else {
			assert.True(t, output.DatabaseInstanceID == dbInstances[0].ID || output.DatabaseInstanceID == dbInstances[2].ID)
		}
	}
	assert.NoError(t, err, "no error expected")
	assert.NotNil(t, outputs)
	assert.Len(t, outputs, 3)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindAllDTOs", 1)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindByID", 1)
	dbInstanceStorage.AssertNumberOfCalls(t, "Update", 1)
	databaseStorage.AssertNumberOfCalls(t, "FindAll", 1)
	databaseStorage.AssertNumberOfCalls(t, "Save", 1)
	databaseStorage.AssertNumberOfCalls(t, "Update", 3)
}

func TestGivenSomeDbInstances_WhenExecuteSyncDatabases_ThenShouldSyncDatabases(t *testing.T) {
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	dbInstances := []*dto.DatabaseInstanceOutputDTO{mocks.BuildAzInstanceDTO()}
	dbInstanceAzure := &entity.DatabaseInstance{
		ID:   uuid.MustParse(dbInstances[0].ID),
		Name: connector.DummyTest + " - Azure",
	}
	dbInstanceStorage.On("FindAllDTOsEnabled", "", "").Return(dbInstances, nil).Once()
	dbInstanceStorage.On("FindByID", dbInstanceAzure.ID.String()).Return(dbInstanceAzure, nil).Once()
	dbInstanceStorage.On("Update", mock.Anything).Return(nil).Once()
	databaseStorage := new(mocks.DatabaseStorageMock)
	databases := mocks.BuildDatabaseList()
	databaseStorage.On("FindAll", dbInstanceAzure.ID.String(), []string{}).Return(databases, nil).Once()
	databaseStorage.On("Save", mock.Anything).Return(nil).Once()
	databaseStorage.On("Update", mock.Anything).Return(nil).Times(3)

	uc := NewSyncDatabasesUseCase(dbInstanceStorage, databaseStorage)
	outputs, err := uc.Execute(dto.SyncDatabasesInputDTO{}, mocks.UserID)

	assert.NoError(t, err, "no error expected")
	assert.NotNil(t, outputs)
	assert.Len(t, outputs, 1)
	assert.Equal(t, dbInstanceAzure.ID.String(), outputs[0].DatabaseInstanceID)
	assert.Equal(t, dbInstanceAzure.Name, outputs[0].Instance)
	assert.Equal(t, 3, outputs[0].TotalDatabases)
	assert.True(t, outputs[0].Success)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindAllDTOsEnabled", 1)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindByID", 1)
	dbInstanceStorage.AssertNumberOfCalls(t, "Update", 1)
	databaseStorage.AssertNumberOfCalls(t, "FindAll", 1)
	databaseStorage.AssertNumberOfCalls(t, "Save", 1)
	databaseStorage.AssertNumberOfCalls(t, "Update", 3)
}

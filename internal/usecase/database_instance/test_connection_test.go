package instance

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

func TestGivenAnErrorInDb_WhenExecuteTestConnection_ThenShouldReturnError(t *testing.T) {
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	dbInstanceStorage.On("FindAllDTOsEnabled", "", "").Return([]*dto.DatabaseInstanceOutputDTO{}, sql.ErrConnDone).Once()

	uc := NewTestConnectionUseCase(dbInstanceStorage)
	outputs, err := uc.Execute(dto.TestConnectionInputDTO{})

	assert.Error(t, err, "error expected when some error in db")
	assert.EqualError(t, err, sql.ErrConnDone.Error())
	assert.Nil(t, outputs)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindAllDTOsEnabled", 1)
}

func TestGivenAnErrorInDbWhileFetchingInstances_WhenExecuteTestConnection_ThenShouldReturnError(t *testing.T) {
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	dbInstanceStorage.On("FindAllDTOs", "", "", []string{"1", "2"}).Return([]*dto.DatabaseInstanceOutputDTO{}, sql.ErrConnDone).Once()

	uc := NewTestConnectionUseCase(dbInstanceStorage)
	outputs, err := uc.Execute(dto.TestConnectionInputDTO{DatabaseInstancesIDs: []string{"1", "2"}})

	assert.Error(t, err, "error expected when some error in db")
	assert.EqualError(t, err, sql.ErrConnDone.Error())
	assert.Nil(t, outputs)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindAllDTOs", 1)
}

func TestGivenInputWithIdsNotExistent_WhenExecuteTestConnection_ThenShouldReturnError(t *testing.T) {
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	dbInstanceStorage.On("FindAllDTOs", "", "", []string{"1", "2"}).Return([]*dto.DatabaseInstanceOutputDTO{}, sql.ErrNoRows).Once()

	uc := NewTestConnectionUseCase(dbInstanceStorage)
	outputs, err := uc.Execute(dto.TestConnectionInputDTO{DatabaseInstancesIDs: []string{"1", "2"}})

	assert.Error(t, err, "error expected when no database instances found")
	assert.EqualError(t, err, ErrNoDatabaseInstancesFound.Error())
	assert.Nil(t, outputs)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindAllDTOs", 1)
}

func TestGivenAnErrorInDbWhileFetchingInstanceToUpdate_WhenExecuteTestConnection_ThenShouldReturnSuccessFalse(t *testing.T) {
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	dbInstances := []*dto.DatabaseInstanceOutputDTO{mocks.BuildAzureInstanceDTO()}
	dbInstanceStorage.On("FindAllDTOsEnabled", "", "").Return(dbInstances, nil).Once()
	dbInstanceStorage.On("FindByID", dbInstances[0].ID).Return(&entity.DatabaseInstance{}, sql.ErrConnDone).Once()

	uc := NewTestConnectionUseCase(dbInstanceStorage)
	outputs, err := uc.Execute(dto.TestConnectionInputDTO{})
	unsuccessfulOutput := outputs[0]

	assert.NoError(t, err, "no error expected")
	assert.NotNil(t, outputs)
	assert.Len(t, outputs, 1)
	assert.False(t, unsuccessfulOutput.Success)
	assert.Equal(t, dbInstances[0].ID, unsuccessfulOutput.DatabaseInstanceID)
	assert.Equal(t, dbInstances[0].Name, unsuccessfulOutput.Instance)
	assert.Equal(t, "connection established successfully! However there was an error while fetching the database instance 1a3af483-89e5-4820-a579-13ed6d90b0cc. Cause: sql: connection is already closed", unsuccessfulOutput.Message)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindAllDTOsEnabled", 1)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindByID", 1)
}

func TestGivenAnErrorInDbWhileUpdateInstance_WhenExecuteTestConnection_ThenShouldReturnSuccessFalse(t *testing.T) {
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	dbInstances := []*dto.DatabaseInstanceOutputDTO{mocks.BuildAzureInstanceDTO()}
	dbInstanceStorage.On("FindAllDTOsEnabled", "", "").Return(dbInstances, nil).Once()
	dbInstanceStorage.On("FindByID", dbInstances[0].ID).Return(&entity.DatabaseInstance{}, nil).Once()
	dbInstanceStorage.On("Update", mock.Anything).Return(sql.ErrConnDone).Once()

	uc := NewTestConnectionUseCase(dbInstanceStorage)
	outputs, err := uc.Execute(dto.TestConnectionInputDTO{})
	unsuccessfulOutput := outputs[0]

	assert.NoError(t, err, "no error expected")
	assert.NotNil(t, outputs)
	assert.Len(t, outputs, 1)
	assert.False(t, unsuccessfulOutput.Success)
	assert.Equal(t, dbInstances[0].ID, unsuccessfulOutput.DatabaseInstanceID)
	assert.Equal(t, dbInstances[0].Name, unsuccessfulOutput.Instance)
	assert.Equal(t, "connection established successfully! However there was an error to update last connection date for database instance 1a3af483-89e5-4820-a579-13ed6d90b0cc. Cause: sql: connection is already closed", unsuccessfulOutput.Message)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindAllDTOsEnabled", 1)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindByID", 1)
}

func TestGivenInstanceWithConnectorNotImplemented_WhenExecuteTestConnection_ThenShouldReturnSuccessFalse(t *testing.T) {
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	dbInstanceMysql := &dto.DatabaseInstanceOutputDTO{
		ID:                        "2",
		Name:                      "MySQL - AWS",
		DatabaseTechnologyName:    "MySQL",
		DatabaseTechnologyVersion: "1",
		Enabled:                   true,
	}
	dbInstances := []*dto.DatabaseInstanceOutputDTO{dbInstanceMysql}
	dbInstanceStorage.On("FindAllDTOsEnabled", "", "").Return(dbInstances, nil).Once()

	uc := NewTestConnectionUseCase(dbInstanceStorage)
	outputs, err := uc.Execute(dto.TestConnectionInputDTO{})

	assert.NoError(t, err, "no error expected")
	assert.NotNil(t, outputs)
	assert.Len(t, outputs, 1)
	assert.False(t, outputs[0].Success)
	assert.Equal(t, dbInstanceMysql.ID, outputs[0].DatabaseInstanceID)
	assert.Equal(t, dbInstanceMysql.Name, outputs[0].Instance)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindAllDTOsEnabled", 1)
}

func TestGivenSomeDbInstances_WhenExecuteTestConnectionByIds_ThenShouldTestConnection(t *testing.T) {
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	dbInstances := mocks.BuildInstancesList()
	var dbInstancesIds []string
	for _, dbInstance := range dbInstances {
		dbInstancesIds = append(dbInstancesIds, dbInstance.ID)
	}
	dbInstanceAws := &entity.DatabaseInstance{
		ID:   uuid.MustParse(dbInstances[1].ID),
		Name: connector.DummyTest + " - Azure",
	}
	dbInstanceDummyError := &entity.DatabaseInstance{
		ID:   uuid.MustParse(dbInstances[2].ID),
		Name: connector.InstanceDummyTestError,
	}
	dbInstanceStorage.On("FindAllDTOs", "", "", dbInstancesIds).Return(dbInstances, nil).Once()
	dbInstanceStorage.On("FindByID", dbInstanceAws.ID.String()).Return(dbInstanceAws, nil).Once()
	dbInstanceStorage.On("FindByID", dbInstanceDummyError.ID.String()).Return(dbInstanceDummyError, nil).Once()
	dbInstanceStorage.On("Update", mock.Anything).Return(nil).Twice()

	uc := NewTestConnectionUseCase(dbInstanceStorage)
	outputs, err := uc.Execute(dto.TestConnectionInputDTO{DatabaseInstancesIDs: dbInstancesIds})

	for _, output := range outputs {
		if output.Success {
			assert.Equal(t, dbInstanceAws.ID.String(), output.DatabaseInstanceID)
			assert.Equal(t, dbInstanceAws.Name, output.Instance)
		} else {
			assert.True(t, output.DatabaseInstanceID == dbInstances[0].ID || output.DatabaseInstanceID == dbInstances[2].ID)
		}
	}
	assert.NoError(t, err, "no error expected")
	assert.NotNil(t, outputs)
	assert.Len(t, outputs, 3)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindAllDTOs", 1)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindByID", 2)
	dbInstanceStorage.AssertNumberOfCalls(t, "Update", 2)
}

func TestGivenErrorWhileTestingConnectionAndUpdatingInstance_WhenExecuteTestConnectionByIds_ThenShouldReturnSuccessFalse(t *testing.T) {
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	dbInstanceDtoDummy := mocks.BuildDummyErrorInstance()
	dbInstanceDummyError := &entity.DatabaseInstance{
		ID:   uuid.MustParse(dbInstanceDtoDummy.ID),
		Name: connector.InstanceDummyTestError,
	}
	dbInstances := []*dto.DatabaseInstanceOutputDTO{dbInstanceDtoDummy}
	dbInstanceStorage.On("FindAllDTOs", "", "", []string{dbInstanceDtoDummy.ID}).Return(dbInstances, nil).Once()
	dbInstanceStorage.On("FindByID", dbInstanceDtoDummy.ID).Return(dbInstanceDummyError, nil).Once()
	dbInstanceStorage.On("Update", mock.Anything).Return(sql.ErrConnDone).Twice()

	uc := NewTestConnectionUseCase(dbInstanceStorage)
	outputs, err := uc.Execute(dto.TestConnectionInputDTO{DatabaseInstancesIDs: []string{dbInstanceDtoDummy.ID}})
	unsuccessfulOutput := outputs[0]

	assert.NoError(t, err, "no error expected")
	assert.NotNil(t, outputs)
	assert.Len(t, outputs, 1)
	assert.False(t, unsuccessfulOutput.Success)
	assert.Equal(t, dbInstanceDtoDummy.ID, unsuccessfulOutput.DatabaseInstanceID)
	assert.Equal(t, dbInstanceDtoDummy.Name, unsuccessfulOutput.Instance)
	assert.Equal(t, "error testing connection with instance-dummy-test-error. There was also an error to update last connection date for database instance ba3aaedd-458f-4582-b763-12c3ae7b27ee. Cause: sql: connection is already closed", unsuccessfulOutput.Message)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindAllDTOs", 1)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindByID", 1)
	dbInstanceStorage.AssertNumberOfCalls(t, "Update", 1)
}

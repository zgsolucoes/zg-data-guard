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

func TestGivenAnErrorInDb_WhenExecutePropagateRoles_ThenShouldReturnError(t *testing.T) {
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	dbInstanceStorage.On("FindAllDTOsEnabled", "", "").Return([]*dto.DatabaseInstanceOutputDTO{}, sql.ErrConnDone).Once()
	roleStorage := new(mocks.DatabaseRoleStorageMock)

	uc := NewPropagateRolesUseCase(dbInstanceStorage, roleStorage)
	outputs, err := uc.Execute(dto.PropagateRolesInputDTO{}, mocks.UserID)

	assert.Error(t, err, "error expected when some error in db")
	assert.EqualError(t, err, sql.ErrConnDone.Error())
	assert.Nil(t, outputs)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindAllDTOsEnabled", 1)
}

func TestGivenAnErrorInDbWhileFetchingInstances_WhenExecutePropagateRoles_ThenShouldReturnError(t *testing.T) {
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	dbInstanceStorage.On("FindAllDTOs", "", "", []string{"1", "2"}).Return([]*dto.DatabaseInstanceOutputDTO{}, sql.ErrConnDone).Once()
	roleStorage := new(mocks.DatabaseRoleStorageMock)

	uc := NewPropagateRolesUseCase(dbInstanceStorage, roleStorage)
	outputs, err := uc.Execute(dto.PropagateRolesInputDTO{DatabaseInstancesIDs: []string{"1", "2"}}, mocks.UserID)

	assert.Error(t, err, "error expected when some error in db")
	assert.EqualError(t, err, sql.ErrConnDone.Error())
	assert.Nil(t, outputs)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindAllDTOs", 1)
}

func TestGivenInputWithIdsNotExistent_WhenExecutePropagateRoles_ThenShouldReturnError(t *testing.T) {
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	dbInstanceStorage.On("FindAllDTOs", "", "", []string{"1", "2"}).Return([]*dto.DatabaseInstanceOutputDTO{}, sql.ErrNoRows).Once()
	roleStorage := new(mocks.DatabaseRoleStorageMock)

	uc := NewPropagateRolesUseCase(dbInstanceStorage, roleStorage)
	outputs, err := uc.Execute(dto.PropagateRolesInputDTO{DatabaseInstancesIDs: []string{"1", "2"}}, mocks.UserID)

	assert.Error(t, err, "error expected when no database instances found")
	assert.EqualError(t, err, ErrNoDatabaseInstancesFound.Error())
	assert.Nil(t, outputs)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindAllDTOs", 1)
}

func TestGivenAnErrorInDbWhileFetchingRoles_WhenExecutePropagateRoles_ThenShouldReturnError(t *testing.T) {
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	dbInstances := []*dto.DatabaseInstanceOutputDTO{mocks.BuildAzInstanceDTO()}
	dbInstanceStorage.On("FindAllDTOsEnabled", "", "").Return(dbInstances, nil).Once()
	roleStorage := new(mocks.DatabaseRoleStorageMock)
	roleStorage.On("FindAll").Return([]*entity.DatabaseRole{}, sql.ErrConnDone).Once()

	uc := NewPropagateRolesUseCase(dbInstanceStorage, roleStorage)
	outputs, err := uc.Execute(dto.PropagateRolesInputDTO{}, mocks.UserID)

	assert.Error(t, err, "error expected when some error in db")
	assert.EqualError(t, err, "error while fetching database roles. Cause: "+sql.ErrConnDone.Error())
	assert.Nil(t, outputs)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindAllDTOsEnabled", 1)
	roleStorage.AssertNumberOfCalls(t, "FindAll", 1)
}

func TestGivenAnErrorInDbWhileFetchingInstanceToUpdate_WhenExecutePropagateRoles_ThenShouldReturnSuccessFalse(t *testing.T) {
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	dbInstances := []*dto.DatabaseInstanceOutputDTO{mocks.BuildAzInstanceDTO()}
	dbInstanceStorage.On("FindAllDTOsEnabled", "", "").Return(dbInstances, nil).Once()
	dbInstanceStorage.On("FindByID", dbInstances[0].ID).Return(&entity.DatabaseInstance{}, sql.ErrConnDone).Once()
	roleStorage := new(mocks.DatabaseRoleStorageMock)
	roleStorage.On("FindAll").Return(mocks.BuildRolesList(), nil).Once()

	uc := NewPropagateRolesUseCase(dbInstanceStorage, roleStorage)
	outputs, err := uc.Execute(dto.PropagateRolesInputDTO{}, mocks.UserID)
	unsuccessfulOutput := outputs[0]

	assert.NoError(t, err, "no error expected")
	assert.NotNil(t, outputs)
	assert.Len(t, outputs, 1)
	assert.False(t, unsuccessfulOutput.Success)
	assert.Equal(t, dbInstances[0].ID, unsuccessfulOutput.DatabaseInstanceID)
	assert.Equal(t, dbInstances[0].Name, unsuccessfulOutput.Instance)
	assert.Equal(t, "roles created in database instance successfully! However there was an error while fetching the database instance 1a3af483-89e5-4820-a579-13ed6d90b0cc. Cause: sql: connection is already closed", unsuccessfulOutput.Message)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindAllDTOsEnabled", 1)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindByID", 1)
	roleStorage.AssertNumberOfCalls(t, "FindAll", 1)
}

func TestGivenAnErrorInDbWhileUpdateInstance_WhenExecutePropagateRoles_ThenShouldReturnSuccessFalse(t *testing.T) {
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	dbInstances := []*dto.DatabaseInstanceOutputDTO{mocks.BuildAzInstanceDTO()}
	dbInstanceStorage.On("FindAllDTOsEnabled", "", "").Return(dbInstances, nil).Once()
	dbInstanceStorage.On("FindByID", dbInstances[0].ID).Return(&entity.DatabaseInstance{}, nil).Once()
	dbInstanceStorage.On("Update", mock.Anything).Return(sql.ErrConnDone).Once()
	roleStorage := new(mocks.DatabaseRoleStorageMock)
	roleStorage.On("FindAll").Return(mocks.BuildRolesList(), nil).Once()

	uc := NewPropagateRolesUseCase(dbInstanceStorage, roleStorage)
	outputs, err := uc.Execute(dto.PropagateRolesInputDTO{}, mocks.UserID)
	unsuccessfulOutput := outputs[0]

	assert.NoError(t, err, "no error expected")
	assert.NotNil(t, outputs)
	assert.Len(t, outputs, 1)
	assert.False(t, unsuccessfulOutput.Success)
	assert.Equal(t, dbInstances[0].ID, unsuccessfulOutput.DatabaseInstanceID)
	assert.Equal(t, dbInstances[0].Name, unsuccessfulOutput.Instance)
	assert.Equal(t, "roles created in database instance successfully! However there was an error to update created roles in database instance 1a3af483-89e5-4820-a579-13ed6d90b0cc. Cause: sql: connection is already closed", unsuccessfulOutput.Message)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindAllDTOsEnabled", 1)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindByID", 1)
	roleStorage.AssertNumberOfCalls(t, "FindAll", 1)
}

func TestGivenInstanceWithConnectorNotImplemented_WhenExecutePropagateRoles_ThenShouldReturnSuccessFalse(t *testing.T) {
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	dbInstanceMysql := &dto.DatabaseInstanceOutputDTO{
		ID:                        "2",
		Name:                      "MySQL - Local",
		DatabaseTechnologyName:    "MySQL",
		DatabaseTechnologyVersion: "1",
		Enabled:                   true,
	}
	dbInstances := []*dto.DatabaseInstanceOutputDTO{dbInstanceMysql}
	dbInstanceStorage.On("FindAllDTOsEnabled", "", "").Return(dbInstances, nil).Once()
	roleStorage := new(mocks.DatabaseRoleStorageMock)
	roleStorage.On("FindAll").Return(mocks.BuildRolesList(), nil).Once()

	uc := NewPropagateRolesUseCase(dbInstanceStorage, roleStorage)
	outputs, err := uc.Execute(dto.PropagateRolesInputDTO{}, mocks.UserID)

	assert.NoError(t, err, "no error expected")
	assert.NotNil(t, outputs)
	assert.Len(t, outputs, 1)
	assert.False(t, outputs[0].Success)
	assert.Equal(t, dbInstanceMysql.ID, outputs[0].DatabaseInstanceID)
	assert.Equal(t, dbInstanceMysql.Name, outputs[0].Instance)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindAllDTOsEnabled", 1)
	roleStorage.AssertNumberOfCalls(t, "FindAll", 1)
}

func TestGivenSomeInstances_WhenExecutePropagateRolesByIds_ThenShouldPropagateRolesOnlyInValidInstances(t *testing.T) {
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
	dbInstanceStorage.On("FindAllDTOs", "", "", dbInstancesIds).Return(dbInstances, nil).Once()
	dbInstanceStorage.On("FindByID", dbInstanceAws.ID.String()).Return(dbInstanceAws, nil).Once()
	dbInstanceStorage.On("Update", mock.Anything).Return(nil).Once()
	roleStorage := new(mocks.DatabaseRoleStorageMock)
	roleStorage.On("FindAll").Return(mocks.BuildRolesList(), nil).Once()

	uc := NewPropagateRolesUseCase(dbInstanceStorage, roleStorage)
	outputs, err := uc.Execute(dto.PropagateRolesInputDTO{DatabaseInstancesIDs: dbInstancesIds}, mocks.UserID)

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
	dbInstanceStorage.AssertNumberOfCalls(t, "FindByID", 1)
	dbInstanceStorage.AssertNumberOfCalls(t, "Update", 1)
	roleStorage.AssertNumberOfCalls(t, "FindAll", 1)
}

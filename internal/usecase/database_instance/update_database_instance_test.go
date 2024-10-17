package instance

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	"github.com/zgsolucoes/zg-data-guard/internal/entity"
	"github.com/zgsolucoes/zg-data-guard/internal/usecase/common"
	"github.com/zgsolucoes/zg-data-guard/testdata/mocks"
)

var (
	updatedInput = dto.DatabaseInstanceInputDTO{
		Name:                 "Elasticsearch",
		Host:                 "10.1.1.1",
		Port:                 "9200",
		HostConnection:       "192.168.1.1",
		PortConnection:       "9201",
		AdminUser:            "elkadm",
		AdminPassword:        "elkpwd",
		EcosystemID:          "7e4f6436-2cf2-4d54-ab09-3bf77be8dec8",
		DatabaseTechnologyID: "5d46f461-1f71-4661-b869-fe1c5517c75b",
		Note:                 "Cool text",
	}
)

func TestGivenANonexistentId_WhenExecuteUpdate_ThenShouldReturnError(t *testing.T) {
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	ecosystemStorage := new(mocks.EcosystemStorageMock)
	techStorage := new(mocks.TechnologyStorageMock)
	dbInstanceStorage.On("FindByID", mocks.DatabaseInstanceId).Return(&entity.DatabaseInstance{}, sql.ErrNoRows).Once()
	uc := NewUpdateDatabaseInstanceUseCase(dbInstanceStorage, ecosystemStorage, techStorage)

	output, err := uc.Execute(validInput, mocks.DatabaseInstanceId, mocks.UserID)

	assert.Error(t, err, "error expected when database instance not found")
	assert.EqualError(t, err, ErrDatabaseInstanceNotFound.Error())
	assert.Nil(t, output)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindByID", 1)
}

func TestGivenAnErrorInDbWhileFetchId_WhenExecuteUpdate_ThenShouldReturnError(t *testing.T) {
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	ecosystemStorage := new(mocks.EcosystemStorageMock)
	techStorage := new(mocks.TechnologyStorageMock)
	dbInstanceStorage.On("FindByID", mocks.DatabaseInstanceId).Return(&entity.DatabaseInstance{}, sql.ErrConnDone).Once()
	uc := NewUpdateDatabaseInstanceUseCase(dbInstanceStorage, ecosystemStorage, techStorage)

	output, err := uc.Execute(validInput, mocks.DatabaseInstanceId, mocks.UserID)

	assert.Error(t, err, "error expected when database instance not found")
	assert.EqualError(t, err, sql.ErrConnDone.Error())
	assert.Nil(t, output)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindByID", 1)
}

func TestGivenAHostAlreadyExistent_WhenExecuteUpdate_ThenShouldReturnError(t *testing.T) {
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	ecosystemStorage := new(mocks.EcosystemStorageMock)
	techStorage := new(mocks.TechnologyStorageMock)
	dbInstance, _ := entity.NewDatabaseInstance(validInput, mocks.UserID)
	dbInstanceStorage.On("FindByID", dbInstance.ID.String()).Return(dbInstance, nil).Once()
	dbInstanceStorage.On("Exists", updatedInput.Host, updatedInput.Port).Return(true, nil).Once()
	uc := NewUpdateDatabaseInstanceUseCase(dbInstanceStorage, ecosystemStorage, techStorage)

	output, err := uc.Execute(updatedInput, dbInstance.ID.String(), mocks.UserID)
	assert.Error(t, err, "error expected with host and port already existent")
	assert.EqualError(t, err, ErrHostAlreadyExists.Error())
	assert.Nil(t, output)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindByID", 1)
	dbInstanceStorage.AssertNumberOfCalls(t, "Exists", 1)
}

func TestGivenAnNonexistentEcosystemId_WhenExecuteUpdate_ThenShouldReturnError(t *testing.T) {
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	ecosystemStorage := new(mocks.EcosystemStorageMock)
	techStorage := new(mocks.TechnologyStorageMock)
	dbInstance, _ := entity.NewDatabaseInstance(validInput, mocks.UserID)
	dbInstanceStorage.On("FindByID", dbInstance.ID.String()).Return(dbInstance, nil).Once()
	ecosystemStorage.On("FindByID", validInput.EcosystemID).Return(&entity.Ecosystem{}, sql.ErrNoRows).Once()
	uc := NewUpdateDatabaseInstanceUseCase(dbInstanceStorage, ecosystemStorage, techStorage)

	output, err := uc.Execute(validInput, dbInstance.ID.String(), mocks.UserID)

	assert.Error(t, err, "error expected due to nonexistent ecosystem")
	assert.EqualError(t, err, common.ErrEcosystemNotFound.Error())
	assert.Nil(t, output)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindByID", 1)
	ecosystemStorage.AssertNumberOfCalls(t, "FindByID", 1)
}

func TestGivenAErrorInDbWhileFetchingEcosystemId_WhenExecuteUpdate_ThenShouldReturnError(t *testing.T) {
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	ecosystemStorage := new(mocks.EcosystemStorageMock)
	techStorage := new(mocks.TechnologyStorageMock)
	dbInstance, _ := entity.NewDatabaseInstance(validInput, mocks.UserID)
	dbInstanceStorage.On("FindByID", dbInstance.ID.String()).Return(dbInstance, nil).Once()
	ecosystemStorage.On("FindByID", validInput.EcosystemID).Return(&entity.Ecosystem{}, sql.ErrConnDone).Once()
	uc := NewUpdateDatabaseInstanceUseCase(dbInstanceStorage, ecosystemStorage, techStorage)

	dbInstanceOutput, err := uc.Execute(validInput, dbInstance.ID.String(), mocks.UserID)

	assert.Error(t, err, "error expected due to error in db")
	assert.EqualError(t, err, sql.ErrConnDone.Error())
	assert.Nil(t, dbInstanceOutput)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindByID", 1)
	ecosystemStorage.AssertNumberOfCalls(t, "FindByID", 1)
}

func TestGivenAnInvalidInput_WhenExecuteUpdate_ThenShouldReturnError(t *testing.T) {
	invalidInput := dto.DatabaseInstanceInputDTO{
		Name:                 "PostgreSQL - Cloud XPTO",
		Host:                 "192.1.1.1",
		Port:                 "5432",
		EcosystemID:          "7e4f6436-2cf2-4d54-ab09-3bf77be8dec8",
		DatabaseTechnologyID: "5d46f461-1f71-4661-b869-fe1c5517c75b",
	}
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	ecosystemStorage := new(mocks.EcosystemStorageMock)
	techStorage := new(mocks.TechnologyStorageMock)
	dbInstance, _ := entity.NewDatabaseInstance(validInput, mocks.UserID)

	dbInstanceStorage.On("FindByID", dbInstance.ID.String()).Return(dbInstance, nil).Once()
	dbInstanceStorage.On("Exists", invalidInput.Host, invalidInput.Port).Return(false, nil).Once()
	ecosystemStorage.On("FindByID", invalidInput.EcosystemID).Return(&entity.Ecosystem{}, nil).Once()
	techStorage.On("FindByID", invalidInput.DatabaseTechnologyID).Return(&entity.DatabaseTechnology{}, nil).Once()

	uc := NewUpdateDatabaseInstanceUseCase(dbInstanceStorage, ecosystemStorage, techStorage)
	dbInstanceOutput, err := uc.Execute(invalidInput, dbInstance.ID.String(), mocks.UserID)

	assert.Error(t, err, "error expected due to invalid host (not informed)")
	assert.EqualError(t, err, entity.ErrInvalidHostConnection.Error())
	assert.Nil(t, dbInstanceOutput)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindByID", 1)
	dbInstanceStorage.AssertNumberOfCalls(t, "Exists", 1)
	ecosystemStorage.AssertNumberOfCalls(t, "FindByID", 1)
	techStorage.AssertNumberOfCalls(t, "FindByID", 1)
}

func TestGivenAValidInput_WhenExecuteUpdate_ThenShouldUpdateDatabaseInstance(t *testing.T) {
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	ecosystemStorage := new(mocks.EcosystemStorageMock)
	techStorage := new(mocks.TechnologyStorageMock)
	dbInstance, _ := entity.NewDatabaseInstance(validInput, mocks.UserID)

	dbInstanceStorage.On("FindByID", dbInstance.ID.String()).Return(dbInstance, nil).Once()
	dbInstanceStorage.On("Exists", updatedInput.Host, updatedInput.Port).Return(false, nil).Once()
	ecosystemStorage.On("FindByID", updatedInput.EcosystemID).Return(&entity.Ecosystem{}, nil).Once()
	techStorage.On("FindByID", updatedInput.DatabaseTechnologyID).Return(&entity.DatabaseTechnology{}, nil).Once()
	dbInstanceStorage.On("UpdateWithHostInfo", mock.Anything).Return(nil).Once()

	uc := NewUpdateDatabaseInstanceUseCase(dbInstanceStorage, ecosystemStorage, techStorage)
	dbInstanceOutput, err := uc.Execute(updatedInput, dbInstance.ID.String(), mocks.UserID)

	assert.NoError(t, err, "no error expected with a valid input and not existent host and port")
	assert.NotNil(t, dbInstanceOutput, "database instance should not be nil")
	assert.Equal(t, updatedInput.Name, dbInstanceOutput.Name)
	assert.Equal(t, updatedInput.Host, dbInstanceOutput.Host)
	assert.Equal(t, updatedInput.Port, dbInstanceOutput.Port)
	assert.Equal(t, updatedInput.HostConnection, dbInstanceOutput.HostConnection)
	assert.Equal(t, updatedInput.PortConnection, dbInstanceOutput.PortConnection)
	assert.Equal(t, updatedInput.AdminUser, dbInstanceOutput.AdminUser)
	assert.Equal(t, updatedInput.EcosystemID, dbInstanceOutput.EcosystemID)
	assert.Equal(t, updatedInput.DatabaseTechnologyID, dbInstanceOutput.DatabaseTechnologyID)
	assert.Equal(t, updatedInput.Note, dbInstanceOutput.Note)
	assert.True(t, dbInstanceOutput.Enabled)
	assert.NotEmpty(t, dbInstanceOutput.ID)
	assert.NotEmpty(t, dbInstanceOutput.CreatedByUserID)
	assert.NotEmpty(t, dbInstanceOutput.CreatedAt)
	assert.NotEmpty(t, dbInstanceOutput.UpdatedAt)
	assert.Empty(t, dbInstanceOutput.DisabledAt)
	assert.NotEqual(t, dbInstanceOutput.CreatedAt, dbInstanceOutput.UpdatedAt)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindByID", 1)
	dbInstanceStorage.AssertNumberOfCalls(t, "Exists", 1)
	dbInstanceStorage.AssertNumberOfCalls(t, "UpdateWithHostInfo", 1)
	ecosystemStorage.AssertNumberOfCalls(t, "FindByID", 1)
	techStorage.AssertNumberOfCalls(t, "FindByID", 1)
}

func TestGivenAnErrorInDbWhileUpdate_WhenExecuteUpdate_ThenShouldReturnError(t *testing.T) {
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	ecosystemStorage := new(mocks.EcosystemStorageMock)
	techStorage := new(mocks.TechnologyStorageMock)
	dbInstance, _ := entity.NewDatabaseInstance(validInput, mocks.UserID)

	dbInstanceStorage.On("FindByID", dbInstance.ID.String()).Return(dbInstance, nil).Once()
	dbInstanceStorage.On("Exists", updatedInput.Host, updatedInput.Port).Return(false, nil).Once()
	ecosystemStorage.On("FindByID", updatedInput.EcosystemID).Return(&entity.Ecosystem{}, nil).Once()
	techStorage.On("FindByID", updatedInput.DatabaseTechnologyID).Return(&entity.DatabaseTechnology{}, nil).Once()
	dbInstanceStorage.On("UpdateWithHostInfo", mock.Anything).Return(sql.ErrConnDone).Once()

	uc := NewUpdateDatabaseInstanceUseCase(dbInstanceStorage, ecosystemStorage, techStorage)
	dbInstanceOutput, err := uc.Execute(updatedInput, dbInstance.ID.String(), mocks.UserID)

	assert.Error(t, err, "error expected when some error in db")
	assert.EqualError(t, err, sql.ErrConnDone.Error())
	assert.Nil(t, dbInstanceOutput)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindByID", 1)
	dbInstanceStorage.AssertNumberOfCalls(t, "Exists", 1)
	dbInstanceStorage.AssertNumberOfCalls(t, "UpdateWithHostInfo", 1)
	ecosystemStorage.AssertNumberOfCalls(t, "FindByID", 1)
	techStorage.AssertNumberOfCalls(t, "FindByID", 1)
}

func TestGivenAnErrorInDbWhileCheckHost_WhenExecuteUpdate_ThenShouldReturnError(t *testing.T) {
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	ecosystemStorage := new(mocks.EcosystemStorageMock)
	techStorage := new(mocks.TechnologyStorageMock)
	dbInstance, _ := entity.NewDatabaseInstance(validInput, mocks.UserID)

	dbInstanceStorage.On("FindByID", dbInstance.ID.String()).Return(dbInstance, nil).Once()
	dbInstanceStorage.On("Exists", updatedInput.Host, updatedInput.Port).Return(false, sql.ErrConnDone).Once()

	uc := NewUpdateDatabaseInstanceUseCase(dbInstanceStorage, ecosystemStorage, techStorage)
	dbInstanceOutput, err := uc.Execute(updatedInput, dbInstance.ID.String(), mocks.UserID)

	assert.Error(t, err, "error expected when some error in db")
	assert.EqualError(t, err, sql.ErrConnDone.Error())
	assert.Nil(t, dbInstanceOutput)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindByID", 1)
	dbInstanceStorage.AssertNumberOfCalls(t, "Exists", 1)
}

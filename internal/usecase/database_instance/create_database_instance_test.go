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

var validInput = mocks.ValidInstanceInput

func TestGivenAHostAlreadyExistent_WhenExecuteCreate_ThenShouldReturnError(t *testing.T) {
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	ecosystemStorage := new(mocks.EcosystemStorageMock)
	techStorage := new(mocks.TechnologyStorageMock)
	dbInstanceStorage.On("Exists", validInput.Host, validInput.Port).Return(true, nil).Once()

	uc := NewCreateDatabaseInstanceUseCase(dbInstanceStorage, ecosystemStorage, techStorage)
	output, err := uc.Execute(validInput, mocks.UserID)

	assert.Error(t, err, "error expected with host and port already existent")
	assert.EqualError(t, err, ErrHostAlreadyExists.Error())
	assert.Nil(t, output)
	dbInstanceStorage.AssertNumberOfCalls(t, "Exists", 1)
}

func TestGivenAnNonexistentEcosystemId_WhenExecuteCreate_ThenShouldReturnError(t *testing.T) {
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	ecosystemStorage := new(mocks.EcosystemStorageMock)
	techStorage := new(mocks.TechnologyStorageMock)
	dbInstanceStorage.On("Exists", validInput.Host, validInput.Port).Return(false, nil).Once()
	ecosystemStorage.On("FindByID", validInput.EcosystemID).Return(&entity.Ecosystem{}, sql.ErrNoRows).Once()

	uc := NewCreateDatabaseInstanceUseCase(dbInstanceStorage, ecosystemStorage, techStorage)
	dbInstanceOutput, err := uc.Execute(validInput, mocks.UserID)

	assert.Error(t, err, "error expected due to nonexistent ecosystem")
	assert.EqualError(t, err, common.ErrEcosystemNotFound.Error())
	assert.Nil(t, dbInstanceOutput)
	dbInstanceStorage.AssertNumberOfCalls(t, "Exists", 1)
	ecosystemStorage.AssertNumberOfCalls(t, "FindByID", 1)
}

func TestGivenAErrorInDbWhileFetchingEcosystemId_WhenExecuteCreate_ThenShouldReturnError(t *testing.T) {
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	ecosystemStorage := new(mocks.EcosystemStorageMock)
	techStorage := new(mocks.TechnologyStorageMock)
	dbInstanceStorage.On("Exists", validInput.Host, validInput.Port).Return(false, nil).Once()
	ecosystemStorage.On("FindByID", validInput.EcosystemID).Return(&entity.Ecosystem{}, sql.ErrConnDone).Once()

	uc := NewCreateDatabaseInstanceUseCase(dbInstanceStorage, ecosystemStorage, techStorage)
	dbInstanceOutput, err := uc.Execute(validInput, mocks.UserID)

	assert.Error(t, err, "error expected due to error in db")
	assert.EqualError(t, err, sql.ErrConnDone.Error())
	assert.Nil(t, dbInstanceOutput)
	dbInstanceStorage.AssertNumberOfCalls(t, "Exists", 1)
	ecosystemStorage.AssertNumberOfCalls(t, "FindByID", 1)
}

func TestGivenAnNonexistentTechnologyId_WhenExecuteCreate_ThenShouldReturnError(t *testing.T) {
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	ecosystemStorage := new(mocks.EcosystemStorageMock)
	techStorage := new(mocks.TechnologyStorageMock)
	dbInstanceStorage.On("Exists", validInput.Host, validInput.Port).Return(false, nil).Once()
	ecosystemStorage.On("FindByID", validInput.EcosystemID).Return(&entity.Ecosystem{}, nil).Once()
	techStorage.On("FindByID", validInput.DatabaseTechnologyID).Return(&entity.DatabaseTechnology{}, sql.ErrNoRows).Once()

	uc := NewCreateDatabaseInstanceUseCase(dbInstanceStorage, ecosystemStorage, techStorage)
	dbInstanceOutput, err := uc.Execute(validInput, mocks.UserID)

	assert.Error(t, err, "error expected due to nonexistent database technology")
	assert.EqualError(t, err, common.ErrTechnologyNotFound.Error())
	assert.Nil(t, dbInstanceOutput)
	dbInstanceStorage.AssertNumberOfCalls(t, "Exists", 1)
	ecosystemStorage.AssertNumberOfCalls(t, "FindByID", 1)
	techStorage.AssertNumberOfCalls(t, "FindByID", 1)
}

func TestGivenAErrorInDbWhileFetchingTechnologyId_WhenExecuteCreate_ThenShouldReturnError(t *testing.T) {
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	ecosystemStorage := new(mocks.EcosystemStorageMock)
	techStorage := new(mocks.TechnologyStorageMock)
	dbInstanceStorage.On("Exists", validInput.Host, validInput.Port).Return(false, nil).Once()
	ecosystemStorage.On("FindByID", validInput.EcosystemID).Return(&entity.Ecosystem{}, nil).Once()
	techStorage.On("FindByID", validInput.DatabaseTechnologyID).Return(&entity.DatabaseTechnology{}, sql.ErrConnDone).Once()

	uc := NewCreateDatabaseInstanceUseCase(dbInstanceStorage, ecosystemStorage, techStorage)
	dbInstanceOutput, err := uc.Execute(validInput, mocks.UserID)

	assert.Error(t, err, "error expected due to error in db")
	assert.EqualError(t, err, sql.ErrConnDone.Error())
	assert.Nil(t, dbInstanceOutput)
	dbInstanceStorage.AssertNumberOfCalls(t, "Exists", 1)
	ecosystemStorage.AssertNumberOfCalls(t, "FindByID", 1)
	techStorage.AssertNumberOfCalls(t, "FindByID", 1)
}

func TestGivenAnInvalidInput_WhenExecuteCreate_ThenShouldReturnError(t *testing.T) {
	input := dto.DatabaseInstanceInputDTO{
		Name: "PostgreSQL - Azure",
	}
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	ecosystemStorage := new(mocks.EcosystemStorageMock)
	techStorage := new(mocks.TechnologyStorageMock)

	uc := NewCreateDatabaseInstanceUseCase(dbInstanceStorage, ecosystemStorage, techStorage)
	dbInstanceOutput, err := uc.Execute(input, mocks.UserID)

	assert.Error(t, err, "error expected due to invalid host (not informed)")
	assert.EqualError(t, err, entity.ErrInvalidHost.Error())
	assert.Nil(t, dbInstanceOutput)
}

func TestGivenAValidInput_WhenExecuteCreate_ThenShouldSaveDatabaseInstance(t *testing.T) {
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	ecosystemStorage := new(mocks.EcosystemStorageMock)
	techStorage := new(mocks.TechnologyStorageMock)

	dbInstanceStorage.On("Exists", validInput.Host, validInput.Port).Return(false, nil).Once()
	ecosystemStorage.On("FindByID", validInput.EcosystemID).Return(&entity.Ecosystem{}, nil).Once()
	techStorage.On("FindByID", validInput.DatabaseTechnologyID).Return(&entity.DatabaseTechnology{}, nil).Once()
	dbInstanceStorage.On("Save", mock.Anything).Return(nil).Once()

	uc := NewCreateDatabaseInstanceUseCase(dbInstanceStorage, ecosystemStorage, techStorage)
	dbInstanceOutput, err := uc.Execute(validInput, mocks.UserID)

	assert.NoError(t, err, "no error expected with a valid input and not existent host and port")
	assert.NotNil(t, dbInstanceOutput, "database instance should not be nil")
	assert.Equal(t, validInput.Name, dbInstanceOutput.Name)
	assert.Equal(t, validInput.Host, dbInstanceOutput.Host)
	assert.Equal(t, validInput.Port, dbInstanceOutput.Port)
	assert.Equal(t, validInput.HostConnection, dbInstanceOutput.HostConnection)
	assert.Equal(t, validInput.PortConnection, dbInstanceOutput.PortConnection)
	assert.Equal(t, validInput.AdminUser, dbInstanceOutput.AdminUser)
	assert.Equal(t, validInput.EcosystemID, dbInstanceOutput.EcosystemID)
	assert.Equal(t, validInput.DatabaseTechnologyID, dbInstanceOutput.DatabaseTechnologyID)
	assert.Equal(t, validInput.Note, dbInstanceOutput.Note)
	assert.NotEmpty(t, dbInstanceOutput.ID)
	assert.NotEmpty(t, dbInstanceOutput.CreatedByUserID)
	assert.NotEmpty(t, dbInstanceOutput.CreatedAt)
	dbInstanceStorage.AssertNumberOfCalls(t, "Exists", 1)
	dbInstanceStorage.AssertNumberOfCalls(t, "Save", 1)
	ecosystemStorage.AssertNumberOfCalls(t, "FindByID", 1)
	techStorage.AssertNumberOfCalls(t, "FindByID", 1)
}

func TestGivenAnErrorInDbWhileSave_WhenExecuteCreate_ThenShouldReturnError(t *testing.T) {
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	ecosystemStorage := new(mocks.EcosystemStorageMock)
	techStorage := new(mocks.TechnologyStorageMock)

	dbInstanceStorage.On("Exists", validInput.Host, validInput.Port).Return(false, nil).Once()
	ecosystemStorage.On("FindByID", validInput.EcosystemID).Return(&entity.Ecosystem{}, nil).Once()
	techStorage.On("FindByID", validInput.DatabaseTechnologyID).Return(&entity.DatabaseTechnology{}, nil).Once()
	dbInstanceStorage.On("Save", mock.Anything).Return(sql.ErrConnDone).Once()

	uc := NewCreateDatabaseInstanceUseCase(dbInstanceStorage, ecosystemStorage, techStorage)
	dbInstanceOutput, err := uc.Execute(validInput, mocks.UserID)

	assert.Error(t, err, "error expected when some error in db")
	assert.EqualError(t, err, sql.ErrConnDone.Error())
	assert.Nil(t, dbInstanceOutput)
	dbInstanceStorage.AssertNumberOfCalls(t, "Exists", 1)
	dbInstanceStorage.AssertNumberOfCalls(t, "Save", 1)
	ecosystemStorage.AssertNumberOfCalls(t, "FindByID", 1)
	techStorage.AssertNumberOfCalls(t, "FindByID", 1)
}

func TestGivenAnErrorInDbWhileCheckHost_WhenExecuteCreate_ThenShouldReturnError(t *testing.T) {
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	ecosystemStorage := new(mocks.EcosystemStorageMock)
	techStorage := new(mocks.TechnologyStorageMock)
	dbInstanceStorage.On("Exists", validInput.Host, validInput.Port).Return(false, sql.ErrConnDone).Once()

	uc := NewCreateDatabaseInstanceUseCase(dbInstanceStorage, ecosystemStorage, techStorage)
	dbInstanceOutput, err := uc.Execute(validInput, mocks.UserID)

	assert.Error(t, err, "error expected when some error in db")
	assert.EqualError(t, err, sql.ErrConnDone.Error())
	assert.Nil(t, dbInstanceOutput)
	dbInstanceStorage.AssertNumberOfCalls(t, "Exists", 1)
}

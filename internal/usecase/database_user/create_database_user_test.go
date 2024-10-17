package dbuser

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	"github.com/zgsolucoes/zg-data-guard/internal/entity"
	"github.com/zgsolucoes/zg-data-guard/testdata/mocks"
)

var validInput = mocks.ValidDBUserInput

func TestGivenAnEmailAlreadyExistent_WhenExecuteCreate_ThenShouldReturnError(t *testing.T) {
	dbUserStorage := new(mocks.DatabaseUserStorageMock)
	dbRoleStorage := new(mocks.DatabaseRoleStorageMock)

	dbUserStorage.On("Exists", validInput.Email).Return(true, nil).Once()

	uc := NewCreateDatabaseUserUseCase(dbUserStorage, dbRoleStorage)
	output, err := uc.Execute(validInput, mocks.UserID)

	assert.Error(t, err, "error expected with e-mail already existent")
	assert.EqualError(t, err, ErrEmailAlreadyExists.Error())
	assert.Nil(t, output)
	dbUserStorage.AssertNumberOfCalls(t, "Exists", 1)
}

func TestGivenAnErrorInDbWhileCheckHost_WhenExecuteCreate_ThenShouldReturnError(t *testing.T) {
	dbUserStorage := new(mocks.DatabaseUserStorageMock)
	dbRoleStorage := new(mocks.DatabaseRoleStorageMock)

	dbUserStorage.On("Exists", validInput.Email).Return(false, sql.ErrConnDone).Once()

	uc := NewCreateDatabaseUserUseCase(dbUserStorage, dbRoleStorage)
	dbUserOutput, err := uc.Execute(validInput, mocks.UserID)

	assert.Error(t, err, "error expected when some error in db")
	assert.EqualError(t, err, sql.ErrConnDone.Error())
	assert.Nil(t, dbUserOutput)
	dbUserStorage.AssertNumberOfCalls(t, "Exists", 1)
}

func TestGivenAnNonexistentDatabaseRoleId_WhenExecuteCreate_ThenShouldReturnError(t *testing.T) {
	dbUserStorage := new(mocks.DatabaseUserStorageMock)
	dbRoleStorage := new(mocks.DatabaseRoleStorageMock)

	dbUserStorage.On("Exists", validInput.Email).Return(false, nil).Once()
	dbRoleStorage.On("FindByID", validInput.DatabaseRoleID).Return(&entity.DatabaseRole{}, sql.ErrNoRows).Once()

	uc := NewCreateDatabaseUserUseCase(dbUserStorage, dbRoleStorage)
	dbUserOutput, err := uc.Execute(validInput, mocks.UserID)

	assert.Error(t, err, "error expected due to nonexistent database role")
	assert.EqualError(t, err, ErrDatabaseRoleNotFound.Error())
	assert.Nil(t, dbUserOutput)
	dbUserStorage.AssertNumberOfCalls(t, "Exists", 1)
	dbRoleStorage.AssertNumberOfCalls(t, "FindByID", 1)
}

func TestGivenAErrorInDbWhileFetchingDatabaseRoleId_WhenExecuteCreate_ThenShouldReturnError(t *testing.T) {
	dbUserStorage := new(mocks.DatabaseUserStorageMock)
	dbRoleStorage := new(mocks.DatabaseRoleStorageMock)

	dbUserStorage.On("Exists", validInput.Email).Return(false, nil).Once()
	dbRoleStorage.On("FindByID", validInput.DatabaseRoleID).Return(&entity.DatabaseRole{}, sql.ErrConnDone).Once()

	uc := NewCreateDatabaseUserUseCase(dbUserStorage, dbRoleStorage)
	dbUserOutput, err := uc.Execute(validInput, mocks.UserID)

	assert.Error(t, err, "error expected due to error in db")
	assert.EqualError(t, err, sql.ErrConnDone.Error())
	assert.Nil(t, dbUserOutput)
	dbUserStorage.AssertNumberOfCalls(t, "Exists", 1)
	dbRoleStorage.AssertNumberOfCalls(t, "FindByID", 1)
}

func TestGivenAnInvalidInput_WhenExecuteCreate_ThenShouldReturnError(t *testing.T) {
	input := dto.DatabaseUserInputDTO{
		Name: "User Test",
	}
	dbUserStorage := new(mocks.DatabaseUserStorageMock)
	dbRoleStorage := new(mocks.DatabaseRoleStorageMock)

	uc := NewCreateDatabaseUserUseCase(dbUserStorage, dbRoleStorage)
	dbUserOutput, err := uc.Execute(input, mocks.UserID)

	assert.Error(t, err, "error expected due to invalid e-mail (not informed)")
	assert.EqualError(t, err, entity.ErrInvalidEmail.Error())
	assert.Nil(t, dbUserOutput)
}

func TestGivenAValidInput_WhenExecuteCreate_ThenShouldSaveDatabaseUser(t *testing.T) {
	dbUserStorage := new(mocks.DatabaseUserStorageMock)
	dbRoleStorage := new(mocks.DatabaseRoleStorageMock)
	dbUserStorage.On("Exists", validInput.Email).Return(false, nil).Once()
	dbRoleStorage.On("FindByID", validInput.DatabaseRoleID).Return(&entity.DatabaseRole{}, nil).Once()
	dbUserStorage.On("Save", mock.Anything).Return(nil).Once()

	uc := NewCreateDatabaseUserUseCase(dbUserStorage, dbRoleStorage)
	dbUserOutput, err := uc.Execute(validInput, mocks.UserID)

	assert.NoError(t, err, "no error expected with a valid input and not existent e-mail")
	assert.NotNil(t, dbUserOutput, "database user should not be nil")
	assert.Equal(t, validInput.Name, dbUserOutput.Name)
	assert.Equal(t, validInput.Email, dbUserOutput.Email)
	assert.Equal(t, "foobar", dbUserOutput.Username)
	assert.Equal(t, validInput.Team, dbUserOutput.Team)
	assert.Equal(t, validInput.Position, dbUserOutput.Position)
	assert.Equal(t, validInput.DatabaseRoleID, dbUserOutput.DatabaseRoleID)
	assert.Equal(t, true, dbUserOutput.Enabled)
	assert.NotEmpty(t, dbUserOutput.Password)
	assert.NotEmpty(t, dbUserOutput.ID)
	assert.NotEmpty(t, dbUserOutput.CreatedByUserID)
	assert.NotEmpty(t, dbUserOutput.CreatedAt)
	dbUserStorage.AssertNumberOfCalls(t, "Exists", 1)
	dbUserStorage.AssertNumberOfCalls(t, "Save", 1)
	dbRoleStorage.AssertNumberOfCalls(t, "FindByID", 1)
}

func TestGivenAnErrorInDbWhileSave_WhenExecuteCreate_ThenShouldReturnError(t *testing.T) {
	dbUserStorage := new(mocks.DatabaseUserStorageMock)
	dbRoleStorage := new(mocks.DatabaseRoleStorageMock)

	dbUserStorage.On("Exists", validInput.Email).Return(false, nil).Once()
	dbRoleStorage.On("FindByID", validInput.DatabaseRoleID).Return(&entity.DatabaseRole{}, nil).Once()
	dbUserStorage.On("Save", mock.Anything).Return(sql.ErrConnDone).Once()

	uc := NewCreateDatabaseUserUseCase(dbUserStorage, dbRoleStorage)
	dbUserOutput, err := uc.Execute(validInput, mocks.UserID)

	assert.Error(t, err, "error expected when some error in db")
	assert.EqualError(t, err, sql.ErrConnDone.Error())
	assert.Nil(t, dbUserOutput)
	dbUserStorage.AssertNumberOfCalls(t, "Exists", 1)
	dbUserStorage.AssertNumberOfCalls(t, "Save", 1)
	dbRoleStorage.AssertNumberOfCalls(t, "FindByID", 1)
}

package tech

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	"github.com/zgsolucoes/zg-data-guard/internal/entity"
	"github.com/zgsolucoes/zg-data-guard/testdata/mocks"
)

var (
	validInput = dto.TechnologyInputDTO{
		Name:    "PostgreSQL",
		Version: "15",
	}
)

func TestGivenACodeAlreadyExistent_WhenExecuteCreate_ThenShouldReturnError(t *testing.T) {
	dbTechnologyStorage := new(mocks.TechnologyStorageMock)
	dbTechnologyStorage.On("Exists", validInput.Name, validInput.Version).Return(true, nil).Once()
	uc := NewCreateTechnologyUseCase(dbTechnologyStorage)

	output, err := uc.Execute(validInput, mocks.UserID)
	assert.Error(t, err, "error expected with name and version already existent")
	assert.EqualError(t, err, ErrTechnologyAlreadyExists.Error())
	assert.Nil(t, output)
	dbTechnologyStorage.AssertNumberOfCalls(t, "Exists", 1)
}

func TestGivenAnInvalidInput_WhenExecuteCreate_ThenShouldReturnError(t *testing.T) {
	input := dto.TechnologyInputDTO{
		Version: "16",
	}
	dbTechnologyStorage := new(mocks.TechnologyStorageMock)
	uc := NewCreateTechnologyUseCase(dbTechnologyStorage)

	technologyOutput, err := uc.Execute(input, mocks.UserID)
	assert.Error(t, err, "error expected due to invalid name (not informed)")
	assert.EqualError(t, err, entity.ErrInvalidName.Error())
	assert.Nil(t, technologyOutput)
}

func TestGivenAValidInput_WhenExecuteCreate_ThenShouldSaveTechnology(t *testing.T) {
	dbTechnologyStorage := new(mocks.TechnologyStorageMock)
	dbTechnologyStorage.On("Exists", validInput.Name, validInput.Version).Return(false, nil).Once()
	dbTechnologyStorage.On("Save", mock.Anything).Return(nil).Once()
	uc := NewCreateTechnologyUseCase(dbTechnologyStorage)

	technologyOutput, err := uc.Execute(validInput, mocks.UserID)
	assert.NoError(t, err, "no error expected with a valid input and not existent name and version")
	assert.NotNil(t, technologyOutput, "technology should not be nil")
	assert.Equal(t, validInput.Name, technologyOutput.Name)
	assert.Equal(t, validInput.Version, technologyOutput.Version)
	assert.NotEmpty(t, technologyOutput.ID)
	assert.NotEmpty(t, technologyOutput.CreatedByUserID)
	assert.NotEmpty(t, technologyOutput.CreatedAt)
	dbTechnologyStorage.AssertNumberOfCalls(t, "Exists", 1)
	dbTechnologyStorage.AssertNumberOfCalls(t, "Save", 1)
}

func TestGivenAnErrorInDbWhileSave_WhenExecuteCreate_ThenShouldReturnError(t *testing.T) {
	dbTechnologyStorage := new(mocks.TechnologyStorageMock)
	dbTechnologyStorage.On("Exists", validInput.Name, validInput.Version).Return(false, nil).Once()
	dbTechnologyStorage.On("Save", mock.Anything).Return(sql.ErrConnDone).Once()
	uc := NewCreateTechnologyUseCase(dbTechnologyStorage)

	technologyOutput, err := uc.Execute(validInput, mocks.UserID)
	assert.Error(t, err, "error expected when some error in db")
	assert.EqualError(t, err, sql.ErrConnDone.Error())
	assert.Nil(t, technologyOutput)
	dbTechnologyStorage.AssertNumberOfCalls(t, "Exists", 1)
	dbTechnologyStorage.AssertNumberOfCalls(t, "Save", 1)
}

func TestGivenAnErrorInDbWhileCheckCode_WhenExecuteCreate_ThenShouldReturnError(t *testing.T) {
	dbTechnologyStorage := new(mocks.TechnologyStorageMock)
	dbTechnologyStorage.On("Exists", validInput.Name, validInput.Version).Return(false, sql.ErrConnDone).Once()
	uc := NewCreateTechnologyUseCase(dbTechnologyStorage)

	technologyOutput, err := uc.Execute(validInput, mocks.UserID)

	assert.Error(t, err, "error expected when some error in db")
	assert.EqualError(t, err, sql.ErrConnDone.Error())
	assert.Nil(t, technologyOutput)
	dbTechnologyStorage.AssertNumberOfCalls(t, "Exists", 1)
}

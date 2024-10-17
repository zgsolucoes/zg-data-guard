package tech

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	"github.com/zgsolucoes/zg-data-guard/internal/entity"
	"github.com/zgsolucoes/zg-data-guard/testdata/mocks"
)

var (
	updatedInput = dto.TechnologyInputDTO{
		Name:    "Oracle",
		Version: "11",
	}
)

func TestGivenAnNonexistentId_WhenExecuteUpdate_ThenShouldReturnError(t *testing.T) {
	technologyStorage := new(mocks.TechnologyStorageMock)
	technologyStorage.On("FindByID", mocks.TechnologyId).Return(&entity.DatabaseTechnology{}, sql.ErrNoRows).Once()
	uc := NewUpdateTechnologyUseCase(technologyStorage)

	technologyOutput, err := uc.Execute(validInput, mocks.TechnologyId, mocks.UserID)

	assert.Error(t, err, "error expected when technology not found")
	assert.EqualError(t, err, ErrTechnologyNotFound.Error())
	assert.Nil(t, technologyOutput)
	technologyStorage.AssertNumberOfCalls(t, "FindByID", 1)
}

func TestGivenAnErrorInDbWhileFindById_WhenExecuteUpdate_ThenShouldReturnError(t *testing.T) {
	technologyStorage := new(mocks.TechnologyStorageMock)
	technologyStorage.On("FindByID", mocks.TechnologyId).Return(&entity.DatabaseTechnology{}, sql.ErrConnDone).Once()
	uc := NewUpdateTechnologyUseCase(technologyStorage)

	technologyOutput, err := uc.Execute(validInput, mocks.TechnologyId, mocks.UserID)

	assert.Error(t, err, "error expected when some error in db")
	assert.EqualError(t, err, sql.ErrConnDone.Error())
	assert.Nil(t, technologyOutput)
	technologyStorage.AssertNumberOfCalls(t, "FindByID", 1)
}

func TestGivenAnErrorWhenCheckingIfNameAndVersionExists_WhenExecuteUpdate_ThenShouldReturnError(t *testing.T) {
	technologyStorage := new(mocks.TechnologyStorageMock)
	technology, _ := entity.NewDatabaseTechnology(validInput.Name, validInput.Version, mocks.UserID)
	technologyStorage.On("FindByID", technology.ID.String()).Return(technology, nil).Once()
	technologyStorage.On("Exists", updatedInput.Name, updatedInput.Version).Return(false, sql.ErrConnDone).Once()
	uc := NewUpdateTechnologyUseCase(technologyStorage)

	technologyOutput, err := uc.Execute(updatedInput, technology.ID.String(), mocks.UserID)

	assert.Error(t, err, "error expected when some error in db")
	assert.EqualError(t, err, sql.ErrConnDone.Error())
	assert.Nil(t, technologyOutput)
	technologyStorage.AssertNumberOfCalls(t, "FindByID", 1)
	technologyStorage.AssertNumberOfCalls(t, "Exists", 1)
}

func TestGivenAnInputWithAlreadyExistingNameAndVersion_WhenExecuteUpdate_ThenShouldReturnError(t *testing.T) {
	technologyStorage := new(mocks.TechnologyStorageMock)
	technology, _ := entity.NewDatabaseTechnology(validInput.Name, validInput.Version, mocks.UserID)
	technologyStorage.On("FindByID", technology.ID.String()).Return(technology, nil).Once()
	technologyStorage.On("Exists", updatedInput.Name, updatedInput.Version).Return(true, nil).Once()
	uc := NewUpdateTechnologyUseCase(technologyStorage)

	technologyOutput, err := uc.Execute(updatedInput, technology.ID.String(), mocks.UserID)

	assert.Error(t, err, "error expected when the name and version already exists in db")
	assert.EqualError(t, err, ErrTechnologyAlreadyExists.Error())
	assert.Nil(t, technologyOutput)
	technologyStorage.AssertNumberOfCalls(t, "FindByID", 1)
	technologyStorage.AssertNumberOfCalls(t, "Exists", 1)
}

func TestGivenAnErrorInDbWhileUpdate_WhenExecuteUpdate_ThenShouldReturnError(t *testing.T) {
	technologyStorage := new(mocks.TechnologyStorageMock)
	technology, _ := entity.NewDatabaseTechnology(validInput.Name, validInput.Version, mocks.UserID)
	technologyStorage.On("FindByID", technology.ID.String()).Return(technology, nil).Once()
	technologyStorage.On("Exists", updatedInput.Name, updatedInput.Version).Return(false, nil).Once()
	technologyStorage.On("Update", technology).Return(sql.ErrConnDone).Once()
	uc := NewUpdateTechnologyUseCase(technologyStorage)

	technologyOutput, err := uc.Execute(updatedInput, technology.ID.String(), mocks.UserID)

	assert.Error(t, err, "error expected when some error in db")
	assert.EqualError(t, err, sql.ErrConnDone.Error())
	assert.Nil(t, technologyOutput)
	technologyStorage.AssertNumberOfCalls(t, "FindByID", 1)
	technologyStorage.AssertNumberOfCalls(t, "Exists", 1)
	technologyStorage.AssertNumberOfCalls(t, "Update", 1)
}
func TestGivenAValidIdAndInput_WhenExecuteUpdate_ThenShouldReturnTechnology(t *testing.T) {
	technologyStorage := new(mocks.TechnologyStorageMock)
	technology, _ := entity.NewDatabaseTechnology(validInput.Name, validInput.Version, mocks.UserID)
	technologyStorage.On("FindByID", technology.ID.String()).Return(technology, nil).Once()
	technologyStorage.On("Exists", updatedInput.Name, updatedInput.Version).Return(false, nil).Once()
	technologyStorage.On("Update", technology).Return(nil).Once()
	uc := NewUpdateTechnologyUseCase(technologyStorage)

	technologyOutput, err := uc.Execute(updatedInput, technology.ID.String(), mocks.UserID)

	assert.NoError(t, err, "no error expected with an existent id and valid input")
	assert.NotNil(t, technologyOutput, "technology should not be nil")
	assert.Equal(t, technology.ID.String(), technologyOutput.ID)
	assert.Equal(t, updatedInput.Name, technologyOutput.Name)
	assert.Equal(t, updatedInput.Version, technologyOutput.Version)
	assert.Equal(t, technology.CreatedAt, technologyOutput.CreatedAt)
	assert.Equal(t, technology.UpdatedAt, *technologyOutput.UpdatedAt)
	technologyStorage.AssertNumberOfCalls(t, "FindByID", 1)
	technologyStorage.AssertNumberOfCalls(t, "Exists", 1)
	technologyStorage.AssertNumberOfCalls(t, "Update", 1)
}

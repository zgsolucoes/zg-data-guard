package tech

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zgsolucoes/zg-data-guard/internal/entity"
	"github.com/zgsolucoes/zg-data-guard/internal/usecase/common"
	"github.com/zgsolucoes/zg-data-guard/testdata/mocks"
)

func TestGivenAnNonexistentId_WhenExecuteGet_ThenShouldReturnError(t *testing.T) {
	technologyStorage := new(mocks.TechnologyStorageMock)
	userStorage := new(mocks.UserStorageMock)
	technologyStorage.On("FindByID", mocks.TechnologyId).Return(&entity.DatabaseTechnology{}, sql.ErrNoRows).Once()
	uc := NewGetTechnologyUseCase(technologyStorage, userStorage)

	technologyOutput, err := uc.Execute(mocks.TechnologyId)

	assert.Error(t, err, "error expected when technology not found")
	assert.EqualError(t, err, common.ErrTechnologyNotFound.Error())
	assert.Nil(t, technologyOutput)
	technologyStorage.AssertNumberOfCalls(t, "FindByID", 1)
}

func TestGivenAnErrorInDb_WhenExecuteGet_ThenShouldReturnError(t *testing.T) {
	technologyStorage := new(mocks.TechnologyStorageMock)
	userStorage := new(mocks.UserStorageMock)
	technologyStorage.On("FindByID", mocks.TechnologyId).Return(&entity.DatabaseTechnology{}, sql.ErrConnDone).Once()
	uc := NewGetTechnologyUseCase(technologyStorage, userStorage)

	technologyOutput, err := uc.Execute(mocks.TechnologyId)

	assert.Error(t, err, "error expected when some error in db")
	assert.EqualError(t, err, sql.ErrConnDone.Error())
	assert.Nil(t, technologyOutput)
	technologyStorage.AssertNumberOfCalls(t, "FindByID", 1)
}

func TestGivenAnErrorWhileFindingUser_WhenExecuteGet_ThenShouldReturnError(t *testing.T) {
	technologyStorage := new(mocks.TechnologyStorageMock)
	userStorage := new(mocks.UserStorageMock)
	technology, _ := entity.NewDatabaseTechnology(validInput.Name, validInput.Version, mocks.UserID)
	technologyStorage.On("FindByID", technology.ID.String()).Return(technology, nil).Once()
	userStorage.On("FindByID", technology.CreatedByUserID).Return(&entity.ApplicationUser{}, sql.ErrConnDone).Once()
	uc := NewGetTechnologyUseCase(technologyStorage, userStorage)

	technologyOutput, err := uc.Execute(technology.ID.String())

	assert.Error(t, err, "error expected when some error in db while finding user")
	assert.EqualError(t, err, sql.ErrConnDone.Error())
	assert.Nil(t, technologyOutput)
	technologyStorage.AssertNumberOfCalls(t, "FindByID", 1)
	userStorage.AssertNumberOfCalls(t, "FindByID", 1)
}

func TestGivenAValidId_WhenExecuteGet_ThenShouldReturnTechnology(t *testing.T) {
	technologyStorage := new(mocks.TechnologyStorageMock)
	userStorage := new(mocks.UserStorageMock)
	user, _ := entity.NewApplicationUser("Foo Bar", "foobar@email.com")
	technology, _ := entity.NewDatabaseTechnology(validInput.Name, validInput.Version, user.ID.String())
	technologyStorage.On("FindByID", technology.ID.String()).Return(technology, nil).Once()

	userStorage.On("FindByID", technology.CreatedByUserID).Return(user, nil).Once()
	uc := NewGetTechnologyUseCase(technologyStorage, userStorage)

	technologyOutput, err := uc.Execute(technology.ID.String())

	assert.NoError(t, err, "no error expected with an existent id")
	assert.NotNil(t, technologyOutput, "technology should not be nil")
	assert.Equal(t, technology.ID.String(), technologyOutput.ID)
	assert.Equal(t, technology.Name, technologyOutput.Name)
	assert.Equal(t, technology.Version, technologyOutput.Version)
	assert.Equal(t, technology.CreatedAt, technologyOutput.CreatedAt)
	assert.Equal(t, user.Name, technologyOutput.CreatedByUser)
	assert.Equal(t, technology.UpdatedAt, *technologyOutput.UpdatedAt)
	technologyStorage.AssertNumberOfCalls(t, "FindByID", 1)
	userStorage.AssertNumberOfCalls(t, "FindByID", 1)
}

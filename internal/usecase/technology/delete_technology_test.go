package tech

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zgsolucoes/zg-data-guard/internal/usecase/common"
	"github.com/zgsolucoes/zg-data-guard/testdata/mocks"
)

func TestGivenAnNonexistentId_WhenExecuteDelete_ThenShouldReturnError(t *testing.T) {
	technologyStorage := new(mocks.TechnologyStorageMock)
	technologyStorage.On("Delete", mocks.TechnologyId).Return(sql.ErrNoRows).Once()
	uc := NewDeleteTechnologyUseCase(technologyStorage)

	err := uc.Execute(mocks.TechnologyId, mocks.UserID)

	assert.Error(t, err, "error expected when technology not found")
	assert.EqualError(t, err, common.ErrTechnologyNotFound.Error())
	technologyStorage.AssertNumberOfCalls(t, "Delete", 1)
}

func TestGivenAnErrorInDb_WhenExecuteDelete_ThenShouldReturnError(t *testing.T) {
	technologyStorage := new(mocks.TechnologyStorageMock)
	technologyStorage.On("Delete", mocks.TechnologyId).Return(sql.ErrConnDone).Once()
	uc := NewDeleteTechnologyUseCase(technologyStorage)

	err := uc.Execute(mocks.TechnologyId, mocks.UserID)

	assert.Error(t, err, "error expected when some error in db")
	assert.EqualError(t, err, sql.ErrConnDone.Error())
	technologyStorage.AssertNumberOfCalls(t, "Delete", 1)
}

func TestGivenAValidId_WhenExecuteDelete_ThenShouldDeleteTechnology(t *testing.T) {
	technologyStorage := new(mocks.TechnologyStorageMock)
	technologyStorage.On("Delete", mocks.TechnologyId).Return(nil).Once()
	uc := NewDeleteTechnologyUseCase(technologyStorage)

	err := uc.Execute(mocks.TechnologyId, mocks.UserID)

	assert.NoError(t, err, "no error expected with an existent id")
	technologyStorage.AssertNumberOfCalls(t, "Delete", 1)
}

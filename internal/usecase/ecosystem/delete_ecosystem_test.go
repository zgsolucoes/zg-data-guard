package ecosystem

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zgsolucoes/zg-data-guard/testdata/mocks"
)

func TestGivenAnNonexistentId_WhenExecuteDelete_ThenShouldReturnError(t *testing.T) {
	ecosystemStorage := new(mocks.EcosystemStorageMock)
	ecosystemStorage.On("Delete", mocks.EcosystemId).Return(sql.ErrNoRows).Once()
	uc := NewDeleteEcosystemUseCase(ecosystemStorage)

	err := uc.Execute(mocks.EcosystemId, mocks.UserID)

	assert.Error(t, err, "error expected when ecosystem not found")
	assert.EqualError(t, err, ErrEcosystemNotFound.Error())
	ecosystemStorage.AssertNumberOfCalls(t, "Delete", 1)
}

func TestGivenAnErrorInDb_WhenExecuteDelete_ThenShouldReturnError(t *testing.T) {
	ecosystemStorage := new(mocks.EcosystemStorageMock)
	ecosystemStorage.On("Delete", mocks.EcosystemId).Return(sql.ErrConnDone).Once()
	uc := NewDeleteEcosystemUseCase(ecosystemStorage)

	err := uc.Execute(mocks.EcosystemId, mocks.UserID)

	assert.Error(t, err, "error expected when some error in db")
	assert.EqualError(t, err, sql.ErrConnDone.Error())
	ecosystemStorage.AssertNumberOfCalls(t, "Delete", 1)
}

func TestGivenAValidId_WhenExecuteDelete_ThenShouldDeleteEcosystem(t *testing.T) {
	ecosystemStorage := new(mocks.EcosystemStorageMock)
	ecosystemStorage.On("Delete", mocks.EcosystemId).Return(nil).Once()
	uc := NewDeleteEcosystemUseCase(ecosystemStorage)

	err := uc.Execute(mocks.EcosystemId, mocks.UserID)

	assert.NoError(t, err, "no error expected with an existent id")
	ecosystemStorage.AssertNumberOfCalls(t, "Delete", 1)
}

package ecosystem

import (
	"database/sql"
	"testing"

	"github.com/zgsolucoes/zg-data-guard/internal/entity"
	"github.com/zgsolucoes/zg-data-guard/testdata/mocks"

	"github.com/stretchr/testify/assert"
)

func TestGivenAnNonexistentId_WhenExecuteGet_ThenShouldReturnError(t *testing.T) {
	ecosystemStorage := new(mocks.EcosystemStorageMock)
	userStorage := new(mocks.UserStorageMock)
	ecosystemStorage.On("FindByID", mocks.EcosystemId).Return(&entity.Ecosystem{}, sql.ErrNoRows).Once()
	uc := NewGetEcosystemUseCase(ecosystemStorage, userStorage)

	ecosystemOutput, err := uc.Execute(mocks.EcosystemId)

	assert.Error(t, err, "error expected when ecosystem not found")
	assert.EqualError(t, err, ErrEcosystemNotFound.Error())
	assert.Nil(t, ecosystemOutput)
	ecosystemStorage.AssertNumberOfCalls(t, "FindByID", 1)
}

func TestGivenAnErrorInDb_WhenExecuteGet_ThenShouldReturnError(t *testing.T) {
	ecosystemStorage := new(mocks.EcosystemStorageMock)
	userStorage := new(mocks.UserStorageMock)
	ecosystemStorage.On("FindByID", mocks.EcosystemId).Return(&entity.Ecosystem{}, sql.ErrConnDone).Once()
	uc := NewGetEcosystemUseCase(ecosystemStorage, userStorage)

	ecosystemOutput, err := uc.Execute(mocks.EcosystemId)

	assert.Error(t, err, "error expected when some error in db")
	assert.EqualError(t, err, sql.ErrConnDone.Error())
	assert.Nil(t, ecosystemOutput)
	ecosystemStorage.AssertNumberOfCalls(t, "FindByID", 1)
}

func TestGivenAnErrorWhileFindingUser_WhenExecuteGet_ThenShouldReturnError(t *testing.T) {
	ecosystemStorage := new(mocks.EcosystemStorageMock)
	userStorage := new(mocks.UserStorageMock)
	ecosystem, _ := entity.NewEcosystem(validInput.Code, validInput.DisplayName, mocks.UserID)
	ecosystemStorage.On("FindByID", ecosystem.ID.String()).Return(ecosystem, nil).Once()
	userStorage.On("FindByID", ecosystem.CreatedByUserID).Return(&entity.ApplicationUser{}, sql.ErrConnDone).Once()
	uc := NewGetEcosystemUseCase(ecosystemStorage, userStorage)

	ecosystemOutput, err := uc.Execute(ecosystem.ID.String())

	assert.Error(t, err, "error expected when some error in db while finding user")
	assert.EqualError(t, err, sql.ErrConnDone.Error())
	assert.Nil(t, ecosystemOutput)
	ecosystemStorage.AssertNumberOfCalls(t, "FindByID", 1)
	userStorage.AssertNumberOfCalls(t, "FindByID", 1)
}

func TestGivenAValidId_WhenExecuteGet_ThenShouldReturnEcosystem(t *testing.T) {
	ecosystemStorage := new(mocks.EcosystemStorageMock)
	userStorage := new(mocks.UserStorageMock)
	user, _ := entity.NewApplicationUser("Foo Bar", "foobar@email.com")
	ecosystem, _ := entity.NewEcosystem(validInput.Code, validInput.DisplayName, user.ID.String())
	ecosystemStorage.On("FindByID", ecosystem.ID.String()).Return(ecosystem, nil).Once()

	userStorage.On("FindByID", ecosystem.CreatedByUserID).Return(user, nil).Once()
	uc := NewGetEcosystemUseCase(ecosystemStorage, userStorage)

	ecosystemOutput, err := uc.Execute(ecosystem.ID.String())

	assert.NoError(t, err, "no error expected with an existent id")
	assert.NotNil(t, ecosystemOutput, "ecosystem should not be nil")
	assert.Equal(t, ecosystem.ID.String(), ecosystemOutput.ID)
	assert.Equal(t, ecosystem.Code, ecosystemOutput.Code)
	assert.Equal(t, ecosystem.DisplayName, ecosystemOutput.DisplayName)
	assert.Equal(t, ecosystem.CreatedAt, ecosystemOutput.CreatedAt)
	assert.Equal(t, user.Name, ecosystemOutput.CreatedByUser)
	assert.Equal(t, ecosystem.UpdatedAt, *ecosystemOutput.UpdatedAt)
	ecosystemStorage.AssertNumberOfCalls(t, "FindByID", 1)
	userStorage.AssertNumberOfCalls(t, "FindByID", 1)
}

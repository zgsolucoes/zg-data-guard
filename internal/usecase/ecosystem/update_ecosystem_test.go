package ecosystem

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	"github.com/zgsolucoes/zg-data-guard/internal/entity"
	"github.com/zgsolucoes/zg-data-guard/testdata/mocks"
)

var (
	updatedInput = dto.EcosystemInputDTO{
		Code:        "foobar",
		DisplayName: "Foo Bar",
	}
)

func TestGivenAnNonexistentId_WhenExecuteUpdate_ThenShouldReturnError(t *testing.T) {
	ecosystemStorage := new(mocks.EcosystemStorageMock)
	ecosystemStorage.On("FindByID", mocks.EcosystemId).Return(&entity.Ecosystem{}, sql.ErrNoRows).Once()
	uc := NewUpdateEcosystemUseCase(ecosystemStorage)

	ecosystemOutput, err := uc.Execute(validInput, mocks.EcosystemId, mocks.UserID)

	assert.Error(t, err, "error expected when ecosystem not found")
	assert.EqualError(t, err, ErrEcosystemNotFound.Error())
	assert.Nil(t, ecosystemOutput)
	ecosystemStorage.AssertNumberOfCalls(t, "FindByID", 1)
}

func TestGivenAnErrorInDbWhileFindById_WhenExecuteUpdate_ThenShouldReturnError(t *testing.T) {
	ecosystemStorage := new(mocks.EcosystemStorageMock)
	ecosystemStorage.On("FindByID", mocks.EcosystemId).Return(&entity.Ecosystem{}, sql.ErrConnDone).Once()
	uc := NewUpdateEcosystemUseCase(ecosystemStorage)

	ecosystemOutput, err := uc.Execute(validInput, mocks.EcosystemId, mocks.UserID)

	assert.Error(t, err, "error expected when some error in db")
	assert.EqualError(t, err, sql.ErrConnDone.Error())
	assert.Nil(t, ecosystemOutput)
	ecosystemStorage.AssertNumberOfCalls(t, "FindByID", 1)
}

func TestGivenAnErrorWhileCheckingCodeExists_WhenExecuteUpdate_ThenShouldReturnError(t *testing.T) {
	ecosystemStorage := new(mocks.EcosystemStorageMock)
	ecosystem, _ := entity.NewEcosystem(validInput.Code, validInput.DisplayName, mocks.UserID)
	ecosystemStorage.On("FindByID", ecosystem.ID.String()).Return(ecosystem, nil).Once()
	ecosystemStorage.On("CheckCodeExists", updatedInput.Code).Return(false, sql.ErrConnDone).Once()
	uc := NewUpdateEcosystemUseCase(ecosystemStorage)

	ecosystemOutput, err := uc.Execute(updatedInput, ecosystem.ID.String(), mocks.UserID)

	assert.Error(t, err, "error expected when some error in db")
	assert.EqualError(t, err, sql.ErrConnDone.Error())
	assert.Nil(t, ecosystemOutput)
	ecosystemStorage.AssertNumberOfCalls(t, "FindByID", 1)
	ecosystemStorage.AssertNumberOfCalls(t, "CheckCodeExists", 1)
}

func TestGivenAnInputWithAlreadyExistingCode_WhenExecuteUpdate_ThenShouldReturnError(t *testing.T) {
	ecosystemStorage := new(mocks.EcosystemStorageMock)
	ecosystem, _ := entity.NewEcosystem(validInput.Code, validInput.DisplayName, mocks.UserID)
	ecosystemStorage.On("FindByID", ecosystem.ID.String()).Return(ecosystem, nil).Once()
	ecosystemStorage.On("CheckCodeExists", updatedInput.Code).Return(true, nil).Once()
	uc := NewUpdateEcosystemUseCase(ecosystemStorage)

	ecosystemOutput, err := uc.Execute(updatedInput, ecosystem.ID.String(), mocks.UserID)

	assert.Error(t, err, "error expected when the code already exists in db")
	assert.EqualError(t, err, ErrCodeAlreadyExists.Error())
	assert.Nil(t, ecosystemOutput)
	ecosystemStorage.AssertNumberOfCalls(t, "FindByID", 1)
	ecosystemStorage.AssertNumberOfCalls(t, "CheckCodeExists", 1)
}

func TestGivenAnErrorInDbWhileUpdate_WhenExecuteUpdate_ThenShouldReturnError(t *testing.T) {
	ecosystemStorage := new(mocks.EcosystemStorageMock)
	ecosystem, _ := entity.NewEcosystem(validInput.Code, validInput.DisplayName, mocks.UserID)
	ecosystemStorage.On("FindByID", ecosystem.ID.String()).Return(ecosystem, nil).Once()
	ecosystemStorage.On("CheckCodeExists", updatedInput.Code).Return(false, nil).Once()
	ecosystemStorage.On("Update", ecosystem).Return(sql.ErrConnDone).Once()
	uc := NewUpdateEcosystemUseCase(ecosystemStorage)

	ecosystemOutput, err := uc.Execute(updatedInput, ecosystem.ID.String(), mocks.UserID)

	assert.Error(t, err, "error expected when some error in db")
	assert.EqualError(t, err, sql.ErrConnDone.Error())
	assert.Nil(t, ecosystemOutput)
	ecosystemStorage.AssertNumberOfCalls(t, "FindByID", 1)
	ecosystemStorage.AssertNumberOfCalls(t, "CheckCodeExists", 1)
	ecosystemStorage.AssertNumberOfCalls(t, "Update", 1)
}
func TestGivenAValidIdAndInput_WhenExecuteUpdate_ThenShouldReturnEcosystem(t *testing.T) {
	ecosystemStorage := new(mocks.EcosystemStorageMock)
	ecosystem, _ := entity.NewEcosystem(validInput.Code, validInput.DisplayName, mocks.UserID)
	ecosystemStorage.On("FindByID", ecosystem.ID.String()).Return(ecosystem, nil).Once()
	ecosystemStorage.On("CheckCodeExists", updatedInput.Code).Return(false, nil).Once()
	ecosystemStorage.On("Update", ecosystem).Return(nil).Once()
	uc := NewUpdateEcosystemUseCase(ecosystemStorage)

	ecosystemOutput, err := uc.Execute(updatedInput, ecosystem.ID.String(), mocks.UserID)

	assert.NoError(t, err, "no error expected with an existent id and valid input")
	assert.NotNil(t, ecosystemOutput, "ecosystem should not be nil")
	assert.Equal(t, ecosystem.ID.String(), ecosystemOutput.ID)
	assert.Equal(t, ecosystemOutput.Code, updatedInput.Code)
	assert.Equal(t, ecosystemOutput.DisplayName, updatedInput.DisplayName)
	assert.Equal(t, ecosystem.CreatedAt, ecosystemOutput.CreatedAt)
	assert.Equal(t, ecosystem.UpdatedAt, *ecosystemOutput.UpdatedAt)
	ecosystemStorage.AssertNumberOfCalls(t, "FindByID", 1)
	ecosystemStorage.AssertNumberOfCalls(t, "CheckCodeExists", 1)
	ecosystemStorage.AssertNumberOfCalls(t, "Update", 1)
}

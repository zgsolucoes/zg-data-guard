package ecosystem

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
	validInput = dto.EcosystemInputDTO{
		Code:        "test",
		DisplayName: "Test",
	}
)

func TestGivenACodeAlreadyExistent_WhenExecuteCreate_ThenShouldReturnError(t *testing.T) {
	ecosystemStorage := new(mocks.EcosystemStorageMock)
	ecosystemStorage.On("CheckCodeExists", validInput.Code).Return(true, nil).Once()
	uc := NewCreateEcosystemUseCase(ecosystemStorage)

	ecosystemOutput, err := uc.Execute(validInput, mocks.UserID)
	assert.Error(t, err, "error expected with code already existent")
	assert.EqualError(t, err, ErrCodeAlreadyExists.Error())
	assert.Nil(t, ecosystemOutput)
	ecosystemStorage.AssertNumberOfCalls(t, "CheckCodeExists", 1)
}

func TestGivenAnInvalidInput_WhenExecuteCreate_ThenShouldReturnError(t *testing.T) {
	input := dto.EcosystemInputDTO{
		DisplayName: "Test",
	}
	ecosystemStorage := new(mocks.EcosystemStorageMock)
	uc := NewCreateEcosystemUseCase(ecosystemStorage)

	ecosystemOutput, err := uc.Execute(input, mocks.UserID)
	assert.Error(t, err, "error expected due to invalid code (not informed)")
	assert.EqualError(t, err, entity.ErrInvalidCode.Error())
	assert.Nil(t, ecosystemOutput)
}

func TestGivenAValidInput_WhenExecuteCreate_ThenShouldSaveEcosystem(t *testing.T) {
	ecosystemStorage := new(mocks.EcosystemStorageMock)
	ecosystemStorage.On("CheckCodeExists", validInput.Code).Return(false, nil).Once()
	ecosystemStorage.On("Save", mock.Anything).Return(nil).Once()
	uc := NewCreateEcosystemUseCase(ecosystemStorage)

	ecosystemOutput, err := uc.Execute(validInput, mocks.UserID)
	assert.NoError(t, err, "no error expected with a valid input and not existent code")
	assert.NotNil(t, ecosystemOutput, "ecosystem should not be nil")
	assert.Equal(t, validInput.Code, ecosystemOutput.Code)
	assert.Equal(t, validInput.DisplayName, ecosystemOutput.DisplayName)
	assert.NotEmpty(t, ecosystemOutput.ID)
	assert.NotEmpty(t, ecosystemOutput.CreatedByUserID)
	assert.NotEmpty(t, ecosystemOutput.CreatedAt)
	ecosystemStorage.AssertNumberOfCalls(t, "CheckCodeExists", 1)
	ecosystemStorage.AssertNumberOfCalls(t, "Save", 1)
}

func TestGivenAnErrorInDbWhileSave_WhenExecuteCreate_ThenShouldReturnError(t *testing.T) {
	ecosystemStorage := new(mocks.EcosystemStorageMock)
	ecosystemStorage.On("CheckCodeExists", validInput.Code).Return(false, nil).Once()
	ecosystemStorage.On("Save", mock.Anything).Return(sql.ErrConnDone).Once()
	uc := NewCreateEcosystemUseCase(ecosystemStorage)

	ecosystemOutput, err := uc.Execute(validInput, mocks.UserID)
	assert.Error(t, err, "error expected when some error in db")
	assert.EqualError(t, err, sql.ErrConnDone.Error())
	assert.Nil(t, ecosystemOutput)
	ecosystemStorage.AssertNumberOfCalls(t, "CheckCodeExists", 1)
	ecosystemStorage.AssertNumberOfCalls(t, "Save", 1)
}

func TestGivenAnErrorInDbWhileCheckCode_WhenExecuteCreate_ThenShouldReturnError(t *testing.T) {
	ecosystemStorage := new(mocks.EcosystemStorageMock)
	ecosystemStorage.On("CheckCodeExists", validInput.Code).Return(false, sql.ErrConnDone).Once()
	uc := NewCreateEcosystemUseCase(ecosystemStorage)

	ecosystemOutput, err := uc.Execute(validInput, mocks.UserID)

	assert.Error(t, err, "error expected when some error in db")
	assert.EqualError(t, err, sql.ErrConnDone.Error())
	assert.Nil(t, ecosystemOutput)
	ecosystemStorage.AssertNumberOfCalls(t, "CheckCodeExists", 1)
}

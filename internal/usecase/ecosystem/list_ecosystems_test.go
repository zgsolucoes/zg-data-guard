package ecosystem

import (
	"database/sql"
	"testing"

	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	"github.com/zgsolucoes/zg-data-guard/testdata/mocks"

	"github.com/stretchr/testify/assert"
)

func TestGivenAnErrorInDb_WhenExecuteList_ThenShouldReturnError(t *testing.T) {
	ecosystemStorage := new(mocks.EcosystemStorageMock)
	ecosystemStorage.On("FindAll", 0, 50).Return([]*dto.EcosystemOutputDTO{}, sql.ErrConnDone).Once()
	uc := NewListEcosystemsUseCase(ecosystemStorage)

	ecosystemsObtained, err := uc.Execute(0, 50)

	assert.Error(t, err, "error expected when some error in db")
	assert.EqualError(t, err, sql.ErrConnDone.Error())
	assert.Nil(t, ecosystemsObtained)
	assert.Equal(t, len(ecosystemsObtained), 0, "0 ecosystems expected")
	ecosystemStorage.AssertNumberOfCalls(t, "FindAll", 1)
}

func TestGivenSomeEcosystems_WhenExecuteList_ThenShouldListAllEcosystems(t *testing.T) {
	ecosystemDto := &dto.EcosystemOutputDTO{
		ID:            "1",
		Code:          "aws",
		DisplayName:   "AWS",
		CreatedByUser: "Foo Bar",
	}
	ecosystemDto2 := &dto.EcosystemOutputDTO{
		ID:            "2",
		Code:          "azure",
		DisplayName:   "Azure",
		CreatedByUser: "John Doe",
	}
	ecosystemDto3 := &dto.EcosystemOutputDTO{
		ID:            "3",
		Code:          "qa",
		DisplayName:   "QA",
		CreatedByUser: "Foo Bar",
	}
	ecosystemsDtos := []*dto.EcosystemOutputDTO{ecosystemDto, ecosystemDto2, ecosystemDto3}

	ecosystemStorage := new(mocks.EcosystemStorageMock)
	ecosystemStorage.On("FindAll", 0, 50).Return(ecosystemsDtos, nil).Once()
	uc := NewListEcosystemsUseCase(ecosystemStorage)

	ecosystemsObtained, err := uc.Execute(0, 50)

	assert.NoError(t, err, "no error expected ")
	assert.Equal(t, len(ecosystemsObtained), len(ecosystemsDtos), "3 ecosystems expected")
	ecosystemStorage.AssertNumberOfCalls(t, "FindAll", 1)
}

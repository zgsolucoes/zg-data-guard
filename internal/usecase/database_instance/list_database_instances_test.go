package instance

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	"github.com/zgsolucoes/zg-data-guard/testdata/mocks"
)

func TestGivenAnErrorInDb_WhenExecuteList_ThenShouldReturnError(t *testing.T) {
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	dbInstanceStorage.On("FindAllDTOs", mocks.EcosystemId, mocks.TechnologyId, []string{}).Return([]*dto.DatabaseInstanceOutputDTO{}, sql.ErrConnDone).Once()

	uc := NewListDatabaseInstancesUseCase(dbInstanceStorage)
	dbInstancesObtained, err := uc.Execute(mocks.EcosystemId, mocks.TechnologyId, false)

	assert.Error(t, err, "error expected when some error in db")
	assert.EqualError(t, err, sql.ErrConnDone.Error())
	assert.Nil(t, dbInstancesObtained)
	assert.Equal(t, len(dbInstancesObtained), 0, "0 database instances expected")
	dbInstanceStorage.AssertNumberOfCalls(t, "FindAllDTOs", 1)
}

func TestGivenSomeDbInstances_WhenExecuteList_ThenShouldListAllDbInstances(t *testing.T) {
	dbInstanceDto := &dto.DatabaseInstanceOutputDTO{
		ID:            "1",
		Name:          "Postgres - QA",
		Host:          "localhost",
		Port:          "5432",
		CreatedByUser: "Foo Bar",
		Enabled:       true,
	}
	dbInstanceDto2 := &dto.DatabaseInstanceOutputDTO{
		ID:            "2",
		Name:          "Postgres - Azure",
		Host:          "127.0.0.1",
		Port:          "5432",
		CreatedByUser: "Luiz Henrique",
		Enabled:       true,
	}
	dbInstanceDto3 := &dto.DatabaseInstanceOutputDTO{
		ID:            "3",
		Name:          "Postgres - AWS",
		Host:          "localhost",
		Port:          "5433",
		CreatedByUser: "John Doe",
		Enabled:       true,
	}
	dbInstancesDtos := []*dto.DatabaseInstanceOutputDTO{dbInstanceDto, dbInstanceDto2, dbInstanceDto3}
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	dbInstanceStorage.On("FindAllDTOsEnabled", mocks.EcosystemId, mocks.TechnologyId).Return(dbInstancesDtos, nil).Once()

	uc := NewListDatabaseInstancesUseCase(dbInstanceStorage)
	dbInstancesObtained, err := uc.Execute(mocks.EcosystemId, mocks.TechnologyId, true)

	assert.NoError(t, err, "no error expected ")
	assert.Equal(t, len(dbInstancesObtained), len(dbInstancesDtos), "3 database instances expected")
	dbInstanceStorage.AssertNumberOfCalls(t, "FindAllDTOsEnabled", 1)
}

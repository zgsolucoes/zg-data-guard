package tech

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	"github.com/zgsolucoes/zg-data-guard/testdata/mocks"
)

func TestGivenAnErrorInDb_WhenExecuteList_ThenShouldReturnError(t *testing.T) {
	technologyStorage := new(mocks.TechnologyStorageMock)
	technologyStorage.On("FindAll", 0, 50).Return([]*dto.TechnologyOutputDTO{}, sql.ErrConnDone).Once()
	uc := NewListTechnologiesUseCase(technologyStorage)

	technologiesObtained, err := uc.Execute(0, 50)

	assert.Error(t, err, "error expected when some error in db")
	assert.EqualError(t, err, sql.ErrConnDone.Error())
	assert.Nil(t, technologiesObtained)
	assert.Equal(t, len(technologiesObtained), 0, "0 technologies expected")
	technologyStorage.AssertNumberOfCalls(t, "FindAll", 1)
}

func TestGivenSomeTechnologies_WhenExecuteList_ThenShouldListAllTechnologies(t *testing.T) {
	technologyDto := &dto.TechnologyOutputDTO{
		ID:            "1",
		Name:          "PostgreSQL",
		Version:       "16",
		CreatedByUser: "Foo Bar",
	}
	technologyDto2 := &dto.TechnologyOutputDTO{
		ID:            "2",
		Name:          "MySQL",
		Version:       "7",
		CreatedByUser: "User Test",
	}
	technologyDto3 := &dto.TechnologyOutputDTO{
		ID:            "3",
		Name:          "Oracle",
		Version:       "11",
		CreatedByUser: "Foo Bar",
	}
	technologiesDtos := []*dto.TechnologyOutputDTO{technologyDto, technologyDto2, technologyDto3}

	technologyStorage := new(mocks.TechnologyStorageMock)
	technologyStorage.On("FindAll", 0, 50).Return(technologiesDtos, nil).Once()
	uc := NewListTechnologiesUseCase(technologyStorage)

	technologiesObtained, err := uc.Execute(0, 50)

	assert.NoError(t, err, "no error expected ")
	assert.Equal(t, len(technologiesObtained), len(technologiesDtos), "3 technologies expected")
	technologyStorage.AssertNumberOfCalls(t, "FindAll", 1)
}

package mocks

import (
	"github.com/stretchr/testify/mock"

	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	"github.com/zgsolucoes/zg-data-guard/internal/entity"
)

const TechnologyId = "57690512-12fc-4714-a3b1-59f1d09edc0d"

type TechnologyStorageMock struct {
	mock.Mock
}

func (m *TechnologyStorageMock) Exists(name, version string) (bool, error) {
	args := m.Called(name, version)
	return args.Bool(0), args.Error(1)
}

func (m *TechnologyStorageMock) Save(dbTechnology *entity.DatabaseTechnology) error {
	args := m.Called(dbTechnology)
	return args.Error(0)
}

func (m *TechnologyStorageMock) FindByID(id string) (*entity.DatabaseTechnology, error) {
	args := m.Called(id)
	return args.Get(0).(*entity.DatabaseTechnology), args.Error(1)
}

func (m *TechnologyStorageMock) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *TechnologyStorageMock) Update(dbTechnology *entity.DatabaseTechnology) error {
	args := m.Called(dbTechnology)
	return args.Error(0)
}

func (m *TechnologyStorageMock) FindAll(page, limit int) ([]*dto.TechnologyOutputDTO, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]*dto.TechnologyOutputDTO), args.Error(1)
}

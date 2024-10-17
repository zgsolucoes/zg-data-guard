package mocks

import (
	"github.com/stretchr/testify/mock"

	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	"github.com/zgsolucoes/zg-data-guard/internal/entity"
)

const EcosystemId = "3919b8ce-5ea8-4395-ba8e-437de0615d9d"

type EcosystemStorageMock struct {
	mock.Mock
}

func (m *EcosystemStorageMock) CheckCodeExists(code string) (bool, error) {
	args := m.Called(code)
	return args.Bool(0), args.Error(1)
}

func (m *EcosystemStorageMock) Save(ecosystem *entity.Ecosystem) error {
	args := m.Called(ecosystem)
	return args.Error(0)
}

func (m *EcosystemStorageMock) Update(ecosystem *entity.Ecosystem) error {
	args := m.Called(ecosystem)
	return args.Error(0)
}

func (m *EcosystemStorageMock) FindByID(id string) (*entity.Ecosystem, error) {
	args := m.Called(id)
	return args.Get(0).(*entity.Ecosystem), args.Error(1)
}

func (m *EcosystemStorageMock) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *EcosystemStorageMock) FindAll(page, limit int) ([]*dto.EcosystemOutputDTO, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]*dto.EcosystemOutputDTO), args.Error(1)
}

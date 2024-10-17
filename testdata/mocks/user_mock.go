package mocks

import (
	"github.com/stretchr/testify/mock"

	"github.com/zgsolucoes/zg-data-guard/internal/entity"
)

const UserID = "dd42cf0c-8a91-42d7-a906-cb9313494e7d"

type UserStorageMock struct {
	mock.Mock
}

func (m *UserStorageMock) FindByEmail(email string) (*entity.ApplicationUser, error) {
	args := m.Called(email)
	return args.Get(0).(*entity.ApplicationUser), args.Error(1)
}

func (m *UserStorageMock) FindByID(id string) (*entity.ApplicationUser, error) {
	args := m.Called(id)
	return args.Get(0).(*entity.ApplicationUser), args.Error(1)
}

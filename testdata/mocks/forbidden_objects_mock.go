package mocks

import (
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/zgsolucoes/zg-data-guard/internal/entity"
)

const forbiddenDatabaseID = "c3dabc2f-55bf-4a53-870d-e998f26ae12a"

type ForbiddenObjectsStorageMock struct {
	mock.Mock
}

func (m *ForbiddenObjectsStorageMock) FindAllDatabases() ([]*entity.ForbiddenDatabase, error) {
	args := m.Called()
	return args.Get(0).([]*entity.ForbiddenDatabase), args.Error(1)
}

func BuildForbiddenDatabasesList() []*entity.ForbiddenDatabase {
	var forbiddenDatabases []*entity.ForbiddenDatabase
	forbiddenDatabases = append(forbiddenDatabases, &entity.ForbiddenDatabase{
		ID:              uuid.MustParse(forbiddenDatabaseID),
		Name:            "dummy-db-1",
		Description:     "Database 1",
		CreatedAt:       time.Now(),
		CreatedByUserID: UserID,
		UpdatedAt:       time.Now(),
	})
	forbiddenDatabases = append(forbiddenDatabases, &entity.ForbiddenDatabase{
		ID:              uuid.MustParse("96cfa8f2-2c91-4630-b556-f7a2eab84e29"),
		Name:            "dummy-db-2",
		Description:     "Database 2",
		CreatedAt:       time.Now(),
		CreatedByUserID: UserID,
		UpdatedAt:       time.Now(),
	})
	return forbiddenDatabases
}

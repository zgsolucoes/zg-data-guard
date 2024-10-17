package mocks

import (
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	"github.com/zgsolucoes/zg-data-guard/internal/entity"
)

const (
	DatabaseID = "63862219-f1c3-41c4-9938-31346773b697"
)

type DatabaseStorageMock struct {
	mock.Mock
}

func (m *DatabaseStorageMock) Save(database *entity.Database) error {
	args := m.Called(database)
	return args.Error(0)
}

func (m *DatabaseStorageMock) Update(database *entity.Database) error {
	args := m.Called(database)
	return args.Error(0)
}

func (m *DatabaseStorageMock) FindDTOByID(id string) (*dto.DatabaseOutputDTO, error) {
	args := m.Called(id)
	return args.Get(0).(*dto.DatabaseOutputDTO), args.Error(1)
}

func (m *DatabaseStorageMock) FindAll(databaseInstanceId string, ids []string) ([]*entity.Database, error) {
	args := m.Called(databaseInstanceId, ids)
	return args.Get(0).([]*entity.Database), args.Error(1)
}

func (m *DatabaseStorageMock) FindAllEnabled(databaseInstanceId string) ([]*entity.Database, error) {
	args := m.Called(databaseInstanceId)
	return args.Get(0).([]*entity.Database), args.Error(1)
}

func (m *DatabaseStorageMock) FindAllDTOs(ecosystemId, databaseInstanceId string) ([]*dto.DatabaseOutputDTO, error) {
	args := m.Called(ecosystemId, databaseInstanceId)
	return args.Get(0).([]*dto.DatabaseOutputDTO), args.Error(1)
}

func (m *DatabaseStorageMock) DeactivateAllByInstance(databaseInstanceId string) error {
	args := m.Called(databaseInstanceId)
	return args.Error(0)
}

func BuildDatabaseListSameInstanceAndOnlyEnabled() []*entity.Database {
	db1 := &entity.Database{
		Name:               "dummy-db-1",
		CurrentSize:        "10GB",
		Enabled:            true,
		RolesConfigured:    true,
		DatabaseInstanceID: DatabaseInstanceId,
	}
	db2 := &entity.Database{
		Name:               "dummy-db-2",
		CurrentSize:        "20MB",
		Enabled:            true,
		RolesConfigured:    true,
		DatabaseInstanceID: DatabaseInstanceId,
	}
	db3 := &entity.Database{
		Name:               "dummy-db-3",
		CurrentSize:        "3GB",
		Enabled:            true,
		RolesConfigured:    true,
		DatabaseInstanceID: DatabaseInstanceId,
	}
	return []*entity.Database{db1, db2, db3}
}

func BuildDatabaseList() []*entity.Database {
	db1 := &entity.Database{
		Name:               "dummy-db-1",
		CurrentSize:        "10GB",
		Enabled:            true,
		DatabaseInstanceID: DatabaseInstanceId,
	}
	db2 := &entity.Database{
		Name:               "dummy-db-2",
		CurrentSize:        "20MB",
		Enabled:            false,
		DatabaseInstanceID: DatabaseInstanceId,
	}
	db3 := &entity.Database{
		Name:               "dummy-db-3",
		CurrentSize:        "3GB",
		Enabled:            true,
		DatabaseInstanceID: "5385173e-17dd-4d29-af89-acdd89452958",
	}
	db4 := &entity.Database{
		Name:               "dummy-db-4",
		CurrentSize:        "4GB",
		Enabled:            false,
		DatabaseInstanceID: DatabaseInstanceId,
	}
	return []*entity.Database{db1, db2, db3, db4}
}

func BuildDatabaseListMixedScenarios() []*entity.Database {
	db1 := &entity.Database{
		ID:                 uuid.MustParse("112c05c4-8e33-49db-b41d-9c3b8d9e2676"),
		Name:               "db-1",
		CurrentSize:        "10GB",
		Enabled:            true,
		DatabaseInstanceID: DatabaseInstanceId,
	}
	db2 := &entity.Database{
		ID:                 uuid.MustParse("0f0e9ec3-f5c1-48fd-b16b-702a0e9f1804"),
		Name:               "db-2",
		CurrentSize:        "20MB",
		Enabled:            false,
		DatabaseInstanceID: DatabaseInstanceId,
	}
	db3 := &entity.Database{
		ID:                 uuid.MustParse("b55a4c72-159c-449f-bad4-0080944eb4da"),
		Name:               "dummy-db-3",
		CurrentSize:        "3GB",
		Enabled:            true,
		DatabaseInstanceID: DummyErrorInstanceId,
	}
	db4 := &entity.Database{
		ID:                 uuid.MustParse("431d710e-b519-4490-8574-ea9fb84a8d33"),
		Name:               "qa-db-4",
		CurrentSize:        "4GB",
		Enabled:            true,
		DatabaseInstanceID: QAInstanceId,
	}
	return []*entity.Database{db1, db2, db3, db4}
}

func BuildSettingsDatabase() *entity.Database {
	return &entity.Database{
		ID:                 uuid.MustParse(DatabaseID),
		Name:               "settings",
		CurrentSize:        "10GB",
		Enabled:            true,
		RolesConfigured:    true,
		DatabaseInstanceID: DatabaseInstanceId,
	}
}

func BuildDatabaseDTOList() []*dto.DatabaseOutputDTO {
	dbDto := &dto.DatabaseOutputDTO{
		ID:            "1",
		Name:          "settings",
		CurrentSize:   "10GB",
		CreatedByUser: "Luiz Henrique",
	}
	dbDto2 := &dto.DatabaseOutputDTO{
		ID:            "2",
		Name:          "bills-reviewer",
		CurrentSize:   "2GB",
		CreatedByUser: "Luiz Henrique",
	}
	dbDto3 := &dto.DatabaseOutputDTO{
		ID:            "3",
		Name:          "jobs",
		CurrentSize:   "10MB",
		CreatedByUser: "Luiz Henrique",
	}
	return []*dto.DatabaseOutputDTO{dbDto, dbDto2, dbDto3}
}

func BuildDatabaseDTOExample() *dto.DatabaseOutputDTO {
	systemTime := time.Now()
	return &dto.DatabaseOutputDTO{
		ID:                        DatabaseID,
		Name:                      "jobs",
		CurrentSize:               "1 GB",
		DatabaseInstanceID:        DatabaseInstanceId,
		DatabaseInstanceName:      "PostgreSQL K8s QA",
		EcosystemID:               EcosystemId,
		EcosystemName:             "QA",
		DatabaseTechnologyID:      "b3b3b3b3-5ea8-4395-ba8e-437de0615d9d",
		DatabaseTechnologyName:    "PostgreSQL",
		DatabaseTechnologyVersion: "13.7",
		Enabled:                   true,
		Description:               "Jobs database",
		CreatedByUserID:           UserID,
		CreatedByUser:             "User Test",
		CreatedAt:                 time.Now(),
		UpdatedAt:                 &systemTime,
		LastDatabaseSync:          &systemTime,
		DisabledAt:                nil,
	}
}

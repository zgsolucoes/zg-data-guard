package mocks

import (
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/zgsolucoes/zg-data-guard/internal/database/connector"
	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	"github.com/zgsolucoes/zg-data-guard/internal/entity"
)

const (
	DatabaseInstanceId   = "1a3af483-89e5-4820-a579-13ed6d90b0cc"
	DummyErrorInstanceId = "ba3aaedd-458f-4582-b763-12c3ae7b27ee"
	QAInstanceId         = "3c2cd2ea-39bf-46d8-b2a2-c52af81d9072"
	encryptedPwd         = "49e5bf3f6a45a75c972c68b39d640e53f050a6a0b4125ff9"
)

var (
	ValidInstanceInput = dto.DatabaseInstanceInputDTO{
		Name:                 "PostgreSQL - Local",
		Host:                 "localhost",
		Port:                 "5432",
		HostConnection:       "host.conn.ip",
		PortConnection:       "5433",
		AdminUser:            "admin",
		AdminPassword:        "pwd",
		EcosystemID:          EcosystemId,
		DatabaseTechnologyID: TechnologyId,
		Note:                 "note",
	}
)

func BuildTestInstance() *entity.DatabaseInstance {
	return &entity.DatabaseInstance{
		ID:                   uuid.MustParse(QAInstanceId),
		Name:                 "Test Local",
		HostConnection:       &entity.HostConnectionInfo{Host: "localhost", Port: "5432"},
		EcosystemID:          EcosystemId,
		DatabaseTechnologyID: TechnologyId,
		Note:                 "note",
		CreatedAt:            time.Now(),
		CreatedByUserID:      UserID,
		Enabled:              true,
	}
}

func BuildInstancesList() []*dto.DatabaseInstanceOutputDTO {
	dbInstanceDto := BuildQAInstanceDTO()
	dbInstanceDto2 := BuildAzInstanceDTO()
	dbInstanceDto3 := BuildDummyErrorInstance()
	dbInstances := []*dto.DatabaseInstanceOutputDTO{dbInstanceDto, dbInstanceDto2, dbInstanceDto3}
	return dbInstances
}

func BuildQAInstanceDTO() *dto.DatabaseInstanceOutputDTO {
	return &dto.DatabaseInstanceOutputDTO{
		ID:                     QAInstanceId,
		EcosystemName:          "QA",
		Name:                   connector.DummyTest + " - QA",
		DatabaseTechnologyName: connector.DummyTest,
		HostConnection:         "localhost",
		PortConnection:         "5432",
		AdminUser:              "dummy",
		AdminPassword:          encryptedPwd,
		CreatedByUser:          "Foo Bar",
		Enabled:                false,
	}
}

func BuildQAInstanceEnabled() *dto.DatabaseInstanceOutputDTO {
	i := BuildQAInstanceDTO()
	i.Enabled = true
	return i
}

func BuildAzInstanceDTO() *dto.DatabaseInstanceOutputDTO {
	return &dto.DatabaseInstanceOutputDTO{
		ID:                        DatabaseInstanceId,
		EcosystemName:             "Azure",
		Name:                      connector.DummyTest + " - Azure",
		DatabaseTechnologyName:    connector.DummyTest,
		DatabaseTechnologyVersion: "2",
		HostConnection:            "10.1.1.1",
		PortConnection:            "5432",
		AdminUser:                 "dummy",
		AdminPassword:             encryptedPwd,
		CreatedByUser:             "John Doe",
		Enabled:                   true,
		RolesCreated:              true,
	}
}

func BuildDummyErrorInstance() *dto.DatabaseInstanceOutputDTO {
	return &dto.DatabaseInstanceOutputDTO{
		ID:                        DummyErrorInstanceId,
		EcosystemName:             "QA",
		Name:                      connector.InstanceDummyTestError,
		DatabaseTechnologyName:    connector.DummyTest,
		DatabaseTechnologyVersion: "1",
		Host:                      "localhost",
		Port:                      "5435",
		AdminUser:                 "dummy",
		AdminPassword:             encryptedPwd,
		CreatedByUser:             "John Doe",
		Enabled:                   true,
		RolesCreated:              true,
	}
}

func BuildConnectorNotImplementedInstance() *dto.DatabaseInstanceOutputDTO {
	return &dto.DatabaseInstanceOutputDTO{
		ID:                        "57200738-9b52-4c31-945b-fb1603df4f37",
		Name:                      "MySQL - AWS",
		EcosystemName:             "AWS",
		DatabaseTechnologyName:    "MySQL",
		DatabaseTechnologyVersion: "5",
		AdminUser:                 "admin",
		AdminPassword:             encryptedPwd,
		Enabled:                   true,
		RolesCreated:              true,
	}
}

func BuildFullDataInstanceExample() *dto.DatabaseInstanceOutputDTO {
	systemTime := time.Now()
	return &dto.DatabaseInstanceOutputDTO{
		ID:                        uuid.New().String(),
		Name:                      ValidInstanceInput.Name,
		Host:                      ValidInstanceInput.Host,
		HostConnection:            ValidInstanceInput.HostConnection,
		Port:                      ValidInstanceInput.Port,
		PortConnection:            ValidInstanceInput.PortConnection,
		AdminUser:                 ValidInstanceInput.AdminUser,
		AdminPassword:             encryptedPwd,
		EcosystemID:               ValidInstanceInput.EcosystemID,
		EcosystemName:             "QA AWS",
		DatabaseTechnologyID:      ValidInstanceInput.DatabaseTechnologyID,
		DatabaseTechnologyName:    "PostgreSQL",
		DatabaseTechnologyVersion: "13.3",
		Note:                      ValidInstanceInput.Note,
		CreatedByUserID:           UserID,
		CreatedByUser:             "User Test",
		CreatedAt:                 systemTime,
		UpdatedAt:                 &systemTime,
		Enabled:                   true,
	}
}

type DatabaseInstanceStorageMock struct {
	mock.Mock
}

func (m *DatabaseInstanceStorageMock) Exists(host, port string) (bool, error) {
	args := m.Called(host, port)
	return args.Bool(0), args.Error(1)
}

func (m *DatabaseInstanceStorageMock) Save(databaseInstance *entity.DatabaseInstance) error {
	args := m.Called(databaseInstance)
	return args.Error(0)
}

func (m *DatabaseInstanceStorageMock) UpdateWithHostInfo(databaseInstance *entity.DatabaseInstance) error {
	args := m.Called(databaseInstance)
	return args.Error(0)
}

func (m *DatabaseInstanceStorageMock) Update(databaseInstance *entity.DatabaseInstance) error {
	args := m.Called(databaseInstance)
	return args.Error(0)
}

func (m *DatabaseInstanceStorageMock) FindByID(id string) (*entity.DatabaseInstance, error) {
	args := m.Called(id)
	return args.Get(0).(*entity.DatabaseInstance), args.Error(1)
}

func (m *DatabaseInstanceStorageMock) FindDTOByID(id string) (*dto.DatabaseInstanceOutputDTO, error) {
	args := m.Called(id)
	return args.Get(0).(*dto.DatabaseInstanceOutputDTO), args.Error(1)
}

func (m *DatabaseInstanceStorageMock) FindAllDTOs(ecosystemId, technologyId string, ids []string) ([]*dto.DatabaseInstanceOutputDTO, error) {
	args := m.Called(ecosystemId, technologyId, ids)
	return args.Get(0).([]*dto.DatabaseInstanceOutputDTO), args.Error(1)
}

func (m *DatabaseInstanceStorageMock) FindAllDTOsEnabled(ecosystemId, technologyId string) ([]*dto.DatabaseInstanceOutputDTO, error) {
	args := m.Called(ecosystemId, technologyId)
	return args.Get(0).([]*dto.DatabaseInstanceOutputDTO), args.Error(1)
}

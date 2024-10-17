package mocks

import (
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	"github.com/zgsolucoes/zg-data-guard/internal/entity"
)

const (
	DefaultPage  = 0
	DefaultLimit = 50
)

type AccessPermissionStorageMock struct {
	mock.Mock
}

func (a *AccessPermissionStorageMock) Save(accessPermission *entity.AccessPermission) error {
	args := a.Called(accessPermission)
	return args.Error(0)
}

func (a *AccessPermissionStorageMock) Exists(databaseId, databaseUserId string) (bool, error) {
	args := a.Called(databaseId, databaseUserId)
	return args.Bool(0), args.Error(1)
}

func (a *AccessPermissionStorageMock) FindAllDTOs(databaseId, databaseUserId, databaseInstanceId string) ([]*dto.AccessPermissionOutputDTO, error) {
	args := a.Called(databaseId, databaseUserId, databaseInstanceId)
	return args.Get(0).([]*dto.AccessPermissionOutputDTO), args.Error(1)
}

func (a *AccessPermissionStorageMock) SaveLog(log *entity.AccessPermissionLog) error {
	args := a.Called(log)
	return args.Error(0)
}

func (a *AccessPermissionStorageMock) FindAllAccessibleInstancesIDsByUser(userID string) ([]string, error) {
	args := a.Called(userID)
	return args.Get(0).([]string), args.Error(1)
}

func (a *AccessPermissionStorageMock) FindAllLogsDTOs(page, limit int) ([]*dto.AccessPermissionLogOutputDTO, error) {
	args := a.Called(page, limit)
	return args.Get(0).([]*dto.AccessPermissionLogOutputDTO), args.Error(1)
}

func (a *AccessPermissionStorageMock) DeleteAllByUserAndInstance(databaseUserID, instanceID string) error {
	args := a.Called(databaseUserID, instanceID)
	return args.Error(0)
}

func (a *AccessPermissionStorageMock) DeleteAllByInstance(instanceID string) error {
	args := a.Called(instanceID)
	return args.Error(0)
}

func (a *AccessPermissionStorageMock) CheckIfUserHasAccessPermission(databaseUserID string) (bool, error) {
	args := a.Called(databaseUserID)
	return args.Bool(0), args.Error(1)
}

func (a *AccessPermissionStorageMock) LogCount() (int, error) {
	args := a.Called()
	return args.Int(0), args.Error(1)
}

func BuildAccessPermissionsDTOList() []*dto.AccessPermissionOutputDTO {
	return []*dto.AccessPermissionOutputDTO{
		{
			ID:             "1",
			DatabaseID:     "1",
			DatabaseUserID: "1",
			GrantedAt:      time.Now(),
		},
		{
			ID:             "2",
			DatabaseID:     "1",
			DatabaseUserID: "2",
			GrantedAt:      time.Now(),
		},
		{
			ID:             "3",
			DatabaseID:     "1",
			DatabaseUserID: "3",
			GrantedAt:      time.Now(),
		},
	}
}

func BuildAccessPermissionLogDTOList() []*dto.AccessPermissionLogOutputDTO {
	strDatabaseUserID := "1"
	strDatabaseUserID2 := "2"
	strDatabaseID := "2"
	return []*dto.AccessPermissionLogOutputDTO{
		{
			ID:                 "1",
			DatabaseInstanceID: "1",
			DatabaseUserID:     &strDatabaseUserID,
			DatabaseID:         &strDatabaseID,
			Message:            "Log message 1",
			Success:            true,
			Date:               time.Now(),
			OperationUserID:    "1",
		},
		{
			ID:                 "2",
			DatabaseInstanceID: "1",
			DatabaseUserID:     &strDatabaseUserID2,
			DatabaseID:         &strDatabaseID,
			Message:            "Log message 2",
			Success:            true,
			Date:               time.Now(),
			OperationUserID:    "1",
		},
		{
			ID:                 "3",
			DatabaseInstanceID: "2",
			DatabaseUserID:     nil,
			DatabaseID:         nil,
			Message:            "Log message 3",
			Success:            false,
			Date:               time.Now(),
			OperationUserID:    "2",
		},
	}
}

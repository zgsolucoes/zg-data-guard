package accesspermission

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	"github.com/zgsolucoes/zg-data-guard/testdata/mocks"
)

func TestGivenAnErrorInDb_WhenExecuteListAccessPermissionLogs_ThenShouldReturnError(t *testing.T) {
	accessStorage := new(mocks.AccessPermissionStorageMock)
	accessStorage.On("FindAllLogsDTOs", mocks.DefaultPage, mocks.DefaultLimit).Return([]*dto.AccessPermissionLogOutputDTO{}, sql.ErrConnDone).Once()

	uc := NewListAccessPermissionLogsUseCase(accessStorage)
	accessObtained, totalCount, err := uc.Execute(mocks.DefaultPage, mocks.DefaultLimit)

	assert.Error(t, err, "error expected when some error in db")
	assert.EqualError(t, err, "error fetching access permission logs! Cause: sql: connection is already closed")
	assert.Nil(t, accessObtained)
	assert.Equal(t, totalCount, 0)
	assert.Equal(t, len(accessObtained), 0, "0 access permission logs expected")
	accessStorage.AssertNumberOfCalls(t, "FindAllLogsDTOs", 1)
}

func TestGivenAnErrorInDb_WhenExecuteLogCount_ThenShouldReturnError(t *testing.T) {
	accessStorage := new(mocks.AccessPermissionStorageMock)
	accessStorage.On("FindAllLogsDTOs", mocks.DefaultPage, mocks.DefaultLimit).Return(mocks.BuildAccessPermissionLogDTOList(), nil).Once()
	accessStorage.On("LogCount").Return(0, sql.ErrConnDone).Once()

	uc := NewListAccessPermissionLogsUseCase(accessStorage)
	accessObtained, totalCount, err := uc.Execute(mocks.DefaultPage, mocks.DefaultLimit)

	assert.Error(t, err, "error expected when some error in db")
	assert.EqualError(t, err, "error fetching access permission logs count! Cause: sql: connection is already closed")
	assert.Nil(t, accessObtained)
	assert.Equal(t, totalCount, 0)
	assert.Equal(t, len(accessObtained), 0, "0 access permission logs expected")
	accessStorage.AssertNumberOfCalls(t, "FindAllLogsDTOs", 1)
	accessStorage.AssertNumberOfCalls(t, "LogCount", 1)
}

func TestGivenSomeLogs_WhenExecuteListAccessPermissionLogs_ThenShouldListAllPermissionLogs(t *testing.T) {
	logList := mocks.BuildAccessPermissionLogDTOList()
	accessStorage := new(mocks.AccessPermissionStorageMock)
	accessStorage.On("FindAllLogsDTOs", mocks.DefaultPage, mocks.DefaultLimit).Return(logList, nil).Once()
	accessStorage.On("LogCount").Return(len(logList), nil).Once()

	uc := NewListAccessPermissionLogsUseCase(accessStorage)
	permissionsObtained, totalCount, err := uc.Execute(mocks.DefaultPage, mocks.DefaultLimit)

	assert.NoError(t, err, "no error expected")
	assert.Equal(t, len(permissionsObtained), len(logList), "3 access permission logs expected")
	assert.Equal(t, totalCount, len(logList), "3 access permission logs expected")
	accessStorage.AssertNumberOfCalls(t, "FindAllLogsDTOs", 1)
	accessStorage.AssertNumberOfCalls(t, "LogCount", 1)
}

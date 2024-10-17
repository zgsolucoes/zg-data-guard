package instance

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/zgsolucoes/zg-data-guard/internal/entity"
	"github.com/zgsolucoes/zg-data-guard/internal/usecase/common"
	"github.com/zgsolucoes/zg-data-guard/testdata"
	"github.com/zgsolucoes/zg-data-guard/testdata/mocks"
)

func TestGivenAnErrorInDbWhileCheckHost_WhenExecuteChangeStatus_ThenShouldReturnError(t *testing.T) {
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	dbInstanceStorage.On("FindByID", mocks.DbUserID).Return(&entity.DatabaseInstance{}, sql.ErrConnDone).Once()

	uc := NewChangeStatusDatabaseInstanceUseCase(dbInstanceStorage, nil, nil)
	dbUserOutput, err := uc.Execute(mocks.DbUserID, true, mocks.UserID)

	assert.Error(t, err, "error expected when some error in db")
	assert.EqualError(t, err, sql.ErrConnDone.Error())
	assert.Nil(t, dbUserOutput)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindByID", 1)
}

func TestGivenANonexistentId_WhenExecuteChangeStatus_ThenShouldReturnError(t *testing.T) {
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	dbInstanceStorage.On("FindByID", mocks.DbUserID).Return(&entity.DatabaseInstance{}, sql.ErrNoRows).Once()

	uc := NewChangeStatusDatabaseInstanceUseCase(dbInstanceStorage, nil, nil)
	dbUserOutput, err := uc.Execute(mocks.DbUserID, true, mocks.UserID)

	assert.Error(t, err)
	assert.EqualError(t, err, ErrDatabaseInstanceNotFound.Error())
	assert.Nil(t, dbUserOutput)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindByID", 1)
}

func TestGivenErrorWhileUpdatingInstance_WhenExecuteChangeStatus_ThenShouldReturnError(t *testing.T) {
	dbInstance := mocks.BuildTestInstance()
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	dbInstanceStorage.On("FindByID", dbInstance.ID.String()).Return(dbInstance, nil).Once()
	dbInstanceStorage.On("Update", dbInstance).Return(sql.ErrConnDone).Once()

	uc := NewChangeStatusDatabaseInstanceUseCase(dbInstanceStorage, nil, nil)
	dbUserOutput, err := uc.Execute(dbInstance.ID.String(), true, mocks.UserID)

	assert.Error(t, err)
	assert.EqualError(t, err, fmt.Errorf("error when updating database instance %s to new status '%t'. Cause: %w", dbInstance.ID.String(), true, sql.ErrConnDone).Error())
	assert.Nil(t, dbUserOutput)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindByID", 1)
	dbInstanceStorage.AssertNumberOfCalls(t, "Update", 1)
}

func TestGivenErrorWhileRevokingAccess_WhenExecuteDisable_ThenShouldReturnError(t *testing.T) {
	dbInstance := mocks.BuildTestInstance()
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	dbInstanceStorage.On("FindByID", dbInstance.ID.String()).Return(dbInstance, nil).Once()
	accessPermissionStorage := new(mocks.AccessPermissionStorageMock)
	accessPermissionStorage.On("DeleteAllByInstance", dbInstance.ID.String()).Return(sql.ErrConnDone).Once()

	uc := NewChangeStatusDatabaseInstanceUseCase(dbInstanceStorage, nil, accessPermissionStorage)
	dbUserOutput, err := uc.Execute(dbInstance.ID.String(), false, mocks.UserID)

	assert.Error(t, err)
	assert.EqualError(t, err, fmt.Errorf("error when revoking access from database instance '%s'. Cause: %w", dbInstance.Name, sql.ErrConnDone).Error())
	assert.Nil(t, dbUserOutput)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindByID", 1)
	accessPermissionStorage.AssertNumberOfCalls(t, "DeleteAllByInstance", 1)
}

func TestGivenErrorWhileDeactivatingDatabases_WhenExecuteDisable_ThenShouldReturnError(t *testing.T) {
	dbInstance := mocks.BuildTestInstance()
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	dbInstanceStorage.On("FindByID", dbInstance.ID.String()).Return(dbInstance, nil).Once()
	accessPermissionStorage := new(mocks.AccessPermissionStorageMock)
	accessPermissionStorage.On("DeleteAllByInstance", dbInstance.ID.String()).Return(nil).Once()
	databaseStorage := new(mocks.DatabaseStorageMock)
	databaseStorage.On("DeactivateAllByInstance", dbInstance.ID.String()).Return(sql.ErrConnDone).Once()

	uc := NewChangeStatusDatabaseInstanceUseCase(dbInstanceStorage, databaseStorage, accessPermissionStorage)
	dbUserOutput, err := uc.Execute(dbInstance.ID.String(), false, mocks.UserID)

	assert.Error(t, err)
	assert.EqualError(t, err, fmt.Errorf("error when deactivating databases from database instance '%s'. Cause: %w", dbInstance.Name, sql.ErrConnDone).Error())
	assert.Nil(t, dbUserOutput)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindByID", 1)
	accessPermissionStorage.AssertNumberOfCalls(t, "DeleteAllByInstance", 1)
	databaseStorage.AssertNumberOfCalls(t, "DeactivateAllByInstance", 1)
}

func TestGivenErrorWhileRegisteringLog_WhenExecuteEnable_ThenShouldReturnError(t *testing.T) {
	dbInstance := mocks.BuildTestInstance()
	dbInstance.Disable()
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	dbInstanceStorage.On("FindByID", dbInstance.ID.String()).Return(dbInstance, nil).Once()
	dbInstanceStorage.On("Update", dbInstance).Return(nil).Once()
	accessPermissionStorage := new(mocks.AccessPermissionStorageMock)
	accessPermissionStorage.On("SaveLog", mock.Anything).Return(sql.ErrConnDone).Once()

	uc := NewChangeStatusDatabaseInstanceUseCase(dbInstanceStorage, nil, accessPermissionStorage)
	dbUserOutput, err := uc.Execute(dbInstance.ID.String(), true, mocks.UserID)

	assert.Error(t, err)
	assert.EqualError(t, err, "error when registering log for database instance 'Test Local'. Cause: error when saving log for instance 'Test Local'. Cause: sql: connection is already closed")
	assert.Nil(t, dbUserOutput)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindByID", 1)
	dbInstanceStorage.AssertNumberOfCalls(t, "Update", 1)
	accessPermissionStorage.AssertNumberOfCalls(t, "SaveLog", 1)
}

func TestRegisterLog_ErrorCreatingLog(t *testing.T) {
	accessPermissionStorage := new(mocks.AccessPermissionStorageMock)
	uc := NewChangeStatusDatabaseInstanceUseCase(nil, nil, accessPermissionStorage)

	err := uc.registerLog("", "TestDB", mocks.UserID, true)

	assert.Error(t, err)
	assert.EqualError(t, err, "error when creating log for instance 'TestDB'. Cause: database instance id not informed")
	accessPermissionStorage.AssertNotCalled(t, "SaveLog")
}

func TestGivenEnableInput_WhenExecuteChangeStatus_ThenShouldEnableInstance(t *testing.T) {
	dbInstance := mocks.BuildTestInstance()
	dbInstance.Disable()
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	dbInstanceStorage.On("FindByID", mocks.DbUserID).Return(dbInstance, nil).Once()
	dbInstanceStorage.On("Update", mock.Anything).Return(nil).Once()
	accessPermissionStorage := new(mocks.AccessPermissionStorageMock)
	log, _ := entity.NewAccessPermissionLog(dbInstance.ID.String(), "", "", fmt.Sprintf(common.InstanceEnabledSuccessMsg, dbInstance.Name), mocks.UserID, true)
	accessPermissionStorage.On("SaveLog", testdata.CompareLogs(log)).Return(nil).Once()

	uc := NewChangeStatusDatabaseInstanceUseCase(dbInstanceStorage, nil, accessPermissionStorage)
	output, err := uc.Execute(mocks.DbUserID, true, mocks.UserID)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, dbInstance.ID.String(), output.ID)
	assert.True(t, dbInstance.Enabled)
	assert.True(t, output.Enabled)
	assert.NotEmpty(t, output.UpdatedAt)
	assert.Empty(t, output.DisabledAt)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindByID", 1)
	dbInstanceStorage.AssertNumberOfCalls(t, "Update", 1)
	accessPermissionStorage.AssertNumberOfCalls(t, "SaveLog", 1)
}

func TestGivenDisableInput_WhenExecuteChangeStatus_ThenShouldDisableInstance(t *testing.T) {
	dbInstance := mocks.BuildTestInstance()
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	dbInstanceStorage.On("FindByID", mocks.DbUserID).Return(dbInstance, nil).Once()
	dbInstanceStorage.On("Update", dbInstance).Return(nil).Once()
	accessPermissionStorage := new(mocks.AccessPermissionStorageMock)
	accessPermissionStorage.On("DeleteAllByInstance", dbInstance.ID.String()).Return(nil).Once()
	log, _ := entity.NewAccessPermissionLog(dbInstance.ID.String(), "", "", fmt.Sprintf(common.InstanceDisabledSuccessMsg, dbInstance.Name), mocks.UserID, true)
	accessPermissionStorage.On("SaveLog", testdata.CompareLogs(log)).Return(nil).Once()
	databaseStorage := new(mocks.DatabaseStorageMock)
	databaseStorage.On("DeactivateAllByInstance", dbInstance.ID.String()).Return(nil).Once()

	uc := NewChangeStatusDatabaseInstanceUseCase(dbInstanceStorage, databaseStorage, accessPermissionStorage)
	output, err := uc.Execute(mocks.DbUserID, false, mocks.UserID)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.False(t, output.Enabled)
	assert.NotEmpty(t, output.DisabledAt)
	assert.NotEmpty(t, output.UpdatedAt)
	assert.False(t, dbInstance.Enabled)
	assert.NotEmpty(t, dbInstance.DisabledAt)
	assert.NotEmpty(t, dbInstance.UpdatedAt)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindByID", 1)
	dbInstanceStorage.AssertNumberOfCalls(t, "Update", 1)
	accessPermissionStorage.AssertNumberOfCalls(t, "DeleteAllByInstance", 1)
	accessPermissionStorage.AssertNumberOfCalls(t, "SaveLog", 1)
	databaseStorage.AssertNumberOfCalls(t, "DeactivateAllByInstance", 1)
}

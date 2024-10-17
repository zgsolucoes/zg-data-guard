package accesspermission

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	"github.com/zgsolucoes/zg-data-guard/internal/entity"
	"github.com/zgsolucoes/zg-data-guard/internal/usecase/common"
	"github.com/zgsolucoes/zg-data-guard/testdata"
	"github.com/zgsolucoes/zg-data-guard/testdata/mocks"
)

func TestGivenAnErrorInDbWhenFetchingDBUser_WhenExecuteRevokeAccess_ThenShouldReturnError(t *testing.T) {
	dbUserStorage := new(mocks.DatabaseUserStorageMock)
	dbUserStorage.On("FindByID", "").Return(&entity.DatabaseUser{}, sql.ErrConnDone).Once()

	uc := NewRevokeAccessPermissionUseCase(nil, nil, dbUserStorage)
	output, err := uc.Execute(dto.RevokeAccessInputDTO{DatabaseUserID: "", DatabaseInstancesIDs: []string{}}, mocks.UserID)

	assert.Error(t, err, "error expected when some error in db")
	assert.EqualError(t, err, sql.ErrConnDone.Error())
	assert.Nil(t, output)
	dbUserStorage.AssertNumberOfCalls(t, "FindByID", 1)
}

func TestGivenAnNotExistingIDWhenFetchingDBUser_WhenExecuteRevokeAccess_ThenShouldReturnError(t *testing.T) {
	dbUserStorage := new(mocks.DatabaseUserStorageMock)
	dbUserStorage.On("FindByID", "").Return(&entity.DatabaseUser{}, sql.ErrNoRows).Once()

	uc := NewRevokeAccessPermissionUseCase(nil, nil, dbUserStorage)
	output, err := uc.Execute(dto.RevokeAccessInputDTO{DatabaseUserID: "", DatabaseInstancesIDs: []string{}}, mocks.UserID)

	assert.Error(t, err, "error expected when database user not found in db")
	assert.EqualError(t, err, common.ErrDatabaseUserNotFound.Error())
	assert.Nil(t, output)
	dbUserStorage.AssertNumberOfCalls(t, "FindByID", 1)
}

func TestGivenAnErrorInDbWhenFetchingAccessibleInstances_WhenExecuteRevokeAccess_ThenShouldReturnError(t *testing.T) {
	dbUserStorage := new(mocks.DatabaseUserStorageMock)
	dbUserStorage.On("FindByID", "").Return(&entity.DatabaseUser{}, nil).Once()
	accessPermissionStorage := new(mocks.AccessPermissionStorageMock)
	accessPermissionStorage.On("FindAllAccessibleInstancesIDsByUser", "").Return([]string{}, sql.ErrConnDone).Once()

	uc := NewRevokeAccessPermissionUseCase(accessPermissionStorage, nil, dbUserStorage)
	output, err := uc.Execute(dto.RevokeAccessInputDTO{DatabaseInstancesIDs: []string{}}, mocks.UserID)

	assert.Error(t, err, "error expected when some error in db")
	assert.EqualError(t, err, sql.ErrConnDone.Error())
	assert.Nil(t, output)
	dbUserStorage.AssertNumberOfCalls(t, "FindByID", 1)
	accessPermissionStorage.AssertNumberOfCalls(t, "FindAllAccessibleInstancesIDsByUser", 1)
}

func TestGivenAnUserWithNoAccessibleInstances_WhenExecuteRevokeAccess_ThenShouldReturnError(t *testing.T) {
	dbUserStorage := new(mocks.DatabaseUserStorageMock)
	dbUserStorage.On("FindByID", "").Return(&entity.DatabaseUser{}, nil).Once()
	accessPermissionStorage := new(mocks.AccessPermissionStorageMock)
	accessPermissionStorage.On("FindAllAccessibleInstancesIDsByUser", "").Return([]string{}, nil).Once()

	uc := NewRevokeAccessPermissionUseCase(accessPermissionStorage, nil, dbUserStorage)
	output, err := uc.Execute(dto.RevokeAccessInputDTO{DatabaseInstancesIDs: []string{}}, mocks.UserID)

	assert.Error(t, err, "error expected when the database user has no accessible instances")
	assert.EqualError(t, err, common.ErrNoAccessibleInstancesFound.Error())
	assert.Nil(t, output)
	dbUserStorage.AssertNumberOfCalls(t, "FindByID", 1)
	accessPermissionStorage.AssertNumberOfCalls(t, "FindAllAccessibleInstancesIDsByUser", 1)
}

func TestGivenNonAccessibleInstancesInput_WhenExecuteRevokeAccess_ThenShouldReturnError(t *testing.T) {
	instance := mocks.BuildQAInstanceDTO()
	dbUser := mocks.BuildDbUserJohn()
	dbUserStorage := new(mocks.DatabaseUserStorageMock)
	dbUserStorage.On("FindByID", dbUser.ID.String()).Return(dbUser, nil).Once()
	accessPermissionStorage := new(mocks.AccessPermissionStorageMock)
	accessPermissionStorage.On("FindAllAccessibleInstancesIDsByUser", dbUser.ID.String()).Return([]string{instance.ID}, nil).Once()

	uc := NewRevokeAccessPermissionUseCase(accessPermissionStorage, nil, dbUserStorage)
	output, err := uc.Execute(dto.RevokeAccessInputDTO{
		DatabaseInstancesIDs: []string{"57200738-9b52-4c31-945b-fb1603df4f37", mocks.DatabaseInstanceId, mocks.DummyErrorInstanceId},
		DatabaseUserID:       dbUser.ID.String()},
		mocks.UserID,
	)

	assert.Error(t, err, "error expected when the selected instances in input are not accessible by the database user")
	assert.EqualError(t, err, common.ErrNoAccessibleInstancesFound.Error())
	assert.Nil(t, output)
	dbUserStorage.AssertNumberOfCalls(t, "FindByID", 1)
	accessPermissionStorage.AssertNumberOfCalls(t, "FindAllAccessibleInstancesIDsByUser", 1)
}

func TestGivenAnErrorInDbWhenFetchingInstances_WhenExecuteRevokeAccess_ThenShouldReturnOutputError(t *testing.T) {
	instance := mocks.BuildQAInstanceDTO()
	dbUser := mocks.BuildDbUserJohn()
	dbUserStorage := new(mocks.DatabaseUserStorageMock)
	dbUserStorage.On("FindByID", dbUser.ID.String()).Return(dbUser, nil).Once()
	accessPermissionStorage := new(mocks.AccessPermissionStorageMock)
	accessPermissionStorage.On("FindAllAccessibleInstancesIDsByUser", dbUser.ID.String()).Return([]string{instance.ID}, nil).Once()
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	dbInstanceStorage.On("FindAllDTOs", "", "", []string{instance.ID}).Return([]*dto.DatabaseInstanceOutputDTO{}, sql.ErrConnDone).Once()

	uc := NewRevokeAccessPermissionUseCase(accessPermissionStorage, dbInstanceStorage, dbUserStorage)
	output, err := uc.Execute(dto.RevokeAccessInputDTO{DatabaseInstancesIDs: []string{instance.ID}, DatabaseUserID: dbUser.ID.String()}, mocks.UserID)

	assert.Error(t, err, "error expected when some error in db")
	assert.EqualError(t, err, sql.ErrConnDone.Error())
	assert.Nil(t, output)
	dbUserStorage.AssertNumberOfCalls(t, "FindByID", 1)
	accessPermissionStorage.AssertNumberOfCalls(t, "FindAllAccessibleInstancesIDsByUser", 1)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindAllDTOs", 1)
}

func TestGivenAnErrorInDbWhenDeletingAccess_WhenExecuteRevokeAccess_ThenShouldReturnOutputError(t *testing.T) {
	instance := mocks.BuildQAInstanceDTO()
	dbUser := mocks.BuildDbUserJohn()
	dbUserID := getDBUserID(dbUser)
	dbUserStorage := new(mocks.DatabaseUserStorageMock)
	dbUserStorage.On("FindByID", dbUserID).Return(dbUser, nil).Once()
	accessPermissionStorage := new(mocks.AccessPermissionStorageMock)
	accessPermissionStorage.On("FindAllAccessibleInstancesIDsByUser", dbUserID).Return([]string{instance.ID}, nil).Once()
	accessPermissionStorage.On("DeleteAllByUserAndInstance", dbUserID, instance.ID).Return(sql.ErrConnDone).Once()
	expectedLogMsg := fmt.Sprintf(ErrDeletingAccessOfUserMsg, dbUser.Username, instance.Name, sql.ErrConnDone.Error())
	expectedLog, _ := entity.NewAccessPermissionLog(instance.ID, dbUserID, "", expectedLogMsg, mocks.UserID, false)
	accessPermissionStorage.On("SaveLog", testdata.CompareLogs(expectedLog)).Return(nil).Once()
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	dbInstanceStorage.On("FindAllDTOs", "", "", []string{instance.ID}).Return([]*dto.DatabaseInstanceOutputDTO{instance}, nil).Once()

	uc := NewRevokeAccessPermissionUseCase(accessPermissionStorage, dbInstanceStorage, dbUserStorage)
	output, err := uc.Execute(dto.RevokeAccessInputDTO{DatabaseInstancesIDs: []string{instance.ID}, DatabaseUserID: dbUserID}, mocks.UserID)

	assert.NoError(t, err, "the error should be in the output, not in the process and has to be logged")
	assert.NotNil(t, output)
	assert.True(t, output.HasErrors)
	assert.Equal(t, SomeErrorsDuringProcessMsg, output.Message)
	dbUserStorage.AssertNumberOfCalls(t, "FindByID", 1)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindAllDTOs", 1)
	accessPermissionStorage.AssertNumberOfCalls(t, "FindAllAccessibleInstancesIDsByUser", 1)
	accessPermissionStorage.AssertNumberOfCalls(t, "DeleteAllByUserAndInstance", 1)
	accessPermissionStorage.AssertNumberOfCalls(t, "SaveLog", 1)
}

func TestGivenAnErrorWhenCreatingLog_WhenExecuteRevokeAccess_ThenShouldReturnOutputError(t *testing.T) {
	instance := mocks.BuildQAInstanceDTO()
	dbUser := mocks.BuildDbUserJohn()
	dbUserID := getDBUserID(dbUser)
	dbUserStorage := new(mocks.DatabaseUserStorageMock)
	dbUserStorage.On("FindByID", dbUserID).Return(dbUser, nil).Once()
	accessPermissionStorage := new(mocks.AccessPermissionStorageMock)
	accessPermissionStorage.On("FindAllAccessibleInstancesIDsByUser", dbUserID).Return([]string{instance.ID}, nil).Once()
	accessPermissionStorage.On("DeleteAllByUserAndInstance", dbUserID, instance.ID).Return(nil).Once()
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	dbInstanceStorage.On("FindAllDTOs", "", "", []string{instance.ID}).Return([]*dto.DatabaseInstanceOutputDTO{instance}, nil).Once()

	uc := NewRevokeAccessPermissionUseCase(accessPermissionStorage, dbInstanceStorage, dbUserStorage)
	output, err := uc.Execute(dto.RevokeAccessInputDTO{DatabaseInstancesIDs: []string{instance.ID}, DatabaseUserID: dbUserID}, "")

	assert.NoError(t, err, "the error should be in the output, not in the process and has to be logged. The user ID is not informed")
	assert.NotNil(t, output)
	assert.True(t, output.HasErrors)
	assert.Equal(t, SomeErrorsDuringProcessMsg, output.Message)
	dbUserStorage.AssertNumberOfCalls(t, "FindByID", 1)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindAllDTOs", 1)
	accessPermissionStorage.AssertNumberOfCalls(t, "FindAllAccessibleInstancesIDsByUser", 1)
	accessPermissionStorage.AssertNumberOfCalls(t, "DeleteAllByUserAndInstance", 1)
}

func TestGivenAnErrorInDbWhenSavingLog_WhenExecuteRevokeAccess_ThenShouldReturnOutputError(t *testing.T) {
	instance := mocks.BuildQAInstanceDTO()
	dbUser := mocks.BuildDbUserJohn()
	dbUserID := getDBUserID(dbUser)
	dbUserStorage := new(mocks.DatabaseUserStorageMock)
	dbUserStorage.On("FindByID", dbUserID).Return(dbUser, nil).Once()
	accessPermissionStorage := new(mocks.AccessPermissionStorageMock)
	accessPermissionStorage.On("FindAllAccessibleInstancesIDsByUser", dbUserID).Return([]string{instance.ID}, nil).Once()
	accessPermissionStorage.On("DeleteAllByUserAndInstance", dbUserID, instance.ID).Return(nil).Once()
	expectedLogMsg := fmt.Sprintf(UserAccessRevokedAndExcludedMsg, dbUser.Username, instance.Name)
	expectedLog, _ := entity.NewAccessPermissionLog(instance.ID, dbUserID, "", expectedLogMsg, mocks.UserID, true)
	accessPermissionStorage.On("SaveLog", testdata.CompareLogs(expectedLog)).Return(sql.ErrConnDone).Once()
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	dbInstanceStorage.On("FindAllDTOs", "", "", []string{instance.ID}).Return([]*dto.DatabaseInstanceOutputDTO{instance}, nil).Once()

	uc := NewRevokeAccessPermissionUseCase(accessPermissionStorage, dbInstanceStorage, dbUserStorage)
	output, err := uc.Execute(dto.RevokeAccessInputDTO{DatabaseInstancesIDs: []string{instance.ID}, DatabaseUserID: dbUserID}, mocks.UserID)

	assert.NoError(t, err, "the error should be in the output, not in the process and has to be logged")
	assert.NotNil(t, output)
	assert.True(t, output.HasErrors)
	assert.Equal(t, SomeErrorsDuringProcessMsg, output.Message)
	dbUserStorage.AssertNumberOfCalls(t, "FindByID", 1)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindAllDTOs", 1)
	accessPermissionStorage.AssertNumberOfCalls(t, "FindAllAccessibleInstancesIDsByUser", 1)
	accessPermissionStorage.AssertNumberOfCalls(t, "DeleteAllByUserAndInstance", 1)
	accessPermissionStorage.AssertNumberOfCalls(t, "SaveLog", 1)
}

func TestGivenAnErrorWhenCreatingConnector_WhenExecuteRevokeAccess_ThenShouldReturnOutputError(t *testing.T) {
	dbUser := mocks.BuildDbUserJohn()
	instance := mocks.BuildConnectorNotImplementedInstance()
	expectedLogMsg := fmt.Sprintf(ErrCreatingConnectorMsg, instance.Name, "the database technology 'mysql' don't have a connector implemented")
	runRevokeLoggingSingleError(t, dbUser, instance, expectedLogMsg)
}

func TestGivenAnErrorWhenRevokingPermissionsInInstance_WhenExecuteRevokeAccess_ThenShouldReturnOutputError(t *testing.T) {
	dbUser := mocks.BuildDbUserJohn()
	instance := mocks.BuildDummyErrorInstance()
	expectedLogMsg := fmt.Sprintf(ErrRevokeAndDropUserFailedMsg, dbUser.Username, instance.Name, "error revoking permissions and removing user: Instance(instance-dummy-test-error) - User(johndoe)")
	runRevokeLoggingSingleError(t, dbUser, instance, expectedLogMsg)
}

func TestGivenValidInput_WhenExecuteRevokeAccess_ThenShouldReturnOutputSuccess(t *testing.T) {
	instance := mocks.BuildQAInstanceDTO()
	dbUser := mocks.BuildDbUserJohn()
	dbUserID := getDBUserID(dbUser)
	dbUserStorage := new(mocks.DatabaseUserStorageMock)
	dbUserStorage.On("FindByID", dbUserID).Return(dbUser, nil).Once()
	accessPermissionStorage := new(mocks.AccessPermissionStorageMock)
	accessPermissionStorage.On("FindAllAccessibleInstancesIDsByUser", dbUserID).Return([]string{instance.ID}, nil).Once()
	accessPermissionStorage.On("DeleteAllByUserAndInstance", dbUserID, instance.ID).Return(nil).Once()
	expectedLogMsg := fmt.Sprintf(UserAccessRevokedAndExcludedMsg, dbUser.Username, instance.Name)
	expectedLog, _ := entity.NewAccessPermissionLog(instance.ID, dbUserID, "", expectedLogMsg, mocks.UserID, true)
	accessPermissionStorage.On("SaveLog", testdata.CompareLogs(expectedLog)).Return(nil).Once()
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	dbInstanceStorage.On("FindAllDTOs", "", "", []string{instance.ID}).Return([]*dto.DatabaseInstanceOutputDTO{instance}, nil).Once()

	uc := NewRevokeAccessPermissionUseCase(accessPermissionStorage, dbInstanceStorage, dbUserStorage)
	output, err := uc.Execute(dto.RevokeAccessInputDTO{DatabaseInstancesIDs: []string{instance.ID}, DatabaseUserID: dbUserID}, mocks.UserID)

	assert.NoError(t, err, "the error should be in the output, not in the process and has to be logged")
	assert.NotNil(t, output)
	assert.False(t, output.HasErrors)
	assert.Equal(t, "Successfully revoked access for user 'johndoe' in 1 database instances!", output.Message)
	dbUserStorage.AssertNumberOfCalls(t, "FindByID", 1)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindAllDTOs", 1)
	accessPermissionStorage.AssertNumberOfCalls(t, "FindAllAccessibleInstancesIDsByUser", 1)
	accessPermissionStorage.AssertNumberOfCalls(t, "DeleteAllByUserAndInstance", 1)
	accessPermissionStorage.AssertNumberOfCalls(t, "SaveLog", 1)
}

func runRevokeLoggingSingleError(t *testing.T, dbUser *entity.DatabaseUser, instance *dto.DatabaseInstanceOutputDTO, expectedLogMsg string) {
	dbUserID := getDBUserID(dbUser)
	dbUserStorage := new(mocks.DatabaseUserStorageMock)
	dbUserStorage.On("FindByID", dbUser.ID.String()).Return(dbUser, nil).Once()
	accessPermissionStorage := new(mocks.AccessPermissionStorageMock)
	accessPermissionStorage.On("FindAllAccessibleInstancesIDsByUser", dbUser.ID.String()).Return([]string{instance.ID}, nil).Once()
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	dbInstanceStorage.On("FindAllDTOs", "", "", []string{instance.ID}).Return([]*dto.DatabaseInstanceOutputDTO{instance}, nil).Once()

	expectedLog, _ := entity.NewAccessPermissionLog(instance.ID, dbUserID, "", expectedLogMsg, mocks.UserID, false)
	accessPermissionStorage.On("SaveLog", testdata.CompareLogs(expectedLog)).Return(nil).Once()

	uc := NewRevokeAccessPermissionUseCase(accessPermissionStorage, dbInstanceStorage, dbUserStorage)
	output, err := uc.Execute(dto.RevokeAccessInputDTO{DatabaseInstancesIDs: []string{instance.ID}, DatabaseUserID: dbUserID}, mocks.UserID)

	assert.NoError(t, err, "the error should be in the output, not in the process and has to be logged")
	assert.NotNil(t, output)
	assert.True(t, output.HasErrors)
	assert.Equal(t, SomeErrorsDuringProcessMsg, output.Message)
	dbUserStorage.AssertNumberOfCalls(t, "FindByID", 1)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindAllDTOs", 1)
	accessPermissionStorage.AssertNumberOfCalls(t, "FindAllAccessibleInstancesIDsByUser", 1)
	accessPermissionStorage.AssertNumberOfCalls(t, "SaveLog", 1)
}

func getDBUserID(dbUser *entity.DatabaseUser) string {
	if dbUser != nil {
		return dbUser.ID.String()
	}
	return ""
}

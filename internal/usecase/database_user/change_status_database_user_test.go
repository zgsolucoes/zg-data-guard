package dbuser

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/zgsolucoes/zg-data-guard/internal/usecase/common"

	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	"github.com/zgsolucoes/zg-data-guard/internal/entity"
	"github.com/zgsolucoes/zg-data-guard/testdata/mocks"
)

type RevokeAccessPermissionUseCaseMock struct {
	mock.Mock
}

func (m *RevokeAccessPermissionUseCaseMock) Execute(input dto.RevokeAccessInputDTO, operationUserID string) (*dto.RevokeAccessOutputDTO, error) {
	args := m.Called(input, operationUserID)
	return args.Get(0).(*dto.RevokeAccessOutputDTO), args.Error(1)
}

func TestGivenAnErrorInDbWhileCheckHost_WhenExecuteChangeStatus_ThenShouldReturnError(t *testing.T) {
	dbUserStorage := new(mocks.DatabaseUserStorageMock)
	dbUserStorage.On("FindByID", mocks.DbUserID).Return(&entity.DatabaseUser{}, sql.ErrConnDone).Once()

	uc := NewChangeStatusDatabaseUserUseCase(dbUserStorage, nil)
	dbUserOutput, err := uc.Execute(mocks.DbUserID, true, mocks.UserID)

	assert.Error(t, err, "error expected when some error in db")
	assert.EqualError(t, err, sql.ErrConnDone.Error())
	assert.Nil(t, dbUserOutput)
	dbUserStorage.AssertNumberOfCalls(t, "FindByID", 1)
}

func TestGivenANonexistentId_WhenExecuteChangeStatus_ThenShouldReturnError(t *testing.T) {
	dbUserStorage := new(mocks.DatabaseUserStorageMock)
	dbUserStorage.On("FindByID", mocks.DbUserID).Return(&entity.DatabaseUser{}, sql.ErrNoRows).Once()

	uc := NewChangeStatusDatabaseUserUseCase(dbUserStorage, nil)
	dbUserOutput, err := uc.Execute(mocks.DbUserID, true, mocks.UserID)

	assert.Error(t, err)
	assert.EqualError(t, err, common.ErrDatabaseUserNotFound.Error())
	assert.Nil(t, dbUserOutput)
	dbUserStorage.AssertNumberOfCalls(t, "FindByID", 1)
}

func TestGivenEnableInput_WhenExecuteChangeStatus_ThenShouldEnableDBUser(t *testing.T) {
	dbUser := mocks.BuildDbUserJohn()
	dbUser.Disable()
	dbUserStorage := new(mocks.DatabaseUserStorageMock)
	dbUserStorage.On("FindByID", dbUser.ID.String()).Return(dbUser, nil).Once()
	dbUserStorage.On("Update", mock.Anything).Return(nil).Once()

	uc := NewChangeStatusDatabaseUserUseCase(dbUserStorage, nil)
	dbUserOutput, err := uc.Execute(dbUser.ID.String(), true, mocks.UserID)

	assert.NoError(t, err)
	assert.NotNil(t, dbUserOutput)
	assert.True(t, dbUserOutput.Enabled)
	assert.Equal(t, dbUserOutput.ID, dbUser.ID.String())
	assert.NotEmpty(t, dbUserOutput.UpdatedAt)
	dbUserStorage.AssertNumberOfCalls(t, "FindByID", 1)
	dbUserStorage.AssertNumberOfCalls(t, "Update", 1)
}

func TestGivenErrorWhileUpdatingUser_WhenExecuteChangeStatus_ThenShouldReturnError(t *testing.T) {
	dbUser := mocks.BuildDbUserJohn()
	dbUser.Disable()
	dbUserStorage := new(mocks.DatabaseUserStorageMock)
	dbUserStorage.On("FindByID", dbUser.ID.String()).Return(dbUser, nil).Once()
	dbUserStorage.On("Update", mock.Anything).Return(sql.ErrConnDone).Once()

	uc := NewChangeStatusDatabaseUserUseCase(dbUserStorage, nil)
	dbUserOutput, err := uc.Execute(dbUser.ID.String(), true, mocks.UserID)

	assert.Error(t, err)
	assert.EqualError(t, err, fmt.Errorf("error when updating database user %s to new status '%t'. Cause: %w", dbUser.ID.String(), true, sql.ErrConnDone).Error())
	assert.Nil(t, dbUserOutput)
	dbUserStorage.AssertNumberOfCalls(t, "FindByID", 1)
	dbUserStorage.AssertNumberOfCalls(t, "Update", 1)
}

func TestGivenAnErrorOnRevokingAccess_WhenExecuteChangeStatusToDisable_ThenShouldReturnError(t *testing.T) {
	dbUser := mocks.BuildDbUserJohn()
	dbUserID := dbUser.ID.String()
	dbUserStorage := new(mocks.DatabaseUserStorageMock)
	mockRevoke := new(RevokeAccessPermissionUseCaseMock)
	dbUserStorage.On("FindByID", dbUserID).Return(dbUser, nil).Once()
	revokeInput := dto.RevokeAccessInputDTO{DatabaseUserID: dbUserID, DatabaseInstancesIDs: []string{}}
	mockRevoke.On("Execute", revokeInput, mocks.UserID).Return(&dto.RevokeAccessOutputDTO{}, sql.ErrConnDone).Once()

	uc := NewChangeStatusDatabaseUserUseCase(dbUserStorage, mockRevoke)
	dbUserOutput, err := uc.Execute(dbUserID, false, mocks.UserID)

	assert.Error(t, err, "error expected when some unexpected error occurs on revoking access")
	assert.Nil(t, dbUserOutput)
	assert.EqualError(t, err, fmt.Errorf("error revoking access from instances of user %s. Cause: %w", dbUserID, sql.ErrConnDone).Error())
	dbUserStorage.AssertNumberOfCalls(t, "FindByID", 1)
	mockRevoke.AssertNumberOfCalls(t, "Execute", 1)
}

func TestExecuteToDisable_WhenNoAccessibleInstancesFound_ShouldDisableUser(t *testing.T) {
	dbUser := mocks.BuildDbUserJohn()
	dbUserID := dbUser.ID.String()
	dbUserStorage := new(mocks.DatabaseUserStorageMock)
	mockRevoke := new(RevokeAccessPermissionUseCaseMock)
	dbUserStorage.On("FindByID", dbUserID).Return(dbUser, nil).Once()
	revokeInput := dto.RevokeAccessInputDTO{DatabaseUserID: dbUserID, DatabaseInstancesIDs: []string{}}
	mockRevoke.On("Execute", revokeInput, mocks.UserID).Return(&dto.RevokeAccessOutputDTO{}, common.ErrNoAccessibleInstancesFound).Once()
	dbUserStorage.On("Update", mock.Anything).Return(nil).Once()

	uc := NewChangeStatusDatabaseUserUseCase(dbUserStorage, mockRevoke)
	dbUserOutput, err := uc.Execute(dbUserID, false, mocks.UserID)

	assert.NoError(t, err, "no error expected when no accessible instances found to revoke access so user can be disabled")
	assert.NotNil(t, dbUserOutput)
	assert.False(t, dbUserOutput.Enabled)
	assert.Equal(t, dbUserOutput.ID, dbUserID)
	assert.NotEmpty(t, dbUserOutput.UpdatedAt)
	dbUserStorage.AssertNumberOfCalls(t, "FindByID", 1)
	mockRevoke.AssertNumberOfCalls(t, "Execute", 1)
	dbUserStorage.AssertNumberOfCalls(t, "Update", 1)
}

func TestExecuteToDisable_WhenRevokeAllAccess_ShouldDisableUser(t *testing.T) {
	dbUser := mocks.BuildDbUserJohn()
	dbUserID := dbUser.ID.String()
	dbUserStorage := new(mocks.DatabaseUserStorageMock)
	mockRevoke := new(RevokeAccessPermissionUseCaseMock)
	dbUserStorage.On("FindByID", dbUserID).Return(dbUser, nil).Once()
	revokeInput := dto.RevokeAccessInputDTO{DatabaseUserID: dbUserID, DatabaseInstancesIDs: []string{}}
	mockRevoke.On("Execute", revokeInput, mocks.UserID).Return(&dto.RevokeAccessOutputDTO{HasErrors: false}, nil).Once()
	dbUserStorage.On("Update", mock.Anything).Return(nil).Once()

	uc := NewChangeStatusDatabaseUserUseCase(dbUserStorage, mockRevoke)
	dbUserOutput, err := uc.Execute(dbUserID, false, mocks.UserID)

	assert.NoError(t, err, "no error expected when all access is revoked so user can be disabled")
	assert.NotNil(t, dbUserOutput)
	assert.False(t, dbUserOutput.Enabled)
	assert.Equal(t, dbUserOutput.ID, dbUserID)
	assert.NotEmpty(t, dbUserOutput.UpdatedAt)
	dbUserStorage.AssertNumberOfCalls(t, "FindByID", 1)
	mockRevoke.AssertNumberOfCalls(t, "Execute", 1)
	dbUserStorage.AssertNumberOfCalls(t, "Update", 1)
}

func TestExecuteToDisable_WhenRevokeAccessHasErrors_ShouldReturnError(t *testing.T) {
	dbUser := mocks.BuildDbUserJohn()
	dbUserID := dbUser.ID.String()
	dbUserStorage := new(mocks.DatabaseUserStorageMock)
	mockRevoke := new(RevokeAccessPermissionUseCaseMock)
	dbUserStorage.On("FindByID", dbUserID).Return(dbUser, nil).Once()
	revokeInput := dto.RevokeAccessInputDTO{DatabaseUserID: dbUserID, DatabaseInstancesIDs: []string{}}
	mockRevoke.On("Execute", revokeInput, mocks.UserID).Return(&dto.RevokeAccessOutputDTO{HasErrors: true}, nil).Once()

	uc := NewChangeStatusDatabaseUserUseCase(dbUserStorage, mockRevoke)
	dbUserOutput, err := uc.Execute(dbUserID, false, mocks.UserID)

	assert.Error(t, err, "error expected when could not revoke all access of the user to disable it")
	assert.Nil(t, dbUserOutput)
	assert.EqualError(t, err, ErrCouldNotRevokeAllAccess.Error())
	dbUserStorage.AssertNumberOfCalls(t, "FindByID", 1)
	mockRevoke.AssertNumberOfCalls(t, "Execute", 1)
}

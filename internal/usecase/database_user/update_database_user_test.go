package dbuser

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zgsolucoes/zg-data-guard/internal/usecase/common"

	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	"github.com/zgsolucoes/zg-data-guard/internal/entity"
	"github.com/zgsolucoes/zg-data-guard/testdata/mocks"
)

func TestGivenAnNonexistentId_WhenExecuteUpdate_ThenShouldReturnError(t *testing.T) {
	dbUserStorage := new(mocks.DatabaseUserStorageMock)
	dbRoleStorage := new(mocks.DatabaseRoleStorageMock)
	dbUserStorage.On("FindByID", mocks.DbUserID).Return(&entity.DatabaseUser{}, sql.ErrNoRows).Once()
	uc := NewUpdateDatabaseUserUseCase(dbUserStorage, dbRoleStorage, nil)

	dbUserOutput, err := uc.Execute(mocks.ValidUpdateDBUserInput, mocks.DbUserID, mocks.UserID)

	assert.Error(t, err, "error expected when dbUser not found")
	assert.EqualError(t, err, common.ErrDatabaseUserNotFound.Error())
	assert.Nil(t, dbUserOutput)
	dbUserStorage.AssertNumberOfCalls(t, "FindByID", 1)
}

func TestGivenAnErrorInDbWhileFindById_WhenExecuteUpdate_ThenShouldReturnError(t *testing.T) {
	dbUserStorage := new(mocks.DatabaseUserStorageMock)
	dbRoleStorage := new(mocks.DatabaseRoleStorageMock)
	dbUserStorage.On("FindByID", mocks.DbUserID).Return(&entity.DatabaseUser{}, sql.ErrConnDone).Once()
	uc := NewUpdateDatabaseUserUseCase(dbUserStorage, dbRoleStorage, nil)

	dbUserOutput, err := uc.Execute(mocks.ValidUpdateDBUserInput, mocks.DbUserID, mocks.UserID)

	assert.Error(t, err, "error expected when some error in db")
	assert.EqualError(t, err, sql.ErrConnDone.Error())
	assert.Nil(t, dbUserOutput)
	dbUserStorage.AssertNumberOfCalls(t, "FindByID", 1)
}

func TestGivenAnErrorInDbWhileCheckingUserAccess_WhenExecuteUpdate_ThenShouldReturnError(t *testing.T) {
	dbUserStorage := new(mocks.DatabaseUserStorageMock)
	accessStorage := new(mocks.AccessPermissionStorageMock)
	dbUserStorage.On("FindByID", mocks.DbUserID).Return(&entity.DatabaseUser{}, nil).Once()
	accessStorage.On("CheckIfUserHasAccessPermission", mocks.DbUserID).Return(false, sql.ErrConnDone).Once()
	uc := NewUpdateDatabaseUserUseCase(dbUserStorage, nil, accessStorage)

	dbUserOutput, err := uc.Execute(mocks.ValidUpdateDBUserInput, mocks.DbUserID, mocks.UserID)

	assert.Error(t, err, "error expected when some error in db")
	assert.EqualError(t, err, sql.ErrConnDone.Error())
	assert.Nil(t, dbUserOutput)
	dbUserStorage.AssertNumberOfCalls(t, "FindByID", 1)
	accessStorage.AssertNumberOfCalls(t, "CheckIfUserHasAccessPermission", 1)
}

func TestGivenAnUserWithAccessPermissionsDefined_WhenExecuteUpdateToADifferentRole_ThenShouldReturnError(t *testing.T) {
	dbUserStorage := new(mocks.DatabaseUserStorageMock)
	accessStorage := new(mocks.AccessPermissionStorageMock)
	dbUser, _ := entity.NewDatabaseUser(
		mocks.ValidDBUserInput.Name,
		mocks.ValidDBUserInput.Email,
		mocks.ValidDBUserInput.Team,
		mocks.ValidDBUserInput.Position,
		mocks.ValidDBUserInput.DatabaseRoleID,
		mocks.UserID)
	dbUserStorage.On("FindByID", dbUser.ID.String()).Return(dbUser, nil).Once()
	accessStorage.On("CheckIfUserHasAccessPermission", dbUser.ID.String()).Return(true, nil).Once()
	dbUserStorage.On("Update", dbUser).Return(sql.ErrConnDone).Once()
	uc := NewUpdateDatabaseUserUseCase(dbUserStorage, nil, accessStorage)

	dbUserOutput, err := uc.Execute(mocks.ValidUpdateDBUserInput, dbUser.ID.String(), mocks.UserID)

	assert.Error(t, err, "error expected when database user has access permissions and cannot have their role changed")
	assert.EqualError(t, err, ErrDatabaseUserHasAccessPermissions.Error())
	assert.Nil(t, dbUserOutput)
	dbUserStorage.AssertNumberOfCalls(t, "FindByID", 1)
	accessStorage.AssertNumberOfCalls(t, "CheckIfUserHasAccessPermission", 1)
}

func TestGivenAnErrorInDbWhileValidateRole_WhenExecuteUpdate_ThenShouldReturnError(t *testing.T) {
	dbUserStorage := new(mocks.DatabaseUserStorageMock)
	dbRoleStorage := new(mocks.DatabaseRoleStorageMock)
	dbUserStorage.On("FindByID", mocks.DbUserID).Return(&entity.DatabaseUser{DatabaseRoleID: mocks.ValidUpdateDBUserInput.DatabaseRoleID}, nil).Once()
	dbRoleStorage.On("FindByID", mocks.ValidUpdateDBUserInput.DatabaseRoleID).Return(&entity.DatabaseRole{}, sql.ErrConnDone).Once()
	uc := NewUpdateDatabaseUserUseCase(dbUserStorage, dbRoleStorage, nil)

	dbUserOutput, err := uc.Execute(mocks.ValidUpdateDBUserInput, mocks.DbUserID, mocks.UserID)

	assert.Error(t, err, "error expected when some error in db while fetching database role")
	assert.EqualError(t, err, sql.ErrConnDone.Error())
	assert.Nil(t, dbUserOutput)
	dbUserStorage.AssertNumberOfCalls(t, "FindByID", 1)
	dbRoleStorage.AssertNumberOfCalls(t, "FindByID", 1)
}

func TestGivenAnNonexistentRoleId_WhenExecuteUpdate_ThenShouldReturnError(t *testing.T) {
	dbUserStorage := new(mocks.DatabaseUserStorageMock)
	dbRoleStorage := new(mocks.DatabaseRoleStorageMock)
	dbUserStorage.On("FindByID", mocks.DbUserID).Return(&entity.DatabaseUser{DatabaseRoleID: mocks.ValidUpdateDBUserInput.DatabaseRoleID}, nil).Once()
	dbRoleStorage.On("FindByID", mocks.ValidUpdateDBUserInput.DatabaseRoleID).Return(&entity.DatabaseRole{}, sql.ErrNoRows).Once()
	uc := NewUpdateDatabaseUserUseCase(dbUserStorage, dbRoleStorage, nil)

	dbUserOutput, err := uc.Execute(mocks.ValidUpdateDBUserInput, mocks.DbUserID, mocks.UserID)

	assert.Error(t, err, "error expected when nonexistent database role")
	assert.EqualError(t, err, ErrDatabaseRoleNotFound.Error())
	assert.Nil(t, dbUserOutput)
	dbUserStorage.AssertNumberOfCalls(t, "FindByID", 1)
	dbRoleStorage.AssertNumberOfCalls(t, "FindByID", 1)
}

func TestGivenAnInvalidUpdateInput_WhenExecuteUpdate_ThenShouldReturnError(t *testing.T) {
	dbUserStorage := new(mocks.DatabaseUserStorageMock)
	dbRoleStorage := new(mocks.DatabaseRoleStorageMock)
	accessStorage := new(mocks.AccessPermissionStorageMock)
	dbUser, _ := entity.NewDatabaseUser(
		mocks.ValidDBUserInput.Name,
		mocks.ValidDBUserInput.Email,
		mocks.ValidDBUserInput.Team,
		mocks.ValidDBUserInput.Position,
		mocks.ValidDBUserInput.DatabaseRoleID,
		mocks.UserID)
	dbUserStorage.On("FindByID", dbUser.ID.String()).Return(dbUser, nil).Once()
	accessStorage.On("CheckIfUserHasAccessPermission", dbUser.ID.String()).Return(false, nil).Once()
	dbRoleStorage.On("FindByID", mocks.ValidUpdateDBUserInput.DatabaseRoleID).Return(&entity.DatabaseRole{}, nil).Once()

	invalidUpdateInput := dto.UpdateDatabaseUserInputDTO{
		DatabaseRoleID: mocks.ValidUpdateDBUserInput.DatabaseRoleID,
	}
	uc := NewUpdateDatabaseUserUseCase(dbUserStorage, dbRoleStorage, accessStorage)

	dbUserOutput, err := uc.Execute(invalidUpdateInput, dbUser.ID.String(), mocks.UserID)

	assert.Error(t, err, "error expected when some error in db")
	assert.EqualError(t, err, entity.ErrInvalidName.Error())
	assert.Nil(t, dbUserOutput)
	dbUserStorage.AssertNumberOfCalls(t, "FindByID", 1)
	accessStorage.AssertNumberOfCalls(t, "CheckIfUserHasAccessPermission", 1)
	dbRoleStorage.AssertNumberOfCalls(t, "FindByID", 1)
}

func TestGivenAnErrorInDbWhileUpdate_WhenExecuteUpdate_ThenShouldReturnError(t *testing.T) {
	dbUserStorage := new(mocks.DatabaseUserStorageMock)
	dbRoleStorage := new(mocks.DatabaseRoleStorageMock)
	accessStorage := new(mocks.AccessPermissionStorageMock)
	dbUser, _ := entity.NewDatabaseUser(
		mocks.ValidDBUserInput.Name,
		mocks.ValidDBUserInput.Email,
		mocks.ValidDBUserInput.Team,
		mocks.ValidDBUserInput.Position,
		mocks.ValidDBUserInput.DatabaseRoleID,
		mocks.UserID)
	dbUserStorage.On("FindByID", dbUser.ID.String()).Return(dbUser, nil).Once()
	accessStorage.On("CheckIfUserHasAccessPermission", dbUser.ID.String()).Return(false, nil).Once()
	dbRoleStorage.On("FindByID", mocks.ValidUpdateDBUserInput.DatabaseRoleID).Return(&entity.DatabaseRole{}, nil).Once()
	dbUserStorage.On("Update", dbUser).Return(sql.ErrConnDone).Once()
	uc := NewUpdateDatabaseUserUseCase(dbUserStorage, dbRoleStorage, accessStorage)

	dbUserOutput, err := uc.Execute(mocks.ValidUpdateDBUserInput, dbUser.ID.String(), mocks.UserID)

	assert.Error(t, err, "error expected when some error in db")
	assert.EqualError(t, err, sql.ErrConnDone.Error())
	assert.Nil(t, dbUserOutput)
	dbUserStorage.AssertNumberOfCalls(t, "FindByID", 1)
	accessStorage.AssertNumberOfCalls(t, "CheckIfUserHasAccessPermission", 1)
	dbRoleStorage.AssertNumberOfCalls(t, "FindByID", 1)
	dbUserStorage.AssertNumberOfCalls(t, "Update", 1)
}

func TestGivenAValidIdAndInput_WhenExecuteUpdate_ThenShouldReturnDatabaseUser(t *testing.T) {
	dbUserStorage := new(mocks.DatabaseUserStorageMock)
	dbRoleStorage := new(mocks.DatabaseRoleStorageMock)
	accessStorage := new(mocks.AccessPermissionStorageMock)
	dbUser, _ := entity.NewDatabaseUser(
		mocks.ValidDBUserInput.Name,
		mocks.ValidDBUserInput.Email,
		mocks.ValidDBUserInput.Team,
		mocks.ValidDBUserInput.Position,
		mocks.ValidDBUserInput.DatabaseRoleID,
		mocks.UserID)
	dbUserStorage.On("FindByID", dbUser.ID.String()).Return(dbUser, nil).Once()
	accessStorage.On("CheckIfUserHasAccessPermission", dbUser.ID.String()).Return(false, nil).Once()
	dbRoleStorage.On("FindByID", mocks.ValidUpdateDBUserInput.DatabaseRoleID).Return(&entity.DatabaseRole{}, nil).Once()
	dbUserStorage.On("Update", dbUser).Return(nil).Once()
	uc := NewUpdateDatabaseUserUseCase(dbUserStorage, dbRoleStorage, accessStorage)
	dbUser.Disable()

	dbUserOutput, err := uc.Execute(mocks.ValidUpdateDBUserInput, dbUser.ID.String(), mocks.UserID)

	assert.NoError(t, err, "no error expected with an existent id and valid input")
	assert.NotNil(t, dbUserOutput, "dbUser should not be nil")
	assert.Equal(t, dbUser.ID.String(), dbUserOutput.ID)
	assert.Equal(t, mocks.ValidUpdateDBUserInput.Name, dbUserOutput.Name)
	assert.Equal(t, mocks.ValidUpdateDBUserInput.Team, dbUserOutput.Team)
	assert.Equal(t, mocks.ValidUpdateDBUserInput.Position, dbUserOutput.Position)
	assert.Equal(t, mocks.ValidUpdateDBUserInput.DatabaseRoleID, dbUserOutput.DatabaseRoleID)
	assert.Equal(t, dbUser.CreatedAt, dbUserOutput.CreatedAt)
	assert.Equal(t, dbUser.UpdatedAt, *dbUserOutput.UpdatedAt)
	assert.NotEmpty(t, dbUser.DisabledAt)
	assert.False(t, dbUser.Enabled)
	dbUserStorage.AssertNumberOfCalls(t, "FindByID", 1)
	accessStorage.AssertNumberOfCalls(t, "CheckIfUserHasAccessPermission", 1)
	dbRoleStorage.AssertNumberOfCalls(t, "FindByID", 1)
	dbUserStorage.AssertNumberOfCalls(t, "Update", 1)
}

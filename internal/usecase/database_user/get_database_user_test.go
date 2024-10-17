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

func TestGivenANonexistentId_WhenExecuteGetDatabaseUser_ThenShouldReturnError(t *testing.T) {
	dbUserStorage := new(mocks.DatabaseUserStorageMock)
	uc := NewGetDatabaseUserUseCase(dbUserStorage)

	dbUserStorage.On("FindDTOByID", mocks.DbUserID).Return(&dto.DatabaseUserOutputDTO{}, sql.ErrNoRows)
	_, err := uc.Execute(mocks.DbUserID)
	assert.Error(t, err)
	assert.EqualError(t, err, common.ErrDatabaseUserNotFound.Error())
	dbUserStorage.AssertNumberOfCalls(t, "FindDTOByID", 1)
}

func TestGivenAnErrorInDb_WhenExecuteGetDatabaseUser_ThenShouldReturnError(t *testing.T) {
	dbUserStorage := new(mocks.DatabaseUserStorageMock)
	uc := NewGetDatabaseUserUseCase(dbUserStorage)

	dbUserStorage.On("FindDTOByID", mocks.DbUserID).Return(&dto.DatabaseUserOutputDTO{}, sql.ErrTxDone)
	_, err := uc.Execute(mocks.DbUserID)
	assert.Error(t, err)
	assert.EqualError(t, err, sql.ErrTxDone.Error())
	dbUserStorage.AssertNumberOfCalls(t, "FindDTOByID", 1)
}

func TestGivenAValidId_WhenExecuteGet_ThenShouldReturnDatabaseUser(t *testing.T) {
	dbUserStorage := new(mocks.DatabaseUserStorageMock)
	uc := NewGetDatabaseUserUseCase(dbUserStorage)

	dbUser := mocks.BuildDbUserFooDTO()
	dbUserStorage.On("FindDTOByID", mocks.DbUserID).Return(dbUser, nil)
	dbDTO, err := uc.Execute(mocks.DbUserID)
	assert.NoError(t, err)
	assert.NotNil(t, dbDTO)
	assert.Equal(t, dbUser.ID, dbDTO.ID)
	assert.Equal(t, dbUser.Name, dbDTO.Name)
	assert.Equal(t, dbUser.Email, dbDTO.Email)
	assert.Equal(t, dbUser.Username, dbDTO.Username)
	assert.Equal(t, dbUser.DatabaseRoleID, dbDTO.DatabaseRoleID)
	assert.Equal(t, dbUser.DatabaseRoleName, dbDTO.DatabaseRoleName)
	assert.Equal(t, dbUser.DatabaseRoleDisplayName, dbDTO.DatabaseRoleDisplayName)
	assert.Equal(t, dbUser.Team, dbDTO.Team)
	assert.Equal(t, dbUser.Position, dbDTO.Position)
	assert.Equal(t, dbUser.Enabled, dbDTO.Enabled)
	assert.Equal(t, dbUser.CreatedByUserID, dbDTO.CreatedByUserID)
	assert.Equal(t, dbUser.CreatedByUser, dbDTO.CreatedByUser)
	assert.Equal(t, dbUser.CreatedAt, dbDTO.CreatedAt)
	assert.Equal(t, dbUser.UpdatedAt, dbDTO.UpdatedAt)
	assert.Equal(t, dbUser.DisabledAt, dbDTO.DisabledAt)
	assert.Empty(t, dbUser.Password)
	dbUserStorage.AssertNumberOfCalls(t, "FindDTOByID", 1)
}

func TestGivenAnNonexistentId_WhenFetchCredentials_ThenShouldReturnError(t *testing.T) {
	dbUserStorage := new(mocks.DatabaseUserStorageMock)
	dbUserStorage.On("FindByID", mocks.DbUserID).Return(&entity.DatabaseUser{}, sql.ErrNoRows).Once()

	uc := NewGetDatabaseUserUseCase(dbUserStorage)
	output, err := uc.FetchCredentials(mocks.DbUserID, mocks.UserID)

	assert.Error(t, err, "error expected when database user not found")
	assert.EqualError(t, err, common.ErrDatabaseUserNotFound.Error())
	assert.Nil(t, output)
	dbUserStorage.AssertNumberOfCalls(t, "FindByID", 1)
}

func TestGivenAnErrorInDb_WhenFetchCredentials_ThenShouldReturnError(t *testing.T) {
	dbUserStorage := new(mocks.DatabaseUserStorageMock)
	dbUserStorage.On("FindByID", mocks.DbUserID).Return(&entity.DatabaseUser{}, sql.ErrConnDone).Once()

	uc := NewGetDatabaseUserUseCase(dbUserStorage)
	output, err := uc.FetchCredentials(mocks.DbUserID, mocks.UserID)

	assert.Error(t, err, "error expected when some error in db")
	assert.EqualError(t, err, sql.ErrConnDone.Error())
	assert.Nil(t, output)
	dbUserStorage.AssertNumberOfCalls(t, "FindByID", 1)
}

func TestGivenAValidId_WhenFetchCredentials_ThenShouldReturnCredentials(t *testing.T) {
	dbUserStorage := new(mocks.DatabaseUserStorageMock)
	dbUser := mocks.BuildDbUserJohn()
	dbUser.CipherPassword = "49e5bf3f6a45a75c972c68b39d640e53f050a6a0b4125ff9"
	dbUserStorage.On("FindByID", dbUser.ID.String()).Return(dbUser, nil).Once()

	uc := NewGetDatabaseUserUseCase(dbUserStorage)
	output, err := uc.FetchCredentials(dbUser.ID.String(), mocks.UserID)

	assert.NoError(t, err, "no error expected with an existent id")
	assert.NotNil(t, output, "credentials DTO should not be nil")
	assert.Equal(t, dbUser.Username, output.User)
	assert.Equal(t, "P6\x10\xbc2.\xad\x82", output.Password)
	dbUserStorage.AssertNumberOfCalls(t, "FindByID", 1)
}

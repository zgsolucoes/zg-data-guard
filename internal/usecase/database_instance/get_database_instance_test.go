package instance

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	"github.com/zgsolucoes/zg-data-guard/testdata/mocks"
)

func TestGivenANonexistentId_WhenExecuteGet_ThenShouldReturnError(t *testing.T) {
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	dbInstanceStorage.On("FindDTOByID", mocks.DatabaseInstanceId).Return(&dto.DatabaseInstanceOutputDTO{}, sql.ErrNoRows).Once()

	uc := NewGetDatabaseInstanceUseCase(dbInstanceStorage)
	output, err := uc.Execute(mocks.DatabaseInstanceId)

	assert.Error(t, err, "error expected when database instance not found")
	assert.EqualError(t, err, ErrDatabaseInstanceNotFound.Error())
	assert.Nil(t, output)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindDTOByID", 1)
}

func TestGivenAnErrorInDb_WhenExecuteGet_ThenShouldReturnError(t *testing.T) {
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	dbInstanceStorage.On("FindDTOByID", mocks.DatabaseInstanceId).Return(&dto.DatabaseInstanceOutputDTO{}, sql.ErrConnDone).Once()

	uc := NewGetDatabaseInstanceUseCase(dbInstanceStorage)
	output, err := uc.Execute(mocks.DatabaseInstanceId)

	assert.Error(t, err, "error expected when some error in db")
	assert.EqualError(t, err, sql.ErrConnDone.Error())
	assert.Nil(t, output)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindDTOByID", 1)
}

func TestGivenAValidId_WhenExecuteGet_ThenShouldReturnDatabaseInstance(t *testing.T) {
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	dbInstanceOutput := mocks.BuildFullDataInstanceExample()
	dbInstanceStorage.On("FindDTOByID", dbInstanceOutput.ID).Return(dbInstanceOutput, nil).Once()

	uc := NewGetDatabaseInstanceUseCase(dbInstanceStorage)
	output, err := uc.Execute(dbInstanceOutput.ID)

	assert.NoError(t, err, "no error expected with an existent id")
	assert.NotNil(t, output, "database instance should not be nil")
	assert.Equal(t, dbInstanceOutput.ID, output.ID)
	assert.Equal(t, dbInstanceOutput.Name, output.Name)
	assert.Equal(t, dbInstanceOutput.Host, output.Host)
	assert.Equal(t, dbInstanceOutput.Port, output.Port)
	assert.Equal(t, dbInstanceOutput.HostConnection, output.HostConnection)
	assert.Equal(t, dbInstanceOutput.PortConnection, output.PortConnection)
	assert.Equal(t, dbInstanceOutput.AdminUser, output.AdminUser)
	assert.Equal(t, dbInstanceOutput.EcosystemID, output.EcosystemID)
	assert.Equal(t, dbInstanceOutput.EcosystemName, output.EcosystemName)
	assert.Equal(t, dbInstanceOutput.DatabaseTechnologyID, output.DatabaseTechnologyID)
	assert.Equal(t, dbInstanceOutput.DatabaseTechnologyName, output.DatabaseTechnologyName)
	assert.Equal(t, dbInstanceOutput.DatabaseTechnologyVersion, output.DatabaseTechnologyVersion)
	assert.Equal(t, dbInstanceOutput.Note, output.Note)
	assert.Equal(t, dbInstanceOutput.CreatedByUserID, output.CreatedByUserID)
	assert.Equal(t, dbInstanceOutput.CreatedByUser, output.CreatedByUser)
	assert.Equal(t, dbInstanceOutput.CreatedAt, output.CreatedAt)
	assert.Equal(t, dbInstanceOutput.UpdatedAt, output.UpdatedAt)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindDTOByID", 1)
}

func TestGivenAnNonexistentId_WhenFetchCredentials_ThenShouldReturnError(t *testing.T) {
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	dbInstanceStorage.On("FindDTOByID", mocks.DatabaseInstanceId).Return(&dto.DatabaseInstanceOutputDTO{}, sql.ErrNoRows).Once()

	uc := NewGetDatabaseInstanceUseCase(dbInstanceStorage)
	output, err := uc.FetchCredentials(mocks.DatabaseInstanceId, mocks.UserID)

	assert.Error(t, err, "error expected when database instance not found")
	assert.EqualError(t, err, ErrDatabaseInstanceNotFound.Error())
	assert.Nil(t, output)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindDTOByID", 1)
}

func TestGivenAnErrorInDb_WhenFetchCredentials_ThenShouldReturnError(t *testing.T) {
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	dbInstanceStorage.On("FindDTOByID", mocks.DatabaseInstanceId).Return(&dto.DatabaseInstanceOutputDTO{}, sql.ErrConnDone).Once()

	uc := NewGetDatabaseInstanceUseCase(dbInstanceStorage)
	output, err := uc.FetchCredentials(mocks.DatabaseInstanceId, mocks.UserID)

	assert.Error(t, err, "error expected when some error in db")
	assert.EqualError(t, err, sql.ErrConnDone.Error())
	assert.Nil(t, output)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindDTOByID", 1)
}

func TestGivenAValidId_WhenFetchCredentials_ThenShouldReturnCredentials(t *testing.T) {
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	dbInstanceOutput := mocks.BuildFullDataInstanceExample()
	dbInstanceOutput.AdminPassword = "49e5bf3f6a45a75c972c68b39d640e53f050a6a0b4125ff9"
	dbInstanceStorage.On("FindDTOByID", dbInstanceOutput.ID).Return(dbInstanceOutput, nil).Once()

	uc := NewGetDatabaseInstanceUseCase(dbInstanceStorage)
	output, err := uc.FetchCredentials(dbInstanceOutput.ID, mocks.UserID)

	assert.NoError(t, err, "no error expected with an existent id")
	assert.NotNil(t, output, "credentials DTO should not be nil")
	assert.Equal(t, dbInstanceOutput.AdminUser, output.User)
	assert.Equal(t, "P6\x10\xbc2.\xad\x82", output.Password)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindDTOByID", 1)
}

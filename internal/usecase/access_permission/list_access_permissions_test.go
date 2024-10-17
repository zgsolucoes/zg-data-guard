package accesspermission

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	"github.com/zgsolucoes/zg-data-guard/testdata/mocks"
)

func TestGivenAnErrorInDb_WhenExecuteListAccessPermissions_ThenShouldReturnError(t *testing.T) {
	accessStorage := new(mocks.AccessPermissionStorageMock)
	accessStorage.On("FindAllDTOs", mocks.DatabaseID, mocks.DbUserID, mocks.DatabaseInstanceId).Return([]*dto.AccessPermissionOutputDTO{}, sql.ErrConnDone).Once()

	uc := NewListAccessPermissionsUseCase(accessStorage)
	accessObtained, err := uc.Execute(mocks.DatabaseID, mocks.DbUserID, mocks.DatabaseInstanceId)

	assert.Error(t, err, "error expected when some error in db")
	assert.EqualError(t, err, sql.ErrConnDone.Error())
	assert.Nil(t, accessObtained)
	assert.Equal(t, len(accessObtained), 0, "0 access permissions expected")
	accessStorage.AssertNumberOfCalls(t, "FindAllDTOs", 1)
}

func TestGivenSomePermissions_WhenExecuteListAccessPermissions_ThenShouldListAllPermissions(t *testing.T) {
	permissionsList := mocks.BuildAccessPermissionsDTOList()
	accessStorage := new(mocks.AccessPermissionStorageMock)
	accessStorage.On("FindAllDTOs", mocks.DatabaseID, mocks.DbUserID, mocks.DatabaseInstanceId).Return(permissionsList, nil).Once()

	uc := NewListAccessPermissionsUseCase(accessStorage)
	permissionsObtained, err := uc.Execute(mocks.DatabaseID, mocks.DbUserID, mocks.DatabaseInstanceId)

	assert.NoError(t, err, "no error expected")
	assert.Equal(t, len(permissionsObtained), len(permissionsList), "3 access permissions expected")
	accessStorage.AssertNumberOfCalls(t, "FindAllDTOs", 1)
}

package accesspermission

import (
	"log"

	"github.com/zgsolucoes/zg-data-guard/internal/database/storage"
	"github.com/zgsolucoes/zg-data-guard/internal/dto"
)

type ListAccessPermissionsUseCase struct {
	AccessPermissionStorage storage.AccessPermissionStorage
}

func NewListAccessPermissionsUseCase(permissionStorage storage.AccessPermissionStorage) *ListAccessPermissionsUseCase {
	return &ListAccessPermissionsUseCase{
		AccessPermissionStorage: permissionStorage,
	}
}

func (uc *ListAccessPermissionsUseCase) Execute(databaseID, databaseUserID, databaseInstanceID string) ([]*dto.AccessPermissionOutputDTO, error) {
	accessDTOs, err := uc.AccessPermissionStorage.FindAllDTOs(databaseID, databaseUserID, databaseInstanceID)
	if err != nil {
		log.Printf("Error fetching access permissions! Cause: %v", err.Error())
		return nil, err
	}
	log.Printf("List of access permissions loaded successfully!")
	return accessDTOs, nil
}

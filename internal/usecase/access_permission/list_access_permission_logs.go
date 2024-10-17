package accesspermission

import (
	"fmt"
	"log"

	"github.com/zgsolucoes/zg-data-guard/internal/database/storage"
	"github.com/zgsolucoes/zg-data-guard/internal/dto"
)

type ListAccessPermissionLogsUseCase struct {
	AccessPermissionStorage storage.AccessPermissionStorage
}

func NewListAccessPermissionLogsUseCase(permissionStorage storage.AccessPermissionStorage) *ListAccessPermissionLogsUseCase {
	return &ListAccessPermissionLogsUseCase{
		AccessPermissionStorage: permissionStorage,
	}
}

func (uc *ListAccessPermissionLogsUseCase) Execute(page, limit int) ([]*dto.AccessPermissionLogOutputDTO, int, error) {
	logsDTOs, err := uc.AccessPermissionStorage.FindAllLogsDTOs(page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("error fetching access permission logs! Cause: %w", err)
	}
	totalCount, err := uc.AccessPermissionStorage.LogCount()
	if err != nil {
		return nil, 0, fmt.Errorf("error fetching access permission logs count! Cause: %w", err)
	}
	log.Printf("List of access permission logs loaded successfully!")
	return logsDTOs, totalCount, nil
}

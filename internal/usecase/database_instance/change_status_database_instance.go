package instance

import (
	"fmt"
	"log"
	"time"

	"github.com/zgsolucoes/zg-data-guard/internal/database/storage"
	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	"github.com/zgsolucoes/zg-data-guard/internal/entity"
	"github.com/zgsolucoes/zg-data-guard/internal/usecase/common"
)

type ChangeStatusDatabaseInstanceUseCase struct {
	DatabaseInstanceStorage storage.DatabaseInstanceStorage
	DatabaseStorage         storage.DatabaseStorage
	AccessPermissionStorage storage.AccessPermissionStorage
}

func NewChangeStatusDatabaseInstanceUseCase(
	dbInstanceStorage storage.DatabaseInstanceStorage,
	databaseStorage storage.DatabaseStorage,
	accessPermissionStorage storage.AccessPermissionStorage) *ChangeStatusDatabaseInstanceUseCase {
	return &ChangeStatusDatabaseInstanceUseCase{
		DatabaseInstanceStorage: dbInstanceStorage,
		DatabaseStorage:         databaseStorage,
		AccessPermissionStorage: accessPermissionStorage,
	}
}

func (uc *ChangeStatusDatabaseInstanceUseCase) Execute(userID string, enabled bool, operationUserID string) (*dto.ChangeStatusOutputDTO, error) {
	dbInstance, err := uc.DatabaseInstanceStorage.FindByID(userID)
	if err != nil {
		return nil, common.HandleFindError(err, ErrDatabaseInstanceNotFound)
	}
	log.Printf("Changing status of database instance '%s' to '%t'. Requester: %s", dbInstance.ID.String(), enabled, operationUserID)
	if err = uc.changeInstanceStatus(dbInstance, enabled, operationUserID); err != nil {
		return nil, err
	}
	if err = uc.DatabaseInstanceStorage.Update(dbInstance); err != nil {
		return nil, fmt.Errorf("error when updating database instance %s to new status '%t'. Cause: %w", dbInstance.ID.String(), enabled, err)
	}
	if err = uc.registerLog(dbInstance.ID.String(), dbInstance.Name, operationUserID, enabled); err != nil {
		return nil, fmt.Errorf("error when registering log for database instance '%s'. Cause: %w", dbInstance.Name, err)
	}
	log.Printf("Database instance '%s' status changed to '%t' successfully. Requester: %s", dbInstance.ID.String(), enabled, operationUserID)
	return uc.buildOutputDTO(dbInstance), nil
}

func (uc *ChangeStatusDatabaseInstanceUseCase) changeInstanceStatus(dbInstance *entity.DatabaseInstance, enabled bool, operationUserID string) error {
	if enabled {
		dbInstance.Enable()
		return nil
	}
	if err := uc.revokeInstanceAccess(dbInstance.ID.String(), operationUserID); err != nil {
		return fmt.Errorf("error when revoking access from database instance '%s'. Cause: %w", dbInstance.Name, err)
	}
	log.Printf("All access from database instance '%s' revoked successfully. Requester: %s", dbInstance.Name, operationUserID)
	if err := uc.deactivateDatabases(dbInstance.ID.String(), operationUserID); err != nil {
		return fmt.Errorf("error when deactivating databases from database instance '%s'. Cause: %w", dbInstance.Name, err)
	}
	log.Printf("All databases from database instance '%s' deactivated successfully. Requester: %s", dbInstance.Name, operationUserID)
	dbInstance.Disable()
	return nil
}

func (uc *ChangeStatusDatabaseInstanceUseCase) revokeInstanceAccess(dbInstanceID, operationUserID string) error {
	log.Printf("Revoking all access from database instance '%s'. Requester: %s", dbInstanceID, operationUserID)
	return uc.AccessPermissionStorage.DeleteAllByInstance(dbInstanceID)
}

func (uc *ChangeStatusDatabaseInstanceUseCase) deactivateDatabases(dbInstanceID, operationUserID string) error {
	log.Printf("Deactivating all databases from database instance '%s'. Requester: %s", dbInstanceID, operationUserID)
	return uc.DatabaseStorage.DeactivateAllByInstance(dbInstanceID)
}

func (uc *ChangeStatusDatabaseInstanceUseCase) registerLog(instanceID, instanceName, operationUserID string, enabled bool) error {
	var message string
	if enabled {
		message = fmt.Sprintf(common.InstanceEnabledSuccessMsg, instanceName)
	} else {
		message = fmt.Sprintf(common.InstanceDisabledSuccessMsg, instanceName)
	}
	grantLog, err := entity.NewAccessPermissionLog(instanceID, "", "", message, operationUserID, true)
	if err != nil {
		return fmt.Errorf("error when creating log for instance '%s'. Cause: %w", instanceName, err)
	}
	err = uc.AccessPermissionStorage.SaveLog(grantLog)
	if err != nil {
		return fmt.Errorf("error when saving log for instance '%s'. Cause: %w", instanceName, err)
	}
	return nil
}

func (uc *ChangeStatusDatabaseInstanceUseCase) buildOutputDTO(dbInstance *entity.DatabaseInstance) *dto.ChangeStatusOutputDTO {
	var disabledAt *time.Time
	if dbInstance.DisabledAt.Valid {
		disabledAt = &dbInstance.DisabledAt.Time
	}
	return &dto.ChangeStatusOutputDTO{
		ID:         dbInstance.ID.String(),
		Enabled:    dbInstance.Enabled,
		UpdatedAt:  dbInstance.UpdatedAt,
		DisabledAt: disabledAt,
	}
}

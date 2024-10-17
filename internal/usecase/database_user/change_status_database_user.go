package dbuser

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/zgsolucoes/zg-data-guard/internal/usecase/common"

	"github.com/zgsolucoes/zg-data-guard/internal/database/storage"
	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	"github.com/zgsolucoes/zg-data-guard/internal/entity"
)

var ErrCouldNotRevokeAllAccess = errors.New("could not revoke all access from user")

type ChangeStatusDatabaseUserUseCase struct {
	DatabaseUserStorage storage.DatabaseUserStorage
	RevokeAccessUseCase common.RevokeAccessPermissionUseCaseInterface
}

func NewChangeStatusDatabaseUserUseCase(
	databaseUserStorage storage.DatabaseUserStorage,
	revokeAccessUseCase common.RevokeAccessPermissionUseCaseInterface) *ChangeStatusDatabaseUserUseCase {
	return &ChangeStatusDatabaseUserUseCase{DatabaseUserStorage: databaseUserStorage, RevokeAccessUseCase: revokeAccessUseCase}
}

func (uc *ChangeStatusDatabaseUserUseCase) Execute(userID string, enabled bool, operationUserID string) (*dto.ChangeStatusOutputDTO, error) {
	dbUser, err := uc.DatabaseUserStorage.FindByID(userID)
	if err != nil {
		return nil, common.HandleFindError(err, common.ErrDatabaseUserNotFound)
	}
	log.Printf("Changing status of database user '%s' to '%t'. Requester: %s", dbUser.ID.String(), enabled, operationUserID)
	if err := uc.changeUserStatus(dbUser, enabled, operationUserID); err != nil {
		return nil, err
	}
	if err := uc.DatabaseUserStorage.Update(dbUser); err != nil {
		return nil, fmt.Errorf("error when updating database user %s to new status '%t'. Cause: %w", dbUser.ID.String(), enabled, err)
	}
	log.Printf("Database user '%s' status changed to '%t' successfully. Requester: %s", dbUser.ID.String(), enabled, operationUserID)
	return uc.buildOutputDTO(dbUser), nil
}

func (uc *ChangeStatusDatabaseUserUseCase) changeUserStatus(dbUser *entity.DatabaseUser, enabled bool, operationUserID string) error {
	if enabled {
		dbUser.Enable()
		return nil
	}
	if err := uc.revokeUserAccess(dbUser, operationUserID); err != nil {
		return err
	}
	dbUser.Disable()
	return nil
}

func (uc *ChangeStatusDatabaseUserUseCase) revokeUserAccess(dbUser *entity.DatabaseUser, operationUserID string) error {
	revokeResult, err := uc.RevokeAccessUseCase.Execute(dto.RevokeAccessInputDTO{DatabaseInstancesIDs: []string{}, DatabaseUserID: dbUser.ID.String()}, operationUserID)
	if err != nil {
		return handleRevokeAccessError(err, dbUser.ID.String())
	}
	if revokeResult.HasErrors {
		return ErrCouldNotRevokeAllAccess
	}
	return nil
}

func (uc *ChangeStatusDatabaseUserUseCase) buildOutputDTO(dbUser *entity.DatabaseUser) *dto.ChangeStatusOutputDTO {
	var disabledAt *time.Time
	if dbUser.DisabledAt.Valid {
		disabledAt = &dbUser.DisabledAt.Time
	}
	return &dto.ChangeStatusOutputDTO{
		ID:         dbUser.ID.String(),
		Enabled:    dbUser.Enabled,
		UpdatedAt:  dbUser.UpdatedAt,
		DisabledAt: disabledAt,
	}
}

func handleRevokeAccessError(err error, dbUserID string) error {
	if errors.Is(err, common.ErrNoAccessibleInstancesFound) {
		log.Printf("No accessible instances found for user %s. Skipping access revocation.", dbUserID)
		return nil
	}
	return fmt.Errorf("error revoking access from instances of user %s. Cause: %w", dbUserID, err)
}

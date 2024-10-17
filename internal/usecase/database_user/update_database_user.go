package dbuser

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/zgsolucoes/zg-data-guard/internal/usecase/common"

	"github.com/zgsolucoes/zg-data-guard/internal/database/storage"
	"github.com/zgsolucoes/zg-data-guard/internal/dto"
)

const errorUpdatingDatabaseUser = "Error updating database user"

var (
	ErrDatabaseUserHasAccessPermissions = errors.New("database user has access permissions and cannot have their role changed")
)

type UpdateDatabaseUserUseCase struct {
	DatabaseUserStorage     storage.DatabaseUserStorage
	DatabaseRoleStorage     storage.DatabaseRoleStorage
	AccessPermissionStorage storage.AccessPermissionStorage
}

func NewUpdateDatabaseUserUseCase(
	databaseUserStorage storage.DatabaseUserStorage,
	databaseRoleStorage storage.DatabaseRoleStorage,
	accessPermissionStorage storage.AccessPermissionStorage,
) *UpdateDatabaseUserUseCase {
	return &UpdateDatabaseUserUseCase{
		DatabaseUserStorage:     databaseUserStorage,
		DatabaseRoleStorage:     databaseRoleStorage,
		AccessPermissionStorage: accessPermissionStorage,
	}
}

func (uc *UpdateDatabaseUserUseCase) Execute(input dto.UpdateDatabaseUserInputDTO, dbUserID, operationUserID string) (*dto.DatabaseUserOutputDTO, error) {
	dbUser, err := uc.DatabaseUserStorage.FindByID(dbUserID)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		logErrorWithID(common.ErrDatabaseUserNotFound, errorUpdatingDatabaseUser, dbUserID)
		return nil, common.ErrDatabaseUserNotFound
	}
	if err != nil {
		logErrorWithID(err, errorUpdatingDatabaseUser, dbUserID)
		return nil, err
	}
	if input.DatabaseRoleID != dbUser.DatabaseRoleID {
		hasAccessPermissions, errCheckingPermission := uc.AccessPermissionStorage.CheckIfUserHasAccessPermission(dbUserID)
		if errCheckingPermission != nil {
			logErrorWithID(errCheckingPermission, errorUpdatingDatabaseUser, dbUserID)
			return nil, errCheckingPermission
		}
		if hasAccessPermissions {
			logErrorWithID(ErrDatabaseUserHasAccessPermissions, errorUpdatingDatabaseUser, dbUserID)
			return nil, ErrDatabaseUserHasAccessPermissions
		}
	}

	err = validateDatabaseRoleExistence(input.DatabaseRoleID, uc.DatabaseRoleStorage, errorUpdatingDatabaseUser)
	if err != nil {
		return nil, err
	}

	err = dbUser.Update(input.Name, input.DatabaseRoleID, input.Team, input.Position)
	if err != nil {
		logErrorWithID(err, errorUpdatingDatabaseUser, dbUserID)
		return nil, err
	}
	err = uc.DatabaseUserStorage.Update(dbUser)
	if err != nil {
		logErrorWithID(err, errorUpdatingDatabaseUser, dbUserID)
		return nil, err
	}

	log.Printf("Database user %v updated successfully by user %s!", dbUser.ID, operationUserID)
	isDisabled := dbUser.DisabledAt.Valid
	var disabledAt *time.Time
	if isDisabled {
		disabledAt = &dbUser.DisabledAt.Time
	}
	return &dto.DatabaseUserOutputDTO{
		ID:              dbUser.ID.String(),
		Name:            dbUser.Name,
		Email:           dbUser.Email,
		Username:        dbUser.Username,
		Team:            dbUser.Team,
		Position:        dbUser.Position,
		DatabaseRoleID:  dbUser.DatabaseRoleID,
		Enabled:         dbUser.Enabled,
		CreatedByUserID: dbUser.CreatedByUserID,
		CreatedAt:       dbUser.CreatedAt,
		UpdatedAt:       &dbUser.UpdatedAt,
		DisabledAt:      disabledAt,
	}, nil
}

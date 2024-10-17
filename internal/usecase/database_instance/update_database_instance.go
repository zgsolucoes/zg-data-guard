package instance

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/zgsolucoes/zg-data-guard/internal/database/storage"
	"github.com/zgsolucoes/zg-data-guard/internal/dto"
)

const errorUpdatingDatabaseInstance = "Error updating database instance"

type UpdateDatabaseInstanceUseCase struct {
	DatabaseInstanceStorage storage.DatabaseInstanceStorage
	EcosystemStorage        storage.EcosystemStorage
	TechnologyStorage       storage.DatabaseTechnologyStorage
}

func NewUpdateDatabaseInstanceUseCase(
	dbInstanceStorage storage.DatabaseInstanceStorage,
	ecosystemStorage storage.EcosystemStorage,
	technologyStorage storage.DatabaseTechnologyStorage,
) *UpdateDatabaseInstanceUseCase {
	return &UpdateDatabaseInstanceUseCase{
		DatabaseInstanceStorage: dbInstanceStorage,
		EcosystemStorage:        ecosystemStorage,
		TechnologyStorage:       technologyStorage,
	}
}

func (uc *UpdateDatabaseInstanceUseCase) Execute(input dto.DatabaseInstanceInputDTO, dbInstanceID, operationUserID string) (*dto.DatabaseInstanceOutputDTO, error) {
	dbInstance, err := uc.DatabaseInstanceStorage.FindByID(dbInstanceID)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		logErrorWithID(ErrDatabaseInstanceNotFound, errorUpdatingDatabaseInstance, dbInstanceID)
		return nil, ErrDatabaseInstanceNotFound
	}
	if err != nil {
		logErrorWithID(err, errorUpdatingDatabaseInstance, dbInstanceID)
		return nil, err
	}

	if dbInstance.HostConnection.Host != input.Host || dbInstance.HostConnection.Port != input.Port {
		dbInstanceExists, err := uc.DatabaseInstanceStorage.Exists(input.Host, input.Port)
		if err != nil {
			log.Printf("Error while checking existance of database instance with host %s and port %s. Cause: %v", input.Host, input.Port, err.Error())
			return nil, err
		}
		if dbInstanceExists {
			logError(ErrHostAlreadyExists, errorUpdatingDatabaseInstance)
			return nil, ErrHostAlreadyExists
		}
	}

	err = validateEcosystemAndTechnologyExistence(input, uc.EcosystemStorage, uc.TechnologyStorage, errorUpdatingDatabaseInstance)
	if err != nil {
		return nil, err
	}

	err = dbInstance.Update(input)
	if err != nil {
		logErrorWithID(err, errorUpdatingDatabaseInstance, dbInstanceID)
		return nil, err
	}
	err = uc.DatabaseInstanceStorage.UpdateWithHostInfo(dbInstance)
	if err != nil {
		logErrorWithID(err, errorUpdatingDatabaseInstance, dbInstanceID)
		return nil, err
	}

	isDisabled := dbInstance.DisabledAt.Valid
	var disabledAt *time.Time
	if isDisabled {
		disabledAt = &dbInstance.DisabledAt.Time
	}
	log.Printf("Database instance %v updated successfully by user %s!", dbInstance.ID, operationUserID)
	return &dto.DatabaseInstanceOutputDTO{
		ID:                   dbInstance.ID.String(),
		Name:                 dbInstance.Name,
		Host:                 dbInstance.HostConnection.Host,
		Port:                 dbInstance.HostConnection.Port,
		HostConnection:       dbInstance.HostConnection.HostConnection,
		PortConnection:       dbInstance.HostConnection.PortConnection,
		AdminUser:            dbInstance.HostConnection.AdminUser,
		EcosystemID:          dbInstance.EcosystemID,
		DatabaseTechnologyID: dbInstance.DatabaseTechnologyID,
		Enabled:              dbInstance.Enabled,
		Note:                 dbInstance.Note,
		CreatedByUserID:      dbInstance.CreatedByUserID,
		CreatedAt:            dbInstance.CreatedAt,
		UpdatedAt:            &dbInstance.UpdatedAt,
		DisabledAt:           disabledAt,
	}, nil
}

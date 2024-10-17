package instance

import (
	"database/sql"
	"errors"
	"log"

	"github.com/zgsolucoes/zg-data-guard/internal/database/storage"
	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	"github.com/zgsolucoes/zg-data-guard/internal/entity"
	"github.com/zgsolucoes/zg-data-guard/internal/usecase/common"
)

const errorCreatingDatabaseInstance = "Error creating database instance"

var (
	ErrHostAlreadyExists = errors.New("this host and port already exists")
)

type CreateDatabaseInstanceUseCase struct {
	DatabaseInstanceStorage storage.DatabaseInstanceStorage
	EcosystemStorage        storage.EcosystemStorage
	TechnologyStorage       storage.DatabaseTechnologyStorage
}

func NewCreateDatabaseInstanceUseCase(
	databaseInstanceStorage storage.DatabaseInstanceStorage,
	ecosystemStorage storage.EcosystemStorage,
	technologyStorage storage.DatabaseTechnologyStorage,
) *CreateDatabaseInstanceUseCase {
	return &CreateDatabaseInstanceUseCase{
		DatabaseInstanceStorage: databaseInstanceStorage,
		EcosystemStorage:        ecosystemStorage,
		TechnologyStorage:       technologyStorage,
	}
}

func (c *CreateDatabaseInstanceUseCase) Execute(input dto.DatabaseInstanceInputDTO, createdByUserID string) (*dto.DatabaseInstanceOutputDTO, error) {
	databaseInstance, err := entity.NewDatabaseInstance(input, createdByUserID)
	if err != nil {
		logError(err, errorCreatingDatabaseInstance)
		return nil, err
	}

	exists, err := c.DatabaseInstanceStorage.Exists(input.Host, input.Port)
	if err != nil {
		logError(err, errorCreatingDatabaseInstance)
		return nil, err
	}
	if exists {
		logError(ErrHostAlreadyExists, errorCreatingDatabaseInstance)
		return nil, ErrHostAlreadyExists
	}

	err = validateEcosystemAndTechnologyExistence(input, c.EcosystemStorage, c.TechnologyStorage, errorCreatingDatabaseInstance)
	if err != nil {
		return nil, err
	}

	err = c.DatabaseInstanceStorage.Save(databaseInstance)
	if err != nil {
		logError(err, errorCreatingDatabaseInstance)
		return nil, err
	}

	log.Printf("Database instance %v created successfully by user %s!", databaseInstance.ID, createdByUserID)
	return &dto.DatabaseInstanceOutputDTO{
		ID:                   databaseInstance.ID.String(),
		Name:                 databaseInstance.Name,
		Host:                 databaseInstance.HostConnection.Host,
		Port:                 databaseInstance.HostConnection.Port,
		HostConnection:       databaseInstance.HostConnection.HostConnection,
		PortConnection:       databaseInstance.HostConnection.PortConnection,
		AdminUser:            databaseInstance.HostConnection.AdminUser,
		EcosystemID:          databaseInstance.EcosystemID,
		DatabaseTechnologyID: databaseInstance.DatabaseTechnologyID,
		Enabled:              databaseInstance.Enabled,
		Note:                 databaseInstance.Note,
		CreatedByUserID:      databaseInstance.CreatedByUserID,
		CreatedAt:            databaseInstance.CreatedAt,
	}, nil
}

func validateEcosystemAndTechnologyExistence(
	input dto.DatabaseInstanceInputDTO,
	ecoStorage storage.EcosystemStorage,
	techStorage storage.DatabaseTechnologyStorage,
	operation string,
) error {
	_, err := ecoStorage.FindByID(input.EcosystemID)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		logError(common.ErrEcosystemNotFound, operation)
		return common.ErrEcosystemNotFound
	}
	if err != nil {
		logError(err, operation)
		return err
	}
	log.Printf("Ecosystem with id %s loaded successfully!", input.EcosystemID)

	_, err = techStorage.FindByID(input.DatabaseTechnologyID)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		logError(common.ErrTechnologyNotFound, operation)
		return common.ErrTechnologyNotFound
	}
	if err != nil {
		logError(err, operation)
		return err
	}
	log.Printf("Database technology with id %s loaded successfully!", input.DatabaseTechnologyID)
	return nil
}

func logError(err error, operationError string) {
	log.Printf("%s. Cause: %v", operationError, err.Error())
}

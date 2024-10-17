package instance

import (
	"database/sql"
	"errors"
	"log"

	"github.com/zgsolucoes/zg-data-guard/config"
	"github.com/zgsolucoes/zg-data-guard/internal/database/storage"
	"github.com/zgsolucoes/zg-data-guard/internal/dto"
)

const errorFetchingDatabaseInstance = "Error fetching database instance"

var (
	ErrDatabaseInstanceNotFound = errors.New("database instance not found")
)

type GetDatabaseInstanceUseCase struct {
	DatabaseInstanceStorage storage.DatabaseInstanceStorage
}

func NewGetDatabaseInstanceUseCase(dbInstanceStorage storage.DatabaseInstanceStorage) *GetDatabaseInstanceUseCase {
	return &GetDatabaseInstanceUseCase{DatabaseInstanceStorage: dbInstanceStorage}
}

func (uc *GetDatabaseInstanceUseCase) Execute(dbInstanceID string) (*dto.DatabaseInstanceOutputDTO, error) {
	dbDTO, err := uc.DatabaseInstanceStorage.FindDTOByID(dbInstanceID)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		logErrorWithID(ErrDatabaseInstanceNotFound, errorFetchingDatabaseInstance, dbInstanceID)
		return nil, ErrDatabaseInstanceNotFound
	}
	if err != nil {
		logErrorWithID(err, errorFetchingDatabaseInstance, dbInstanceID)
		return nil, err
	}
	dbDTO.AdminUser = ""
	dbDTO.AdminPassword = ""
	log.Printf("DatabaseInstance with id %s loaded successfully!", dbInstanceID)
	return dbDTO, nil
}

func (uc *GetDatabaseInstanceUseCase) FetchCredentials(dbInstanceID, userID string) (*dto.DatabaseInstanceCredentialsOutputDTO, error) {
	instanceDTO, err := uc.DatabaseInstanceStorage.FindDTOByID(dbInstanceID)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		logErrorWithID(ErrDatabaseInstanceNotFound, errorFetchingDatabaseInstance, dbInstanceID)
		return nil, ErrDatabaseInstanceNotFound
	}
	if err != nil {
		logErrorWithID(err, errorFetchingDatabaseInstance, dbInstanceID)
		return nil, err
	}
	log.Printf("Credentials of database instance with id %s - %s accessed by user with id %s", dbInstanceID, instanceDTO.Name, userID)
	plaintTextPasswd, err := config.GetCryptoHelper().Decrypt(instanceDTO.AdminPassword)
	if err != nil {
		logErrorWithID(err, errorFetchingDatabaseInstance, dbInstanceID)
		return nil, err
	}
	return &dto.DatabaseInstanceCredentialsOutputDTO{
		User:     instanceDTO.AdminUser,
		Password: plaintTextPasswd,
	}, nil
}

func logErrorWithID(err error, operationError, id string) {
	log.Printf("%s with id %s. Cause: %v", operationError, id, err.Error())
}

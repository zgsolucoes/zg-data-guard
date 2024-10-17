package dbuser

import (
	"database/sql"
	"errors"
	"log"

	"github.com/zgsolucoes/zg-data-guard/internal/usecase/common"

	"github.com/zgsolucoes/zg-data-guard/internal/database/storage"
	"github.com/zgsolucoes/zg-data-guard/internal/dto"
)

const errorFetchingDatabaseUser = "Error fetching database user"

type GetDatabaseUserUseCase struct {
	DatabaseUserStorage storage.DatabaseUserStorage
}

func NewGetDatabaseUserUseCase(dbUserStorage storage.DatabaseUserStorage) *GetDatabaseUserUseCase {
	return &GetDatabaseUserUseCase{DatabaseUserStorage: dbUserStorage}
}

func (uc *GetDatabaseUserUseCase) Execute(dbUserID string) (*dto.DatabaseUserOutputDTO, error) {
	dbUserDTO, err := uc.DatabaseUserStorage.FindDTOByID(dbUserID)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		logErrorWithID(common.ErrDatabaseUserNotFound, errorFetchingDatabaseUser, dbUserID)
		return nil, common.ErrDatabaseUserNotFound
	}
	if err != nil {
		logErrorWithID(err, errorFetchingDatabaseUser, dbUserID)
		return nil, err
	}
	dbUserDTO.Password = ""
	log.Printf("Database user with id %s loaded successfully!", dbUserID)
	return dbUserDTO, nil
}

func (uc *GetDatabaseUserUseCase) FetchCredentials(dbUserID, userID string) (*dto.DatabaseUserCredentialsOutputDTO, error) {
	dbUser, err := uc.DatabaseUserStorage.FindByID(dbUserID)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		logErrorWithID(common.ErrDatabaseUserNotFound, errorFetchingDatabaseUser, dbUserID)
		return nil, common.ErrDatabaseUserNotFound
	}
	if err != nil {
		logErrorWithID(err, errorFetchingDatabaseUser, dbUserID)
		return nil, err
	}
	log.Printf("Credentials of database user with id %s - %s accessed by user with id %s", dbUserID, dbUser.Name, userID)
	err = dbUser.DecryptPassword()
	if err != nil {
		logErrorWithID(err, errorFetchingDatabaseUser, dbUserID)
		return nil, err
	}
	return &dto.DatabaseUserCredentialsOutputDTO{
		User:     dbUser.Username,
		Password: dbUser.Password,
	}, nil
}

func logErrorWithID(err error, operationError, id string) {
	log.Printf("%s with id %s. Cause: %v", operationError, id, err.Error())
}

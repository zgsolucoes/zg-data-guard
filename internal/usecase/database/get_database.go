package database

import (
	"database/sql"
	"errors"
	"log"

	"github.com/zgsolucoes/zg-data-guard/internal/database/storage"
	"github.com/zgsolucoes/zg-data-guard/internal/dto"
)

const errorFetchingDatabase = "Error fetching database"

var (
	ErrDatabaseNotFound = errors.New("database not found")
)

type GetDatabaseUseCase struct {
	DatabaseStorage storage.DatabaseStorage
}

func NewGetDatabaseUseCase(databaseStorage storage.DatabaseStorage) *GetDatabaseUseCase {
	return &GetDatabaseUseCase{DatabaseStorage: databaseStorage}
}

func (uc *GetDatabaseUseCase) Execute(databaseID string) (*dto.DatabaseOutputDTO, error) {
	dbDTO, err := uc.DatabaseStorage.FindDTOByID(databaseID)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		logErrorWithID(ErrDatabaseNotFound, errorFetchingDatabase, databaseID)
		return nil, ErrDatabaseNotFound
	}
	if err != nil {
		logErrorWithID(err, errorFetchingDatabase, databaseID)
		return nil, err
	}
	log.Printf("Database with id %s loaded successfully!", databaseID)
	return dbDTO, nil
}
func logErrorWithID(err error, operationError, id string) {
	log.Printf("%s with id %s. Cause: %v", operationError, id, err.Error())
}

package database

import (
	"log"

	"github.com/zgsolucoes/zg-data-guard/internal/database/storage"
	"github.com/zgsolucoes/zg-data-guard/internal/dto"
)

type ListDatabasesUseCase struct {
	DatabaseStorage storage.DatabaseStorage
}

func NewListDatabasesUseCase(databaseStorage storage.DatabaseStorage) *ListDatabasesUseCase {
	return &ListDatabasesUseCase{
		DatabaseStorage: databaseStorage,
	}
}

func (uc *ListDatabasesUseCase) Execute(ecosystemID, databaseInstanceID string) ([]*dto.DatabaseOutputDTO, error) {
	databaseDTOs, err := uc.DatabaseStorage.FindAllDTOs(ecosystemID, databaseInstanceID)
	if err != nil {
		log.Printf("Error fetching databases! Cause: %v", err.Error())
		return nil, err
	}
	log.Printf("List of databases loaded successfully!")
	return databaseDTOs, nil
}

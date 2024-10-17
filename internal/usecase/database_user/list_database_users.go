package dbuser

import (
	"log"

	"github.com/zgsolucoes/zg-data-guard/internal/database/storage"
	"github.com/zgsolucoes/zg-data-guard/internal/dto"
)

type ListDatabaseUsersUseCase struct {
	DatabaseUserStorage storage.DatabaseUserStorage
}

func NewListDatabaseUsersUseCase(databaseUserStorage storage.DatabaseUserStorage) *ListDatabaseUsersUseCase {
	return &ListDatabaseUsersUseCase{
		DatabaseUserStorage: databaseUserStorage,
	}
}

func (uc *ListDatabaseUsersUseCase) Execute(onlyEnabled bool) ([]*dto.DatabaseUserOutputDTO, error) {
	var dbUsersDTO []*dto.DatabaseUserOutputDTO
	var err error
	if onlyEnabled {
		dbUsersDTO, err = uc.DatabaseUserStorage.FindAllDTOsEnabled()
	} else {
		dbUsersDTO, err = uc.DatabaseUserStorage.FindAllDTOs([]string{})
	}
	// Hiding passwords from the output
	for _, dbUser := range dbUsersDTO {
		dbUser.Password = ""
	}
	if err != nil {
		log.Printf("Error fetching database users! Cause: %v", err.Error())
		return nil, err
	}
	log.Printf("List of database users loaded successfully!")
	return dbUsersDTO, nil
}

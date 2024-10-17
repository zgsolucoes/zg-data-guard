package instance

import (
	"log"

	"github.com/zgsolucoes/zg-data-guard/internal/database/storage"
	"github.com/zgsolucoes/zg-data-guard/internal/dto"
)

type ListDatabaseInstancesUseCase struct {
	DatabaseInstanceStorage storage.DatabaseInstanceStorage
}

func NewListDatabaseInstancesUseCase(databaseInstanceStorage storage.DatabaseInstanceStorage) *ListDatabaseInstancesUseCase {
	return &ListDatabaseInstancesUseCase{
		DatabaseInstanceStorage: databaseInstanceStorage,
	}
}

func (uc *ListDatabaseInstancesUseCase) Execute(ecosystemID, technologyID string, onlyEnabled bool) ([]*dto.DatabaseInstanceOutputDTO, error) {
	var dbInstancesDTO []*dto.DatabaseInstanceOutputDTO
	var err error
	if onlyEnabled {
		dbInstancesDTO, err = uc.DatabaseInstanceStorage.FindAllDTOsEnabled(ecosystemID, technologyID)
	} else {
		dbInstancesDTO, err = uc.DatabaseInstanceStorage.FindAllDTOs(ecosystemID, technologyID, []string{})
	}
	if err != nil {
		log.Printf("Error fetching database instances! Cause: %v", err.Error())
		return nil, err
	}
	for _, dbInstance := range dbInstancesDTO {
		dbInstance.AdminUser = ""
		dbInstance.AdminPassword = ""
	}
	log.Printf("List of database instances loaded successfully!")
	return dbInstancesDTO, nil
}

package ecosystem

import (
	"log"

	"github.com/zgsolucoes/zg-data-guard/internal/database/storage"
	"github.com/zgsolucoes/zg-data-guard/internal/dto"
)

type ListEcosystemsUseCase struct {
	EcosystemStorage storage.EcosystemStorage
}

func NewListEcosystemsUseCase(ecosystemStorage storage.EcosystemStorage) *ListEcosystemsUseCase {
	return &ListEcosystemsUseCase{
		EcosystemStorage: ecosystemStorage,
	}
}

func (uc *ListEcosystemsUseCase) Execute(page, limit int) ([]*dto.EcosystemOutputDTO, error) {
	ecosystems, err := uc.EcosystemStorage.FindAll(page, limit)
	if err != nil {
		log.Printf("Error fetching ecosystems! Cause: %v", err.Error())
		return nil, err
	}
	log.Printf("All ecosystems from page %d and limit %d loaded successfully!", page, limit)
	return ecosystems, nil
}

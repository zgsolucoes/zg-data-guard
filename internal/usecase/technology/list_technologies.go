package tech

import (
	"log"

	"github.com/zgsolucoes/zg-data-guard/internal/database/storage"
	"github.com/zgsolucoes/zg-data-guard/internal/dto"
)

type ListTechnologiesUseCase struct {
	TechnologyStorage storage.DatabaseTechnologyStorage
}

func NewListTechnologiesUseCase(technologyStorage storage.DatabaseTechnologyStorage) *ListTechnologiesUseCase {
	return &ListTechnologiesUseCase{
		TechnologyStorage: technologyStorage,
	}
}

func (uc *ListTechnologiesUseCase) Execute(page, limit int) ([]*dto.TechnologyOutputDTO, error) {
	technologies, err := uc.TechnologyStorage.FindAll(page, limit)
	if err != nil {
		log.Printf("Error fetching technologies! Cause: %v", err.Error())
		return nil, err
	}
	log.Printf("All technologies from page %d and limit %d loaded successfully!", page, limit)
	return technologies, nil
}

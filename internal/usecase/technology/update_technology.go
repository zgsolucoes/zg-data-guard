package tech

import (
	"database/sql"
	"errors"
	"log"

	"github.com/zgsolucoes/zg-data-guard/internal/database/storage"
	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	"github.com/zgsolucoes/zg-data-guard/internal/usecase/common"
)

const errorUpdatingTechnology = "Error updating technology"

type UpdateTechnologyUseCase struct {
	TechnologyStorage storage.DatabaseTechnologyStorage
}

func NewUpdateTechnologyUseCase(technologyStorage storage.DatabaseTechnologyStorage) *UpdateTechnologyUseCase {
	return &UpdateTechnologyUseCase{
		TechnologyStorage: technologyStorage,
	}
}

func (uc *UpdateTechnologyUseCase) Execute(input dto.TechnologyInputDTO, technologyID, operationUserID string) (*dto.TechnologyOutputDTO, error) {
	technology, err := uc.TechnologyStorage.FindByID(technologyID)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		logErrorWithID(common.ErrTechnologyNotFound, errorUpdatingTechnology, technologyID)
		return nil, common.ErrTechnologyNotFound
	}
	if err != nil {
		logErrorWithID(err, errorUpdatingTechnology, technologyID)
		return nil, err
	}

	if technology.Name != input.Name || technology.Version != input.Version {
		technologyExists, err := uc.TechnologyStorage.Exists(input.Name, input.Version)
		if err != nil {
			log.Printf("Error while checking existance of technology with name %s and version %s. Cause: %v", input.Name, input.Version, err.Error())
			return nil, err
		}
		if technologyExists {
			logError(ErrTechnologyAlreadyExists, errorUpdatingTechnology)
			return nil, ErrTechnologyAlreadyExists
		}
	}

	technology.Update(input.Name, input.Version)
	err = uc.TechnologyStorage.Update(technology)
	if err != nil {
		logErrorWithID(err, errorUpdatingTechnology, technologyID)
		return nil, err
	}

	log.Printf("Technology %v updated successfully by user %s!", technology.ID, operationUserID)
	return &dto.TechnologyOutputDTO{
		ID:        technology.ID.String(),
		Name:      technology.Name,
		Version:   technology.Version,
		CreatedAt: technology.CreatedAt,
		UpdatedAt: &technology.UpdatedAt,
	}, nil
}

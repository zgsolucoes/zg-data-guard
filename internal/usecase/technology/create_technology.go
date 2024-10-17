package tech

import (
	"errors"
	"log"

	"github.com/zgsolucoes/zg-data-guard/internal/database/storage"
	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	"github.com/zgsolucoes/zg-data-guard/internal/entity"
)

const errorCreatingTechnology = "Error creating technology"

var (
	ErrTechnologyAlreadyExists = errors.New("database technology already exists with this name and version")
)

type CreateTechnologyUseCase struct {
	TechnologyStorage storage.DatabaseTechnologyStorage
}

func NewCreateTechnologyUseCase(technologyStorage storage.DatabaseTechnologyStorage) *CreateTechnologyUseCase {
	return &CreateTechnologyUseCase{
		TechnologyStorage: technologyStorage,
	}
}

func (c *CreateTechnologyUseCase) Execute(input dto.TechnologyInputDTO, createdByUserID string) (*dto.TechnologyOutputDTO, error) {
	technology, err := entity.NewDatabaseTechnology(input.Name, input.Version, createdByUserID)
	if err != nil {
		logError(err, errorCreatingTechnology)
		return nil, err
	}

	exists, err := c.TechnologyStorage.Exists(input.Name, input.Version)
	if err != nil {
		logError(err, errorCreatingTechnology)
		return nil, err
	}
	if exists {
		logError(ErrTechnologyAlreadyExists, errorCreatingTechnology)
		return nil, ErrTechnologyAlreadyExists
	}

	err = c.TechnologyStorage.Save(technology)
	if err != nil {
		logError(err, errorCreatingTechnology)
		return nil, err
	}

	log.Printf("Database technology %v created successfully by user %s!", technology.ID, createdByUserID)
	return &dto.TechnologyOutputDTO{
		ID:              technology.ID.String(),
		Name:            technology.Name,
		Version:         technology.Version,
		CreatedByUserID: technology.CreatedByUserID,
		CreatedAt:       technology.CreatedAt,
	}, nil
}

func logError(err error, operationError string) {
	log.Printf("%s. Cause: %v", operationError, err.Error())
}

package tech

import (
	"database/sql"
	"errors"
	"log"

	"github.com/zgsolucoes/zg-data-guard/internal/database/storage"
	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	"github.com/zgsolucoes/zg-data-guard/internal/usecase/common"
)

const errorFetchingTechnology = "Error fetching technology"

type GetTechnologyUseCase struct {
	TechnologyStorage storage.DatabaseTechnologyStorage
	UserStorage       storage.ApplicationUserStorage
}

func NewGetTechnologyUseCase(
	technologyStorage storage.DatabaseTechnologyStorage,
	userStorage storage.ApplicationUserStorage,
) *GetTechnologyUseCase {
	return &GetTechnologyUseCase{
		TechnologyStorage: technologyStorage,
		UserStorage:       userStorage,
	}
}

func (uc *GetTechnologyUseCase) Execute(technologyID string) (*dto.TechnologyOutputDTO, error) {
	technology, err := uc.TechnologyStorage.FindByID(technologyID)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		logErrorWithID(common.ErrTechnologyNotFound, errorFetchingTechnology, technologyID)
		return nil, common.ErrTechnologyNotFound
	}
	if err != nil {
		logErrorWithID(err, errorFetchingTechnology, technologyID)
		return nil, err
	}
	user, err := uc.UserStorage.FindByID(technology.CreatedByUserID)
	if err != nil {
		log.Printf("Error fetching user with id %s. Cause: %v", technology.CreatedByUserID, err.Error())
		return nil, err
	}
	log.Printf("Technology with id %s loaded successfully!", technologyID)
	return &dto.TechnologyOutputDTO{
		ID:            technology.ID.String(),
		Name:          technology.Name,
		Version:       technology.Version,
		CreatedAt:     technology.CreatedAt,
		CreatedByUser: user.Name,
		UpdatedAt:     &technology.UpdatedAt,
	}, nil
}

func logErrorWithID(err error, operationError, id string) {
	log.Printf("%s with id %s. Cause: %v", operationError, id, err.Error())
}

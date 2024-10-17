package tech

import (
	"database/sql"
	"errors"
	"log"

	"github.com/zgsolucoes/zg-data-guard/internal/database/storage"
	"github.com/zgsolucoes/zg-data-guard/internal/usecase/common"
)

const errorDeletingTechnology = "Error deleting technology"

type DeleteTechnologyUseCase struct {
	TechnologyStorage storage.DatabaseTechnologyStorage
}

func NewDeleteTechnologyUseCase(technologyStorage storage.DatabaseTechnologyStorage) *DeleteTechnologyUseCase {
	return &DeleteTechnologyUseCase{
		TechnologyStorage: technologyStorage,
	}
}

func (uc *DeleteTechnologyUseCase) Execute(technologyID string, operationUserID string) error {
	err := uc.TechnologyStorage.Delete(technologyID)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		logErrorWithID(common.ErrTechnologyNotFound, errorDeletingTechnology, technologyID)
		return common.ErrTechnologyNotFound
	}
	if err != nil {
		logErrorWithID(err, errorDeletingTechnology, technologyID)
		return err
	}
	log.Printf("Technology %s deleted successfully by user %s!", technologyID, operationUserID)
	return nil
}

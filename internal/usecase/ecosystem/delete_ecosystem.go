package ecosystem

import (
	"database/sql"
	"errors"
	"log"

	"github.com/zgsolucoes/zg-data-guard/internal/database/storage"
	"github.com/zgsolucoes/zg-data-guard/internal/usecase/common"
)

type DeleteEcosystemUseCase struct {
	EcosystemStorage storage.EcosystemStorage
}

func NewDeleteEcosystemUseCase(ecosystemStorage storage.EcosystemStorage) *DeleteEcosystemUseCase {
	return &DeleteEcosystemUseCase{
		EcosystemStorage: ecosystemStorage,
	}
}

func (uc *DeleteEcosystemUseCase) Execute(ecosystemID string, operationUserID string) error {
	err := uc.EcosystemStorage.Delete(ecosystemID)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		log.Printf("Ecosystem with id %s not found in database!", ecosystemID)
		return common.ErrEcosystemNotFound
	}
	if err != nil {
		log.Printf("Error deleting ecosystem with id %s. Cause: %v", ecosystemID, err.Error())
		return err
	}
	log.Printf("Ecosystem %s deleted successfully by user %s!", ecosystemID, operationUserID)
	return nil
}

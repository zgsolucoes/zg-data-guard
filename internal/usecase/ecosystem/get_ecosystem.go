package ecosystem

import (
	"database/sql"
	"errors"
	"log"

	"github.com/zgsolucoes/zg-data-guard/internal/database/storage"
	"github.com/zgsolucoes/zg-data-guard/internal/dto"
)

const errorFetchingEcosystem = "Error fetching ecosystem"

var (
	ErrEcosystemNotFound = errors.New("ecosystem not found")
)

type GetEcosystemUseCase struct {
	EcosystemStorage storage.EcosystemStorage
	UserStorage      storage.ApplicationUserStorage
}

func NewGetEcosystemUseCase(
	ecosystemStorage storage.EcosystemStorage,
	userStorage storage.ApplicationUserStorage,
) *GetEcosystemUseCase {
	return &GetEcosystemUseCase{
		EcosystemStorage: ecosystemStorage,
		UserStorage:      userStorage,
	}
}

func (uc *GetEcosystemUseCase) Execute(ecosystemID string) (*dto.EcosystemOutputDTO, error) {
	ecosystem, err := uc.EcosystemStorage.FindByID(ecosystemID)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		logErrorWithID(err, errorFetchingEcosystem, ecosystemID)
		return nil, ErrEcosystemNotFound
	}
	if err != nil {
		logErrorWithID(err, errorFetchingEcosystem, ecosystemID)
		return nil, err
	}
	user, err := uc.UserStorage.FindByID(ecosystem.CreatedByUserID)
	if err != nil {
		log.Printf("Error fetching user with id %s. Cause: %v", ecosystem.CreatedByUserID, err.Error())
		return nil, err
	}
	log.Printf("Ecosystem with id %s loaded successfully!", ecosystemID)
	return &dto.EcosystemOutputDTO{
		ID:            ecosystem.ID.String(),
		Code:          ecosystem.Code,
		DisplayName:   ecosystem.DisplayName,
		CreatedAt:     ecosystem.CreatedAt,
		CreatedByUser: user.Name,
		UpdatedAt:     &ecosystem.UpdatedAt,
	}, nil
}

func logErrorWithID(err error, operationError, id string) {
	log.Printf("%s with id %s. Cause: %v", operationError, id, err.Error())
}

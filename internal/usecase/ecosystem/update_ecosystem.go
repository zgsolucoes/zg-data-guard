package ecosystem

import (
	"database/sql"
	"errors"
	"log"

	"github.com/zgsolucoes/zg-data-guard/internal/database/storage"
	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	"github.com/zgsolucoes/zg-data-guard/internal/usecase/common"
)

type UpdateEcosystemUseCase struct {
	EcosystemStorage storage.EcosystemStorage
}

func NewUpdateEcosystemUseCase(ecosystemStorage storage.EcosystemStorage) *UpdateEcosystemUseCase {
	return &UpdateEcosystemUseCase{
		EcosystemStorage: ecosystemStorage,
	}
}

func (uc *UpdateEcosystemUseCase) Execute(input dto.EcosystemInputDTO, ecosystemID, operationUserID string) (*dto.EcosystemOutputDTO, error) {
	ecosystem, err := uc.EcosystemStorage.FindByID(ecosystemID)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		log.Printf("Ecosystem with id %s not found in database!", ecosystemID)
		return nil, common.ErrEcosystemNotFound
	}
	if err != nil {
		log.Printf("Error fetching ecosystem with id %s. Cause: %v", ecosystemID, err.Error())
		return nil, err
	}

	if ecosystem.Code != input.Code {
		codeExists, err := uc.EcosystemStorage.CheckCodeExists(input.Code)
		if err != nil {
			log.Printf("Error checking existance of ecosystem with code %s. Cause: %v", input.Code, err.Error())
			return nil, err
		}
		if codeExists {
			log.Printf("Error updating ecosystem. Cause: ecosystem with code %s already exists", input.Code)
			return nil, ErrCodeAlreadyExists
		}
	}

	ecosystem.Update(input.Code, input.DisplayName)
	err = uc.EcosystemStorage.Update(ecosystem)
	if err != nil {
		log.Printf("error updating ecosystem: %v", err.Error())
		return nil, err
	}

	log.Printf("Ecosystem %v updated successfully by user %s!", ecosystem.ID, operationUserID)
	return &dto.EcosystemOutputDTO{
		ID:          ecosystem.ID.String(),
		Code:        ecosystem.Code,
		DisplayName: ecosystem.DisplayName,
		CreatedAt:   ecosystem.CreatedAt,
		UpdatedAt:   &ecosystem.UpdatedAt,
	}, nil
}

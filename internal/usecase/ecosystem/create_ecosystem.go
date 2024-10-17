package ecosystem

import (
	"errors"
	"log"

	"github.com/zgsolucoes/zg-data-guard/internal/database/storage"
	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	"github.com/zgsolucoes/zg-data-guard/internal/entity"
)

var (
	ErrCodeAlreadyExists = errors.New("code already exists")
)

type CreateEcosystemUseCase struct {
	EcosystemStorage storage.EcosystemStorage
}

func NewCreateEcosystemUseCase(ecosystemStorage storage.EcosystemStorage) *CreateEcosystemUseCase {
	return &CreateEcosystemUseCase{
		EcosystemStorage: ecosystemStorage,
	}
}

func (c *CreateEcosystemUseCase) Execute(input dto.EcosystemInputDTO, createdByUserID string) (*dto.EcosystemOutputDTO, error) {
	ecosystem, err := entity.NewEcosystem(input.Code, input.DisplayName, createdByUserID)
	if err != nil {
		log.Printf("Error creating ecosystem. Cause: %v", err.Error())
		return nil, err
	}

	codeExists, err := c.EcosystemStorage.CheckCodeExists(input.Code)
	if err != nil {
		log.Printf("Error checking existance of ecosystem with code %s. Cause: %v", input.Code, err.Error())
		return nil, err
	}
	if codeExists {
		log.Printf("Error creating ecosystem. Cause: ecosystem with code %s already exists", input.Code)
		return nil, ErrCodeAlreadyExists
	}

	err = c.EcosystemStorage.Save(ecosystem)
	if err != nil {
		log.Printf("error creating ecosystem: %v", err.Error())
		return nil, err
	}

	log.Printf("Ecosystem %v created successfully by user %s!", ecosystem.ID, createdByUserID)
	return &dto.EcosystemOutputDTO{
		ID:              ecosystem.ID.String(),
		Code:            ecosystem.Code,
		DisplayName:     ecosystem.DisplayName,
		CreatedByUserID: ecosystem.CreatedByUserID,
		CreatedAt:       ecosystem.CreatedAt,
	}, nil
}

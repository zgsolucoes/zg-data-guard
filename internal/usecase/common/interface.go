package common

import (
	"github.com/zgsolucoes/zg-data-guard/internal/dto"
)

type RevokeAccessPermissionUseCaseInterface interface {
	Execute(input dto.RevokeAccessInputDTO, operationUserID string) (*dto.RevokeAccessOutputDTO, error)
}

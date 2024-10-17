package role

import (
	"log"

	"github.com/zgsolucoes/zg-data-guard/internal/database/storage"
	"github.com/zgsolucoes/zg-data-guard/internal/dto"
)

type ListDatabaseRolesUseCase struct {
	DatabaseRoleStorage storage.DatabaseRoleStorage
}

func NewListDatabaseRolesUseCase(databaseRoleStorage storage.DatabaseRoleStorage) *ListDatabaseRolesUseCase {
	return &ListDatabaseRolesUseCase{
		DatabaseRoleStorage: databaseRoleStorage,
	}
}

func (uc *ListDatabaseRolesUseCase) Execute() ([]*dto.DatabaseRoleOutputDTO, error) {
	roles, err := uc.DatabaseRoleStorage.FindAll()
	if err != nil {
		log.Printf("Error fetching database roles! Cause: %v", err.Error())
		return nil, err
	}
	rolesDTO := make([]*dto.DatabaseRoleOutputDTO, 0, len(roles))
	for _, role := range roles {
		rolesDTO = append(rolesDTO, &dto.DatabaseRoleOutputDTO{
			ID:              role.ID.String(),
			Name:            string(role.Name),
			DisplayName:     role.DisplayName,
			Description:     role.Description,
			ReadOnly:        role.ReadOnly,
			CreatedAt:       role.CreatedAt,
			CreatedByUserID: role.CreatedByUserID,
		})
	}
	log.Printf("List of database roles loaded successfully!")
	return rolesDTO, nil
}

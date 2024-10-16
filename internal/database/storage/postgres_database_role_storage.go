package storage

import (
	"database/sql"

	"github.com/zgsolucoes/zg-data-guard/internal/entity"
)

type PostgresDatabaseRoleStorage struct {
	db *sql.DB
}

func NewPostgresDatabaseRoleStorage(db *sql.DB) *PostgresDatabaseRoleStorage {
	return &PostgresDatabaseRoleStorage{db: db}
}

func (r *PostgresDatabaseRoleStorage) FindAll() ([]*entity.DatabaseRole, error) {
	query := `SELECT id, name, display_name, description, read_only, created_at, created_by_user_id FROM database_roles`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []*entity.DatabaseRole
	for rows.Next() {
		var role entity.DatabaseRole
		err := rows.Scan(&role.ID, &role.Name, &role.DisplayName, &role.Description, &role.ReadOnly, &role.CreatedAt, &role.CreatedByUserID)
		if err != nil {
			return nil, err
		}
		roles = append(roles, &role)
	}

	return roles, nil
}

func (r *PostgresDatabaseRoleStorage) FindByID(id string) (*entity.DatabaseRole, error) {
	query := `SELECT id, name, display_name, description, read_only, created_at, created_by_user_id FROM database_roles WHERE id = $1`

	var role entity.DatabaseRole
	err := r.db.QueryRow(query, id).Scan(&role.ID, &role.Name, &role.DisplayName, &role.Description, &role.ReadOnly, &role.CreatedAt, &role.CreatedByUserID)
	if err != nil {
		return nil, err
	}

	return &role, nil
}

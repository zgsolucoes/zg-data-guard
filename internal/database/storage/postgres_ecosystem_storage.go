package storage

import (
	"database/sql"

	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	"github.com/zgsolucoes/zg-data-guard/internal/entity"
)

type PostgresEcosystemStorage struct {
	DB *sql.DB
}

func NewPostgresEcosystemStorage(db *sql.DB) *PostgresEcosystemStorage {
	return &PostgresEcosystemStorage{DB: db}
}

func (er *PostgresEcosystemStorage) CheckCodeExists(code string) (bool, error) {
	var codeExists bool
	err := er.DB.QueryRow("SELECT EXISTS( SELECT 1 FROM ecosystems WHERE code LIKE $1) AS exists", code).Scan(&codeExists)
	if err != nil {
		return false, err
	}
	return codeExists, nil
}

func (er *PostgresEcosystemStorage) Save(e *entity.Ecosystem) error {
	stmt, err := er.DB.Prepare("INSERT INTO ecosystems (id, code, display_name, created_at, updated_at, created_by_user_id) VALUES ($1, $2, $3, $4, $5, $6)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(e.ID, e.Code, e.DisplayName, e.CreatedAt, e.UpdatedAt, e.CreatedByUserID)
	if err != nil {
		return err
	}
	return nil
}

func (er *PostgresEcosystemStorage) Update(e *entity.Ecosystem) error {
	stmt, err := er.DB.Prepare("UPDATE ecosystems SET code = $2, display_name = $3, updated_at = $4 WHERE id = $1")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(e.ID, e.Code, e.DisplayName, e.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (er *PostgresEcosystemStorage) FindByID(id string) (*entity.Ecosystem, error) {
	var e entity.Ecosystem
	err := er.DB.QueryRow("SELECT id, code, display_name, created_at, updated_at, created_by_user_id FROM ecosystems WHERE id = $1", id).
		Scan(&e.ID, &e.Code, &e.DisplayName, &e.CreatedAt, &e.UpdatedAt, &e.CreatedByUserID)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (er *PostgresEcosystemStorage) Delete(id string) error {
	_, err := er.FindByID(id)
	if err != nil {
		return err
	}
	stmt, err := er.DB.Prepare("DELETE FROM ecosystems WHERE id = $1")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return nil
}

func (er *PostgresEcosystemStorage) FindAll(page, limit int) ([]*dto.EcosystemOutputDTO, error) {
	var ecosystems []*dto.EcosystemOutputDTO
	rows, err := er.DB.Query(`
SELECT e.id,
       code,
       display_name,
       e.created_at,
       e.updated_at,
       u.name,
       created_by_user_id 
FROM ecosystems e 
	JOIN application_users u ON e.created_by_user_id = u.id 
ORDER BY e.display_name OFFSET $1 LIMIT $2`, (page-1)*limit, limit)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var e dto.EcosystemOutputDTO
		err := rows.Scan(&e.ID, &e.Code, &e.DisplayName, &e.CreatedAt, &e.UpdatedAt, &e.CreatedByUser, &e.CreatedByUserID)
		if err != nil {
			return nil, err
		}
		ecosystems = append(ecosystems, &e)
	}
	return ecosystems, nil
}

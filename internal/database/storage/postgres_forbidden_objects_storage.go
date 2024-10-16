package storage

import (
	"database/sql"

	"github.com/zgsolucoes/zg-data-guard/internal/entity"
)

type PostgresForbiddenObjectsStorage struct {
	db *sql.DB
}

func NewPostgresForbiddenObjectsStorage(db *sql.DB) *PostgresForbiddenObjectsStorage {
	return &PostgresForbiddenObjectsStorage{db: db}
}

func (r *PostgresForbiddenObjectsStorage) FindAllDatabases() ([]*entity.ForbiddenDatabase, error) {
	query := `SELECT id, database_name, description, created_at, created_by_user_id, updated_at FROM forbidden_databases`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var forbiddenDbs []*entity.ForbiddenDatabase

	for rows.Next() {
		var forbiddenDB entity.ForbiddenDatabase
		err := rows.Scan(&forbiddenDB.ID, &forbiddenDB.Name, &forbiddenDB.Description, &forbiddenDB.CreatedAt, &forbiddenDB.CreatedByUserID, &forbiddenDB.UpdatedAt)
		if err != nil {
			return nil, err
		}
		forbiddenDbs = append(forbiddenDbs, &forbiddenDB)
	}
	return forbiddenDbs, nil
}

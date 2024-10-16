package storage

import (
	"database/sql"

	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	"github.com/zgsolucoes/zg-data-guard/internal/entity"
)

type PostgresDatabaseTechnologyStorage struct {
	DB *sql.DB
}

func NewPostgresDatabaseTechnologyStorage(db *sql.DB) *PostgresDatabaseTechnologyStorage {
	return &PostgresDatabaseTechnologyStorage{DB: db}
}

func (dtr *PostgresDatabaseTechnologyStorage) Exists(name, version string) (bool, error) {
	var exists bool
	err := dtr.DB.QueryRow("SELECT EXISTS( SELECT 1 FROM database_technologies WHERE name LIKE $1 AND version LIKE $2) AS exists", name, version).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (dtr *PostgresDatabaseTechnologyStorage) Save(databaseTechnology *entity.DatabaseTechnology) error {
	stmt, err := dtr.DB.Prepare("INSERT INTO database_technologies (id, name, version, created_at, updated_at, created_by_user_id) VALUES ($1, $2, $3, $4, $5, $6)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(databaseTechnology.ID, databaseTechnology.Name, databaseTechnology.Version, databaseTechnology.CreatedAt, databaseTechnology.UpdatedAt, databaseTechnology.CreatedByUserID)
	if err != nil {
		return err
	}
	return nil
}

func (dtr *PostgresDatabaseTechnologyStorage) Update(databaseTechnology *entity.DatabaseTechnology) error {
	stmt, err := dtr.DB.Prepare("UPDATE database_technologies SET name = $2, version = $3, updated_at = $4 WHERE id = $1")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(databaseTechnology.ID, databaseTechnology.Name, databaseTechnology.Version, databaseTechnology.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (dtr *PostgresDatabaseTechnologyStorage) FindByID(id string) (*entity.DatabaseTechnology, error) {
	var databaseTechnology entity.DatabaseTechnology
	err := dtr.DB.QueryRow("SELECT id, name, version, created_at, updated_at, created_by_user_id FROM database_technologies WHERE id = $1", id).
		Scan(&databaseTechnology.ID, &databaseTechnology.Name, &databaseTechnology.Version, &databaseTechnology.CreatedAt, &databaseTechnology.UpdatedAt, &databaseTechnology.CreatedByUserID)
	if err != nil {
		return nil, err
	}
	return &databaseTechnology, nil
}

func (dtr *PostgresDatabaseTechnologyStorage) Delete(id string) error {
	_, err := dtr.FindByID(id)
	if err != nil {
		return err
	}
	stmt, err := dtr.DB.Prepare("DELETE FROM database_technologies WHERE id = $1")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return nil
}

func (dtr *PostgresDatabaseTechnologyStorage) FindAll(page, limit int) ([]*dto.TechnologyOutputDTO, error) {
	var databaseTechnologies []*dto.TechnologyOutputDTO
	rows, err := dtr.DB.Query(`
SELECT dbtech.id,
       dbtech.name,
       dbtech.version,
       dbtech.created_at,
       dbtech.updated_at,
       u.name,
       dbtech.created_by_user_id
FROM database_technologies dbtech
	JOIN application_users u ON dbtech.created_by_user_id = u.id 
ORDER BY dbtech.name, dbtech.version LIMIT $1 OFFSET $2`, limit, (page-1)*limit)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var tec dto.TechnologyOutputDTO
		err = rows.Scan(&tec.ID, &tec.Name, &tec.Version, &tec.CreatedAt, &tec.UpdatedAt, &tec.CreatedByUser, &tec.CreatedByUserID)
		if err != nil {
			return nil, err
		}
		databaseTechnologies = append(databaseTechnologies, &tec)
	}
	return databaseTechnologies, nil
}

package storage

import (
	"database/sql"
	"time"

	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	"github.com/zgsolucoes/zg-data-guard/internal/entity"
)

type PostgresDatabaseStorage struct {
	db *sql.DB
}

func NewPostgresDatabaseStorage(db *sql.DB) *PostgresDatabaseStorage {
	return &PostgresDatabaseStorage{db: db}
}

func (r *PostgresDatabaseStorage) Save(database *entity.Database) error {
	query := `INSERT INTO databases (id, name, description, current_size, enabled, roles_configured, database_instance_id, created_at, created_by_user_id, updated_at)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := r.db.Exec(
		query,
		database.ID,
		database.Name,
		database.Description,
		database.CurrentSize,
		database.Enabled,
		database.RolesConfigured,
		database.DatabaseInstanceID,
		database.CreatedAt,
		database.CreatedByUserID,
		database.UpdatedAt)
	return err
}

func (r *PostgresDatabaseStorage) Update(database *entity.Database) error {
	query := `
UPDATE 
	databases
SET name = $1,
    description = $2,
    current_size = $3,
    enabled = $4,
    database_instance_id = $5,
    updated_at = $6,
    disabled_at = $7,
    roles_configured = $8 
WHERE id = $9`
	_, err := r.db.Exec(
		query,
		database.Name,
		database.Description,
		database.CurrentSize,
		database.Enabled,
		database.DatabaseInstanceID,
		database.UpdatedAt,
		database.DisabledAt,
		database.RolesConfigured,
		database.ID)
	return err
}

/** func (r *DatabaseStorage) FindByID(id string) (*entity.Database, error) {
	query := `SELECT id, name, description, current_size, enabled, roles_configured, database_instance_id, created_at, created_by_user_id, updated_at, disabled_at FROM databases WHERE id = $1`
	row := r.db.QueryRow(query, id)

	var database entity.Database
	err := row.Scan(
		&database.ID,
		&database.Name,
		&database.Description,
		&database.CurrentSize,
		&database.Enabled,
		&database.RolesConfigured,
		&database.DatabaseInstanceID,
		&database.CreatedAt,
		&database.CreatedByUserID,
		&database.UpdatedAt,
		&database.DisabledAt)
	if err != nil {
		return nil, err
	}

	return &database, nil
} */

func (r *PostgresDatabaseStorage) FindDTOByID(id string) (*dto.DatabaseOutputDTO, error) {
	var output dto.DatabaseOutputDTO
	baseQuery, err := ReadSQLFile("internal/database/sqls/select_database_by_id.sql")
	if err != nil {
		return nil, err
	}
	row := r.db.QueryRow(baseQuery, id)
	err = row.Scan(
		&output.ID,
		&output.Name,
		&output.CurrentSize,
		&output.DatabaseInstanceID,
		&output.DatabaseInstanceName,
		&output.EcosystemID,
		&output.EcosystemName,
		&output.DatabaseTechnologyID,
		&output.DatabaseTechnologyName,
		&output.DatabaseTechnologyVersion,
		&output.Enabled,
		&output.RolesConfigured,
		&output.Description,
		&output.CreatedByUserID,
		&output.CreatedByUser,
		&output.CreatedAt,
		&output.UpdatedAt,
		&output.LastDatabaseSync,
		&output.DisabledAt)
	if err != nil {
		return nil, err
	}
	return &output, nil
}

func (r *PostgresDatabaseStorage) FindAll(databaseInstanceID string, ids []string) ([]*entity.Database, error) {
	baseQuery := `
SELECT id,
       name,
       description,
       current_size,
       enabled,
       roles_configured,
       database_instance_id,
       created_at,
       created_by_user_id,
       updated_at,
       disabled_at
FROM databases db
WHERE 1 = 1`

	var args []any
	baseQuery, args = addFilterCondition(baseQuery, args, "db.database_instance_id", databaseInstanceID)
	baseQuery, args = appendFilterIdsInQuery(baseQuery, "db", ids, args)
	baseQuery += " ORDER BY db.name, db.database_instance_id"
	rows, err := executeSQLQuery(r.db, baseQuery, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var databases []*entity.Database
	for rows.Next() {
		database := &entity.Database{}
		err := rows.Scan(
			&database.ID,
			&database.Name,
			&database.Description,
			&database.CurrentSize,
			&database.Enabled,
			&database.RolesConfigured,
			&database.DatabaseInstanceID,
			&database.CreatedAt,
			&database.CreatedByUserID,
			&database.UpdatedAt,
			&database.DisabledAt)
		if err != nil {
			return nil, err
		}
		databases = append(databases, database)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return databases, nil
}

func (r *PostgresDatabaseStorage) FindAllEnabled(databaseInstanceID string) ([]*entity.Database, error) {
	databases, err := r.FindAll(databaseInstanceID, nil)
	if err != nil {
		return nil, err
	}

	var enabledDbs []*entity.Database
	for _, database := range databases {
		if database.Enabled {
			enabledDbs = append(enabledDbs, database)
		}
	}
	return enabledDbs, nil
}

func (r *PostgresDatabaseStorage) FindAllDTOs(ecosystemID, databaseInstanceID string) ([]*dto.DatabaseOutputDTO, error) {
	baseQuery, err := ReadSQLFile("internal/database/sqls/select_databases.sql")
	if err != nil {
		return nil, err
	}

	var args []any
	baseQuery, args = addFilterCondition(baseQuery, args, "di.ecosystem_id", ecosystemID)
	baseQuery, args = addFilterCondition(baseQuery, args, "db.database_instance_id", databaseInstanceID)
	baseQuery += " ORDER BY db.name, di.name"
	rows, err := r.db.Query(baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var databases []*dto.DatabaseOutputDTO
	for rows.Next() {
		output, err := scanDatabase(rows)
		if err != nil {
			return nil, err
		}
		databases = append(databases, &output)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return databases, nil
}

func (r *PostgresDatabaseStorage) DeactivateAllByInstance(databaseInstanceID string) error {
	currentTime := time.Now()
	query := `UPDATE databases SET enabled = FALSE, disabled_at = $1, updated_at = $2 WHERE database_instance_id = $3`
	_, err := r.db.Exec(query, currentTime, currentTime, databaseInstanceID)
	return err
}

func scanDatabase(rows *sql.Rows) (dto.DatabaseOutputDTO, error) {
	var output dto.DatabaseOutputDTO
	err := rows.Scan(&output.ID,
		&output.Name,
		&output.CurrentSize,
		&output.DatabaseInstanceID,
		&output.DatabaseInstanceName,
		&output.EcosystemID,
		&output.EcosystemName,
		&output.DatabaseTechnologyID,
		&output.DatabaseTechnologyName,
		&output.DatabaseTechnologyVersion,
		&output.Enabled,
		&output.RolesConfigured,
		&output.Description,
		&output.CreatedByUserID,
		&output.CreatedByUser,
		&output.CreatedAt,
		&output.UpdatedAt,
		&output.LastDatabaseSync,
		&output.DisabledAt)
	return output, err
}

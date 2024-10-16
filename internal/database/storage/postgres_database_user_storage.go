package storage

import (
	"database/sql"

	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	"github.com/zgsolucoes/zg-data-guard/internal/entity"
)

type PostgresDatabaseUserStorage struct {
	db *sql.DB
}

func NewPostgresDatabaseUserStorage(db *sql.DB) *PostgresDatabaseUserStorage {
	return &PostgresDatabaseUserStorage{db: db}
}

func (dur *PostgresDatabaseUserStorage) Save(d *entity.DatabaseUser) error {
	query := `INSERT INTO database_users (id, name, email, username, password, database_role_id, enabled, team, position, created_at, created_by_user_id, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	_, err := dur.db.Exec(
		query,
		d.ID,
		d.Name,
		d.Email,
		d.Username,
		d.CipherPassword,
		d.DatabaseRoleID,
		d.Enabled,
		d.Team,
		d.Position,
		d.CreatedAt,
		d.CreatedByUserID,
		d.UpdatedAt)
	return err
}

func (dur *PostgresDatabaseUserStorage) FindByID(id string) (*entity.DatabaseUser, error) {
	query := `SELECT id, name, email, username, password, database_role_id, enabled, team, position, created_at, created_by_user_id, updated_at, disabled_at, expired, expires_at
			FROM database_users WHERE id = $1`
	row := dur.db.QueryRow(query, id)

	var d entity.DatabaseUser
	err := row.Scan(
		&d.ID,
		&d.Name,
		&d.Email,
		&d.Username,
		&d.CipherPassword,
		&d.DatabaseRoleID,
		&d.Enabled,
		&d.Team,
		&d.Position,
		&d.CreatedAt,
		&d.CreatedByUserID,
		&d.UpdatedAt,
		&d.DisabledAt,
		&d.Expired,
		&d.ExpiresAt)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (dur *PostgresDatabaseUserStorage) Update(d *entity.DatabaseUser) error {
	query := `UPDATE database_users SET name = $1, team = $2, position = $3, database_role_id = $4, enabled = $5, updated_at = $6, disabled_at = $7 WHERE id = $8`
	_, err := dur.db.Exec(
		query,
		d.Name,
		d.Team,
		d.Position,
		d.DatabaseRoleID,
		d.Enabled,
		d.UpdatedAt,
		d.DisabledAt,
		d.ID)
	return err
}

func (dur *PostgresDatabaseUserStorage) Exists(email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM database_users WHERE email ILIKE $1)`

	var exists bool
	err := dur.db.QueryRow(query, email).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (dur *PostgresDatabaseUserStorage) FindDTOByID(id string) (*dto.DatabaseUserOutputDTO, error) {
	var query = dur.baseQueryDTO() + ` WHERE du.id = $1`
	row := dur.db.QueryRow(query, id)

	var d dto.DatabaseUserOutputDTO
	err := row.Scan(
		&d.ID,
		&d.Name,
		&d.Email,
		&d.Username,
		&d.Password,
		&d.DatabaseRoleID,
		&d.DatabaseRoleName,
		&d.DatabaseRoleDisplayName,
		&d.Enabled,
		&d.Team,
		&d.Position,
		&d.CreatedByUserID,
		&d.CreatedByUser,
		&d.CreatedAt,
		&d.UpdatedAt,
		&d.DisabledAt)

	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (dur *PostgresDatabaseUserStorage) FindAll(ids []string) ([]*entity.DatabaseUser, error) {
	baseQuery := dur.baseQuery()
	var args []any
	baseQuery, args = appendFilterIdsInQuery(baseQuery, "du", ids, args)
	baseQuery += " ORDER BY du.name, du.email"
	rows, err := dur.db.Query(baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var dbUsers []*entity.DatabaseUser
	for rows.Next() {
		var d entity.DatabaseUser
		err := rows.Scan(
			&d.ID,
			&d.Name,
			&d.Email,
			&d.Username,
			&d.CipherPassword,
			&d.DatabaseRoleID,
			&d.Enabled,
			&d.Team,
			&d.Position,
			&d.CreatedAt,
			&d.CreatedByUserID,
			&d.UpdatedAt,
			&d.DisabledAt,
			&d.Expired,
			&d.ExpiresAt)
		if err != nil {
			return nil, err
		}
		dbUsers = append(dbUsers, &d)
	}
	return dbUsers, nil
}

func (dur *PostgresDatabaseUserStorage) FindAllDTOs(ids []string) ([]*dto.DatabaseUserOutputDTO, error) {
	baseQuery := dur.baseQueryDTO()
	var args []any
	baseQuery, args = appendFilterIdsInQuery(baseQuery, "du", ids, args)
	baseQuery += " ORDER BY du.name, du.email"
	rows, err := dur.db.Query(baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dbUsersDTO []*dto.DatabaseUserOutputDTO
	for rows.Next() {
		var d dto.DatabaseUserOutputDTO
		err := rows.Scan(
			&d.ID,
			&d.Name,
			&d.Email,
			&d.Username,
			&d.Password,
			&d.DatabaseRoleID,
			&d.DatabaseRoleName,
			&d.DatabaseRoleDisplayName,
			&d.Enabled,
			&d.Team,
			&d.Position,
			&d.CreatedByUserID,
			&d.CreatedByUser,
			&d.CreatedAt,
			&d.UpdatedAt,
			&d.DisabledAt)
		if err != nil {
			return nil, err
		}
		dbUsersDTO = append(dbUsersDTO, &d)
	}

	return dbUsersDTO, nil
}

func (dur *PostgresDatabaseUserStorage) FindAllDTOsEnabled() ([]*dto.DatabaseUserOutputDTO, error) {
	dbUsers, err := dur.FindAllDTOs(nil)
	if err != nil {
		return nil, err
	}

	var enabledDBUsers []*dto.DatabaseUserOutputDTO
	for _, dbUser := range dbUsers {
		if dbUser.Enabled {
			enabledDBUsers = append(enabledDBUsers, dbUser)
		}
	}
	return enabledDBUsers, nil
}

func (dur *PostgresDatabaseUserStorage) baseQuery() string {
	return `
SELECT id, 
       name, 
       email,
       username, 
       password, 
       database_role_id, 
       enabled, 
       team, 
       position, 
       created_at, 
       created_by_user_id, 
       updated_at,
       disabled_at,
       expired, 
       expires_at
FROM database_users du
WHERE 1 = 1`
}

func (dur *PostgresDatabaseUserStorage) baseQueryDTO() string {
	return `
SELECT 
       du.id, 
       du.name, 
       du.email, 
       du.username,
       du.password,
       du.database_role_id, 
       dr.name,
       dr.display_name,
       du.enabled, 
       du.team, 
       du.position,
       du.created_by_user_id,
       au.name,
       du.created_at,
       du.updated_at,
       du.disabled_at
FROM database_users du
	JOIN database_roles dr 
		ON du.database_role_id = dr.id
	JOIN application_users au 
		ON du.created_by_user_id = au.id
WHERE 1 = 1`
}

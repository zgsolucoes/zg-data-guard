package storage

import (
	"database/sql"

	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	"github.com/zgsolucoes/zg-data-guard/internal/entity"
)

type PostgresAccessPermissionStorage struct {
	db *sql.DB
}

func NewPostgresAccessPermissionStorage(db *sql.DB) *PostgresAccessPermissionStorage {
	return &PostgresAccessPermissionStorage{db: db}
}

func (ar *PostgresAccessPermissionStorage) Save(d *entity.AccessPermission) error {
	query := `INSERT INTO access_permissions (id, database_id, database_user_id, granted_by_user_id, granted_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := ar.db.Exec(
		query,
		d.ID,
		d.DatabaseID,
		d.DatabaseUserID,
		d.GrantedByUserID,
		d.GrantedAt)
	return err
}

func (ar *PostgresAccessPermissionStorage) Exists(databaseID, databaseUserID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM access_permissions WHERE database_id = $1 AND database_user_id = $2)`

	var exists bool
	err := ar.db.QueryRow(query, databaseID, databaseUserID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (ar *PostgresAccessPermissionStorage) FindAllDTOs(databaseID, databaseUserID, databaseInstanceID string) ([]*dto.AccessPermissionOutputDTO, error) {
	baseQuery := ar.baseQueryDTO()
	var args []any
	baseQuery, args = addFilterCondition(baseQuery, args, "ap.database_id", databaseID)
	baseQuery, args = addFilterCondition(baseQuery, args, "ap.database_user_id", databaseUserID)
	baseQuery, args = addFilterCondition(baseQuery, args, "di.id", databaseInstanceID)
	baseQuery += " ORDER BY db_user.name, ap.granted_at"
	rows, err := ar.db.Query(baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accessDTOs []*dto.AccessPermissionOutputDTO
	for rows.Next() {
		var d dto.AccessPermissionOutputDTO
		err := rows.Scan(
			&d.ID,
			&d.DatabaseUserID,
			&d.DatabaseUserName,
			&d.DatabaseUserEmail,
			&d.DatabaseRoleID,
			&d.DatabaseRoleName,
			&d.EcosystemID,
			&d.EcosystemName,
			&d.DatabaseInstanceID,
			&d.DatabaseInstanceName,
			&d.DatabaseID,
			&d.DatabaseName,
			&d.GrantedByUserID,
			&d.GrantedByUserName,
			&d.GrantedAt)
		if err != nil {
			return nil, err
		}
		accessDTOs = append(accessDTOs, &d)
	}

	return accessDTOs, nil
}

func (ar *PostgresAccessPermissionStorage) DeleteAllByUserAndInstance(databaseUserID, instanceID string) error {
	query := `DELETE FROM access_permissions WHERE database_user_id = $1 AND database_id IN (SELECT id FROM databases WHERE database_instance_id = $2)`
	_, err := ar.db.Exec(query, databaseUserID, instanceID)
	return err
}

func (ar *PostgresAccessPermissionStorage) DeleteAllByInstance(instanceID string) error {
	query := `DELETE FROM access_permissions WHERE database_id IN (SELECT id FROM databases WHERE database_instance_id = $1)`
	_, err := ar.db.Exec(query, instanceID)
	return err
}

func (ar *PostgresAccessPermissionStorage) SaveLog(log *entity.AccessPermissionLog) error {
	query := `INSERT INTO access_permission_log (id, database_instance_id, database_user_id, database_id, message, success, date, user_id) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := ar.db.Exec(
		query,
		log.ID,
		log.DatabaseInstanceID,
		log.DatabaseUserID,
		log.DatabaseID,
		log.Message,
		log.Success,
		log.Date,
		log.UserID)
	return err
}

func (ar *PostgresAccessPermissionStorage) FindAllAccessibleInstancesIDsByUser(userID string) ([]string, error) {
	query := `
SELECT DISTINCT di.id
FROM access_permissions ap
	JOIN databases db
		ON ap.database_id = db.id
	JOIN database_instances di
		ON db.database_instance_id = di.id
WHERE ap.database_user_id = $1`

	rows, err := ar.db.Query(query, userID)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var instanceIDs []string
	for rows.Next() {
		var instanceID string
		err := rows.Scan(&instanceID)
		if err != nil {
			return nil, err
		}
		instanceIDs = append(instanceIDs, instanceID)
	}

	return instanceIDs, nil
}

func (ar *PostgresAccessPermissionStorage) FindAllLogsDTOs(page, limit int) ([]*dto.AccessPermissionLogOutputDTO, error) {
	var logs []*dto.AccessPermissionLogOutputDTO
	query := ar.baseQueryLogDTO()
	rows, err := ar.db.Query(query, (page-1)*limit, limit)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	for rows.Next() {
		var log dto.AccessPermissionLogOutputDTO
		err := rows.Scan(
			&log.ID,
			&log.DatabaseUserID,
			&log.DatabaseUserName,
			&log.DatabaseUserEmail,
			&log.DatabaseInstanceID,
			&log.DatabaseInstanceName,
			&log.DatabaseID,
			&log.DatabaseName,
			&log.Message,
			&log.Success,
			&log.Date,
			&log.OperationUserID,
			&log.OperationUserName,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, &log)
	}
	return logs, nil
}

func (ar *PostgresAccessPermissionStorage) CheckIfUserHasAccessPermission(databaseUserID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM access_permissions WHERE database_user_id = $1)`
	var exists bool
	err := ar.db.QueryRow(query, databaseUserID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (ar *PostgresAccessPermissionStorage) LogCount() (int, error) {
	query := `SELECT COUNT(*) FROM access_permission_log`
	var count int
	err := ar.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (ar *PostgresAccessPermissionStorage) baseQueryDTO() string {
	return `
SELECT
       ap.id,
       ap.database_user_id,
       db_user.name,
       db_user.email,
       COALESCE(ap.database_role_id, db_user.database_role_id),
       COALESCE(role_from_access.display_name, role.display_name),
       e.id,
       e.display_name,
       di.id,
       di.name,
       ap.database_id,
       db.name,
       ap.granted_by_user_id,
       op_user.name,
       ap.granted_at
FROM access_permissions ap
	JOIN databases db
		ON ap.database_id = db.id
	JOIN database_instances di
		ON db.database_instance_id = di.id
	JOIN ecosystems e
		ON di.ecosystem_id = e.id
	JOIN database_users db_user
		ON ap.database_user_id = db_user.id
	JOIN database_roles role
		ON db_user.database_role_id = role.id
	LEFT JOIN database_roles role_from_access
		ON ap.database_role_id = role_from_access.id
	JOIN application_users op_user
		ON ap.granted_by_user_id = op_user.id
WHERE 1 = 1`
}

func (ar *PostgresAccessPermissionStorage) baseQueryLogDTO() string {
	return `
SELECT
	log.id,
	log.database_user_id,
	db_user.name,
	db_user.email,
	log.database_instance_id,
	di.name,
	log.database_id,
	db.name,
	log.message,
	log.success,
	log.date,
	log.user_id,
	op_user.name
FROM access_permission_log log
	JOIN database_instances di
		ON log.database_instance_id = di.id
	JOIN application_users op_user
		ON log.user_id = op_user.id
	LEFT JOIN database_users db_user
		ON log.database_user_id = db_user.id
	LEFT JOIN databases db
		ON log.database_id = db.id
ORDER BY log.date DESC OFFSET $1 LIMIT $2`
}

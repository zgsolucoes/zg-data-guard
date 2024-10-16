package storage

import (
	"database/sql"
	"log"

	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	"github.com/zgsolucoes/zg-data-guard/internal/entity"
)

type PostgresDatabaseInstanceStorage struct {
	Uow *UnitOfWork
}

func NewPostgresInstanceStorage(db *sql.DB) *PostgresDatabaseInstanceStorage {
	uow := NewUnitOfWork(db)
	return &PostgresDatabaseInstanceStorage{
		Uow: uow,
	}
}

func (dir *PostgresDatabaseInstanceStorage) Save(databaseInstance *entity.DatabaseInstance) error {
	saveOperation := func() error {
		err := dir.insertHostConnectionInfo(databaseInstance.HostConnection)
		if err != nil {
			return err
		}
		log.Printf("Host connection inserted with id: %s", databaseInstance.HostConnection.ID)
		err = dir.insertDatabaseInstance(databaseInstance)
		if err != nil {
			return err
		}
		return nil
	}

	return dir.Uow.ExecuteInTransaction(saveOperation)
}

func (dir *PostgresDatabaseInstanceStorage) UpdateWithHostInfo(databaseInstance *entity.DatabaseInstance) error {
	updateOperation := func() error {
		err := dir.updateHostConnectionInfo(databaseInstance.HostConnection)
		if err != nil {
			return err
		}
		log.Printf("Host connection updated with id: %s", databaseInstance.HostConnection.ID)
		err = dir.updateDatabaseInstance(databaseInstance)
		if err != nil {
			return err
		}
		return nil
	}

	return dir.Uow.ExecuteInTransaction(updateOperation)
}

func (dir *PostgresDatabaseInstanceStorage) Exists(host string, port string) (bool, error) {
	var count int
	err := dir.Uow.db.QueryRow("SELECT COUNT(*) FROM host_connection_info WHERE host = $1 AND port = $2", host, port).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (dir *PostgresDatabaseInstanceStorage) Update(databaseInstance *entity.DatabaseInstance) error {
	query := `
UPDATE database_instances 
SET name = $2, 
    ecosystem_id = $3, 
    database_technology_id = $4, 
    enabled = $5, 
    note = $6, 
    updated_at = $7, 
    disabled_at = $8, 
    last_database_sync = $9,
    connection_status = $10,
    last_connection_test = $11,
    last_connection_result = $12,
    roles_created = $13
WHERE id = $1`
	_, err := dir.Uow.db.Exec(
		query,
		databaseInstance.ID,
		databaseInstance.Name,
		databaseInstance.EcosystemID,
		databaseInstance.DatabaseTechnologyID,
		databaseInstance.Enabled,
		databaseInstance.Note,
		databaseInstance.UpdatedAt,
		databaseInstance.DisabledAt,
		databaseInstance.LastDatabaseSync,
		databaseInstance.ConnectionStatus,
		databaseInstance.LastConnectionTest,
		databaseInstance.LastConnectionResult,
		databaseInstance.RolesCreated,
	)
	return err
}

func (dir *PostgresDatabaseInstanceStorage) FindByID(id string) (*entity.DatabaseInstance, error) {
	var databaseInstance entity.DatabaseInstance
	var hostConnection entity.HostConnectionInfo
	err := dir.Uow.db.QueryRow(`
SELECT di.id, 
       name, 
       host_connection_info_id,
       hci.host,
       hci.port,
       hci.host_connection,
       hci.port_connection,
       hci.admin_username,
       hci.admin_password,
       ecosystem_id, 
       database_technology_id, 
       enabled,
       roles_created,
       note, 
       created_at, 
       created_by_user_id, 
       updated_at, 
       disabled_at,
       last_database_sync,
       connection_status,
       last_connection_test,
       last_connection_result
FROM database_instances di
	JOIN host_connection_info hci ON di.host_connection_info_id = hci.id
WHERE di.id = $1`, id).
		Scan(&databaseInstance.ID,
			&databaseInstance.Name,
			&hostConnection.ID,
			&hostConnection.Host,
			&hostConnection.Port,
			&hostConnection.HostConnection,
			&hostConnection.PortConnection,
			&hostConnection.AdminUser,
			&hostConnection.AdminPassword,
			&databaseInstance.EcosystemID,
			&databaseInstance.DatabaseTechnologyID,
			&databaseInstance.Enabled,
			&databaseInstance.RolesCreated,
			&databaseInstance.Note,
			&databaseInstance.CreatedAt,
			&databaseInstance.CreatedByUserID,
			&databaseInstance.UpdatedAt,
			&databaseInstance.DisabledAt,
			&databaseInstance.LastDatabaseSync,
			&databaseInstance.ConnectionStatus,
			&databaseInstance.LastConnectionTest,
			&databaseInstance.LastConnectionResult)
	if err != nil {
		return nil, err
	}
	databaseInstance.HostConnection = &hostConnection
	return &databaseInstance, nil
}

func (dir *PostgresDatabaseInstanceStorage) FindDTOByID(id string) (*dto.DatabaseInstanceOutputDTO, error) {
	var output dto.DatabaseInstanceOutputDTO
	baseQuery, err := ReadSQLFile("internal/database/sqls/select_database_instance_by_id.sql")
	if err != nil {
		return nil, err
	}
	err = dir.Uow.db.QueryRow(baseQuery, id).
		Scan(&output.ID,
			&output.Name,
			&output.Host,
			&output.Port,
			&output.HostConnection,
			&output.PortConnection,
			&output.AdminUser,
			&output.AdminPassword,
			&output.EcosystemID,
			&output.EcosystemName,
			&output.DatabaseTechnologyID,
			&output.DatabaseTechnologyName,
			&output.DatabaseTechnologyVersion,
			&output.Enabled,
			&output.RolesCreated,
			&output.Note,
			&output.CreatedAt,
			&output.CreatedByUserID,
			&output.CreatedByUser,
			&output.UpdatedAt,
			&output.DisabledAt,
			&output.LastDatabaseSync,
			&output.ConnectionStatus,
			&output.LastConnectionTest,
			&output.LastConnectionResult)
	if err != nil {
		return nil, err
	}
	return &output, nil
}

func (dir *PostgresDatabaseInstanceStorage) FindAllDTOs(ecosystemID, technologyID string, ids []string) ([]*dto.DatabaseInstanceOutputDTO, error) {
	baseQuery, err := ReadSQLFile("internal/database/sqls/select_database_instances.sql")
	if err != nil {
		return nil, err
	}

	var args []any
	baseQuery, args = addFilterCondition(baseQuery, args, "di.ecosystem_id", ecosystemID)
	baseQuery, args = addFilterCondition(baseQuery, args, "di.database_technology_id", technologyID)
	baseQuery, args = appendFilterIdsInQuery(baseQuery, "di", ids, args)
	baseQuery += " ORDER BY di.name, e.display_name, dt.name, dt.version"
	rows, err := executeSQLQuery(dir.Uow.db, baseQuery, args)
	if err != nil {
		return nil, err
	}

	databaseInstances, err := fetchDataFromRows(rows)
	if err != nil {
		return nil, err
	}

	return databaseInstances, nil
}

func (dir *PostgresDatabaseInstanceStorage) FindAllDTOsEnabled(ecosystemID, technologyID string) ([]*dto.DatabaseInstanceOutputDTO, error) {
	databaseInstances, err := dir.FindAllDTOs(ecosystemID, technologyID, nil)
	if err != nil {
		return nil, err
	}

	var enabledInstances []*dto.DatabaseInstanceOutputDTO
	for _, instance := range databaseInstances {
		if instance.Enabled {
			enabledInstances = append(enabledInstances, instance)
		}
	}
	return enabledInstances, nil
}

func (dir *PostgresDatabaseInstanceStorage) insertHostConnectionInfo(h *entity.HostConnectionInfo) error {
	stmt, err := dir.Uow.transaction.Prepare("INSERT INTO host_connection_info (id, host, port, host_connection, port_connection, admin_username, admin_password) VALUES ($1, $2, $3, $4, $5, $6, $7)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(h.ID, h.Host, h.Port, h.HostConnection, h.PortConnection, h.AdminUser, h.AdminPassword)
	if err != nil {
		return err
	}
	return nil
}

func (dir *PostgresDatabaseInstanceStorage) insertDatabaseInstance(di *entity.DatabaseInstance) error {
	stmt, err := dir.Uow.transaction.Prepare("INSERT INTO database_instances (id, name, host_connection_info_id, ecosystem_id, database_technology_id, enabled, note, created_at, created_by_user_id, updated_at, connection_status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(di.ID, di.Name, di.HostConnection.ID, di.EcosystemID, di.DatabaseTechnologyID, di.Enabled, di.Note, di.CreatedAt, di.CreatedByUserID, di.UpdatedAt, di.ConnectionStatus)
	if err != nil {
		return err
	}
	return nil
}

func (dir *PostgresDatabaseInstanceStorage) updateHostConnectionInfo(h *entity.HostConnectionInfo) error {
	stmt, err := dir.Uow.transaction.Prepare("UPDATE host_connection_info SET host = $2, port = $3, host_connection = $4, port_connection = $5, admin_username = $6, admin_password = $7 WHERE id = $1")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(h.ID, h.Host, h.Port, h.HostConnection, h.PortConnection, h.AdminUser, h.AdminPassword)
	if err != nil {
		return err
	}
	return nil
}

func (dir *PostgresDatabaseInstanceStorage) updateDatabaseInstance(di *entity.DatabaseInstance) error {
	stmt, err := dir.Uow.transaction.Prepare(`
UPDATE database_instances 
SET name = $2, 
	ecosystem_id = $3, 
	database_technology_id = $4, 
	enabled = $5, 
	note = $6,
	updated_at = $7, 
	disabled_at = $8
WHERE id = $1`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(di.ID, di.Name, di.EcosystemID, di.DatabaseTechnologyID, di.Enabled, di.Note, di.UpdatedAt, di.DisabledAt)
	if err != nil {
		return err
	}
	return nil

}

func fetchDataFromRows(rows *sql.Rows) ([]*dto.DatabaseInstanceOutputDTO, error) {
	var databaseInstances []*dto.DatabaseInstanceOutputDTO
	defer rows.Close()
	for rows.Next() {
		dbInstance, err := scanDatabaseInstance(rows)
		if err != nil {
			return nil, err
		}
		databaseInstances = append(databaseInstances, &dbInstance)
	}
	return databaseInstances, nil
}

func scanDatabaseInstance(rows *sql.Rows) (dto.DatabaseInstanceOutputDTO, error) {
	var dbInstance dto.DatabaseInstanceOutputDTO
	err := rows.Scan(&dbInstance.ID,
		&dbInstance.Name,
		&dbInstance.Host,
		&dbInstance.Port,
		&dbInstance.HostConnection,
		&dbInstance.PortConnection,
		&dbInstance.AdminUser,
		&dbInstance.AdminPassword,
		&dbInstance.EcosystemID,
		&dbInstance.EcosystemName,
		&dbInstance.DatabaseTechnologyID,
		&dbInstance.DatabaseTechnologyName,
		&dbInstance.DatabaseTechnologyVersion,
		&dbInstance.Enabled,
		&dbInstance.RolesCreated,
		&dbInstance.Note,
		&dbInstance.CreatedAt,
		&dbInstance.CreatedByUserID,
		&dbInstance.CreatedByUser,
		&dbInstance.UpdatedAt,
		&dbInstance.DisabledAt,
		&dbInstance.LastDatabaseSync,
		&dbInstance.ConnectionStatus,
		&dbInstance.LastConnectionTest,
		&dbInstance.LastConnectionResult)
	return dbInstance, err
}

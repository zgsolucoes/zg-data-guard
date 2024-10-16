package connector

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/lib/pq"

	"github.com/zgsolucoes/zg-data-guard/internal/database/storage"

	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	"github.com/zgsolucoes/zg-data-guard/internal/entity"
)

const (
	postgres          = "postgres"
	connectionTimeout = 30 * time.Second
	maxRetries        = 3
	retryInterval     = 500 * time.Millisecond
	sqlFilePath       = "internal/database/connector/scripts/postgres"
)

var (
	ErrWhileExecutingStatementPostgres = errors.New("error occurred while attempting to execute the statement on the target PostgreSQL instance")
)

type PostgresConnector struct {
	ConnectionData dto.ConnectionInputDTO
}

func newPostgresConnector(connectionData dto.ConnectionInputDTO) *PostgresConnector {
	return &PostgresConnector{ConnectionData: connectionData}
}

func (pc *PostgresConnector) TestConnection() error {
	return pc.executeWithTimeout(context.Background(), func(ctx context.Context, db *sql.DB) error {
		return db.PingContext(ctx)
	})
}

// CreateRoles godoc
// Create Data Guard roles in the database instance and setup zgbd user
func (pc *PostgresConnector) CreateRoles(roles []*DatabaseRole) error {
	return pc.executeWithTimeout(context.Background(), func(ctx context.Context, db *sql.DB) error {
		if err := pc.createRolesInDB(ctx, db, roles); err != nil {
			return err
		}
		if err := pc.allowReplicationSlotCreation(ctx, db); err != nil {
			return err
		}
		if err := pc.grantApplicationRole(ctx, db); err != nil {
			return err
		}
		return nil
	})
}

func (pc *PostgresConnector) createRolesInDB(ctx context.Context, db *sql.DB, roles []*DatabaseRole) error {
	sqlFilePath := filepath.Join(sqlFilePath, "create_role_if_not_exists.sql")
	createRoleTemplate, err := storage.ReadSQLFile(sqlFilePath)
	if err != nil {
		return err
	}
	for _, role := range roles {
		createRoleFunction := fmt.Sprintf(createRoleTemplate, role.Name)
		if _, err := db.ExecContext(ctx, createRoleFunction); err != nil {
			return err
		}
	}
	return nil
}

func (pc *PostgresConnector) allowReplicationSlotCreation(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, "ALTER ROLE zgbd REPLICATION")
	return err
}

func (pc *PostgresConnector) grantApplicationRole(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, "GRANT application TO zgbd")
	return err
}

// SetupGrantsToRoles godoc
// Execute a procedure that grants privileges to roles in the all schemas of the database
// It's necessary to execute this procedure after creating the roles to grant the necessary privileges to them
// The database user will have one of this roles
func (pc *PostgresConnector) SetupGrantsToRoles() error {
	return pc.executeWithTimeout(context.Background(), func(ctx context.Context, db *sql.DB) error {
		sqlFilePath := filepath.Join(sqlFilePath, "setup_grants_roles_database.sql")
		setupFunction, err := storage.ReadSQLFile(sqlFilePath)
		if err != nil {
			return err
		}
		_, err = db.ExecContext(ctx, setupFunction)
		if err != nil {
			return err
		}
		return nil
	})
}

func (pc *PostgresConnector) ListDatabases() ([]*Database, error) {
	dbConn, err := sql.Open(pc.Driver(), pc.URL())
	if err != nil {
		return nil, err
	}
	defer dbConn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), connectionTimeout)
	defer cancel()
	errorChan := make(chan error, 1)
	databasesChan := make(chan []*Database, 1)
	go listDatabasesWithTimeout(ctx, dbConn, errorChan, databasesChan)

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("connection failed! Cause: connection timed out after %s", connectionTimeout.String())
	case err := <-errorChan:
		return nil, err
	case databases := <-databasesChan:
		return databases, nil
	}
}

func (pc *PostgresConnector) UserExists(username string) (bool, error) {
	result, err := pc.queryWithTimeout(context.Background(), func(ctx context.Context, db *sql.DB) (any, error) {
		query := `SELECT EXISTS(SELECT 1 FROM pg_roles WHERE rolname=$1)`

		var exists bool
		err := db.QueryRowContext(ctx, query, username).Scan(&exists)
		if err != nil {
			return nil, err
		}
		return exists, nil
	})

	if err != nil {
		return false, err
	}

	return result.(bool), nil
}

func (pc *PostgresConnector) CreateUser(user *DatabaseUser) error {
	return pc.executeWithTimeout(context.Background(), func(ctx context.Context, db *sql.DB) error {
		var err error
		if entity.CheckRoleApplication(user.Role) {
			// Application role requires password encryption to be set to 'md5' for compatibility purposes
			if _, err = db.ExecContext(ctx, `SET password_encryption = 'md5'`); err != nil {
				return err
			}
		}
		_, err = db.ExecContext(ctx, fmt.Sprintf(`CREATE USER "%s" WITH LOGIN PASSWORD '%s' IN ROLE "%s"`, user.Username, user.Password, user.Role))
		return err
	})
}

func (pc *PostgresConnector) GrantConnect(username string) error {
	return pc.executeWithTimeout(context.Background(), func(ctx context.Context, db *sql.DB) error {
		stmt := fmt.Sprintf(`GRANT CONNECT ON DATABASE "%s" TO "%s"`, pc.Database(), username)
		var err error
		for i := 0; i < maxRetries; i++ {
			_, err = db.ExecContext(ctx, stmt)
			if err != nil && isConcurrentError(err) {
				log.Printf("Concurrent error detected on grant connect to '%s' in '%s', retrying... [%d/%d]", username, pc.Database(), i+1, maxRetries)
				time.Sleep(retryInterval)
				continue
			}
			break
		}
		return err
	})
}

// RevokeUserPrivilegesAndRemove godoc
// Revokes all privileges from a user and removes it from the database
// It's necessary to revoke all privileges before removing the user, but if the user owns objects, it's necessary to transfer ownership to another user before removing it.
// The function returns an error if the user doesn't exist or if the user owns objects.
func (pc *PostgresConnector) RevokeUserPrivilegesAndRemove(username string) error {
	return pc.executeWithTimeout(context.Background(), func(ctx context.Context, db *sql.DB) error {
		sqlFilePath := filepath.Join(sqlFilePath, "revoke_user_privileges_and_exclude.sql")
		removeUserFunc, err := storage.ReadSQLFile(sqlFilePath)
		if err != nil {
			return err
		}
		stmt := fmt.Sprintf(removeUserFunc, username)
		_, err = db.ExecContext(ctx, stmt)
		if err != nil {
			return err
		}

		userStillExists, err := pc.UserExists(username)
		if err != nil {
			return err
		}
		if userStillExists {
			return fmt.Errorf("user '%s' still exists after removal attempt", username)
		}
		return nil
	})
}

func (pc *PostgresConnector) Driver() string {
	return "postgres"
}

func (pc *PostgresConnector) DefaultDatabase() string {
	return "postgres"
}

func (pc *PostgresConnector) Database() string {
	databaseName := pc.ConnectionData.Database
	if databaseName == "" {
		databaseName = pc.DefaultDatabase()
	}
	return databaseName
}

func (pc *PostgresConnector) URL() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		url.QueryEscape(pc.ConnectionData.User),
		url.QueryEscape(pc.ConnectionData.Password),
		pc.ConnectionData.Host,
		pc.ConnectionData.Port,
		pc.Database(),
	)
}

func (pc *PostgresConnector) executeWithTimeout(ctx context.Context, operation func(context.Context, *sql.DB) error) error {
	dbConn, err := sql.Open(pc.Driver(), pc.URL())
	if err != nil {
		return err
	}
	defer dbConn.Close()

	ctx, cancel := context.WithTimeout(ctx, connectionTimeout)
	defer cancel()

	result := make(chan error, 1)
	go func() { result <- operation(ctx, dbConn) }()

	select {
	case <-ctx.Done():
		log.Printf("Connection failed! Cause: connection timed out after %s", connectionTimeout)
		return fmt.Errorf("connection failed! Cause: connection timed out after %s", connectionTimeout)
	case err := <-result:
		if err != nil {
			log.Printf("Error while executing statement on database %s. Cause: %v", pc.Database(), err)
			return fmt.Errorf("%w! Database: %s. Cause: %w", ErrWhileExecutingStatementPostgres, pc.Database(), err)
		}
		return nil
	}
}

func (pc *PostgresConnector) queryWithTimeout(ctx context.Context, operation func(context.Context, *sql.DB) (any, error)) (any, error) {
	dbConn, err := sql.Open(pc.Driver(), pc.URL())
	if err != nil {
		return nil, err
	}
	defer dbConn.Close()

	ctx, cancel := context.WithTimeout(ctx, connectionTimeout)
	defer cancel()

	result := make(chan any, 1)
	errChan := make(chan error, 1)
	go func() {
		res, err := operation(ctx, dbConn)
		if err != nil {
			errChan <- err
		} else {
			result <- res
		}
	}()

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("connection failed! Cause: connection timed out after %s", connectionTimeout)
	case err := <-errChan:
		return nil, err
	case res := <-result:
		return res, nil
	}
}

func listDatabasesWithTimeout(ctx context.Context, db *sql.DB, errorChan chan<- error, databasesChan chan<- []*Database) {
	rows, err := db.QueryContext(ctx, "SELECT datname AS name, PG_SIZE_PRETTY(PG_DATABASE_SIZE(datname)) AS current_size FROM pg_database WHERE datistemplate = FALSE AND datname != 'postgres'")
	if err != nil {
		errorChan <- err
		return
	}
	defer rows.Close()

	var databases []*Database
	for rows.Next() {
		var database Database
		err := rows.Scan(&database.Name, &database.CurrentSize)
		if err != nil {
			errorChan <- err
			return
		}
		databases = append(databases, &database)
	}
	databasesChan <- databases
}

func isConcurrentError(err error) bool {
	return strings.Contains(err.Error(), "tuple concurrently updated")
}

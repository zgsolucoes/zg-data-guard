package config

import (
	"database/sql"
	"errors"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

const migrationDirectory = "file://./internal/database/migrations"

func runMigrations(databaseName string, dbConn *sql.DB) {
	// Starting run migrations
	log.Println("Running migrations...")
	driver, err := postgres.WithInstance(dbConn, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		migrationDirectory,
		databaseName, driver)
	if err != nil {
		log.Fatal(err)
	}

	// Run migrations
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal(err)
	} else if errors.Is(err, migrate.ErrNoChange) {
		log.Println("No migrations to run!")
	} else {
		log.Println("Migrations ran successfully!")
	}
}

package config

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
)

func newDBConnection(databaseName string) (*sql.DB, error) {
	connURL := buildPostgresURL(databaseName)
	dbConn, err := sql.Open(getPostgresDriver(), connURL)
	if err != nil {
		return nil, err
	}
	if err = dbConn.Ping(); err != nil {
		return nil, err
	}
	log.Printf("Connected to PostgreSQL database: %s", databaseName)

	return dbConn, nil
}

func buildPostgresURL(databaseName string) string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		url.QueryEscape(os.Getenv("DATABASE_USER")),
		url.QueryEscape(os.Getenv("DATABASE_PASSWORD")),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		databaseName,
		func() string {
			if sslmode := os.Getenv("DATABASE_SSLMODE"); sslmode != "" {
				return sslmode
			}
			return "disable"
		}(),
	)
}

func getPostgresDriver() string {
	postgresDriver := os.Getenv("DATABASE_DRIVER")
	if postgresDriver == "" {
		return "postgres"
	}
	return postgresDriver
}

func closeDatabaseConnection() {
	if dbConn != nil {
		err := dbConn.Close()
		if err != nil {
			log.Fatal(err, "Error closing database connection")
		}
		log.Println("Database connection closed successfully!")
	}
}

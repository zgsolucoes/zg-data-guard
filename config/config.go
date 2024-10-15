//go:build !cover

package config

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
)

const (
	EnvDevelopment = "development"
)

var (
	dbConn  *sql.DB
	webPort string
	appURL  string
)

func Init() {
	var err error
	// Load .env file
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	initializeLogFile()
	initializeBuildData()
	dbName := os.Getenv("DATABASE_NAME")
	dbConn, err = newDBConnection(dbName)
	if err != nil {
		log.Fatal(err, " error initializing postgres")
	}

	runMigrations(dbName, dbConn)
	initializeJwt()
	initializeCryptography()
}

func Cleanup() {
	closeDatabaseConnection()
	closeLogFile(logFile)
}

func GetDBConn() *sql.DB {
	return dbConn
}

func GetWebPort() string {
	webPort = os.Getenv("WEBSERVER_PORT")
	if webPort == "" {
		return "8081"
	}
	return webPort
}

func GetEnvironment() string {
	return os.Getenv("ENVIRONMENT")
}

func GetExternalHost() string {
	return os.Getenv("EXTERNAL_HOST")
}

func GetAppName() string {
	return os.Getenv("APP_NAME")
}

func SetApplicationURL(url string) {
	appURL = url
}

func GetApplicationURL() string {
	return appURL
}

func GetAppContextPath() string {
	if GetEnvironment() == EnvDevelopment {
		return "/"
	}
	return "/" + GetAppName()
}

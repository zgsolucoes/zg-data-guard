package connector

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/zgsolucoes/zg-data-guard/config"
	"github.com/zgsolucoes/zg-data-guard/internal/dto"
)

const (
	ClusterConnectorPrefix = "[CLUSTER]"
)

var (
	ErrEmptyPasswordAfterDecrypt = errors.New("unexpected empty password for instance after decrypt")
)

type DatabaseTCPConnectorInterface interface {
	TestConnection() error
	Driver() string
	URL() string
	Database() string
	DefaultDatabase() string
	ListDatabases() ([]*Database, error)
	CreateRoles([]*DatabaseRole) error
	SetupGrantsToRoles() error
	UserExists(string) (bool, error)
	CreateUser(*DatabaseUser) error
	RevokeUserPrivilegesAndRemove(string) error
	GrantConnect(string) error
}

func NewDatabaseConnector(instanceData *dto.DatabaseInstanceOutputDTO, databaseName string) (DatabaseTCPConnectorInterface, error) {
	technologyName := strings.ToLower(instanceData.DatabaseTechnologyName)
	plainTextPasswd, err := config.GetCryptoHelper().Decrypt(instanceData.AdminPassword)
	if err != nil {
		log.Printf("Error decrypting password for database instance %s - %s. Cause: %v", instanceData.ID, instanceData.Name, err)
		return nil, err
	}
	if plainTextPasswd == "" {
		log.Printf("Unexpected empty password after decrypt for database instance %s - %s", instanceData.ID, instanceData.Name)
		return nil, ErrEmptyPasswordAfterDecrypt
	}
	connectionData := buildConnectionData(instanceData, databaseName, plainTextPasswd)
	switch {
	case strings.Contains(technologyName, postgres):
		return newPostgresConnector(connectionData), nil
	case strings.Contains(technologyName, DummyTest):
		return newDummyTestConnector(connectionData), nil
	default:
		return nil, fmt.Errorf("the database technology '%s' don't have a connector implemented", technologyName)
	}
}

func buildConnectionData(instanceData *dto.DatabaseInstanceOutputDTO, databaseName string, plainTextPasswd string) dto.ConnectionInputDTO {
	return dto.ConnectionInputDTO{
		ID:         instanceData.ID,
		Host:       instanceData.HostConnection,
		Port:       instanceData.PortConnection,
		User:       instanceData.AdminUser,
		Password:   plainTextPasswd,
		Database:   databaseName,
		Instance:   instanceData.Name,
		Ecosystem:  instanceData.EcosystemName,
		Technology: fmt.Sprintf("%s %s", instanceData.DatabaseTechnologyName, instanceData.DatabaseTechnologyVersion),
	}
}

package entity

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/zgsolucoes/zg-data-guard/config"
	"github.com/zgsolucoes/zg-data-guard/internal/dto"
)

type ConnectionStatus string

const (
	StatusOnline          ConnectionStatus = "ONLINE"
	StatusOffline         ConnectionStatus = "OFFLINE"
	StatusNotTested       ConnectionStatus = "NOT_TESTED"
	StatusDeactivated     ConnectionStatus = "DEACTIVATED"
	connectionEstablished                  = "connection established successfully!"
)

var (
	ErrInvalidHost               = errors.New("invalid host")
	ErrInvalidPort               = errors.New("invalid port")
	ErrInvalidHostConnection     = errors.New("invalid host connection")
	ErrInvalidPortConnection     = errors.New("invalid port connection")
	ErrInvalidAdminUser          = errors.New("invalid admin user")
	ErrInvalidAdminPassword      = errors.New("invalid admin password")
	ErrInvalidEcosystem          = errors.New("invalid ecosystem")
	ErrInvalidDatabaseTechnology = errors.New("invalid database technology")
	ErrInvalidConnectionStatus   = errors.New("invalid connection status")
)

type HostConnectionInfo struct {
	ID             uuid.UUID
	Host           string
	Port           string
	HostConnection string
	PortConnection string
	AdminUser      string
	AdminPassword  string
}

type DatabaseInstance struct {
	ID                   uuid.UUID
	Name                 string
	HostConnection       *HostConnectionInfo
	EcosystemID          string
	DatabaseTechnologyID string
	Enabled              bool
	RolesCreated         bool
	Note                 string
	CreatedAt            time.Time
	CreatedByUserID      string
	LastDatabaseSync     sql.NullTime
	ConnectionStatus     ConnectionStatus
	LastConnectionTest   sql.NullTime
	LastConnectionResult sql.NullString
	UpdatedAt            time.Time
	DisabledAt           sql.NullTime
}

func NewDatabaseInstance(input dto.DatabaseInstanceInputDTO, createdByUserID string) (*DatabaseInstance, error) {
	currentTime := time.Now()

	e := &DatabaseInstance{
		ID:   uuid.New(),
		Name: input.Name,
		HostConnection: &HostConnectionInfo{
			ID:             uuid.New(),
			Host:           input.Host,
			Port:           input.Port,
			HostConnection: input.HostConnection,
			PortConnection: input.PortConnection,
			AdminUser:      input.AdminUser,
			AdminPassword:  input.AdminPassword,
		},
		EcosystemID:          input.EcosystemID,
		DatabaseTechnologyID: input.DatabaseTechnologyID,
		Enabled:              true,
		RolesCreated:         false,
		Note:                 input.Note,
		ConnectionStatus:     StatusNotTested,
		CreatedAt:            currentTime,
		CreatedByUserID:      createdByUserID,
		UpdatedAt:            currentTime,
	}
	err := e.Validate()
	if err != nil {
		return nil, err
	}
	cipherPasswordHex, err := config.GetCryptoHelper().Encrypt(input.AdminPassword)
	if err != nil {
		return nil, err
	}
	e.HostConnection.AdminPassword = cipherPasswordHex
	return e, nil
}

func (dbi *DatabaseInstance) Validate() error {
	if dbi.Name == "" {
		return ErrInvalidName
	}
	if dbi.HostConnection.Host == "" {
		return ErrInvalidHost
	}
	if dbi.HostConnection.Port == "" {
		return ErrInvalidPort
	}
	if dbi.HostConnection.HostConnection == "" {
		return ErrInvalidHostConnection
	}
	if dbi.HostConnection.PortConnection == "" {
		return ErrInvalidPortConnection
	}
	if dbi.HostConnection.AdminUser == "" {
		return ErrInvalidAdminUser
	}
	if dbi.HostConnection.AdminPassword == "" {
		return ErrInvalidAdminPassword
	}
	if dbi.EcosystemID == "" {
		return ErrInvalidEcosystem
	}
	if dbi.DatabaseTechnologyID == "" {
		return ErrInvalidDatabaseTechnology
	}
	if dbi.ConnectionStatus == "" {
		return ErrInvalidConnectionStatus
	}
	if dbi.CreatedByUserID == "" {
		return ErrCreatedByUserNotInformed
	}
	return nil
}

func (dbi *DatabaseInstance) Update(updatedData dto.DatabaseInstanceInputDTO) error {
	dbi.Name = updatedData.Name
	dbi.HostConnection.Host = updatedData.Host
	dbi.HostConnection.Port = updatedData.Port
	dbi.HostConnection.HostConnection = updatedData.HostConnection
	dbi.HostConnection.PortConnection = updatedData.PortConnection
	dbi.EcosystemID = updatedData.EcosystemID
	dbi.DatabaseTechnologyID = updatedData.DatabaseTechnologyID
	dbi.Note = updatedData.Note
	dbi.UpdatedAt = time.Now()

	if updatedData.AdminUser != "" {
		dbi.HostConnection.AdminUser = updatedData.AdminUser
	}
	if updatedData.AdminPassword != "" {
		cipherPasswordHex, err := config.GetCryptoHelper().Encrypt(updatedData.AdminPassword)
		if err != nil {
			return err
		}
		dbi.HostConnection.AdminPassword = cipherPasswordHex
	}

	return dbi.Validate()
}

func (dbi *DatabaseInstance) Enable() {
	dbi.Enabled = true
	dbi.UpdatedAt = time.Now()
	dbi.DisabledAt = sql.NullTime{}
	dbi.ConnectionStatus = StatusNotTested
}

func (dbi *DatabaseInstance) Disable() {
	currentTime := time.Now()
	dbi.Enabled = false
	dbi.UpdatedAt = currentTime
	dbi.DisabledAt = sql.NullTime{Time: currentTime, Valid: true}
	dbi.ConnectionStatus = StatusDeactivated
}

func (dbi *DatabaseInstance) RefreshLastDatabaseSync() {
	dbi.LastDatabaseSync = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	dbi.RefreshLastConnectionTest(true, connectionEstablished)
}

func (dbi *DatabaseInstance) RefreshLastConnectionTest(success bool, resultMsg string) {
	dbi.LastConnectionTest = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	dbi.LastConnectionResult = sql.NullString{String: resultMsg, Valid: true}
	if success {
		dbi.turnOnline()
	} else {
		dbi.turnOffline()
	}
}

func (dbi *DatabaseInstance) CreateRoles() {
	dbi.RolesCreated = true
	dbi.RefreshLastConnectionTest(true, connectionEstablished)
	dbi.UpdatedAt = time.Now()
}

func (dbi *DatabaseInstance) turnOnline() {
	dbi.ConnectionStatus = StatusOnline
}

func (dbi *DatabaseInstance) turnOffline() {
	dbi.ConnectionStatus = StatusOffline
}

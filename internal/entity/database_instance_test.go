package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zgsolucoes/zg-data-guard/internal/dto"
)

const (
	dbName          = "PostgreSQL - Local"
	dbHost          = "localhost"
	dbPort          = "5432"
	dbHostConn      = "host.conn.ip"
	dbPortConn      = "5433"
	dbAdminUser     = "admin"
	dbAdminPassword = "pwd"
	dbEcosystemId   = "ecosystem_id"
	dbTechnologyId  = "technology_id"
	dbNote          = "note"
	userID          = "dd42cf0c-8a91-42d7-a906-cb9313494e7d"
)

var validInput = dto.DatabaseInstanceInputDTO{
	Name:                 dbName,
	Host:                 dbHost,
	Port:                 dbPort,
	HostConnection:       dbHostConn,
	PortConnection:       dbPortConn,
	AdminUser:            dbAdminUser,
	AdminPassword:        dbAdminPassword,
	EcosystemID:          dbEcosystemId,
	DatabaseTechnologyID: dbTechnologyId,
	Note:                 dbNote,
}

func TestGivenAnEmptyRequiredParam_WhenValidateDatabaseInstance_ThenShouldReceiveAnError(t *testing.T) {
	dbInstance := &DatabaseInstance{}
	assertValidate(t, dbInstance, ErrInvalidName)

	dbInstance = buildDatabaseInstance(dbName, "", "", "", "", "", "", "", "", "", "")
	assertValidate(t, dbInstance, ErrInvalidHost)

	dbInstance = buildDatabaseInstance(dbName, dbHost, "", "", "", "", "", "", "", "", "")
	assertValidate(t, dbInstance, ErrInvalidPort)

	dbInstance = buildDatabaseInstance(dbName, dbHost, dbPort, "", "", "", "", "", "", "", "")
	assertValidate(t, dbInstance, ErrInvalidHostConnection)

	dbInstance = buildDatabaseInstance(dbName, dbHost, dbPort, dbHostConn, "", "", "", "", "", "", "")
	assertValidate(t, dbInstance, ErrInvalidPortConnection)

	dbInstance = buildDatabaseInstance(dbName, dbHost, dbPort, dbHostConn, dbPortConn, "", "", "", "", "", "")
	assertValidate(t, dbInstance, ErrInvalidAdminUser)

	dbInstance = buildDatabaseInstance(dbName, dbHost, dbPort, dbHostConn, dbPortConn, dbAdminUser, "", "", "", "", "")
	assertValidate(t, dbInstance, ErrInvalidAdminPassword)

	dbInstance = buildDatabaseInstance(dbName, dbHost, dbPort, dbHostConn, dbPortConn, dbAdminUser, dbAdminPassword, "", "", "", "")
	assertValidate(t, dbInstance, ErrInvalidEcosystem)

	dbInstance = buildDatabaseInstance(dbName, dbHost, dbPort, dbHostConn, dbPortConn, dbAdminUser, dbAdminPassword, dbEcosystemId, "", "", "")
	assertValidate(t, dbInstance, ErrInvalidDatabaseTechnology)

	dbInstance = buildDatabaseInstance(dbName, dbHost, dbPort, dbHostConn, dbPortConn, dbAdminUser, dbAdminPassword, dbEcosystemId, dbTechnologyId, "", "")
	assertValidate(t, dbInstance, ErrInvalidConnectionStatus)

	dbInstance = buildDatabaseInstance(dbName, dbHost, dbPort, dbHostConn, dbPortConn, dbAdminUser, dbAdminPassword, dbEcosystemId, dbTechnologyId, "", StatusNotTested)
	assertValidate(t, dbInstance, ErrCreatedByUserNotInformed)
}

func TestGivenAValidParams_WhenValidateDatabaseInstance_ThenShouldNotReceiveAnError(t *testing.T) {
	dbInstance := buildDatabaseInstance(dbName, dbHost, dbPort, dbHostConn, dbPortConn, dbAdminUser, dbAdminPassword, dbEcosystemId, dbTechnologyId, userID, StatusNotTested)
	assert.NoError(t, dbInstance.Validate())
}

func TestGivenAnInvalidParams_WhenCreateNewDatabaseInstance_ThenShouldReturnAnError(t *testing.T) {
	invalidInput := dto.DatabaseInstanceInputDTO{
		Name: "",
		Host: dbHost,
		Port: dbPort,
	}
	db, err := NewDatabaseInstance(invalidInput, "")
	assert.Error(t, err)
	assert.EqualError(t, err, ErrInvalidName.Error())
	assert.Nil(t, db, "DatabaseInstance should be nil")
}

func TestGivenAnInstance_WhenRefreshLastConnectionTest_ThenShouldUpdateConnectionStatus(t *testing.T) {
	dbInstance := buildDatabaseInstance(dbName, dbHost, dbPort, dbHostConn, dbPortConn, dbAdminUser, dbAdminPassword, dbEcosystemId, dbTechnologyId, userID, StatusNotTested)
	assert.Equal(t, StatusNotTested, dbInstance.ConnectionStatus)
	assert.Empty(t, dbInstance.LastConnectionTest)

	const msgFailed = "connection failed!"
	dbInstance.RefreshLastConnectionTest(false, msgFailed)
	assert.Equal(t, StatusOffline, dbInstance.ConnectionStatus)
	assert.NotEmpty(t, dbInstance.LastConnectionTest)
	assert.Equal(t, msgFailed, dbInstance.LastConnectionResult.String)

	const msgSuccess = "connection established successfully!"
	dbInstance.RefreshLastConnectionTest(true, msgSuccess)
	assert.Equal(t, StatusOnline, dbInstance.ConnectionStatus)
	assert.NotEmpty(t, dbInstance.LastConnectionTest)
	assert.Equal(t, msgSuccess, dbInstance.LastConnectionResult.String)
}

func TestGivenAnInstance_WhenRefreshLastDatabaseSync_ThenShouldUpdateLastSync(t *testing.T) {
	dbInstance := buildDatabaseInstance(dbName, dbHost, dbPort, dbHostConn, dbPortConn, dbAdminUser, dbAdminPassword, dbEcosystemId, dbTechnologyId, userID, StatusNotTested)
	assert.Empty(t, dbInstance.LastDatabaseSync)

	dbInstance.RefreshLastDatabaseSync()
	assert.NotEmpty(t, dbInstance.LastDatabaseSync)
}

func TestGivenAValidParams_WhenCreateNewDatabaseInstance_ThenShouldReturnADatabaseInstance(t *testing.T) {
	dbInstance, err := NewDatabaseInstance(validInput, userID)
	assert.NoError(t, err)
	assert.NotNil(t, dbInstance, "DatabaseInstance should not be nil")
	assert.NotEmpty(t, dbInstance.ID, "DatabaseInstance id should not be empty")
	assert.Equal(t, dbName, dbInstance.Name)
	assert.Equal(t, dbHost, dbInstance.HostConnection.Host)
	assert.Equal(t, dbPort, dbInstance.HostConnection.Port)
	assert.Equal(t, dbHostConn, dbInstance.HostConnection.HostConnection)
	assert.Equal(t, dbPortConn, dbInstance.HostConnection.PortConnection)
	assert.Equal(t, dbAdminUser, dbInstance.HostConnection.AdminUser)
	assert.NotEmpty(t, dbInstance.HostConnection.AdminPassword)
	assert.Equal(t, dbEcosystemId, dbInstance.EcosystemID)
	assert.Equal(t, dbTechnologyId, dbInstance.DatabaseTechnologyID)
	assert.Equal(t, true, dbInstance.Enabled)
	assert.Equal(t, dbNote, dbInstance.Note)
	assert.Equal(t, userID, dbInstance.CreatedByUserID)
	assert.Equal(t, StatusNotTested, dbInstance.ConnectionStatus)
	assert.True(t, dbInstance.Enabled)
	assert.False(t, dbInstance.RolesCreated)
	assert.NotEmpty(t, dbInstance.CreatedAt)
	assert.NotEmpty(t, dbInstance.UpdatedAt)
	assert.Equal(t, dbInstance.CreatedAt, dbInstance.UpdatedAt)
}

func TestGivenValidInput_WhenUpdateDatabaseInstance_ThenShouldUpdateDatabaseInstance(t *testing.T) {
	dbInstance, _ := NewDatabaseInstance(validInput, userID)
	updatedInput := dto.DatabaseInstanceInputDTO{
		Name:                 "PostgreSQL - Local - Updated",
		Host:                 "192.152.1.1",
		Port:                 "5555",
		HostConnection:       "host-updated.conn.ip",
		PortConnection:       "4444",
		AdminUser:            "updated_admin",
		AdminPassword:        "updated_pwd",
		EcosystemID:          "ecosystem_id_updated",
		DatabaseTechnologyID: "technology_id_updated",
		Note:                 "note updated",
	}

	err := dbInstance.Update(updatedInput)
	dbInstance.CreateRoles()

	assert.NoError(t, err)
	assert.Equal(t, updatedInput.Name, dbInstance.Name)
	assert.Equal(t, updatedInput.Host, dbInstance.HostConnection.Host)
	assert.Equal(t, updatedInput.Port, dbInstance.HostConnection.Port)
	assert.Equal(t, updatedInput.HostConnection, dbInstance.HostConnection.HostConnection)
	assert.Equal(t, updatedInput.PortConnection, dbInstance.HostConnection.PortConnection)
	assert.Equal(t, updatedInput.AdminUser, dbInstance.HostConnection.AdminUser)
	assert.NotEmpty(t, dbInstance.HostConnection.AdminPassword)
	assert.Equal(t, updatedInput.EcosystemID, dbInstance.EcosystemID)
	assert.Equal(t, updatedInput.DatabaseTechnologyID, dbInstance.DatabaseTechnologyID)
	assert.Equal(t, updatedInput.Note, dbInstance.Note)
	assert.True(t, dbInstance.Enabled)
	assert.True(t, dbInstance.RolesCreated)
	assert.NotEmpty(t, dbInstance.UpdatedAt)
	assert.Empty(t, dbInstance.DisabledAt)
	assert.NotEqual(t, dbInstance.CreatedAt, dbInstance.UpdatedAt)

	err = dbInstance.Update(updatedInput)
	assert.NoError(t, err)
	assert.Empty(t, dbInstance.DisabledAt)
}

func TestDatabaseInstance_Enable(t *testing.T) {
	dbInstance := &DatabaseInstance{Enabled: false}

	dbInstance.Enable()

	assert.True(t, dbInstance.Enabled)
	assert.WithinDuration(t, time.Now(), dbInstance.UpdatedAt, time.Second)
	assert.False(t, dbInstance.DisabledAt.Valid)
	assert.Equal(t, StatusNotTested, dbInstance.ConnectionStatus)
}

func TestDatabaseInstance_Disable(t *testing.T) {
	dbInstance := &DatabaseInstance{Enabled: true}

	dbInstance.Disable()

	assert.False(t, dbInstance.Enabled)
	assert.WithinDuration(t, time.Now(), dbInstance.UpdatedAt, time.Second)
	assert.True(t, dbInstance.DisabledAt.Valid)
	assert.WithinDuration(t, time.Now(), dbInstance.DisabledAt.Time, time.Second)
	assert.Equal(t, StatusDeactivated, dbInstance.ConnectionStatus)
}

func buildDatabaseInstance(name, host, port, hostConn, portConn, admUser, admPwd, ecosystemId, techId, userId string, status ConnectionStatus) *DatabaseInstance {
	return &DatabaseInstance{
		Name: name,
		HostConnection: &HostConnectionInfo{
			Host:           host,
			Port:           port,
			HostConnection: hostConn,
			PortConnection: portConn,
			AdminUser:      admUser,
			AdminPassword:  admPwd,
		},
		EcosystemID:          ecosystemId,
		DatabaseTechnologyID: techId,
		CreatedByUserID:      userId,
		ConnectionStatus:     status,
	}
}

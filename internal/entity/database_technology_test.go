package entity

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const (
	postgresName    = "PostgreSQL"
	postgresVersion = "15"
)

func TestGivenAnEmptyName_WhenValidateDatabaseTechnology_ThenShouldReceiveAnError(t *testing.T) {
	e := DatabaseTechnology{}
	err := e.Validate()
	assert.Error(t, err)
	assert.EqualError(t, err, ErrInvalidName.Error())

	e = DatabaseTechnology{Version: postgresVersion}
	err = e.Validate()
	assert.Error(t, err)
	assert.EqualError(t, err, ErrInvalidName.Error())
}

func TestGivenAnEmptyVersion_WhenValidateDatabaseTechnology_ThenShouldReceiveAnError(t *testing.T) {
	u := DatabaseTechnology{Name: postgresName}
	err := u.Validate()
	assert.Error(t, err)
	assert.EqualError(t, err, ErrInvalidVersion.Error())
}

func TestGivenAnEmptyCreatedByUser_WhenValidateDatabaseTechnology_ThenShouldReceiveAnError(t *testing.T) {
	u := DatabaseTechnology{Name: postgresName, Version: postgresVersion}
	err := u.Validate()
	assert.Error(t, err)
	assert.EqualError(t, err, ErrCreatedByUserNotInformed.Error())
}

func TestGivenAValidParams_WhenValidateDatabaseTechnology_ThenShouldNotReceiveAnError(t *testing.T) {
	e := DatabaseTechnology{Name: postgresName, Version: postgresVersion, CreatedByUserID: uuid.New().String()}
	assert.NoError(t, e.Validate())
}

func TestGivenAnInvalidParams_WhenCreateNewDatabaseTechnology_ThenShouldReturnAnError(t *testing.T) {
	name := ""
	version := postgresVersion
	createdBy := uuid.New().String()
	e, err := NewDatabaseTechnology(name, version, createdBy)
	assert.Error(t, err)
	assert.EqualError(t, err, ErrInvalidName.Error())
	assert.Nil(t, e, "DatabaseTechnology should be nil")
}

func TestGivenAValidParams_WhenCreateNewDatabaseTechnology_ThenShouldReturnADatabaseTechnology(t *testing.T) {
	createdBy := uuid.New().String()
	e, err := NewDatabaseTechnology(postgresName, postgresVersion, createdBy)
	assert.NoError(t, err)
	assert.NotNil(t, e, "DatabaseTechnology should not be nil")
	assert.NotEmpty(t, e.ID, "DatabaseTechnology id should not be empty")
	assert.Equal(t, e.Name, postgresName)
	assert.Equal(t, e.Version, postgresVersion)
	assert.Equal(t, e.CreatedByUserID, createdBy)
	assert.NotEmpty(t, e.CreatedAt)
	assert.NotEmpty(t, e.UpdatedAt)
	assert.Equal(t, e.CreatedAt, e.UpdatedAt)
}

func TestGivenAnDatabaseTechnology_WhenUpdate_ThenShouldBeUpdated(t *testing.T) {
	e, err := NewDatabaseTechnology(postgresName, postgresVersion, uuid.New().String())
	assert.NoError(t, err)

	newName := "Postgres"
	newVersion := "16"
	e.Update(newName, newVersion)

	assert.Equal(t, newName, e.Name)
	assert.Equal(t, newVersion, e.Version)
	assert.NotEmpty(t, e.UpdatedAt)
	assert.NotEqual(t, e.CreatedAt, e.UpdatedAt)
}

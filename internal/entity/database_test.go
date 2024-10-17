package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	databaseName         = "zg-data-guard"
	dbDesc               = "ZG Data Guard Database"
	dbCurrentSize        = "15MB"
	dbDatabaseInstanceId = "dd42cf0c-8a91-42d7-a906-cb9313494e7d"
)

func TestGivenAnEmptyRequiredParam_WhenValidateDatabase_ThenShouldReceiveAnError(t *testing.T) {
	db := &Database{}
	assertValidate(t, db, ErrInvalidName)

	db = buildDatabase(databaseName, "", "")
	assertValidate(t, db, ErrDatabaseInstanceIDNotInformed)

	db = buildDatabase(databaseName, dbDatabaseInstanceId, "")
	assertValidate(t, db, ErrCreatedByUserNotInformed)
}

func TestGivenAValidParams_WhenValidateDatabase_ThenShouldNotReceiveAnError(t *testing.T) {
	db := buildDatabase(databaseName, dbDatabaseInstanceId, userID)
	assert.NoError(t, db.Validate())
}

func TestGivenAnInvalidParams_WhenCreateNewDatabase_ThenShouldReturnAnError(t *testing.T) {
	db, err := NewDatabase("", "", "", "", "")
	assert.Error(t, err)
	assert.EqualError(t, err, ErrInvalidName.Error())
	assert.Nil(t, db, "Database should be nil")
}

func TestGivenAValidParams_WhenCreateNewDatabase_ThenShouldReturnADatabase(t *testing.T) {
	db, err := NewDatabase(databaseName, dbDesc, dbDatabaseInstanceId, dbCurrentSize, userID)
	assert.NoError(t, err)
	assert.NotNil(t, db, "Database should not be nil")
	assert.NotEmpty(t, db.ID, "Database id should not be empty")
	assert.Equal(t, databaseName, db.Name)
	assert.Equal(t, dbDesc, db.Description)
	assert.Equal(t, dbDatabaseInstanceId, db.DatabaseInstanceID)
	assert.Equal(t, dbCurrentSize, db.CurrentSize)
	assert.Equal(t, true, db.Enabled)
	assert.Equal(t, userID, db.CreatedByUserID)
	assert.True(t, db.Enabled)
	assert.NotEmpty(t, db.CreatedAt)
	assert.NotEmpty(t, db.UpdatedAt)
	assert.Equal(t, db.CreatedAt, db.UpdatedAt)
	assert.Empty(t, db.DisabledAt)
}

func TestGivenValidParams_WhenUpdateDatabase_ThenShouldUpdateDatabase(t *testing.T) {
	db, _ := NewDatabase(databaseName, dbDesc, dbDatabaseInstanceId, dbCurrentSize, userID)
	updatedSize := "20MB"

	db.Update(updatedSize)

	assert.Equal(t, databaseName, db.Name)
	assert.Equal(t, dbDesc, db.Description)
	assert.Equal(t, updatedSize, db.CurrentSize)
	assert.True(t, db.Enabled)
	assert.NotEmpty(t, db.UpdatedAt)
	assert.NotEqual(t, db.CreatedAt, db.UpdatedAt)
	assert.Empty(t, db.DisabledAt)
}

func TestGivenEnabledDatabase_WhenEnableDatabase_ThenShouldEnableDatabase(t *testing.T) {
	db, _ := NewDatabase(databaseName, dbDesc, dbDatabaseInstanceId, dbCurrentSize, userID)
	db.Enabled = false

	db.Enable()

	assert.True(t, db.Enabled)
	assert.NotEmpty(t, db.UpdatedAt)
	assert.Empty(t, db.DisabledAt)
}

func TestGivenEnabledDatabase_WhenDisableDatabase_ThenShouldDisableDatabase(t *testing.T) {
	db, _ := NewDatabase(databaseName, dbDesc, dbDatabaseInstanceId, dbCurrentSize, userID)

	db.Disable()

	assert.False(t, db.Enabled)
	assert.NotEmpty(t, db.UpdatedAt)
	assert.NotEmpty(t, db.DisabledAt.Time)
	assert.True(t, db.DisabledAt.Valid)
	assert.Equal(t, db.UpdatedAt, db.DisabledAt.Time)
}

func TestGivenValidParams_WhenConfigureRoles_ThenShouldConfigureRoles(t *testing.T) {
	db, _ := NewDatabase(databaseName, dbDesc, dbDatabaseInstanceId, dbCurrentSize, userID)

	db.ConfigureRoles()

	assert.True(t, db.RolesConfigured)
	assert.NotEmpty(t, db.UpdatedAt)
}

func buildDatabase(name, databaseInstanceId, userId string) *Database {
	return &Database{
		Name:               name,
		DatabaseInstanceID: databaseInstanceId,
		CreatedByUserID:    userId,
	}
}

package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	dbUserName     = "Foo Bar"
	dbUserEmail    = "foo.bar@email.com"
	dbUserUsername = "foo.bar"
	dbUserPassword = "passwd"
	dbUserTeam     = "team"
	dbUserPosition = "developer"
	dbRoleID       = "bb1bcb7c-aac0-4f4c-ac1c-eb326b47f588"
)

func TestGivenAnEmptyRequiredParam_WhenValidateDatabaseUser_ThenShouldReceiveAnError(t *testing.T) {
	dbUser := &DatabaseUser{}
	assertValidate(t, dbUser, ErrInvalidName)

	dbUser = buildDatabaseUser(dbUserName, "", "", "", "", "")
	assertValidate(t, dbUser, ErrInvalidEmail)

	dbUser = buildDatabaseUser(dbUserName, "invalid-email.com", "", "", "", "")
	assertValidate(t, dbUser, ErrInvalidEmail)

	dbUser = buildDatabaseUser(dbUserName, dbUserEmail, "", "", "", "")
	assertValidate(t, dbUser, ErrInvalidUsername)

	dbUser = buildDatabaseUser(dbUserName, dbUserEmail, dbUserUsername, "", "", "")
	assertValidate(t, dbUser, ErrInvalidPassword)

	dbUser = buildDatabaseUser(dbUserName, dbUserEmail, dbUserUsername, dbUserPassword, "", "")
	assertValidate(t, dbUser, ErrDatabaseRoleIDNotInformed)

	dbUser = buildDatabaseUser(dbUserName, dbUserEmail, dbUserUsername, dbUserPassword, dbRoleID, "")
	assertValidate(t, dbUser, ErrCreatedByUserNotInformed)
}

func TestGivenAValidParams_WhenValidateDatabaseUser_ThenShouldNotReceiveAnError(t *testing.T) {
	db := buildDatabaseUser(dbUserName, dbUserEmail, dbUserUsername, dbUserPassword, dbRoleID, userID)
	assert.NoError(t, db.Validate())
}

func TestGivenAnInvalidParams_WhenCreateNewDatabaseUser_ThenShouldReturnAnError(t *testing.T) {
	db, err := NewDatabaseUser("", "", "", "", "", "")
	assert.Error(t, err)
	assert.EqualError(t, err, ErrInvalidName.Error())
	assert.Nil(t, db, "Database user should be nil")
}

func TestGivenAValidParams_WhenCreateNewDatabaseUser_ThenShouldReturnADatabaseUser(t *testing.T) {
	db, err := NewDatabaseUser(dbUserName, dbUserEmail, dbUserTeam, dbUserPosition, dbRoleID, userID)
	assert.NoError(t, err)
	assert.NotNil(t, db, "Database user should not be nil")
	assert.NotEmpty(t, db.ID, "Database user id should not be empty")
	assert.NotEmpty(t, db.Password, "Database user id should not be empty")
	assert.Equal(t, dbUserName, db.Name)
	assert.Equal(t, dbUserEmail, db.Email)
	assert.Equal(t, dbUserUsername, db.Username)
	assert.Equal(t, dbUserTeam, db.Team)
	assert.Equal(t, dbUserPosition, db.Position)
	assert.Equal(t, dbRoleID, db.DatabaseRoleID)
	assert.Equal(t, true, db.Enabled)
	assert.Equal(t, userID, db.CreatedByUserID)
	assert.True(t, db.Enabled)
	assert.NotEmpty(t, db.CreatedAt)
	assert.NotEmpty(t, db.UpdatedAt)
	assert.Equal(t, db.CreatedAt, db.UpdatedAt)
	assert.Empty(t, db.DisabledAt)
}

func TestGivenAValidParams_WhenCreateNewDatabaseUser_ThenShouldLowerCaseUsername(t *testing.T) {
	const name = "LUIZ HENRIQUE"
	db, err := NewDatabaseUser(name, "LUIZHENRIQUE@EMAIL.COM", dbUserTeam, dbUserPosition, dbRoleID, userID)
	assert.NoError(t, err)
	assert.NotNil(t, db, "Database user should not be nil")
	assert.NotEmpty(t, db.ID, "Database user id should not be empty")
	assert.NotEmpty(t, db.Password, "Database user id should not be empty")
	assert.Equal(t, name, db.Name)
	assert.Equal(t, "luizhenrique@email.com", db.Email)
	assert.Equal(t, "luizhenrique", db.Username)
	assert.Equal(t, dbUserTeam, db.Team)
	assert.Equal(t, dbUserPosition, db.Position)
	assert.Equal(t, dbRoleID, db.DatabaseRoleID)
	assert.Equal(t, true, db.Enabled)
	assert.Equal(t, userID, db.CreatedByUserID)
	assert.True(t, db.Enabled)
	assert.NotEmpty(t, db.CreatedAt)
	assert.NotEmpty(t, db.UpdatedAt)
	assert.Equal(t, db.CreatedAt, db.UpdatedAt)
	assert.Empty(t, db.DisabledAt)
}

func TestGivenValidParams_WhenUpdateDatabaseUser_ThenShouldUpdateDatabaseUser(t *testing.T) {
	dbUser, _ := NewDatabaseUser(dbUserName, dbUserEmail, dbUserTeam, dbUserPosition, dbRoleID, userID)
	updatedName := "Bar Foo"
	updatedRoleID := "84b91cf3-9b93-4b9f-b153-ae09a8ecc7aa"
	updatedTeam := "team2"
	updatedPosition := "tester"

	err := dbUser.Update(updatedName, updatedRoleID, updatedTeam, updatedPosition)

	assert.NoError(t, err)
	assert.Equal(t, updatedName, dbUser.Name)
	assert.Equal(t, updatedRoleID, dbUser.DatabaseRoleID)
	assert.Equal(t, updatedTeam, dbUser.Team)
	assert.Equal(t, updatedPosition, dbUser.Position)
	assert.Equal(t, dbUserEmail, dbUser.Email)
	assert.Equal(t, dbUserUsername, dbUser.Username)
	assert.True(t, dbUser.Enabled)
	assert.NotEmpty(t, dbUser.Password)
	assert.NotEmpty(t, dbUser.UpdatedAt)
	assert.NotEqual(t, dbUser.CreatedAt, dbUser.UpdatedAt)
	assert.Empty(t, dbUser.DisabledAt)
}

func TestGivenEnabledDatabase_WhenEnableDatabaseUser_ThenShouldEnableDatabaseUser(t *testing.T) {
	db, _ := NewDatabaseUser(dbUserName, dbUserEmail, dbUserTeam, dbUserPosition, dbRoleID, userID)
	db.Enabled = false

	db.Enable()

	assert.True(t, db.Enabled)
	assert.NotEmpty(t, db.UpdatedAt)
	assert.Empty(t, db.DisabledAt)
}

func TestGivenEnabledDatabase_WhenDisableDatabaseUser_ThenShouldDisableDatabaseUser(t *testing.T) {
	db, _ := NewDatabaseUser(dbUserName, dbUserEmail, dbUserTeam, dbUserPosition, dbRoleID, userID)

	db.Disable()

	assert.False(t, db.Enabled)
	assert.NotEmpty(t, db.UpdatedAt)
	assert.NotEmpty(t, db.DisabledAt.Time)
	assert.True(t, db.DisabledAt.Valid)
	assert.Equal(t, db.UpdatedAt, db.DisabledAt.Time)
}

func TestGivenAnEncryptedPassword_WhenDecrypt_ThenShouldDecryptPassword(t *testing.T) {
	db, _ := NewDatabaseUser(dbUserName, dbUserEmail, dbUserTeam, dbUserPosition, dbRoleID, userID)
	passwordOnCreate := db.Password

	err := db.DecryptPassword()

	assert.NoError(t, err)
	assert.NotEmpty(t, db.Password)
	assert.NotEmpty(t, db.CipherPassword)
	assert.Equal(t, passwordOnCreate, db.Password)
}

func TestGivenAnInvalidEmail_WhenGenerateUsername_ThenShouldGenerateEmptyUsername(t *testing.T) {
	generatedUsername := generateUsername("invalid-email.com")

	assert.Empty(t, generatedUsername)
}

func buildDatabaseUser(name, email, username, password, roleID, userID string) *DatabaseUser {
	return &DatabaseUser{
		Name:            name,
		Email:           email,
		Username:        username,
		Password:        password,
		DatabaseRoleID:  roleID,
		CreatedByUserID: userID,
	}
}

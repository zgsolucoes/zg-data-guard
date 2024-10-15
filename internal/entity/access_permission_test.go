package entity

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const (
	databaseID      = "1"
	databaseUserID  = "2"
	grantedByUserID = "3"
)

func TestGivenAnEmptyRequiredParam_WhenValidateAccessPermission_ThenShouldReceiveAnError(t *testing.T) {
	a := &AccessPermission{}
	assertValidate(t, a, ErrDatabaseIDNotInformed)

	a = &AccessPermission{DatabaseID: databaseID}
	assertValidate(t, a, ErrDatabaseUserIDNotInformed)

	a = &AccessPermission{DatabaseID: databaseID, DatabaseUserID: databaseUserID}
	assertValidate(t, a, ErrGrantedByUserIDNotInformed)
}

func TestGivenValidParams_WhenValidateAccessPermission_ThenShouldNotReceiveAnError(t *testing.T) {
	a := &AccessPermission{DatabaseID: databaseID, DatabaseUserID: databaseUserID, GrantedByUserID: grantedByUserID}
	assert.NoError(t, a.Validate())
}

func TestGivenAnInvalidParams_WhenCreateNewAccessPermission_ThenShouldReturnAnError(t *testing.T) {
	a, err := NewAccessPermission(databaseID, "", "")
	assert.Error(t, err)
	assert.EqualError(t, err, ErrDatabaseUserIDNotInformed.Error())
	assert.Nil(t, a, "AccessPermission should be nil")
}

func TestGivenAValidParams_WhenCreateNewAccessPermission_ThenShouldReturnAccessPermission(t *testing.T) {
	grantedBy := uuid.New().String()
	a, err := NewAccessPermission(databaseID, databaseUserID, grantedBy)
	assert.NoError(t, err)
	assert.NotNil(t, a, "AccessPermission should not be nil")
	assert.NotEmpty(t, a.ID, "AccessPermission id should not be empty")
	assert.Equal(t, a.DatabaseID, databaseID)
	assert.Equal(t, a.DatabaseUserID, databaseUserID)
	assert.Equal(t, a.GrantedByUserID, grantedBy)
	assert.NotEmpty(t, a.GrantedAt)
}

func assertValidate(t *testing.T, entity Validator, expectedError error) {
	err := entity.Validate()
	assert.Error(t, err)
	assert.EqualError(t, err, expectedError.Error())
}

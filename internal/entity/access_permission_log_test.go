package entity

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const (
	instanceID = "1"
	message    = "test message"
)

func TestGivenAnEmptyRequiredParam_WhenValidateAccessPermissionLog_ThenShouldReceiveAnError(t *testing.T) {
	g := &AccessPermissionLog{}
	assertValidate(t, g, ErrDatabaseInstanceIDNotInformed)

	g = &AccessPermissionLog{DatabaseInstanceID: instanceID}
	assertValidate(t, g, ErrMessageNotInformed)

	g = &AccessPermissionLog{DatabaseInstanceID: instanceID, Message: message}
	assertValidate(t, g, ErrUserIDNotInformed)
}

func TestGivenValidParams_WhenValidateAccessPermissionLog_ThenShouldNotReceiveAnError(t *testing.T) {
	g := &AccessPermissionLog{DatabaseInstanceID: instanceID, Message: message, UserID: uuid.New().String()}
	assert.NoError(t, g.Validate())
}

func TestGivenAnInvalidParams_WhenCreateNewAccessPermissionLog_ThenShouldReturnAnError(t *testing.T) {
	a, err := NewAccessPermissionLog(instanceID, "", "", "", "", false)
	assert.Error(t, err)
	assert.EqualError(t, err, ErrMessageNotInformed.Error())
	assert.Nil(t, a, "AccessPermissionLog should be nil")
}

func TestGivenAValidParams_WhenCreateNewAccessPermissionLog_ThenShouldReturnAnAccessPermissionLog(t *testing.T) {
	g, err := NewAccessPermissionLog(instanceID, userID, databaseID, message, userID, true)
	assert.NoError(t, err)
	assert.NotNil(t, g, "AccessPermissionLog should not be nil")
	assert.NotEmpty(t, g.ID, "AccessPermissionLog id should not be empty")
	assert.Equal(t, g.DatabaseUserID.String, userID)
	assert.Equal(t, g.DatabaseInstanceID, instanceID)
	assert.Equal(t, g.DatabaseID.String, databaseID)
	assert.Equal(t, g.Message, message)
	assert.Equal(t, g.UserID, userID)
	assert.True(t, g.Success)
	assert.NotEmpty(t, g.Date)
}

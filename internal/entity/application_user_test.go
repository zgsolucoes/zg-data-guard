package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGivenAnEmptyName_WhenValidateUser_ThenShouldReceiveAnError(t *testing.T) {
	u := ApplicationUser{}
	err := u.Validate()
	assert.Error(t, err)
	assert.EqualError(t, err, ErrInvalidName.Error())
}

func TestGivenAnEmptyEmail_WhenValidateUser_ThenShouldReceiveAnError(t *testing.T) {
	u := ApplicationUser{Name: "Foo Bar"}
	err := u.Validate()
	assert.Error(t, err)
	assert.EqualError(t, err, ErrInvalidEmail.Error())
}

func TestGivenAValidParams_WhenValidateUser_ThenShouldNotReceiveAnError(t *testing.T) {
	u := ApplicationUser{Name: "Foo Bar", Email: "foobar@email.com"}
	assert.NoError(t, u.Validate())
}

func TestGivenAnInvalidParams_WhenCreateNewUser_ThenShouldReturnAnError(t *testing.T) {
	name := "Foo Bar"
	email := ""
	u, err := NewApplicationUser(name, email)
	assert.Error(t, err)
	assert.EqualError(t, err, ErrInvalidEmail.Error())
	assert.Nil(t, u, "user should be nil")
}

func TestGivenAValidParams_WhenCreateNewUser_ThenShouldReturnAUser(t *testing.T) {
	name := "Foo Bar"
	email := "foobar@email.com"
	u, err := NewApplicationUser(name, email)
	assert.NoError(t, err)
	assert.NotNil(t, u, "user should not be nil")
	assert.NotEmpty(t, u.ID, "user id should not be empty")
	assert.Equal(t, u.Name, name)
	assert.Equal(t, u.Email, email)
	assert.True(t, u.Enabled)
	assert.NotEmpty(t, u.CreatedAt)
	assert.NotEmpty(t, u.UpdatedAt)
	assert.Equal(t, u.CreatedAt, u.UpdatedAt)
}

func TestGivenAUser_WheDisable_ThenShouldBeDisabled(t *testing.T) {
	u, err := NewApplicationUser("Foo Bar", "foobar@email.com")
	assert.NoError(t, err)
	assert.True(t, u.Enabled)

	u.Disable()
	assert.False(t, u.Enabled)
	assert.NotEmpty(t, u.DisabledAt)
}

func TestGivenAUser_WhenEnable_ThenShouldBeEnabled(t *testing.T) {
	u, err := NewApplicationUser("Foo Bar", "foobar@email.com")
	assert.NoError(t, err)
	assert.True(t, u.Enabled)

	u.Disable()
	assert.False(t, u.Enabled)
	assert.NotEmpty(t, u.DisabledAt)

	u.Enable()
	assert.True(t, u.Enabled)
	assert.Empty(t, u.DisabledAt)
}

package entity

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const (
	fooCode        = "foo bar"
	fooDisplayName = "Foo Bar"
)

func TestGivenAnEmptyCode_WhenValidateEcosystem_ThenShouldReceiveAnError(t *testing.T) {
	e := Ecosystem{}
	err := e.Validate()
	assert.Error(t, err)
	assert.EqualError(t, err, ErrInvalidCode.Error())

	e = Ecosystem{DisplayName: fooDisplayName}
	err = e.Validate()
	assert.Error(t, err)
	assert.EqualError(t, err, ErrInvalidCode.Error())
}

func TestGivenAnEmptyDisplayName_WhenValidateEcosystem_ThenShouldReceiveAnError(t *testing.T) {
	u := Ecosystem{Code: fooCode}
	err := u.Validate()
	assert.Error(t, err)
	assert.EqualError(t, err, ErrInvalidDisplayName.Error())
}

func TestGivenAnEmptyCreatedByUser_WhenValidateEcosystem_ThenShouldReceiveAnError(t *testing.T) {
	u := Ecosystem{Code: fooCode, DisplayName: fooDisplayName}
	err := u.Validate()
	assert.Error(t, err)
	assert.EqualError(t, err, ErrCreatedByUserNotInformed.Error())
}

func TestGivenAValidParams_WhenValidateEcosystem_ThenShouldNotReceiveAnError(t *testing.T) {
	e := Ecosystem{Code: fooCode, DisplayName: fooDisplayName, CreatedByUserID: uuid.New().String()}
	assert.NoError(t, e.Validate())
}

func TestGivenAnInvalidParams_WhenCreateNewEcosystem_ThenShouldReturnAnError(t *testing.T) {
	createdBy := uuid.New().String()
	e, err := NewEcosystem("", fooDisplayName, createdBy)
	assert.Error(t, err)
	assert.EqualError(t, err, ErrInvalidCode.Error())
	assert.Nil(t, e, "Ecosystem should be nil")
}

func TestGivenAValidParams_WhenCreateNewEcosystem_ThenShouldReturnAEcosystem(t *testing.T) {
	createdBy := uuid.New().String()
	e, err := NewEcosystem(fooCode, fooDisplayName, createdBy)
	assert.NoError(t, err)
	assert.NotNil(t, e, "Ecosystem should not be nil")
	assert.NotEmpty(t, e.ID, "Ecosystem id should not be empty")
	assert.Equal(t, e.Code, fooCode)
	assert.Equal(t, e.DisplayName, fooDisplayName)
	assert.Equal(t, e.CreatedByUserID, createdBy)
	assert.NotEmpty(t, e.CreatedAt)
	assert.NotEmpty(t, e.UpdatedAt)
	assert.Equal(t, e.CreatedAt, e.UpdatedAt)
}

func TestGivenAnEcosystem_WhenUpdate_ThenShouldBeUpdated(t *testing.T) {
	e, err := NewEcosystem(fooCode, fooDisplayName, uuid.New().String())
	assert.NoError(t, err)

	newCode := "bar foo"
	newDisplayName := "Bar Foo"
	e.Update(newCode, newDisplayName)

	assert.Equal(t, newCode, e.Code)
	assert.Equal(t, newDisplayName, e.DisplayName)
	assert.NotEmpty(t, e.UpdatedAt)
	assert.NotEqual(t, e.CreatedAt, e.UpdatedAt)
}

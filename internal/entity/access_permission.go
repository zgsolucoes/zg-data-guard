package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrDatabaseIDNotInformed      = errors.New("database id not informed")
	ErrDatabaseUserIDNotInformed  = errors.New("database user id not informed")
	ErrGrantedByUserIDNotInformed = errors.New("granted by user id not informed")
)

type AccessPermission struct {
	ID              uuid.UUID
	DatabaseID      string
	DatabaseUserID  string
	GrantedByUserID string
	GrantedAt       time.Time
}

func NewAccessPermission(databaseID, databaseUserID, grantedByUserID string) (*AccessPermission, error) {
	a := &AccessPermission{
		ID:              uuid.New(),
		DatabaseID:      databaseID,
		DatabaseUserID:  databaseUserID,
		GrantedByUserID: grantedByUserID,
		GrantedAt:       time.Now(),
	}
	if err := a.Validate(); err != nil {
		return nil, err
	}
	return a, nil
}

func (a *AccessPermission) Validate() error {
	if a.DatabaseID == "" {
		return ErrDatabaseIDNotInformed
	}
	if a.DatabaseUserID == "" {
		return ErrDatabaseUserIDNotInformed
	}
	if a.GrantedByUserID == "" {
		return ErrGrantedByUserIDNotInformed
	}
	return nil
}

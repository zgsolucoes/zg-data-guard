package entity

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrUserIDNotInformed  = errors.New("user id not informed")
	ErrMessageNotInformed = errors.New("message not informed")
)

type AccessPermissionLog struct {
	ID                 uuid.UUID
	DatabaseUserID     sql.NullString
	DatabaseInstanceID string
	DatabaseID         sql.NullString
	Message            string
	Success            bool
	Date               time.Time
	UserID             string
}

func NewAccessPermissionLog(databaseInstanceID, databaseUserID, databaseID, message, operationUserID string, success bool) (*AccessPermissionLog, error) {
	databaseUserIDFormatted := sql.NullString{
		String: databaseUserID,
		Valid:  databaseUserID != "",
	}
	databaseIDFormatted := sql.NullString{
		String: databaseID,
		Valid:  databaseID != "",
	}
	g := &AccessPermissionLog{
		ID:                 uuid.New(),
		DatabaseUserID:     databaseUserIDFormatted,
		DatabaseInstanceID: databaseInstanceID,
		DatabaseID:         databaseIDFormatted,
		Message:            message,
		Date:               time.Now(),
		UserID:             operationUserID,
		Success:            success,
	}
	err := g.Validate()
	if err != nil {
		return nil, err
	}
	return g, nil
}

func (g *AccessPermissionLog) Validate() error {
	if g.DatabaseInstanceID == "" {
		return ErrDatabaseInstanceIDNotInformed
	}
	if g.Message == "" {
		return ErrMessageNotInformed
	}
	if g.UserID == "" {
		return ErrUserIDNotInformed
	}
	return nil
}

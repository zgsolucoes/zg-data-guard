package common

import (
	"database/sql"
	"errors"
)

var (
	ErrTechnologyNotFound         = errors.New("technology not found")
	ErrEcosystemNotFound          = errors.New("ecosystem not found")
	ErrDatabaseUserNotFound       = errors.New("database user not found")
	ErrNoAccessibleInstancesFound = errors.New("no accessible instances (clusters) found for the user with the provided IDs")
)

func HandleFindError(err error, entityNotFoundError error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return entityNotFoundError
	}
	return err
}

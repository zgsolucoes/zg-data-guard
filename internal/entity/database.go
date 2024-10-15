package entity

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrDatabaseInstanceIDNotInformed = errors.New("database instance id not informed")

type Database struct {
	ID                 uuid.UUID
	Name               string
	Description        string
	CurrentSize        string
	Enabled            bool
	RolesConfigured    bool
	DatabaseInstanceID string
	CreatedAt          time.Time
	CreatedByUserID    string
	UpdatedAt          time.Time
	DisabledAt         sql.NullTime
}

func NewDatabase(name, description, databaseInstanceID, currentSize, createdByUserID string) (*Database, error) {
	currentTime := time.Now()
	d := &Database{
		ID:                 uuid.New(),
		Name:               name,
		Description:        description,
		CurrentSize:        currentSize,
		Enabled:            true,
		RolesConfigured:    false,
		DatabaseInstanceID: databaseInstanceID,
		CreatedAt:          currentTime,
		CreatedByUserID:    createdByUserID,
		UpdatedAt:          currentTime,
	}
	err := d.Validate()
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (d *Database) Enable() {
	d.Enabled = true
	d.UpdatedAt = time.Now()
	d.DisabledAt = sql.NullTime{}
}

func (d *Database) Disable() {
	currentTime := time.Now()
	d.Enabled = false
	d.UpdatedAt = currentTime
	d.DisabledAt = sql.NullTime{Time: currentTime, Valid: true}
}

func (d *Database) Validate() error {
	if d.Name == "" {
		return ErrInvalidName
	}
	if d.DatabaseInstanceID == "" {
		return ErrDatabaseInstanceIDNotInformed
	}
	if d.CreatedByUserID == "" {
		return ErrCreatedByUserNotInformed
	}
	return nil
}

func (d *Database) Update(currentSize string) {
	d.CurrentSize = currentSize
	d.UpdatedAt = time.Now()
}

func (d *Database) ConfigureRoles() {
	d.RolesConfigured = true
	d.UpdatedAt = time.Now()
}

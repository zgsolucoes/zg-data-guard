package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidVersion = errors.New("invalid version")
)

type DatabaseTechnology struct {
	ID              uuid.UUID
	Name            string
	Version         string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	CreatedByUserID string
}

func NewDatabaseTechnology(name, version string, createdByID string) (*DatabaseTechnology, error) {
	currentTime := time.Now()
	e := &DatabaseTechnology{
		ID:              uuid.New(),
		Name:            name,
		Version:         version,
		CreatedAt:       currentTime,
		UpdatedAt:       currentTime,
		CreatedByUserID: createdByID,
	}
	err := e.Validate()
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (e *DatabaseTechnology) Update(name, version string) {
	e.Name = name
	e.Version = version
	e.UpdatedAt = time.Now()
}

func (e *DatabaseTechnology) Validate() error {
	if e.Name == "" {
		return ErrInvalidName
	}
	if e.Version == "" {
		return ErrInvalidVersion
	}
	if e.CreatedByUserID == "" {
		return ErrCreatedByUserNotInformed
	}
	return nil
}

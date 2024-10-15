package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidCode              = errors.New("invalid code")
	ErrInvalidDisplayName       = errors.New("invalid display name")
	ErrCreatedByUserNotInformed = errors.New("created by user not informed")
)

type Ecosystem struct {
	ID              uuid.UUID
	Code            string
	DisplayName     string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	CreatedByUserID string
}

func NewEcosystem(code, displayName string, createdByID string) (*Ecosystem, error) {
	currentTime := time.Now()
	e := &Ecosystem{
		ID:              uuid.New(),
		Code:            code,
		DisplayName:     displayName,
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

func (e *Ecosystem) Update(code string, displayName string) {
	e.Code = code
	e.DisplayName = displayName
	e.UpdatedAt = time.Now()
}

func (e *Ecosystem) Validate() error {
	if e.Code == "" {
		return ErrInvalidCode
	}
	if e.DisplayName == "" {
		return ErrInvalidDisplayName
	}
	if e.CreatedByUserID == "" {
		return ErrCreatedByUserNotInformed
	}
	return nil
}

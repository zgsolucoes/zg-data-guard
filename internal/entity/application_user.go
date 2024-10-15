package entity

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidName  = errors.New("invalid name")
	ErrInvalidEmail = errors.New("invalid e-mail")
)

type ApplicationUser struct {
	ID         uuid.UUID
	Name       string
	Email      string
	Enabled    bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DisabledAt sql.NullTime
}

func NewApplicationUser(name, email string) (*ApplicationUser, error) {
	currentTime := time.Now()

	u := &ApplicationUser{
		ID:        uuid.New(),
		Name:      name,
		Email:     email,
		Enabled:   true,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}
	err := u.Validate()
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (u *ApplicationUser) Validate() error {
	if u.Name == "" {
		return ErrInvalidName
	}
	if u.Email == "" {
		return ErrInvalidEmail
	}
	return nil
}

func (u *ApplicationUser) Enable() {
	u.Enabled = true
	u.UpdatedAt = time.Now()
	u.DisabledAt = sql.NullTime{}
}

func (u *ApplicationUser) Disable() {
	u.Enabled = false
	u.UpdatedAt = time.Now()
	u.DisabledAt = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
}

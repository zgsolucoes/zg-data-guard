package entity

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/zgsolucoes/zg-data-guard/config"
	"github.com/zgsolucoes/zg-data-guard/pkg/utils"
)

var (
	ErrInvalidUsername           = errors.New("invalid username")
	ErrInvalidPassword           = errors.New("invalid password")
	ErrDatabaseRoleIDNotInformed = errors.New("database role id not informed")
)

const passwordLength = 16

type DatabaseUser struct {
	ID              uuid.UUID
	Name            string
	Email           string
	Username        string
	Password        string
	CipherPassword  string
	DatabaseRoleID  string
	Team            string
	Position        string
	Enabled         bool
	CreatedAt       time.Time
	CreatedByUserID string
	UpdatedAt       time.Time
	DisabledAt      sql.NullTime
	Expired         bool
	ExpiresAt       sql.NullTime
}

func NewDatabaseUser(name, email, team, position, databaseRoleID, createdByUserID string) (*DatabaseUser, error) {
	currentTime := time.Now()
	d := &DatabaseUser{
		ID:              uuid.New(),
		Name:            name,
		Email:           strings.ToLower(email),
		Username:        generateUsername(email),
		Password:        utils.GenerateRandomString(passwordLength),
		Team:            team,
		Position:        position,
		DatabaseRoleID:  databaseRoleID,
		Enabled:         true,
		CreatedAt:       currentTime,
		CreatedByUserID: createdByUserID,
		UpdatedAt:       currentTime,
	}
	err := d.Validate()
	if err != nil {
		return nil, err
	}
	err = d.encryptPassword()
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (d *DatabaseUser) Validate() error {
	if d.Name == "" {
		return ErrInvalidName
	}
	if d.Email == "" || !utils.ValidEmail(d.Email) {
		return ErrInvalidEmail
	}
	if d.Username == "" {
		return ErrInvalidUsername
	}
	if d.Password == "" && d.CipherPassword == "" {
		return ErrInvalidPassword
	}
	if d.DatabaseRoleID == "" {
		return ErrDatabaseRoleIDNotInformed
	}
	if d.CreatedByUserID == "" {
		return ErrCreatedByUserNotInformed
	}
	return nil
}

func (d *DatabaseUser) Update(name, databaseRoleID, team, position string) error {
	d.Name = name
	d.DatabaseRoleID = databaseRoleID
	d.Team = team
	d.Position = position
	d.UpdatedAt = time.Now()
	return d.Validate()
}

func (d *DatabaseUser) encryptPassword() error {
	cipherPasswordHex, err := config.GetCryptoHelper().Encrypt(d.Password)
	if err != nil {
		return err
	}
	d.CipherPassword = cipherPasswordHex
	return nil
}

func (d *DatabaseUser) DecryptPassword() error {
	password, err := config.GetCryptoHelper().Decrypt(d.CipherPassword)
	if err != nil {
		return err
	}
	d.Password = password
	return nil
}

func (d *DatabaseUser) Enable() {
	d.Enabled = true
	d.UpdatedAt = time.Now()
	d.DisabledAt = sql.NullTime{}
}

func (d *DatabaseUser) Disable() {
	currentTime := time.Now()
	d.Enabled = false
	d.UpdatedAt = currentTime
	d.DisabledAt = sql.NullTime{Time: currentTime, Valid: true}
}

func generateUsername(email string) string {
	const atSign = "@"
	parts := strings.Split(email, atSign)
	if strings.Contains(email, atSign) && len(parts) > 0 {
		return strings.ToLower(parts[0])
	}
	return ""
}

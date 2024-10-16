package storage

import (
	"database/sql"

	"github.com/zgsolucoes/zg-data-guard/internal/entity"
)

type PostgresApplicationUserStorage struct {
	DB *sql.DB
}

func NewPostgresApplicationUserStorage(db *sql.DB) *PostgresApplicationUserStorage {
	return &PostgresApplicationUserStorage{
		DB: db,
	}
}

func (a *PostgresApplicationUserStorage) FindByEmail(email string) (*entity.ApplicationUser, error) {
	var user entity.ApplicationUser
	err := a.DB.QueryRow("SELECT id, name, email, enabled, created_at, updated_at, disabled_at FROM application_users WHERE email = $1", email).
		Scan(&user.ID, &user.Name, &user.Email, &user.Enabled, &user.CreatedAt, &user.UpdatedAt, &user.DisabledAt)
	return &user, err
}

func (a *PostgresApplicationUserStorage) FindByID(id string) (*entity.ApplicationUser, error) {
	var user entity.ApplicationUser
	err := a.DB.QueryRow("SELECT id, name, email, enabled, created_at, updated_at, disabled_at FROM application_users WHERE id = $1", id).
		Scan(&user.ID, &user.Name, &user.Email, &user.Enabled, &user.CreatedAt, &user.UpdatedAt, &user.DisabledAt)
	return &user, err
}

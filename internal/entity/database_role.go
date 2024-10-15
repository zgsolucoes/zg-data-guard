package entity

import (
	"time"

	"github.com/google/uuid"
)

type RoleName string

// Any new role should be added here and the corresponding role should be added to the grant script below
// internal/database/connector/scripts/postgres/setup_grants_roles_database.sql
const (
	UserRO      RoleName = "user_ro"
	Developer   RoleName = "developer"
	DevOps      RoleName = "devops"
	Application RoleName = "application"
)

type DatabaseRole struct {
	ID              uuid.UUID
	Name            RoleName
	DisplayName     string
	Description     string
	ReadOnly        bool
	CreatedAt       time.Time
	CreatedByUserID string
}

func (d *DatabaseRole) IsUserRO() bool {
	return d.Name == UserRO
}

func (d *DatabaseRole) IsDeveloper() bool {
	return d.Name == Developer
}

func (d *DatabaseRole) IsDevOps() bool {
	return d.Name == DevOps
}

func (d *DatabaseRole) IsApplication() bool {
	return d.Name == Application
}

func ValidateRoleName(currentRole string) bool {
	switch RoleName(currentRole) {
	case UserRO, Developer, DevOps, Application:
		return true
	default:
		return false
	}
}

func CheckRoleApplication(currentRole string) bool {
	return CheckRole(currentRole, Application)
}

func CheckRole(currentRole string, roleToCheck RoleName) bool {
	if !ValidateRoleName(currentRole) {
		return false
	}
	return RoleName(currentRole) == roleToCheck
}

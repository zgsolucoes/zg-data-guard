package dto

import "time"

type ApplicationUserOutputDTO struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Enabled bool   `json:"enabled,omitempty"`
}

type EcosystemOutputDTO struct {
	ID              string     `json:"id"`
	Code            string     `json:"code"`
	DisplayName     string     `json:"displayName"`
	CreatedByUserID string     `json:"createdByUserId,omitempty"`
	CreatedByUser   string     `json:"createdByUser,omitempty"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       *time.Time `json:"updatedAt,omitempty"`
}

type TechnologyOutputDTO struct {
	ID              string     `json:"id"`
	Name            string     `json:"name"`
	Version         string     `json:"version"`
	CreatedByUserID string     `json:"createdByUserId,omitempty"`
	CreatedByUser   string     `json:"createdByUser,omitempty"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       *time.Time `json:"updatedAt,omitempty"`
}

type DatabaseInstanceOutputDTO struct {
	ID                        string     `json:"id"`
	Name                      string     `json:"name"`
	Host                      string     `json:"host"`
	Port                      string     `json:"port"`
	HostConnection            string     `json:"hostConnection"`
	PortConnection            string     `json:"portConnection"`
	AdminUser                 string     `json:"adminUser,omitempty"`
	AdminPassword             string     `json:"adminPassword,omitempty"`
	EcosystemID               string     `json:"ecosystemId"`
	EcosystemName             string     `json:"ecosystemName,omitempty"`
	DatabaseTechnologyID      string     `json:"databaseTechnologyId"`
	DatabaseTechnologyName    string     `json:"databaseTechnologyName,omitempty"`
	DatabaseTechnologyVersion string     `json:"databaseTechnologyVersion,omitempty"`
	Enabled                   bool       `json:"enabled"`
	RolesCreated              bool       `json:"rolesCreated"`
	Note                      string     `json:"note"`
	CreatedByUserID           string     `json:"createdByUserId,omitempty"`
	CreatedByUser             string     `json:"createdByUser,omitempty"`
	CreatedAt                 time.Time  `json:"createdAt"`
	UpdatedAt                 *time.Time `json:"updatedAt,omitempty"`
	ConnectionStatus          string     `json:"connectionStatus"`
	LastConnectionTest        *time.Time `json:"lastConnectionTest,omitempty"`
	LastConnectionResult      *string    `json:"lastConnectionResult,omitempty"`
	LastDatabaseSync          *time.Time `json:"lastDatabaseSync,omitempty"`
	DisabledAt                *time.Time `json:"disabledAt,omitempty"`
}

type DatabaseOutputDTO struct {
	ID                        string     `json:"id"`
	Name                      string     `json:"name"`
	CurrentSize               string     `json:"currentSize"`
	DatabaseInstanceID        string     `json:"databaseInstanceId"`
	DatabaseInstanceName      string     `json:"databaseInstanceName"`
	EcosystemID               string     `json:"ecosystemId"`
	EcosystemName             string     `json:"ecosystemName"`
	DatabaseTechnologyID      string     `json:"databaseTechnologyId"`
	DatabaseTechnologyName    string     `json:"databaseTechnologyName"`
	DatabaseTechnologyVersion string     `json:"databaseTechnologyVersion"`
	Enabled                   bool       `json:"enabled"`
	RolesConfigured           bool       `json:"rolesConfigured"`
	Description               string     `json:"description"`
	CreatedByUserID           string     `json:"createdByUserId"`
	CreatedByUser             string     `json:"createdByUser"`
	CreatedAt                 time.Time  `json:"createdAt"`
	UpdatedAt                 *time.Time `json:"updatedAt"`
	LastDatabaseSync          *time.Time `json:"lastDatabaseSync,omitempty"`
	DisabledAt                *time.Time `json:"disabledAt,omitempty"`
}

type DatabaseInstanceCredentialsOutputDTO struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

type TestConnectionOutputDTO struct {
	DatabaseInstanceID string `json:"databaseInstanceId"`
	Ecosystem          string `json:"ecosystem,omitempty"`
	Instance           string `json:"instance,omitempty"`
	Technology         string `json:"technology,omitempty"`
	Success            bool   `json:"success"`
	Message            string `json:"message"`
}

type SyncDatabasesOutputDTO struct {
	DatabaseInstanceID string `json:"databaseInstanceId"`
	Ecosystem          string `json:"ecosystem,omitempty"`
	Instance           string `json:"instance,omitempty"`
	Technology         string `json:"technology,omitempty"`
	Success            bool   `json:"success"`
	Message            string `json:"message"`
	TotalDatabases     int    `json:"totalDatabases,omitempty"`
}

type PropagateRolesOutputDTO struct {
	DatabaseInstanceID string `json:"databaseInstanceId"`
	Ecosystem          string `json:"ecosystem,omitempty"`
	Instance           string `json:"instance,omitempty"`
	Technology         string `json:"technology,omitempty"`
	Success            bool   `json:"success"`
	Message            string `json:"message"`
}

type DatabaseRoleOutputDTO struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	DisplayName     string    `json:"displayName"`
	Description     string    `json:"description"`
	ReadOnly        bool      `json:"readOnly"`
	CreatedAt       time.Time `json:"createdAt"`
	CreatedByUserID string    `json:"createdByUserId"`
}

type SetupRolesOutputDTO struct {
	DatabaseID         string `json:"databaseId"`
	DatabaseName       string `json:"databaseName"`
	DatabaseInstanceID string `json:"databaseInstanceId"`
	Ecosystem          string `json:"ecosystem,omitempty"`
	Technology         string `json:"technology,omitempty"`
	Instance           string `json:"instance,omitempty"`
	Success            bool   `json:"success"`
	Message            string `json:"message"`
}

type DatabaseUserOutputDTO struct {
	ID                      string     `json:"id"`
	Name                    string     `json:"name"`
	Email                   string     `json:"email"`
	Username                string     `json:"username"`
	Password                string     `json:"password,omitempty"`
	DatabaseRoleID          string     `json:"databaseRoleId"`
	DatabaseRoleName        string     `json:"databaseRoleName,omitempty"`
	DatabaseRoleDisplayName string     `json:"databaseRoleDisplayName,omitempty"`
	Team                    string     `json:"team,omitempty"`
	Position                string     `json:"position,omitempty"`
	Enabled                 bool       `json:"enabled"`
	CreatedAt               time.Time  `json:"createdAt"`
	CreatedByUserID         string     `json:"createdByUserId"`
	CreatedByUser           string     `json:"createdByUser,omitempty"`
	UpdatedAt               *time.Time `json:"updatedAt,omitempty"`
	DisabledAt              *time.Time `json:"disabledAt,omitempty"`
}

type HealthCheckOutputDTO struct {
	ServiceName    string `json:"serviceName"`
	ServiceVersion string `json:"serviceVersion"`
	ContextPath    string `json:"contextPath"`
	BuildTime      string `json:"buildTime"`
}

type DatabaseUserCredentialsOutputDTO struct {
	User     string `json:"username"`
	Password string `json:"password"`
}

type AccessPermissionOutputDTO struct {
	ID                   string    `json:"id"`
	DatabaseUserID       string    `json:"databaseUserId"`
	DatabaseUserName     string    `json:"databaseUserName"`
	DatabaseUserEmail    string    `json:"databaseUserEmail"`
	DatabaseRoleID       string    `json:"databaseRoleId"`
	DatabaseRoleName     string    `json:"databaseRoleName"`
	EcosystemID          string    `json:"ecosystemId"`
	EcosystemName        string    `json:"ecosystemName"`
	DatabaseInstanceID   string    `json:"databaseInstanceId"`
	DatabaseInstanceName string    `json:"databaseInstanceName"`
	DatabaseID           string    `json:"databaseId"`
	DatabaseName         string    `json:"databaseName"`
	GrantedByUserID      string    `json:"grantedByUserId"`
	GrantedByUserName    string    `json:"grantedByUserName"`
	GrantedAt            time.Time `json:"grantedAt"`
}

type GrantAccessOutputDTO struct {
	HasErrors bool   `json:"hasErrors"`
	Message   string `json:"message"`
}

type RevokeAccessOutputDTO struct {
	HasErrors bool   `json:"hasErrors"`
	Message   string `json:"message"`
}

type AccessPermissionLogOutputDTO struct {
	ID                   string    `json:"id"`
	DatabaseUserID       *string   `json:"databaseUserId,omitempty"`
	DatabaseUserName     *string   `json:"databaseUserName,omitempty"`
	DatabaseUserEmail    *string   `json:"databaseUserEmail,omitempty"`
	DatabaseInstanceID   string    `json:"databaseInstanceId"`
	DatabaseInstanceName string    `json:"databaseInstanceName"`
	Message              string    `json:"message"`
	Success              bool      `json:"success"`
	DatabaseID           *string   `json:"databaseId,omitempty"`
	DatabaseName         *string   `json:"databaseName,omitempty"`
	OperationUserID      string    `json:"operationUserId"`
	OperationUserName    string    `json:"operationUserName"`
	Date                 time.Time `json:"date"`
}

type ChangeStatusOutputDTO struct {
	ID         string     `json:"id"`
	Enabled    bool       `json:"enabled"`
	UpdatedAt  time.Time  `json:"updatedAt"`
	DisabledAt *time.Time `json:"disabledAt,omitempty"`
}

package dto

import (
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/zgsolucoes/zg-data-guard/pkg/utils"
)

const (
	emptyString = ""
	typeString  = "string"
	typeUUID    = "UUID"
	typeBoolean = "boolean"
)

var (
	ErrArrayDatabaseUsersIdsEmpty = errors.New("param: databaseUsersIds (type: []string) cannot be empty")
	ErrArrayInstancesDataEmpty    = errors.New("param: instancesData (type: []InstanceDataDTO) cannot be empty")
)

type InputValidator interface {
	Validate() error
}

type EcosystemInputDTO struct {
	Code        string `json:"code"`
	DisplayName string `json:"displayName"`
}

func (e *EcosystemInputDTO) Validate() error {
	if e.Code == emptyString {
		return errParamIsRequired("code", typeString)
	}
	if e.DisplayName == emptyString {
		return errParamIsRequired("displayName", typeString)
	}
	return nil
}

type TechnologyInputDTO struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func (t *TechnologyInputDTO) Validate() error {
	if t.Name == emptyString {
		return errParamIsRequired("name", typeString)
	}
	if t.Version == emptyString {
		return errParamIsRequired("version", typeString)
	}
	return nil
}

type DatabaseInstanceInputDTO struct {
	Name                 string `json:"name"`
	Host                 string `json:"host"`
	Port                 string `json:"port"`
	HostConnection       string `json:"hostConnection"`
	PortConnection       string `json:"portConnection"`
	AdminUser            string `json:"adminUser"`
	AdminPassword        string `json:"adminPassword"`
	EcosystemID          string `json:"ecosystemId"`
	DatabaseTechnologyID string `json:"databaseTechnologyId"`
	Note                 string `json:"note"`
}

func (d *DatabaseInstanceInputDTO) Validate() error {
	if d.Name == emptyString {
		return errParamIsRequired("name", typeString)
	}
	if d.Host == emptyString {
		return errParamIsRequired("host", typeString)
	}
	if d.Port == emptyString {
		return errParamIsRequired("port", typeString)
	}
	if d.HostConnection == emptyString {
		return errParamIsRequired("hostConnection", typeString)
	}
	if d.PortConnection == emptyString {
		return errParamIsRequired("portConnection", typeString)
	}
	if d.AdminUser == emptyString {
		return errParamIsRequired("adminUser", typeString)
	}
	if d.AdminPassword == emptyString {
		return errParamIsRequired("adminPassword", typeString)
	}
	if d.EcosystemID == emptyString {
		return errParamIsRequired("ecosystemId", typeUUID)
	}
	if !validUUID(d.EcosystemID) {
		return errParamIsInvalid("ecosystemId", typeUUID)
	}
	if d.DatabaseTechnologyID == emptyString {
		return errParamIsRequired("databaseTechnologyId", typeUUID)
	}
	if !validUUID(d.DatabaseTechnologyID) {
		return errParamIsInvalid("databaseTechnologyId", typeUUID)
	}
	return nil
}

func validUUID(id string) bool {
	_, err := uuid.Parse(id)
	return err == nil
}

type ConnectionInputDTO struct {
	ID         string `json:"id"`
	Host       string `json:"host"`
	Port       string `json:"port"`
	User       string `json:"user"`
	Password   string `json:"password"`
	Database   string `json:"database"`
	Instance   string `json:"instance"`
	Ecosystem  string `json:"ecosystem"`
	Technology string `json:"technology"`
}

type PropagateRolesInputDTO struct {
	DatabaseInstancesIDs []string `json:"databaseInstancesIds"`
}

type TestConnectionInputDTO struct {
	DatabaseInstancesIDs []string `json:"databaseInstancesIds"`
}

type SyncDatabasesInputDTO struct {
	DatabaseInstancesIDs []string `json:"databaseInstancesIds"`
}

type SetupRolesInputDTO struct {
	DatabaseInstanceID string   `json:"databaseInstanceId"`
	DatabasesIDs       []string `json:"databasesIds"`
}

type DatabaseUserInputDTO struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	Team           string `json:"team"`
	Position       string `json:"position"`
	DatabaseRoleID string `json:"databaseRoleId"`
}

func (d *DatabaseUserInputDTO) Validate() error {
	if d.Name == emptyString {
		return errParamIsRequired("name", typeString)
	}
	if d.Email == emptyString {
		return errParamIsRequired("email", typeString)
	}
	if !utils.ValidEmail(d.Email) {
		return errParamIsInvalid("email", typeString)
	}
	if d.DatabaseRoleID == emptyString {
		return errParamIsRequired("databaseRoleId", typeUUID)
	}
	if !validUUID(d.DatabaseRoleID) {
		return errParamIsInvalid("databaseRoleId", typeUUID)
	}
	return nil
}

type UpdateDatabaseUserInputDTO struct {
	Name           string `json:"name"`
	DatabaseRoleID string `json:"databaseRoleId"`
	Team           string `json:"team"`
	Position       string `json:"position"`
}

func (d *UpdateDatabaseUserInputDTO) Validate() error {
	if d.Name == emptyString {
		return errParamIsRequired("name", typeString)
	}
	if d.DatabaseRoleID == emptyString {
		return errParamIsRequired("databaseRoleId", typeUUID)
	}
	if !validUUID(d.DatabaseRoleID) {
		return errParamIsInvalid("databaseRoleId", typeUUID)
	}
	return nil
}

type InstanceDataDTO struct {
	DatabaseInstanceID string   `json:"databaseInstanceId"`
	DatabasesIDs       []string `json:"databasesIds"`
}

type GrantAccessInputDTO struct {
	DatabaseUsersIDs []string          `json:"databaseUsersIds"`
	InstancesData    []InstanceDataDTO `json:"instancesData"`
}

func (g *GrantAccessInputDTO) Validate() error {
	if len(g.DatabaseUsersIDs) == 0 {
		return ErrArrayDatabaseUsersIdsEmpty
	}
	if len(g.InstancesData) == 0 {
		return ErrArrayInstancesDataEmpty
	}
	for _, id := range g.DatabaseUsersIDs {
		if !validUUID(id) {
			return errParamIsInvalid("databaseUsersIds", typeUUID)
		}
	}
	for _, instanceData := range g.InstancesData {
		if !validUUID(instanceData.DatabaseInstanceID) {
			return errParamIsInvalid("instancesData", typeUUID)
		}
		for _, id := range instanceData.DatabasesIDs {
			if !validUUID(id) {
				return errParamIsInvalid("instancesData", typeUUID)
			}
		}
	}
	return nil
}

type RevokeAccessInputDTO struct {
	DatabaseUserID       string   `json:"databaseUserId"`
	DatabaseInstancesIDs []string `json:"databaseInstancesIds"`
}

func (r *RevokeAccessInputDTO) Validate() error {
	if r.DatabaseUserID == emptyString {
		return errParamIsRequired("databaseUserId", typeUUID)
	}
	if !validUUID(r.DatabaseUserID) {
		return errParamIsInvalid("databaseUserId", typeUUID)
	}
	for _, id := range r.DatabaseInstancesIDs {
		if !validUUID(id) {
			return errParamIsInvalid("databaseInstancesIds", typeUUID)
		}
	}
	return nil
}

type ChangeStatusInputDTO struct {
	ID      string `json:"id"`
	Enabled *bool  `json:"enabled"`
}

func (cs *ChangeStatusInputDTO) Validate() error {
	if cs.ID == emptyString {
		return errParamIsRequired("id", typeUUID)
	}
	if !validUUID(cs.ID) {
		return errParamIsInvalid("id", typeUUID)
	}
	if cs.Enabled == nil {
		return errParamIsRequired("enabled", typeBoolean)
	}
	return nil
}

func errParamIsRequired(name, typ string) error {
	return fmt.Errorf("param: %s (type: %s) is required", name, typ)
}

func errParamIsInvalid(name, typ string) error {
	return fmt.Errorf("param: %s (type: %s) is invalid", name, typ)
}

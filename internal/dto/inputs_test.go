package dto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateEcosystemInputDTO(t *testing.T) {
	i := &EcosystemInputDTO{}
	assertValidate(t, i, errParamIsRequired("code", typeString))

	i = &EcosystemInputDTO{Code: "poc"}
	assertValidate(t, i, errParamIsRequired("displayName", typeString))

	i = &EcosystemInputDTO{Code: "poc", DisplayName: "POC"}
	assert.NoError(t, i.Validate())
}

func TestValidateTechnologyInputDTO(t *testing.T) {
	i := &TechnologyInputDTO{}
	assertValidate(t, i, errParamIsRequired("name", typeString))

	i = &TechnologyInputDTO{Name: "PostgreSQL"}
	assertValidate(t, i, errParamIsRequired("version", typeString))

	i = &TechnologyInputDTO{Name: "PostgreSQL", Version: "13.4"}
	assert.NoError(t, i.Validate())
}

func TestValidateDatabaseInstanceInputDTO(t *testing.T) {
	i := &DatabaseInstanceInputDTO{}
	assertValidate(t, i, errParamIsRequired("name", typeString))

	i = &DatabaseInstanceInputDTO{Name: "PostgreSQL Cloud A"}
	assertValidate(t, i, errParamIsRequired("host", typeString))

	i = &DatabaseInstanceInputDTO{Name: "PostgreSQL", Host: "127.0.0.1"}
	assertValidate(t, i, errParamIsRequired("port", typeString))

	i = &DatabaseInstanceInputDTO{Name: "PostgreSQL", Host: "127.0.0.1", Port: "5432"}
	assertValidate(t, i, errParamIsRequired("hostConnection", typeString))

	i = &DatabaseInstanceInputDTO{Name: "PostgreSQL", Host: "127.0.0.1", Port: "5432", HostConnection: "host.conn.ip"}
	assertValidate(t, i, errParamIsRequired("portConnection", typeString))

	i = &DatabaseInstanceInputDTO{Name: "PostgreSQL", Host: "127.0.0.1", Port: "5432", HostConnection: "host.conn.ip", PortConnection: "5433"}
	assertValidate(t, i, errParamIsRequired("adminUser", typeString))

	i = &DatabaseInstanceInputDTO{Name: "PostgreSQL", Host: "127.0.0.1", Port: "5432", HostConnection: "host.conn.ip", PortConnection: "5433", AdminUser: "admin"}
	assertValidate(t, i, errParamIsRequired("adminPassword", typeString))

	i = &DatabaseInstanceInputDTO{Name: "PostgreSQL", Host: "127.0.0.1", Port: "5432", HostConnection: "host.conn.ip", PortConnection: "5433", AdminUser: "admin", AdminPassword: "pwd"}
	assertValidate(t, i, errParamIsRequired("ecosystemId", typeUUID))

	i = &DatabaseInstanceInputDTO{Name: "PostgreSQL", Host: "127.0.0.1", Port: "5432", HostConnection: "host.conn.ip", PortConnection: "5433", AdminUser: "admin", AdminPassword: "pwd", EcosystemID: "123"}
	assertValidate(t, i, errParamIsInvalid("ecosystemId", typeUUID))

	i = &DatabaseInstanceInputDTO{Name: "PostgreSQL", Host: "127.0.0.1", Port: "5432", HostConnection: "host.conn.ip", PortConnection: "5433", AdminUser: "admin", AdminPassword: "pwd", EcosystemID: "dd42cf0c-8a91-42d7-a906-cb9313494e7d"}
	assertValidate(t, i, errParamIsRequired("databaseTechnologyId", typeUUID))

	i = &DatabaseInstanceInputDTO{Name: "PostgreSQL", Host: "127.0.0.1", Port: "5432", HostConnection: "host.conn.ip", PortConnection: "5433", AdminUser: "admin", AdminPassword: "pwd", EcosystemID: "dd42cf0c-8a91-42d7-a906-cb9313494e7d", DatabaseTechnologyID: "123"}
	assertValidate(t, i, errParamIsInvalid("databaseTechnologyId", typeUUID))

	i = &DatabaseInstanceInputDTO{Name: "PostgreSQL", Host: "127.0.0.1", Port: "5432", HostConnection: "host.conn.ip", PortConnection: "5433", AdminUser: "admin", AdminPassword: "pwd", EcosystemID: "dd42cf0c-8a91-42d7-a906-cb9313494e7d", DatabaseTechnologyID: "1eb93da6-e739-4396-902f-19f79aa74e39"}
	assert.NoError(t, i.Validate())
}

func TestValidateDatabaseUserInputDTO(t *testing.T) {
	i := &DatabaseUserInputDTO{}
	assertValidate(t, i, errParamIsRequired("name", typeString))

	i = &DatabaseUserInputDTO{Name: "Foo Bar"}
	assertValidate(t, i, errParamIsRequired("email", typeString))

	i = &DatabaseUserInputDTO{Name: "Foo Bar", Email: "test.email"}
	assertValidate(t, i, errParamIsInvalid("email", typeString))

	i = &DatabaseUserInputDTO{Name: "Foo Bar", Email: "foobar@email.com"}
	assertValidate(t, i, errParamIsRequired("databaseRoleId", typeUUID))

	i = &DatabaseUserInputDTO{Name: "admin", Email: "foobar@email.com", DatabaseRoleID: "abc"}
	assertValidate(t, i, errParamIsInvalid("databaseRoleId", typeUUID))

	i = &DatabaseUserInputDTO{Name: "admin", Email: "foobar@email.com", DatabaseRoleID: "1eb93da6-e739-4396-902f-19f79aa74e39"}
	assert.NoError(t, i.Validate())
}

func TestValidateUpdateDatabaseUserInputDTO(t *testing.T) {
	i := &UpdateDatabaseUserInputDTO{}
	assertValidate(t, i, errParamIsRequired("name", typeString))

	i = &UpdateDatabaseUserInputDTO{Name: "Foo Bar"}
	assertValidate(t, i, errParamIsRequired("databaseRoleId", typeUUID))

	i = &UpdateDatabaseUserInputDTO{Name: "Foo Bar", DatabaseRoleID: "abc"}
	assertValidate(t, i, errParamIsInvalid("databaseRoleId", typeUUID))

	i = &UpdateDatabaseUserInputDTO{Name: "admin", DatabaseRoleID: "1eb93da6-e739-4396-902f-19f79aa74e39"}
	assert.NoError(t, i.Validate())
}

func TestValidateGrantAccessInputDTO(t *testing.T) {
	i := &GrantAccessInputDTO{}
	assertValidate(t, i, ErrArrayDatabaseUsersIdsEmpty)

	i = &GrantAccessInputDTO{DatabaseUsersIDs: []string{"1"}}
	assertValidate(t, i, ErrArrayInstancesDataEmpty)

	i = &GrantAccessInputDTO{DatabaseUsersIDs: []string{"1"}, InstancesData: []InstanceDataDTO{{DatabaseInstanceID: "1", DatabasesIDs: []string{"1"}}}}
	assertValidate(t, i, errParamIsInvalid("databaseUsersIds", typeUUID))

	i = &GrantAccessInputDTO{DatabaseUsersIDs: []string{"1eb93da6-e739-4396-902f-19f79aa74e39"}, InstancesData: []InstanceDataDTO{{DatabaseInstanceID: "1", DatabasesIDs: []string{"1"}}}}
	assertValidate(t, i, errParamIsInvalid("instancesData", typeUUID))

	i = &GrantAccessInputDTO{DatabaseUsersIDs: []string{"1eb93da6-e739-4396-902f-19f79aa74e39"},
		InstancesData: []InstanceDataDTO{{DatabaseInstanceID: "1eb93da6-e739-4396-902f-19f79aa74e39", DatabasesIDs: []string{"1"}}}}
	assertValidate(t, i, errParamIsInvalid("instancesData", typeUUID))

	i = &GrantAccessInputDTO{DatabaseUsersIDs: []string{"1eb93da6-e739-4396-902f-19f79aa74e39", "123"},
		InstancesData: []InstanceDataDTO{{DatabaseInstanceID: "1", DatabasesIDs: []string{"1"}}}}
	assertValidate(t, i, errParamIsInvalid("databaseUsersIds", typeUUID))

	i = &GrantAccessInputDTO{DatabaseUsersIDs: []string{"1eb93da6-e739-4396-902f-19f79aa74e39"},
		InstancesData: []InstanceDataDTO{
			{DatabaseInstanceID: "1eb93da6-e739-4396-902f-19f79aa74e39", DatabasesIDs: []string{"1eb93da6-e739-4396-902f-19f79aa74e39"}},
			{DatabaseInstanceID: "1", DatabasesIDs: []string{"1eb93da6-e739-4396-902f-19f79aa74e39"}},
		}}
	assertValidate(t, i, errParamIsInvalid("instancesData", typeUUID))

	i = &GrantAccessInputDTO{DatabaseUsersIDs: []string{"1eb93da6-e739-4396-902f-19f79aa74e39"},
		InstancesData: []InstanceDataDTO{
			{DatabaseInstanceID: "1eb93da6-e739-4396-902f-19f79aa74e39", DatabasesIDs: []string{"1eb93da6-e739-4396-902f-19f79aa74e39"}},
			{DatabaseInstanceID: "1eb93da6-e739-4396-902f-19f79aa74e39", DatabasesIDs: []string{"1"}},
		}}
	assertValidate(t, i, errParamIsInvalid("instancesData", typeUUID))

	i = &GrantAccessInputDTO{DatabaseUsersIDs: []string{"1eb93da6-e739-4396-902f-19f79aa74e39"},
		InstancesData: []InstanceDataDTO{{DatabaseInstanceID: "1eb93da6-e739-4396-902f-19f79aa74e39", DatabasesIDs: []string{"1eb93da6-e739-4396-902f-19f79aa74e39"}}}}
	assert.NoError(t, i.Validate())
}

func TestValidateRevokeAccessInputDTO(t *testing.T) {
	i := &RevokeAccessInputDTO{}
	assertValidate(t, i, errParamIsRequired("databaseUserId", typeUUID))

	i = &RevokeAccessInputDTO{DatabaseUserID: "1"}
	assertValidate(t, i, errParamIsInvalid("databaseUserId", typeUUID))

	i = &RevokeAccessInputDTO{DatabaseUserID: "1eb93da6-e739-4396-902f-19f79aa74e39", DatabaseInstancesIDs: []string{"1"}}
	assertValidate(t, i, errParamIsInvalid("databaseInstancesIds", typeUUID))

	i = &RevokeAccessInputDTO{
		DatabaseUserID:       "1eb93da6-e739-4396-902f-19f79aa74e39",
		DatabaseInstancesIDs: []string{"96cfa8f2-2c91-4630-b556-f7a2eab84e29", "63862219-f1c3-41c4-9938-31346773b697"},
	}
	assert.NoError(t, i.Validate())
}

func TestValidateChangeStatusDBUserInputDTO(t *testing.T) {
	i := &ChangeStatusInputDTO{}
	assertValidate(t, i, errParamIsRequired("id", typeUUID))

	i = &ChangeStatusInputDTO{ID: "1"}
	assertValidate(t, i, errParamIsInvalid("id", typeUUID))

	i = &ChangeStatusInputDTO{ID: "1eb93da6-e739-4396-902f-19f79aa74e39"}
	assertValidate(t, i, errParamIsRequired("enabled", typeBoolean))

	trueParam := true
	i = &ChangeStatusInputDTO{
		ID:      "1eb93da6-e739-4396-902f-19f79aa74e39",
		Enabled: &trueParam,
	}
	assert.NoError(t, i.Validate())
}

func assertValidate(t *testing.T, dto InputValidator, expectedError error) {
	err := dto.Validate()
	assert.Error(t, err)
	assert.EqualError(t, err, expectedError.Error())
}

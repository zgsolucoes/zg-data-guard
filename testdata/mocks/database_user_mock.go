package mocks

import (
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/zgsolucoes/zg-data-guard/internal/database/connector"
	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	"github.com/zgsolucoes/zg-data-guard/internal/entity"
)

const (
	DbUserID  = "1eb93da6-e739-4396-902f-19f79aa74e39"
	dbUserID2 = "cd7f93a4-a2ff-41db-9ad2-6dd67dd285c7"
	roleID1   = "cd7f93a4-a2ff-41db-9ad2-6dd67dd285c7"
	roleID2   = "1eb93da6-e739-4396-902f-19f79aa74e39"
)

var (
	ValidDBUserInput = dto.DatabaseUserInputDTO{
		Name:           "Foo Bar",
		Email:          "foobar@email.com",
		Team:           "Team A",
		Position:       "Developer",
		DatabaseRoleID: roleID1,
	}

	ValidUpdateDBUserInput = dto.UpdateDatabaseUserInputDTO{
		Name:           "Foo Bar Updated",
		DatabaseRoleID: roleID2,
		Team:           "Team C",
		Position:       "Customer Success",
	}
)

type DatabaseUserStorageMock struct {
	mock.Mock
}

func (m *DatabaseUserStorageMock) FindAll(ids []string) ([]*entity.DatabaseUser, error) {
	args := m.Called(ids)
	return args.Get(0).([]*entity.DatabaseUser), args.Error(1)
}
func (m *DatabaseUserStorageMock) FindAllDTOs(ids []string) ([]*dto.DatabaseUserOutputDTO, error) {
	args := m.Called(ids)
	return args.Get(0).([]*dto.DatabaseUserOutputDTO), args.Error(1)
}

func (m *DatabaseUserStorageMock) FindAllDTOsEnabled() ([]*dto.DatabaseUserOutputDTO, error) {
	args := m.Called()
	return args.Get(0).([]*dto.DatabaseUserOutputDTO), args.Error(1)
}

func (m *DatabaseUserStorageMock) FindByID(id string) (*entity.DatabaseUser, error) {
	args := m.Called(id)
	return args.Get(0).(*entity.DatabaseUser), args.Error(1)
}

func (m *DatabaseUserStorageMock) FindDTOByID(id string) (*dto.DatabaseUserOutputDTO, error) {
	args := m.Called(id)
	return args.Get(0).(*dto.DatabaseUserOutputDTO), args.Error(1)
}

func (m *DatabaseUserStorageMock) Save(d *entity.DatabaseUser) error {
	args := m.Called(d)
	return args.Error(0)
}

func (m *DatabaseUserStorageMock) Update(d *entity.DatabaseUser) error {
	args := m.Called(d)
	return args.Error(0)
}

func (m *DatabaseUserStorageMock) Exists(email string) (bool, error) {
	args := m.Called(email)
	return args.Bool(0), args.Error(1)
}

func BuildDbUserDTOList() []*dto.DatabaseUserOutputDTO {
	userFoo := BuildDbUserFooDTO()
	userJohn := BuildDbUserJohnDTO()
	dbUsers := []*dto.DatabaseUserOutputDTO{userFoo, userJohn}
	return dbUsers
}

func BuildDbUserFooDTO() *dto.DatabaseUserOutputDTO {
	return &dto.DatabaseUserOutputDTO{
		ID:                      DbUserID,
		Name:                    ValidDBUserInput.Name,
		Email:                   ValidDBUserInput.Email,
		Username:                "foobar",
		Password:                "123456",
		DatabaseRoleID:          ValidDBUserInput.DatabaseRoleID,
		DatabaseRoleName:        "developer",
		DatabaseRoleDisplayName: "Developer",
		Enabled:                 true,
		Team:                    ValidDBUserInput.Team,
		Position:                ValidDBUserInput.Position,
		CreatedByUserID:         UserID,
		CreatedByUser:           "Susan Smith",
		CreatedAt:               time.Now(),
	}
}

func BuildDbUserJohn() *entity.DatabaseUser {
	return &entity.DatabaseUser{
		ID:              uuid.MustParse(dbUserID2),
		Name:            "John Doe",
		Email:           "johndoe@email.com",
		Username:        "johndoe",
		Password:        "postgres",
		DatabaseRoleID:  roleID2,
		Enabled:         true,
		Team:            "Team B",
		Position:        "DevOps",
		CreatedByUserID: UserID,
		CreatedAt:       time.Now(),
	}
}

func BuildDbUserJohnDTO() *dto.DatabaseUserOutputDTO {
	return &dto.DatabaseUserOutputDTO{
		ID:                      dbUserID2,
		Name:                    "John Doe",
		Email:                   "johndoe@email.com",
		Username:                "johndoe",
		Password:                "654321",
		DatabaseRoleID:          roleID2,
		DatabaseRoleName:        "devops",
		DatabaseRoleDisplayName: "DevOps",
		Enabled:                 true,
		Team:                    "Team B",
		Position:                "DevOps",
		CreatedByUserID:         UserID,
		CreatedByUser:           "Susan Smith",
		CreatedAt:               time.Now(),
	}
}

func BuildDisabledDbUserJohnDTO() *dto.DatabaseUserOutputDTO {
	user := BuildDbUserJohnDTO()
	user.Enabled = false
	return user
}

func BuildDbUserDummyDTO() *dto.DatabaseUserOutputDTO {
	return &dto.DatabaseUserOutputDTO{
		ID:               "dummy-id",
		Name:             "Dummy User",
		Email:            "dummy@email.com",
		Username:         connector.DummyTestUser,
		Password:         "654321",
		DatabaseRoleName: "user_ro",
		Enabled:          true,
	}
}

func BuildDbUserDummyErrorCreateDTO() *dto.DatabaseUserOutputDTO {
	user := BuildDbUserDummyDTO()
	user.Username = connector.DummyTestUserErrorOnCreate
	return user
}

func BuildDbUserDummyErrorGrantDTO() *dto.DatabaseUserOutputDTO {
	user := BuildDbUserDummyDTO()
	user.Username = connector.DummyTestUserErrorOnGrant
	return user
}

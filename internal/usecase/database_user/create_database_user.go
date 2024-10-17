package dbuser

import (
	"database/sql"
	"errors"
	"log"

	"github.com/zgsolucoes/zg-data-guard/internal/database/storage"
	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	"github.com/zgsolucoes/zg-data-guard/internal/entity"
)

const errorCreatingDatabaseUser = "Error creating database user"

var (
	ErrEmailAlreadyExists   = errors.New("this e-mail already exists for a database user")
	ErrDatabaseRoleNotFound = errors.New("database role not found")
)

type CreateDatabaseUserUseCase struct {
	DatabaseUserStorage storage.DatabaseUserStorage
	DatabaseRoleStorage storage.DatabaseRoleStorage
}

func NewCreateDatabaseUserUseCase(
	databaseUserStorage storage.DatabaseUserStorage,
	databaseRoleStorage storage.DatabaseRoleStorage,
) *CreateDatabaseUserUseCase {
	return &CreateDatabaseUserUseCase{
		DatabaseUserStorage: databaseUserStorage,
		DatabaseRoleStorage: databaseRoleStorage,
	}
}

func (c *CreateDatabaseUserUseCase) Execute(input dto.DatabaseUserInputDTO, createdByUserID string) (*dto.DatabaseUserOutputDTO, error) {
	databaseUser, err := entity.NewDatabaseUser(input.Name, input.Email, input.Team, input.Position, input.DatabaseRoleID, createdByUserID)
	if err != nil {
		logError(err, errorCreatingDatabaseUser)
		return nil, err
	}

	exists, err := c.DatabaseUserStorage.Exists(input.Email)
	if err != nil {
		logError(err, errorCreatingDatabaseUser)
		return nil, err
	}
	if exists {
		logError(ErrEmailAlreadyExists, errorCreatingDatabaseUser)
		return nil, ErrEmailAlreadyExists
	}

	err = validateDatabaseRoleExistence(input.DatabaseRoleID, c.DatabaseRoleStorage, errorCreatingDatabaseUser)
	if err != nil {
		return nil, err
	}

	err = c.DatabaseUserStorage.Save(databaseUser)
	if err != nil {
		logError(err, errorCreatingDatabaseUser)
		return nil, err
	}

	log.Printf("Database user %v created successfully by user %s!", databaseUser.ID, createdByUserID)
	return &dto.DatabaseUserOutputDTO{
		ID:              databaseUser.ID.String(),
		Name:            databaseUser.Name,
		Email:           databaseUser.Email,
		Username:        databaseUser.Username,
		Password:        databaseUser.Password,
		DatabaseRoleID:  databaseUser.DatabaseRoleID,
		Team:            databaseUser.Team,
		Position:        databaseUser.Position,
		Enabled:         databaseUser.Enabled,
		CreatedAt:       databaseUser.CreatedAt,
		CreatedByUserID: databaseUser.CreatedByUserID,
	}, nil
}

func validateDatabaseRoleExistence(
	databaseRoleID string,
	roleStorage storage.DatabaseRoleStorage,
	operation string,
) error {
	_, err := roleStorage.FindByID(databaseRoleID)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		logError(ErrDatabaseRoleNotFound, operation)
		return ErrDatabaseRoleNotFound
	}
	if err != nil {
		logError(err, operation)
		return err
	}
	log.Printf("Database Role with ID %s loaded successfully!", databaseRoleID)
	return nil
}

func logError(err error, operationError string) {
	log.Printf("%s. Cause: %v", operationError, err.Error())
}

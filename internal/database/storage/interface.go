package storage

import (
	"database/sql"

	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	"github.com/zgsolucoes/zg-data-guard/internal/entity"
)

type DBInterface interface {
	Query(query string, args ...any) (*sql.Rows, error)
}

type ApplicationUserStorage interface {
	FindByEmail(email string) (*entity.ApplicationUser, error)
	FindByID(id string) (*entity.ApplicationUser, error)
}

type EcosystemStorage interface {
	Save(ecosystem *entity.Ecosystem) error
	Update(ecosystem *entity.Ecosystem) error
	FindByID(id string) (*entity.Ecosystem, error)
	FindAll(page, limit int) ([]*dto.EcosystemOutputDTO, error)
	Delete(id string) error
	CheckCodeExists(code string) (bool, error)
}

type DatabaseTechnologyStorage interface {
	Save(databaseTechnology *entity.DatabaseTechnology) error
	Update(databaseTechnology *entity.DatabaseTechnology) error
	Exists(name, version string) (bool, error)
	FindByID(id string) (*entity.DatabaseTechnology, error)
	Delete(id string) error
	FindAll(page, limit int) ([]*dto.TechnologyOutputDTO, error)
}

type DatabaseInstanceStorage interface {
	Save(databaseInstance *entity.DatabaseInstance) error
	UpdateWithHostInfo(databaseInstance *entity.DatabaseInstance) error
	Update(databaseInstance *entity.DatabaseInstance) error
	Exists(host, port string) (bool, error)
	FindByID(id string) (*entity.DatabaseInstance, error)
	FindDTOByID(id string) (*dto.DatabaseInstanceOutputDTO, error)
	FindAllDTOs(ecosystemID, technologyID string, ids []string) ([]*dto.DatabaseInstanceOutputDTO, error)
	FindAllDTOsEnabled(ecosystemID, technologyID string) ([]*dto.DatabaseInstanceOutputDTO, error)
}

type DatabaseRoleStorage interface {
	FindAll() ([]*entity.DatabaseRole, error)
	FindByID(id string) (*entity.DatabaseRole, error)
}

type DatabaseStorage interface {
	Save(database *entity.Database) error
	Update(database *entity.Database) error
	FindDTOByID(id string) (*dto.DatabaseOutputDTO, error)
	FindAll(databaseInstanceID string, ids []string) ([]*entity.Database, error)
	FindAllEnabled(databaseInstanceID string) ([]*entity.Database, error)
	FindAllDTOs(ecosystemID, databaseInstanceID string) ([]*dto.DatabaseOutputDTO, error)
	DeactivateAllByInstance(databaseInstanceID string) error
}

type DatabaseUserStorage interface {
	Save(d *entity.DatabaseUser) error
	Update(d *entity.DatabaseUser) error
	Exists(email string) (bool, error)
	FindByID(id string) (*entity.DatabaseUser, error)
	FindDTOByID(id string) (*dto.DatabaseUserOutputDTO, error)
	FindAll(ids []string) ([]*entity.DatabaseUser, error)
	FindAllDTOs(ids []string) ([]*dto.DatabaseUserOutputDTO, error)
	FindAllDTOsEnabled() ([]*dto.DatabaseUserOutputDTO, error)
}

type AccessPermissionStorage interface {
	Save(d *entity.AccessPermission) error
	Exists(databaseID, databaseUserID string) (bool, error)
	DeleteAllByInstance(instanceID string) error
	DeleteAllByUserAndInstance(databaseUserID, instanceID string) error
	FindAllDTOs(databaseID, databaseUserID, databaseInstanceID string) ([]*dto.AccessPermissionOutputDTO, error)
	SaveLog(log *entity.AccessPermissionLog) error
	FindAllAccessibleInstancesIDsByUser(userID string) ([]string, error)
	FindAllLogsDTOs(page, limit int) ([]*dto.AccessPermissionLogOutputDTO, error)
	CheckIfUserHasAccessPermission(databaseUserID string) (bool, error)
	LogCount() (int, error)
}

type ForbiddenObjectsStorage interface {
	FindAllDatabases() ([]*entity.ForbiddenDatabase, error)
}

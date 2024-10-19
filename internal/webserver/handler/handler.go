package handler

import (
	"github.com/zgsolucoes/zg-data-guard/config"
	database "github.com/zgsolucoes/zg-data-guard/internal/database/storage"
	permissionUsecase "github.com/zgsolucoes/zg-data-guard/internal/usecase/access_permission"
	databaseUsecase "github.com/zgsolucoes/zg-data-guard/internal/usecase/database"
	dbInstanceUsecase "github.com/zgsolucoes/zg-data-guard/internal/usecase/database_instance"
	roleUsecase "github.com/zgsolucoes/zg-data-guard/internal/usecase/database_role"
	databaseUserUsecase "github.com/zgsolucoes/zg-data-guard/internal/usecase/database_user"
	ecosystemUsecase "github.com/zgsolucoes/zg-data-guard/internal/usecase/ecosystem"
	technologyUsecase "github.com/zgsolucoes/zg-data-guard/internal/usecase/technology"
	userUsecase "github.com/zgsolucoes/zg-data-guard/internal/usecase/user"
)

var (
	// Storage
	appUserStorage          database.ApplicationUserStorage
	ecosystemStorage        database.EcosystemStorage
	technologyStorage       database.DatabaseTechnologyStorage
	instanceStorage         database.DatabaseInstanceStorage
	databaseStorage         database.DatabaseStorage
	roleStorage             database.DatabaseRoleStorage
	dbUserStorage           database.DatabaseUserStorage
	accessStorage           database.AccessPermissionStorage
	forbiddenObjectsStorage database.ForbiddenObjectsStorage
)

func InitializeAPIDependencies() {
	initializeEntityStorages()
	initializeUseCases()
}

func initializeEntityStorages() {
	db := config.GetDBConn()
	appUserStorage = database.NewPostgresApplicationUserStorage(db)
	ecosystemStorage = database.NewPostgresEcosystemStorage(db)
	technologyStorage = database.NewPostgresDatabaseTechnologyStorage(db)
	instanceStorage = database.NewPostgresInstanceStorage(db)
	databaseStorage = database.NewPostgresDatabaseStorage(db)
	roleStorage = database.NewPostgresDatabaseRoleStorage(db)
	dbUserStorage = database.NewPostgresDatabaseUserStorage(db)
	accessStorage = database.NewPostgresAccessPermissionStorage(db)
	forbiddenObjectsStorage = database.NewPostgresForbiddenObjectsStorage(db)
}

func initializeUseCases() {
	initializeUserUseCases(appUserStorage)
	initializeEcosystemUseCases(ecosystemStorage, appUserStorage)
	initializeTechnologyUseCases(technologyStorage, appUserStorage)
	initializeDatabaseInstanceUseCases(instanceStorage, ecosystemStorage, technologyStorage, databaseStorage, roleStorage, accessStorage)
	initializeDatabaseUseCases(instanceStorage, databaseStorage)
	initializeDatabaseRoleUseCases(roleStorage)
	initializeAccessPermissionUseCases(accessStorage, dbUserStorage, instanceStorage, databaseStorage, forbiddenObjectsStorage)
	initializeDatabaseUserUseCases(dbUserStorage, roleStorage, accessStorage)
}

func initializeUserUseCases(appUserStorage database.ApplicationUserStorage) {
	getUserUC = userUsecase.NewGetUserUseCase(appUserStorage)
}

func initializeEcosystemUseCases(ecosystemStorage database.EcosystemStorage, appUserStorage database.ApplicationUserStorage) {
	createEcosystemUC = ecosystemUsecase.NewCreateEcosystemUseCase(ecosystemStorage)
	getEcosystemUC = ecosystemUsecase.NewGetEcosystemUseCase(ecosystemStorage, appUserStorage)
	deleteEcosystemUC = ecosystemUsecase.NewDeleteEcosystemUseCase(ecosystemStorage)
	updateEcosystemUC = ecosystemUsecase.NewUpdateEcosystemUseCase(ecosystemStorage)
	listEcosystemsUC = ecosystemUsecase.NewListEcosystemsUseCase(ecosystemStorage)
}

func initializeTechnologyUseCases(technologyStorage database.DatabaseTechnologyStorage, appUserStorage database.ApplicationUserStorage) {
	createTechnologyUC = technologyUsecase.NewCreateTechnologyUseCase(technologyStorage)
	getTechnologyUC = technologyUsecase.NewGetTechnologyUseCase(technologyStorage, appUserStorage)
	deleteTechnologyUC = technologyUsecase.NewDeleteTechnologyUseCase(technologyStorage)
	updateTechnologyUC = technologyUsecase.NewUpdateTechnologyUseCase(technologyStorage)
	listTechnologiesUC = technologyUsecase.NewListTechnologiesUseCase(technologyStorage)
}

func initializeDatabaseInstanceUseCases(
	dbInstanceStorage database.DatabaseInstanceStorage,
	ecosystemStorage database.EcosystemStorage,
	technologyStorage database.DatabaseTechnologyStorage,
	databaseStorage database.DatabaseStorage,
	roleStorage database.DatabaseRoleStorage,
	accessStorage database.AccessPermissionStorage,
) {
	createDBInstanceUC = dbInstanceUsecase.NewCreateDatabaseInstanceUseCase(dbInstanceStorage, ecosystemStorage, technologyStorage)
	getDBInstanceUC = dbInstanceUsecase.NewGetDatabaseInstanceUseCase(dbInstanceStorage)
	listDBInstancesUC = dbInstanceUsecase.NewListDatabaseInstancesUseCase(dbInstanceStorage)
	updateDBInstanceUC = dbInstanceUsecase.NewUpdateDatabaseInstanceUseCase(dbInstanceStorage, ecosystemStorage, technologyStorage)
	testConnectionUC = dbInstanceUsecase.NewTestConnectionUseCase(dbInstanceStorage)
	syncDatabasesUC = databaseUsecase.NewSyncDatabasesUseCase(dbInstanceStorage, databaseStorage)
	propagateRolesUC = dbInstanceUsecase.NewPropagateRolesUseCase(dbInstanceStorage, roleStorage)
	changeStatusInstanceUC = dbInstanceUsecase.NewChangeStatusDatabaseInstanceUseCase(dbInstanceStorage, databaseStorage, accessStorage)
}

func initializeDatabaseUseCases(dbInstanceStorage database.DatabaseInstanceStorage, databaseStorage database.DatabaseStorage) {
	getDatabaseUC = databaseUsecase.NewGetDatabaseUseCase(databaseStorage)
	listDatabasesUC = databaseUsecase.NewListDatabasesUseCase(databaseStorage)
	setupRolesUC = databaseUsecase.NewSetupRolesInDatabasesUseCase(dbInstanceStorage, databaseStorage)
}

func initializeDatabaseRoleUseCases(roleStorage database.DatabaseRoleStorage) {
	listDatabaseRolesUC = roleUsecase.NewListDatabaseRolesUseCase(roleStorage)
}

func initializeDatabaseUserUseCases(
	dbUserStorage database.DatabaseUserStorage,
	roleStorage database.DatabaseRoleStorage,
	accessPermissionStorage database.AccessPermissionStorage,
) {
	createDBUserUC = databaseUserUsecase.NewCreateDatabaseUserUseCase(dbUserStorage, roleStorage)
	getDBUserUC = databaseUserUsecase.NewGetDatabaseUserUseCase(dbUserStorage)
	updateDBUserUC = databaseUserUsecase.NewUpdateDatabaseUserUseCase(dbUserStorage, roleStorage, accessPermissionStorage)
	listDatabaseUsersUC = databaseUserUsecase.NewListDatabaseUsersUseCase(dbUserStorage)
	changeStatusDBUserUC = databaseUserUsecase.NewChangeStatusDatabaseUserUseCase(dbUserStorage, revokeAccessPermissionUC)
}

func initializeAccessPermissionUseCases(
	accessStorage database.AccessPermissionStorage,
	dbUserStorage database.DatabaseUserStorage,
	dbInstanceStorage database.DatabaseInstanceStorage,
	databaseStorage database.DatabaseStorage,
	forbiddenStorage database.ForbiddenObjectsStorage,
) {
	grantAccessPermissionUC = permissionUsecase.NewGrantAccessPermissionUseCase(accessStorage, dbUserStorage, dbInstanceStorage, databaseStorage, forbiddenStorage)
	listAccessPermissionsUC = permissionUsecase.NewListAccessPermissionsUseCase(accessStorage)
	listAccessPermissionLogsUC = permissionUsecase.NewListAccessPermissionLogsUseCase(accessStorage)
	revokeAccessPermissionUC = permissionUsecase.NewRevokeAccessPermissionUseCase(accessStorage, dbInstanceStorage, dbUserStorage)
}

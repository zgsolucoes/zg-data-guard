package accesspermission

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/zgsolucoes/zg-data-guard/config"
	"github.com/zgsolucoes/zg-data-guard/internal/database/connector"
	"github.com/zgsolucoes/zg-data-guard/internal/database/storage"
	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	"github.com/zgsolucoes/zg-data-guard/internal/entity"
)

const (
	AccessGrantedMsg           = "Access permissions created successfully."
	SomeErrorsDuringProcessMsg = "Some errors occurred during the process. Check the logs for more details."
)

type GrantAccessPermissionUseCase struct {
	AccessPermissionStorage storage.AccessPermissionStorage
	DatabaseUserStorage     storage.DatabaseUserStorage
	DatabaseInstanceStorage storage.DatabaseInstanceStorage
	DatabaseStorage         storage.DatabaseStorage
	ForbiddenObjectsStorage storage.ForbiddenObjectsStorage
}

func NewGrantAccessPermissionUseCase(
	accessPermissionStorage storage.AccessPermissionStorage,
	databaseUserStorage storage.DatabaseUserStorage,
	databaseInstanceStorage storage.DatabaseInstanceStorage,
	databaseStorage storage.DatabaseStorage,
	forbiddenObjectsStorage storage.ForbiddenObjectsStorage,
) *GrantAccessPermissionUseCase {
	return &GrantAccessPermissionUseCase{
		AccessPermissionStorage: accessPermissionStorage,
		DatabaseUserStorage:     databaseUserStorage,
		DatabaseInstanceStorage: databaseInstanceStorage,
		DatabaseStorage:         databaseStorage,
		ForbiddenObjectsStorage: forbiddenObjectsStorage,
	}
}

// Execute godoc
/** Responsible for create access permissions for the provided users in the selected instances and databases.
If the user is already created in the instance, it grants the connect permission to the user in the selected databases.
It returns an output DTO that has a flag indicating if the process has errors and a message with the result.
For each error that occurs inside the instance context during the process, it's logged, persisted and the process continues.
The process is divided into four main contexts: global, instance, user and database.
The process occurs concurrently for each instance, user and database. I.e., if there are 3 instances, 5 users and 4 databases, there will be 60 goroutines running concurrently. */
func (useCase *GrantAccessPermissionUseCase) Execute(input dto.GrantAccessInputDTO, operationUserID string) (*dto.GrantAccessOutputDTO, error) {
	start := time.Now()
	dbUsers, err := useCase.DatabaseUserStorage.FindAllDTOs(input.DatabaseUsersIDs)
	if err != nil {
		return nil, err
	}

	dbInstancesIds, dbIdsByInstance := useCase.prepareInstanceData(input.InstancesData)
	dbInstances, err := useCase.DatabaseInstanceStorage.FindAllDTOs("", "", dbInstancesIds)
	if err != nil {
		return nil, err
	}

	forbiddenDatabaseMap, err := useCase.fetchForbiddenDatabases()
	if err != nil {
		return nil, err
	}
	instancesQty := len(dbInstances)
	usersQty := len(dbUsers)
	globalCtx := newGrantAccessGlobalContext(dbUsers, dbIdsByInstance, operationUserID, forbiddenDatabaseMap, instancesQty, usersQty)
	log.Printf("Starting to process grant permissions to %d users in %d instances", usersQty, instancesQty)
	var wg sync.WaitGroup
	wg.Add(instancesQty)
	for idx, dbInstance := range dbInstances {
		go func(instanceDTO *dto.DatabaseInstanceOutputDTO, instanceIdx int) {
			defer wg.Done()
			instanceCtx := newGrantAccessInstanceContext(globalCtx, instanceDTO, instanceIdx)
			errValidating := useCase.validateInstance(instanceCtx)
			if errValidating != nil {
				globalCtx.GlobalErrChan <- errValidating
				return
			}
			if err := useCase.processInstance(instanceCtx); err != nil {
				globalCtx.GlobalErrChan <- err
			}
		}(dbInstance, idx)
	}

	go func() {
		wg.Wait()
		close(globalCtx.GlobalErrChan)
	}()
	log.Printf("All %d instances processed. Elapsed time: %s", instancesQty, time.Since(start))
	return buildGrantAccessOutput(globalCtx.GlobalErrChan), nil
}

func (useCase *GrantAccessPermissionUseCase) prepareInstanceData(instancesData []dto.InstanceDataDTO) ([]string, map[string][]string) {
	dbInstancesIds := make([]string, 0, len(instancesData))
	databasesIdsByInstance := make(map[string][]string)

	for _, instanceData := range instancesData {
		dbInstancesIds = append(dbInstancesIds, instanceData.DatabaseInstanceID)
		databasesIdsByInstance[instanceData.DatabaseInstanceID] = instanceData.DatabasesIDs
	}
	return dbInstancesIds, databasesIdsByInstance
}

func (useCase *GrantAccessPermissionUseCase) fetchForbiddenDatabases() (map[string]bool, error) {
	forbiddenDatabases, err := useCase.ForbiddenObjectsStorage.FindAllDatabases()
	if err != nil {
		return nil, fmt.Errorf("error when fetching forbidden databases. Cause: %v", err)
	}
	forbiddenDatabasesMap := make(map[string]bool, len(forbiddenDatabases))
	for _, db := range forbiddenDatabases {
		forbiddenDatabasesMap[db.Name] = true
	}
	return forbiddenDatabasesMap, nil
}

func (useCase *GrantAccessPermissionUseCase) validateInstance(instanceCtx *instanceContextOnGrant) error {
	if !instanceCtx.Instance.Enabled {
		return useCase.registerInstanceValidationError(instanceCtx, fmt.Sprintf(ErrInstanceDisabledMsg, instanceCtx.Instance.Name), ErrInstanceDisabled)
	}
	if !instanceCtx.Instance.RolesCreated {
		return useCase.registerInstanceValidationError(instanceCtx, fmt.Sprintf(ErrRolesNotCreatedMsg, instanceCtx.Instance.Name), ErrRolesNotCreated)
	}

	return nil
}

func (useCase *GrantAccessPermissionUseCase) processInstance(instanceCtx *instanceContextOnGrant) error {
	targetInstance, err := connector.NewDatabaseConnector(instanceCtx.Instance, "")
	if err != nil {
		return useCase.registerInstanceValidationError(instanceCtx, fmt.Sprintf(ErrCreatingConnectorMsg, instanceCtx.Instance.Name, err.Error()), err)
	}

	var wg sync.WaitGroup
	wg.Add(instanceCtx.GlobalCtx.UsersQty)
	logInstanceContextWithIndex(instanceCtx, fmt.Sprintf("starting processing %d users...", instanceCtx.GlobalCtx.UsersQty), false)
	defer logInstanceContextWithIndex(instanceCtx, fmt.Sprintf("finished processing %d users", instanceCtx.GlobalCtx.UsersQty), false)
	for idx, userDTO := range instanceCtx.GlobalCtx.DBUsers {
		go func(userDTO *dto.DatabaseUserOutputDTO, userIndex int) {
			defer wg.Done()

			userCtx := newGrantAccessUserContext(instanceCtx, targetInstance, userDTO, userIndex)
			errValidating := useCase.validateUser(userCtx)
			if errValidating != nil {
				instanceCtx.GlobalCtx.GlobalErrChan <- errValidating
				return
			}

			if err := useCase.processUser(userCtx); err != nil {
				instanceCtx.GlobalCtx.GlobalErrChan <- err
			}
		}(userDTO, idx)
	}
	wg.Wait()

	return nil
}

func (useCase *GrantAccessPermissionUseCase) validateUser(userCtx *userContextOnGrant) error {
	if !userCtx.DBUser.Enabled {
		return useCase.registerUserValidationError(userCtx, fmt.Sprintf(ErrUserDisabledMsg, userCtx.DBUser.Username), ErrUserDisabled)
	}

	return nil
}

func (useCase *GrantAccessPermissionUseCase) processUser(userCtx *userContextOnGrant) error {
	logUserContextWithIndex(userCtx, "validating existence of user", false)
	userExists, err := userCtx.TargetInstance.UserExists(userCtx.DBUser.Username)
	if err != nil {
		errConnection := fmt.Errorf("connection failed with instance. Cause: %v", err)
		return useCase.registerUserValidationError(userCtx, fmt.Sprintf(ErrConnectionFailedMsg, userCtx.InstanceCtx.Instance.Name, err.Error()), errConnection)
	}

	if userExists {
		logUserContextWithIndex(userCtx, "user already exists in instance", false)
	} else {
		if errCreating := useCase.createUser(userCtx); errCreating != nil {
			return errCreating
		}
	}

	return useCase.processDatabases(userCtx)
}

func (useCase *GrantAccessPermissionUseCase) createUser(userCtx *userContextOnGrant) error {
	if !entity.ValidateRoleName(userCtx.DBUser.DatabaseRoleName) {
		return useCase.registerUserValidationError(userCtx, fmt.Sprintf(ErrInvalidRoleMsg, userCtx.DBUser.DatabaseRoleName, userCtx.DBUser.Username), ErrInvalidRole)
	}

	decryptedPwd, err := config.GetCryptoHelper().Decrypt(userCtx.DBUser.Password)
	if err != nil {
		errDecrypting := fmt.Errorf("error decrypting password. Cause: %v", err)
		return useCase.registerUserValidationError(userCtx, fmt.Sprintf(ErrInvalidUserMsg, userCtx.DBUser.Username, err.Error()), errDecrypting)
	}
	logUserContextWithIndex(userCtx, "creating user in instance", false)
	userToCreate := &connector.DatabaseUser{
		Username: userCtx.DBUser.Username,
		Password: decryptedPwd,
		Role:     userCtx.DBUser.DatabaseRoleName,
	}
	err = userCtx.TargetInstance.CreateUser(userToCreate)
	if err != nil {
		errConnection := fmt.Errorf("connection failed with instance. Cause: %v", err)
		return useCase.registerUserValidationError(userCtx, fmt.Sprintf(ErrConnectionFailedMsg, userCtx.InstanceCtx.Instance.Name, err.Error()), errConnection)
	}

	logUserContextWithIndex(userCtx, "user created successfully!", false)
	logMsg := fmt.Sprintf(UserCreatedMsg, userCtx.DBUser.Username, userCtx.InstanceCtx.Instance.Name)
	errLog := useCase.newLog(userCtx.InstanceCtx.Instance.ID, userCtx.DBUser.ID, "", userCtx.OperationUserID, logMsg, true)
	if errLog != nil {
		return errLog
	}

	return nil
}

func (useCase *GrantAccessPermissionUseCase) processDatabases(userCtx *userContextOnGrant) error {
	databases, err := useCase.fetchDatabases(userCtx)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	databasesQty := len(databases)
	wg.Add(databasesQty)
	logUserContextWithIndex(userCtx, fmt.Sprintf("starting processing %d databases...", databasesQty), false)
	defer logUserContextWithIndex(userCtx, fmt.Sprintf("finished processing %d databases", databasesQty), false)
	for idx, database := range databases {
		go func(database *entity.Database, databaseIndex int) {
			defer wg.Done()

			databaseCtx := newGrantAccessDatabaseContext(userCtx, database, databaseIndex, databasesQty)
			if err := useCase.processDatabase(databaseCtx); err != nil {
				userCtx.InstanceCtx.GlobalCtx.GlobalErrChan <- err
			}
		}(database, idx)
	}
	wg.Wait()

	return nil
}

func (useCase *GrantAccessPermissionUseCase) fetchDatabases(userCtx *userContextOnGrant) ([]*entity.Database, error) {
	var databases []*entity.Database
	var err error
	instance := userCtx.InstanceCtx.Instance
	databaseIds := userCtx.InstanceCtx.GlobalCtx.DBIdsByInstance[instance.ID]
	if len(databaseIds) == 0 {
		databases, err = useCase.DatabaseStorage.FindAllEnabled(instance.ID)
	} else {
		databases, err = useCase.DatabaseStorage.FindAll(instance.ID, databaseIds)
	}
	if err != nil {
		errFetchingDBs := fmt.Errorf("error when fetching databases for instance '%s'. Cause: %v", instance.Name, err)
		return nil, useCase.registerUserValidationError(userCtx, fmt.Sprintf(ErrFetchingDatabasesMsg, instance.Name, err.Error()), errFetchingDBs)
	}
	if len(databases) == 0 {
		return nil, useCase.registerUserValidationError(userCtx, fmt.Sprintf(ErrNoDatabasesFoundMsg, instance.Name), fmt.Errorf("no databases found for instance '%s'", instance.Name))
	}

	return databases, err
}

func (useCase *GrantAccessPermissionUseCase) processDatabase(databaseCtx *databaseContextOnGrant) error {
	if err := useCase.validateDatabase(databaseCtx); err != nil {
		// If the database is forbidden, it's not considered an error and should be ignored for granting permissions.
		if errors.Is(err, ErrDatabaseForbidden) {
			return nil
		}
		return err
	}

	dbUserDTO := databaseCtx.UserCtx.DBUser
	instanceDTO := databaseCtx.UserCtx.InstanceCtx.Instance
	targetDatabase, _ := connector.NewDatabaseConnector(instanceDTO, databaseCtx.Database.Name)

	logDatabaseContextWithIndex(databaseCtx, "granting connect permission to user", false)
	err := targetDatabase.GrantConnect(dbUserDTO.Username)
	if err != nil {
		logMsgPt := fmt.Sprintf(ErrGrantConnectFailedMsg, dbUserDTO.Username, databaseCtx.Database.Name, instanceDTO.Name, err.Error())
		return useCase.registerDatabaseValidationError(databaseCtx, logMsgPt, fmt.Errorf("grant connect failed with instance. Cause: %v", err))
	}

	logDatabaseContextWithIndex(databaseCtx, "connect permission granted to user successfully!", false)
	msgGranted := fmt.Sprintf(PermissionGrantedMsg, dbUserDTO.Username, databaseCtx.Database.Name, instanceDTO.Name)
	err = useCase.newLog(instanceDTO.ID, dbUserDTO.ID, databaseCtx.Database.ID.String(), databaseCtx.OperationUserID, msgGranted, true)
	if err != nil {
		logDatabaseContextWithIndex(databaseCtx, fmt.Sprintf("could not create log. Cause: %v", err), true)
		return err
	}

	accessPermission, err := entity.NewAccessPermission(databaseCtx.Database.ID.String(), dbUserDTO.ID, databaseCtx.OperationUserID)
	if err != nil {
		logDatabaseContextWithIndex(databaseCtx, fmt.Sprintf("could not create access permission. Cause: %v", err), true)
		return err
	}

	return useCase.AccessPermissionStorage.Save(accessPermission)
}

func (useCase *GrantAccessPermissionUseCase) validateDatabase(databaseCtx *databaseContextOnGrant) error {
	currentDBName := databaseCtx.Database.Name
	currentUser := databaseCtx.UserCtx.DBUser.Username
	if databaseCtx.UserCtx.InstanceCtx.GlobalCtx.ForbiddenDatabases[currentDBName] && !entity.CheckRoleApplication(databaseCtx.UserCtx.DBUser.DatabaseRoleName) {
		logMsgPt := fmt.Sprintf(ErrDatabaseForbiddenMsg, currentDBName, currentUser)
		return useCase.registerDatabaseValidationError(databaseCtx, logMsgPt, ErrDatabaseForbidden)
	}
	exists, err := useCase.AccessPermissionStorage.Exists(databaseCtx.Database.ID.String(), databaseCtx.UserCtx.DBUser.ID)
	if err != nil {
		logDatabaseContextWithIndex(databaseCtx, fmt.Sprintf("could not check if user already has permission. Cause: %v", err), true)
		return err
	}
	instanceFromDB := databaseCtx.UserCtx.InstanceCtx.Instance
	if exists {
		logMsgPt := fmt.Sprintf(ErrUserAlreadyHasPermissionMsg, currentUser, currentDBName, instanceFromDB.Name)
		return useCase.registerDatabaseValidationError(databaseCtx, logMsgPt, ErrUserAlreadyHasPermission)
	}
	if !databaseCtx.Database.Enabled {
		logMsgPt := fmt.Sprintf(ErrDatabaseDisabledMsg, currentDBName, instanceFromDB.Name)
		return useCase.registerDatabaseValidationError(databaseCtx, logMsgPt, ErrDatabaseDisabled)
	}
	if !databaseCtx.Database.RolesConfigured {
		logMsgPt := fmt.Sprintf(ErrRolesNotConfiguredMsg, currentDBName, instanceFromDB.Name)
		return useCase.registerDatabaseValidationError(databaseCtx, logMsgPt, ErrRolesNotConfigured)
	}

	return nil
}

func (useCase *GrantAccessPermissionUseCase) registerInstanceValidationError(instanceCtx *instanceContextOnGrant, logMsgPt string, errorToThrow error) error {
	logInstanceContextWithIndex(instanceCtx, errorToThrow.Error(), true)
	return useCase.registerLogAndThrowError(instanceCtx.Instance.ID, "", "", instanceCtx.GlobalCtx.OperationUserID, logMsgPt, errorToThrow)
}

func (useCase *GrantAccessPermissionUseCase) registerUserValidationError(userCtx *userContextOnGrant, logMsgPt string, errorToThrow error) error {
	logUserContextWithIndex(userCtx, errorToThrow.Error(), true)
	return useCase.registerLogAndThrowError(userCtx.InstanceCtx.Instance.ID, userCtx.DBUser.ID, "", userCtx.OperationUserID, logMsgPt, errorToThrow)
}

func (useCase *GrantAccessPermissionUseCase) registerDatabaseValidationError(databaseCtx *databaseContextOnGrant, logMsgPt string, errorToThrow error) error {
	logDatabaseContextWithIndex(databaseCtx, errorToThrow.Error(), true)
	instanceID := databaseCtx.UserCtx.InstanceCtx.Instance.ID
	return useCase.registerLogAndThrowError(instanceID, databaseCtx.UserCtx.DBUser.ID, databaseCtx.Database.ID.String(), databaseCtx.OperationUserID, logMsgPt, errorToThrow)
}

func (useCase *GrantAccessPermissionUseCase) registerLogAndThrowError(instanceID, dbUserID, databaseID, operationUserID, logMsgPt string, errorToThrow error) error {
	errLog := useCase.newLog(instanceID, dbUserID, databaseID, operationUserID, logMsgPt, false)
	if errLog != nil {
		return errLog
	}
	return errorToThrow
}

func (useCase *GrantAccessPermissionUseCase) newLog(instanceID, dbUserID, databaseID, operationUserID, message string, success bool) error {
	grantLog, err := entity.NewAccessPermissionLog(instanceID, dbUserID, databaseID, message, operationUserID, success)
	if err != nil {
		log.Printf("Error when creating grant log for instance %s and database user %s. Cause: %v", instanceID, dbUserID, err)
		return err
	}
	err = useCase.AccessPermissionStorage.SaveLog(grantLog)
	if err != nil {
		log.Printf("Error when saving grant log for instance %s and database user %s. Cause: %v", instanceID, dbUserID, err)
		return err
	}

	return nil
}

func logInstanceContextWithIndex(instanceCtx *instanceContextOnGrant, message string, isError bool) {
	if isError {
		message = fmt.Sprintf("ERROR: %s", message)
	}
	log.Printf("[%d/%d] instance | [%s]: %s", instanceCtx.InstanceIndex+1, instanceCtx.GlobalCtx.InstancesQty, instanceCtx.Instance.Name, message)
}

func logUserContextWithIndex(userCtx *userContextOnGrant, message string, isError bool) {
	if isError {
		message = fmt.Sprintf("ERROR: %s", message)
	}
	log.Printf("[%d/%d] instance [%d/%d] user | [%s # %s]: %s", userCtx.InstanceCtx.InstanceIndex+1, userCtx.InstanceCtx.GlobalCtx.InstancesQty,
		userCtx.UserIndex+1, userCtx.InstanceCtx.GlobalCtx.UsersQty, userCtx.InstanceCtx.Instance.Name, userCtx.DBUser.Username, message)
}

func logDatabaseContextWithIndex(databaseCtx *databaseContextOnGrant, message string, isError bool) {
	instanceCtx := databaseCtx.UserCtx.InstanceCtx
	globalCtx := instanceCtx.GlobalCtx
	if isError {
		message = fmt.Sprintf("ERROR: %s", message)
	}
	log.Printf("[%d/%d] instance [%d/%d] user [%d/%d] database | [%s # %s # %s]: %s",
		instanceCtx.InstanceIndex+1, globalCtx.InstancesQty,
		databaseCtx.UserCtx.UserIndex+1, globalCtx.UsersQty,
		databaseCtx.DatabaseIndex+1, databaseCtx.DatabasesQty,
		instanceCtx.Instance.Name, databaseCtx.UserCtx.DBUser.Username, databaseCtx.Database.Name, message)
}

func buildGrantAccessOutput(errCh chan error) *dto.GrantAccessOutputDTO {
	output := &dto.GrantAccessOutputDTO{
		HasErrors: false,
		Message:   AccessGrantedMsg,
	}
	for err := range errCh {
		if err != nil {
			output.HasErrors = true
			output.Message = SomeErrorsDuringProcessMsg
		}
	}
	return output
}

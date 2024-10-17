package accesspermission

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/zgsolucoes/zg-data-guard/internal/database/connector"
	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	"github.com/zgsolucoes/zg-data-guard/internal/entity"
	"github.com/zgsolucoes/zg-data-guard/testdata"
	"github.com/zgsolucoes/zg-data-guard/testdata/mocks"
)

func TestGivenAnErrorInDbWhenFetchingDBUsers_WhenExecuteGrantAccess_ThenShouldReturnError(t *testing.T) {
	dbUserStorage := new(mocks.DatabaseUserStorageMock)
	dbUserStorage.On("FindAllDTOs", []string{}).Return([]*dto.DatabaseUserOutputDTO{}, sql.ErrConnDone).Once()

	uc := NewGrantAccessPermissionUseCase(nil, dbUserStorage, nil, nil, nil)
	output, err := uc.Execute(dto.GrantAccessInputDTO{DatabaseUsersIDs: []string{}}, mocks.UserID)

	assert.Error(t, err, "error expected when some error in db")
	assert.EqualError(t, err, sql.ErrConnDone.Error())
	assert.Nil(t, output)
	dbUserStorage.AssertNumberOfCalls(t, "FindAllDTOs", 1)
}

func TestGivenAnErrorInDbWhenFetchingInstances_WhenExecuteGrantAccess_ThenShouldReturnError(t *testing.T) {
	dbUserStorage := new(mocks.DatabaseUserStorageMock)
	dbUserStorage.On("FindAllDTOs", []string{}).Return([]*dto.DatabaseUserOutputDTO{}, nil).Once()
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	dbInstanceStorage.On("FindAllDTOs", "", "", []string{}).Return([]*dto.DatabaseInstanceOutputDTO{}, sql.ErrConnDone).Once()

	uc := NewGrantAccessPermissionUseCase(nil, dbUserStorage, dbInstanceStorage, nil, nil)
	output, err := uc.Execute(dto.GrantAccessInputDTO{DatabaseUsersIDs: []string{}}, mocks.UserID)

	assert.Error(t, err, "error expected when some error in db")
	assert.EqualError(t, err, sql.ErrConnDone.Error())
	assert.Nil(t, output)
	dbUserStorage.AssertNumberOfCalls(t, "FindAllDTOs", 1)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindAllDTOs", 1)
}

func TestGivenAnErrorInDbWhenFetchingForbiddenDatabases_WhenExecuteGrantAccess_ThenShouldReturnError(t *testing.T) {
	dbUserStorage := new(mocks.DatabaseUserStorageMock)
	dbUserStorage.On("FindAllDTOs", []string{}).Return([]*dto.DatabaseUserOutputDTO{}, nil).Once()
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	dbInstanceStorage.On("FindAllDTOs", "", "", []string{}).Return([]*dto.DatabaseInstanceOutputDTO{}, nil).Once()
	forbiddenObjStorage := new(mocks.ForbiddenObjectsStorageMock)
	forbiddenObjStorage.On("FindAllDatabases").Return(mocks.BuildForbiddenDatabasesList(), sql.ErrConnDone).Once()

	uc := NewGrantAccessPermissionUseCase(nil, dbUserStorage, dbInstanceStorage, nil, forbiddenObjStorage)
	output, err := uc.Execute(dto.GrantAccessInputDTO{DatabaseUsersIDs: []string{}}, mocks.UserID)

	assert.Error(t, err, "error expected when some error in db")
	assert.EqualError(t, err, fmt.Errorf("error when fetching forbidden databases. Cause: %v", sql.ErrConnDone).Error())
	assert.Nil(t, output)
	dbUserStorage.AssertNumberOfCalls(t, "FindAllDTOs", 1)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindAllDTOs", 1)
	forbiddenObjStorage.AssertNumberOfCalls(t, "FindAllDatabases", 1)
}

func TestGivenAnErrorInDbWhenSavingLog_WhenExecuteGrantAccess_ThenShouldReturnOutputError(t *testing.T) {
	instance := mocks.BuildQAInstanceDTO()
	dbUserStorage := new(mocks.DatabaseUserStorageMock)
	dbUserStorage.On("FindAllDTOs", []string{}).Return([]*dto.DatabaseUserOutputDTO{}, nil).Once()
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	dbInstanceStorage.On("FindAllDTOs", "", "", []string{instance.ID}).Return([]*dto.DatabaseInstanceOutputDTO{instance}, nil).Once()
	accessPermissionStorage := new(mocks.AccessPermissionStorageMock)
	accessPermissionStorage.On("SaveLog", mock.Anything).Return(sql.ErrConnDone).Once()
	forbiddenObjStorage := new(mocks.ForbiddenObjectsStorageMock)
	forbiddenObjStorage.On("FindAllDatabases").Return(mocks.BuildForbiddenDatabasesList(), nil).Once()

	uc := NewGrantAccessPermissionUseCase(accessPermissionStorage, dbUserStorage, dbInstanceStorage, nil, forbiddenObjStorage)
	output, err := uc.Execute(dto.GrantAccessInputDTO{DatabaseUsersIDs: []string{}, InstancesData: []dto.InstanceDataDTO{{DatabaseInstanceID: instance.ID}}}, mocks.UserID)

	assert.NoError(t, err, "the error should be in the output, not in the process and has to be logged")
	assert.NotNil(t, output)
	assert.True(t, output.HasErrors)
	assert.Equal(t, SomeErrorsDuringProcessMsg, output.Message)
	dbUserStorage.AssertNumberOfCalls(t, "FindAllDTOs", 1)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindAllDTOs", 1)
	forbiddenObjStorage.AssertNumberOfCalls(t, "FindAllDatabases", 1)
	accessPermissionStorage.AssertNumberOfCalls(t, "SaveLog", 1)
}

func TestGivenADisabledInstance_WhenExecuteGrantAccess_ThenShouldReturnOutputError(t *testing.T) {
	instance := mocks.BuildQAInstanceDTO()
	expectedLogMsg := fmt.Sprintf(ErrInstanceDisabledMsg, instance.Name)
	runGrantLoggingSingleError(t, nil, instance, nil, expectedLogMsg, false)
}

func TestGivenAnInstanceWithoutRolesCreated_WhenExecuteGrantAccess_ThenShouldReturnOutputError(t *testing.T) {
	instance := mocks.BuildAzInstanceDTO()
	instance.RolesCreated = false
	expectedLogMsg := fmt.Sprintf(ErrRolesNotCreatedMsg, instance.Name)
	runGrantLoggingSingleError(t, nil, instance, nil, expectedLogMsg, false)
}

func TestGivenAnErrorWhenCreatingConnector_WhenExecuteGrantAccess_ThenShouldReturnOutputError(t *testing.T) {
	instance := mocks.BuildConnectorNotImplementedInstance()
	expectedLogMsg := fmt.Sprintf(ErrCreatingConnectorMsg, instance.Name, "the database technology 'mysql' don't have a connector implemented")
	runGrantLoggingSingleError(t, nil, instance, nil, expectedLogMsg, false)
}

func TestGivenADisabledUser_WhenExecuteGrantAccess_ThenShouldReturnOutputError(t *testing.T) {
	dbUser := mocks.BuildDisabledDbUserJohnDTO()
	instance := mocks.BuildAzInstanceDTO()
	expectedLogMsg := fmt.Sprintf(ErrUserDisabledMsg, dbUser.Username)
	runGrantLoggingSingleError(t, dbUser, instance, nil, expectedLogMsg, false)
}

func TestGivenAUserWithInvalidRole_WhenExecuteGrantAccess_ThenShouldReturnOutputError(t *testing.T) {
	dbUser := mocks.BuildDbUserDummyDTO()
	dbUser.DatabaseRoleName = "invalid"
	instance := mocks.BuildAzInstanceDTO()
	expectedLogMsg := fmt.Sprintf(ErrInvalidRoleMsg, dbUser.DatabaseRoleName, dbUser.Username)
	runGrantLoggingSingleError(t, dbUser, instance, nil, expectedLogMsg, false)
}

func TestGivenAnErrorConnectingToInstance_WhenExecuteGrantAccess_ThenShouldReturnOutputError(t *testing.T) {
	dbUser := mocks.BuildDbUserJohnDTO()
	instance := mocks.BuildDummyErrorInstance()
	expectedLogMsg := fmt.Sprintf(ErrConnectionFailedMsg, instance.Name, "error checking user existence in "+instance.Name)
	runGrantLoggingSingleError(t, dbUser, instance, nil, expectedLogMsg, false)
}

func TestGivenAnErrorWhenCreatingUser_WhenExecuteGrantAccess_ThenShouldReturnOutputError(t *testing.T) {
	dbUser := mocks.BuildDbUserDummyErrorCreateDTO()
	instance := mocks.BuildAzInstanceDTO()
	expectedLogMsg := fmt.Sprintf(ErrConnectionFailedMsg, instance.Name, fmt.Errorf("%w: Instance(%s) - User(%s)", connector.ErrCreateUser, instance.Name, dbUser.Username).Error())
	runGrantLoggingSingleError(t, dbUser, instance, nil, expectedLogMsg, false)
}

func TestGivenAnErrorInDBWhenFetchingDatabases_WhenExecuteGrantAccess_ThenShouldReturnOutputError(t *testing.T) {
	dbUser := mocks.BuildDbUserDummyDTO()
	instance := mocks.BuildAzInstanceDTO()
	dbUserStorage := new(mocks.DatabaseUserStorageMock)
	dbUserStorage.On("FindAllDTOs", []string{dbUser.ID}).Return([]*dto.DatabaseUserOutputDTO{dbUser}, nil).Once()
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	dbInstanceStorage.On("FindAllDTOs", "", "", []string{instance.ID}).Return([]*dto.DatabaseInstanceOutputDTO{instance}, nil).Once()
	accessPermissionStorage := new(mocks.AccessPermissionStorageMock)
	userCreatedLog, _ := entity.NewAccessPermissionLog(instance.ID, dbUser.ID, "", fmt.Sprintf(UserCreatedMsg, dbUser.Username, instance.Name), mocks.UserID, true)
	errFetchingDbsLog, _ := entity.NewAccessPermissionLog(instance.ID, dbUser.ID, "", fmt.Sprintf(ErrFetchingDatabasesMsg, instance.Name, sql.ErrConnDone.Error()), mocks.UserID, false)
	accessPermissionStorage.On("SaveLog", testdata.CompareLogs(userCreatedLog)).Return(nil).Once()
	accessPermissionStorage.On("SaveLog", testdata.CompareLogs(errFetchingDbsLog)).Return(nil).Once()
	dbStorage := new(mocks.DatabaseStorageMock)
	dbStorage.On("FindAllEnabled", instance.ID).Return([]*entity.Database{}, sql.ErrConnDone).Once()
	forbiddenObjStorage := new(mocks.ForbiddenObjectsStorageMock)
	forbiddenObjStorage.On("FindAllDatabases").Return(mocks.BuildForbiddenDatabasesList(), nil).Once()

	uc := NewGrantAccessPermissionUseCase(accessPermissionStorage, dbUserStorage, dbInstanceStorage, dbStorage, forbiddenObjStorage)
	output, err := uc.Execute(dto.GrantAccessInputDTO{DatabaseUsersIDs: []string{dbUser.ID}, InstancesData: []dto.InstanceDataDTO{{DatabaseInstanceID: instance.ID}}}, mocks.UserID)

	assert.NoError(t, err, "the error should be in the output, not in the process and has to be logged")
	assert.NotNil(t, output)
	assert.True(t, output.HasErrors)
	assert.Equal(t, SomeErrorsDuringProcessMsg, output.Message)
	dbUserStorage.AssertNumberOfCalls(t, "FindAllDTOs", 1)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindAllDTOs", 1)
	dbStorage.AssertNumberOfCalls(t, "FindAllEnabled", 1)
	forbiddenObjStorage.AssertNumberOfCalls(t, "FindAllDatabases", 1)
	accessPermissionStorage.AssertNumberOfCalls(t, "SaveLog", 2)
}

func TestGivenAnInstanceWithoutDatabases_WhenExecuteGrantAccess_ThenShouldReturnOutputError(t *testing.T) {
	dbUser := mocks.BuildDbUserJohnDTO()
	instance := mocks.BuildAzInstanceDTO()
	expectedLogMsg := fmt.Sprintf(ErrNoDatabasesFoundMsg, instance.Name)
	runGrantLoggingSingleError(t, dbUser, instance, nil, expectedLogMsg, false)
}

func TestGivenADisabledDatabase_WhenExecuteGrantAccess_ThenShouldReturnOutputError(t *testing.T) {
	dbUser := mocks.BuildDbUserJohnDTO()
	instance := mocks.BuildAzInstanceDTO()
	database := mocks.BuildSettingsDatabase()
	database.Enabled = false
	expectedLogMsg := fmt.Sprintf(ErrDatabaseDisabledMsg, database.Name, instance.Name)
	runGrantLoggingSingleError(t, dbUser, instance, database, expectedLogMsg, false)
}

func TestGivenADatabaseWithoutRolesConfigured_WhenExecuteGrantAccess_ThenShouldReturnOutputError(t *testing.T) {
	dbUser := mocks.BuildDbUserJohnDTO()
	instance := mocks.BuildAzInstanceDTO()
	database := mocks.BuildSettingsDatabase()
	database.RolesConfigured = false
	expectedLogMsg := fmt.Sprintf(ErrRolesNotConfiguredMsg, database.Name, instance.Name)
	runGrantLoggingSingleError(t, dbUser, instance, database, expectedLogMsg, false)
}

func TestGivenAnAccessAlreadyExistent_WhenExecuteGrantAccess_ThenShouldReturnOutputError(t *testing.T) {
	dbUser := mocks.BuildDbUserJohnDTO()
	instance := mocks.BuildAzInstanceDTO()
	database := mocks.BuildSettingsDatabase()
	expectedLogMsg := fmt.Sprintf(ErrUserAlreadyHasPermissionMsg, dbUser.Username, database.Name, instance.Name)
	runGrantLoggingSingleError(t, dbUser, instance, database, expectedLogMsg, true)
}

func TestGivenAnErrorWhenGrantingConnect_WhenExecuteGrantAccess_ThenShouldReturnOutputError(t *testing.T) {
	dbUser := mocks.BuildDbUserDummyErrorGrantDTO()
	instance := mocks.BuildAzInstanceDTO()
	database := mocks.BuildSettingsDatabase()
	errMsg := fmt.Errorf("%w: Instance(%s) - User(%s)", connector.ErrGrantConnect, instance.Name, dbUser.Username).Error()
	expectedLogMsg := fmt.Sprintf(ErrGrantConnectFailedMsg, dbUser.Username, database.Name, instance.Name, errMsg)
	runGrantLoggingSingleError(t, dbUser, instance, database, expectedLogMsg, false)
}

func TestGivenAnErrorWhenCheckingUserPermission_WhenExecuteGrantAccess_ThenShouldReturnOutputError(t *testing.T) {
	dbUser := mocks.BuildDbUserJohnDTO()
	instance := mocks.BuildAzInstanceDTO()
	database := mocks.BuildSettingsDatabase()
	dbID := database.ID.String()
	dbUserStorage := new(mocks.DatabaseUserStorageMock)
	dbUserStorage.On("FindAllDTOs", []string{dbUser.ID}).Return([]*dto.DatabaseUserOutputDTO{dbUser}, nil).Once()
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	dbInstanceStorage.On("FindAllDTOs", "", "", []string{instance.ID}).Return([]*dto.DatabaseInstanceOutputDTO{instance}, nil).Once()
	accessPermissionStorage := new(mocks.AccessPermissionStorageMock)
	accessPermissionStorage.On("Exists", dbID, dbUser.ID).Return(false, sql.ErrConnDone).Once()
	dbStorage := new(mocks.DatabaseStorageMock)
	dbStorage.On("FindAll", instance.ID, []string{dbID}).Return([]*entity.Database{database}, nil).Once()
	forbiddenObjStorage := new(mocks.ForbiddenObjectsStorageMock)
	forbiddenObjStorage.On("FindAllDatabases").Return(mocks.BuildForbiddenDatabasesList(), nil).Once()

	uc := NewGrantAccessPermissionUseCase(accessPermissionStorage, dbUserStorage, dbInstanceStorage, dbStorage, forbiddenObjStorage)
	output, err := uc.Execute(buildGrantInput(dbUser, instance, []*entity.Database{database}), mocks.UserID)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.True(t, output.HasErrors)
	assert.Equal(t, SomeErrorsDuringProcessMsg, output.Message)
	dbUserStorage.AssertNumberOfCalls(t, "FindAllDTOs", 1)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindAllDTOs", 1)
	forbiddenObjStorage.AssertNumberOfCalls(t, "FindAllDatabases", 1)
	dbStorage.AssertNumberOfCalls(t, "FindAll", 1)
	accessPermissionStorage.AssertNumberOfCalls(t, "Exists", 1)
}

func TestGivenAnErrorWhenCreatingLog_WhenExecuteGrantAccess_ThenShouldReturnOutputError(t *testing.T) {
	dbUser := mocks.BuildDbUserJohnDTO()
	instance := mocks.BuildAzInstanceDTO()
	database := mocks.BuildSettingsDatabase()
	dbID := database.ID.String()
	dbUserStorage := new(mocks.DatabaseUserStorageMock)
	dbUserStorage.On("FindAllDTOs", []string{dbUser.ID}).Return([]*dto.DatabaseUserOutputDTO{dbUser}, nil).Once()
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	dbInstanceStorage.On("FindAllDTOs", "", "", []string{instance.ID}).Return([]*dto.DatabaseInstanceOutputDTO{instance}, nil).Once()
	accessPermissionStorage := new(mocks.AccessPermissionStorageMock)
	accessPermissionStorage.On("Exists", dbID, dbUser.ID).Return(false, nil).Once()
	dbStorage := new(mocks.DatabaseStorageMock)
	dbStorage.On("FindAll", instance.ID, []string{dbID}).Return([]*entity.Database{database}, nil).Once()
	forbiddenObjStorage := new(mocks.ForbiddenObjectsStorageMock)
	forbiddenObjStorage.On("FindAllDatabases").Return(mocks.BuildForbiddenDatabasesList(), nil).Once()

	uc := NewGrantAccessPermissionUseCase(accessPermissionStorage, dbUserStorage, dbInstanceStorage, dbStorage, forbiddenObjStorage)
	output, err := uc.Execute(buildGrantInput(dbUser, instance, []*entity.Database{database}), "")

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.True(t, output.HasErrors)
	assert.Equal(t, SomeErrorsDuringProcessMsg, output.Message)
	dbUserStorage.AssertNumberOfCalls(t, "FindAllDTOs", 1)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindAllDTOs", 1)
	forbiddenObjStorage.AssertNumberOfCalls(t, "FindAllDatabases", 1)
	dbStorage.AssertNumberOfCalls(t, "FindAll", 1)
	accessPermissionStorage.AssertNumberOfCalls(t, "Exists", 1)
}

func TestGivenValidInputWithForbiddenDatabases_WhenExecuteGrantAccess_ThenShouldReturnOutputSuccess(t *testing.T) {
	dbUser := mocks.BuildDbUserJohnDTO()
	instance := mocks.BuildAzInstanceDTO()
	databases := mocks.BuildDatabaseListSameInstanceAndOnlyEnabled()
	notForbiddenDB := databases[2]
	allowedDbID := notForbiddenDB.ID.String()
	forbiddenDbID := databases[0].ID.String()
	forbiddenDbID2 := databases[1].ID.String()
	dbUserStorage := new(mocks.DatabaseUserStorageMock)
	dbUserStorage.On("FindAllDTOs", []string{dbUser.ID}).Return([]*dto.DatabaseUserOutputDTO{dbUser}, nil).Once()
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	dbInstanceStorage.On("FindAllDTOs", "", "", []string{instance.ID}).Return([]*dto.DatabaseInstanceOutputDTO{instance}, nil).Once()
	accessPermissionStorage := new(mocks.AccessPermissionStorageMock)
	logMsg := fmt.Sprintf(PermissionGrantedMsg, dbUser.Username, notForbiddenDB.Name, instance.Name)
	logForbiddenMsg := fmt.Sprintf(ErrDatabaseForbiddenMsg, databases[0].Name, dbUser.Username)
	logForbiddenMsg2 := fmt.Sprintf(ErrDatabaseForbiddenMsg, databases[1].Name, dbUser.Username)
	expectedLog, _ := entity.NewAccessPermissionLog(instance.ID, dbUser.ID, allowedDbID, logMsg, mocks.UserID, true)
	forbiddenLog, _ := entity.NewAccessPermissionLog(instance.ID, dbUser.ID, forbiddenDbID, logForbiddenMsg, mocks.UserID, false)
	forbiddenLog2, _ := entity.NewAccessPermissionLog(instance.ID, dbUser.ID, forbiddenDbID2, logForbiddenMsg2, mocks.UserID, false)
	accessPermissionStorage.On("Exists", allowedDbID, dbUser.ID).Return(false, nil).Once()
	accessPermissionStorage.On("SaveLog", testdata.CompareLogs(expectedLog)).Return(nil).Once()
	accessPermissionStorage.On("SaveLog", testdata.CompareLogs(forbiddenLog)).Return(nil).Once()
	accessPermissionStorage.On("SaveLog", testdata.CompareLogs(forbiddenLog2)).Return(nil).Once()
	expectedAccess, _ := entity.NewAccessPermission(allowedDbID, dbUser.ID, mocks.UserID)
	accessPermissionStorage.On("Save", compareAccessPermission(expectedAccess)).Return(nil).Once()
	dbStorage := new(mocks.DatabaseStorageMock)
	dbStorage.On("FindAll", instance.ID, []string{allowedDbID, forbiddenDbID, forbiddenDbID2}).Return(databases, nil).Once()
	forbiddenObjStorage := new(mocks.ForbiddenObjectsStorageMock)
	forbiddenObjStorage.On("FindAllDatabases").Return(mocks.BuildForbiddenDatabasesList(), nil).Once()

	uc := NewGrantAccessPermissionUseCase(accessPermissionStorage, dbUserStorage, dbInstanceStorage, dbStorage, forbiddenObjStorage)
	output, err := uc.Execute(buildGrantInput(dbUser, instance, databases), mocks.UserID)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.False(t, output.HasErrors)
	assert.Equal(t, AccessGrantedMsg, output.Message)
	dbUserStorage.AssertNumberOfCalls(t, "FindAllDTOs", 1)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindAllDTOs", 1)
	forbiddenObjStorage.AssertNumberOfCalls(t, "FindAllDatabases", 1)
	dbStorage.AssertNumberOfCalls(t, "FindAll", 1)
	accessPermissionStorage.AssertNumberOfCalls(t, "Exists", 1)
	accessPermissionStorage.AssertNumberOfCalls(t, "SaveLog", 3)
	accessPermissionStorage.AssertNumberOfCalls(t, "Save", 1)
}

func TestGivenValidInput_WhenExecuteGrantAccess_ThenShouldReturnOutputSuccess(t *testing.T) {
	dbUser := mocks.BuildDbUserJohnDTO()
	instance := mocks.BuildAzInstanceDTO()
	database := mocks.BuildSettingsDatabase()
	dbID := database.ID.String()
	dbUserStorage := new(mocks.DatabaseUserStorageMock)
	dbUserStorage.On("FindAllDTOs", []string{dbUser.ID}).Return([]*dto.DatabaseUserOutputDTO{dbUser}, nil).Once()
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	dbInstanceStorage.On("FindAllDTOs", "", "", []string{instance.ID}).Return([]*dto.DatabaseInstanceOutputDTO{instance}, nil).Once()
	accessPermissionStorage := new(mocks.AccessPermissionStorageMock)
	logMsg := fmt.Sprintf(PermissionGrantedMsg, dbUser.Username, database.Name, instance.Name)
	expectedLog, _ := entity.NewAccessPermissionLog(instance.ID, dbUser.ID, dbID, logMsg, mocks.UserID, true)
	accessPermissionStorage.On("Exists", dbID, dbUser.ID).Return(false, nil).Once()
	accessPermissionStorage.On("SaveLog", testdata.CompareLogs(expectedLog)).Return(nil).Once()
	expectedAccess, _ := entity.NewAccessPermission(dbID, dbUser.ID, mocks.UserID)
	accessPermissionStorage.On("Save", compareAccessPermission(expectedAccess)).Return(nil).Once()
	dbStorage := new(mocks.DatabaseStorageMock)
	dbStorage.On("FindAll", instance.ID, []string{dbID}).Return([]*entity.Database{database}, nil).Once()
	forbiddenObjStorage := new(mocks.ForbiddenObjectsStorageMock)
	forbiddenObjStorage.On("FindAllDatabases").Return(mocks.BuildForbiddenDatabasesList(), nil).Once()

	uc := NewGrantAccessPermissionUseCase(accessPermissionStorage, dbUserStorage, dbInstanceStorage, dbStorage, forbiddenObjStorage)
	output, err := uc.Execute(buildGrantInput(dbUser, instance, []*entity.Database{database}), mocks.UserID)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.False(t, output.HasErrors)
	assert.Equal(t, AccessGrantedMsg, output.Message)
	dbUserStorage.AssertNumberOfCalls(t, "FindAllDTOs", 1)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindAllDTOs", 1)
	forbiddenObjStorage.AssertNumberOfCalls(t, "FindAllDatabases", 1)
	dbStorage.AssertNumberOfCalls(t, "FindAll", 1)
	accessPermissionStorage.AssertNumberOfCalls(t, "Exists", 1)
	accessPermissionStorage.AssertNumberOfCalls(t, "SaveLog", 1)
	accessPermissionStorage.AssertNumberOfCalls(t, "Save", 1)
}

func runGrantLoggingSingleError(t *testing.T, dbUser *dto.DatabaseUserOutputDTO, instance *dto.DatabaseInstanceOutputDTO, db *entity.Database, expectedLogMsg string, accessExists bool) {
	dbUserID := getDBUserIDFromDTO(dbUser)
	dbID := getDbID(db)
	dbUserStorage := new(mocks.DatabaseUserStorageMock)
	dbUserStorage.On("FindAllDTOs", []string{dbUserID}).Return([]*dto.DatabaseUserOutputDTO{dbUser}, nil).Once()
	dbInstanceStorage := new(mocks.DatabaseInstanceStorageMock)
	dbInstanceStorage.On("FindAllDTOs", "", "", []string{instance.ID}).Return([]*dto.DatabaseInstanceOutputDTO{instance}, nil).Once()
	accessPermissionStorage := new(mocks.AccessPermissionStorageMock)
	expectedLog, _ := entity.NewAccessPermissionLog(instance.ID, dbUserID, dbID, expectedLogMsg, mocks.UserID, false)
	accessPermissionStorage.On("Exists", dbID, dbUserID).Return(accessExists, nil).Once()
	accessPermissionStorage.On("SaveLog", testdata.CompareLogs(expectedLog)).Return(nil).Once()
	forbiddenObjStorage := new(mocks.ForbiddenObjectsStorageMock)
	forbiddenObjStorage.On("FindAllDatabases").Return(mocks.BuildForbiddenDatabasesList(), nil).Once()
	dbStorage := new(mocks.DatabaseStorageMock)
	if db != nil {
		dbStorage.On("FindAllEnabled", instance.ID).Return([]*entity.Database{db}, nil).Once()
	} else {
		dbStorage.On("FindAllEnabled", instance.ID).Return([]*entity.Database{}, nil).Once()
	}

	uc := NewGrantAccessPermissionUseCase(accessPermissionStorage, dbUserStorage, dbInstanceStorage, dbStorage, forbiddenObjStorage)
	output, err := uc.Execute(dto.GrantAccessInputDTO{DatabaseUsersIDs: []string{dbUserID}, InstancesData: []dto.InstanceDataDTO{{DatabaseInstanceID: instance.ID}}}, mocks.UserID)

	assert.NoError(t, err, "the error should be in the output, not in the process and has to be logged")
	assert.NotNil(t, output)
	assert.True(t, output.HasErrors)
	assert.Equal(t, SomeErrorsDuringProcessMsg, output.Message)
	dbUserStorage.AssertNumberOfCalls(t, "FindAllDTOs", 1)
	dbInstanceStorage.AssertNumberOfCalls(t, "FindAllDTOs", 1)
	forbiddenObjStorage.AssertNumberOfCalls(t, "FindAllDatabases", 1)
	accessPermissionStorage.AssertNumberOfCalls(t, "SaveLog", 1)
	if db != nil {
		accessPermissionStorage.AssertNumberOfCalls(t, "Exists", 1)
		dbStorage.AssertNumberOfCalls(t, "FindAllEnabled", 1)
	}
}

func buildGrantInput(dbUser *dto.DatabaseUserOutputDTO, instance *dto.DatabaseInstanceOutputDTO, databases []*entity.Database) dto.GrantAccessInputDTO {
	databaseIDs := make([]string, len(databases))
	for i, db := range databases {
		databaseIDs[i] = db.ID.String()
	}
	return dto.GrantAccessInputDTO{
		DatabaseUsersIDs: []string{dbUser.ID},
		InstancesData: []dto.InstanceDataDTO{
			{DatabaseInstanceID: instance.ID, DatabasesIDs: databaseIDs},
		},
	}
}

func getDBUserIDFromDTO(dbUser *dto.DatabaseUserOutputDTO) string {
	if dbUser != nil {
		return dbUser.ID
	}
	return ""
}

func getDbID(db *entity.Database) string {
	if db != nil {
		return db.ID.String()
	}
	return ""
}

func compareAccessPermission(expectedAccess *entity.AccessPermission) any {
	return mock.MatchedBy(func(resultAccess *entity.AccessPermission) bool {
		// Ignore ID and Date fields for comparison
		expectedAccess.ID = resultAccess.ID
		expectedAccess.GrantedAt = resultAccess.GrantedAt
		return assert.ObjectsAreEqual(expectedAccess, resultAccess)
	})
}

package accesspermission

import (
	"errors"

	"github.com/zgsolucoes/zg-data-guard/internal/database/connector"
	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	"github.com/zgsolucoes/zg-data-guard/internal/entity"
)

var (
	ErrInstanceDisabled         = errors.New("instance is disabled")
	ErrRolesNotCreated          = errors.New("roles not created yet in instance")
	ErrUserDisabled             = errors.New("user is disabled")
	ErrDatabaseDisabled         = errors.New("database is disabled")
	ErrRolesNotConfigured       = errors.New("roles not configured yet in database")
	ErrDatabaseForbidden        = errors.New("database access is forbidden")
	ErrInvalidRole              = errors.New("invalid role defined for user")
	ErrUserAlreadyHasPermission = errors.New("user already has access permission")
)

const (
	ErrCreatingConnectorMsg         = "error creating connector with instance '%s'. Details: %s"
	ErrInstanceDisabledMsg          = "the instance '%s' is disabled"
	ErrRolesNotCreatedMsg           = "the roles have not been properly created in instance '%s' yet"
	ErrRolesNotConfiguredMsg        = "the roles have not been properly configured in the database '%s' of instance '%s' yet"
	ErrUserDisabledMsg              = "the user '%s' is disabled"
	ErrInvalidUserMsg               = "the user '%s' is invalid. Details: %s"
	ErrInvalidRoleMsg               = "the role '%s' defined for user '%s' is invalid"
	ErrFetchingDatabasesMsg         = "error fetching databases of instance '%s'. Details: %s"
	ErrNoDatabasesFoundMsg          = "no databases found for instance '%s'"
	ErrUserAlreadyHasPermissionMsg  = "the user '%s' already has access permission to the database '%s' of instance '%s'"
	ErrDatabaseDisabledMsg          = "the database '%s' is disabled in instance '%s'"
	ErrConnectionFailedMsg          = "failed to connect to instance '%s'. Details: %s"
	ErrGrantConnectFailedMsg        = "failed to grant access permission to user '%s' on database '%s' of instance '%s'. Details: %s"
	ErrRevokeAndDropUserFailedMsg   = "failed to revoke and remove user '%s' from instance '%s'. Details: %s"
	ErrDeletingAccessOfUserMsg      = "failed to delete all access for user '%s' in instance '%s'. Details: %s"
	ErrDatabaseForbiddenMsg         = "the database '%s' is blacklisted, therefore access is blocked for user '%s'"
	UserAccessRevokedAndExcludedMsg = "the user '%s' has had their access revoked and was successfully removed from instance '%s'"
	UserCreatedMsg                  = "the user '%s' was successfully created in instance '%s'"
	PermissionGrantedMsg            = "access permission granted to user '%s' on database '%s' of instance '%s'"
	InstanceDisabledSuccessMsg      = "instance '%s' has been disabled"
	InstanceEnabledSuccessMsg       = "instance '%s' has been enabled"
)

type globalContextOnGrant struct {
	DBUsers            []*dto.DatabaseUserOutputDTO
	DBIdsByInstance    map[string][]string
	OperationUserID    string
	ForbiddenDatabases map[string]bool
	GlobalErrChan      chan error
	InstancesQty       int
	UsersQty           int
}

func newGrantAccessGlobalContext(
	dbUsers []*dto.DatabaseUserOutputDTO,
	databaseIdsByInstance map[string][]string,
	operationUserID string,
	forbiddenDatabases map[string]bool,
	instancesQty, usersQty int) *globalContextOnGrant {
	bufferSize := instancesQty * usersQty
	return &globalContextOnGrant{
		DBUsers:            dbUsers,
		DBIdsByInstance:    databaseIdsByInstance,
		OperationUserID:    operationUserID,
		ForbiddenDatabases: forbiddenDatabases,
		GlobalErrChan:      make(chan error, bufferSize),
		InstancesQty:       instancesQty,
		UsersQty:           usersQty,
	}
}

type instanceContextOnGrant struct {
	GlobalCtx     *globalContextOnGrant
	Instance      *dto.DatabaseInstanceOutputDTO
	InstanceIndex int
}

func newGrantAccessInstanceContext(
	globalCtx *globalContextOnGrant,
	instanceDTO *dto.DatabaseInstanceOutputDTO,
	instanceIndex int) *instanceContextOnGrant {
	return &instanceContextOnGrant{
		GlobalCtx:     globalCtx,
		Instance:      instanceDTO,
		InstanceIndex: instanceIndex,
	}
}

type userContextOnGrant struct {
	InstanceCtx     *instanceContextOnGrant
	TargetInstance  connector.DatabaseTCPConnectorInterface
	DBUser          *dto.DatabaseUserOutputDTO
	UserIndex       int
	OperationUserID string
}

func newGrantAccessUserContext(
	instanceCtx *instanceContextOnGrant,
	targetInstance connector.DatabaseTCPConnectorInterface,
	userDTO *dto.DatabaseUserOutputDTO,
	userIndex int) *userContextOnGrant {
	return &userContextOnGrant{
		InstanceCtx:     instanceCtx,
		TargetInstance:  targetInstance,
		DBUser:          userDTO,
		UserIndex:       userIndex,
		OperationUserID: instanceCtx.GlobalCtx.OperationUserID,
	}
}

type databaseContextOnGrant struct {
	UserCtx         *userContextOnGrant
	Database        *entity.Database
	DatabaseIndex   int
	DatabasesQty    int
	OperationUserID string
}

func newGrantAccessDatabaseContext(
	userCtx *userContextOnGrant,
	database *entity.Database,
	databaseIndex, databasesQty int) *databaseContextOnGrant {
	return &databaseContextOnGrant{
		UserCtx:         userCtx,
		Database:        database,
		DatabaseIndex:   databaseIndex,
		DatabasesQty:    databasesQty,
		OperationUserID: userCtx.OperationUserID,
	}
}

type revokeAccessContext struct {
	Instance        *dto.DatabaseInstanceOutputDTO
	User            *entity.DatabaseUser
	InstancesQty    int
	InstanceIndex   int
	OperationUserID string
}

func newRevokeAccessContext(instance *dto.DatabaseInstanceOutputDTO, user *entity.DatabaseUser, instancesQty, instanceIndex int, operationUserID string) *revokeAccessContext {
	return &revokeAccessContext{
		Instance:        instance,
		User:            user,
		InstancesQty:    instancesQty,
		InstanceIndex:   instanceIndex,
		OperationUserID: operationUserID,
	}
}

type loggableRevokeResult struct {
	Err          error
	LogMessagePt string
	RevokeCtx    *revokeAccessContext
}

package accesspermission

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/zgsolucoes/zg-data-guard/internal/database/connector"
	"github.com/zgsolucoes/zg-data-guard/internal/database/storage"
	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	"github.com/zgsolucoes/zg-data-guard/internal/entity"
	"github.com/zgsolucoes/zg-data-guard/internal/usecase/common"
	"github.com/zgsolucoes/zg-data-guard/pkg/utils"
)

type RevokeAccessPermissionUseCase struct {
	AccessPermissionStorage storage.AccessPermissionStorage
	DatabaseInstanceStorage storage.DatabaseInstanceStorage
	DatabaseUserStorage     storage.DatabaseUserStorage
}

func NewRevokeAccessPermissionUseCase(
	accessStorage storage.AccessPermissionStorage,
	instanceStorage storage.DatabaseInstanceStorage,
	dbUserStorage storage.DatabaseUserStorage,
) *RevokeAccessPermissionUseCase {
	return &RevokeAccessPermissionUseCase{
		AccessPermissionStorage: accessStorage, DatabaseInstanceStorage: instanceStorage, DatabaseUserStorage: dbUserStorage}
}

// Execute godoc
/** Responsible for revoking access from the selected database instances or from all instances accessible by the user.
It revokes the user's access concurrently.
It returns an output DTO that has a flag indicating if the process has errors and a message with the result.
For each error that occurs inside the instance context during the process, it's logged, persisted and the process continues.
*/
func (useCase *RevokeAccessPermissionUseCase) Execute(input dto.RevokeAccessInputDTO, operationUserID string) (*dto.RevokeAccessOutputDTO, error) {
	userToRevoke, err := useCase.fetchDatabaseUser(input.DatabaseUserID)
	if err != nil {
		return nil, err
	}

	instancesToRevoke, err := useCase.determineInstancesToRevoke(input, operationUserID, userToRevoke)
	if err != nil {
		return nil, err
	}
	return useCase.revokeAccess(instancesToRevoke, userToRevoke, operationUserID)
}

func (useCase *RevokeAccessPermissionUseCase) fetchDatabaseUser(userID string) (*entity.DatabaseUser, error) {
	userToRevoke, err := useCase.DatabaseUserStorage.FindByID(userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, common.ErrDatabaseUserNotFound
		}
		return nil, err
	}
	return userToRevoke, nil
}

func (useCase *RevokeAccessPermissionUseCase) determineInstancesToRevoke(input dto.RevokeAccessInputDTO, opUserID string, userToRevoke *entity.DatabaseUser) ([]*dto.DatabaseInstanceOutputDTO, error) {
	var accessibleInstancesIDs, instancesIDsToRevoke []string
	accessibleInstancesIDs, err := useCase.AccessPermissionStorage.FindAllAccessibleInstancesIDsByUser(input.DatabaseUserID)
	if err != nil {
		return nil, err
	}

	if len(input.DatabaseInstancesIDs) > 0 {
		log.Printf("Revoking access from database user '%s' to the selected %d database instances. Requester: %s", input.DatabaseUserID, len(input.DatabaseInstancesIDs), opUserID)
		instancesIDsToRevoke = filterAccessibleInstances(input.DatabaseInstancesIDs, accessibleInstancesIDs, userToRevoke.Username)
	} else {
		log.Printf("Revoking all access from database user '%s' to all database instances accessible by him. Requester: %s", input.DatabaseUserID, opUserID)
		instancesIDsToRevoke = accessibleInstancesIDs
	}

	if len(instancesIDsToRevoke) == 0 {
		return nil, common.ErrNoAccessibleInstancesFound
	}
	instancesToRevoke, err := useCase.DatabaseInstanceStorage.FindAllDTOs("", "", instancesIDsToRevoke)
	if err != nil {
		return nil, err
	}
	return instancesToRevoke, nil
}

func (useCase *RevokeAccessPermissionUseCase) revokeAccess(dbInstances []*dto.DatabaseInstanceOutputDTO, dbUser *entity.DatabaseUser, operationUserID string) (*dto.RevokeAccessOutputDTO, error) {
	instancesQty := len(dbInstances)
	resultCh := make(chan *loggableRevokeResult, instancesQty)
	output := &dto.RevokeAccessOutputDTO{
		HasErrors: false,
		Message:   fmt.Sprintf("Successfully revoked access for user '%s' in %d database instances!", dbUser.Username, instancesQty),
	}

	var wg sync.WaitGroup
	wg.Add(instancesQty)
	for idx, instance := range dbInstances {
		go func(instance *dto.DatabaseInstanceOutputDTO, instanceIndex int, resultCh chan<- *loggableRevokeResult) {
			defer wg.Done()
			revokeCtx := newRevokeAccessContext(instance, dbUser, instancesQty, instanceIndex, operationUserID)
			resultCh <- useCase.revokeUserAccessAndRemoveFromInstance(revokeCtx)
		}(instance, idx, resultCh)
	}

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	useCase.processResult(resultCh, output)
	log.Printf("Revoke access process finished for user '%s' in %d database instances", dbUser.Username, instancesQty)
	return output, nil
}

func (useCase *RevokeAccessPermissionUseCase) revokeUserAccessAndRemoveFromInstance(revokeCtx *revokeAccessContext) *loggableRevokeResult {
	result := &loggableRevokeResult{RevokeCtx: revokeCtx}

	targetInstance, err := connector.NewDatabaseConnector(revokeCtx.Instance, "")
	if err != nil {
		result.Err = fmt.Errorf("could not create connector. Details: %w", err)
		result.LogMessagePt = fmt.Sprintf(ErrCreatingConnectorMsg, revokeCtx.Instance.Name, err.Error())
		return result
	}
	logRevokeContextWithIndex(revokeCtx, fmt.Sprintf("%s Revoking connection grants and removing user from instance", connector.ClusterConnectorPrefix))
	err = targetInstance.RevokeUserPrivilegesAndRemove(revokeCtx.User.Username)
	if err != nil {
		result.Err = fmt.Errorf("%s could not revoke and drop user. Details: %w", connector.ClusterConnectorPrefix, err)
		result.LogMessagePt = fmt.Sprintf(ErrRevokeAndDropUserFailedMsg, revokeCtx.User.Username, revokeCtx.Instance.Name, err.Error())
		return result
	}

	logRevokeContextWithIndex(revokeCtx, fmt.Sprintf("%s User access revoked and removed from instance", connector.ClusterConnectorPrefix))
	return result
}

func (useCase *RevokeAccessPermissionUseCase) processResult(resultCh <-chan *loggableRevokeResult, output *dto.RevokeAccessOutputDTO) {
	for loggableResult := range resultCh {
		if loggableResult.Err != nil {
			output.HasErrors = true
			output.Message = SomeErrorsDuringProcessMsg
		} else {
			dbUserID := loggableResult.RevokeCtx.User.ID.String()
			instanceID := loggableResult.RevokeCtx.Instance.ID
			err := useCase.AccessPermissionStorage.DeleteAllByUserAndInstance(dbUserID, instanceID)
			if err != nil {
				output.HasErrors = true
				output.Message = SomeErrorsDuringProcessMsg
				loggableResult.Err = fmt.Errorf("could not delete all access from database user '%s' in database instance '%s'. Cause: %w", dbUserID, instanceID, err)
				loggableResult.LogMessagePt = fmt.Sprintf(ErrDeletingAccessOfUserMsg, loggableResult.RevokeCtx.User.Username, loggableResult.RevokeCtx.Instance.Name, err.Error())
			} else {
				logRevokeContextWithIndex(loggableResult.RevokeCtx, "All user access to the respective instance has been successfully deleted!")
				loggableResult.LogMessagePt = fmt.Sprintf(UserAccessRevokedAndExcludedMsg, loggableResult.RevokeCtx.User.Username, loggableResult.RevokeCtx.Instance.Name)
			}
		}
		useCase.persistLog(*loggableResult, output)
	}
}

func (useCase *RevokeAccessPermissionUseCase) persistLog(loggableResult loggableRevokeResult, output *dto.RevokeAccessOutputDTO) {
	if loggableResult.Err != nil {
		logRevokeContextWithIndex(loggableResult.RevokeCtx, fmt.Sprintf("Error: %s", loggableResult.Err.Error()))
	}
	accessLog, err := newLog(loggableResult)
	if err != nil {
		log.Printf("Error: could not create access log. Cause: %s", err.Error())
		output.HasErrors = true
		output.Message = SomeErrorsDuringProcessMsg
		return
	}
	err = useCase.AccessPermissionStorage.SaveLog(accessLog)
	if err != nil {
		log.Printf("Error: could not save access log. Cause: %s", err.Error())
		output.HasErrors = true
		output.Message = SomeErrorsDuringProcessMsg
	}
}

func filterAccessibleInstances(selectedInstancesIDs, accessibleInstancesIDs []string, username string) []string {
	instancesToRevoke := make([]string, 0, len(selectedInstancesIDs))
	for _, instanceID := range selectedInstancesIDs {
		if !utils.Contains(accessibleInstancesIDs, instanceID) {
			log.Printf("WARN: Instance '%s' is not accessible by the user '%s'. Ignoring...", instanceID, username)
			continue
		}
		instancesToRevoke = append(instancesToRevoke, instanceID)
	}
	return instancesToRevoke
}

func newLog(loggableResult loggableRevokeResult) (*entity.AccessPermissionLog, error) {
	return entity.NewAccessPermissionLog(
		loggableResult.RevokeCtx.Instance.ID,
		loggableResult.RevokeCtx.User.ID.String(),
		"",
		loggableResult.LogMessagePt,
		loggableResult.RevokeCtx.OperationUserID,
		loggableResult.Err == nil,
	)
}

func logRevokeContextWithIndex(revokeCtx *revokeAccessContext, message string) {
	log.Printf("user {%s} | [%d/%d] instance {%s # %s}: %s",
		revokeCtx.User.Username, revokeCtx.InstanceIndex+1, revokeCtx.InstancesQty, revokeCtx.Instance.EcosystemName, revokeCtx.Instance.Name, message)
}

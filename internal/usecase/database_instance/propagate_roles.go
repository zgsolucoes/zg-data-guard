package instance

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
)

type PropagateRolesUseCase struct {
	DatabaseInstanceStorage storage.DatabaseInstanceStorage
	DatabaseRoleStorage     storage.DatabaseRoleStorage
}

func NewPropagateRolesUseCase(
	instanceStorage storage.DatabaseInstanceStorage,
	rolesStorage storage.DatabaseRoleStorage,
) *PropagateRolesUseCase {
	return &PropagateRolesUseCase{DatabaseInstanceStorage: instanceStorage, DatabaseRoleStorage: rolesStorage}
}

// Execute godoc
/** Responsible for creating all roles existing in zg-data-guard in all enabled database instances or in the selected database instances.
It creates the roles concurrently in all instances.
It returns a list of results for each database instance, indicating if the roles were created successfully or not.
*/
func (tc *PropagateRolesUseCase) Execute(input dto.PropagateRolesInputDTO, operationUserID string) ([]*dto.PropagateRolesOutputDTO, error) {
	if len(input.DatabaseInstancesIDs) > 0 {
		log.Printf("Propagating roles to the selected %d database instances. Requester: %s", len(input.DatabaseInstancesIDs), operationUserID)
		selectedInstances, err := tc.DatabaseInstanceStorage.FindAllDTOs("", "", input.DatabaseInstancesIDs)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		if len(selectedInstances) == 0 {
			return nil, ErrNoDatabaseInstancesFound
		}
		return tc.propagateRolesInInstances(selectedInstances)
	}

	log.Printf("Propagating roles to all enabled database instances. Requester: %s", operationUserID)
	enabledInstances, err := tc.DatabaseInstanceStorage.FindAllDTOsEnabled("", "")
	if err != nil {
		return nil, err
	}
	return tc.propagateRolesInInstances(enabledInstances)
}

func (tc *PropagateRolesUseCase) propagateRolesInInstances(dbInstances []*dto.DatabaseInstanceOutputDTO) ([]*dto.PropagateRolesOutputDTO, error) {
	roles, err := tc.DatabaseRoleStorage.FindAll()
	if err != nil {
		return nil, fmt.Errorf("error while fetching database roles. Cause: %w", err)
	}

	instancesQty := len(dbInstances)
	resultsChan := make(chan *dto.PropagateRolesOutputDTO, instancesQty)
	var wg sync.WaitGroup
	wg.Add(instancesQty)

	for idx, instance := range dbInstances {
		idx := idx
		go func(instance *dto.DatabaseInstanceOutputDTO) {
			defer wg.Done()
			output := tc.propagateRolesToInstance(instance, roles, idx, instancesQty)
			resultsChan <- output
		}(instance)
	}

	// Wait for all goroutines to finish and close the channels
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	return tc.buildResultingList(resultsChan), nil
}

func (tc *PropagateRolesUseCase) propagateRolesToInstance(instanceDto *dto.DatabaseInstanceOutputDTO, roles []*entity.DatabaseRole, idx int, instancesQty int) *dto.PropagateRolesOutputDTO {
	output := dto.PropagateRolesOutputDTO{
		DatabaseInstanceID: instanceDto.ID,
		Success:            false,
		Instance:           instanceDto.Name,
		Ecosystem:          instanceDto.EcosystemName,
		Technology:         instanceDto.DatabaseTechnologyName,
	}
	if !instanceDto.Enabled {
		output.Message = "database instance is disabled"
		return &output
	}

	c, err := connector.NewDatabaseConnector(instanceDto, "")
	if err != nil {
		output.Message = err.Error()
		return &output
	}
	log.Printf("Creating roles in instance [%s] %s (%d/%d)", instanceDto.EcosystemName, instanceDto.Name, idx+1, instancesQty)
	rolesForCreation := make([]*connector.DatabaseRole, 0, len(roles))
	for _, role := range roles {
		rolesForCreation = append(rolesForCreation, &connector.DatabaseRole{Name: role.Name})
	}
	err = c.CreateRoles(rolesForCreation)
	if err != nil {
		log.Printf("Error creating roles in instance [%s] %s. Cause: %v", instanceDto.EcosystemName, instanceDto.Name, err)
		output.Message = err.Error()
		return &output
	}

	output.Message = "roles created in database instance successfully!"
	if err = tc.updateInstanceProperties(instanceDto.ID); err != nil {
		output.Message = fmt.Sprintf("%s However there was an %v", output.Message, err.Error())
		return &output
	}

	output.Success = true
	return &output
}

func (tc *PropagateRolesUseCase) updateInstanceProperties(instanceID string) error {
	instance, err := tc.DatabaseInstanceStorage.FindByID(instanceID)
	if err != nil {
		return fmt.Errorf("error while fetching the database instance %s. Cause: %w", instanceID, err)
	}
	instance.CreateRoles()
	if err := tc.DatabaseInstanceStorage.Update(instance); err != nil {
		return fmt.Errorf("error to update created roles in database instance %s. Cause: %w", instanceID, err)
	}
	return nil
}

func (tc *PropagateRolesUseCase) buildResultingList(resultsChan chan *dto.PropagateRolesOutputDTO) []*dto.PropagateRolesOutputDTO {
	var connectionOutputs []*dto.PropagateRolesOutputDTO
	for result := range resultsChan {
		connectionOutputs = append(connectionOutputs, result)
	}

	return connectionOutputs
}

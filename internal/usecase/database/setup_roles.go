package database

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
	"github.com/zgsolucoes/zg-data-guard/pkg/utils"
)

var ErrNoDatabasesFound = fmt.Errorf("no databases found with the provided IDs")

type SetupRolesInDatabasesUseCase struct {
	DatabaseInstanceStorage storage.DatabaseInstanceStorage
	DatabaseStorage         storage.DatabaseStorage
}

func NewSetupRolesInDatabasesUseCase(
	instanceStorage storage.DatabaseInstanceStorage,
	databaseStorage storage.DatabaseStorage,
) *SetupRolesInDatabasesUseCase {
	return &SetupRolesInDatabasesUseCase{DatabaseInstanceStorage: instanceStorage, DatabaseStorage: databaseStorage}
}

// Execute godoc
/** Responsible for applying grants to roles in all enabled databases or in the selected databases and/or database instance.
It groups the databases by instance and applies the grants to roles in all databases of each instance concurrently.
I.e. if there are 3 instances with 20 databases each, it will apply the grants to roles in all 60 databases concurrently.
It returns a list of results for each database, indicating if the grants were applied successfully or not.
The grant script for PostgreSQL is read from a file: internal/database/connector/scripts/postgres/setup_grants_roles_database.sql */
func (uc *SetupRolesInDatabasesUseCase) Execute(input dto.SetupRolesInputDTO, operationUserID string) ([]*dto.SetupRolesOutputDTO, error) {
	if len(input.DatabasesIDs) > 0 || input.DatabaseInstanceID != "" {
		log.Printf("Applying grants to roles in the selected %d databases and/or database instance %s. Requester: %s", len(input.DatabasesIDs), input.DatabaseInstanceID, operationUserID)
		selectedDatabases, err := uc.DatabaseStorage.FindAll(input.DatabaseInstanceID, input.DatabasesIDs)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		if len(selectedDatabases) == 0 {
			return nil, ErrNoDatabasesFound
		}
		return uc.setupRoles(selectedDatabases)
	}

	log.Printf("Applying grants to roles in all enabled databases. Requester: %s", operationUserID)
	enabledDbs, err := uc.DatabaseStorage.FindAllEnabled("")
	if err != nil {
		return nil, err
	}
	return uc.setupRoles(enabledDbs)
}

func (uc *SetupRolesInDatabasesUseCase) setupRoles(databases []*entity.Database) ([]*dto.SetupRolesOutputDTO, error) {
	groupedByInstance := utils.GroupByProperty(databases, func(d *entity.Database) string {
		return d.DatabaseInstanceID
	})
	resultsChan := make(chan *dto.SetupRolesOutputDTO, len(databases))

	var wg sync.WaitGroup
	instancesQty := len(groupedByInstance)
	wg.Add(instancesQty)

	index := 0
	for instanceID, instanceDatabases := range groupedByInstance {
		go uc.executeSetupRolesForInstance(instanceID, instanceDatabases, &wg, resultsChan, index, instancesQty)
		index++
	}

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	var results []*dto.SetupRolesOutputDTO
	for result := range resultsChan {
		results = append(results, result)
	}
	return results, nil
}

func (uc *SetupRolesInDatabasesUseCase) executeSetupRolesForInstance(
	instanceID string,
	databases []*entity.Database,
	wg *sync.WaitGroup,
	resultsChan chan *dto.SetupRolesOutputDTO,
	index, instancesQty int) {
	defer wg.Done()

	instanceDto, err := uc.DatabaseInstanceStorage.FindDTOByID(instanceID)
	if err != nil {
		uc.markErrorInAllDatabases(instanceID, "", databases, resultsChan, err.Error())
		return
	}
	log.Printf("[%d/%d] | [%s # %s]: Applying grants to roles in all databases of instance", index+1, instancesQty, instanceDto.EcosystemName, instanceDto.Name)

	if !instanceDto.Enabled {
		log.Printf("[%d/%d] | [%s # %s]: Instance is disabled so no grants will be applied to roles in its databases.", index+1, instancesQty, instanceDto.EcosystemName, instanceDto.Name)
		uc.markErrorInAllDatabases(instanceID, instanceDto.Name, databases, resultsChan, "this database belongs to a disabled instance.")
		return
	}

	databasesQty := len(databases)
	for dbIndex, db := range databases {
		wg.Add(1)
		go func(db *entity.Database, dbIndex, databasesQty int) {
			defer wg.Done()
			result := uc.setupRolesForDatabase(instanceDto, db, dbIndex, databasesQty)
			resultsChan <- result
		}(db, dbIndex, databasesQty)
	}
}

func (uc *SetupRolesInDatabasesUseCase) markErrorInAllDatabases(
	databaseInstanceID,
	instanceName string,
	databasesOfInstance []*entity.Database,
	resultsChan chan *dto.SetupRolesOutputDTO,
	errMsg string) {
	for _, database := range databasesOfInstance {
		resultsChan <- &dto.SetupRolesOutputDTO{
			DatabaseID:         database.ID.String(),
			DatabaseName:       database.Name,
			DatabaseInstanceID: databaseInstanceID,
			Instance:           instanceName,
			Success:            false,
			Message:            errMsg,
		}
	}
}

func (uc *SetupRolesInDatabasesUseCase) setupRolesForDatabase(instanceDto *dto.DatabaseInstanceOutputDTO, database *entity.Database, dbIndex, databaseQty int) *dto.SetupRolesOutputDTO {
	output := dto.SetupRolesOutputDTO{
		DatabaseID:         database.ID.String(),
		DatabaseName:       database.Name,
		DatabaseInstanceID: instanceDto.ID,
		Instance:           instanceDto.Name,
		Ecosystem:          instanceDto.EcosystemName,
		Technology:         fmt.Sprintf("%s %s", instanceDto.DatabaseTechnologyName, instanceDto.DatabaseTechnologyVersion),
		Success:            false,
	}
	if !database.Enabled {
		output.Message = "this database is disabled."
		return &output
	}

	log.Printf("[%d/%d] | [%s # %s # %s]: Applying grants to roles in database", dbIndex+1, databaseQty, instanceDto.EcosystemName, instanceDto.Name, database.Name)

	c, err := connector.NewDatabaseConnector(instanceDto, database.Name)
	if err != nil {
		output.Message = err.Error()
		return &output
	}
	err = c.SetupGrantsToRoles()
	if err != nil {
		output.Message = err.Error()
		return &output
	}
	err = uc.updateDatabaseRolesConfigured(database)
	if err != nil {
		output.Message = err.Error()
		return &output
	}
	log.Printf("[%d/%d] | [%s # %s # %s]: Grants applied successfully to roles in database", dbIndex+1, databaseQty, instanceDto.EcosystemName, instanceDto.Name, database.Name)
	output.Success = true
	output.Message = "grants applied successfully to roles in database!"
	return &output
}

func (uc *SetupRolesInDatabasesUseCase) updateDatabaseRolesConfigured(database *entity.Database) error {
	database.ConfigureRoles()
	return uc.DatabaseStorage.Update(database)
}

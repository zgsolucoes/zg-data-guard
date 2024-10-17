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
)

var ErrNoDatabaseInstancesFound = fmt.Errorf("no database instances found with the provided IDs")

type SyncDatabasesUseCase struct {
	DatabaseInstanceStorage storage.DatabaseInstanceStorage
	DatabaseStorage         storage.DatabaseStorage
}

func NewSyncDatabasesUseCase(dbInstanceStorage storage.DatabaseInstanceStorage, databaseStorage storage.DatabaseStorage) *SyncDatabasesUseCase {
	return &SyncDatabasesUseCase{
		DatabaseInstanceStorage: dbInstanceStorage,
		DatabaseStorage:         databaseStorage,
	}
}

// Execute godoc
/** Responsible for synchronizing the databases of all enabled database instances or of the selected database instances.
It groups the databases by instance and synchronizes the databases of each instance concurrently.
The synchronization process consists of comparing the databases of the instance with the databases of the zg-data-guard:
- If a database exists in the instance and not in the zg-data-guard, it is created.
- If a database exists in the instance and in the zg-data-guard, but with different sizes, the size is updated.
- If a database exists in the zg-data-guard and not in the instance, it is disabled.
It returns a list of results for each database instance, indicating if the databases were synchronized successfully or not.
*/
func (uc *SyncDatabasesUseCase) Execute(input dto.SyncDatabasesInputDTO, operationUserID string) ([]*dto.SyncDatabasesOutputDTO, error) {
	if len(input.DatabaseInstancesIDs) > 0 {
		log.Printf("Synchronizing databases of the selected %d database instances by user %s", len(input.DatabaseInstancesIDs), operationUserID)
		selectedInstances, err := uc.DatabaseInstanceStorage.FindAllDTOs("", "", input.DatabaseInstancesIDs)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		if len(selectedInstances) == 0 {
			return nil, ErrNoDatabaseInstancesFound
		}
		return uc.syncDatabasesFromInstances(operationUserID, selectedInstances)
	}

	log.Printf("Synchronizing databases of all enabled database instances by user %s", operationUserID)
	enabledInstances, err := uc.DatabaseInstanceStorage.FindAllDTOsEnabled("", "")
	if err != nil {
		return nil, err
	}
	return uc.syncDatabasesFromInstances(operationUserID, enabledInstances)
}

func (uc *SyncDatabasesUseCase) syncDatabasesFromInstances(operationUserID string, dbInstances []*dto.DatabaseInstanceOutputDTO) ([]*dto.SyncDatabasesOutputDTO, error) {
	instancesQty := len(dbInstances)
	resultsChan := make(chan *dto.SyncDatabasesOutputDTO, instancesQty)

	var wg sync.WaitGroup
	wg.Add(instancesQty)

	// Start a goroutine for each instance to sync the databases with it and send the results to the channels to be processed later
	for idx, instance := range dbInstances {
		go func(idx int, instance *dto.DatabaseInstanceOutputDTO) {
			defer wg.Done()
			output := uc.syncDatabases(instance, operationUserID, idx, instancesQty)
			resultsChan <- output
		}(idx, instance)
	}

	// Wait for all goroutines to finish and close the channels
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	return uc.buildResultingList(resultsChan), nil
}

func (uc *SyncDatabasesUseCase) syncDatabases(instanceDto *dto.DatabaseInstanceOutputDTO, operationUserID string, idx int, instancesQty int) *dto.SyncDatabasesOutputDTO {
	output := dto.SyncDatabasesOutputDTO{
		DatabaseInstanceID: instanceDto.ID,
		Success:            false,
		Instance:           instanceDto.Name,
		Ecosystem:          instanceDto.EcosystemName,
		Technology:         fmt.Sprintf("%s %s", instanceDto.DatabaseTechnologyName, instanceDto.DatabaseTechnologyVersion),
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
	log.Printf("[%d/%d] | [%s # %s]: Synchronizing databases", idx+1, instancesQty, instanceDto.EcosystemName, instanceDto.Name)
	instanceDbs, err := c.ListDatabases()
	if err != nil {
		output.Message = err.Error()
		return &output
	}
	dataguardDbs, err := uc.DatabaseStorage.FindAll(instanceDto.ID, []string{})
	if err != nil {
		output.Message = err.Error()
		return &output
	}
	err = uc.matchDatabases(instanceDto.ID, operationUserID, instanceDbs, dataguardDbs)
	if err != nil {
		output.Message = err.Error()
		return &output
	}
	err = uc.updateInstanceLastDatabaseSync(instanceDto.ID)
	if err != nil {
		output.Message = err.Error()
		return &output
	}
	totalDatabases := len(instanceDbs)
	log.Printf("[%d/%d] | [%s # %s]: %d databases synchronized successfully by user %s", idx+1, instancesQty, instanceDto.EcosystemName, instanceDto.Name, totalDatabases, operationUserID)
	output.Success = true
	output.TotalDatabases = totalDatabases
	output.Message = fmt.Sprintf("%d databases synchronized successfully!", totalDatabases)
	return &output
}

func (uc *SyncDatabasesUseCase) matchDatabases(dbInstanceID, operationUserID string, instanceDatabases []*connector.Database, dataguardDatabases []*entity.Database) error {
	dataguardDBMap := make(map[string]*entity.Database)
	for _, db := range dataguardDatabases {
		dataguardDBMap[db.Name] = db
	}

	for _, databaseFromInstance := range instanceDatabases {
		err := uc.createOrUpdateDatabase(dbInstanceID, operationUserID, dataguardDBMap, databaseFromInstance)
		if err != nil {
			return err
		}

		// Remove the processed database from the map
		delete(dataguardDBMap, databaseFromInstance.Name)
	}

	err := uc.disableNonExistentDbs(dataguardDBMap)
	if err != nil {
		return err
	}

	return nil
}

func (uc *SyncDatabasesUseCase) createOrUpdateDatabase(dbInstanceID string, operationUserID string, dataguardDBMap map[string]*entity.Database, instanceDB *connector.Database) error {
	dataguardDB, exists := dataguardDBMap[instanceDB.Name]

	if !exists {
		err := uc.createNewDBAndSave(dbInstanceID, operationUserID, instanceDB)
		if err != nil {
			return err
		}
	} else if dataguardDB.CurrentSize != instanceDB.CurrentSize || !dataguardDB.Enabled {
		if !dataguardDB.Enabled {
			dataguardDB.Enable()
			log.Printf("Database '%s' re-enabled in database instance %s", instanceDB.Name, dbInstanceID)
		}
		dataguardDB.Update(instanceDB.CurrentSize)
		err := uc.DatabaseStorage.Update(dataguardDB)
		if err != nil {
			return err
		}
		log.Printf("Database '%s' updated successfully in database instance %s", instanceDB.Name, dbInstanceID)
	}
	return nil
}

func (uc *SyncDatabasesUseCase) createNewDBAndSave(dbInstanceID string, operationUserID string, instanceDB *connector.Database) error {
	newDatabase, err := entity.NewDatabase(instanceDB.Name, "", dbInstanceID, instanceDB.CurrentSize, operationUserID)
	if err != nil {
		return err
	}
	err = uc.DatabaseStorage.Save(newDatabase)
	if err != nil {
		log.Printf("Error creating database %s in database instance %s. Cause: %v", instanceDB.Name, dbInstanceID, err)
		return err
	}
	log.Printf("Database '%s' with size %s created successfully in database instance %s", instanceDB.Name, instanceDB.CurrentSize, dbInstanceID)
	return nil
}

func (uc *SyncDatabasesUseCase) disableNonExistentDbs(dataguardDBMap map[string]*entity.Database) error {
	for _, dataguardDB := range dataguardDBMap {
		if !dataguardDB.Enabled {
			continue
		}
		dataguardDB.Disable()
		err := uc.DatabaseStorage.Update(dataguardDB)
		if err != nil {
			return err
		}
		log.Printf("Database %s no longer exists in the database instance %s and was disabled", dataguardDB.Name, dataguardDB.DatabaseInstanceID)
	}
	return nil
}

func (uc *SyncDatabasesUseCase) buildResultingList(resultsChan chan *dto.SyncDatabasesOutputDTO) []*dto.SyncDatabasesOutputDTO {
	var connectionOutputs []*dto.SyncDatabasesOutputDTO
	for result := range resultsChan {
		connectionOutputs = append(connectionOutputs, result)
	}

	return connectionOutputs
}

func (uc *SyncDatabasesUseCase) updateInstanceLastDatabaseSync(instanceID string) error {
	instance, err := uc.DatabaseInstanceStorage.FindByID(instanceID)
	if err != nil {
		return fmt.Errorf("databases synchronized but the is an error fetching the database instance %s. Cause: %w", instanceID, err)
	}
	instance.RefreshLastDatabaseSync()
	if err := uc.DatabaseInstanceStorage.Update(instance); err != nil {
		return fmt.Errorf("databases synchronized but the is an error updating last sync date for database instance %s. Cause: %w", instanceID, err)
	}
	return nil
}

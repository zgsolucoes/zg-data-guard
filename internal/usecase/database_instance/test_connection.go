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
)

var ErrNoDatabaseInstancesFound = fmt.Errorf("no database instances found with the provided IDs")

type TestConnectionUseCase struct {
	DatabaseInstanceStorage storage.DatabaseInstanceStorage
}

func NewTestConnectionUseCase(repo storage.DatabaseInstanceStorage) *TestConnectionUseCase {
	return &TestConnectionUseCase{DatabaseInstanceStorage: repo}
}

func (tc *TestConnectionUseCase) Execute(input dto.TestConnectionInputDTO) ([]*dto.TestConnectionOutputDTO, error) {
	if len(input.DatabaseInstancesIDs) > 0 {
		log.Printf("Testing connections with the selected %d database instances", len(input.DatabaseInstancesIDs))
		selectedInstances, err := tc.DatabaseInstanceStorage.FindAllDTOs("", "", input.DatabaseInstancesIDs)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		if len(selectedInstances) == 0 {
			return nil, ErrNoDatabaseInstancesFound
		}
		return tc.testConnectionsFromInstances(selectedInstances)
	}

	log.Printf("Testing connections with all enabled database instances")
	enabledInstances, err := tc.DatabaseInstanceStorage.FindAllDTOsEnabled("", "")
	if err != nil {
		return nil, err
	}
	return tc.testConnectionsFromInstances(enabledInstances)
}

func (tc *TestConnectionUseCase) testConnectionsFromInstances(dbInstances []*dto.DatabaseInstanceOutputDTO) ([]*dto.TestConnectionOutputDTO, error) {
	instancesQty := len(dbInstances)
	resultsChan := make(chan *dto.TestConnectionOutputDTO, instancesQty)
	var wg sync.WaitGroup
	wg.Add(instancesQty)

	// Start a goroutine for each instance to test the connection with it and send the results to the channels to be processed later
	for idx, instance := range dbInstances {
		idx := idx
		go func(instance *dto.DatabaseInstanceOutputDTO) {
			defer wg.Done()
			output := tc.testInstanceConnection(instance, idx, instancesQty)
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

func (tc *TestConnectionUseCase) testInstanceConnection(instanceDto *dto.DatabaseInstanceOutputDTO, idx int, instancesQty int) *dto.TestConnectionOutputDTO {
	output := dto.TestConnectionOutputDTO{
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
	log.Printf("Testing connection with [%s] %s (%d/%d)", instanceDto.EcosystemName, instanceDto.Name, idx+1, instancesQty)
	err = c.TestConnection()
	if err != nil {
		output.Message = err.Error()
		if err2 := tc.updateInstanceLastConnectionTest(instanceDto.ID, output.Message, false); err2 != nil {
			output.Message = fmt.Sprintf("%s. There was also an %v", err.Error(), err2.Error())
			return &output
		}
		return &output
	}

	output.Message = "connection established successfully!"
	if err = tc.updateInstanceLastConnectionTest(instanceDto.ID, output.Message, true); err != nil {
		output.Message = fmt.Sprintf("%s However there was an %v", output.Message, err.Error())
		return &output
	}

	output.Success = true
	return &output
}

func (tc *TestConnectionUseCase) updateInstanceLastConnectionTest(instanceID, resultMsg string, success bool) error {
	instance, err := tc.DatabaseInstanceStorage.FindByID(instanceID)
	if err != nil {
		return fmt.Errorf("error while fetching the database instance %s. Cause: %w", instanceID, err)
	}
	instance.RefreshLastConnectionTest(success, resultMsg)
	if err := tc.DatabaseInstanceStorage.Update(instance); err != nil {
		return fmt.Errorf("error to update last connection date for database instance %s. Cause: %w", instanceID, err)
	}
	return nil
}

func (tc *TestConnectionUseCase) buildResultingList(resultsChan chan *dto.TestConnectionOutputDTO) []*dto.TestConnectionOutputDTO {
	var connectionOutputs []*dto.TestConnectionOutputDTO
	for result := range resultsChan {
		connectionOutputs = append(connectionOutputs, result)
	}

	return connectionOutputs
}

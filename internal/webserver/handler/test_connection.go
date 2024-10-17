package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	dbUsecase "github.com/zgsolucoes/zg-data-guard/internal/usecase/database_instance"
)

const (
	opTestConnection = "test-connection"
)

var testConnectionUC *dbUsecase.TestConnectionUseCase

// TestConnectionHandler godoc
// @BasePath /api/v1
// @Summary Test connection with the selected database instances, if instances ids are not provided, test connection with all enabled instances
// @Description Test connection with the selected database instances, if instances ids are not provided, test connection with all enabled instances
// @Tags Database Instance
// @Accept json
// @Produce json
// @Param request body dto.TestConnectionInputDTO false "Request body"
// @Success 200 {object} TestConnectionResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /database-instance/test-connection [post]
// @Security ApiKeyAuth
func TestConnectionHandler(w http.ResponseWriter, r *http.Request) {
	var input dto.TestConnectionInputDTO
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		log.Printf("error decoding request body: %v", err.Error())
		sendError(w, http.StatusBadRequest, "error decoding request body")
		return
	}

	if input.DatabaseInstancesIDs != nil {
		for _, id := range input.DatabaseInstancesIDs {
			if !validateUUID(w, id, "databaseInstancesIds list contains a value that") {
				return
			}
		}
	}

	connectionOutputs, err := testConnectionUC.Execute(input)
	if err != nil && errors.Is(err, dbUsecase.ErrNoDatabaseInstancesFound) {
		sendError(w, http.StatusNotFound, buildErrorMessage(opTestConnection, err))
		return
	}
	if err != nil {
		log.Printf("error in operation %s: %v", opTestConnection, err)
		sendError(w, http.StatusInternalServerError, buildErrorMessage(opTestConnection, err))
		return
	}

	sendSuccessList(w, opTestConnection, connectionOutputs, len(connectionOutputs), 0, 0)
}

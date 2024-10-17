package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	dbUsecase "github.com/zgsolucoes/zg-data-guard/internal/usecase/database"
)

const (
	opSyncDatabases = "sync-databases"
)

var syncDatabasesUC *dbUsecase.SyncDatabasesUseCase

// SyncDatabasesHandler godoc
// @BasePath /api/v1
// @Summary Sync databases from selected database instances, if instances ids are not provided, sync all enabled instances
// @Description Sync databases from selected database instances, if instances ids are not provided, sync all enabled instances
// @Tags Database Instance
// @Accept json
// @Produce json
// @Param request body dto.SyncDatabasesInputDTO false "Request body"
// @Success 200 {object} SyncDatabasesResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /database-instance/sync-databases [post]
// @Security ApiKeyAuth
func SyncDatabasesHandler(w http.ResponseWriter, r *http.Request) {
	userID, hasError := getUserIDFromAuthenticatedRequest(w, r)
	if hasError || userID == emptyString {
		return
	}

	var input dto.SyncDatabasesInputDTO
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

	syncOutputs, err := syncDatabasesUC.Execute(input, userID)
	if err != nil && errors.Is(err, dbUsecase.ErrNoDatabaseInstancesFound) {
		sendError(w, http.StatusNotFound, buildErrorMessage(opSyncDatabases, err))
		return
	}
	if err != nil {
		log.Printf("error in operation %s: %v", opSyncDatabases, err)
		sendError(w, http.StatusInternalServerError, buildErrorMessage(opSyncDatabases, err))
		return
	}

	sendSuccessList(w, opSyncDatabases, syncOutputs, len(syncOutputs), 0, 0)
}

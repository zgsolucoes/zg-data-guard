package handler

import (
	"net/http"

	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	dbUsecase "github.com/zgsolucoes/zg-data-guard/internal/usecase/database"
)

const (
	opListDatabases         = "list-databases"
	paramDatabaseInstanceID = "databaseInstanceId"
)

var listDatabasesUC *dbUsecase.ListDatabasesUseCase

// ListDatabasesHandler godoc
// @BasePath /api/v1
// @Summary List all existing databases
// @Description List all existing databases
// @Tags Database
// @Accept json
// @Produce json
// @Param ecosystemId query string false "Ecosystem ID"
// @Param databaseInstanceId query string false "Database Instance ID"
// @Success 200 {object} ListDatabasesResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /databases [get]
// @Security ApiKeyAuth
func ListDatabasesHandler(w http.ResponseWriter, r *http.Request) {
	ecosystemID := r.URL.Query().Get(paramEcosystemID)
	databaseInstanceID := r.URL.Query().Get(paramDatabaseInstanceID)
	if (ecosystemID != "" && !validateUUID(w, ecosystemID, paramEcosystemID)) || (databaseInstanceID != "" && !validateUUID(w, databaseInstanceID, paramDatabaseInstanceID)) {
		return
	}

	outputDTOs, err := listDatabasesUC.Execute(ecosystemID, databaseInstanceID)
	if err != nil {
		sendError(w, http.StatusInternalServerError, buildErrorMessage(opListDatabases, err))
		return
	}
	if outputDTOs == nil {
		outputDTOs = make([]*dto.DatabaseOutputDTO, 0)
	}

	sendSuccessList(w, opListDatabases, outputDTOs, len(outputDTOs), 0, 0)
}

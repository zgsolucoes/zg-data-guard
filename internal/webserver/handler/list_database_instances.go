package handler

import (
	"net/http"

	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	dbInstanceUsecase "github.com/zgsolucoes/zg-data-guard/internal/usecase/database_instance"
)

const (
	opListDatabaseInstances = "list-database-instances"
	paramEcosystemID        = "ecosystemId"
	paramTechnologyID       = "technologyId"
)

var listDBInstancesUC *dbInstanceUsecase.ListDatabaseInstancesUseCase

// ListDatabaseInstancesHandler godoc
// @BasePath /api/v1
// @Summary List all existing database instances
// @Description List all existing database instances
// @Tags Database Instance
// @Accept json
// @Produce json
// @Param ecosystemId query string false "Ecosystem ID"
// @Param technologyId query string false "Database Technology ID"
// @Param onlyEnabled query bool false "Only Enabled Database Instances"
// @Success 200 {object} ListDatabaseInstancesResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /database-instances [get]
// @Security ApiKeyAuth
func ListDatabaseInstancesHandler(w http.ResponseWriter, r *http.Request) {
	ecosystemID := r.URL.Query().Get(paramEcosystemID)
	technologyID := r.URL.Query().Get(paramTechnologyID)
	onlyEnabledParam := r.URL.Query().Get(paramOnlyEnabled)

	if (ecosystemID != "" && !validateUUID(w, ecosystemID, paramEcosystemID)) ||
		(technologyID != "" && !validateUUID(w, technologyID, paramTechnologyID) ||
			!validateBoolQueryParam(w, onlyEnabledParam, paramOnlyEnabled)) {
		return
	}

	onlyEnabledParamBool := getQueryParamBoolValue(onlyEnabledParam)
	outputDBInstances, err := listDBInstancesUC.Execute(ecosystemID, technologyID, onlyEnabledParamBool)
	if err != nil {
		sendError(w, http.StatusInternalServerError, buildErrorMessage(opListDatabaseInstances, err))
		return
	}
	if outputDBInstances == nil {
		outputDBInstances = make([]*dto.DatabaseInstanceOutputDTO, 0)
	}

	sendSuccessList(w, opListDatabaseInstances, outputDBInstances, len(outputDBInstances), 0, 0)
}

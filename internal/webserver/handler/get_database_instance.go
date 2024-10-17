package handler

import (
	"errors"
	"net/http"

	dbUsecase "github.com/zgsolucoes/zg-data-guard/internal/usecase/database_instance"
)

const opGetDatabaseInstance = "get-database-instance"

var getDBInstanceUC *dbUsecase.GetDatabaseInstanceUseCase

// GetDatabaseInstanceHandler godoc
// @BasePath /api/v1
// @Summary Get a database instance
// @Description Get an existing database instance
// @Tags Database Instance
// @Accept json
// @Produce json
// @Param id query string true "Database Instance ID"
// @Success 200 {object} GetDatabaseInstanceResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /database-instance [get]
// @Security ApiKeyAuth
func GetDatabaseInstanceHandler(w http.ResponseWriter, r *http.Request) {
	id, hasError := getIDFromQueryParamsAndValidate(w, r)
	if hasError {
		return
	}

	outputDatabaseInstance, err := getDBInstanceUC.Execute(id)
	if err != nil && errors.Is(err, dbUsecase.ErrDatabaseInstanceNotFound) {
		sendError(w, http.StatusNotFound, buildErrorMessage(opGetDatabaseInstance, err))
		return

	}
	if err != nil {
		sendError(w, http.StatusInternalServerError, buildErrorMessage(opGetDatabaseInstance, err))
		return
	}

	sendSuccess(w, opGetDatabaseInstance, outputDatabaseInstance)
}

package handler

import (
	"errors"
	"net/http"

	dbUsecase "github.com/zgsolucoes/zg-data-guard/internal/usecase/database_instance"
)

const opGetDatabaseInstanceCred = "get-database-instance-credentials"

// GetDatabaseInstanceCredentialsHandler godoc
// @BasePath /api/v1
// @Summary Get credentials of a specific database instance
// @Description Get credentials of an existing database instance
// @Tags Database Instance
// @Accept json
// @Produce json
// @Param id query string true "Database Instance ID"
// @Success 200 {object} GetDatabaseInstanceCredentialsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /database-instance/credentials [get]
// @Security ApiKeyAuth
func GetDatabaseInstanceCredentialsHandler(w http.ResponseWriter, r *http.Request) {
	id, hasError := getIDFromQueryParamsAndValidate(w, r)
	if hasError {
		return
	}
	userID, hasError := getUserIDFromAuthenticatedRequest(w, r)
	if hasError || userID == emptyString {
		return
	}

	outputDatabaseInstance, err := getDBInstanceUC.FetchCredentials(id, userID)
	if err != nil && errors.Is(err, dbUsecase.ErrDatabaseInstanceNotFound) {
		sendError(w, http.StatusNotFound, buildErrorMessage(opGetDatabaseInstanceCred, err))
		return

	}
	if err != nil {
		sendError(w, http.StatusInternalServerError, buildErrorMessage(opGetDatabaseInstanceCred, err))
		return
	}

	sendSuccess(w, opGetDatabaseInstanceCred, outputDatabaseInstance)
}

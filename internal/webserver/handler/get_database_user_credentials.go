package handler

import (
	"errors"
	"net/http"

	"github.com/zgsolucoes/zg-data-guard/internal/usecase/common"
)

const opGetDatabaseUserCred = "get-database-user-credentials"

// GetDatabaseUserCredentialsHandler godoc
// @BasePath /api/v1
// @Summary Get credentials of a specific database user
// @Description Get credentials of an existing database user
// @Tags Database User
// @Accept json
// @Produce json
// @Param id query string true "Database User ID"
// @Success 200 {object} GetDatabaseUserCredentialsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /database-user/credentials [get]
// @Security ApiKeyAuth
func GetDatabaseUserCredentialsHandler(w http.ResponseWriter, r *http.Request) {
	id, hasError := getIDFromQueryParamsAndValidate(w, r)
	if hasError {
		return
	}
	userID, hasError := getUserIDFromAuthenticatedRequest(w, r)
	if hasError || userID == emptyString {
		return
	}

	outputDatabaseUser, err := getDBUserUC.FetchCredentials(id, userID)
	if err != nil && errors.Is(err, common.ErrDatabaseUserNotFound) {
		sendError(w, http.StatusNotFound, buildErrorMessage(opGetDatabaseUserCred, err))
		return

	}
	if err != nil {
		sendError(w, http.StatusInternalServerError, buildErrorMessage(opGetDatabaseUserCred, err))
		return
	}

	sendSuccess(w, opGetDatabaseUserCred, outputDatabaseUser)
}

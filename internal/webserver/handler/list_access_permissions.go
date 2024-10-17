package handler

import (
	"net/http"

	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	usecase "github.com/zgsolucoes/zg-data-guard/internal/usecase/access_permission"
)

const (
	opListAccessPermissions = "list-access-permissions"
	paramDatabaseID         = "databaseId"
	paramDatabaseUserID     = "databaseUserId"
)

var listAccessPermissionsUC *usecase.ListAccessPermissionsUseCase

// ListAccessPermissionsHandler godoc
// @BasePath /api/v1
// @Summary List all existing access permissions established between users and databases
// @Description List all existing access permissions established between users and databases
// @Tags Access Permission
// @Accept json
// @Produce json
// @Param databaseId query string false "Database ID"
// @Param databaseUserId query string false "Database User ID"
// @Param databaseInstanceId query string false "Database Instance ID"
// @Success 200 {object} ListAccessPermissionsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /access-permissions [get]
// @Security ApiKeyAuth
func ListAccessPermissionsHandler(w http.ResponseWriter, r *http.Request) {
	databaseID := r.URL.Query().Get(paramDatabaseID)
	databaseUserID := r.URL.Query().Get(paramDatabaseUserID)
	databaseInstanceID := r.URL.Query().Get(paramDatabaseInstanceID)
	if validateParams(w, databaseUserID, databaseID, databaseInstanceID) {
		return
	}

	outputDTOs, err := listAccessPermissionsUC.Execute(databaseID, databaseUserID, databaseInstanceID)
	if err != nil {
		sendError(w, http.StatusInternalServerError, buildErrorMessage(opListAccessPermissions, err))
		return
	}
	if outputDTOs == nil {
		outputDTOs = make([]*dto.AccessPermissionOutputDTO, 0)
	}

	sendSuccessList(w, opListAccessPermissions, outputDTOs, len(outputDTOs), 0, 0)
}

func validateParams(w http.ResponseWriter, databaseUserID, databaseID, databaseInstanceID string) bool {
	return validateUUIDParam(w, databaseUserID, paramDatabaseUserID) ||
		validateUUIDParam(w, databaseID, paramDatabaseID) ||
		validateUUIDParam(w, databaseInstanceID, paramDatabaseInstanceID)
}

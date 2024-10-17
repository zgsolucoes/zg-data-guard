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
	opSetupRoles = "setup-roles"
)

var setupRolesUC *dbUsecase.SetupRolesInDatabasesUseCase

// SetupRolesHandler godoc
// @BasePath /api/v1
// @Summary Setup roles (applying grants) in the selected databases, if databases ids are not provided, setup roles in all enabled databases belonging to the enabled instances
// @Description Setup roles (applying grants) in the selected databases, if databases ids are not provided, setup roles in all enabled databases belonging to the enabled instances
// @Tags Database
// @Accept json
// @Produce json
// @Param request body dto.SetupRolesInputDTO false "Request body"
// @Success 200 {object} SetupRolesResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /database/setup-roles [post]
// @Security ApiKeyAuth
func SetupRolesHandler(w http.ResponseWriter, r *http.Request) {
	userID, hasError := getUserIDFromAuthenticatedRequest(w, r)
	if hasError || userID == emptyString {
		return
	}

	var input dto.SetupRolesInputDTO
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		log.Printf("error decoding request body: %v", err.Error())
		sendError(w, http.StatusBadRequest, "error decoding request body")
		return
	}

	if input.DatabasesIDs != nil {
		for _, id := range input.DatabasesIDs {
			if !validateUUID(w, id, "databaseIds list contains a value that") {
				return
			}
		}
	}
	if input.DatabaseInstanceID != emptyString && !validateUUID(w, input.DatabaseInstanceID, "databaseInstanceId") {
		return
	}

	outputs, err := setupRolesUC.Execute(input, userID)
	if err != nil && errors.Is(err, dbUsecase.ErrNoDatabasesFound) {
		sendError(w, http.StatusNotFound, buildErrorMessage(opSetupRoles, err))
		return
	}
	if err != nil {
		log.Printf("error in operation %s: %v", opSetupRoles, err)
		sendError(w, http.StatusInternalServerError, buildErrorMessage(opSetupRoles, err))
		return
	}

	sendSuccessList(w, opSetupRoles, outputs, len(outputs), 0, 0)
}

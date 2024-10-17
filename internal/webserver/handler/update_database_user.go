package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	"github.com/zgsolucoes/zg-data-guard/internal/usecase/common"
	dbUsecase "github.com/zgsolucoes/zg-data-guard/internal/usecase/database_user"
)

const opUpdateDatabaseUser = "update-database-user"

var updateDBUserUC *dbUsecase.UpdateDatabaseUserUseCase

// UpdateDatabaseUserHandler godoc
// @BasePath /api/v1
// @Summary Update a database user
// @Description Update an existing database user
// @Tags Database User
// @Accept json
// @Produce json
// @Param id query string true "Database User ID"
// @Param request body dto.UpdateDatabaseUserInputDTO true "Request body"
// @Success 200 {object} CreateDatabaseUserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /database-user [put]
// @Security ApiKeyAuth
func UpdateDatabaseUserHandler(w http.ResponseWriter, r *http.Request) {
	id, hasError := getIDFromQueryParamsAndValidate(w, r)
	if hasError {
		return
	}

	userID, hasError := getUserIDFromAuthenticatedRequest(w, r)
	if hasError || userID == emptyString {
		return
	}

	var input dto.UpdateDatabaseUserInputDTO
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		log.Printf("error decoding request body: %v", err.Error())
		sendError(w, http.StatusBadRequest, "error decoding request body")
		return
	}

	if err = input.Validate(); err != nil {
		log.Printf("validation error: %v", err.Error())
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	output, err := updateDBUserUC.Execute(input, id, userID)
	if err != nil && errors.Is(err, common.ErrDatabaseUserNotFound) {
		sendError(w, http.StatusNotFound, buildErrorMessage(opUpdateDatabaseUser, err))
		return
	}
	if err != nil && errors.Is(err, dbUsecase.ErrDatabaseUserHasAccessPermissions) {
		sendError(w, http.StatusConflict, buildErrorMessage(opUpdateDatabaseUser, err))
		return
	}

	if err != nil && (errors.Is(err, dbUsecase.ErrDatabaseRoleNotFound)) {
		sendError(w, http.StatusNotFound, buildErrorMessage(opUpdateDatabaseUser, err))
		return
	}
	if err != nil {
		sendError(w, http.StatusInternalServerError, buildErrorMessage(opUpdateDatabaseUser, err))
		return
	}

	sendSuccess(w, opUpdateDatabaseUser, output)
}

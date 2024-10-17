package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	"github.com/zgsolucoes/zg-data-guard/internal/usecase/common"
	dbUserUsecase "github.com/zgsolucoes/zg-data-guard/internal/usecase/database_user"
)

var changeStatusDBUserUC *dbUserUsecase.ChangeStatusDatabaseUserUseCase

const opChangeStatus = "change-status-database-user"

// ChangeStatusDatabaseUserHandler godoc
// @BasePath /api/v1
// @Summary Change the status of a database user, enabling or disabling it. If disabled, it revokes access from all instances accessible by the user.
// @Description Change the status of a database user, enabling or disabling it. If disabled, it revokes access from all instances accessible by the user.
// @Tags Database User
// @Accept json
// @Produce json
// @Param request body dto.ChangeStatusInputDTO true "Request body"
// @Success 200 {object} ChangeStatusResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /database-user/change-status [patch]
// @Security ApiKeyAuth
func ChangeStatusDatabaseUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, hasError := getUserIDFromAuthenticatedRequest(w, r)
	if hasError || userID == emptyString {
		return
	}

	var input dto.ChangeStatusInputDTO
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

	output, err := changeStatusDBUserUC.Execute(input.ID, *input.Enabled, userID)
	if err != nil && errors.Is(err, common.ErrDatabaseUserNotFound) {
		sendError(w, http.StatusNotFound, buildErrorMessage(opChangeStatus, err))
		return
	}
	if err != nil && errors.Is(err, dbUserUsecase.ErrCouldNotRevokeAllAccess) {
		sendError(w, http.StatusConflict, buildErrorMessage(opChangeStatus, err))
		return
	}
	if err != nil {
		sendError(w, http.StatusInternalServerError, buildErrorMessage(opChangeStatus, err))
		return
	}
	sendSuccess(w, opChangeStatus, output)
}

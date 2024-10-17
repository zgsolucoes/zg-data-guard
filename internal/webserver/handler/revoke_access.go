package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	accessPermissionUsecase "github.com/zgsolucoes/zg-data-guard/internal/usecase/access_permission"
	"github.com/zgsolucoes/zg-data-guard/internal/usecase/common"
)

var revokeAccessPermissionUC *accessPermissionUsecase.RevokeAccessPermissionUseCase

const opRevokeAccess = "revoke-access"

// RevokeAccessHandler godoc
// @BasePath /api/v1
// @Summary Revoke connection access to a specific user from a set of instances and their respective databases
// @Description If no instance is provided, it revokes access from all instances accessible by the user.
// @Tags Access Permission
// @Accept json
// @Produce json
// @Param request body dto.RevokeAccessInputDTO true "Request body"
// @Success 200 {object} RevokeAccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /access-permission/revoke [post]
// @Security ApiKeyAuth
func RevokeAccessHandler(w http.ResponseWriter, r *http.Request) {
	userID, hasError := getUserIDFromAuthenticatedRequest(w, r)
	if hasError || userID == emptyString {
		return
	}

	var input dto.RevokeAccessInputDTO
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

	output, err := revokeAccessPermissionUC.Execute(input, userID)
	if err != nil && errors.Is(err, common.ErrDatabaseUserNotFound) {
		sendError(w, http.StatusNotFound, buildErrorMessage(opRevokeAccess, err))
		return
	}
	if err != nil && errors.Is(err, common.ErrNoAccessibleInstancesFound) {
		sendError(w, http.StatusBadRequest, buildErrorMessage(opRevokeAccess, err))
		return
	}
	if err != nil {
		sendError(w, http.StatusInternalServerError, buildErrorMessage(opRevokeAccess, err))
		return
	}

	sendSuccess(w, opRevokeAccess, output)
}

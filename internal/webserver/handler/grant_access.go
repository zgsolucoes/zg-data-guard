package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	accessPermissionUsecase "github.com/zgsolucoes/zg-data-guard/internal/usecase/access_permission"
)

var grantAccessPermissionUC *accessPermissionUsecase.GrantAccessPermissionUseCase

const opGrantAccess = "grant-access"

// GrantAccessHandler godoc
// @BasePath /api/v1
// @Summary Grant connection access to a set of users to a set of instances and their respective databases
// @Description Grant connection access to a set of users to a set of instances and their respective databases
// @Tags Access Permission
// @Accept json
// @Produce json
// @Param request body dto.GrantAccessInputDTO true "Request body"
// @Success 200 {object} GrantAccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /access-permission/grant [post]
// @Security ApiKeyAuth
func GrantAccessHandler(w http.ResponseWriter, r *http.Request) {
	userID, hasError := getUserIDFromAuthenticatedRequest(w, r)
	if hasError || userID == emptyString {
		return
	}

	var input dto.GrantAccessInputDTO
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

	output, err := grantAccessPermissionUC.Execute(input, userID)
	if err != nil {
		log.Printf("error granting access: %v", err.Error())
		sendError(w, http.StatusInternalServerError, buildErrorMessage(opGrantAccess, err))
		return
	}

	sendSuccess(w, opGrantAccess, output)
}

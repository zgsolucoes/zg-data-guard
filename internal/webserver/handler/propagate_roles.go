package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	usecase "github.com/zgsolucoes/zg-data-guard/internal/usecase/database_instance"
)

const (
	opPropagateRoles = "propagate-roles"
)

var propagateRolesUC *usecase.PropagateRolesUseCase

// PropagateRolesHandler godoc
// @BasePath /api/v1
// @Summary Propagates all database role records to the selected database instances, if instances ids are not provided, propagate to all enabled instances
// @Description Propagates all database role records to the selected database instances, if instances ids are not provided, propagate to all enabled instances
// @Tags Database Instance
// @Accept json
// @Produce json
// @Param request body dto.PropagateRolesInputDTO false "Request body"
// @Success 200 {object} PropagateRolesResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /database-instance/propagate-roles [post]
// @Security ApiKeyAuth
func PropagateRolesHandler(w http.ResponseWriter, r *http.Request) {
	userID, hasError := getUserIDFromAuthenticatedRequest(w, r)
	if hasError || userID == emptyString {
		return
	}

	var input dto.PropagateRolesInputDTO
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		log.Printf("error decoding request body: %v", err.Error())
		sendError(w, http.StatusBadRequest, "error decoding request body")
		return
	}

	if input.DatabaseInstancesIDs != nil {
		for _, id := range input.DatabaseInstancesIDs {
			if !validateUUID(w, id, "databaseInstancesIds list contains a value that") {
				return
			}
		}
	}

	outputs, err := propagateRolesUC.Execute(input, userID)
	if err != nil && errors.Is(err, usecase.ErrNoDatabaseInstancesFound) {
		sendError(w, http.StatusNotFound, buildErrorMessage(opPropagateRoles, err))
		return
	}
	if err != nil {
		log.Printf("error in operation %s: %v", opPropagateRoles, err)
		sendError(w, http.StatusInternalServerError, buildErrorMessage(opPropagateRoles, err))
		return
	}

	sendSuccessList(w, opPropagateRoles, outputs, len(outputs), 0, 0)
}

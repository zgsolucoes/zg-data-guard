package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	instanceUC "github.com/zgsolucoes/zg-data-guard/internal/usecase/database_instance"
)

var changeStatusInstanceUC *instanceUC.ChangeStatusDatabaseInstanceUseCase

const opChangeStatusInstance = "change-status-database-instance"

// ChangeStatusDatabaseInstanceHandler godoc
// @BasePath /api/v1
// @Summary Change the status of a database instance (cluster), enabling or disabling it. If disabled, it removes all access permissions from all database users that have access to the instance. It also deactivates all databases from the instance.
// @Description Change the status of a database instance (cluster), enabling or disabling it. If disabled, it removes all access permissions from all database users that have access to the instance. It also deactivates all databases from the instance.
// @Tags Database Instance
// @Accept json
// @Produce json
// @Param request body dto.ChangeStatusInputDTO true "Request body"
// @Success 200 {object} ChangeStatusResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /database-instance/change-status [patch]
// @Security ApiKeyAuth
func ChangeStatusDatabaseInstanceHandler(w http.ResponseWriter, r *http.Request) {
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

	output, err := changeStatusInstanceUC.Execute(input.ID, *input.Enabled, userID)
	if err != nil && errors.Is(err, instanceUC.ErrDatabaseInstanceNotFound) {
		sendError(w, http.StatusNotFound, buildErrorMessage(opChangeStatusInstance, err))
		return
	}
	if err != nil {
		log.Printf("error in operation %s: %v", opChangeStatusInstance, err)
		sendError(w, http.StatusInternalServerError, buildErrorMessage(opChangeStatusInstance, err))
		return
	}
	sendSuccess(w, opChangeStatusInstance, output)
}

package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	"github.com/zgsolucoes/zg-data-guard/internal/usecase/common"
	dbUsecase "github.com/zgsolucoes/zg-data-guard/internal/usecase/database_instance"
)

const opUpdateDatabaseInstance = "update-database-instance"

var updateDBInstanceUC *dbUsecase.UpdateDatabaseInstanceUseCase

// UpdateDatabaseInstanceHandler godoc
// @BasePath /api/v1
// @Summary Update a database instance
// @Description Update an existing database instance
// @Tags Database Instance
// @Accept json
// @Produce json
// @Param id query string true "Database Instance ID"
// @Param request body dto.DatabaseInstanceInputDTO true "Request body"
// @Success 200 {object} CreateDatabaseInstanceResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /database-instance [put]
// @Security ApiKeyAuth
func UpdateDatabaseInstanceHandler(w http.ResponseWriter, r *http.Request) {
	id, hasError := getIDFromQueryParamsAndValidate(w, r)
	if hasError {
		return
	}

	userID, hasError := getUserIDFromAuthenticatedRequest(w, r)
	if hasError || userID == emptyString {
		return
	}

	var input dto.DatabaseInstanceInputDTO
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

	output, err := updateDBInstanceUC.Execute(input, id, userID)
	if err != nil && errors.Is(err, dbUsecase.ErrDatabaseInstanceNotFound) {
		sendError(w, http.StatusNotFound, buildErrorMessage(opUpdateDatabaseInstance, err))
		return

	}
	if err != nil && errors.Is(err, dbUsecase.ErrHostAlreadyExists) {
		sendError(w, http.StatusConflict, buildErrorMessage(opUpdateDatabaseInstance, err))
		return
	}
	if err != nil && (errors.Is(err, common.ErrEcosystemNotFound) || errors.Is(err, common.ErrTechnologyNotFound)) {
		sendError(w, http.StatusNotFound, buildErrorMessage(opUpdateDatabaseInstance, err))
		return
	}
	if err != nil {
		sendError(w, http.StatusInternalServerError, buildErrorMessage(opUpdateDatabaseInstance, err))
		return
	}

	sendSuccess(w, opUpdateDatabaseInstance, output)
}

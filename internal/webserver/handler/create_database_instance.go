package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	"github.com/zgsolucoes/zg-data-guard/internal/usecase/common"
	dbInstanceUsecase "github.com/zgsolucoes/zg-data-guard/internal/usecase/database_instance"
)

var createDBInstanceUC *dbInstanceUsecase.CreateDatabaseInstanceUseCase

const opCreateDatabaseInstance = "create-database-instance"

// CreateDatabaseInstanceHandler godoc
// @BasePath /api/v1
// @Summary Create a database instance
// @Description Create a new database instance
// @Tags Database Instance
// @Accept json
// @Produce json
// @Param request body dto.DatabaseInstanceInputDTO true "Request body"
// @Success 201 {object} CreateDatabaseInstanceResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /database-instance [post]
// @Security ApiKeyAuth
func CreateDatabaseInstanceHandler(w http.ResponseWriter, r *http.Request) {
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

	outputDatabaseInstance, err := createDBInstanceUC.Execute(input, userID)
	if err != nil && errors.Is(err, dbInstanceUsecase.ErrHostAlreadyExists) {
		sendError(w, http.StatusConflict, buildErrorMessage(opCreateDatabaseInstance, err))
		return
	}
	if err != nil && (errors.Is(err, common.ErrEcosystemNotFound) || errors.Is(err, common.ErrTechnologyNotFound)) {
		sendError(w, http.StatusNotFound, buildErrorMessage(opCreateDatabaseInstance, err))
		return
	}
	if err != nil {
		sendError(w, http.StatusInternalServerError, buildErrorMessage(opCreateDatabaseInstance, err))
		return
	}

	sendCreated(w, opCreateDatabaseInstance, outputDatabaseInstance)
}

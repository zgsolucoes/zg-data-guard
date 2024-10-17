package handler

import (
	"errors"
	"net/http"

	dbUsecase "github.com/zgsolucoes/zg-data-guard/internal/usecase/database"
)

const opGetDatabase = "get-database"

var getDatabaseUC *dbUsecase.GetDatabaseUseCase

// GetDatabaseHandler godoc
// @BasePath /api/v1
// @Summary Get a database
// @Description Get an existing database
// @Tags Database
// @Accept json
// @Produce json
// @Param id query string true "Database ID"
// @Success 200 {object} GetDatabaseResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /database [get]
// @Security ApiKeyAuth
func GetDatabaseHandler(w http.ResponseWriter, r *http.Request) {
	id, hasError := getIDFromQueryParamsAndValidate(w, r)
	if hasError {
		return
	}

	outputDatabase, err := getDatabaseUC.Execute(id)
	if err != nil && errors.Is(err, dbUsecase.ErrDatabaseNotFound) {
		sendError(w, http.StatusNotFound, buildErrorMessage(opGetDatabase, err))
		return

	}
	if err != nil {
		sendError(w, http.StatusInternalServerError, buildErrorMessage(opGetDatabase, err))
		return
	}

	sendSuccess(w, opGetDatabase, outputDatabase)
}

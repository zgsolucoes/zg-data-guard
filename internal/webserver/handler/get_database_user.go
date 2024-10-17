package handler

import (
	"errors"
	"net/http"

	"github.com/zgsolucoes/zg-data-guard/internal/usecase/common"
	dbUserUsecase "github.com/zgsolucoes/zg-data-guard/internal/usecase/database_user"
)

const opGetDatabaseUser = "get-database-user"

var getDBUserUC *dbUserUsecase.GetDatabaseUserUseCase

// GetDatabaseUserHandler godoc
// @BasePath /api/v1
// @Summary Get a database user
// @Description Get an existing database user
// @Tags Database User
// @Accept json
// @Produce json
// @Param id query string true "Database User ID"
// @Success 200 {object} GetDatabaseUserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /database-user [get]
// @Security ApiKeyAuth
func GetDatabaseUserHandler(w http.ResponseWriter, r *http.Request) {
	id, hasError := getIDFromQueryParamsAndValidate(w, r)
	if hasError {
		return
	}

	output, err := getDBUserUC.Execute(id)
	if err != nil && errors.Is(err, common.ErrDatabaseUserNotFound) {
		sendError(w, http.StatusNotFound, buildErrorMessage(opGetDatabaseUser, err))
		return
	}
	if err != nil {
		sendError(w, http.StatusInternalServerError, buildErrorMessage(opGetDatabaseUser, err))
		return
	}

	sendSuccess(w, opGetDatabaseUser, output)
}

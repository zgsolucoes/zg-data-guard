package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	dbUserUsecase "github.com/zgsolucoes/zg-data-guard/internal/usecase/database_user"
)

var createDBUserUC *dbUserUsecase.CreateDatabaseUserUseCase

const opCreateDatabaseUser = "create-database-user"

// CreateDatabaseUserHandler godoc
// @BasePath /api/v1
// @Summary Create a database user
// @Description Create a new database user
// @Tags Database User
// @Accept json
// @Produce json
// @Param request body dto.DatabaseUserInputDTO true "Request body"
// @Success 201 {object} CreateDatabaseUserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /database-user [post]
// @Security ApiKeyAuth
func CreateDatabaseUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, hasError := getUserIDFromAuthenticatedRequest(w, r)
	if hasError || userID == emptyString {
		return
	}

	var input dto.DatabaseUserInputDTO
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

	output, err := createDBUserUC.Execute(input, userID)
	if err != nil && errors.Is(err, dbUserUsecase.ErrEmailAlreadyExists) {
		sendError(w, http.StatusConflict, buildErrorMessage(opCreateDatabaseUser, err))
		return
	}
	if err != nil && (errors.Is(err, dbUserUsecase.ErrDatabaseRoleNotFound)) {
		sendError(w, http.StatusNotFound, buildErrorMessage(opCreateDatabaseUser, err))
		return
	}
	if err != nil {
		sendError(w, http.StatusInternalServerError, buildErrorMessage(opCreateDatabaseUser, err))
		return
	}

	sendCreated(w, opCreateDatabaseUser, output)
}

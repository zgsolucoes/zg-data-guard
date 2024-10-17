package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	ecosystemUsecase "github.com/zgsolucoes/zg-data-guard/internal/usecase/ecosystem"
)

const opCreateEcosystem = "create-ecosystem"

var createEcosystemUC *ecosystemUsecase.CreateEcosystemUseCase

// CreateEcosystemHandler godoc
// @BasePath /api/v1
// @Summary Create an ecosystem
// @Description Create a new ecosystem
// @Tags Ecosystem
// @Accept json
// @Produce json
// @Param request body dto.EcosystemInputDTO true "Request body"
// @Success 201 {object} CreateEcosystemResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /ecosystem [post]
// @Security ApiKeyAuth
func CreateEcosystemHandler(w http.ResponseWriter, r *http.Request) {
	userID, hasError := getUserIDFromAuthenticatedRequest(w, r)
	if hasError || userID == emptyString {
		return
	}

	var input dto.EcosystemInputDTO
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

	outputEcosystem, err := createEcosystemUC.Execute(input, userID)
	if err != nil && errors.Is(err, ecosystemUsecase.ErrCodeAlreadyExists) {
		sendError(w, http.StatusConflict, buildErrorMessage(opCreateEcosystem, err))
		return
	}
	if err != nil {
		sendError(w, http.StatusInternalServerError, buildErrorMessage(opCreateEcosystem, err))
		return
	}

	sendCreated(w, opCreateEcosystem, outputEcosystem)
}

package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	"github.com/zgsolucoes/zg-data-guard/internal/usecase/common"
	ecosystemUsecase "github.com/zgsolucoes/zg-data-guard/internal/usecase/ecosystem"
)

const opUpdateEcosystem = "update-ecosystem"

var updateEcosystemUC *ecosystemUsecase.UpdateEcosystemUseCase

// UpdateEcosystemHandler godoc
// @BasePath /api/v1
// @Summary Update an ecosystem
// @Description Update an existing ecosystem
// @Tags Ecosystem
// @Accept json
// @Produce json
// @Param id query string true "Ecosystem ID"
// @Param request body dto.EcosystemInputDTO true "Request body"
// @Success 200 {object} UpdateEcosystemResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /ecosystem [put]
// @Security ApiKeyAuth
func UpdateEcosystemHandler(w http.ResponseWriter, r *http.Request) {
	id, hasError := getIDFromQueryParamsAndValidate(w, r)
	if hasError {
		return
	}

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

	outputEcosystem, err := updateEcosystemUC.Execute(input, id, userID)
	if err != nil && errors.Is(err, common.ErrEcosystemNotFound) {
		sendError(w, http.StatusNotFound, buildErrorMessage(opUpdateEcosystem, err))
		return

	}
	if err != nil && errors.Is(err, ecosystemUsecase.ErrCodeAlreadyExists) {
		sendError(w, http.StatusConflict, buildErrorMessage(opUpdateEcosystem, err))
		return
	}
	if err != nil {
		sendError(w, http.StatusInternalServerError, buildErrorMessage(opUpdateEcosystem, err))
		return
	}

	sendSuccess(w, opUpdateEcosystem, outputEcosystem)
}

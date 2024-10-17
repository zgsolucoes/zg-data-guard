package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	"github.com/zgsolucoes/zg-data-guard/internal/usecase/common"
	technologyUsecase "github.com/zgsolucoes/zg-data-guard/internal/usecase/technology"
)

const opUpdateTechnology = "update-technology"

var updateTechnologyUC *technologyUsecase.UpdateTechnologyUseCase

// UpdateTechnologyHandler godoc
// @BasePath /api/v1
// @Summary Update a technology
// @Description Update an existing technology
// @Tags Technology
// @Accept json
// @Produce json
// @Param id query string true "Technology ID"
// @Param request body dto.TechnologyInputDTO true "Request body"
// @Success 200 {object} UpdateTechnologyResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /technology [put]
// @Security ApiKeyAuth
func UpdateTechnologyHandler(w http.ResponseWriter, r *http.Request) {
	id, hasError := getIDFromQueryParamsAndValidate(w, r)
	if hasError {
		return
	}

	userID, hasError := getUserIDFromAuthenticatedRequest(w, r)
	if hasError || userID == emptyString {
		return
	}

	var input dto.TechnologyInputDTO
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

	outputTechnology, err := updateTechnologyUC.Execute(input, id, userID)
	if err != nil && errors.Is(err, common.ErrTechnologyNotFound) {
		sendError(w, http.StatusNotFound, buildErrorMessage(opUpdateTechnology, err))
		return

	}
	if err != nil && errors.Is(err, technologyUsecase.ErrTechnologyAlreadyExists) {
		sendError(w, http.StatusConflict, buildErrorMessage(opUpdateTechnology, err))
		return
	}
	if err != nil {
		sendError(w, http.StatusInternalServerError, buildErrorMessage(opUpdateTechnology, err))
		return
	}

	sendSuccess(w, opUpdateTechnology, outputTechnology)
}

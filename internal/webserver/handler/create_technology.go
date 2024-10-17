package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	technologyUsecase "github.com/zgsolucoes/zg-data-guard/internal/usecase/technology"
)

const opCreateTechnology = "create-technology"

var createTechnologyUC *technologyUsecase.CreateTechnologyUseCase

// CreateTechnologyHandler godoc
// @BasePath /api/v1
// @Summary Create a technology
// @Description Create a new technology
// @Tags Technology
// @Accept json
// @Produce json
// @Param request body dto.TechnologyInputDTO true "Request body"
// @Success 201 {object} CreateTechnologyResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /technology [post]
// @Security ApiKeyAuth
func CreateTechnologyHandler(w http.ResponseWriter, r *http.Request) {
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

	outputTechnology, err := createTechnologyUC.Execute(input, userID)
	if err != nil && errors.Is(err, technologyUsecase.ErrTechnologyAlreadyExists) {
		sendError(w, http.StatusConflict, buildErrorMessage(opCreateTechnology, err))
		return
	}
	if err != nil {
		sendError(w, http.StatusInternalServerError, buildErrorMessage(opCreateTechnology, err))
		return
	}

	sendCreated(w, opCreateTechnology, outputTechnology)
}

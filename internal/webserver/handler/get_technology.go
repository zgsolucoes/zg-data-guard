package handler

import (
	"errors"
	"net/http"

	"github.com/zgsolucoes/zg-data-guard/internal/usecase/common"
	technologyUsecase "github.com/zgsolucoes/zg-data-guard/internal/usecase/technology"
)

const opGetTechnology = "get-technology"

var getTechnologyUC *technologyUsecase.GetTechnologyUseCase

// GetTechnologyHandler godoc
// @BasePath /api/v1
// @Summary Get a technology
// @Description Get an existing technology
// @Tags Technology
// @Accept json
// @Produce json
// @Param id query string true "Technology ID"
// @Success 200 {object} GetTechnologyResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /technology [get]
// @Security ApiKeyAuth
func GetTechnologyHandler(w http.ResponseWriter, r *http.Request) {
	id, hasError := getIDFromQueryParamsAndValidate(w, r)
	if hasError {
		return
	}

	outputTechnology, err := getTechnologyUC.Execute(id)
	if err != nil && errors.Is(err, common.ErrTechnologyNotFound) {
		sendError(w, http.StatusNotFound, buildErrorMessage(opGetTechnology, err))
		return

	}
	if err != nil {
		sendError(w, http.StatusInternalServerError, buildErrorMessage(opGetTechnology, err))
		return
	}

	sendSuccess(w, opGetTechnology, outputTechnology)
}

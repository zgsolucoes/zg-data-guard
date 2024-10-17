package handler

import (
	"errors"
	"net/http"

	"github.com/zgsolucoes/zg-data-guard/internal/usecase/common"
	ecosystemUsecase "github.com/zgsolucoes/zg-data-guard/internal/usecase/ecosystem"
)

const opGetEcosystem = "get-ecosystem"

var getEcosystemUC *ecosystemUsecase.GetEcosystemUseCase

// GetEcosystemHandler godoc
// @BasePath /api/v1
// @Summary Get an ecosystem
// @Description Get an existing ecosystem
// @Tags Ecosystem
// @Accept json
// @Produce json
// @Param id query string true "Ecosystem ID"
// @Success 200 {object} GetEcosystemResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /ecosystem [get]
// @Security ApiKeyAuth
func GetEcosystemHandler(w http.ResponseWriter, r *http.Request) {
	id, hasError := getIDFromQueryParamsAndValidate(w, r)
	if hasError {
		return
	}

	outputEcosystem, err := getEcosystemUC.Execute(id)
	if err != nil && errors.Is(err, common.ErrEcosystemNotFound) {
		sendError(w, http.StatusNotFound, buildErrorMessage(opGetEcosystem, err))
		return

	}
	if err != nil {
		sendError(w, http.StatusInternalServerError, buildErrorMessage(opGetEcosystem, err))
		return
	}

	sendSuccess(w, opGetEcosystem, outputEcosystem)
}

package handler

import (
	"net/http"

	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	ecosystemUsecase "github.com/zgsolucoes/zg-data-guard/internal/usecase/ecosystem"
)

const opListEcosystems = "list-ecosystems"

var listEcosystemsUC *ecosystemUsecase.ListEcosystemsUseCase

// ListEcosystemsHandler godoc
// @BasePath /api/v1
// @Summary List all existing ecosystems
// @Description List all existing ecosystems
// @Tags Ecosystem
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Limit per page"
// @Success 200 {object} ListEcosystemsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /ecosystems [get]
// @Security ApiKeyAuth
func ListEcosystemsHandler(w http.ResponseWriter, r *http.Request) {
	page, limit := getQueryParamPageAndLimit(r)
	outputEcosystem, err := listEcosystemsUC.Execute(page, limit)
	if err != nil {
		sendError(w, http.StatusInternalServerError, buildErrorMessage(opListEcosystems, err))
		return
	}
	if outputEcosystem == nil {
		outputEcosystem = make([]*dto.EcosystemOutputDTO, 0)
	}

	sendSuccessList(w, opListEcosystems, outputEcosystem, len(outputEcosystem), limit, page)
}

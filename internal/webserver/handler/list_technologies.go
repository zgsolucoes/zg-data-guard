package handler

import (
	"net/http"

	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	technologyUsecase "github.com/zgsolucoes/zg-data-guard/internal/usecase/technology"
)

const opListTechnologies = "list-technologies"

var listTechnologiesUC *technologyUsecase.ListTechnologiesUseCase

// ListTechnologiesHandler godoc
// @BasePath /api/v1
// @Summary List all existing technologies
// @Description List all existing technologies
// @Tags Technology
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Limit per page"
// @Success 200 {object} ListTechnologiesResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /technologies [get]
// @Security ApiKeyAuth
func ListTechnologiesHandler(w http.ResponseWriter, r *http.Request) {
	page, limit := getQueryParamPageAndLimit(r)
	outputTechnologies, err := listTechnologiesUC.Execute(page, limit)
	if err != nil {
		sendError(w, http.StatusInternalServerError, buildErrorMessage(opListTechnologies, err))
		return
	}
	if outputTechnologies == nil {
		outputTechnologies = make([]*dto.TechnologyOutputDTO, 0)
	}

	sendSuccessList(w, opListTechnologies, outputTechnologies, len(outputTechnologies), limit, page)
}

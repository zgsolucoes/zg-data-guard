package handler

import (
	"errors"
	"net/http"

	"github.com/zgsolucoes/zg-data-guard/internal/usecase/common"
	technologyUsecase "github.com/zgsolucoes/zg-data-guard/internal/usecase/technology"
)

const opDeleteTechnology = "delete-technology"

var deleteTechnologyUC *technologyUsecase.DeleteTechnologyUseCase

// DeleteTechnologyHandler godoc
// @Summary Delete a technology
// @Description Delete an existing technology
// @Tags Technology
// @Accept json
// @Produce json
// @Param id query string true "Technology ID"
// @Success 200 {object} DeleteTechnologyResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /technology [delete]
// @Security ApiKeyAuth
func DeleteTechnologyHandler(w http.ResponseWriter, r *http.Request) {
	id, hasError := getIDFromQueryParamsAndValidate(w, r)
	if hasError {
		return
	}

	userID, hasError := getUserIDFromAuthenticatedRequest(w, r)
	if hasError || userID == emptyString {
		return
	}

	err := deleteTechnologyUC.Execute(id, userID)
	if err != nil && errors.Is(err, common.ErrTechnologyNotFound) {
		sendError(w, http.StatusNotFound, buildErrorMessage(opDeleteTechnology, err))
		return

	}
	if err != nil {
		sendError(w, http.StatusInternalServerError, buildErrorMessage(opDeleteTechnology, err))
		return
	}

	sendSuccess(w, opDeleteTechnology, nil)
}

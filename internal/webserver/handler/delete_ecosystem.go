package handler

import (
	"errors"
	"net/http"

	"github.com/zgsolucoes/zg-data-guard/internal/usecase/common"
	ecosystemUsecase "github.com/zgsolucoes/zg-data-guard/internal/usecase/ecosystem"
)

const opDeleteEcosystem = "delete-ecosystem"

var deleteEcosystemUC *ecosystemUsecase.DeleteEcosystemUseCase

// DeleteEcosystemHandler godoc
// @Summary Delete an ecosystem
// @Description Delete an existing ecosystem
// @Tags Ecosystem
// @Accept json
// @Produce json
// @Param id query string true "Ecosystem ID"
// @Success 200 {object} DeleteEcosystemResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /ecosystem [delete]
// @Security ApiKeyAuth
func DeleteEcosystemHandler(w http.ResponseWriter, r *http.Request) {
	id, hasError := getIDFromQueryParamsAndValidate(w, r)
	if hasError {
		return
	}

	userID, hasError := getUserIDFromAuthenticatedRequest(w, r)
	if hasError || userID == emptyString {
		return
	}

	err := deleteEcosystemUC.Execute(id, userID)
	if err != nil && errors.Is(err, common.ErrEcosystemNotFound) {
		sendError(w, http.StatusNotFound, buildErrorMessage(opDeleteEcosystem, err))
		return

	}
	if err != nil {
		sendError(w, http.StatusInternalServerError, buildErrorMessage(opDeleteEcosystem, err))
		return
	}

	sendSuccess(w, opDeleteEcosystem, nil)
}

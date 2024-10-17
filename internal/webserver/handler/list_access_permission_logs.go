package handler

import (
	"net/http"

	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	usecase "github.com/zgsolucoes/zg-data-guard/internal/usecase/access_permission"
)

const (
	opListAccessPermissionLogs = "list-access-permission-logs"
)

var listAccessPermissionLogsUC *usecase.ListAccessPermissionLogsUseCase

// ListAccessPermissionLogsHandler godoc
// @BasePath /api/v1
// @Summary List all existing access permission logs
// @Description List all existing access permission logs
// @Tags Access Permission
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Limit per page"
// @Success 200 {object} ListAccessPermissionLogsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /access-permission/logs [get]
// @Security ApiKeyAuth
func ListAccessPermissionLogsHandler(w http.ResponseWriter, r *http.Request) {
	page, limit := getQueryParamPageAndLimit(r)
	logsDTOs, totalLogsCount, err := listAccessPermissionLogsUC.Execute(page, limit)
	if err != nil {
		sendError(w, http.StatusInternalServerError, buildErrorMessage(opListAccessPermissionLogs, err))
		return
	}
	if logsDTOs == nil {
		logsDTOs = make([]*dto.AccessPermissionLogOutputDTO, 0)
	}

	sendSuccessList(w, opListAccessPermissionLogs, logsDTOs, totalLogsCount, limit, page)
}

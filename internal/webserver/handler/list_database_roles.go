package handler

import (
	"net/http"

	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	dbUsecase "github.com/zgsolucoes/zg-data-guard/internal/usecase/database_role"
)

const (
	opListDatabaseRoles = "list-database-roles"
)

var listDatabaseRolesUC *dbUsecase.ListDatabaseRolesUseCase

// ListDatabaseRolesHandler godoc
// @BasePath /api/v1
// @Summary List all existing database roles
// @Description List all existing database roles
// @Tags Database Role
// @Accept json
// @Produce json
// @Success 200 {object} ListDatabaseRolesResponse
// @Failure 500 {object} ErrorResponse
// @Router /database-roles [get]
// @Security ApiKeyAuth
func ListDatabaseRolesHandler(w http.ResponseWriter, _ *http.Request) {
	outputDTOs, err := listDatabaseRolesUC.Execute()
	if err != nil {
		sendError(w, http.StatusInternalServerError, buildErrorMessage(opListDatabaseRoles, err))
		return
	}
	if outputDTOs == nil {
		outputDTOs = make([]*dto.DatabaseRoleOutputDTO, 0)
	}

	sendSuccessList(w, opListDatabaseRoles, outputDTOs, len(outputDTOs), 0, 0)
}

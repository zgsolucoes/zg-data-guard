package handler

import (
	"net/http"

	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	uc "github.com/zgsolucoes/zg-data-guard/internal/usecase/database_user"
)

const (
	opListDatabaseUsers = "list-database-users"
	paramOnlyEnabled    = "onlyEnabled"
)

var listDatabaseUsersUC *uc.ListDatabaseUsersUseCase

// ListDatabaseUsersHandler godoc
// @BasePath /api/v1
// @Summary List all existing database users
// @Description List all existing database users
// @Tags Database User
// @Accept json
// @Produce json
// @Param onlyEnabled query bool false "Only Enabled Users"
// @Success 200 {object} ListDatabaseUsersResponse
// @Failure 500 {object} ErrorResponse
// @Router /database-users [get]
// @Security ApiKeyAuth
func ListDatabaseUsersHandler(w http.ResponseWriter, r *http.Request) {
	onlyEnabledParam := r.URL.Query().Get(paramOnlyEnabled)
	if !validateBoolQueryParam(w, onlyEnabledParam, paramOnlyEnabled) {
		return
	}
	onlyEnabledParamBool := getQueryParamBoolValue(onlyEnabledParam)
	outputDTOs, err := listDatabaseUsersUC.Execute(onlyEnabledParamBool)
	if err != nil {
		sendError(w, http.StatusInternalServerError, buildErrorMessage(opListDatabaseUsers, err))
		return
	}
	if outputDTOs == nil {
		outputDTOs = make([]*dto.DatabaseUserOutputDTO, 0)
	}

	sendSuccessList(w, opListDatabaseUsers, outputDTOs, len(outputDTOs), 0, 0)
}

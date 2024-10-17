package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/zgsolucoes/zg-data-guard/internal/dto"
	security "github.com/zgsolucoes/zg-data-guard/pkg/security/jwt"
)

func sendError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	errorResponse := map[string]any{
		"message":   msg,
		"errorCode": code,
	}
	_ = json.NewEncoder(w).Encode(errorResponse)
}

func sendCreated(w http.ResponseWriter, operation string, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	sendSuccessfulContent(w, operation, data, 0, 0, 0)
}

func sendSuccess(w http.ResponseWriter, operation string, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	sendSuccessfulContent(w, operation, data, 0, 0, 0)
}

func sendSuccessList(w http.ResponseWriter, operation string, data any, total, limit, page int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	sendSuccessfulContent(w, operation, data, total, limit, page)
}

func sendSuccessfulContent(w http.ResponseWriter, operation string, data any, total, limit, page int) {
	successResponse := map[string]any{
		"message": fmt.Sprintf("operation from handler: %s successful", operation),
	}
	if data != nil {
		successResponse["data"] = data
	}
	if total > 0 {
		successResponse["total"] = total
	}
	if limit > 0 {
		successResponse["limit"] = limit
	}
	if page > 0 {
		successResponse["page"] = page
	}
	_ = json.NewEncoder(w).Encode(successResponse)
}

func buildErrorMessage(operation string, err error) string {
	return fmt.Sprintf("error in operation %s! Cause: %s", operation, err.Error())
}

type ErrorResponse struct {
	Message   string `json:"message"`
	ErrorCode string `json:"errorCode"`
}

type GenerateTokenResponse struct {
	Message string            `json:"message"`
	Data    security.JwtToken `json:"data"`
}

type CreateEcosystemResponse struct {
	Message string                 `json:"message"`
	Data    dto.EcosystemOutputDTO `json:"data"`
}

type UpdateEcosystemResponse struct {
	Message string                 `json:"message"`
	Data    dto.EcosystemOutputDTO `json:"data"`
}

type GetEcosystemResponse struct {
	Message string                 `json:"message"`
	Data    dto.EcosystemOutputDTO `json:"data"`
}

type DeleteEcosystemResponse struct {
	Message string `json:"message"`
}

type ListEcosystemsResponse struct {
	Message string                   `json:"message"`
	Data    []dto.EcosystemOutputDTO `json:"data"`
	Total   int                      `json:"total"`
	Limit   int                      `json:"limit"`
	Page    int                      `json:"page"`
}

type CreateTechnologyResponse struct {
	Message string                  `json:"message"`
	Data    dto.TechnologyOutputDTO `json:"data"`
}

type GetTechnologyResponse struct {
	Message string                  `json:"message"`
	Data    dto.TechnologyOutputDTO `json:"data"`
}

type DeleteTechnologyResponse struct {
	Message string `json:"message"`
}

type UpdateTechnologyResponse struct {
	Message string                  `json:"message"`
	Data    dto.TechnologyOutputDTO `json:"data"`
}

type ListTechnologiesResponse struct {
	Message string                    `json:"message"`
	Data    []dto.TechnologyOutputDTO `json:"data"`
	Total   int                       `json:"total"`
	Limit   int                       `json:"limit"`
	Page    int                       `json:"page"`
}

type CreateDatabaseInstanceResponse struct {
	Message string                        `json:"message"`
	Data    dto.DatabaseInstanceOutputDTO `json:"data"`
}

type GetDatabaseInstanceResponse struct {
	Message string                        `json:"message"`
	Data    dto.DatabaseInstanceOutputDTO `json:"data"`
}

type GetDatabaseResponse struct {
	Message string                `json:"message"`
	Data    dto.DatabaseOutputDTO `json:"data"`
}

type GetDatabaseUserResponse struct {
	Message string                    `json:"message"`
	Data    dto.DatabaseUserOutputDTO `json:"data"`
}

type ListDatabaseInstancesResponse struct {
	Message string                          `json:"message"`
	Data    []dto.DatabaseInstanceOutputDTO `json:"data"`
	Total   int                             `json:"total"`
}

type ListDatabasesResponse struct {
	Message string                  `json:"message"`
	Data    []dto.DatabaseOutputDTO `json:"data"`
	Total   int                     `json:"total"`
}

type TestConnectionResponse struct {
	Message string                        `json:"message"`
	Data    []dto.TestConnectionOutputDTO `json:"data"`
	Total   int                           `json:"total"`
}

type GetDatabaseInstanceCredentialsResponse struct {
	Message string                                   `json:"message"`
	Data    dto.DatabaseInstanceCredentialsOutputDTO `json:"data"`
}

type GetDatabaseUserCredentialsResponse struct {
	Message string                               `json:"message"`
	Data    dto.DatabaseUserCredentialsOutputDTO `json:"data"`
}

type SyncDatabasesResponse struct {
	Message string                       `json:"message"`
	Data    []dto.SyncDatabasesOutputDTO `json:"data"`
	Total   int                          `json:"total"`
}

type PropagateRolesResponse struct {
	Message string                        `json:"message"`
	Data    []dto.PropagateRolesOutputDTO `json:"data"`
	Total   int                           `json:"total"`
}

type ListDatabaseRolesResponse struct {
	Message string                      `json:"message"`
	Data    []dto.DatabaseRoleOutputDTO `json:"data"`
	Total   int                         `json:"total"`
}

type SetupRolesResponse struct {
	Message string                    `json:"message"`
	Data    []dto.SetupRolesOutputDTO `json:"data"`
	Total   int                       `json:"total"`
}

type CreateDatabaseUserResponse struct {
	Message string                    `json:"message"`
	Data    dto.DatabaseUserOutputDTO `json:"data"`
}

type ListDatabaseUsersResponse struct {
	Message string                      `json:"message"`
	Data    []dto.DatabaseUserOutputDTO `json:"data"`
	Total   int                         `json:"total"`
}

type GrantAccessResponse struct {
	Message string                   `json:"message"`
	Data    dto.GrantAccessOutputDTO `json:"data"`
}

type ListAccessPermissionsResponse struct {
	Message string                          `json:"message"`
	Data    []dto.AccessPermissionOutputDTO `json:"data"`
	Total   int                             `json:"total"`
}

type RevokeAccessResponse struct {
	Message string                    `json:"message"`
	Data    dto.RevokeAccessOutputDTO `json:"data"`
}

type ListAccessPermissionLogsResponse struct {
	Message string                             `json:"message"`
	Data    []dto.AccessPermissionLogOutputDTO `json:"data"`
	Total   int                                `json:"total"`
}

type ChangeStatusResponse struct {
	Message string                    `json:"message"`
	Data    dto.ChangeStatusOutputDTO `json:"data"`
}

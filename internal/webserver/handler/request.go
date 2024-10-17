package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/jwtauth"
	"github.com/google/uuid"

	security "github.com/zgsolucoes/zg-data-guard/pkg/security/jwt"
)

const (
	emptyString  = ""
	defaultPage  = 1
	defaultLimit = 250
	trueString   = "true"
	falseString  = "false"
)

func getUserIDFromAuthenticatedRequest(w http.ResponseWriter, r *http.Request) (string, bool) {
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		log.Printf("error getting user id from authenticated request: %v", err.Error())
		sendError(w, http.StatusInternalServerError, "error getting user id from authenticated request")
		return emptyString, true
	}
	userID := claims[security.UserIDCtxKey].(string)
	return userID, false
}

func getIDFromQueryParamsAndValidate(w http.ResponseWriter, r *http.Request) (string, bool) {
	id := r.URL.Query().Get("id")
	if id == emptyString {
		sendError(w, http.StatusBadRequest, "id is required")
		return emptyString, true
	}
	_, err := uuid.Parse(id)
	if err != nil {
		log.Printf("error validanting id: %v", err.Error())
		sendError(w, http.StatusBadRequest, "id is not a valid UUID")
		return emptyString, true
	}
	return id, false
}

func getQueryParamPageAndLimit(r *http.Request) (int, int) {
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = defaultPage
	} else if pageInt <= 0 {
		pageInt = defaultPage
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = defaultLimit
	} else if limitInt <= 0 {
		limitInt = defaultLimit
	}
	return pageInt, limitInt
}

func validateUUIDParam(w http.ResponseWriter, id, param string) bool {
	return id != "" && !validateUUID(w, id, param)
}

func validateUUID(w http.ResponseWriter, id, param string) bool {
	_, err := uuid.Parse(id)
	if err != nil {
		log.Printf("error validanting %s: %v", param, err.Error())
		sendError(w, http.StatusBadRequest, fmt.Sprintf("%s is not a valid UUID", param))
		return false
	}
	return true
}

func validateBoolQueryParam(w http.ResponseWriter, value, param string) bool {
	if value == emptyString {
		return true
	}
	if value != trueString && value != falseString {
		sendError(w, http.StatusBadRequest, fmt.Sprintf("%s must be a boolean value", param))
		return false
	}
	return true
}

func getQueryParamBoolValue(value string) bool {
	return value == trueString
}

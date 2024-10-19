package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/zgsolucoes/zg-data-guard/config"
	userUsecase "github.com/zgsolucoes/zg-data-guard/internal/usecase/user"
)

const zgInternalUserEmail = "zg-service@email.com"

var (
	getUserUC *userUsecase.GetUserUseCase
)

// InternalUserAuthHandler godoc
// Handler for internal user authentication must be used just for testing purposes
func InternalUserAuthHandler(w http.ResponseWriter, _ *http.Request) {
	log.Printf("Token JWT for internal user '%s' was solicited...", zgInternalUserEmail)
	user, err := getUserUC.FindEnabledUserByEmail(zgInternalUserEmail)
	if err != nil && (errors.Is(err, userUsecase.ErrUserNotFound) || errors.Is(err, userUsecase.ErrUserDisabled)) {
		sendError(w, http.StatusForbidden, "You are not allowed to access this application")
		return
	}
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Error when fetching user in database")
		return
	}

	accessToken, err := config.GetJwtHelper().GenerateJwt(user)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Error in jwt token generation")
		return
	}

	log.Printf("Token JWT for internal user '%s' was generated successfully", zgInternalUserEmail)
	sendSuccess(w, "internal-user-login", accessToken)
}

package config

import (
	"os"
	"strconv"

	"github.com/go-chi/jwtauth"

	security "github.com/zgsolucoes/zg-data-guard/pkg/security/jwt"
)

const defaultJwtExpiresIn = 3600             // 1 hour
const developmentJwtExpiresIn = 60 * 60 * 12 // 12 hours

var jwtHelper *security.JwtHelper

func GetJwtHelper() *security.JwtHelper {
	return jwtHelper
}

func initializeJwt() {
	jwtSecret := os.Getenv("JWT_TOKEN_SECRET")
	expiresIn := getJwtExpiresIn()

	jwtAuth := jwtauth.New("HS256", []byte(jwtSecret), nil)
	jwtHelper = security.NewJwtHelper(jwtAuth, expiresIn)
}

func getJwtExpiresIn() int {
	if GetEnvironment() == EnvDevelopment {
		return developmentJwtExpiresIn
	}
	expiresIn := os.Getenv("JWT_EXPIRES_IN")
	if expiresIn == "" {
		return defaultJwtExpiresIn
	}
	intExpiresIn, err := strconv.Atoi(expiresIn)
	if err != nil {
		return defaultJwtExpiresIn
	}
	return intExpiresIn
}

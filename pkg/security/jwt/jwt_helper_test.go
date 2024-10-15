package security

import (
	"testing"

	"github.com/go-chi/jwtauth"
	"github.com/stretchr/testify/assert"

	"github.com/zgsolucoes/zg-data-guard/internal/dto"
)

var (
	jwtAuth   = jwtauth.New("HS256", []byte("secret"), nil)
	jwtHelper = NewJwtHelper(jwtAuth, 10)
)

func TestGivenAnEmptyId_WhenGenerateJwt_ThenShouldReceiveAnError(t *testing.T) {
	accessToken, err := jwtHelper.GenerateJwt(&dto.ApplicationUserOutputDTO{})

	assert.Error(t, err)
	assert.EqualError(t, err, ErrIDEmpty.Error())
	assert.Empty(t, accessToken, "access token should be empty")
}

func TestGivenAValidId_WhenGenerateJwt_ThenShouldReceiveAnError(t *testing.T) {
	accessToken, err := jwtHelper.GenerateJwt(&dto.ApplicationUserOutputDTO{ID: "a271b5d9-0894-4c25-9c69-2805f94a7ec1", Name: "Foo Bar"})

	assert.NoError(t, err)
	assert.NotNil(t, accessToken, "access token should be created")
	assert.NotEmpty(t, accessToken.AccessToken, "access token should not be empty")
	assert.NotEmpty(t, accessToken.ExpiresAt, "expires at should not be empty")
}

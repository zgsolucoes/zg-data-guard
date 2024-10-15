package security

import (
	"errors"
	"log"
	"time"

	"github.com/go-chi/jwtauth"

	"github.com/zgsolucoes/zg-data-guard/internal/dto"
)

var (
	ErrIDEmpty = errors.New("ID is empty, cannot generate token")
)

const UserIDCtxKey = "sub"
const UserNameCtxKey = "name"

type JwtToken struct {
	AccessToken string `json:"accessToken"`
	ExpiresAt   int64  `json:"expiresAt"`
}

type JwtHelper struct {
	Jwt          *jwtauth.JWTAuth
	JwtExpiresIn int
}

func NewJwtHelper(jwt *jwtauth.JWTAuth, jwtExpiresIn int) *JwtHelper {
	return &JwtHelper{
		Jwt:          jwt,
		JwtExpiresIn: jwtExpiresIn,
	}
}

func (helper *JwtHelper) GenerateJwt(user *dto.ApplicationUserOutputDTO) (JwtToken, error) {
	if user == nil || user.ID == "" {
		log.Println("ID is empty, cannot generate token")
		return JwtToken{}, ErrIDEmpty
	}
	log.Println("Generating token JWT for user with ID", user.ID)
	expires := time.Now().Add(time.Second * time.Duration(helper.JwtExpiresIn)).Unix()
	_, tokenString, err := helper.Jwt.Encode(map[string]any{
		UserIDCtxKey:   user.ID,
		UserNameCtxKey: user.Name,
		"exp":          expires,
	})
	if err != nil {
		return JwtToken{}, err
	}

	return JwtToken{AccessToken: "Bearer " + tokenString, ExpiresAt: expires}, nil
}

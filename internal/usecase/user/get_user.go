package user

import (
	"database/sql"
	"errors"
	"log"

	"github.com/zgsolucoes/zg-data-guard/internal/database/storage"
	"github.com/zgsolucoes/zg-data-guard/internal/dto"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserDisabled = errors.New("user disabled")
)

type GetUserUseCase struct {
	UserStorage storage.ApplicationUserStorage
}

func NewGetUserUseCase(userStorage storage.ApplicationUserStorage) *GetUserUseCase {
	return &GetUserUseCase{
		UserStorage: userStorage,
	}
}

func (uc *GetUserUseCase) FindByEmail(email string) (*dto.ApplicationUserOutputDTO, error) {
	user, err := uc.UserStorage.FindByEmail(email)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		log.Printf("User with email %s not found in database", email)
		return nil, ErrUserNotFound
	}
	if err != nil {
		log.Printf("Error fetching user with email %s. Cause: %v", email, err.Error())
		return nil, err
	}

	log.Printf("User %v loaded successfully!", user.ID)
	return &dto.ApplicationUserOutputDTO{
		ID:      user.ID.String(),
		Name:    user.Name,
		Email:   user.Email,
		Enabled: user.Enabled,
	}, nil
}

func (uc *GetUserUseCase) FindEnabledUserByEmail(email string) (*dto.ApplicationUserOutputDTO, error) {
	user, err := uc.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	if !user.Enabled {
		log.Printf("User %v is disabled!", user.ID)
		return nil, ErrUserDisabled
	}
	return user, nil
}

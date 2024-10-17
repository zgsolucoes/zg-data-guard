package user

import (
	"database/sql"
	"testing"

	"github.com/zgsolucoes/zg-data-guard/internal/entity"
	"github.com/zgsolucoes/zg-data-guard/testdata/mocks"

	"github.com/stretchr/testify/assert"
)

func TestGivenAnEmptyEmail_WhenFindByEmail_ThenShouldReturnError(t *testing.T) {
	userStorage := new(mocks.UserStorageMock)
	userStorage.On("FindByEmail", "").Return(&entity.ApplicationUser{}, sql.ErrNoRows).Once()
	uc := NewGetUserUseCase(userStorage)

	appUserObtained, err := uc.FindByEmail("")
	assert.Error(t, err, "user not found error expected")
	assert.EqualError(t, err, ErrUserNotFound.Error())
	assert.Nil(t, appUserObtained)
	userStorage.AssertNumberOfCalls(t, "FindByEmail", 1)
}

func TestGivenAValidEmail_WhenFindByEmail_ThenShouldReturnUser(t *testing.T) {
	userStorage := new(mocks.UserStorageMock)
	email := "foo@email.com.br"
	appUser, err := entity.NewApplicationUser("Foo Bar", email)
	userStorage.On("FindByEmail", email).Return(appUser, nil).Once()
	uc := NewGetUserUseCase(userStorage)

	appUserObtained, err := uc.FindByEmail(email)
	assert.NoError(t, err, "no error expected with a valid user email")
	assert.NotNil(t, appUserObtained, "user should not be nil")
	assert.Equal(t, appUser.ID.String(), appUserObtained.ID)
	assert.Equal(t, appUser.Name, appUserObtained.Name)
	assert.Equal(t, appUser.Email, appUserObtained.Email)
	assert.Equal(t, appUser.Enabled, appUserObtained.Enabled)

	userStorage.AssertNumberOfCalls(t, "FindByEmail", 1)
}

func TestGivenADisabledUser_WhenFindEnabledUserByEmail_ThenShouldReturnError(t *testing.T) {
	userStorage := new(mocks.UserStorageMock)
	email := "foo@email.com.br"
	appUser, _ := entity.NewApplicationUser("Foo Bar", email)
	appUser.Disable()
	userStorage.On("FindByEmail", email).Return(appUser, nil).Once()
	uc := NewGetUserUseCase(userStorage)

	appUserObtained, err := uc.FindEnabledUserByEmail(email)
	assert.Error(t, err, "user disabled error expected")
	assert.EqualError(t, err, ErrUserDisabled.Error())
	assert.Nil(t, appUserObtained)

	userStorage.AssertNumberOfCalls(t, "FindByEmail", 1)
}

func TestGivenAnEnabledUser_WhenFindEnabledUserByEmail_ThenShouldReturnUser(t *testing.T) {
	userStorage := new(mocks.UserStorageMock)
	email := "foo@email.com.br"
	appUser, err := entity.NewApplicationUser("Foo Bar", email)
	userStorage.On("FindByEmail", email).Return(appUser, nil).Once()
	uc := NewGetUserUseCase(userStorage)

	appUserObtained, err := uc.FindEnabledUserByEmail(email)
	assert.NoError(t, err, "no error expected with a enabled user")
	assert.NotNil(t, appUserObtained, "user should not be nil")
	assert.Equal(t, appUser.ID.String(), appUserObtained.ID)
	assert.Equal(t, appUser.Name, appUserObtained.Name)
	assert.Equal(t, appUser.Email, appUserObtained.Email)
	assert.Equal(t, appUser.Enabled, appUserObtained.Enabled)
}

func TestGivenAnErrorInDb_WhenFindEnabledUserByEmail_ThenShouldReturnError(t *testing.T) {
	userStorage := new(mocks.UserStorageMock)
	email := "foo@email.com.br"
	userStorage.On("FindByEmail", email).Return(&entity.ApplicationUser{}, sql.ErrConnDone).Once()
	uc := NewGetUserUseCase(userStorage)

	appUserObtained, err := uc.FindEnabledUserByEmail(email)
	assert.Error(t, err, "error expected when some error in db")
	assert.EqualError(t, err, sql.ErrConnDone.Error())
	assert.Nil(t, appUserObtained)
	userStorage.AssertNumberOfCalls(t, "FindByEmail", 1)
}

package auth

import (
	"app-ez-pwd/internal/logger"
	"app-ez-pwd/internal/settings"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type DoLoginForm struct {
	Username     string `json:"username"`
	PasswordHash string `json:"passwordHash"` // the sha256 pwd
	UserId       int
}

func (f *DoLoginForm) ValidateFront() error {
	return validation.ValidateStruct(f,
		validation.Field(&f.Username, validation.Required, validation.Length(3, 50)),
		validation.Field(&f.PasswordHash, validation.Required, is.Hexadecimal))
}

func (f *DoLoginForm) Validate() map[string]string {
	formErrors := make(map[string]string)

	userId, bCryptPasswordHash, err := GetPasswordHashDB(strings.ToUpper(f.Username))
	if err != nil {
		formErrors["passwordHash"] = "internal error"
	}

	if userId == 0 {
		formErrors["username"] = "username does not exists"
	} else {
		if err = bcrypt.CompareHashAndPassword([]byte(bCryptPasswordHash), []byte(f.PasswordHash)); err != nil {
			formErrors["passwordHash"] = "invalid password"
		} else {
			f.UserId = userId
		}
	}

	return formErrors
}

func (f *DoLoginForm) AuthToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": f.UserId,
	})
	strToken, err := token.SignedString([]byte(settings.Settings.SecretHex))
	if err != nil {
		logger.Logger.Error("err sign token", zap.Error(err))
	}
	return strToken, err
}

type NewAccountForm struct {
	Username     string `json:"username"`
	PasswordHash string `json:"passwordHash"` // sha256
}

func (f NewAccountForm) ValidateFront() error {
	return validation.ValidateStruct(&f,
		validation.Field(&f.Username, validation.Required, validation.Length(3, 50)),
		validation.Field(&f.PasswordHash, validation.Required, is.Hexadecimal))
}

func (f NewAccountForm) Validate() map[string]string {
	formErrors := make(map[string]string)

	existsUser, err := ExistsUsernameDB(f.Username)
	if err != nil {
		formErrors["username"] = "error"
	}
	if existsUser {
		formErrors["username"] = "username already exists"
	}

	return formErrors
}

func (f NewAccountForm) Save() error {
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(f.PasswordHash), bcrypt.DefaultCost)
	err := SaveNewAccount(f.Username, string(passwordHash))
	return err
}

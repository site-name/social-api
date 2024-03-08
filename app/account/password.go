package account

import (
	"errors"
	"net/http"
	"strings"

	"github.com/sitename/sitename/model_helper"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword generates a hash using the bcrypt.GenerateFromPassword
func HashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		panic(err)
	}

	return string(hash)
}

func ComparePassword(hash string, password string) error {
	if password == "" || hash == "" {
		return errors.New("empty password or hash")
	}

	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func (a *ServiceAccount) isPasswordValid(password string) *model_helper.AppError {
	if *a.srv.Config().ServiceSettings.EnableDeveloper {
		return nil
	}

	if err := IsPasswordValidWithSettings(password, &a.srv.Config().PasswordSettings); err != nil {
		return model_helper.NewAppError("User.IsValid", err.Id(), map[string]any{"Min": *a.srv.Config().PasswordSettings.MinimumLength}, "", http.StatusBadRequest)
	}

	return nil
}

// IsPasswordValidWithSettings is a utility functions that checks if the given password
// comforms to the password settings. It returns the error id as error value.
func IsPasswordValidWithSettings(password string, settings *model_helper.PasswordSettings) *ErrInvalidPassword {
	id := "model.user.is_valid.pwd"
	isError := false

	if len(password) < *settings.MinimumLength || len(password) > model_helper.PASSWORD_MAXIMUM_LENGTH {
		isError = true
	}

	if *settings.Lowercase {
		if !strings.ContainsAny(password, model_helper.LOWERCASE_LETTERS) {
			isError = true
		}

		id = id + "_lowercase"
	}

	if *settings.Uppercase {
		if !strings.ContainsAny(password, model_helper.UPPERCASE_LETTERS) {
			isError = true
		}

		id = id + "_uppercase"
	}

	if *settings.Number {
		if !strings.ContainsAny(password, model_helper.NUMBERS) {
			isError = true
		}

		id = id + "_number"
	}

	if *settings.Symbol {
		if !strings.ContainsAny(password, model_helper.SYMBOLS) {
			isError = true
		}

		id = id + "_symbol"
	}

	if isError {
		return NewErrInvalidPassword(id + ".app_error")
	}

	return nil
}

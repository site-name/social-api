package account

import (
	"errors"
	"net/http"
	"strings"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"golang.org/x/crypto/bcrypt"
)

func CheckUserPassword(user *account.User, password string) error {
	if err := ComparePassword(user.Password, password); err != nil {
		return NewErrInvalidPassword("")
	}

	return nil
}

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

func (a *AppAccount) isPasswordValid(password string) *model.AppError {

	if *a.Config().ServiceSettings.EnableDeveloper {
		return nil
	}

	if err := IsPasswordValidWithSettings(password, &a.Config().PasswordSettings); err != nil {
		var invErr *ErrInvalidPassword
		switch {
		case errors.As(err, &invErr):
			return model.NewAppError("User.IsValid", invErr.Id(), map[string]interface{}{"Min": *a.Config().PasswordSettings.MinimumLength}, "", http.StatusBadRequest)
		default:
			return model.NewAppError("User.IsValid", "app.valid_password_generic.app_error", nil, err.Error(), http.StatusInternalServerError)
		}
	}

	return nil
}

// IsPasswordValidWithSettings is a utility functions that checks if the given password
// comforms to the password settings. It returns the error id as error value.
func IsPasswordValidWithSettings(password string, settings *model.PasswordSettings) error {
	id := "model.user.is_valid.pwd"
	isError := false

	if len(password) < *settings.MinimumLength || len(password) > model.PASSWORD_MAXIMUM_LENGTH {
		isError = true
	}

	if *settings.Lowercase {
		if !strings.ContainsAny(password, model.LOWERCASE_LETTERS) {
			isError = true
		}

		id = id + "_lowercase"
	}

	if *settings.Uppercase {
		if !strings.ContainsAny(password, model.UPPERCASE_LETTERS) {
			isError = true
		}

		id = id + "_uppercase"
	}

	if *settings.Number {
		if !strings.ContainsAny(password, model.NUMBERS) {
			isError = true
		}

		id = id + "_number"
	}

	if *settings.Symbol {
		if !strings.ContainsAny(password, model.SYMBOLS) {
			isError = true
		}

		id = id + "_symbol"
	}

	if isError {
		return NewErrInvalidPassword(id + ".app_error")
	}

	return nil
}

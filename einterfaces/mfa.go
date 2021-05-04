package einterfaces

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
)

type MfaInterface interface {
	GenerateSecret(user *account.User) (string, []byte, *model.AppError)
	Activate(user *account.User, token string) *model.AppError
	Deactivate(userID string) *model.AppError
	ValidateToken(secret, token string) (bool, *model.AppError)
}

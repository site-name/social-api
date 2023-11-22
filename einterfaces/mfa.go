package einterfaces

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
)

type MfaInterface interface {
	GenerateSecret(user *model.User) (string, []byte, *model_helper.AppError)
	Activate(user *model.User, token string) *model_helper.AppError
	Deactivate(userID string) *model_helper.AppError
	ValidateToken(secret, token string) (bool, *model_helper.AppError)
}

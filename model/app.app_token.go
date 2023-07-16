package model

import (
	"strings"

	"gorm.io/gorm"
)

// max lengths for some fields
const (
	APP_TOKEN_NAME_MAX_LENGTH       = 128
	APP_TOKEN_AUTH_TOKEN_MAX_LENGTH = 30
)

type AppToken struct {
	Id        string `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	AppId     string `json:"app_id" gorm:"type:uuid;column:AppId"`
	Name      string `json:"name" gorm:"type:varchar(128);column:Name"`
	AuthToken string `json:"auth_token" gorm:"type:varchar(30);column:AuthToken"`
}

func (a *AppToken) BeforeCreate(_ *gorm.DB) error {
	a.commonPre()
	return a.IsValid()
}

func (a *AppToken) BeforeUpdate(_ *gorm.DB) error {
	a.commonPre()
	return a.IsValid()
}

func (a *AppToken) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"model.app_token.is_valid.%s.app_error",
		"app_token_id=",
		"AppToken.IsValid",
	)
	if !IsValidId(a.AppId) {
		return outer("app_id", nil)
	}

	return nil
}

func (a *AppToken) commonPre() {
	a.Name = SanitizeUnicode(a.Name)
	if a.AuthToken == "" {
		a.AuthToken = strings.ReplaceAll(NewId(), "-", "")[0:APP_TOKEN_AUTH_TOKEN_MAX_LENGTH]
	}
}

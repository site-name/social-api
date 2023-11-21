package model

import (
	"net/http"
	"strings"

	"gorm.io/gorm"
)

// max lengths for some fields
const (
	APP_TOKEN_AUTH_TOKEN_MAX_LENGTH = 30
)

type AppToken struct {
	Id        string `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	AppId     string `json:"app_id" gorm:"type:uuid;column:AppId"`
	Name      string `json:"name" gorm:"type:varchar(128);column:Name"`           // varchar(128)
	AuthToken string `json:"auth_token" gorm:"type:varchar(30);column:AuthToken"` // varchar(30)
}

func (a *AppToken) BeforeCreate(_ *gorm.DB) error { a.commonPre(); return a.IsValid() }
func (a *AppToken) BeforeUpdate(_ *gorm.DB) error { a.commonPre(); return a.IsValid() }
func (*AppToken) TableName() string               { return "AppTokens" }

func (a *AppToken) IsValid() *AppError {
	if !IsValidId(a.AppId) {
		return NewAppError("AppToken.IsValid", "model.app_token.is_valid.app_id.app_error", nil, "please provide valid app id", http.StatusBadRequest)
	}

	return nil
}

func (a *AppToken) commonPre() {
	a.Name = SanitizeUnicode(a.Name)
	if a.AuthToken == "" {
		a.AuthToken = strings.ReplaceAll(NewId(), "-", "")[0:APP_TOKEN_AUTH_TOKEN_MAX_LENGTH]
	}
}

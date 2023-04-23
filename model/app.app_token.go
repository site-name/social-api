package model

import (
	"strings"
	"unicode/utf8"
)

// max lengths for some fields
const (
	APP_TOKEN_NAME_MAX_LENGTH       = 128
	APP_TOKEN_AUTH_TOKEN_MAX_LENGTH = 30
)

type AppToken struct {
	Id        string `json:"id"`
	AppId     string `json:"app_id"`
	Name      string `json:"name"`
	AuthToken string `json:"auth_token"`
}

func (a *AppToken) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"model.app_token.is_valid.%s.app_error",
		"app_token_id=",
		"AppToken.IsValid",
	)
	if !IsValidId(a.Id) {
		return outer("id", nil)
	}
	if !IsValidId(a.AppId) {
		return outer("app_id", nil)
	}
	if utf8.RuneCountInString(a.Name) > APP_TOKEN_NAME_MAX_LENGTH {
		return outer("name", &a.Id)
	}
	if len(a.AuthToken) > APP_TOKEN_AUTH_TOKEN_MAX_LENGTH {
		return outer("auth_token", &a.Id)
	}

	return nil
}

func (a *AppToken) ToJSON() string {
	return ModelToJson(a)
}

func (a *AppToken) PreSave() {
	if a.Id == "" {
		a.Id = NewId()
	}
	a.Name = SanitizeUnicode(a.Name)
	if a.AuthToken == "" {
		a.AuthToken = strings.ReplaceAll(NewId(), "-", "")[0:APP_TOKEN_AUTH_TOKEN_MAX_LENGTH]
	}
}

func (a *AppToken) PreUpdate() {
	a.Name = SanitizeUnicode(a.Name)
}

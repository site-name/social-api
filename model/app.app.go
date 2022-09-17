package model

import (
	"strings"
	"unicode/utf8"
)

// max lengths for some fields
const (
	APP_NAME_MAX_LENGTH       = 60
	APP_IDENTIFIER_MAX_LENGTH = 256
	APP_VERSION_MAX_LENGTH    = 60
)

// app type's choices
const (
	LOCAL      = "local"
	THIRDPARTY = "thirdparty"
)

var AppTypeChoiceStrings = map[string]string{
	LOCAL:      "local",
	THIRDPARTY: "thirdparty",
}

type App struct {
	Id               string        `json:"id"`
	Name             string        `json:"name"`
	CreateAt         int64         `json:"create_at"`
	IsActive         *bool         `json:"is_active"`
	Type             string        `json:"type"`
	Identifier       *string       `json:"identifier"`
	Permissions      []*Permission `json:"permissions" db:"-"`
	AboutApp         *string       `json:"about_app"`
	DataPrivacy      *string       `json:"data_privacy"`
	DataPrivacyUrl   *string       `json:"data_privacy_url"`
	HomePageUrl      *string       `json:"homepage_url"`
	SupportUrl       *string       `json:"support_url"`
	ConfigurationUrl *string       `json:"configuration_url"`
	AppUrl           *string       `json:"app_url"`
	Version          *string       `json:"version"`
}

func (a *App) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"app.is_valid.%s.app_error",
		"app_id=",
		"App.IsValid",
	)
	if !IsValidId(a.Id) {
		return outer("id", nil)
	}
	if utf8.RuneCountInString(a.Name) > APP_NAME_MAX_LENGTH {
		return outer("name", &a.Id)
	}
	if AppTypeChoiceStrings[strings.ToLower(a.Type)] == "" {
		return outer("type", &a.Id)
	}
	if a.Identifier != nil && utf8.RuneCountInString(*a.Identifier) > APP_IDENTIFIER_MAX_LENGTH {
		return outer("identifier", &a.Id)
	}
	if a.Version != nil && len(*a.Version) > APP_VERSION_MAX_LENGTH {
		return outer("version", &a.Id)
	}
	if a.CreateAt == 0 {
		return outer("create_at", &a.Id)
	}

	return nil
}

func (a *App) ToJSON() string {
	return ModelToJson(a)
}

func (a *App) PreSave() {
	if a.Id == "" {
		a.Id = NewId()
	}
	a.CreateAt = GetMillis()
	a.Name = SanitizeUnicode(a.Name)
	if a.Identifier != nil {
		a.Identifier = NewString(SanitizeUnicode(*a.Identifier))
	}
	if a.IsActive == nil {
		a.IsActive = NewBool(true)
	}
	if a.AboutApp != nil {
		a.AboutApp = NewString(SanitizeUnicode(*a.AboutApp))
	}
}

func (a *App) PreUpdate() {
	a.Name = SanitizeUnicode(a.Name)
	if a.Identifier != nil {
		a.Identifier = NewString(SanitizeUnicode(*a.Identifier))
	}
	if a.AboutApp != nil {
		a.AboutApp = NewString(SanitizeUnicode(*a.AboutApp))
	}
}

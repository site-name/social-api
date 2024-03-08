package model

import (
	"net/http"

	"gorm.io/gorm"
)

type AppType string

func (a *AppType) IsValid() bool {
	return *a == APP_TYPE_LOCAL || *a == APP_TYPE_THIRDPARTY
}

// app type's choices
const (
	APP_TYPE_LOCAL      AppType = "local"
	APP_TYPE_THIRDPARTY AppType = "thirdparty"
)

type App struct {
	Id               string        `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	Name             string        `json:"name" gorm:"type:varchar(60);column:Name"`
	CreateAt         int64         `json:"create_at" gorm:"autoCreateTime:milli;type:bigint;column:CreateAt"`
	IsActive         *bool         `json:"is_active" gorm:"default:true;column:IsActive"` // default true
	Type             AppType       `json:"type" gorm:"type:varchar(15);column:Type"`
	Identifier       *string       `json:"identifier" gorm:"type:varchar(256);column:Identifier"`
	Permissions      []*Permission `json:"permissions" gorm:"-"`
	AboutApp         *string       `json:"about_app" gorm:"column:AboutApp"`
	DataPrivacy      *string       `json:"data_privacy" gorm:"column:DataPrivacy"`
	DataPrivacyUrl   *string       `json:"data_privacy_url" gorm:"column:DataPrivacyUrl"`
	HomePageUrl      *string       `json:"homepage_url" gorm:"column:HomePageUrl"`
	SupportUrl       *string       `json:"support_url" gorm:"column:SupportUrl"`
	ConfigurationUrl *string       `json:"configuration_url" gorm:"column:ConfigurationUrl"`
	AppUrl           *string       `json:"app_url" gorm:"column:AppUrl"`
	Version          *string       `json:"version" gorm:"type:varchar(60);column:Version"`
}

func (a *App) BeforeCreate(_ *gorm.DB) error { a.commonPre(); return a.IsValid() }
func (a *App) BeforeUpdate(_ *gorm.DB) error { a.commonPre(); return a.IsValid() }
func (*App) TableName() string               { return "Apps" }

func (a *App) IsValid() *AppError {
	if !a.Type.IsValid() {
		return NewAppError("App.IsValid", "model.app.is_valid.type.app_error", nil, "please provide valid app type", http.StatusBadRequest)
	}

	return nil
}

func (a *App) commonPre() {
	a.Name = SanitizeUnicode(a.Name)
	if a.Identifier != nil {
		*a.Identifier = SanitizeUnicode(*a.Identifier)
	}
	if a.AboutApp != nil {
		*a.AboutApp = SanitizeUnicode(*a.AboutApp)
	}
}

package model

import (
	"net/http"

	"gorm.io/gorm"
)

type MenuItemTranslation struct {
	Id           string           `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	LanguageCode LanguageCodeEnum `json:"language_code" gorm:"type:varchar(3);column:LanguageCode"`
	MenuItemID   string           `json:"menu_item_id" gorm:"type:uuid;column:MenuItemID"`
	Name         string           `json:"name" gorm:"type:varchar(128);column:Name"`
}

func (c *MenuItemTranslation) BeforeCreate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *MenuItemTranslation) BeforeUpdate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *MenuItemTranslation) TableName() string             { return MenuItemTranslationTableName }

func (m *MenuItemTranslation) IsValid() *AppError {
	if !IsValidId(m.MenuItemID) {
		return NewAppError("MenuItemTranslation.IsValid", "model.menu_item_translation.is_valid.menu_item_id.app_error", nil, "please provide valid menu item id", http.StatusBadRequest)
	}
	if !m.LanguageCode.IsValid() {
		return NewAppError("MenuItemTranslation.IsValid", "model.menu_item_translation.is_valid.language_code.app_error", nil, "please provide valid language code", http.StatusBadRequest)
	}

	return nil
}

func (m *MenuItemTranslation) commonPre() {
	m.Name = SanitizeUnicode(m.Name)
}

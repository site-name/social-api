package model

import (
	"net/http"

	"gorm.io/gorm"
)

const (
	SHOP_TRANSLATION_NAME_MAX_LENGTH        = 110
	SHOP_TRANSLATION_DESCRIPTION_MAX_LENGTH = 220
)

type ShopTranslation struct {
	Id           string           `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	LanguageCode LanguageCodeEnum `json:"language_code" gorm:"type:varchar(5);column:LanguageCode;unique"`
	Name         string           `json:"name" gorm:"type:varchar(250);column:Name"`
	Description  string           `json:"description" gorm:"type:varchar(1000);column:Description"`
	CreateAt     int64            `json:"create_at" gorm:"type:bigint;autoCreateTime:milli;column:CreateAt"`
	UpdateAt     int64            `json:"update_at" gorm:"type:bigint;autoCreateTime:milli;autoUpdateTime:milli;column:UpdateAt"`
}

func (c *ShopTranslation) BeforeCreate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *ShopTranslation) BeforeUpdate(_ *gorm.DB) error {
	c.commonPre()
	c.CreateAt = 0 // prevent update manually
	return c.IsValid()
}
func (c *ShopTranslation) TableName() string { return ShopTranslationTableName }

func (s *ShopTranslation) commonPre() {
	s.Name = SanitizeUnicode(s.Name)
	s.Description = SanitizeUnicode(s.Description)
}

func (s *ShopTranslation) IsValid() *AppError {
	if !s.LanguageCode.IsValid() {
		return NewAppError("ShopTranslation.IsValid", "model.shop_translation.is_valid.language_code.app_error", nil, "please provide valid language code", http.StatusBadRequest)
	}
	return nil
}

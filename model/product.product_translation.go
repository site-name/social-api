package model

import (
	"net/http"

	"github.com/Masterminds/squirrel"
	"gorm.io/gorm"
)

type ProductTranslation struct {
	Id           UUID             `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	LanguageCode LanguageCodeEnum `json:"language_code" gorm:"type:varchar(5);column:LanguageCode;index:languagecode_productid_key"`
	ProductID    UUID             `json:"product_id" gorm:"type:uuid;column:ProductID;index:languagecode_productid_key"`
	Name         string           `json:"name" gorm:"type:varchar(250);column:Name"`
	Description  StringInterface  `json:"description" gorm:"type:jsonb;column:Description"`
	SeoTranslation
}

func (c *ProductTranslation) BeforeCreate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *ProductTranslation) BeforeUpdate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *ProductTranslation) TableName() string             { return ProductTranslationTableName }
func (p *ProductTranslation) String() string                { return p.Name }

// ProductTranslationFilterOption is used to build squirrel sql queries
type ProductTranslationFilterOption struct {
	Conditions squirrel.Sqlizer
}

func (p *ProductTranslation) IsValid() *AppError {
	if !IsValidId(p.ProductID) {
		return NewAppError("ProductTranslation.IsValid", "model.product_translation.is_valid.product_id.app_error", nil, "please provide valid product id", http.StatusBadRequest)
	}
	if !p.LanguageCode.IsValid() {
		return NewAppError("ProductTranslation.IsValid", "model.product_translation.is_valid.language_code.app_error", nil, "please provide valid language code", http.StatusBadRequest)
	}

	return nil
}

func (p *ProductTranslation) commonPre() {
	p.SeoTranslation.commonPre()
	p.Name = SanitizeUnicode(p.Name)
	if p.Description == nil {
		p.Description = StringInterface{}
	}
}

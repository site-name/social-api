package model

import (
	"net/http"

	"gorm.io/gorm"
)

// unique together language_code, page_id
type PageTranslation struct {
	Id           string           `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	LanguageCode LanguageCodeEnum `json:"language_code" gorm:"type:varchar(5);column:LanguageCode;index:languagecode_pageid_key"` // unique with page_id
	PageID       string           `json:"page_id" gorm:"type:uuid;column:PageID;index:languagecode_pageid_key"`
	Title        string           `json:"title" gorm:"type:varchar(255);column:Title"`
	Content      *StringInterface `json:"content" gorm:"type:jsonb;column:Content"`
	SeoTranslation
}

func (c *PageTranslation) BeforeCreate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *PageTranslation) BeforeUpdate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *PageTranslation) TableName() string             { return PageTranslationTableName }

func (p *PageTranslation) IsValid() *AppError {
	if !IsValidId(p.PageID) {
		return NewAppError("PageTranslation.IsValid", "model.page_translation.is_valid.page_id.app_error", nil, "please provide valid page id", http.StatusBadRequest)
	}
	if !p.LanguageCode.IsValid() {
		return NewAppError("PageTranslation.IsValid", "model.page_translation.is_valid.language_code.app_error", nil, "please provide valid language code", http.StatusBadRequest)
	}

	return nil
}

func (p *PageTranslation) commonPre() {
	p.Title = SanitizeUnicode(p.Title)
}

func (p *PageTranslation) String() string {
	if p.Title != "" {
		return p.Title
	}
	return p.Id
}

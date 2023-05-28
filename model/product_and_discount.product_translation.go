package model

import (
	"strings"
	"unicode/utf8"

	"github.com/Masterminds/squirrel"
	"golang.org/x/text/language"
)

type ProductTranslation struct {
	Id           string           `json:"id"`
	LanguageCode string           `json:"language_code"`
	ProductID    string           `json:"product_id"`
	Name         string           `json:"name"`
	Description  *StringInterface `json:"description"`
	SeoTranslation
}

// ProductTranslationFilterOption is used to build squirrel sql queries
type ProductTranslationFilterOption struct {
	Id           squirrel.Sqlizer
	LanguageCode squirrel.Sqlizer
	ProductID    squirrel.Sqlizer
	Name         squirrel.Sqlizer
}

func (p *ProductTranslation) String() string {
	return p.Name
}

func (p *ProductTranslation) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"model.product_translation.is_valid.%s.app_error",
		"product_translation_id=",
		"ProductTranslation.IsValid",
	)
	if !IsValidId(p.Id) {
		return outer("id", nil)
	}
	if !IsValidId(p.ProductID) {
		return outer("product_id", &p.Id)
	}
	if tag, err := language.Parse(p.LanguageCode); err != nil || !strings.EqualFold(tag.String(), p.LanguageCode) {
		return outer("language_code", &p.Id)
	}
	if utf8.RuneCountInString(p.Name) > PRODUCT_NAME_MAX_LENGTH {
		return outer("name", &p.Id)
	}

	return nil
}

func (p *ProductTranslation) ToJSON() string {
	return ModelToJson(p)
}

func (p *ProductTranslation) PreSave() {
	if p.Id == "" {
		p.Id = NewId()
	}
	p.commonPre()
}

func (p *ProductTranslation) commonPre() {
	p.LanguageCode = strings.ToLower(p.LanguageCode)
	p.Name = SanitizeUnicode(p.Name)
	if p.SeoTitle != nil {
		*p.SeoTitle = SanitizeUnicode(*p.SeoTitle)
	}
	if p.SeoDescription != nil {
		*p.SeoDescription = SanitizeUnicode(*p.SeoDescription)
	}
}

func (p *ProductTranslation) PreUpdate() {
	p.commonPre()
}

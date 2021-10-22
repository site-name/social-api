package product_and_discount

import (
	"strings"
	"unicode/utf8"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/seo"
	"golang.org/x/text/language"
)

type ProductTranslation struct {
	Id           string                 `json:"id"`
	LanguageCode string                 `json:"language_code"`
	ProductID    string                 `json:"product_id"`
	Name         string                 `json:"name"`
	Description  *model.StringInterface `json:"description"`
	seo.SeoTranslation
}

// ProductTranslationFilterOption is used to build squirrel sql queries
type ProductTranslationFilterOption struct {
	Id           *model.StringFilter
	LanguageCode *model.StringFilter
	ProductID    *model.StringFilter
	Name         *model.StringFilter
}

func (p *ProductTranslation) String() string {
	return p.Name
}

func (p *ProductTranslation) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.product_translation.is_valid.%s.app_error",
		"product_translation_id=",
		"ProductTranslation.IsValid",
	)
	if !model.IsValidId(p.Id) {
		return outer("id", nil)
	}
	if !model.IsValidId(p.ProductID) {
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
	return model.ModelToJson(p)
}

func (p *ProductTranslation) PreSave() {
	if p.Id == "" {
		p.Id = model.NewId()
	}
	p.commonPre()
}

func (p *ProductTranslation) commonPre() {
	p.LanguageCode = strings.ToLower(p.LanguageCode)
	p.Name = model.SanitizeUnicode(p.Name)
	if p.SeoTitle != nil {
		*p.SeoTitle = model.SanitizeUnicode(*p.SeoTitle)
	}
	if p.SeoDescription != nil {
		*p.SeoDescription = model.SanitizeUnicode(*p.SeoDescription)
	}
}

func (p *ProductTranslation) PreUpdate() {
	p.commonPre()
}

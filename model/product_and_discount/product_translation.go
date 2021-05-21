package product_and_discount

import (
	"io"
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
	if tag, err := language.Parse(p.LanguageCode); err != nil || tag.String() != p.LanguageCode {
		return outer("language_code", &p.Id)
	}
	if utf8.RuneCountInString(p.Name) > PRODUCT_NAME_MAX_LENGTH {
		return outer("name", &p.Id)
	}

	return nil
}

func (p *ProductTranslation) ToJson() string {
	return model.ModelToJson(p)
}

func ProductTranslationFromJson(data io.Reader) *ProductTranslation {
	var p ProductTranslation
	model.ModelFromJson(&p, data)
	return &p
}

func (p *ProductTranslation) PreSave() {
	if p.Id == "" {
		p.Id = model.NewId()
	}
}

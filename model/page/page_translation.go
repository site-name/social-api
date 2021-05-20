package page

import (
	"io"
	"strings"
	"unicode/utf8"

	"github.com/sitename/sitename/model"
	"golang.org/x/text/language"
)

// max lengths for page translation' fields
const (
	PAGE_TRANSLATION_TITLE_MAX_LENGTH = 255
)

// unique together language_code, page_id
type PageTranslation struct {
	Id           string                 `json:"id"`
	LanguageCode string                 `json:"language_code"`
	PageID       string                 `json:"page_id"`
	Title        string                 `json:"title"` // unique
	Content      *model.StringInterface `json:"content"`
}

func (p *PageTranslation) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.page_translation.is_valid.%s.app_error",
		"page_translation_id=",
		"PageTranslation.IsValid",
	)
	if !model.IsValidId(p.Id) {
		return outer("Id", nil)
	}
	if !model.IsValidId(p.PageID) {
		return outer("page_id", &p.Id)
	}
	if utf8.RuneCountInString(p.Title) > PAGE_TRANSLATION_TITLE_MAX_LENGTH {
		return outer("title", &p.Id)
	}
	if tag, err := language.Parse(p.LanguageCode); err != nil || !strings.EqualFold(tag.String(), p.LanguageCode) {
		return outer("language_code", &p.Id)
	}

	return nil
}

func (p *PageTranslation) ToJson() string {
	return model.ModelToJson(p)
}

func PageTranslationFromJson(data io.Reader) *PageTranslation {
	var pt PageTranslation
	model.ModelFromJson(&pt, data)
	return &pt
}

func (p *PageTranslation) PreSave() {
	p.Title = model.SanitizeUnicode(p.Title)
}

func (p *PageTranslation) PreUpdate() {
	p.Title = model.SanitizeUnicode(p.Title)
}

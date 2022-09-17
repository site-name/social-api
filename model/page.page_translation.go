package model

import (
	"strings"
	"unicode/utf8"

	"golang.org/x/text/language"
)

// max lengths for page translation' fields
const (
	PAGE_TRANSLATION_TITLE_MAX_LENGTH = 255
)

// unique together language_code, page_id
type PageTranslation struct {
	Id           string           `json:"id"`
	LanguageCode string           `json:"language_code"`
	PageID       string           `json:"page_id"`
	Title        string           `json:"title"` // unique
	Content      *StringInterface `json:"content"`
	SeoTranslation
}

func (p *PageTranslation) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"page_translation.is_valid.%s.app_error",
		"page_translation_id=",
		"PageTranslation.IsValid",
	)
	if !IsValidId(p.Id) {
		return outer("Id", nil)
	}
	if !IsValidId(p.PageID) {
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

func (p *PageTranslation) ToJSON() string {
	return ModelToJson(p)
}

func (p *PageTranslation) PreSave() {
	p.Title = SanitizeUnicode(p.Title)
}

func (p *PageTranslation) PreUpdate() {
	p.Title = SanitizeUnicode(p.Title)
}

func (p *PageTranslation) String() string {
	if p.Title != "" {
		return p.Title
	}
	return p.Id
}

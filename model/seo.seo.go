package model

import (
	"unicode/utf8"
)

// max lengths for seo's fields
const (
	SEO_TITLE_MAX_LENGTH       = 70
	SEO_DESCRIPTION_MAX_LENGTH = 300
)

type Seo struct {
	SeoTitle       *string `json:"seo_title"`
	SeoDescription *string `json:"seo_description"`
}

func (s *Seo) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"seo.is_valid.%s.app_error",
		"seo_id=",
		"Seo.IsValid",
	)

	if s.SeoTitle != nil && utf8.RuneCountInString(*s.SeoTitle) > SEO_TITLE_MAX_LENGTH {
		return outer("seo_title", nil)
	}
	if s.SeoDescription != nil && utf8.RuneCountInString(*s.SeoDescription) > SEO_DESCRIPTION_MAX_LENGTH {
		return outer("seo_description", nil)
	}

	return nil
}

func (s *Seo) PreSave() {
	if s.SeoTitle != nil {
		st := SanitizeUnicode(*s.SeoTitle)
		s.SeoTitle = &st
	}
	if s.SeoDescription != nil {
		st := SanitizeUnicode(*s.SeoDescription)
		s.SeoDescription = &st
	}
}

// SeoTranslation represents translation for Seo
type SeoTranslation struct {
	SeoTitle       *string `json:"seo_title"`
	SeoDescription *string `json:"seo_description"`
}

func (s *SeoTranslation) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"seo_translation.is_valid.%s.app_error",
		"seo_translation_id=",
		"SeoTranslation.IsValid")

	if s.SeoTitle != nil && utf8.RuneCountInString(*s.SeoTitle) > SEO_TITLE_MAX_LENGTH {
		return outer("seo_title", nil)
	}
	if s.SeoDescription != nil && utf8.RuneCountInString(*s.SeoDescription) > SEO_DESCRIPTION_MAX_LENGTH {
		return outer("seo_description", nil)
	}

	return nil
}

func (s *SeoTranslation) PreSave() {
	if s.SeoTitle != nil {
		st := SanitizeUnicode(*s.SeoTitle)
		s.SeoTitle = &st
	}
	if s.SeoDescription != nil {
		st := SanitizeUnicode(*s.SeoDescription)
		s.SeoDescription = &st
	}
}
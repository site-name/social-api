package model

import (
	"io"
	"net/http"
	"unicode/utf8"

	"github.com/sitename/sitename/modules/json"
)

const (
	SEO_TITLE_MAX_LENGTH       = 70
	SEO_DESCRIPTION_MAX_LENGTH = 300
)

type Seo struct {
	SeoTitle       *string `json:"seo_title"`
	SeoDescription *string `json:"seo_description"`
}

func (s *Seo) IsValid() *AppError {
	if s.SeoTitle != nil && utf8.RuneCountInString(*s.SeoTitle) > SEO_TITLE_MAX_LENGTH {
		return NewAppError("Seo.IsValid", "model.seo_is_valid.seo_title.app_error", nil, "", http.StatusBadRequest)
	}
	if s.SeoDescription != nil && utf8.RuneCountInString(*s.SeoDescription) > SEO_DESCRIPTION_MAX_LENGTH {
		return NewAppError("Seo.IsValid", "model.seo_is_valid.seo_description.app_error", nil, "", http.StatusBadRequest)
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

func (s *Seo) ToJson() string {
	b, _ := json.JSON.Marshal(s)
	return string(b)
}

func SeoFromJson(data io.Reader) *Seo {
	var seo *Seo
	json.JSON.NewDecoder(data).Decode(seo)

	return seo
}

// SeoTranslation represents translation for Seo
type SeoTranslation struct {
	SeoTitle       *string `json:"seo_title"`
	SeoDescription *string `json:"seo_description"`
}

func (s *SeoTranslation) IsValid() *AppError {
	if s.SeoTitle != nil && utf8.RuneCountInString(*s.SeoTitle) > SEO_TITLE_MAX_LENGTH {
		return NewAppError("SeoTranslation.IsValid", "model.seo_translation.is_valid.seo_title.app_error", nil, "", http.StatusBadRequest)
	}
	if s.SeoDescription != nil && utf8.RuneCountInString(*s.SeoDescription) > SEO_DESCRIPTION_MAX_LENGTH {
		return NewAppError("SeoTranslation.IsValid", "model.seo_translation.is_valid.seo_description.app_error", nil, "", http.StatusBadRequest)
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

func (s *SeoTranslation) ToJson() string {
	b, _ := json.JSON.Marshal(s)
	return string(b)
}

func SeoTranslationFromJson(data io.Reader) *SeoTranslation {
	var seo *SeoTranslation
	json.JSON.NewDecoder(data).Decode(seo)

	return seo
}

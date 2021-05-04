package seo

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/json"
)

const (
	SEO_TITLE_MAX_LENGTH       = 70
	SEO_DESCRIPTION_MAX_LENGTH = 300
)

type Seo struct {
	Id             string  `json:"id"`
	SeoTitle       *string `json:"seo_title"`
	SeoDescription *string `json:"seo_description"`
}

func (s *Seo) createAppError(field string) *model.AppError {
	id := fmt.Sprintf("model.seo.is_valid.%s.app_error", field)
	var details string
	if !strings.EqualFold(field, "id") {
		details = "seo_id=" + s.Id
	}

	return model.NewAppError("Seo.IsValid", id, nil, details, http.StatusBadRequest)
}

func (s *Seo) IsValid() *model.AppError {
	if s.Id == "" {
		return s.createAppError("id")
	}
	if s.SeoTitle != nil && utf8.RuneCountInString(*s.SeoTitle) > SEO_TITLE_MAX_LENGTH {
		return s.createAppError("seo_title")
	}
	if s.SeoDescription != nil && utf8.RuneCountInString(*s.SeoDescription) > SEO_DESCRIPTION_MAX_LENGTH {
		return s.createAppError("seo_description")
	}

	return nil
}

func (s *Seo) PreSave() {
	if s.Id == "" {
		s.Id = model.NewId()
	}
	if s.SeoTitle != nil {
		st := model.SanitizeUnicode(*s.SeoTitle)
		s.SeoTitle = &st
	}
	if s.SeoDescription != nil {
		st := model.SanitizeUnicode(*s.SeoDescription)
		s.SeoDescription = &st
	}
}

func (s *Seo) ToJson() string {
	b, _ := json.JSON.Marshal(s)
	return string(b)
}

func SeoFromJson(data io.Reader) *Seo {
	var seo Seo
	err := json.JSON.NewDecoder(data).Decode(&seo)
	if err != nil {
		return nil
	}
	return &seo
}

// SeoTranslation represents translation for Seo
type SeoTranslation struct {
	Id             string  `json:"is"`
	SeoTitle       *string `json:"seo_title"`
	SeoDescription *string `json:"seo_description"`
}

func (s *SeoTranslation) createAppError(field string) *model.AppError {
	id := fmt.Sprintf("model.seo_translation.is_valid.%s.app_error", field)
	var details string
	if !strings.EqualFold(field, "id") {
		details = "seo_id=" + s.Id
	}

	return model.NewAppError("SeoTranslation.IsValid", id, nil, details, http.StatusBadRequest)
}

func (s *SeoTranslation) IsValid() *model.AppError {
	if s.Id == "" {
		return s.createAppError("id")
	}
	if s.SeoTitle != nil && utf8.RuneCountInString(*s.SeoTitle) > SEO_TITLE_MAX_LENGTH {
		return s.createAppError("seo_title")
	}
	if s.SeoDescription != nil && utf8.RuneCountInString(*s.SeoDescription) > SEO_DESCRIPTION_MAX_LENGTH {
		return s.createAppError("seo_description")
	}

	return nil
}

func (s *SeoTranslation) PreSave() {
	if s.Id == "" {
		s.Id = model.NewId()
	}
	if s.SeoTitle != nil {
		st := model.SanitizeUnicode(*s.SeoTitle)
		s.SeoTitle = &st
	}
	if s.SeoDescription != nil {
		st := model.SanitizeUnicode(*s.SeoDescription)
		s.SeoDescription = &st
	}
}

func (s *SeoTranslation) ToJson() string {
	b, _ := json.JSON.Marshal(s)
	return string(b)
}

func SeoTranslationFromJson(data io.Reader) *SeoTranslation {
	var seo SeoTranslation
	err := json.JSON.NewDecoder(data).Decode(&seo)
	if err != nil {
		return nil
	}
	return &seo
}

package seo

import (
	"io"
	"unicode/utf8"

	"github.com/sitename/sitename/model"
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

func (s *Seo) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.seo.is_valid.%s.app_error",
		"seo_id=",
		"Seo.IsValid")
	if !model.IsValidId(s.Id) {
		return outer("id", nil)
	}
	if s.SeoTitle != nil && utf8.RuneCountInString(*s.SeoTitle) > SEO_TITLE_MAX_LENGTH {
		return outer("seo_title", &s.Id)
	}
	if s.SeoDescription != nil && utf8.RuneCountInString(*s.SeoDescription) > SEO_DESCRIPTION_MAX_LENGTH {
		return outer("seo_description", &s.Id)
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
	return model.ModelToJson(s)
}

func SeoFromJson(data io.Reader) *Seo {
	var seo Seo
	model.ModelFromJson(&seo, data)
	return &seo
}

// SeoTranslation represents translation for Seo
type SeoTranslation struct {
	Id             string  `json:"is"`
	SeoTitle       *string `json:"seo_title"`
	SeoDescription *string `json:"seo_description"`
}

func (s *SeoTranslation) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.seo_translation.is_valid.%s.app_error",
		"seo_translation_id=",
		"SeoTranslation.IsValid")
	if s.Id == "" {
		return outer("id", nil)
	}
	if s.SeoTitle != nil && utf8.RuneCountInString(*s.SeoTitle) > SEO_TITLE_MAX_LENGTH {
		return outer("seo_title", &s.Id)
	}
	if s.SeoDescription != nil && utf8.RuneCountInString(*s.SeoDescription) > SEO_DESCRIPTION_MAX_LENGTH {
		return outer("seo_description", &s.Id)
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
	return model.ModelToJson(s)
}

func SeoTranslationFromJson(data io.Reader) *SeoTranslation {
	var seo SeoTranslation
	model.ModelFromJson(&seo, data)
	return &seo
}

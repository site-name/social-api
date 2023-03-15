package model

import (
	"unicode/utf8"
)

const (
	SHOP_TRANSLATION_NAME_MAX_LENGTH        = 110
	SHOP_TRANSLATION_DESCRIPTION_MAX_LENGTH = 220
)

type ShopTranslation struct {
	Id           string           `json:"id"`
	ShopID       string           `json:"shop_id"`
	LanguageCode LanguageCodeEnum `json:"language_code"`
	Name         string           `json:"name"`
	Description  string           `json:"description"`
	CreateAt     int64            `json:"create_at"`
	UpdateAt     int64            `json:"update_at"`
}

func (s *ShopTranslation) PreSave() {
	if s.Id == "" {
		s.Id = NewId()
	}
	if s.CreateAt == 0 {
		s.CreateAt = GetMillis()
	}
	s.UpdateAt = s.CreateAt
	s.Name = SanitizeUnicode(s.Name)
	s.Description = SanitizeUnicode(s.Description)
}

func (s *ShopTranslation) PreUpdate() {
	s.UpdateAt = GetMillis()
	s.Name = SanitizeUnicode(s.Name)
	s.Description = SanitizeUnicode(s.Description)
}

func (s *ShopTranslation) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"shop_translation.is_valid.%s.app_error",
		"shop_translation_id=",
		"ShopTranslation.IsValid",
	)
	if !IsValidId(s.Id) {
		return outer("id", nil)
	}
	if !IsValidId(s.ShopID) {
		return outer("shop_id", &s.Id)
	}
	if !s.LanguageCode.IsValid() {
		return outer("language_code", &s.Id)
	}
	if utf8.RuneCountInString(s.Name) > SHOP_TRANSLATION_NAME_MAX_LENGTH {
		return outer("name", &s.Id)
	}
	if utf8.RuneCountInString(s.Description) > SHOP_TRANSLATION_DESCRIPTION_MAX_LENGTH {
		return outer("description", &s.Id)
	}

	return nil
}

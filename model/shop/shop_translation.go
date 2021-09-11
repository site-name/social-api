package shop

import (
	"strings"
	"unicode/utf8"

	"github.com/sitename/sitename/model"
)

const (
	SHOP_TRANSLATION_NAME_MAX_LENGTH        = 110
	SHOP_TRANSLATION_DESCRIPTION_MAX_LENGTH = 220
)

type ShopTranslation struct {
	Id           string `json:"id"`
	ShopID       string `json:"shop_id"`
	LanguageCode string `json:"language_code"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	CreateAt     int64  `json:"create_at"`
	UpdateAt     int64  `json:"update_at"`
}

func (s *ShopTranslation) PreSave() {
	if s.Id == "" {
		s.Id = model.NewId()
	}
	if s.CreateAt == 0 {
		s.CreateAt = model.GetMillis()
	}
	s.UpdateAt = s.CreateAt
	s.Name = model.SanitizeUnicode(s.Name)
	s.Description = model.SanitizeUnicode(s.Description)
}

func (s *ShopTranslation) PreUpdate() {
	s.UpdateAt = model.GetMillis()
	s.Name = model.SanitizeUnicode(s.Name)
	s.Description = model.SanitizeUnicode(s.Description)
}

func (s *ShopTranslation) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.shop_translation.is_valid.%s.app_error",
		"shop_translation_id=",
		"ShopTranslation.IsValid",
	)
	if !model.IsValidId(s.Id) {
		return outer("id", nil)
	}
	if !model.IsValidId(s.ShopID) {
		return outer("shop_id", &s.Id)
	}
	if len(s.LanguageCode) > model.LANGUAGE_CODE_MAX_LENGTH || model.Languages[strings.ToUpper(s.LanguageCode)] == "" {
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

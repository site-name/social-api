package menu

import (
	"io"
	"strings"
	"unicode/utf8"

	"github.com/sitename/sitename/model"
	"golang.org/x/text/language"
)

type MenuItemTranslation struct {
	Id           string `json:"id"`
	LanguageCode string `json:"language_code"`
	MenuItemID   string `json:"menu_item_id"`
	Name         string `json:"name"`
}

func (m *MenuItemTranslation) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"moodel.menu_item_translation.is_valid.%s.app_error",
		"menu_item_id=",
		"MenuItemTranslation.IsValid",
	)
	if !model.IsValidId(m.Id) {
		return outer("id", nil)
	}
	if !model.IsValidId(m.MenuItemID) {
		return outer("menu_item_id", &m.Id)
	}
	if utf8.RuneCountInString(m.Name) > MENU_ITEM_NAME_MAX_LENGTH {
		return outer("name", &m.Id)
	}
	if tag, err := language.Parse(m.LanguageCode); err != nil || strings.EqualFold(tag.String(), m.LanguageCode) {
		return outer("language_code", &m.Id)
	}

	return nil
}

func (m *MenuItemTranslation) PreSave() {
	if m.Id == "" {
		m.Id = model.NewId()
	}
	m.Name = model.SanitizeUnicode(m.Name)
}

func (m *MenuItemTranslation) PreUpdate() {
	m.Name = model.SanitizeUnicode(m.Name)
}

func (m *MenuItemTranslation) ToJSON() string {
	return model.ModelToJson(m)
}

func MenuItemTranslationFromJson(data io.Reader) *MenuItemTranslation {
	var m MenuItemTranslation
	model.ModelFromJson(&m, data)
	return &m
}

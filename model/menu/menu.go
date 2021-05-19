package menu

import (
	"io"
	"unicode/utf8"

	"github.com/gosimple/slug"
	"github.com/sitename/sitename/model"
)

// max lengths for some manu's fields
const (
	MENU_NAME_MAX_LENGTH = 250
	MENU_SLUG_MAX_LENGTH = 255
)

type Menu struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
	model.ModelMetadata
}

func (m *Menu) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.menu.is_valid.%s.app_error",
		"menu_id=",
		"Menu.IsValid",
	)
	if !model.IsValidId(m.Id) {
		return outer("id", nil)
	}
	if utf8.RuneCountInString(m.Name) > MENU_NAME_MAX_LENGTH {
		return outer("Name", &m.Id)
	}
	if utf8.RuneCountInString(m.Slug) > MENU_SLUG_MAX_LENGTH {
		return outer("Slug", &m.Id)
	}

	return nil
}

func (m *Menu) PreSave() {
	if m.Id == "" {
		m.Id = model.NewId()
	}
	m.Name = model.SanitizeUnicode(m.Name)
	m.Slug = slug.Make(m.Name)
}

func (m *Menu) PreUpdate() {
	m.Name = model.SanitizeUnicode(m.Name)
	m.Slug = slug.Make(m.Name)
}

func (m *Menu) ToJson() string {
	return model.ModelToJson(m)
}

func MenuFromJson(data io.Reader) *Menu {
	var m Menu
	model.ModelFromJson(&m, data)
	return &m
}

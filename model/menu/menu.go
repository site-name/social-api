package menu

import (
	"unicode/utf8"

	"github.com/Masterminds/squirrel"
	"github.com/gosimple/slug"
	"github.com/sitename/sitename/model"
)

// max lengths for some manu's fields
const (
	MENU_NAME_MAX_LENGTH = 250
	MENU_SLUG_MAX_LENGTH = 255
)

type Menu struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	CreateAt int64  `json:"create_at"` // this field can be used for ordering
	model.ModelMetadata
}

type MenuFilterOptions struct {
	Id   squirrel.Sqlizer
	Name squirrel.Sqlizer
	Slug squirrel.Sqlizer
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
	if m.CreateAt <= 0 {
		return outer("create_at", &m.Id)
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
	if m.CreateAt == 0 {
		m.CreateAt = model.GetMillis()
	}
}

func (m *Menu) PreUpdate() {
	m.Name = model.SanitizeUnicode(m.Name)
}

func (m *Menu) ToJSON() string {
	return model.ModelToJson(m)
}

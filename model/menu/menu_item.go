package menu

import (
	"io"
	"unicode/utf8"

	"github.com/sitename/sitename/model"
)

// max length for some menu item's fiedlds
const (
	MENU_ITEM_NAME_MAX_LENGTH = 128
	MENU_ITEM_URL_MAX_LENGTH  = 256
)

type MenuItem struct {
	Id           string  `json:"id"`
	MenuID       string  `json:"menu_id"`
	Name         string  `json:"name"`
	ParentID     *string `json:"parent_id"`
	Url          *string `json:"url"`
	CategoryID   *string `json:"category_id"`
	CollectionID *string `json:"collection_id"`
	PageID       *string `json:"page_id"`
	model.ModelMetadata
	model.Sortable
}

func (m *MenuItem) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.menu_item.is_valid.%s.app_error",
		"menu_item_id=",
		"MenuItem.IsValid",
	)
	if !model.IsValidId(m.Id) {
		return outer("id", nil)
	}
	if !model.IsValidId(m.MenuID) {
		return outer("menu_id", &m.Id)
	}
	if utf8.RuneCountInString(m.Name) > MENU_ITEM_NAME_MAX_LENGTH {
		return outer("name", &m.Id)
	}
	if m.Url != nil && len(*m.Url) > MENU_ITEM_URL_MAX_LENGTH {
		return outer("url", &m.Id)
	}
	if m.ParentID != nil && !model.IsValidId(*m.ParentID) {
		return outer("parent_id", &m.Id)
	}
	if m.CategoryID != nil && !model.IsValidId(*m.CategoryID) {
		return outer("category_id", &m.Id)
	}
	if m.PageID != nil && !model.IsValidId(*m.PageID) {
		return outer("page_id", &m.Id)
	}
	if m.CollectionID != nil && !model.IsValidId(*m.CollectionID) {
		return outer("collection_id", &m.Id)
	}

	return nil
}

func (m *MenuItem) PreSave() {
	if m.Id == "" {
		m.Id = model.NewId()
	}
	m.Name = model.SanitizeUnicode(m.Name)
}

func (m *MenuItem) PreUpdate() {
	m.Name = model.SanitizeUnicode(m.Name)
}

func (m *MenuItem) ToJSON() string {
	return model.ModelToJson(m)
}

func MenuItemFromJson(data io.Reader) *MenuItem {
	var m MenuItem
	model.ModelFromJson(&m, data)
	return &m
}

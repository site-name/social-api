package model

import (
	"unicode/utf8"

	"github.com/Masterminds/squirrel"
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
	ModelMetadata
	Sortable
}

type MenuItemFilterOptions struct {
	Id     squirrel.Sqlizer
	Name   squirrel.Sqlizer
	MenuID squirrel.Sqlizer
}

func (m *MenuItem) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"menu_item.is_valid.%s.app_error",
		"menu_item_id=",
		"MenuItem.IsValid",
	)
	if !IsValidId(m.Id) {
		return outer("id", nil)
	}
	if !IsValidId(m.MenuID) {
		return outer("menu_id", &m.Id)
	}
	if utf8.RuneCountInString(m.Name) > MENU_ITEM_NAME_MAX_LENGTH {
		return outer("name", &m.Id)
	}
	if m.Url != nil && len(*m.Url) > MENU_ITEM_URL_MAX_LENGTH {
		return outer("url", &m.Id)
	}
	if m.ParentID != nil && !IsValidId(*m.ParentID) {
		return outer("parent_id", &m.Id)
	}
	if m.CategoryID != nil && !IsValidId(*m.CategoryID) {
		return outer("category_id", &m.Id)
	}
	if m.PageID != nil && !IsValidId(*m.PageID) {
		return outer("page_id", &m.Id)
	}
	if m.CollectionID != nil && !IsValidId(*m.CollectionID) {
		return outer("collection_id", &m.Id)
	}

	return nil
}

func (m *MenuItem) PreSave() {
	if m.Id == "" {
		m.Id = NewId()
	}
	m.Name = SanitizeUnicode(m.Name)
}

func (m *MenuItem) PreUpdate() {
	m.Name = SanitizeUnicode(m.Name)
}

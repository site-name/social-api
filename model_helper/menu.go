package model_helper

import (
	"net/http"

	"github.com/gosimple/slug"
	"github.com/sitename/sitename/model"
)

func MenuCommonPre(m *model.Menu) {
	m.Name = SanitizeUnicode(m.Name)
	m.Slug = slug.Make(m.Name)
}

func MenuPreSave(m *model.Menu) {
	if m.ID == "" {
		m.ID = NewId()
	}
	if m.CreatedAt == 0 {
		m.CreatedAt = GetMillis()
	}
	MenuCommonPre(m)
}

func MenuIsValid(m model.Menu) *AppError {
	if !IsValidId(m.ID) {
		return NewAppError("MenuIsValid", "model.menu.is_valid.id.app_error", nil, "", http.StatusBadRequest)
	}
	if m.Name == "" {
		return NewAppError("MenuIsValid", "model.menu.is_valid.name.app_error", nil, "", http.StatusBadRequest)
	}
	if !slug.IsSlug(m.Slug) {
		return NewAppError("MenuIsValid", "model.menu.is_valid.slug.app_error", nil, "", http.StatusBadRequest)
	}
	if m.CreatedAt <= 0 {
		return NewAppError("MenuIsValid", "model.menu.is_valid.created_at.app_error", nil, "", http.StatusBadRequest)
	}

	return nil
}

func MenuItemCommonPre(m *model.MenuItem) {
	m.Name = SanitizeUnicode(m.Name)
}

func MenuItemIsValid(m model.MenuItem) *AppError {
	if !IsValidId(m.ID) {
		return NewAppError("MenuItemIsValid", "model.menu_item.is_valid.id.app_error", nil, "", http.StatusBadRequest)
	}
	if m.Name == "" {
		return NewAppError("MenuItemIsValid", "model.menu_item.is_valid.name.app_error", nil, "", http.StatusBadRequest)
	}
	if !IsValidId(m.MenuID) {
		return NewAppError("MenuItemIsValid", "model.menu_item.is_valid.menu_id.app_error", nil, "", http.StatusBadRequest)
	}
	if !m.CategoryID.IsNil() && !IsValidId(*m.CategoryID.String) {
		return NewAppError("MenuItemIsValid", "model.menu_item.is_valid.category_id.app_error", nil, "", http.StatusBadRequest)
	}
	if !m.ParentID.IsNil() && !IsValidId(*m.ParentID.String) {
		return NewAppError("MenuItemIsValid", "model.menu_item.is_valid.parent_id.app_error", nil, "", http.StatusBadRequest)
	}
	if !m.CollectionID.IsNil() && !IsValidId(*m.CollectionID.String) {
		return NewAppError("MenuItemIsValid", "model.menu_item.is_valid.collection_id.app_error", nil, "", http.StatusBadRequest)
	}
	if !m.PageID.IsNil() && !IsValidId(*m.PageID.String) {
		return NewAppError("MenuItemIsValid", "model.menu_item.is_valid.page_id.app_error", nil, "", http.StatusBadRequest)
	}

	return nil
}

type MenuItemFilterOptions struct {
	CommonQueryOptions
}

type MenuFilterOptions struct {
	CommonQueryOptions
}

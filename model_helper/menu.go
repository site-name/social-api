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

func MenuIsValid(m model.Menu) *AppError {
	if m.ID != "" && !IsValidId(m.ID) {
		return NewAppError("MenuIsValid", "model.menu.is_valid.id.app_error", nil, "", http.StatusBadRequest)
	}
	return nil
}

func MenuItemCommonPre(m *model.MenuItem) {
	m.Name = SanitizeUnicode(m.Name)
}

func MenuItemIsValid(m model.MenuItem) *AppError {
	if m.ID != "" && !IsValidId(m.ID) {
		return NewAppError("MenuItemIsValid", "model.menu_item.is_valid.idd.app_error", nil, "", http.StatusBadRequest)
	}
	if !IsValidId(m.MenuID) {
		return NewAppError("MenuItemIsValid", "model.menu_item.is_valid.menu_id.app_error", nil, "", http.StatusBadRequest)
	}

	return nil
}

package menu

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
)

func (s *ServiceMenu) MenuItemsByOptions(options model_helper.MenuItemFilterOptions) (model.MenuItemSlice, *model_helper.AppError) {
	items, err := s.srv.Store.MenuItem().FilterByOptions(options)
	if err != nil {
		return nil, model_helper.NewAppError("MenuItemsByOptions", "app.menu.menu_items_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return items, nil
}

func (s *ServiceMenu) UpsertMenuItem(item model.MenuItem) (*model.MenuItem, *model_helper.AppError) {
	menuItem, err := s.srv.Store.MenuItem().Upsert(item)
	if err != nil {
		return nil, model_helper.NewAppError("UpsertMenuItem", "app.menu.upsert_item.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return menuItem, nil
}

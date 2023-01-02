package menu

import (
	"net/http"

	"github.com/sitename/sitename/model"
)

func (s *ServiceMenu) MenuItemsByOptions(options *model.MenuItemFilterOptions) ([]*model.MenuItem, *model.AppError) {
	items, err := s.srv.Store.MenuItem().FilterByOptions(options)
	if err != nil {
		return nil, model.NewAppError("MenuItemsByOptions", "app.menu.menu_items_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return items, nil
}

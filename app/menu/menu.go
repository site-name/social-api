/*
NOTE: This package is initialized during server startup (modules/imports does that)
so the init() function get the chance to register a function to create `ServiceAccount`
*/
package menu

import (
	"net/http"

	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

type ServiceMenu struct {
	srv *app.Server
}

type ServiceMenuConfig struct {
	Server *app.Server
}

func init() {
	app.RegisterService(func(s *app.Server) error {
		s.Menu = &ServiceMenu{s}
		return nil
	})
}

func (s *ServiceMenu) MenuByOptions(options *model.MenuFilterOptions) (*model.Menu, *model.AppError) {
	mnu, err := s.srv.Store.Menu().GetByOptions(options)
	if err != nil {
		var statucCode int = http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statucCode = http.StatusNotFound
		}

		return nil, model.NewAppError("MenuByOptions", "app.menu.missing_menu.app_error", nil, err.Error(), statucCode)
	}

	return mnu, nil
}

func (s *ServiceMenu) MenusByOptions(options *model.MenuFilterOptions) ([]*model.Menu, *model.AppError) {
	mnu, err := s.srv.Store.Menu().FilterByOptions(options)
	if err != nil {
		return nil, model.NewAppError("MenuByOptions", "app.menu.missing_menu.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return mnu, nil
}

func (s *ServiceMenu) UpsertMenu(menu *model.Menu) (*model.Menu, *model.AppError) {
	menu, err := s.srv.Store.Menu().Save(menu)
	if err != nil {
		var statusCode = http.StatusInternalServerError
		if _, ok := err.(*store.ErrInvalidInput); ok {
			statusCode = http.StatusBadRequest
		}

		return nil, model.NewAppError("UpsertMenu", "app.menu.upsert_menu.app_error", nil, err.Error(), statusCode)
	}

	return menu, nil
}

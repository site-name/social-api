/*
NOTE: This package is initialized during server startup (modules/imports does that)
so the init() function get the chance to register a function to create `ServiceAccount`
*/
package menu

import (
	"net/http"

	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/menu"
	"github.com/sitename/sitename/store"
)

const (
	missingMenuErrId = "app.menu.missing_menu.app_error"
)

type ServiceMenu struct {
	srv *app.Server
}

type ServiceMenuConfig struct {
	Server *app.Server
}

func init() {
	app.RegisterMenuService(func(s *app.Server) (sub_app_iface.MenuService, error) {
		return &ServiceMenu{
			srv: s,
		}, nil
	})
}

func (s *ServiceMenu) MenuByOptions(options *menu.MenuFilterOptions) (*menu.Menu, *model.AppError) {
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

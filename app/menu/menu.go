/*
	NOTE: This package is initialized during server startup (modules/imports does that)
	so the init() function get the chance to register a function to create `ServiceAccount`
*/
package menu

import (
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

func (a *ServiceMenu) MenuById(id string) (*menu.Menu, *model.AppError) {
	mnu, err := a.srv.Store.Menu().GetById(id)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("MenuById", missingMenuErrId, err)
	}

	return mnu, nil
}

func (a *ServiceMenu) MenuByName(name string) (*menu.Menu, *model.AppError) {
	mnu, err := a.srv.Store.Menu().GetByName(name)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("MenuByName", missingMenuErrId, err)
	}

	return mnu, nil
}

func (a *ServiceMenu) MenuBySlug(slug string) (*menu.Menu, *model.AppError) {
	mnu, err := a.srv.Store.Menu().GetBySlug(slug)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("MenuBySlug", missingMenuErrId, err)
	}

	return mnu, nil
}

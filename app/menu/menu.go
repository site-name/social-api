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

func NewServiceMenuConfig(config *ServiceMenuConfig) sub_app_iface.MenuService {
	return &ServiceMenu{
		srv: config.Server,
	}
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

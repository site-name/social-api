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

type AppMenu struct {
	app.AppIface
}

func init() {
	app.RegisterMenuApp(func(a app.AppIface) sub_app_iface.MenuApp {
		return &AppMenu{a}
	})
}

func (a *AppMenu) MenuById(id string) (*menu.Menu, *model.AppError) {
	mnu, err := a.Srv().Store.Menu().GetById(id)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("MenuById", missingMenuErrId, err)
	}

	return mnu, nil
}

func (a *AppMenu) MenuByName(name string) (*menu.Menu, *model.AppError) {
	mnu, err := a.Srv().Store.Menu().GetByName(name)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("MenuByName", missingMenuErrId, err)
	}

	return mnu, nil
}

func (a *AppMenu) MenuBySlug(slug string) (*menu.Menu, *model.AppError) {
	mnu, err := a.Srv().Store.Menu().GetBySlug(slug)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("MenuBySlug", missingMenuErrId, err)
	}

	return mnu, nil
}

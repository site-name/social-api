package menu

import (
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
)

type AppMenu struct {
	app.AppIface
}

func init() {
	app.RegisterMenuApp(func(a app.AppIface) sub_app_iface.MenuApp {
		return &AppMenu{a}
	})
}

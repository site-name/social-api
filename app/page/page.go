package page

import (
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
)

type AppPage struct {
	app.AppIface
}

func init() {
	app.RegisterPageApp(func(a app.AppIface) sub_app_iface.PageApp {
		return &AppPage{a}
	})
}

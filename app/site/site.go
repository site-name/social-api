package site

import (
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
)

type AppSite struct {
	app.AppIface
}

func init() {
	app.RegisterSiteApp(func(a app.AppIface) sub_app_iface.SiteApp {
		return &AppSite{a}
	})
}

package seo

import (
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
)

type AppSeo struct {
	app.AppIface
}

func init() {
	app.RegisterSeoApp(func(a app.AppIface) sub_app_iface.SeoApp {
		return &AppSeo{a}
	})
}

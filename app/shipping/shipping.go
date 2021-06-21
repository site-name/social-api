package shipping

import (
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
)

type AppShipping struct {
	app.AppIface
}

func init() {
	app.RegisterShippingApp(func(a app.AppIface) sub_app_iface.ShippingApp {
		return &AppShipping{a}
	})
}

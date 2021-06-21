package checkout

import (
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
)

type AppCheckout struct {
	app.AppIface
}

func init() {
	app.RegisterCheckoutApp(func(a app.AppIface) sub_app_iface.CheckoutApp {
		return &AppCheckout{a}
	})
}

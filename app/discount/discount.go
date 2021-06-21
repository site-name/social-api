package discount

import (
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
)

type AppDiscount struct {
	app.AppIface
}

func init() {
	app.RegisterDiscountApp(func(a app.AppIface) sub_app_iface.DiscountApp {
		return &AppDiscount{a}
	})
}

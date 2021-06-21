package warehouse

import (
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
)

type AppWarehouse struct {
	app.AppIface
}

func init() {
	app.RegisterWarehouseApp(func(a app.AppIface) sub_app_iface.WarehouseApp {
		return &AppWarehouse{a}
	})
}

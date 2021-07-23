package attribute

import (
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
)

type AppAttribute struct {
	app app.AppIface
}

func init() {
	app.RegisterAttributeApp(func(a app.AppIface) sub_app_iface.AttributeApp {
		return &AppAttribute{
			app: a,
		}
	})
}

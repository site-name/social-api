package invoice

import (
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
)

type AppInvoice struct {
	app.AppIface
}

func init() {
	app.RegisterInvoiceApp(func(a app.AppIface) sub_app_iface.InvoiceApp {
		return &AppInvoice{a}
	})
}

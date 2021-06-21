package csv

import (
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
)

type AppCsv struct {
	app.AppIface
}

func init() {
	app.RegisterCsvApp(func(a app.AppIface) sub_app_iface.CsvApp {
		return &AppCsv{a}
	})
}

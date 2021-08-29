package csv

import (
	"sync"

	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
)

type AppCsv struct {
	app.AppIface
	sync.WaitGroup
	sync.Mutex
}

func init() {
	app.RegisterCsvApp(func(a app.AppIface) sub_app_iface.CsvApp {
		return &AppCsv{
			AppIface: a,
		}
	})
}

/*
	NOTE: This package is initialized during server startup (modules/imports does that)
	so the init() function get the chance to register a function to create `ServiceAccount`
*/
package csv

import (
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
)

type ServiceCsv struct {
	srv *app.Server
}

func init() {
	app.RegisterCsvService(func(s *app.Server) (sub_app_iface.CsvService, error) {
		return &ServiceCsv{
			srv: s,
		}, nil
	})
}

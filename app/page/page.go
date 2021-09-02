/*
	NOTE: This package is initialized during server startup (modules/imports does that)
	so the init() function get the chance to register a function to create `ServiceAccount`
*/
package page

import (
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
)

type ServicePage struct {
	srv *app.Server
}

func init() {
	app.RegisterPageApp(func(s *app.Server) (sub_app_iface.PageService, error) {
		return &ServicePage{
			srv: s,
		}, nil
	})
}

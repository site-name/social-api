/*
	NOTE: This package is initialized during server startup (modules/imports does that)
	so the init() function get the chance to register a function to create `ServiceAccount`
*/
package shipping

import (
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
)

type ServiceShipping struct {
	srv *app.Server
}

func init() {
	app.RegisterShippingService(func(s *app.Server) (sub_app_iface.ShippingService, error) {
		return &ServiceShipping{
			srv: s,
		}, nil
	})
}

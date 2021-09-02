package shipping

import (
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
)

type ServiceShipping struct {
	srv *app.Server
}

type ServiceShippingConfig struct {
	Server *app.Server
}

func NewServiceShipping(config *ServiceShippingConfig) sub_app_iface.ShippingService {
	return &ServiceShipping{
		srv: config.Server,
	}
}

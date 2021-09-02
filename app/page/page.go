package page

import (
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
)

type ServicePage struct {
	srv *app.Server
}

type ServicePageConfig struct {
	Server *app.Server
}

func NewServicePage(config *ServicePageConfig) sub_app_iface.PageService {
	return &ServicePage{
		srv: config.Server,
	}
}

package seo

import (
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
)

type ServiceSeo struct {
	srv *app.Server
}

type ServiceSeoConfig struct {
	Server *app.Server
}

func NewServiceSeo(config *ServiceSeoConfig) sub_app_iface.SeoService {
	return &ServiceSeo{
		srv: config.Server,
	}
}

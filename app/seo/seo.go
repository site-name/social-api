/*
	NOTE: This package is initialized during server startup (modules/imports does that)
	so the init() function get the chance to register a function to create `ServiceAccount`
*/
package seo

import (
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
)

type ServiceSeo struct {
	srv *app.Server
}

func init() {
	app.RegisterSeoService(func(s *app.Server) (sub_app_iface.SeoService, error) {
		return &ServiceSeo{
			srv: s,
		}, nil
	})
}

type ServiceSeoConfig struct {
	Server *app.Server
}

func NewServiceSeo(config *ServiceSeoConfig) sub_app_iface.SeoService {
	return &ServiceSeo{
		srv: config.Server,
	}
}

/*
	NOTE: This package is initialized during server startup (modules/imports does that)
	so the init() function get the chance to register a function to create `ServiceAccount`
*/
package webhook

import (
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
)

type ServiceWebhook struct {
	srv *app.Server
}

func init() {
	app.RegisterWebhookService(func(s *app.Server) (sub_app_iface.WebhookService, error) {
		return &ServiceWebhook{
			srv: s,
		}, nil
	})
}

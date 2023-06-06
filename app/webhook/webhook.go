/*
NOTE: This package is initialized during server startup (modules/imports does that)
so the init() function get the chance to register a function to create `ServiceAccount`
*/
package webhook

import (
	"github.com/sitename/sitename/app"
)

type ServiceWebhook struct {
	srv *app.Server
}

func init() {
	app.RegisterService(func(s *app.Server) error {
		s.Webhook = &ServiceWebhook{s}
		return nil
	})
}

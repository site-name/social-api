package webhook

import (
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
)

type ServiceWebhook struct {
	srv *app.Server
}

type ServiceWebhookConfig struct {
	Server *app.Server
}

func NewServiceWebhook(config *ServiceWebhookConfig) sub_app_iface.WebhookService {
	return &ServiceWebhook{
		srv: config.Server,
	}
}

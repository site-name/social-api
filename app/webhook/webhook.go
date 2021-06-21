package webhook

import (
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
)

type AppWebhook struct {
	app.AppIface
}

func init() {
	app.RegisterWebhookApp(func(a app.AppIface) sub_app_iface.WebhookApp {
		return &AppWebhook{a}
	})
}

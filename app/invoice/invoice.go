package invoice

import (
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
)

type ServiceInvoice struct {
	srv *app.Server
}

type ServiceInvoiceConfig struct {
	Server *app.Server
}

func NewServiceInvoice(config *ServiceInvoiceConfig) sub_app_iface.InvoiceService {
	return &ServiceInvoice{
		srv: config.Server,
	}
}

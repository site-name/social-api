/*
NOTE: This package is initialized during server startup (modules/imports does that)
so the init() function get the chance to register a function to create `ServiceAccount`
*/
package invoice

import (
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
	"github.com/sitename/sitename/model"
)

type ServiceInvoice struct {
	srv *app.Server
}

func init() {
	app.RegisterInvoiceService(func(s *app.Server) (sub_app_iface.InvoiceService, error) {
		return &ServiceInvoice{
			srv: s,
		}, nil
	})
}

func (s *ServiceInvoice) UpdateInvoice(inVoice *model.Invoice, number *string, url *string) {
	if number != nil {
		inVoice.Number = *number
	}
	if url != nil {
		inVoice.ExternalUrl = *url
	}
}

/*
NOTE: This package is initialized during server startup (modules/imports does that)
so the init() function get the chance to register a function to create `ServiceAccount`
*/
package invoice

import (
	"net/http"

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

func (s *ServiceInvoice) FilterInvoicesByOptions(options *model.InvoiceFilterOptions) ([]*model.Invoice, *model.AppError) {
	invoices, err := s.srv.Store.Invoice().FilterByOptions(options)
	if err != nil {
		return nil, model.NewAppError("FilterInvoicesByOptions", "app.invoice.invoices_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return invoices, nil
}

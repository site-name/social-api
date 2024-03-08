/*
NOTE: This package is initialized during server startup (modules/imports does that)
so the init() function get the chance to register a function to create `ServiceAccount`
*/
package invoice

import (
	"net/http"

	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
)

type ServiceInvoice struct {
	srv *app.Server
}

func init() {
	app.RegisterService(func(s *app.Server) error {
		s.Invoice = &ServiceInvoice{s}
		return nil
	})
}

func (s *ServiceInvoice) FilterInvoicesByOptions(options *model.InvoiceFilterOptions) ([]*model.Invoice, *model_helper.AppError) {
	invoices, err := s.srv.Store.Invoice().FilterByOptions(options)
	if err != nil {
		return nil, model_helper.NewAppError("FilterInvoicesByOptions", "app.invoice.invoices_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return invoices, nil
}

func (s *ServiceInvoice) UpsertInvoice(invoice *model.Invoice) (*model.Invoice, *model_helper.AppError) {
	res, err := s.srv.Store.Invoice().Upsert(invoice)
	if err != nil {
		if appErr, ok := err.(*model_helper.AppError); ok {
			return nil, appErr
		}
		if _, ok := err.(*store.ErrNotFound); ok {
			return nil, model_helper.NewAppError("UpsertInvoice", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "invoice.Id"}, "", http.StatusBadRequest)
		}
		return nil, model_helper.NewAppError("UpsertInvoice", "app.invoice.upserting_invoice.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return res, nil
}

func (s *ServiceInvoice) GetInvoiceByOptions(options *model.InvoiceFilterOptions) (*model.Invoice, *model_helper.AppError) {
	invoice, err := s.srv.Store.Invoice().GetbyOptions(options)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}

		return nil, model_helper.NewAppError("GetInvoiceByOptions", "app.invoice.invoice_by_options.app_error", nil, err.Error(), statusCode)
	}

	return invoice, nil
}

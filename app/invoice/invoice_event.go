package invoice

import (
	"net/http"

	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/invoice"
	"github.com/sitename/sitename/store"
)

// UpsertInvoiceEvent is shortcut for creating invoice events
func (a *AppInvoice) UpsertInvoiceEvent(option *invoice.InvoiceEventOption) (*invoice.InvoiceEvent, *model.AppError) {
	invoiceEvent := new(invoice.InvoiceEvent)

	invoiceEvent.Type = option.Type
	if option.UserID != nil {
		invoiceEvent.UserID = option.UserID
	}
	if option.OrderID != nil {
		invoiceEvent.OrderID = option.OrderID
	}
	if option.InvoiceID != nil {
		invoiceEvent.InvoiceID = option.InvoiceID
	}
	if option.Parameters != nil {
		invoiceEvent.Parameters = option.Parameters
	}

	invoiceEvent, err := a.Srv().Store.InvoiceEvent().Upsert(invoiceEvent)
	if err != nil {
		if appErr, ok := err.(*model.AppError); ok {
			return nil, appErr
		}
		if _, ok := err.(*store.ErrNotFound); ok {
			return nil, model.NewAppError("UpsertInvoiceEvent", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "invoiceEvent.Id"}, "", http.StatusBadRequest)
		}
		return nil, model.NewAppError("UpsertInvoiceEvent", "app.invoice.error_upserting_invoice_event.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return invoiceEvent, nil
}
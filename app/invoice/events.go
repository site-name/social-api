package invoice

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
)

func (a *ServiceInvoice) UpsertInvoiceEvent(option model_helper.InvoiceEventCreationOptions) (*model.InvoiceEvent, *model_helper.AppError) {
	invoiceEvent := model.InvoiceEvent{}

	invoiceEvent.Type = option.Type
	if option.UserID != nil {
		invoiceEvent.UserID.String = option.UserID
	}
	if option.OrderID != nil {
		invoiceEvent.OrderID.String = option.OrderID
	}
	if option.InvoiceID != nil {
		invoiceEvent.InvoiceID.String = option.InvoiceID
	}
	if option.Parameters != nil {
		invoiceEvent.Parameters = option.Parameters
	}

	upsertedInvoiceEvent, err := a.srv.Store.InvoiceEvent().Upsert(invoiceEvent)
	if err != nil {
		if appErr, ok := err.(*model_helper.AppError); ok {
			return nil, appErr
		}
		if _, ok := err.(*store.ErrNotFound); ok {
			return nil, model_helper.NewAppError("UpsertInvoiceEvent", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "invoiceEvent.Id"}, "", http.StatusBadRequest)
		}
		return nil, model_helper.NewAppError("UpsertInvoiceEvent", "app.invoice.error_upserting_invoice_event.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return upsertedInvoiceEvent, nil
}

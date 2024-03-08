package model_helper

import (
	"net/http"
	"net/url"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/model_types"
)

func InvoicePreSave(invoice *model.Invoice) {
	if invoice.ID == "" {
		invoice.ID = NewId()
	}
	if invoice.CreatedAt == 0 {
		invoice.CreatedAt = GetMillis()
	}
	invoice.UpdatedAt = invoice.CreatedAt
}

func InvoicePreUpdate(invoice *model.Invoice) {
	invoice.UpdatedAt = GetMillis()
}

func InvoiceIsValid(invoice model.Invoice) *AppError {
	if !IsValidId(invoice.ID) {
		return NewAppError("InvoiceIsValid", "model.invoice.is_valid.id.app_error", nil, "please provide valid id", http.StatusBadRequest)
	}
	if !invoice.OrderID.IsNil() && !IsValidId(*invoice.OrderID.String) {
		return NewAppError("InvoiceIsValid", "model.invoice.is_valid.order_id.app_error", nil, "please provide valid order id", http.StatusBadRequest)
	}
	if invoice.CreatedAt <= 0 {
		return NewAppError("InvoiceIsValid", "model.invoice.is_valid.created_at.app_error", nil, "please provide valid created at", http.StatusBadRequest)
	}
	if invoice.UpdatedAt <= 0 {
		return NewAppError("InvoiceIsValid", "model.invoice.is_valid.updated_at.app_error", nil, "please provide valid updated at", http.StatusBadRequest)
	}
	if invoice.ExternalURL != "" {
		_, err := url.Parse(invoice.ExternalURL)
		if err != nil {
			return NewAppError("InvoiceIsValid", "model.invoice.is_valid.external_url.app_error", nil, err.Error(), http.StatusBadRequest)
		}
	}

	return nil
}

func InvoiceEventPreSave(invoiceEvent *model.InvoiceEvent) {
	if invoiceEvent.ID == "" {
		invoiceEvent.ID = NewId()
	}
	if invoiceEvent.CreatedAt == 0 {
		invoiceEvent.CreatedAt = GetMillis()
	}
}

func InvoiceEventIsValid(invoiceEvent model.InvoiceEvent) *AppError {
	if !IsValidId(invoiceEvent.ID) {
		return NewAppError("InvoiceEventIsValid", "model.invoice_event.is_valid.id.app_error", nil, "please provide valid id", http.StatusBadRequest)
	}
	if !invoiceEvent.InvoiceID.IsNil() && !IsValidId(*invoiceEvent.InvoiceID.String) {
		return NewAppError("InvoiceEventIsValid", "model.invoice_event.is_valid.invoice_id.app_error", nil, "please provide valid invoice id", http.StatusBadRequest)
	}
	if !invoiceEvent.OrderID.IsNil() && !IsValidId(*invoiceEvent.OrderID.String) {
		return NewAppError("InvoiceEventIsValid", "model.invoice_event.is_valid.order_id.app_error", nil, "please provide valid order id", http.StatusBadRequest)
	}
	if !invoiceEvent.UserID.IsNil() && !IsValidId(*invoiceEvent.UserID.String) {
		return NewAppError("InvoiceEventIsValid", "model.invoice_event.is_valid.user_id.app_error", nil, "please provide valid user id", http.StatusBadRequest)
	}
	if invoiceEvent.CreatedAt <= 0 {
		return NewAppError("InvoiceEventIsValid", "model.invoice_event.is_valid.created_at.app_error", nil, "please provide valid created at", http.StatusBadRequest)
	}
	if invoiceEvent.Type.IsValid() != nil {
		return NewAppError("InvoiceEventIsValid", "model.invoice_event.is_valid.type.app_error", nil, "please provide valid type", http.StatusBadRequest)
	}

	return nil
}

type InvoiceFilterOption struct {
	CommonQueryOptions
}

type InvoiceEventCreationOptions struct {
	Type      model.InvoiceEventType
	InvoiceID *string
	OrderID   *string
	UserID    *string
	// if provided, it should contains below keys:
	//  "number", "url", "invoice_id"
	Parameters model_types.JSONString
}

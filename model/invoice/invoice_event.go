package invoice

import (
	"strings"

	"github.com/sitename/sitename/model"
)

// events for invoice
const (
	REQUESTED          = "requested"
	REQUESTED_DELETION = "requested_deletion"
	CREATED            = "created"
	DELETED            = "deleted"
	SENT               = "sent"
)

var InvoiceEventStrings = map[string]string{
	REQUESTED:          "The invoice was requested",
	REQUESTED_DELETION: "The invoice was requested for deletion",
	CREATED:            "The invoice was created",
	DELETED:            "The invoice was deleted",
	SENT:               "The invoice has been sent",
}

type InvoiceEvent struct {
	Id         string          `json:"id"`
	CreateAt   int64           `json:"create_at"`
	Type       string          `json:"type"`
	InvoiceID  *string         `json:"invoice_id"`
	OrderID    *string         `json:"order_id"`
	UserID     *string         `json:"user_id"`
	Parameters model.StringMap `json:"parameters"`
}

func (i *InvoiceEvent) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.invoice_event.is_valid.%s.app_error",
		"invoice_event_id=",
		"InvoiceEvent.IsValid",
	)
	if !model.IsValidId(i.Id) {
		return outer("id", nil)
	}
	if i.CreateAt == 0 {
		return outer("create_at", &i.Id)
	}
	if InvoiceEventStrings[strings.ToLower(i.Type)] == "" {
		return outer("type", &i.Id)
	}
	if i.UserID != nil && !model.IsValidId(*i.UserID) {
		return outer("user_id", &i.Id)
	}
	if i.OrderID != nil && !model.IsValidId(*i.OrderID) {
		return outer("order_id", &i.Id)
	}
	if i.InvoiceID != nil && !model.IsValidId(*i.InvoiceID) {
		return outer("invoice_id", &i.Id)
	}

	return nil
}

func (i *InvoiceEvent) PreSave() {
	if i.Id == "" {
		i.Id = model.NewId()
	}
	i.CreateAt = model.GetMillis()
}
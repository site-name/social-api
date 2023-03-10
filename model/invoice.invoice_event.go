package model

// max lengths for some invoice event's fields
const (
	INVOICE_EVENT_TYPE_MAX_LENGTH = 255
)

// types for invoice events
const (
	REQUESTED          = "requested"
	REQUESTED_DELETION = "requested_deletion"
	CREATED            = "created"
	DELETED            = "deleted"
	SENT               = "sent"
)

var InVoiceEventTypeString = map[string]string{
	REQUESTED:          "The invoice was requested",
	REQUESTED_DELETION: "The invoice was requested for deletion",
	CREATED:            "The invoice was created",
	DELETED:            "The invoice was deleted",
	SENT:               "The invoice has been sent",
}

// Model used to store events that happened during the invoice lifecycle.
type InvoiceEvent struct {
	Id         string    `json:"id"`
	CreateAt   int64     `json:"create_at"`
	Type       string    `json:"type"`
	InvoiceID  *string   `json:"invoice_id"`
	OrderID    *string   `json:"order_id"`
	UserID     *string   `json:"user_id"`
	Parameters StringMap `json:"parameters"`
}

// InvoiceEventOption is used for creating new invoice events
type InvoiceEventOption struct {
	Type       string
	InvoiceID  *string
	OrderID    *string
	UserID     *string
	Parameters StringMap // keys should be "number", "url", "invoice_id"
}

func (i *InvoiceEvent) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"model.invoice_event.is_valid.%s.app_error",
		"invoice_event_id=",
		"InvoiceEvent.IsValid",
	)
	if !IsValidId(i.Id) {
		return outer("id", nil)
	}
	if i.CreateAt == 0 {
		return outer("create_at", &i.Id)
	}
	if _, exist := InVoiceEventTypeString[i.Type]; !exist {
		return outer("type", &i.Id)
	}
	if i.UserID != nil && !IsValidId(*i.UserID) {
		return outer("user_id", &i.Id)
	}
	if i.OrderID != nil && !IsValidId(*i.OrderID) {
		return outer("order_id", &i.Id)
	}
	if i.InvoiceID != nil && !IsValidId(*i.InvoiceID) {
		return outer("invoice_id", &i.Id)
	}

	return nil
}

func (i *InvoiceEvent) PreSave() {
	if i.Id == "" {
		i.Id = NewId()
	}
	i.CreateAt = GetMillis()
	if i.Parameters == nil {
		i.Parameters = make(StringMap)
	}
}

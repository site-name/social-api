package model

import (
	"net/http"

	"gorm.io/gorm"
)

type InvoiceEventType string

// types for invoice events
const (
	INVOICE_EVENT_TYPE_REQUESTED          InvoiceEventType = "requested"
	INVOICE_EVENT_TYPE_REQUESTED_DELETION InvoiceEventType = "requested_deletion"
	INVOICE_EVENT_TYPE_CREATED            InvoiceEventType = "created"
	INVOICE_EVENT_TYPE_DELETED            InvoiceEventType = "deleted"
	INVOICE_EVENT_TYPE_SENT               InvoiceEventType = "sent"
)

func (e InvoiceEventType) IsValid() bool {
	return InVoiceEventTypeString[e] != ""
}

var InVoiceEventTypeString = map[InvoiceEventType]string{
	INVOICE_EVENT_TYPE_REQUESTED:          "The invoice was requested",
	INVOICE_EVENT_TYPE_REQUESTED_DELETION: "The invoice was requested for deletion",
	INVOICE_EVENT_TYPE_CREATED:            "The invoice was created",
	INVOICE_EVENT_TYPE_DELETED:            "The invoice was deleted",
	INVOICE_EVENT_TYPE_SENT:               "The invoice has been sent",
}

// Model used to store events that happened during the invoice lifecycle.
type InvoiceEvent struct {
	Id        UUID             `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	CreateAt  int64            `json:"create_at" gorm:"type:bigint;autoCreateTime:milli;column:CreateAt"`
	Type      InvoiceEventType `json:"type" gorm:"type:varchar(255);column:Type"`
	InvoiceID *UUID            `json:"invoice_id" gorm:"type:uuid;column:InvoiceID"`
	OrderID   *UUID            `json:"order_id" gorm:"type:uuid;column:OrderID"`
	UserID    *UUID            `json:"user_id" gorm:"type:uuid;column:UserID"`
	// if provided, it should contains below keys:
	//  "number", "url", "invoice_id"
	Parameters StringMap `json:"parameters" gorm:"type:jsonb;column:Parameters"`
}

func (c *InvoiceEvent) BeforeCreate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *InvoiceEvent) BeforeUpdate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *InvoiceEvent) TableName() string             { return InvoiceEventTableName }

// InvoiceEventCreationOptions is used for creating new invoice events
type InvoiceEventCreationOptions struct {
	Type      InvoiceEventType
	InvoiceID *string
	OrderID   *string
	UserID    *string
	// if provided, it should contains below keys:
	//  "number", "url", "invoice_id"
	Parameters StringMap
}

func (i *InvoiceEvent) IsValid() *AppError {
	if !i.Type.IsValid() {
		return NewAppError("InvoiceEvent.IsValid", "model.invoice_event.is_valid.type.app_error", nil, "please provide valid invoice event type", http.StatusBadRequest)
	}
	if i.UserID != nil && !IsValidId(*i.UserID) {
		return NewAppError("InvoiceEvent.IsValid", "model.invoice_event.is_valid.user_id.app_error", nil, "please provide valid invoice event user id", http.StatusBadRequest)
	}
	if i.OrderID != nil && !IsValidId(*i.OrderID) {
		return NewAppError("InvoiceEvent.IsValid", "model.invoice_event.is_valid.order_id.app_error", nil, "please provide valid invoice event order id", http.StatusBadRequest)
	}
	if i.InvoiceID != nil && !IsValidId(*i.InvoiceID) {
		return NewAppError("InvoiceEvent.IsValid", "model.invoice_event.is_valid.invoice_id.app_error", nil, "please provide valid invoice event invoice id", http.StatusBadRequest)
	}

	return nil
}

func (i *InvoiceEvent) commonPre() {
	if i.Parameters == nil {
		i.Parameters = make(StringMap)
	}
}

package order

import (
	"io"
	"strings"

	"github.com/sitename/sitename/model"
)

// valid values for order event's type
const (
	ORDER_EVENT_TYPE__CONFIRMED                            = "confirmed"
	ORDER_EVENT_TYPE__DRAFT_CREATED                        = "draft_created"
	ORDER_EVENT_TYPE__DRAFT_CREATED_FROM_REPLACE           = "draft_created_from_replace"
	ORDER_EVENT_TYPE__ADDED_PRODUCTS                       = "added_products"
	ORDER_EVENT_TYPE__REMOVED_PRODUCTS                     = "removed_products"
	ORDER_EVENT_TYPE__PLACED                               = "placed"
	ORDER_EVENT_TYPE__PLACED_FROM_DRAFT                    = "placed_from_draft"
	ORDER_EVENT_TYPE__OVERSOLD_ITEMS                       = "oversold_items"
	ORDER_EVENT_TYPE__CANCELED                             = "canceled"
	ORDER_EVENT_TYPE__ORDER_MARKED_AS_PAID                 = "order_marked_as_paid"
	ORDER_EVENT_TYPE__ORDER_FULLY_PAID                     = "order_fully_paid"
	ORDER_EVENT_TYPE__ORDER_REPLACEMENT_CREATED            = "order_replacement_created"
	ORDER_EVENT_TYPE__ORDER_DISCOUNT_ADDED                 = "order_discount_added"
	ORDER_EVENT_TYPE__ORDER_DISCOUNT_AUTOMATICALLY_UPDATED = "order_discount_automatically_updated"
	ORDER_EVENT_TYPE__ORDER_DISCOUNT_UPDATED               = "order_discount_updated"
	ORDER_EVENT_TYPE__ORDER_DISCOUNT_DELETED               = "order_discount_deleted"
	ORDER_EVENT_TYPE__ORDER_LINE_DISCOUNT_UPDATED          = "order_line_discount_updated"
	ORDER_EVENT_TYPE__ORDER_LINE_DISCOUNT_REMOVED          = "order_line_discount_removed"
	ORDER_EVENT_TYPE__UPDATED_ADDRESS                      = "updated_address"
	ORDER_EVENT_TYPE__EMAIL_SENT                           = "email_sent"
	ORDER_EVENT_TYPE__PAYMENT_AUTHORIZED                   = "payment_authorized"
	ORDER_EVENT_TYPE__PAYMENT_CAPTURED                     = "payment_captured"
	ORDER_EVENT_TYPE__PAYMENT_REFUNDED                     = "payment_refunded"
	ORDER_EVENT_TYPE__PAYMENT_VOIDED                       = "payment_voided"
	ORDER_EVENT_TYPE__PAYMENT_FAILED                       = "payment_failed"
	ORDER_EVENT_TYPE__EXTERNAL_SERVICE_NOTIFICATION        = "external_service_notification"
	ORDER_EVENT_TYPE__INVOICE_REQUESTED                    = "invoice_requested"
	ORDER_EVENT_TYPE__INVOICE_GENERATED                    = "invoice_generated"
	ORDER_EVENT_TYPE__INVOICE_UPDATED                      = "invoice_updated"
	ORDER_EVENT_TYPE__INVOICE_SENT                         = "invoice_sent"
	ORDER_EVENT_TYPE__FULFILLMENT_CANCELED                 = "fulfillment_canceled"
	ORDER_EVENT_TYPE__FULFILLMENT_RESTOCKED_ITEMS          = "fulfillment_restocked_items"
	ORDER_EVENT_TYPE__FULFILLMENT_FULFILLED_ITEMS          = "fulfillment_fulfilled_items"
	ORDER_EVENT_TYPE__FULFILLMENT_REFUNDED                 = "fulfillment_refunded"
	ORDER_EVENT_TYPE__FULFILLMENT_RETURNED                 = "fulfillment_returned"
	ORDER_EVENT_TYPE__FULFILLMENT_REPLACED                 = "fulfillment_replaced"
	ORDER_EVENT_TYPE__TRACKING_UPDATED                     = "tracking_updated"
	ORDER_EVENT_TYPE__NOTE_ADDED                           = "note_added"
	ORDER_EVENT_TYPE__OTHER                                = "other" // Used mostly for importing legacy data from before Enum-based events
)

var (
	OrderEventTypeStrings map[string]string
)

func init() {
	OrderEventTypeStrings = map[string]string{
		ORDER_EVENT_TYPE__DRAFT_CREATED:                        "The draft order was created",
		ORDER_EVENT_TYPE__DRAFT_CREATED_FROM_REPLACE:           "The draft order with replace lines was created",
		ORDER_EVENT_TYPE__ADDED_PRODUCTS:                       "Some products were added to the order",
		ORDER_EVENT_TYPE__REMOVED_PRODUCTS:                     "Some products were removed from the order",
		ORDER_EVENT_TYPE__PLACED:                               "The order was placed",
		ORDER_EVENT_TYPE__PLACED_FROM_DRAFT:                    "The draft order was placed",
		ORDER_EVENT_TYPE__OVERSOLD_ITEMS:                       "The draft order was placed with oversold items",
		ORDER_EVENT_TYPE__CANCELED:                             "The order was canceled",
		ORDER_EVENT_TYPE__ORDER_MARKED_AS_PAID:                 "The order was manually marked as fully paid",
		ORDER_EVENT_TYPE__ORDER_FULLY_PAID:                     "The order was fully paid",
		ORDER_EVENT_TYPE__ORDER_REPLACEMENT_CREATED:            "The draft order was created based on this order.",
		ORDER_EVENT_TYPE__ORDER_DISCOUNT_ADDED:                 "New order discount applied to this order.",
		ORDER_EVENT_TYPE__ORDER_DISCOUNT_AUTOMATICALLY_UPDATED: "Order discount was automatically updated after the changes in order.",
		ORDER_EVENT_TYPE__ORDER_DISCOUNT_UPDATED:               "Order discount was updated for this order.",
		ORDER_EVENT_TYPE__ORDER_DISCOUNT_DELETED:               "Order discount was deleted for this order.",
		ORDER_EVENT_TYPE__ORDER_LINE_DISCOUNT_UPDATED:          "Order line was discounted.",
		ORDER_EVENT_TYPE__ORDER_LINE_DISCOUNT_REMOVED:          "The discount for order line was removed.",
		ORDER_EVENT_TYPE__UPDATED_ADDRESS:                      "The address from the placed order was updated",
		ORDER_EVENT_TYPE__EMAIL_SENT:                           "The email was sent",
		ORDER_EVENT_TYPE__CONFIRMED:                            "Order was confirmed",
		ORDER_EVENT_TYPE__PAYMENT_AUTHORIZED:                   "The payment was authorized",
		ORDER_EVENT_TYPE__PAYMENT_CAPTURED:                     "The payment was captured",
		ORDER_EVENT_TYPE__EXTERNAL_SERVICE_NOTIFICATION:        "Notification from external service",
		ORDER_EVENT_TYPE__PAYMENT_REFUNDED:                     "The payment was refunded",
		ORDER_EVENT_TYPE__PAYMENT_VOIDED:                       "The payment was voided",
		ORDER_EVENT_TYPE__PAYMENT_FAILED:                       "The payment was failed",
		ORDER_EVENT_TYPE__INVOICE_REQUESTED:                    "An invoice was requested",
		ORDER_EVENT_TYPE__INVOICE_GENERATED:                    "An invoice was generated",
		ORDER_EVENT_TYPE__INVOICE_UPDATED:                      "An invoice was updated",
		ORDER_EVENT_TYPE__INVOICE_SENT:                         "An invoice was sent",
		ORDER_EVENT_TYPE__FULFILLMENT_CANCELED:                 "A fulfillment was canceled",
		ORDER_EVENT_TYPE__FULFILLMENT_RESTOCKED_ITEMS:          "The items of the fulfillment were restocked",
		ORDER_EVENT_TYPE__FULFILLMENT_FULFILLED_ITEMS:          "Some items were fulfilled",
		ORDER_EVENT_TYPE__FULFILLMENT_REFUNDED:                 "Some items were refunded",
		ORDER_EVENT_TYPE__FULFILLMENT_RETURNED:                 "Some items were returned",
		ORDER_EVENT_TYPE__FULFILLMENT_REPLACED:                 "Some items were replaced",
		ORDER_EVENT_TYPE__TRACKING_UPDATED:                     "The fulfillment's tracking code was updated",
		ORDER_EVENT_TYPE__NOTE_ADDED:                           "A note was added to the order",
		ORDER_EVENT_TYPE__OTHER:                                "An unknown order event containing a message",
	}
}

// max lengths for some order event's type
const (
	ORDER_EVENT_TYPE_MAX_LENGTH = 255
)

// Model used to store events that happened during the order lifecycle.
type OrderEvent struct {
	Id         string                 `json:"id"`
	CreateAt   int64                  `json:"create_at"`
	Type       string                 `json:"type"`
	OrderID    string                 `json:"order_id"`
	Parameters *model.StringInterface `json:"parameters"`
	UserID     *string                `json:"user_id"`
}

// OrderEventOption contains parameters to create new order event instance
type OrderEventOption struct {
	OrderID    string
	Parameters *model.StringInterface
	Type       string
	UserID     *string
}

func (o *OrderEvent) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.order_event.is_valid.%s.app_error",
		"order_event_id=",
		"OrderEvent.IsValid",
	)
	if !model.IsValidId(o.Id) {
		return outer("id", nil)
	}
	if o.CreateAt == 0 {
		return outer("create_st", &o.Id)
	}
	if o.UserID != nil && !model.IsValidId(*o.UserID) {
		return outer("user_id", &o.Id)
	}
	if !model.IsValidId(o.OrderID) {
		return outer("order_id", &o.Id)
	}
	if len(o.Type) > ORDER_EVENT_TYPE_MAX_LENGTH || OrderEventTypeStrings[strings.ToLower(o.Type)] == "" {
		return outer("type", &o.Id)
	}

	return nil
}

func (o *OrderEvent) ToJson() string {
	return model.ModelToJson(o)
}

func OrderEventFromJson(data io.Reader) *OrderEvent {
	var o OrderEvent
	model.ModelFromJson(&o, data)
	return &o
}

func (o *OrderEvent) PreSave() {
	if o.Id == "" {
		o.Id = model.NewId()
	}
	o.CreateAt = model.GetMillis()
}

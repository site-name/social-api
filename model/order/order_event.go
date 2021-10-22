package order

import (
	"github.com/sitename/sitename/model"
)

// OrderEvents represents type of order event
type OrderEvents string

// valid values for order event's type
const (
	CONFIRMED                            OrderEvents = "confirmed"
	DRAFT_CREATED                        OrderEvents = "draft_created"
	DRAFT_CREATED_FROM_REPLACE           OrderEvents = "draft_created_from_replace"
	ADDED_PRODUCTS                       OrderEvents = "added_products"
	REMOVED_PRODUCTS                     OrderEvents = "removed_products"
	PLACED                               OrderEvents = "placed"
	PLACED_FROM_DRAFT                    OrderEvents = "placed_from_draft"
	OVERSOLD_ITEMS                       OrderEvents = "oversold_items"
	CANCELED_                            OrderEvents = "canceled"
	ORDER_MARKED_AS_PAID                 OrderEvents = "order_marked_as_paid"
	ORDER_FULLY_PAID                     OrderEvents = "order_fully_paid"
	ORDER_REPLACEMENT_CREATED            OrderEvents = "order_replacement_created"
	ORDER_DISCOUNT_ADDED                 OrderEvents = "order_discount_added"
	ORDER_DISCOUNT_AUTOMATICALLY_UPDATED OrderEvents = "order_discount_automatically_updated"
	ORDER_DISCOUNT_UPDATED               OrderEvents = "order_discount_updated"
	ORDER_DISCOUNT_DELETED               OrderEvents = "order_discount_deleted"
	ORDER_LINE_DISCOUNT_UPDATED          OrderEvents = "order_line_discount_updated"
	ORDER_LINE_DISCOUNT_REMOVED          OrderEvents = "order_line_discount_removed"
	UPDATED_ADDRESS                      OrderEvents = "updated_address"
	EMAIL_SENT                           OrderEvents = "email_sent"
	PAYMENT_AUTHORIZED                   OrderEvents = "payment_authorized"
	PAYMENT_CAPTURED                     OrderEvents = "payment_captured"
	PAYMENT_REFUNDED                     OrderEvents = "payment_refunded"
	PAYMENT_VOIDED                       OrderEvents = "payment_voided"
	PAYMENT_FAILED                       OrderEvents = "payment_failed"
	EXTERNAL_SERVICE_NOTIFICATION        OrderEvents = "external_service_notification"
	INVOICE_REQUESTED                    OrderEvents = "invoice_requested"
	INVOICE_GENERATED                    OrderEvents = "invoice_generated"
	INVOICE_UPDATED                      OrderEvents = "invoice_updated"
	INVOICE_SENT                         OrderEvents = "invoice_sent"
	FULFILLMENT_CANCELED_                OrderEvents = "fulfillment_canceled"
	FULFILLMENT_RESTOCKED_ITEMS          OrderEvents = "fulfillment_restocked_items"
	FULFILLMENT_FULFILLED_ITEMS          OrderEvents = "fulfillment_fulfilled_items"
	FULFILLMENT_REFUNDED_                OrderEvents = "fulfillment_refunded"
	FULFILLMENT_RETURNED_                OrderEvents = "fulfillment_returned"
	FULFILLMENT_REPLACED_                OrderEvents = "fulfillment_replaced"
	FULFILLMENT_AWAITS_APPROVAL          OrderEvents = "fulfillment_awaits_approval"
	TRACKING_UPDATED                     OrderEvents = "tracking_updated"
	NOTE_ADDED                           OrderEvents = "note_added"
	OTHER                                OrderEvents = "other" // Used mostly for importing legacy data from before Enum-based events
)

var OrderEventTypeStrings = map[OrderEvents]string{
	DRAFT_CREATED:                        "The draft order was created",
	DRAFT_CREATED_FROM_REPLACE:           "The draft order with replace lines was created",
	ADDED_PRODUCTS:                       "Some products were added to the order",
	REMOVED_PRODUCTS:                     "Some products were removed from the order",
	PLACED:                               "The order was placed",
	PLACED_FROM_DRAFT:                    "The draft order was placed",
	OVERSOLD_ITEMS:                       "The draft order was placed with oversold items",
	CANCELED_:                            "The order was canceled",
	ORDER_MARKED_AS_PAID:                 "The order was manually marked as fully paid",
	ORDER_FULLY_PAID:                     "The order was fully paid",
	ORDER_REPLACEMENT_CREATED:            "The draft order was created based on this order.",
	ORDER_DISCOUNT_ADDED:                 "New order discount applied to this order.",
	ORDER_DISCOUNT_AUTOMATICALLY_UPDATED: "Order discount was automatically updated after the changes in order.",
	ORDER_DISCOUNT_UPDATED:               "Order discount was updated for this order.",
	ORDER_DISCOUNT_DELETED:               "Order discount was deleted for this order.",
	ORDER_LINE_DISCOUNT_UPDATED:          "Order line was discounted.",
	ORDER_LINE_DISCOUNT_REMOVED:          "The discount for order line was removed.",
	UPDATED_ADDRESS:                      "The address from the placed order was updated",
	EMAIL_SENT:                           "The email was sent",
	CONFIRMED:                            "Order was confirmed",
	PAYMENT_AUTHORIZED:                   "The payment was authorized",
	PAYMENT_CAPTURED:                     "The payment was captured",
	EXTERNAL_SERVICE_NOTIFICATION:        "Notification from external service",
	PAYMENT_REFUNDED:                     "The payment was refunded",
	PAYMENT_VOIDED:                       "The payment was voided",
	PAYMENT_FAILED:                       "The payment was failed",
	INVOICE_REQUESTED:                    "An invoice was requested",
	INVOICE_GENERATED:                    "An invoice was generated",
	INVOICE_UPDATED:                      "An invoice was updated",
	INVOICE_SENT:                         "An invoice was sent",
	FULFILLMENT_CANCELED_:                "A fulfillment was canceled",
	FULFILLMENT_RESTOCKED_ITEMS:          "The items of the fulfillment were restocked",
	FULFILLMENT_FULFILLED_ITEMS:          "Some items were fulfilled",
	FULFILLMENT_REFUNDED_:                "Some items were refunded",
	FULFILLMENT_RETURNED_:                "Some items were returned",
	FULFILLMENT_REPLACED_:                "Some items were replaced",
	FULFILLMENT_AWAITS_APPROVAL:          "Fulfillments awaits approval",
	TRACKING_UPDATED:                     "The fulfillment's tracking code was updated",
	NOTE_ADDED:                           "A note was added to the order",
	OTHER:                                "An unknown order event containing a message",
}

// max lengths for some order event's type
const (
	ORDER_EVENT_TYPE_MAX_LENGTH = 255
)

// Model used to store events that happened during the order lifecycle.
type OrderEvent struct {
	Id         string                `json:"id"`
	CreateAt   int64                 `json:"create_at"`
	Type       OrderEvents           `json:"type"`
	OrderID    string                `json:"order_id"`
	Parameters model.StringInterface `json:"parameters"`
	UserID     *string               `json:"user_id"`
}

// OrderEventOption contains parameters to create new order event instance
type OrderEventOption struct {
	OrderID    string
	Parameters model.StringInterface
	Type       OrderEvents
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
	if len(o.Type) > ORDER_EVENT_TYPE_MAX_LENGTH || OrderEventTypeStrings[o.Type] == "" {
		return outer("type", &o.Id)
	}

	return nil
}

func (o *OrderEvent) ToJSON() string {
	return model.ModelToJson(o)
}

func (o *OrderEvent) PreSave() {
	if o.Id == "" {
		o.Id = model.NewId()
	}
	o.CreateAt = model.GetMillis()
}

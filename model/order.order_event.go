package model

import "github.com/Masterminds/squirrel"

// OrderEventType represents type of order event
type OrderEventType string

// valid values for order event's type
const (
	CONFIRMED                            OrderEventType = "confirmed"
	DRAFT_CREATED                        OrderEventType = "draft_created"
	DRAFT_CREATED_FROM_REPLACE           OrderEventType = "draft_created_from_replace"
	ADDED_PRODUCTS                       OrderEventType = "added_products"
	REMOVED_PRODUCTS                     OrderEventType = "removed_products"
	PLACED                               OrderEventType = "placed"
	PLACED_FROM_DRAFT                    OrderEventType = "placed_from_draft"
	OVERSOLD_ITEMS                       OrderEventType = "oversold_items"
	CANCELED_                            OrderEventType = "canceled"
	ORDER_MARKED_AS_PAID                 OrderEventType = "order_marked_as_paid"
	ORDER_FULLY_PAID                     OrderEventType = "order_fully_paid"
	ORDER_REPLACEMENT_CREATED            OrderEventType = "order_replacement_created"
	ORDER_DISCOUNT_ADDED                 OrderEventType = "order_discount_added"
	ORDER_DISCOUNT_AUTOMATICALLY_UPDATED OrderEventType = "order_discount_automatically_updated"
	ORDER_DISCOUNT_UPDATED               OrderEventType = "order_discount_updated"
	ORDER_DISCOUNT_DELETED               OrderEventType = "order_discount_deleted"
	ORDER_LINE_DISCOUNT_UPDATED          OrderEventType = "order_line_discount_updated"
	ORDER_LINE_DISCOUNT_REMOVED          OrderEventType = "order_line_discount_removed"
	UPDATED_ADDRESS                      OrderEventType = "updated_address"
	EMAIL_SENT                           OrderEventType = "email_sent"
	PAYMENT_AUTHORIZED                   OrderEventType = "payment_authorized"
	PAYMENT_CAPTURED                     OrderEventType = "payment_captured"
	PAYMENT_REFUNDED                     OrderEventType = "payment_refunded"
	PAYMENT_VOIDED                       OrderEventType = "payment_voided"
	PAYMENT_FAILED                       OrderEventType = "payment_failed"
	EXTERNAL_SERVICE_NOTIFICATION        OrderEventType = "external_service_notification"
	INVOICE_REQUESTED                    OrderEventType = "invoice_requested"
	INVOICE_GENERATED                    OrderEventType = "invoice_generated"
	INVOICE_UPDATED                      OrderEventType = "invoice_updated"
	INVOICE_SENT                         OrderEventType = "invoice_sent"
	FULFILLMENT_CANCELED_                OrderEventType = "fulfillment_canceled"
	FULFILLMENT_RESTOCKED_ITEMS          OrderEventType = "fulfillment_restocked_items"
	FULFILLMENT_FULFILLED_ITEMS          OrderEventType = "fulfillment_fulfilled_items"
	FULFILLMENT_REFUNDED_                OrderEventType = "fulfillment_refunded"
	FULFILLMENT_RETURNED_                OrderEventType = "fulfillment_returned"
	FULFILLMENT_REPLACED_                OrderEventType = "fulfillment_replaced"
	FULFILLMENT_AWAITS_APPROVAL          OrderEventType = "fulfillment_awaits_approval"
	TRACKING_UPDATED                     OrderEventType = "tracking_updated"
	NOTE_ADDED                           OrderEventType = "note_added"
	OTHER                                OrderEventType = "other" // Used mostly for importing legacy data from before Enum-based events
)

var OrderEventTypeStrings = map[OrderEventType]string{
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
	Id       string         `json:"id"`
	CreateAt int64          `json:"create_at"`
	Type     OrderEventType `json:"type"`
	OrderID  string         `json:"order_id"`
	// To reduce number of type casting, checking steps, below are
	// possible keys and their according values TYPES you must follow when storing things into this field:
	//  "email": string
	//  "email_type": string
	//  "amount": float64
	//  "payment_id": string
	//  "quantity": int
	//  "message": string
	//  "composedID": string
	//  "oversold_items": []string
	//  "invoice_number": string
	//  "transaction_reference": string
	//  "shipping_costs_included": bool
	//  "related_order_pk": string
	//  "warehouse": string
	//  "fulfilled_items": []string
	Parameters StringInterface `json:"parameters"`
	UserID     *string         `json:"user_id"`
}

// OrderEventOption contains parameters to create new order event instance
type OrderEventOption struct {
	OrderID    string
	Parameters StringInterface // should contains keys in ["invoice_number"]
	Type       OrderEventType
	UserID     *string
}

type OrderEventFilterOptions struct {
	Id      squirrel.Sqlizer
	Type    squirrel.Sqlizer
	OrderID squirrel.Sqlizer
}

func (o *OrderEvent) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"model.order_event.is_valid.%s.app_error",
		"order_event_id=",
		"OrderEvent.IsValid",
	)
	if !IsValidId(o.Id) {
		return outer("id", nil)
	}
	if o.CreateAt == 0 {
		return outer("create_st", &o.Id)
	}
	if o.UserID != nil && !IsValidId(*o.UserID) {
		return outer("user_id", &o.Id)
	}
	if !IsValidId(o.OrderID) {
		return outer("order_id", &o.Id)
	}
	if len(o.Type) > ORDER_EVENT_TYPE_MAX_LENGTH || OrderEventTypeStrings[o.Type] == "" {
		return outer("type", &o.Id)
	}

	return nil
}

func (o *OrderEvent) ToJSON() string {
	return ModelToJson(o)
}

func (o *OrderEvent) PreSave() {
	if o.Id == "" {
		o.Id = NewId()
	}
	o.CreateAt = GetMillis()
}

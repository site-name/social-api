package model

import (
	"net/http"

	"github.com/Masterminds/squirrel"
	"gorm.io/gorm"
)

// OrderEventType represents type of order event
type OrderEventType string

func (o OrderEventType) IsValid() bool {
	return OrderEventTypeStrings[o] != ""
}

// valid values for order event's type
const (
	ORDER_EVENT_TYPE_CONFIRMED                            OrderEventType = "confirmed"
	ORDER_EVENT_TYPE_DRAFT_CREATED                        OrderEventType = "draft_created"
	ORDER_EVENT_TYPE_DRAFT_CREATED_FROM_REPLACE           OrderEventType = "draft_created_from_replace"
	ORDER_EVENT_TYPE_ADDED_PRODUCTS                       OrderEventType = "added_products"
	ORDER_EVENT_TYPE_REMOVED_PRODUCTS                     OrderEventType = "removed_products"
	ORDER_EVENT_TYPE_PLACED                               OrderEventType = "placed"
	ORDER_EVENT_TYPE_PLACED_FROM_DRAFT                    OrderEventType = "placed_from_draft"
	ORDER_EVENT_TYPE_OVERSOLD_ITEMS                       OrderEventType = "oversold_items"
	ORDER_EVENT_TYPE_CANCELED                             OrderEventType = "canceled"
	ORDER_EVENT_TYPE_ORDER_MARKED_AS_PAID                 OrderEventType = "order_marked_as_paid"
	ORDER_EVENT_TYPE_ORDER_FULLY_PAID                     OrderEventType = "order_fully_paid"
	ORDER_EVENT_TYPE_ORDER_REPLACEMENT_CREATED            OrderEventType = "order_replacement_created"
	ORDER_EVENT_TYPE_ORDER_DISCOUNT_ADDED                 OrderEventType = "order_discount_added"
	ORDER_EVENT_TYPE_ORDER_DISCOUNT_AUTOMATICALLY_UPDATED OrderEventType = "order_discount_automatically_updated"
	ORDER_EVENT_TYPE_ORDER_DISCOUNT_UPDATED               OrderEventType = "order_discount_updated"
	ORDER_EVENT_TYPE_ORDER_DISCOUNT_DELETED               OrderEventType = "order_discount_deleted"
	ORDER_EVENT_TYPE_ORDER_LINE_DISCOUNT_UPDATED          OrderEventType = "order_line_discount_updated"
	ORDER_EVENT_TYPE_ORDER_LINE_DISCOUNT_REMOVED          OrderEventType = "order_line_discount_removed"
	ORDER_EVENT_TYPE_UPDATED_ADDRESS                      OrderEventType = "updated_address"
	ORDER_EVENT_TYPE_EMAIL_SENT                           OrderEventType = "email_sent"
	ORDER_EVENT_TYPE_PAYMENT_AUTHORIZED                   OrderEventType = "payment_authorized"
	ORDER_EVENT_TYPE_PAYMENT_CAPTURED                     OrderEventType = "payment_captured"
	ORDER_EVENT_TYPE_PAYMENT_REFUNDED                     OrderEventType = "payment_refunded"
	ORDER_EVENT_TYPE_PAYMENT_VOIDED                       OrderEventType = "payment_voided"
	ORDER_EVENT_TYPE_PAYMENT_FAILED                       OrderEventType = "payment_failed"
	ORDER_EVENT_TYPE_EXTERNAL_SERVICE_NOTIFICATION        OrderEventType = "external_service_notification"
	ORDER_EVENT_TYPE_INVOICE_REQUESTED                    OrderEventType = "invoice_requested"
	ORDER_EVENT_TYPE_INVOICE_GENERATED                    OrderEventType = "invoice_generated"
	ORDER_EVENT_TYPE_INVOICE_UPDATED                      OrderEventType = "invoice_updated"
	ORDER_EVENT_TYPE_INVOICE_SENT                         OrderEventType = "invoice_sent"
	ORDER_EVENT_TYPE_FULFILLMENT_CANCELED                 OrderEventType = "fulfillment_canceled"
	ORDER_EVENT_TYPE_FULFILLMENT_RESTOCKED_ITEMS          OrderEventType = "fulfillment_restocked_items"
	ORDER_EVENT_TYPE_FULFILLMENT_FULFILLED_ITEMS          OrderEventType = "fulfillment_fulfilled_items"
	ORDER_EVENT_TYPE_FULFILLMENT_REFUNDED                 OrderEventType = "fulfillment_refunded"
	ORDER_EVENT_TYPE_FULFILLMENT_RETURNED                 OrderEventType = "fulfillment_returned"
	ORDER_EVENT_TYPE_FULFILLMENT_REPLACED                 OrderEventType = "fulfillment_replaced"
	ORDER_EVENT_TYPE_FULFILLMENT_AWAITS_APPROVAL          OrderEventType = "fulfillment_awaits_approval"
	ORDER_EVENT_TYPE_TRACKING_UPDATED                     OrderEventType = "tracking_updated"
	ORDER_EVENT_TYPE_NOTE_ADDED                           OrderEventType = "note_added"
	ORDER_EVENT_TYPE_OTHER                                OrderEventType = "other" // Used mostly for importing legacy data from before Enum-based events
)

var OrderEventTypeStrings = map[OrderEventType]string{
	ORDER_EVENT_TYPE_DRAFT_CREATED:                        "The draft order was created",
	ORDER_EVENT_TYPE_DRAFT_CREATED_FROM_REPLACE:           "The draft order with replace lines was created",
	ORDER_EVENT_TYPE_ADDED_PRODUCTS:                       "Some products were added to the order",
	ORDER_EVENT_TYPE_REMOVED_PRODUCTS:                     "Some products were removed from the order",
	ORDER_EVENT_TYPE_PLACED:                               "The order was placed",
	ORDER_EVENT_TYPE_PLACED_FROM_DRAFT:                    "The draft order was placed",
	ORDER_EVENT_TYPE_OVERSOLD_ITEMS:                       "The draft order was placed with oversold items",
	ORDER_EVENT_TYPE_CANCELED:                             "The order was canceled",
	ORDER_EVENT_TYPE_ORDER_MARKED_AS_PAID:                 "The order was manually marked as fully paid",
	ORDER_EVENT_TYPE_ORDER_FULLY_PAID:                     "The order was fully paid",
	ORDER_EVENT_TYPE_ORDER_REPLACEMENT_CREATED:            "The draft order was created based on this order.",
	ORDER_EVENT_TYPE_ORDER_DISCOUNT_ADDED:                 "New order discount applied to this order.",
	ORDER_EVENT_TYPE_ORDER_DISCOUNT_AUTOMATICALLY_UPDATED: "Order discount was automatically updated after the changes in order.",
	ORDER_EVENT_TYPE_ORDER_DISCOUNT_UPDATED:               "Order discount was updated for this order.",
	ORDER_EVENT_TYPE_ORDER_DISCOUNT_DELETED:               "Order discount was deleted for this order.",
	ORDER_EVENT_TYPE_ORDER_LINE_DISCOUNT_UPDATED:          "Order line was discounted.",
	ORDER_EVENT_TYPE_ORDER_LINE_DISCOUNT_REMOVED:          "The discount for order line was removed.",
	ORDER_EVENT_TYPE_UPDATED_ADDRESS:                      "The address from the placed order was updated",
	ORDER_EVENT_TYPE_EMAIL_SENT:                           "The email was sent",
	ORDER_EVENT_TYPE_CONFIRMED:                            "Order was confirmed",
	ORDER_EVENT_TYPE_PAYMENT_AUTHORIZED:                   "The payment was authorized",
	ORDER_EVENT_TYPE_PAYMENT_CAPTURED:                     "The payment was captured",
	ORDER_EVENT_TYPE_EXTERNAL_SERVICE_NOTIFICATION:        "Notification from external service",
	ORDER_EVENT_TYPE_PAYMENT_REFUNDED:                     "The payment was refunded",
	ORDER_EVENT_TYPE_PAYMENT_VOIDED:                       "The payment was voided",
	ORDER_EVENT_TYPE_PAYMENT_FAILED:                       "The payment was failed",
	ORDER_EVENT_TYPE_INVOICE_REQUESTED:                    "An invoice was requested",
	ORDER_EVENT_TYPE_INVOICE_GENERATED:                    "An invoice was generated",
	ORDER_EVENT_TYPE_INVOICE_UPDATED:                      "An invoice was updated",
	ORDER_EVENT_TYPE_INVOICE_SENT:                         "An invoice was sent",
	ORDER_EVENT_TYPE_FULFILLMENT_CANCELED:                 "A fulfillment was canceled",
	ORDER_EVENT_TYPE_FULFILLMENT_RESTOCKED_ITEMS:          "The items of the fulfillment were restocked",
	ORDER_EVENT_TYPE_FULFILLMENT_FULFILLED_ITEMS:          "Some items were fulfilled",
	ORDER_EVENT_TYPE_FULFILLMENT_REFUNDED:                 "Some items were refunded",
	ORDER_EVENT_TYPE_FULFILLMENT_RETURNED:                 "Some items were returned",
	ORDER_EVENT_TYPE_FULFILLMENT_REPLACED:                 "Some items were replaced",
	ORDER_EVENT_TYPE_FULFILLMENT_AWAITS_APPROVAL:          "Fulfillments awaits approval",
	ORDER_EVENT_TYPE_TRACKING_UPDATED:                     "The fulfillment's tracking code was updated",
	ORDER_EVENT_TYPE_NOTE_ADDED:                           "A note was added to the order",
	ORDER_EVENT_TYPE_OTHER:                                "An unknown order event containing a message",
}

// Model used to store events that happened during the order lifecycle.
type OrderEvent struct {
	Id       string         `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	CreateAt int64          `json:"create_at" gorm:"type:bigint;column:CreateAt;autoCreateTime:milli"`
	Type     OrderEventType `json:"type" gorm:"type:varchar(255);column:Type"`
	OrderID  string         `json:"order_id" gorm:"type:uuid;column:OrderID"`
	// To reduce number of type assertion steps, below are
	// possible keys and their according values TYPES you must follow when storing things into this field:
	//  "email": string
	//  "email_type": string
	//  "amount": float64
	//  "payment_id": string
	//  "payment_gateway": string
	//  "quantity": int
	//  "message": string
	//  "composed_id": string
	//  "oversold_items": []string
	//  "invoice_number": string
	//  "transaction_reference": string
	//  "shipping_costs_included": bool
	//  "related_order_pk": string
	//  "warehouse": string
	//  "fulfilled_items": []string // ids of fulfillment lines
	//  "lines": []map[string]any{
	//      "quantity": int,
	//      "line_pk": string,
	//      "item": string,
	//      "discount": map[string]any{ // NOTE: Remember to check nil
	//        "value": float64,
	//        "amount_value": float64,
	//        "currency": string,
	//        "value_type": string,
	//        "reason": string,
	//        "old_value": float64,
	//        "old_value_type": string,
	//        "old_amount_value": float64,
	//      }
	//    }
	//  "url": string
	//  "status": string
	//  "gateway": string
	//  "awaiting_fulfillments": []string // ids of fulfillment lines
	//  "tracking_number": string
	//  "fulfillment": string
	//  "discount": map[string]any{
	//    "value": string,
	//    "amount_value": float64,
	//    "currency": string,
	//    "value_type": string,
	//    "reason": string,
	//    "old_value": float64,
	//    "old_value_type": string,
	//    "old_amount_value": float64,
	//  }
	Parameters StringInterface `json:"parameters" gorm:"type:jsonb;column:Parameters"`
	UserID     *string         `json:"user_id" gorm:"type:uuid;column:UserID"`
}

func (c *OrderEvent) BeforeCreate(_ *gorm.DB) error { return c.IsValid() }
func (c *OrderEvent) BeforeUpdate(_ *gorm.DB) error { return c.IsValid() }
func (c *OrderEvent) TableName() string             { return OrderEventTableName }

// OrderEventOption contains parameters to create new order event instance
type OrderEventOption struct {
	OrderID string
	// To reduce number of type assertion steps, below are
	// possible keys and their according values TYPES you must follow when storing things into this field:
	//  "email": string
	//  "email_type": string
	//  "amount": float64
	//  "payment_id": string
	//  "payment_gateway": string
	//  "quantity": int
	//  "message": string
	//  "composed_id": string
	//  "oversold_items": []string
	//  "invoice_number": string
	//  "transaction_reference": string
	//  "shipping_costs_included": bool
	//  "related_order_pk": string
	//  "warehouse": string
	//  "fulfilled_items": []string // ids of fulfillment lines
	//  "lines": []map[string]any{
	//      "quantity": int,
	//      "line_pk": string,
	//      "item": string,
	//      "discount": map[string]any{ // NOTE: Remember to check nil
	//        "value": float64,
	//        "amount_value": float64,
	//        "currency": string,
	//        "value_type": string,
	//        "reason": string,
	//        "old_value": float64,
	//        "old_value_type": string,
	//        "old_amount_value": float64,
	//      }
	//    }
	//  "url": string
	//  "status": string
	//  "gateway": string
	//  "awaiting_fulfillments": []string // ids of fulfillment lines
	//  "tracking_number": string
	//  "fulfillment": string
	//  "discount": map[string]any{
	//    "value": string,
	//    "amount_value": float64,
	//    "currency": string,
	//    "value_type": string,
	//    "reason": string,
	//    "old_value": float64,
	//    "old_value_type": string,
	//    "old_amount_value": float64,
	//  }
	Parameters StringInterface
	Type       OrderEventType
	UserID     *string
}

type OrderEventFilterOptions struct {
	Conditions squirrel.Sqlizer
}

func (o *OrderEvent) IsValid() *AppError {
	if o.UserID != nil && !IsValidId(*o.UserID) {
		return NewAppError("OrderEvent.IsValid", "model.order_event.is_valid.user_id.app_error", nil, "please provide valid user id", http.StatusBadRequest)
	}
	if !IsValidId(o.OrderID) {
		return NewAppError("OrderEvent.IsValid", "model.order_event.is_valid.order_id.app_error", nil, "please provide valid order id", http.StatusBadRequest)
	}
	if !o.Type.IsValid() {
		return NewAppError("OrderEvent.IsValid", "model.order_event.is_valid.type.app_error", nil, "please provide valid type", http.StatusBadRequest)
	}

	return nil
}

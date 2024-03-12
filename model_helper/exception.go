package model_helper

import (
	"net/http"
	"strings"

	"github.com/sitename/sitename/model"
)

type CheckoutErrorCode string

// checkout error codes
const (
	BILLING_ADDRESS_NOT_SET_         CheckoutErrorCode = "billing_address_not_set"
	CHECKOUT_NOT_FULLY_PAID          CheckoutErrorCode = "checkout_not_fully_paid"
	GRAPHQL_ERROR__                  CheckoutErrorCode = "graphql_error"
	PRODUCT_NOT_PUBLISHED            CheckoutErrorCode = "product_not_published"
	PRODUCT_UNAVAILABLE_FOR_PURCHASE CheckoutErrorCode = "product_unavailable_for_purchase"
	INSUFFICIENT_STOCK               CheckoutErrorCode = "insufficient_stock"
	INVALID__                        CheckoutErrorCode = "invalid"
	INVALID_SHIPPING_METHOD_         CheckoutErrorCode = "invalid_shipping_method"
	NOT_FOUND__                      CheckoutErrorCode = "not_found"
	PAYMENT_ERROR_                   CheckoutErrorCode = "payment_error"
	QUANTITY_GREATER_THAN_LIMIT      CheckoutErrorCode = "quantity_greater_than_limit"
	REQUIRED__                       CheckoutErrorCode = "required"
	SHIPPING_ADDRESS_NOT_SET_        CheckoutErrorCode = "shipping_address_not_set"
	SHIPPING_METHOD_NOT_APPLICABLE   CheckoutErrorCode = "shipping_method_not_applicable"
	DELIVERY_METHOD_NOT_APPLICABLE   CheckoutErrorCode = "delivery_method_not_applicable"
	SHIPPING_METHOD_NOT_SET_         CheckoutErrorCode = "shipping_method_not_set"
	SHIPPING_NOT_REQUIRED            CheckoutErrorCode = "shipping_not_required"
	TAX_ERROR                        CheckoutErrorCode = "tax_error"
	UNIQUE__                         CheckoutErrorCode = "unique"
	VOUCHER_NOT_APPLICABLE           CheckoutErrorCode = "voucher_not_applicable"
	GIFT_CARD_NOT_APPLICABLE         CheckoutErrorCode = "gift_card_not_applicable"
	ZERO_QUANTITY                    CheckoutErrorCode = "zero_quantity"
	MISSING_CHANNEL_SLUG             CheckoutErrorCode = "missing_channel_slug"
	CHANNEL_INACTIVE_                CheckoutErrorCode = "channel_inactive"
	UNAVAILABLE_VARIANT_IN_CHANNEL   CheckoutErrorCode = "unavailable_variant_in_channel"
)

type TaxError struct {
	Where   string
	Message string
}

func (t *TaxError) Error() string {
	return t.Where + ": " + t.Message
}

// InsufficientStockData is an error type
type InsufficientStockData struct {
	Variant           model.ProductVariant // Product variant
	OrderLine         *model.OrderLine     // can be nil
	WarehouseID       *string              // can be nil
	AvailableQuantity *int                 // can be nil
}

// InsufficientStock is an error indicating stock is insufficient
type InsufficientStock struct {
	Items []*InsufficientStockData
	Code  CheckoutErrorCode
}

func (err *InsufficientStock) ToAppError(where string) *AppError {
	var warehouseIDs = make([]string, 0, len(err.Items))
	var orderLineIDs = make([]string, 0, len(err.Items))
	for _, item := range err.Items {
		if item.WarehouseID != nil {
			warehouseIDs = append(warehouseIDs, *item.WarehouseID)
		}
		if item.OrderLine != nil {
			orderLineIDs = append(orderLineIDs, item.OrderLine.ID)
		}
	}

	return NewAppError(where, "app.order.insufficient_stock.app_error", map[string]any{"orderLines": strings.Join(orderLineIDs, ", "), "warehouses": strings.Join(warehouseIDs, ", ")}, "insufficient product stock", http.StatusNotAcceptable)
}

func NewInsufficientStock(items []*InsufficientStockData) *InsufficientStock {
	return &InsufficientStock{
		Items: items,
		Code:  INSUFFICIENT_STOCK,
	}
}

func (i *InsufficientStock) VariantIDs() []string {
	res := []string{}
	for _, item := range i.Items {
		res = append(res, item.Variant.ID)
	}

	return res
}

func (i *InsufficientStock) Error() string {
	var res string
	for idx, item := range i.Items {
		if idx < len(i.Items)-1 {
			res += ProductVariantString(item.Variant) + ", "
			continue
		}

		res += ProductVariantString(item.Variant)
	}

	return res
}

type PreorderAllocationError struct {
	OrderLine *model.OrderLine
	Message   string
}

func (p *PreorderAllocationError) Error() string {
	return p.Message
}

func NewPreorderAllocationError(orderLine model.OrderLine) *PreorderAllocationError {
	return &PreorderAllocationError{
		OrderLine: &orderLine,
		Message:   "Unable to allocate in stock for line " + OrderLineString(orderLine),
	}
}

package exception

import (
	"strings"

	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/model/product_and_discount"
)

type CheckoutErrorCode string

// checkout error codes
const (
	BILLING_ADDRESS_NOT_SET          CheckoutErrorCode = "billing_address_not_set"
	CHECKOUT_NOT_FULLY_PAID          CheckoutErrorCode = "checkout_not_fully_paid"
	GRAPHQL_ERROR                    CheckoutErrorCode = "graphql_error"
	PRODUCT_NOT_PUBLISHED            CheckoutErrorCode = "product_not_published"
	PRODUCT_UNAVAILABLE_FOR_PURCHASE CheckoutErrorCode = "product_unavailable_for_purchase"
	INSUFFICIENT_STOCK               CheckoutErrorCode = "insufficient_stock"
	INVALID                          CheckoutErrorCode = "invalid"
	INVALID_SHIPPING_METHOD          CheckoutErrorCode = "invalid_shipping_method"
	NOT_FOUND                        CheckoutErrorCode = "not_found"
	PAYMENT_ERROR                    CheckoutErrorCode = "payment_error"
	QUANTITY_GREATER_THAN_LIMIT      CheckoutErrorCode = "quantity_greater_than_limit"
	REQUIRED                         CheckoutErrorCode = "required"
	SHIPPING_ADDRESS_NOT_SET         CheckoutErrorCode = "shipping_address_not_set"
	SHIPPING_METHOD_NOT_APPLICABLE   CheckoutErrorCode = "shipping_method_not_applicable"
	DELIVERY_METHOD_NOT_APPLICABLE   CheckoutErrorCode = "delivery_method_not_applicable"
	SHIPPING_METHOD_NOT_SET          CheckoutErrorCode = "shipping_method_not_set"
	SHIPPING_NOT_REQUIRED            CheckoutErrorCode = "shipping_not_required"
	TAX_ERROR                        CheckoutErrorCode = "tax_error"
	UNIQUE                           CheckoutErrorCode = "unique"
	VOUCHER_NOT_APPLICABLE           CheckoutErrorCode = "voucher_not_applicable"
	ZERO_QUANTITY                    CheckoutErrorCode = "zero_quantity"
	MISSING_CHANNEL_SLUG             CheckoutErrorCode = "missing_channel_slug"
	CHANNEL_INACTIVE                 CheckoutErrorCode = "channel_inactive"
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
	Variant           product_and_discount.ProductVariant // Product variant
	OrderLine         *order.OrderLine                    // can be nil
	WarehouseID       *string                             // can be nil
	AvailableQuantity *int                                // can be nil
}

// InsufficientStock is an error indicating stock is insufficient
type InsufficientStock struct {
	Items []*InsufficientStockData
	Code  CheckoutErrorCode
}

func NewInsufficientStock(items []*InsufficientStockData) *InsufficientStock {
	return &InsufficientStock{
		Items: items,
		Code:  INSUFFICIENT_STOCK,
	}
}

func (i *InsufficientStock) Error() string {
	var builder strings.Builder

	builder.WriteString("Insufficient stock for ")
	for idx, item := range i.Items {
		builder.WriteString(item.Variant.String())
		if idx == 0 {
			continue
		}
		if idx == len(i.Items)-1 {
			break
		}
		builder.WriteString(", ")
	}

	return builder.String()
}

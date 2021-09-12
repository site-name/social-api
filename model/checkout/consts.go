package checkout

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

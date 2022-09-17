package model

// plugin event types
const (
	ACCOUNT_CONFIRMATION           = "account_confirmation"
	ACCOUNT_PASSWORD_RESET         = "account_password_reset"
	ACCOUNT_CHANGE_EMAIL_REQUEST   = "account_change_email_request"
	ACCOUNT_CHANGE_EMAIL_CONFIRM   = "account_change_email_confirm"
	ACCOUNT_DELETE                 = "account_delete"
	ACCOUNT_SET_CUSTOMER_PASSWORD  = "account_set_customer_password"
	INVOICE_READY                  = "invoice_ready"
	ORDER_CONFIRMATION             = "order_confirmation"
	ORDER_CONFIRMED                = "order_confirmed"
	ORDER_FULFILLMENT_CONFIRMATION = "order_fulfillment_confirmation"
	ORDER_FULFILLMENT_UPDATE       = "order_fulfillment_update"
	ORDER_PAYMENT_CONFIRMATION     = "order_payment_confirmation"
	ORDER_CANCELED                 = "order_canceled"
	ORDER_REFUND_CONFIRMATION      = "order_refund_confirmation"
	SEND_GIFT_CARD                 = "send_gift_card"
)

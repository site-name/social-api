package sub_app_iface

import (
	"context"

	"github.com/shopspring/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/channel"
	"github.com/sitename/sitename/model/checkout"
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/model/payment"
)

// GiftCardApp defines methods for giftcard app
type GiftcardApp interface {
	Save(id string) error
}

// PaymentApp defines methods for payment app
type PaymentApp interface {
	GetAllPaymentsByOrderId(orderID string) ([]*payment.Payment, *model.AppError)                // GetAllPaymentsByOrderId returns all payments that belong to order with given orderID
	GetLastOrderPayment(orderID string) (*payment.Payment, *model.AppError)                      // GetLastOrderPayment get most recent payment made for given order
	GetAllPaymentTransactions(paymentID string) ([]*payment.PaymentTransaction, *model.AppError) // GetAllPaymentTransactions returns all transactions belong to given payment
	GetLastPaymentTransaction(paymentID string) (*payment.PaymentTransaction, *model.AppError)   // GetLastPaymentTransaction return most recent transaction made for given payment
	PaymentIsAuthorized(paymentID string) (bool, *model.AppError)                                // PaymentIsAuthorized checks if given payment is authorized
	PaymentGetAuthorizedAmount(pm *payment.Payment) (*goprices.Money, *model.AppError)           // PaymentGetAuthorizedAmount calculates authorized amount
	PaymentCanVoid(pm *payment.Payment) (bool, *model.AppError)                                  // PaymentCanVoid check if payment can void
	// Extract order information along with payment details. Returns information required to process payment and additional billing/shipping addresses for optional fraud-prevention mechanisms.
	CreatePaymentInformation(payment *payment.Payment, paymentToken *string, amount *decimal.Decimal, customerId *string, storeSource bool, additionalData map[string]string) (*payment.PaymentData, *model.AppError)
	GetAlreadyProcessedTransaction(payment *payment.Payment, gatewayResponse *payment.GatewayResponse) (*payment.PaymentTransaction, *model.AppError) // GetAlreadyProcessedTransaction returns most recent processed transaction made for given payment
	// CreatePayment creates new payment inside database with given data and returned it
	CreatePayment(gateway, currency, email, customerIpAddress, paymentToken, returnUrl, externalReference string, total decimal.Decimal, extraData map[string]string, checkOut *checkout.Checkout, orDer *order.Order) (*payment.Payment, *model.AppError)
}

// CheckoutApp
type CheckoutApp interface {
}

// CheckoutApp
type AccountApp interface {
	AddressById(id string) (*account.Address, *model.AppError)                                    // GetAddressById returns address with given id. If not found returns nil and concret error
	UserById(ctx context.Context, userID string) (*account.User, *model.AppError)                 // GetUserById get user from database with given userId
	CustomerEventsByUser(userID string) ([]*account.CustomerEvent, *model.AppError)               // CustomerEventsByUser returns all customer event(s) belong to given user
	AddressesByUserId(userID string) ([]*account.Address, *model.AppError)                        // AddressesByUserId returns list of address(es) (if found) that belong to given user
	UserSetDefaultAddress(userID, addressID, addressType string) (*account.User, *model.AppError) // UserSetDefaultAddress set given address to be default for given user
	AddressDeleteForUser(userID, addressID string) *model.AppError                                // AddressDeleteForUser delete given address from addresses list of given user
}

type ProductApp interface {
}

type WishlistApp interface {
}

type AttributeApp interface {
}

type InvoiceApp interface {
}

type ChannelApp interface {
	// GetChannelBySlug returns a channel (if found) from database with given slug
	GetChannelBySlug(slug string) (*channel.Channel, *model.AppError)
	// GetDefaultChannel get random channel that is active
	GetDefaultActiveChannel() (*channel.Channel, *model.AppError)
	// CleanChannel performs:
	//
	// 1) If given slug is not nil, try getting a channel with that slug.
	//   +) if found, check if channel is active
	//
	// 2) If given slug if nil, it try
	CleanChannel(channelSlug *string) (*channel.Channel, *model.AppError)
}

type WarehouseApp interface {
}

type DiscountApp interface {
}

type OrderApp interface {
	GetAllOrderLinesByOrderId(orderID string) ([]*order.OrderLine, *model.AppError) // GetAllOrderLinesByOrderId returns a slice of order lines that belong to given order
	// OrderShippingIsRequired checks if an order requires ship or not by:
	//
	// 1) Find all child order lines that belong to given order
	//
	// 2) iterates over resulting slice to check if at least one order line requires shipping
	OrderShippingIsRequired(orderID string) (bool, *model.AppError)
	OrderTotalQuantity(orderID string) (int, *model.AppError)                                   // OrderTotalQuantity return total quantity of given order
	UpdateOrderTotalPaid(orderID string) *model.AppError                                        // UpdateOrderTotalPaid update given order's total paid amount
	OrderIsPreAuthorized(orderID string) (bool, *model.AppError)                                // OrderIsPreAuthorized checks if order is pre-authorized
	OrderIsCaptured(orderID string) (bool, *model.AppError)                                     // OrderIsCaptured checks if given order is captured
	OrderSubTotal(orderID string, orderCurrency string) (*goprices.TaxedMoney, *model.AppError) // OrderSubTotal returns sum of TotalPrice of all order lines that belong to given order
	OrderCanCancel(ord *order.Order) (bool, *model.AppError)                                    // OrderCanCalcel checks if given order can be canceled
	OrderCanCapture(ord *order.Order, payment *payment.Payment) (bool, *model.AppError)         // OrderCanCapture checks if given order can capture.
	OrderCanVoid(ord *order.Order, payment *payment.Payment) (bool, *model.AppError)            // OrderCanVoid checks if given order can void
	OrderCanRefund(ord *order.Order, payments []*payment.Payment) (bool, *model.AppError)       // OrderCanRefund checks if order can refund
	CanMarkOrderAsPaid(ord *order.Order, payments []*payment.Payment) (bool, *model.AppError)   // CanMarkOrderAsPaid checks if given order can be marked as paid.
	OrderTotalAuthorized(ord *order.Order) (*goprices.Money, *model.AppError)                   // OrderTotalAuthorized returns order's total authorized amount
	GetOrderCountryCode(ord *order.Order) (string, *model.AppError)                             // GetOrderCountryCode is helper function, returns contry code of given order
	OrderLineById(id string) (*order.OrderLine, *model.AppError)                                // OrderLineById returns order line with id of given id
}

type MenuApp interface {
}

type AppApp interface {
}

type CsvApp interface {
}

type SiteApp interface {
}

type ShippingApp interface {
}

type WebhookApp interface {
}

type PageApp interface {
}

type SeoApp interface {
}

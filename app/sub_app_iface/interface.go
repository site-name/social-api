package sub_app_iface

import (
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/channel"
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/model/payment"
)

// GiftCardApp defines methods for giftcard app
type GiftcardApp interface {
	Save(id string) error
}

// PaymentApp defines methods for payment app
type PaymentApp interface {
	// GetAllPaymentsByOrderId returns all payments that belong to order with given orderID
	GetAllPaymentsByOrderId(orderID string) ([]*payment.Payment, *model.AppError)
	// GetLastOrderPayment get most recent payment made for given order
	GetLastOrderPayment(orderID string) (*payment.Payment, *model.AppError)
}

// CheckoutApp
type CheckoutApp interface {
}

type ProductApp interface {
}

type WishlistApp interface {
}

type AttributeApp interface {
}

type ChannelApp interface {
	// GetChannelBySlug get a channel from database with given slug
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
	// GetAllOrderLinesByOrderId returns a slice of order lines that belong to given order
	GetAllOrderLinesByOrderId(orderID string) ([]*order.OrderLine, *model.AppError)
	// OrderShippingIsRequired checks if an order requires ship or not by:
	//
	// 1) Find all child order lines that belong to given order
	//
	// 2) iterates over resulting slice to check if at least one order line requires shipping
	OrderShippingIsRequired(orderID string) (bool, *model.AppError)
	// OrderTotalQuantity return total quantity of given order
	OrderTotalQuantity(orderID string) (int, *model.AppError)
	// UpdateOrderTotalPaid update given order's total paid amount
	UpdateOrderTotalPaid(orderID string) *model.AppError
	// OrderIsPreAuthorized checks if order is pre-authorized
	OrderIsPreAuthorized(orderID string) (bool, *model.AppError)
	// OrderIsCaptured checks if given order is captured
	OrderIsCaptured(orderID string) (bool, *model.AppError)
	// OrderSubTotal returns sum of TotalPrice of all order lines that belong to given order
	OrderSubTotal(orderID string, orderCurrency string) (*goprices.TaxedMoney, *model.AppError)
	// OrderCanCalcel checks if given order can be canceled
	OrderCanCancel(ord *order.Order) (bool, *model.AppError)
	// OrderCanCapture
	OrderCanCapture(ord *order.Order, payment *payment.Payment) (bool, *model.AppError)
	// OrderCanVoid
	OrderCanVoid(ord *order.Order, payment *payment.Payment) (bool, *model.AppError)
	// OrderCanRefund checks if order can refund
	OrderCanRefund(ord *order.Order, payments []*payment.Payment) (bool, *model.AppError)
	// CanMarkOrderAsPaid checks if given order can be marked as paid.
	CanMarkOrderAsPaid(ord *order.Order, payments []*payment.Payment) (bool, *model.AppError)
	// OrderTotalAuthorized returns order's total authorized amount
	OrderTotalAuthorized(ord *order.Order) (*goprices.Money, *model.AppError)
	// GetOrderCountryCode is helper function, returns contry code of given order
	GetOrderCountryCode(ord *order.Order) (string, *model.AppError)
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

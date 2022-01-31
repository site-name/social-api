package order

import (
	"github.com/site-name/decimal"
	"github.com/sitename/sitename/app/plugin/interfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/order"
)

// SendPaymentConfirmation sends notification with the payment confirmation
func (s *ServiceOrder) SendPaymentConfirmation(orDer order.Order, manager interfaces.PluginManagerInterface) *model.AppError {
	panic("not implemented")
}

func (s *ServiceOrder) SendOrderCancelledConfirmation(orDer *order.Order, user *account.User, _, manager interfaces.PluginManagerInterface) *model.AppError {
	panic("not implemented")
}

// SendOrderConfirmation sends notification with order confirmation
func (s *ServiceOrder) SendOrderConfirmation(orDer *order.Order, redirectURL string, manager interfaces.PluginManagerInterface) *model.AppError {
	panic("not implemented")
}

// SendFulfillmentConfirmationToCustomer
//
// NOTE: user can be nil
func (s *ServiceOrder) SendFulfillmentConfirmationToCustomer(orDer *order.Order, fulfillment *order.Fulfillment, user *account.User, _, manager interfaces.PluginManagerInterface) *model.AppError {
	panic("not implemented")
}

// SendOrderConfirmed Send email which tells customer that order has been confirmed
func (s *ServiceOrder) SendOrderConfirmed(orDer order.Order, user *account.User, _ interface{}, manager interfaces.PluginManagerInterface) {

}

func (s *ServiceOrder) SendOrderRefundedConfirmation(orDer order.Order, user *account.User, _ interface{}, amount decimal.Decimal, currency string, manager interfaces.PluginManagerInterface) *model.AppError {
	panic("not implemented")
}

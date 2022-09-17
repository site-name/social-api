package order

import (
	"github.com/site-name/decimal"
	"github.com/sitename/sitename/app/plugin/interfaces"
	"github.com/sitename/sitename/model"
)

// SendPaymentConfirmation sends notification with the payment confirmation
func (s *ServiceOrder) SendPaymentConfirmation(orDer model.Order, manager interfaces.PluginManagerInterface) *model.AppError {
	panic("not implemented")
}

func (s *ServiceOrder) SendOrderCancelledConfirmation(orDer *model.Order, user *model.User, _, manager interfaces.PluginManagerInterface) *model.AppError {
	panic("not implemented")
}

// SendOrderConfirmation sends notification with order confirmation
func (s *ServiceOrder) SendOrderConfirmation(orDer *model.Order, redirectURL string, manager interfaces.PluginManagerInterface) *model.AppError {
	panic("not implemented")
}

// SendFulfillmentConfirmationToCustomer
//
// NOTE: user can be nil
func (s *ServiceOrder) SendFulfillmentConfirmationToCustomer(orDer *model.Order, fulfillment *model.Fulfillment, user *model.User, _, manager interfaces.PluginManagerInterface) *model.AppError {
	panic("not implemented")
}

// SendOrderConfirmed Send email which tells customer that order has been confirmed
func (s *ServiceOrder) SendOrderConfirmed(orDer model.Order, user *model.User, _ interface{}, manager interfaces.PluginManagerInterface) {

}

func (s *ServiceOrder) SendOrderRefundedConfirmation(orDer model.Order, user *model.User, _ interface{}, amount decimal.Decimal, currency string, manager interfaces.PluginManagerInterface) *model.AppError {
	panic("not implemented")
}

package order

import (
	"github.com/site-name/decimal"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/model/payment"
)

// OrderCreated
//
// fromDraft is default to false
func (a *AppOrder) OrderCreated(ord *order.Order, user *account.User, manager interface{}, fromDraft bool) *model.AppError {
	panic("not implemented")
}

// OrderConfirmed Trigger event, plugin hooks and optionally confirmation email.
func (a *AppOrder) OrderConfirmed(ord *order.Order, user *account.User, manager interface{}, sendConfirmationEmail bool) *model.AppError {
	panic("not implemented")
}

// HandleFullyPaidOrder
//
// user can be nil
func (a *AppOrder) HandleFullyPaidOrder(manager interface{}, ord *order.Order, user *account.User) *model.AppError {
	panic("not implemented")
}

// CancelOrder Release allocation of unfulfilled order items.
func (a *AppOrder) CancelOrder(ord *order.Order, user *account.User, manager interface{}) *model.AppError {
	panic("not implemented")
}

// OrderRefunded
func (a *AppOrder) OrderRefunded(ord *order.Order, user *account.User, amount *decimal.Decimal, payMent *payment.Payment, manager interface{}) *model.AppError {
	panic("not implemented")
}

// OrderVoided
func (a *AppOrder) OrderVoided(ord *order.Order, user *account.User, payMent *payment.Payment, manager interface{}) *model.AppError {
	panic("not implemented")
}

// OrderReturned
func (a *AppOrder) OrderReturned(ord *order.Order, user *account.User, returnedLines []*QuantityOrderLine) *model.AppError {
	var userID *string
	if user == nil {
		userID = nil
	} else {
		userID = &user.Id
	}

	_, appErr := a.CommonCreateOrderEvent(&order.OrderEventOption{
		OrderID: ord.Id,
		Type:    order.ORDER_EVENT_TYPE__FULFILLMENT_RETURNED,
		UserID:  userID,
		Parameters: &model.StringInterface{
			"lines": linesPerQuantityToLineObjectList(returnedLines),
		},
	})
	if appErr != nil {
		appErr.Where = "OrderReturned"
		return appErr
	}

	appErr = a.UpdateOrderStatus(ord)
	if appErr != nil {
		appErr.Where = "OrderReturned"
		return appErr
	}

	return nil
}

// OrderFulfilled
//
// notifyCustomer default to true
func (a *AppOrder) OrderFulfilled(fulfillments []*order.Fulfillment, user *account.User, fulfillmentLines []*order.FulfillmentLine, manager interface{}, notifyCustomer bool) *model.AppError {
	panic("not implemented")
}

// OrderShippingUpdated
func (a *AppOrder) OrderShippingUpdated(ord *order.Order, manager interface{}) *model.AppError {
	panic("not implemented")
}

// OrderAuthorized
func (a *AppOrder) OrderAuthorized(ord *order.Order, user *account.User, amount *decimal.Decimal, payMent *payment.Payment, manager interface{}) *model.AppError {
	panic("not implemented")
}

// OrderCaptured
func (a *AppOrder) OrderCaptured(ord *order.Order, user *account.User, amount *decimal.Decimal, payMent *payment.Payment, manager interface{}) *model.AppError {
	panic("not implemented")
}

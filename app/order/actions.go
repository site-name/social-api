package order

import (
	"net/http"

	"github.com/site-name/decimal"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/model/payment"
	"github.com/sitename/sitename/model/warehouse"
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

// FulfillmentTrackingUpdated
func (a *AppOrder) FulfillmentTrackingUpdated(fulfillment *order.Fulfillment, user *account.User, trackingNumber string, manager interface{}) *model.AppError {
	panic("not implemented")
}

// CancelFulfillment Return products to corresponding stocks.
func (a *AppOrder) CancelFulfillment(fulfillment *order.Fulfillment, user *account.User, warehouse *warehouse.WareHouse, manager interface{}) *model.AppError {
	panic("not implemented")
}

// Mark order as paid.
//
// Allows to create a payment for an order without actually performing any
// payment by the gateway.
//
// externalReference can be empty
func (a *AppOrder) MarkOrderAsPaid(ord *order.Order, requestUser *account.User, manager interface{}, externalReference string) *model.AppError {
	panic("not implemented")
}

// CleanMarkOrderAsPaid Check if an order can be marked as paid.
func (a *AppOrder) CleanMarkOrderAsPaid(ord *order.Order) *model.AppError {
	paymentsForOrder, appErr := a.PaymentApp().PaymentsByOption(&payment.PaymentFilterOption{
		OrderID: ord.Id,
	})
	if appErr != nil {
		appErr.Where = "CleanMarkOrderAsPaid"
	}

	if len(paymentsForOrder) > 0 {
		return model.NewAppError("CleanMarkOrderAsPaid", "app.order.order_with_payments_can_not_be_marked_as_paid.app_error", nil, "", http.StatusNotAcceptable)
	}

	return nil
}

// FulfillOrderLines Fulfill order line with given quantity
func (a *AppOrder) FulfillOrderLines(orderLineInfos []*order.OrderLineData) *model.AppError {
	orderLineInfosToDecreaseStock := a.WarehouseApp().GetOrderLinesWithTrackInventory(orderLineInfos)
	if len(orderLineInfosToDecreaseStock) > 0 {
		appErr := a.WarehouseApp().DecreaseStock(orderLineInfosToDecreaseStock, true)
		if appErr != nil {
			appErr.Where = "FulfillOrderLines"
			return appErr
		}
	}

	orderLines := []*order.OrderLine{}
	for _, lineInfo := range orderLineInfos {
		lineInfo.Line.QuantityFulfilled += lineInfo.Quantity
		orderLines = append(orderLines, &lineInfo.Line)
	}

	appErr := a.BulkUpsertOrderLines(orderLines)
	if appErr != nil {
		appErr.Where = "FulfillOrderLines"
		return appErr
	}

	return nil
}

// AutomaticallyFulfillDigitalLines
// Fulfill all digital lines which have enabled automatic fulfillment setting.
//
// Send confirmation email afterward.
func (a *AppOrder) AutomaticallyFulfillDigitalLines(ord *order.Order, manager interface{}) (appErr *model.AppError) {
	// find order lines of given order that are:
	// 1) NOT require shipping
	// 2) has ProductVariant attached AND that productVariant has a digitalContent accompanies
	defer func() {
		if appErr != nil {
			appErr.Where = "AutomaticallyFulfillDigitalLines"
		}
	}()

	digitalOrderLinesOfOrder, appErr := a.OrderLinesByOption(&order.OrderLineFilterOption{
		OrderID: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: ord.Id,
			},
		},
		IsShippingRequired: model.NewBool(false),

		Either: order.Either{
			VariantDigitalContentID: &model.StringFilter{
				StringOption: &model.StringOption{
					NULL: model.NewBool(false),
				},
			},
		},
		PrefetchRelated: true,
	})

	if appErr != nil {
		return
	}

	if len(digitalOrderLinesOfOrder) == 0 {
		return nil
	}

	var (
		fulfillment *order.Fulfillment
	)

	// try finding fulfillments that belong to given order
	fulfillmentsOfOrder, appErr := a.FulfillmentsByOption(&order.FulfillmentFilterOption{
		OrderID: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: ord.Id,
			},
		},
	})
	if appErr != nil {
		return
	}

	// if there is no fulfillment that belong to given order yet. We have to create a new one
	if len(fulfillmentsOfOrder) == 0 {
		fulfillment, appErr = a.UpsertFulfillment(&order.Fulfillment{
			OrderID: ord.Id,
		})
		if appErr != nil {
			return
		}
	} else {
		fulfillment = fulfillmentsOfOrder[0]
	}

	// finding shop that hold this order:
	ownerShopOfOrder, appErr := a.ShopApp().ShopById(ord.ShopID)
	if appErr != nil {
		return appErr
	}
	shopDefaultDigitalContentSettings := a.ProductApp().GetDefaultDigitalContentSettings(ownerShopOfOrder)

	var (
		fulfillmentLines []*order.FulfillmentLine
		orderLineDatas   []*order.OrderLineData
	)

	for _, orderLine := range digitalOrderLinesOfOrder {
		orderLineNeedsAutomaticFulfillment, appErr := a.OrderLineNeedsAutomaticFulfillment(orderLine, shopDefaultDigitalContentSettings)
		if appErr != nil {
			return appErr // must return if error occured
		}
		if !orderLineNeedsAutomaticFulfillment {
			continue
		}

		if orderLine.ProductVariant != nil || orderLine.VariantID != nil {

		}

		fulfillmentLines = append(fulfillmentLines, &order.FulfillmentLine{
			FulfillmentID: fulfillment.Id,
			OrderLineID:   orderLine.Id,
			Quantity:      orderLine.Quantity,
		})

		orderLineDatas = append(orderLineDatas, &order.OrderLineData{
			Line:     *orderLine,
			Quantity: orderLine.Quantity,
			Variant:  orderLine.ProductVariant,
			// WarehouseID: ,
		})
	}
}

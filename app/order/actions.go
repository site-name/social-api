package order

import (
	"net/http"

	"github.com/site-name/decimal"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/model/payment"
	"github.com/sitename/sitename/model/product_and_discount"
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
		VariantDigitalContentID: &model.StringFilter{
			StringOption: &model.StringOption{
				NULL: model.NewBool(false),
			},
		},
		PrefetchRelated: order.OrderLinePrefetchRelated{
			VariantDigitalContent: true, // this tell store to prefetch related product variants, digital contents too
		},
	})
	if appErr != nil {
		return
	}

	if len(digitalOrderLinesOfOrder) == 0 {
		return nil
	}

	fulfillment, appErr := a.GetOrCreateFulfillment(&order.FulfillmentFilterOption{
		OrderID: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: ord.Id,
			},
		},
	})
	if appErr != nil {
		return
	}

	// finding shop that hold this order:
	ownerShopOfOrder, appErr := a.ShopApp().ShopById(ord.ShopID)
	if appErr != nil {
		return
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

		if orderLine.ProductVariant != nil { // ProductVariant is available to use, prefetch option is enabled above
			_, appErr = a.ProductApp().CreateDigitalContentURL(&product_and_discount.DigitalContentUrl{
				LineID: &orderLine.Id,
			})
			if appErr != nil {
				return appErr
			}
		}

		fulfillmentLines = append(fulfillmentLines, &order.FulfillmentLine{
			FulfillmentID: fulfillment.Id,
			OrderLineID:   orderLine.Id,
			Quantity:      orderLine.Quantity,
		})

		allocationsOfOrderLine, appErr := a.WarehouseApp().AllocationsByOption(&warehouse.AllocationFilterOption{
			OrderLineID: &model.StringFilter{
				StringOption: &model.StringOption{
					Eq: orderLine.Id,
				},
			},
		})
		if appErr != nil {
			return appErr
		}

		stock, appErr := a.WarehouseApp().GetStockByOption(&warehouse.StockFilterOption{
			Id: &model.StringFilter{
				StringOption: &model.StringOption{
					Eq: allocationsOfOrderLine[0].StockID,
				},
			},
		})
		if appErr != nil {
			return appErr
		}

		orderLineDatas = append(orderLineDatas, &order.OrderLineData{
			Line:        *orderLine,
			Quantity:    orderLine.Quantity,
			Variant:     orderLine.ProductVariant,
			WarehouseID: stock.WarehouseID,
		})
	}

	// TODO: fixme
	panic("not implemented")
}

// Modify stocks and allocations. Return list of unsaved FulfillmentLines.
//
//     Args:
//         fulfillment (Fulfillment): Fulfillment to create lines
//         warehouse_pk (str): Warehouse to fulfill order.
//         lines_data (List[Dict]): List with information from which system
//             create FulfillmentLines. Example:
//                 [
//                     {
//                         "order_line": (OrderLine),
//                         "quantity": (int),
//                     },
//                     ...
//                 ]
//         channel_slug (str): Channel for which fulfillment lines should be created.
//
//     Return:
//         List[FulfillmentLine]: Unsaved fulfillmet lines created for this fulfillment
//             based on information form `lines`
//
//     Raise:
//         InsufficientStock: If system hasn't containt enough item in stock for any line.
func (a *AppOrder) createFulfillmentLines(fulfillment *order.Fulfillment, warehouseID string, lineDatas []map[string]*order.OrderLine, channelSlug string) ([]*order.FulfillmentLine, *model.AppError) {
	panic("not impl")
}

// Fulfill order.
//
//     Function create fulfillments with lines.
//     Next updates Order based on created fulfillments.
//
//     Args:
//         requester (User): Requester who trigger this action.
//         order (Order): Order to fulfill
//         fulfillment_lines_for_warehouses (Dict): Dict with information from which
//             system create fulfillments. Example:
//                 {
//                     (Warehouse.pk): [
//                         {
//                             "order_line": (OrderLine),
//                             "quantity": (int),
//                         },
//                         ...
//                     ]
//                 }
//         manager (PluginsManager): Base manager for handling plugins logic.
//         notify_customer (bool): If `True` system send email about
//             fulfillments to customer.
//
//     Return:
//         List[Fulfillment]: Fulfillmet with lines created for this order
//             based on information form `fulfillment_lines_for_warehouses`
//
//
//     Raise:
//         InsufficientStock: If system hasn't containt enough item in stock for any line.
func (a *AppOrder) CreateFulfillments(requester *account.User, ord *order.Order, fulfillmentLinesForWarehouse interface{}, manager interface{}, notifyCustomer bool) ([]*order.Fulfillment, *model.AppError) {
	panic("not impl")
}

// getFulfillmentLineIfExists
//
// NOTE: stockID can be empty
func (a *AppOrder) getFulfillmentLineIfExists(fulfillmentLines []*order.FulfillmentLine, orderLineID string, stockID string) *order.FulfillmentLine {
	for _, line := range fulfillmentLines {
		if line.OrderLineID == orderLineID && (line.StockID != nil && *line.StockID == stockID) {
			return line
		}
	}

	return nil
}

// getFulfillmentLine Get fulfillment line if extists or create new fulfillment line object.
//
// NOTE: stockID can be empty
func (a *AppOrder) getFulfillmentLine(targetFulfillment *order.Fulfillment, linesInTargetFulfillment []*order.FulfillmentLine, orderLineID string, stockID string) *struct {
	order.FulfillmentLine
	bool
} {
	// Check if line for order_line_id and stock_id does not exist in DB.
	movedFulfillmentLine := a.getFulfillmentLineIfExists(linesInTargetFulfillment, orderLineID, stockID)

	fulfillmentLineExisted := true

	var stockIdPointer *string
	if model.IsValidId(stockID) {
		stockIdPointer = &stockID
	}

	if movedFulfillmentLine == nil {
		// Create new not saved FulfillmentLine object and assign it to target fulfillment
		fulfillmentLineExisted = false
		movedFulfillmentLine = &order.FulfillmentLine{
			FulfillmentID: targetFulfillment.Id,
			OrderLineID:   orderLineID,
			StockID:       stockIdPointer,
			Quantity:      0,
		}
	}

	return &struct {
		order.FulfillmentLine
		bool
	}{
		*movedFulfillmentLine,
		fulfillmentLineExisted,
	}
}

// moveOrderLinesTotargetFulfillment Move order lines with given quantity to the target fulfillment
func (a *AppOrder) moveOrderLinesTotargetFulfillment(orderLinesToMove []*order.OrderLineData, targetFulfillment *order.Fulfillment) ([]*order.FulfillmentLine, *model.AppError) {
	panic("not implt")
}

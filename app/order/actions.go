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
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/modules/util"
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
func (a *AppOrder) FulfillOrderLines(orderLineInfos []*order.OrderLineData) (appErr *model.AppError) {
	defer func() {
		if appErr != nil {
			appErr.Where = "FulfillOrderLines"
		}
	}()

	orderLineInfosToDecreaseStock := a.WarehouseApp().GetOrderLinesWithTrackInventory(orderLineInfos)
	if len(orderLineInfosToDecreaseStock) > 0 {
		appErr := a.WarehouseApp().DecreaseStock(orderLineInfosToDecreaseStock, true)
		if appErr != nil {
			return appErr
		}
	}

	orderLines := []*order.OrderLine{}
	for _, lineInfo := range orderLineInfos {
		lineInfo.Line.QuantityFulfilled += lineInfo.Quantity
		orderLines = append(orderLines, &lineInfo.Line)
	}

	_, appErr = a.BulkUpsertOrderLines(orderLines)
	return appErr
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
func (a *AppOrder) getFulfillmentLineIfExists(fulfillmentLines []*order.FulfillmentLine, orderLineID string, stockID *string) *order.FulfillmentLine {
	for _, line := range fulfillmentLines {
		if line.OrderLineID == orderLineID &&
			(line.StockID != nil && stockID != nil && *line.StockID == *stockID) {
			return line
		}
	}

	return nil
}

type AResult struct {
	MovedFulfillmentLine *order.FulfillmentLine
	FulfillmentLineExist bool
}

// getFulfillmentLine Get fulfillment line if extists or create new fulfillment line object.
//
// NOTE: stockID can be empty
func (a *AppOrder) getFulfillmentLine(targetFulfillment *order.Fulfillment, linesInTargetFulfillment []*order.FulfillmentLine, orderLineID string, stockID *string) *AResult {
	// Check if line for order_line_id and stock_id does not exist in DB.
	movedFulfillmentLine := a.getFulfillmentLineIfExists(linesInTargetFulfillment, orderLineID, stockID)

	fulfillmentLineExisted := true

	if movedFulfillmentLine == nil {
		// Create new not saved FulfillmentLine object and assign it to target fulfillment
		fulfillmentLineExisted = false
		movedFulfillmentLine = &order.FulfillmentLine{
			FulfillmentID: targetFulfillment.Id,
			OrderLineID:   orderLineID,
			StockID:       stockID,
			Quantity:      0,
		}
	}

	return &AResult{
		MovedFulfillmentLine: movedFulfillmentLine,
		FulfillmentLineExist: fulfillmentLineExisted,
	}
}

// moveOrderLinesTotargetFulfillment Move order lines with given quantity to the target fulfillment
func (a *AppOrder) moveOrderLinesTotargetFulfillment(orderLinesToMove []*order.OrderLineData, targetFulfillment *order.Fulfillment) (fulfillmentLineToCreate []*order.FulfillmentLine, appErr *model.AppError) {

	defer func() {
		if appErr != nil {
			appErr.Where = "moveOrderLinesTotargetFulfillment"
		}
	}()

	var (
		orderLinesToUpdate        []*order.OrderLine
		orderLineDatasToDeAlocate []*order.OrderLineData
	)

	for _, lineData := range orderLinesToMove {
		// calculate the quantity fulfilled/unfulfilled to move
		unFulfilledToMove := util.Min(lineData.Line.QuantityUnFulfilled(), lineData.Quantity)
		lineData.Line.QuantityFulfilled += unFulfilledToMove

		// update current lines with new value of quantity
		orderLinesToUpdate = append(orderLinesToUpdate, &lineData.Line)
		fulfillmentLineToCreate = append(fulfillmentLineToCreate, &order.FulfillmentLine{
			FulfillmentID: targetFulfillment.Id,
			OrderLineID:   lineData.Line.Id,
			StockID:       nil,
			Quantity:      unFulfilledToMove,
		})

		allocationsOfOrderLine, appErr := a.WarehouseApp().AllocationsByOption(&warehouse.AllocationFilterOption{
			OrderLineID: &model.StringFilter{
				StringOption: &model.StringOption{
					Eq: lineData.Line.Id,
				},
			},
		})
		if appErr != nil {
			if appErr.StatusCode == http.StatusInternalServerError {
				return nil, appErr
			}
		}

		if len(allocationsOfOrderLine) > 0 {
			orderLineDatasToDeAlocate = append(orderLineDatasToDeAlocate, &order.OrderLineData{
				Line:     lineData.Line,
				Quantity: unFulfilledToMove,
			})
		}
	}

	if len(orderLineDatasToDeAlocate) > 0 {
		allocationErr, appErr := a.WarehouseApp().DeallocateStock(orderLineDatasToDeAlocate)
		if appErr != nil {
			return nil, appErr
		}

		if allocationErr != nil {
			slog.Warn("Unable to deallocate stock for order lines", slog.String("order_lines", allocationErr.OrderLineIDs()))
		}
	}

	fulfillmentLineToCreate, appErr = a.BulkUpsertFulfillmentLines(fulfillmentLineToCreate)
	if appErr != nil {
		return nil, appErr
	}

	_, appErr = a.BulkUpsertOrderLines(orderLinesToUpdate)
	if appErr != nil {
		return nil, appErr
	}

	return fulfillmentLineToCreate, nil
}

// moveFulfillmentLinesToTargetFulfillment Move fulfillment lines with given quantity to the target fulfillment
func (a *AppOrder) moveFulfillmentLinesToTargetFulfillment(fulfillmentLinesToMove []*order.FulfillmentLineData, linesInTargetFulfillment []*order.FulfillmentLine, targetFulfillment *order.Fulfillment) (appErr *model.AppError) {

	defer func() {
		if appErr != nil {
			appErr.Where = "moveFulfillmentLinesToTargetFulfillment"
		}
	}()

	var (
		fulfillmentLinesToCreate      []*order.FulfillmentLine
		fulfillmentLinesToUpdate      []*order.FulfillmentLine
		emptyFulfillmentLinesToDelete []*order.FulfillmentLine
	)

	for _, fulfillmentLineData := range fulfillmentLinesToMove {
		fulfillmentLine := fulfillmentLineData.Line
		quantityToMove := fulfillmentLineData.Quantity

		res := a.getFulfillmentLine(targetFulfillment, linesInTargetFulfillment, fulfillmentLine.OrderLineID, fulfillmentLine.StockID)

		// calculate the quantity fulfilled/unfulfilled/to move
		fulfilledToMove := util.Min(fulfillmentLine.Quantity, quantityToMove)
		quantityToMove -= fulfilledToMove
		res.MovedFulfillmentLine.Quantity += fulfilledToMove
		fulfillmentLine.Quantity -= fulfilledToMove

		if fulfillmentLine.Quantity == 0 {
			// the fulfillment line without any items will be deleted
			emptyFulfillmentLinesToDelete = append(emptyFulfillmentLinesToDelete, &fulfillmentLine)
		} else {
			// update with new quantity value
			fulfillmentLinesToUpdate = append(fulfillmentLinesToUpdate, &fulfillmentLine)
		}

		if res.MovedFulfillmentLine.Quantity > 0 && !res.FulfillmentLineExist {
			// If this is new type of (order_line, stock) then we create new fulfillment line
			fulfillmentLinesToCreate = append(fulfillmentLinesToCreate, res.MovedFulfillmentLine)
		} else if res.FulfillmentLineExist {
			// if target fulfillment already have the same line, we  just update the quantity
			fulfillmentLinesToUpdate = append(fulfillmentLinesToUpdate, res.MovedFulfillmentLine)
		}
	}

	// update the fulfillment lines with new values

	setAppErr := func(err *model.AppError) {
		a.mutex.Lock()
		if err != nil {
			appErr = err
		}
		a.mutex.Unlock()
	}

	a.wg.Add(3)
	go func() {
		defer a.wg.Done()
		_, err := a.BulkUpsertFulfillmentLines(fulfillmentLinesToUpdate)
		setAppErr(err)
	}()

	go func() {
		defer a.wg.Done()
		_, err := a.BulkUpsertFulfillmentLines(fulfillmentLinesToCreate)
		setAppErr(err)
	}()

	go func() {
		defer a.wg.Done()
		err := a.DeleteFulfillmentLinesByOption(&order.FulfillmentLineFilterOption{
			Id: &model.StringFilter{
				StringOption: &model.StringOption{
					In: order.FulfillmentLines(emptyFulfillmentLinesToDelete).IDs(),
				},
			},
		})
		setAppErr(err)
	}()

	a.wg.Done()

	return
}

func getShippingRefundAmount(refundShippingCosts bool, refundAmount *decimal.Decimal, shippingPrice *decimal.Decimal) *decimal.Decimal {
	// We set shipping refund amount only when refund amount is calculated
	var shippingRefundAmount *decimal.Decimal
	if refundShippingCosts && refundAmount == nil {
		shippingRefundAmount = shippingPrice
	}
	return shippingRefundAmount
}

// Proceed with all steps required for refunding products.
//
// Calculate refunds for products based on the order's lines and fulfillment
// lines.  The logic takes the list of order lines, fulfillment lines, and their
// quantities which is used to create the refund fulfillment. The stock for
// unfulfilled lines will be deallocated.
//
// NOTE: `refundShippingCosts` default to false
func (a *AppOrder) CreateRefundFulfillment(requester *account.User, ord *order.Order, payMent *payment.Payment, orderLinesToRefund []*order.OrderLineData, manager interface{}, amount *decimal.Decimal, refundShippingCosts bool) (interface{}, *model.AppError) {
	panic("not implt")
}

// populateReplaceOrderFields create new order based on the state of given originalOrder
//
// If original order has shippingAddress/billingAddress, the new order copy these address(es) and change their IDs
func (a *AppOrder) populateReplaceOrderFields(originalOrder *order.Order) (replaceOrder *order.Order, appErr *model.AppError) {
	replaceOrder = &order.Order{
		Status:             order.STATUS_DRAFT,
		UserID:             originalOrder.UserID,
		LanguageCode:       originalOrder.LanguageCode,
		UserEmail:          originalOrder.UserEmail,
		Currency:           originalOrder.Currency,
		ChannelID:          originalOrder.ChannelID,
		DisplayGrossPrices: originalOrder.DisplayGrossPrices,
		RedirectUrl:        originalOrder.RedirectUrl,
		OriginalID:         &originalOrder.Id,
		Origin:             order.REISSUE,
		ModelMetadata: model.ModelMetadata{
			Metadata:        originalOrder.Metadata,
			PrivateMetadata: originalOrder.PrivateMetadata,
		},
	}

	originalOrderAddressIDs := []string{}
	if originalOrder.BillingAddressID != nil {
		originalOrderAddressIDs = append(originalOrderAddressIDs, *originalOrder.BillingAddressID)
	}
	if originalOrder.ShippingAddressID != nil {
		originalOrderAddressIDs = append(originalOrderAddressIDs, *originalOrder.ShippingAddressID)
	}

	if len(originalOrderAddressIDs) > 0 {
		addressesOfOriginalOrder, appErr := a.AccountApp().AddressesByOption(&account.AddressFilterOption{
			Id: &model.StringFilter{
				StringOption: &model.StringOption{
					In: originalOrderAddressIDs,
				},
			},
		})
		if appErr != nil {
			return nil, appErr
		}

		for _, address := range addressesOfOriginalOrder {
			originalOrderAddressID := address.Id
			address.Id = ""
			newAddress, appErr := a.AccountApp().UpsertAddress(address)
			if appErr != nil {
				return nil, appErr
			}

			if originalOrder.BillingAddressID != nil && originalOrderAddressID == *originalOrder.BillingAddressID {
				replaceOrder.BillingAddressID = &newAddress.Id
			} else if originalOrder.ShippingAddressID != nil && originalOrderAddressID == *originalOrder.ShippingAddressID {
				replaceOrder.ShippingAddressID = &newAddress.Id
			}
		}
	}

	return a.UpsertOrder(replaceOrder)
}

// CreateReplaceOrder Create draft order with lines to replace
func (a *AppOrder) CreateReplaceOrder(requester *account.User, originalOrder *order.Order, orderLinesToReplace []*order.OrderLineData, fulfillmentLinesToReplace []*order.FulfillmentLineData) (replaceOrder *order.Order, appErr *model.AppError) {
	defer func() {
		if appErr != nil {
			appErr.Where = "CreateReplaceOrder"
		}
	}()

	replaceOrder, appErr = a.populateReplaceOrderFields(originalOrder)
	if appErr != nil {
		return
	}

	orderLinesToCreateMap := map[string]*order.OrderLine{}

	// iterate over lines without fulfillment to get the items for replace.
	// deepcopy to not lose the reference for lines assigned to original order
	for _, orderLineData := range order.OrderLineDatas(orderLinesToReplace).DeepCopy() {
		orderLine := orderLineData.Line
		orderLineID := orderLine.Id

		orderLine.Id = ""
		orderLine.OrderID = replaceOrder.Id
		orderLine.Quantity = orderLineData.Quantity
		orderLine.QuantityFulfilled = 0
		// we set order_line_id as a key to use it for iterating over fulfillment items
		orderLinesToCreateMap[orderLineID] = &orderLine
	}

	orderLineWithFulfillmentIDs := []string{}
	for _, lineData := range fulfillmentLinesToReplace {
		orderLineWithFulfillmentIDs = append(orderLineWithFulfillmentIDs, lineData.Line.OrderLineID)
	}

	orderLinesWithFulfillment, appErr := a.OrderLinesByOption(&order.OrderLineFilterOption{
		Id: &model.StringFilter{
			StringOption: &model.StringOption{
				In: orderLineWithFulfillmentIDs,
			},
		},
	})
	if appErr != nil {
		return
	}

	orderLinesWithFulfillmentMap := map[string]*order.OrderLine{}
	for _, id := range orderLineWithFulfillmentIDs {
		for _, orderLine := range orderLinesWithFulfillment {
			if id == orderLine.Id {
				orderLinesWithFulfillmentMap[id] = orderLine
			}
		}
	}

	for _, fulfillmentLineData := range fulfillmentLinesToReplace {
		fulfillmentLine := fulfillmentLineData.Line
		orderLineID := fulfillmentLine.OrderLineID

		// if order_line_id exists in order_line_to_create, it means that we already have
		// prepared new order_line for this fulfillment. In that case we need to increase
		// quantity amount of new order_line by fulfillment_line.quantity
		if item, exist := orderLinesToCreateMap[orderLineID]; exist && item != nil {
			orderLinesToCreateMap[orderLineID].Quantity += fulfillmentLineData.Quantity
			continue
		}

		orderLine := orderLinesWithFulfillmentMap[orderLineID]
		orderLineID = orderLine.Id
		orderLine.Id = ""
		orderLine.OrderID = replaceOrder.Id
		orderLine.Quantity = fulfillmentLineData.Quantity
		orderLine.QuantityFulfilled = 0
		orderLinesToCreateMap[orderLineID] = orderLine
	}

	orderLinesToCreate := []*order.OrderLine{}
	for _, orderLine := range orderLinesToCreateMap {
		orderLinesToCreate = append(orderLinesToCreate, orderLine)
	}

	_, appErr = a.BulkUpsertOrderLines(orderLinesToCreate)
	if appErr != nil {
		return
	}

	appErr = a.RecalculateOrder(replaceOrder, nil)
	if appErr != nil {
		return
	}

	var userID *string
	if requester != nil && model.IsValidId(requester.Id) {
		userID = &requester.Id
	}

	_, appErr = a.CommonCreateOrderEvent(&order.OrderEventOption{
		OrderID: replaceOrder.Id,
		Type:    order.ORDER_EVENT_TYPE__DRAFT_CREATED_FROM_REPLACE,
		UserID:  userID,
		Parameters: &model.StringInterface{
			"related_order_pk": originalOrder.Id,
			"lines":            linesPerQuantityToLineObjectList(orderLinesToQuantityOrderLine(orderLinesToCreate)),
		},
	})

	return
}

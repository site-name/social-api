package order

import (
	"net/http"
	"strings"

	"github.com/mattermost/gorp"
	"github.com/site-name/decimal"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/channel"
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/model/payment"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/model/warehouse"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/modules/util"
)

// OrderCreated. `fromDraft` is default to false
func (a *ServiceOrder) OrderCreated(ord *order.Order, user *account.User, manager interface{}, fromDraft bool) *model.AppError {
	panic("not implemented")
}

// OrderConfirmed Trigger event, plugin hooks and optionally confirmation email.
func (a *ServiceOrder) OrderConfirmed(ord *order.Order, user *account.User, manager interface{}, sendConfirmationEmail bool) *model.AppError {
	panic("not implemented")
}

// HandleFullyPaidOrder
//
// user can be nil
func (a *ServiceOrder) HandleFullyPaidOrder(manager interface{}, ord *order.Order, user *account.User) *model.AppError {
	panic("not implemented")
}

// CancelOrder Release allocation of unfulfilled order items.
func (a *ServiceOrder) CancelOrder(ord *order.Order, user *account.User, manager interface{}) *model.AppError {
	transaction, err := a.srv.Store.GetMaster().Begin()
	if err != nil {
		return model.NewAppError("CancelOrder", app.ErrorCreatingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}
	defer a.srv.Store.FinalizeTransaction(transaction)

	if err = transaction.Commit(); err != nil {
		return model.NewAppError("CancelOrder", app.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}
	// FIXME
	panic("not implemented")
}

// OrderRefunded
func (a *ServiceOrder) OrderRefunded(ord *order.Order, user *account.User, amount *decimal.Decimal, payMent *payment.Payment, manager interface{}) *model.AppError {
	panic("not implemented")
}

// OrderVoided
func (a *ServiceOrder) OrderVoided(ord *order.Order, user *account.User, payMent *payment.Payment, manager interface{}) *model.AppError {
	panic("not implemented")
}

// OrderReturned
func (a *ServiceOrder) OrderReturned(transaction *gorp.Transaction, ord *order.Order, user *account.User, returnedLines []*order.QuantityOrderLine) *model.AppError {
	var userID *string
	if user == nil {
		userID = nil
	} else {
		userID = &user.Id
	}

	_, appErr := a.CommonCreateOrderEvent(transaction, &order.OrderEventOption{
		OrderID: ord.Id,
		Type:    order.ORDER_EVENT_TYPE__FULFILLMENT_RETURNED,
		UserID:  userID,
		Parameters: &model.StringInterface{
			"lines": linesPerQuantityToLineObjectList(returnedLines),
		},
	})
	if appErr != nil {
		return appErr
	}

	appErr = a.UpdateOrderStatus(transaction, ord)
	if appErr != nil {
		return appErr
	}

	return nil
}

// OrderFulfilled
func (a *ServiceOrder) OrderFulfilled(fulfillments []*order.Fulfillment, user *account.User, fulfillmentLines []*order.FulfillmentLine, manager interface{}, notifyCustomer bool) *model.AppError {
	transaction, err := a.srv.Store.GetMaster().Begin()
	if err != nil {
		return model.NewAppError("CancelOrder", app.ErrorCreatingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}
	defer a.srv.Store.FinalizeTransaction(transaction)

	if err = transaction.Commit(); err != nil {
		return model.NewAppError("CancelOrder", app.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}
	panic("not implemented")
}

// OrderShippingUpdated
func (a *ServiceOrder) OrderShippingUpdated(ord *order.Order, manager interface{}) *model.AppError {
	panic("not implemented")
}

// OrderAuthorized
func (a *ServiceOrder) OrderAuthorized(ord *order.Order, user *account.User, amount *decimal.Decimal, payMent *payment.Payment, manager interface{}) *model.AppError {
	panic("not implemented")
}

// OrderCaptured
func (a *ServiceOrder) OrderCaptured(ord *order.Order, user *account.User, amount *decimal.Decimal, payMent *payment.Payment, manager interface{}) *model.AppError {
	panic("not implemented")
}

// FulfillmentTrackingUpdated
func (a *ServiceOrder) FulfillmentTrackingUpdated(fulfillment *order.Fulfillment, user *account.User, trackingNumber string, manager interface{}) *model.AppError {
	panic("not implemented")
}

// CancelFulfillment Return products to corresponding stocks.
func (a *ServiceOrder) CancelFulfillment(fulfillment *order.Fulfillment, user *account.User, warehouse *warehouse.WareHouse, manager interface{}) *model.AppError {
	// initialize a transaction
	transaction, err := a.srv.Store.GetMaster().Begin()
	if err != nil {
		return model.NewAppError("CancelOrder", app.ErrorCreatingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}
	defer a.srv.Store.FinalizeTransaction(transaction)

	fulfillment, appErr := a.FulfillmentByOption(transaction, &order.FulfillmentFilterOption{
		Id: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: fulfillment.Id,
			},
		},
		SelectForUpdate:    true, // this tells store to add `FOR UPDATE` to select query
		SelectRelatedOrder: true, // this make you free to acces `Order` of returning `fulfillment`
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return appErr
		}
		// if error is not found error, this mean given `fulfillment` is not valid
		return model.NewAppError("CancelFulfillment", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "fulfillment"}, appErr.DetailedError, http.StatusBadRequest)
	}

	appErr = a.RestockFulfillmentLines(transaction, fulfillment, warehouse)
	if appErr != nil {
		return appErr
	}

	var userID *string
	if user != nil && model.IsValidId(user.Id) {
		userID = &user.Id
	}
	_, appErr = a.CommonCreateOrderEvent(transaction, &order.OrderEventOption{
		OrderID: fulfillment.OrderID,
		UserID:  userID,
		Type:    order.ORDER_EVENT_TYPE__FULFILLMENT_CANCELED,
		Parameters: &model.StringInterface{
			"composed_id": fulfillment.ComposedId(),
		},
	})
	if appErr != nil {
		return appErr
	}

	// get total wuantity for order of given fulfillment:
	orderTotalQuantity, appErr := a.OrderTotalQuantity(fulfillment.OrderID)
	if appErr != nil {
		return appErr
	}
	_, appErr = a.CommonCreateOrderEvent(transaction, &order.OrderEventOption{
		OrderID: fulfillment.OrderID,
		UserID:  userID,
		Type:    order.ORDER_EVENT_TYPE__FULFILLMENT_RESTOCKED_ITEMS,
		Parameters: &model.StringInterface{
			"quantity":  orderTotalQuantity,
			"warehouse": warehouse.Id,
		},
	})
	if appErr != nil {
		return appErr
	}

	fulfillment.Status = order.FULFILLMENT_CANCELED
	_, appErr = a.UpsertFulfillment(transaction, fulfillment)
	if appErr != nil {
		return appErr
	}

	appErr = a.UpdateOrderStatus(transaction, fulfillment.Order) // you can access order here since store attached it to fulfillment above
	if appErr != nil {
		return appErr
	}

	//--------------------
	panic("not implemented")

	// commit transaction:
	if err = transaction.Commit(); err != nil {
		return model.NewAppError("CancelOrder", app.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}

// Mark order as paid.
//
// Allows to create a payment for an order without actually performing any
// payment by the gateway.
//
// externalReference can be empty
func (a *ServiceOrder) MarkOrderAsPaid(ord *order.Order, requestUser *account.User, manager interface{}, externalReference string) *model.AppError {
	transaction, err := a.srv.Store.GetMaster().Begin()
	if err != nil {
		return model.NewAppError("CancelOrder", app.ErrorCreatingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}
	defer a.srv.Store.FinalizeTransaction(transaction)

	if err = transaction.Commit(); err != nil {
		return model.NewAppError("CancelOrder", app.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	panic("not implemented")
}

// CleanMarkOrderAsPaid Check if an order can be marked as paid.
func (a *ServiceOrder) CleanMarkOrderAsPaid(ord *order.Order) (*payment.PaymentError, *model.AppError) {
	paymentsForOrder, appErr := a.srv.PaymentService().PaymentsByOption(&payment.PaymentFilterOption{
		OrderID: ord.Id,
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr
		}
		return nil, nil
	}

	if len(paymentsForOrder) > 0 {
		return payment.NewPaymentError("CleanMarkOrderAsPaid", "Orders with payments can not be manually marked as paid.", payment.INVALID), nil
	}

	return nil, nil
}

// FulfillOrderLines Fulfill order line with given quantity
func (a *ServiceOrder) FulfillOrderLines(orderLineInfos []*order.OrderLineData, manager interface{}) (*warehouse.InsufficientStock, *model.AppError) {
	transaction, err := a.srv.Store.GetMaster().Begin()
	if err != nil {
		return nil, model.NewAppError("AutomaticallyFulfillDigitalLines", app.ErrorCreatingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}
	defer a.srv.Store.FinalizeTransaction(transaction)

	orderLineInfosToDecreaseStock := a.srv.WarehouseService().GetOrderLinesWithTrackInventory(orderLineInfos)
	if len(orderLineInfosToDecreaseStock) > 0 {
		insufficientErr, appErr := a.srv.WarehouseService().DecreaseStock(orderLineInfosToDecreaseStock, manager, true)
		if appErr != nil || insufficientErr != nil {
			return insufficientErr, appErr
		}
	}

	orderLines := []*order.OrderLine{}
	for _, lineInfo := range orderLineInfos {
		lineInfo.Line.QuantityFulfilled += lineInfo.Quantity
		orderLines = append(orderLines, &lineInfo.Line)
	}

	_, appErr := a.BulkUpsertOrderLines(transaction, orderLines)
	if appErr != nil {
		return nil, appErr
	}

	if err = transaction.Commit(); err != nil {
		return nil, model.NewAppError("AutomaticallyFulfillDigitalLines", app.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return nil, nil
}

// AutomaticallyFulfillDigitalLines
// Fulfill all digital lines which have enabled automatic fulfillment setting. Send confirmation email afterward.
func (a *ServiceOrder) AutomaticallyFulfillDigitalLines(ord *order.Order, manager interface{}) *model.AppError {
	transaction, err := a.srv.Store.GetMaster().Begin()
	if err != nil {
		return model.NewAppError("AutomaticallyFulfillDigitalLines", app.ErrorCreatingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}
	defer a.srv.Store.FinalizeTransaction(transaction)

	// find order lines of given order that are:
	// 1) NOT require shipping
	// 2) has ProductVariant attached AND that productVariant has a digitalContent accompanies
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
		return appErr
	}

	if digitalOrderLinesOfOrder == nil || len(digitalOrderLinesOfOrder) == 0 {
		return nil
	}

	fulfillment, appErr := a.GetOrCreateFulfillment(transaction, &order.FulfillmentFilterOption{
		OrderID: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: ord.Id,
			},
		},
	})
	if appErr != nil {
		return appErr
	}

	// finding shop that hold this order:
	ownerShopOfOrder, appErr := a.srv.ShopService().ShopById(ord.ShopID)
	if appErr != nil {
		return appErr
	}
	shopDefaultDigitalContentSettings := a.srv.ProductService().GetDefaultDigitalContentSettings(ownerShopOfOrder)

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
			_, appErr = a.srv.ProductService().UpsertDigitalContentURL(&product_and_discount.DigitalContentUrl{
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

		allocationsOfOrderLine, appErr := a.srv.WarehouseService().AllocationsByOption(transaction, &warehouse.AllocationFilterOption{
			OrderLineID: &model.StringFilter{
				StringOption: &model.StringOption{
					Eq: orderLine.Id,
				},
			},
		})
		if appErr != nil {
			return appErr
		}

		stock, appErr := a.srv.WarehouseService().GetStockById(allocationsOfOrderLine[0].StockID)
		if appErr != nil {
			return appErr
		}

		orderLineDatas = append(orderLineDatas, &order.OrderLineData{
			Line:        *orderLine,
			Quantity:    orderLine.Quantity,
			Variant:     orderLine.ProductVariant,
			WarehouseID: &stock.WarehouseID,
		})
	}

	if err = transaction.Commit(); err != nil {
		return model.NewAppError("AutomaticallyFulfillDigitalLines", app.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
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
func (a *ServiceOrder) createFulfillmentLines(fulfillment *order.Fulfillment, warehouseID string, lineDatas order.QuantityOrderLines, channelSlug string, manager interface{}, decreaseStock bool) ([]*order.FulfillmentLine, *warehouse.InsufficientStock, *model.AppError) {

	var (
		variantIDs          = lineDatas.OrderLines().ProductVariantIDs()
		appError            *model.AppError
		stocksChan          = make(chan []*warehouse.Stock)
		productVariantsChan = make(chan []*product_and_discount.ProductVariant)
		syncSetAppError     = func(err *model.AppError) {
			a.mutex.Lock()
			defer a.mutex.Unlock()

			if err != nil && appError == nil {
				appError = err
			}
		}
	)
	defer func() {
		close(stocksChan)
		close(productVariantsChan)
	}()

	go func() {
		stocks, appErr := a.srv.WarehouseService().FilterStocksForChannel(&warehouse.StockFilterForChannelOption{
			ChannelSlug: channelSlug,
			WarehouseID: &model.StringFilter{
				StringOption: &model.StringOption{
					Eq: warehouseID,
				},
			},
			ProductVariantID: &model.StringFilter{
				StringOption: &model.StringOption{
					In: variantIDs,
				},
			},
			SelectRelatedProductVariant: true,
		})
		if appErr != nil {
			syncSetAppError(appErr)
		}
		stocksChan <- stocks

	}()

	go func() {
		productVariants, appErr := a.srv.ProductService().ProductVariantsByOption(&product_and_discount.ProductVariantFilterOption{
			Id: &model.StringFilter{
				StringOption: &model.StringOption{
					In: variantIDs,
				},
			},
			SelectRelatedDigitalContent: true, // NOTE: this asks store to populate related DigitalContent data to returned product variants
		})
		if appErr != nil {
			syncSetAppError(appErr)
		}
		productVariantsChan <- productVariants
	}()

	// lock here:
	productVariants := <-productVariantsChan
	stocks := <-stocksChan

	if appError != nil {
		return nil, nil, appError
	}

	// productVariantsMap has keys are product variant ids
	productVariantsMap := map[string]*product_and_discount.ProductVariant{}
	for _, variant := range productVariants {
		productVariantsMap[variant.Id] = variant
	}

	// variantToStock map has keys are product variant ids
	variantToStock := map[string][]*warehouse.Stock{}
	for _, stock := range stocks {
		variantToStock[stock.ProductVariantID] = append(variantToStock[stock.ProductVariantID], stock)
	}

	var (
		insufficientStocks              []*warehouse.InsufficientStockData
		fulfillmentLines                []*order.FulfillmentLine
		linesInfo                       []*order.OrderLineData
		quantity                        int
		orderLine                       *order.OrderLine
		productVariantOfOrderLine       product_and_discount.ProductVariant
		productVariantOfOrderLineIsReal bool
	)

	for _, line := range lineDatas {

		productVariantOfOrderLineIsReal = false
		quantity = line.Quantity
		orderLine = line.OrderLine
		productVariantOfOrderLine = product_and_discount.ProductVariant{}

		if orderLine.VariantID != nil && productVariantsMap[*orderLine.VariantID] != nil {
			productVariantOfOrderLine = *(productVariantsMap[*orderLine.VariantID])
			productVariantOfOrderLineIsReal = true
		}

		if quantity > 0 {
			if orderLine.VariantID == nil || variantToStock[*orderLine.VariantID] == nil {
				insufficientStocks = append(insufficientStocks, &warehouse.InsufficientStockData{
					Variant:     productVariantOfOrderLine,
					OrderLine:   orderLine,
					WarehouseID: &warehouseID,
				})
				continue
			}

			stock := variantToStock[*orderLine.VariantID][0]
			linesInfo = append(linesInfo, &order.OrderLineData{
				Line:        *orderLine,
				Quantity:    line.Quantity,
				Variant:     &productVariantOfOrderLine,
				WarehouseID: &warehouseID,
			})

			orderLineIsDigital, appErr := a.srv.OrderService().OrderLineIsDigital(orderLine)
			if appErr != nil {
				return nil, nil, appErr
			}
			if orderLineIsDigital && productVariantOfOrderLineIsReal {

				_, appErr = a.srv.ProductService().UpsertDigitalContentURL(&product_and_discount.DigitalContentUrl{
					ContentID: productVariantOfOrderLine.DigitalContent.Id, // check out 2nd goroutine above to see why is it possible to access DigitalContent.
					LineID:    &orderLine.Id,
				})
				if appErr != nil {
					return nil, nil, appErr
				}
			}

			fulfillmentLines = append(fulfillmentLines, &order.FulfillmentLine{
				OrderLineID:   orderLine.Id,
				FulfillmentID: fulfillment.Id,
				Quantity:      line.Quantity,
				StockID:       &stock.Id,
			})
		}
	}

	if len(insufficientStocks) > 0 {
		return nil,
			&warehouse.InsufficientStock{
				Items: insufficientStocks,
			},
			nil
	}

	if len(linesInfo) > 0 {
		insufficientStockErr, appErr := a.FulfillOrderLines(linesInfo, manager)
		if insufficientStockErr != nil || appErr != nil {
			return nil, insufficientStockErr, appErr
		}
	}

	return fulfillmentLines, nil, nil
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
func (a *ServiceOrder) CreateFulfillments(requester *account.User, orDer *order.Order, fulfillmentLinesForWarehouses map[string][]*order.QuantityOrderLine, manager interface{}, notifyCustomer bool, approved bool) ([]*order.Fulfillment, *warehouse.InsufficientStock, *model.AppError) {
	transaction, err := a.srv.Store.GetMaster().Begin()
	if err != nil {
		return nil, nil, model.NewAppError("CreateFulfillments", app.ErrorCreatingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}
	defer a.srv.Store.FinalizeTransaction(transaction)

	var (
		fulfillments     []*order.Fulfillment
		fulfillmentLines []*order.FulfillmentLine
	)

	channel, appErr := a.srv.ChannelService().ChannelByOption(&channel.ChannelFilterOption{
		Id: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: orDer.ChannelID,
			},
		},
	})
	if appErr != nil {
		return nil, nil, appErr
	}

	for warehouseID, quantityOrderLine := range fulfillmentLinesForWarehouses {
		fulfillment, appErr := a.UpsertFulfillment(transaction, &order.Fulfillment{
			OrderID: orDer.Id,
		})
		if appErr != nil {
			return nil, nil, appErr
		}

		fulfillments = append(fulfillments, fulfillment)

		filmentLines, insufficientStockErr, appErr := a.createFulfillmentLines(
			fulfillment,
			warehouseID,
			quantityOrderLine,
			channel.Slug,
			manager,
			approved,
		)
		if insufficientStockErr != nil || appErr != nil {
			return nil, insufficientStockErr, appErr
		}

		fulfillmentLines = append(fulfillmentLines, filmentLines...)
	}

	fulfillmentLines, appErr = a.BulkUpsertFulfillmentLines(transaction, fulfillmentLines)
	if appErr != nil {
		return nil, nil, appErr
	}

	// commit transaction
	if err = transaction.Commit(); err != nil {
		return nil, nil, model.NewAppError("CreateFulfillments", app.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	appErr = a.OrderFulfilled(fulfillments, requester, fulfillmentLines, manager, notifyCustomer)
	if appErr != nil {
		return nil, nil, appErr
	}

	return fulfillments, nil, nil
}

// getFulfillmentLineIfExists
//
// NOTE: stockID can be empty
func (a *ServiceOrder) getFulfillmentLineIfExists(fulfillmentLines []*order.FulfillmentLine, orderLineID string, stockID *string) *order.FulfillmentLine {
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
func (a *ServiceOrder) getFulfillmentLine(targetFulfillment *order.Fulfillment, linesInTargetFulfillment []*order.FulfillmentLine, orderLineID string, stockID *string) *AResult {
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

// moveOrderLinesToTargetFulfillment Move order lines with given quantity to the target fulfillment
func (a *ServiceOrder) moveOrderLinesToTargetFulfillment(orderLinesToMove []*order.OrderLineData, targetFulfillment *order.Fulfillment) (fulfillmentLineToCreate []*order.FulfillmentLine, appErr *model.AppError) {
	transaction, err := a.srv.Store.GetMaster().Begin()
	if err != nil {
		return nil, model.NewAppError("moveOrderLinesToTargetFulfillment", app.ErrorCreatingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}
	defer a.srv.Store.FinalizeTransaction(transaction)

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

		allocationsOfOrderLine, appErr := a.srv.WarehouseService().AllocationsByOption(transaction, &warehouse.AllocationFilterOption{
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

		if allocationsOfOrderLine != nil && len(allocationsOfOrderLine) > 0 {
			orderLineDatasToDeAlocate = append(orderLineDatasToDeAlocate, &order.OrderLineData{
				Line:     lineData.Line,
				Quantity: unFulfilledToMove,
			})
		}
	}

	if len(orderLineDatasToDeAlocate) > 0 {
		allocationErr, appErr := a.srv.WarehouseService().DeallocateStock(orderLineDatasToDeAlocate)
		if appErr != nil {
			return nil, appErr
		}

		if allocationErr != nil {
			slog.Warn("Unable to deallocate stock for order lines", slog.String("lines", strings.Join(allocationErr.OrderLines.IDs(), ", ")))
		}
	}

	fulfillmentLineToCreate, appErr = a.BulkUpsertFulfillmentLines(transaction, fulfillmentLineToCreate)
	if appErr != nil {
		return nil, appErr
	}

	_, appErr = a.BulkUpsertOrderLines(transaction, orderLinesToUpdate)
	if appErr != nil {
		return nil, appErr
	}

	if err = transaction.Commit(); err != nil {
		return nil, model.NewAppError("moveOrderLinesToTargetFulfillment", app.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return fulfillmentLineToCreate, nil
}

// moveFulfillmentLinesToTargetFulfillment Move fulfillment lines with given quantity to the target fulfillment
func (a *ServiceOrder) moveFulfillmentLinesToTargetFulfillment(fulfillmentLinesToMove []*order.FulfillmentLineData, linesInTargetFulfillment []*order.FulfillmentLine, targetFulfillment *order.Fulfillment) *model.AppError {
	transaction, err := a.srv.Store.GetMaster().Begin()
	if err != nil {
		return model.NewAppError("moveFulfillmentLinesToTargetFulfillment", app.ErrorCreatingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}
	defer a.srv.Store.FinalizeTransaction(transaction)

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
	var appError *model.AppError
	setAppErr := func(err *model.AppError) {
		a.mutex.Lock()
		if err != nil && appError == nil {
			appError = err
		}
		a.mutex.Unlock()
	}

	a.wg.Add(3)
	go func() {
		defer a.wg.Done()
		_, err := a.BulkUpsertFulfillmentLines(transaction, fulfillmentLinesToUpdate)
		setAppErr(err)
	}()

	go func() {
		defer a.wg.Done()
		_, err := a.BulkUpsertFulfillmentLines(transaction, fulfillmentLinesToCreate)
		setAppErr(err)
	}()

	go func() {
		defer a.wg.Done()
		err := a.DeleteFulfillmentLinesByOption(transaction, &order.FulfillmentLineFilterOption{
			Id: &model.StringFilter{
				StringOption: &model.StringOption{
					In: order.FulfillmentLines(emptyFulfillmentLinesToDelete).IDs(),
				},
			},
		})
		setAppErr(err)
	}()

	a.wg.Done()

	if appError != nil {
		return appError
	}

	if err = transaction.Commit(); err != nil {
		return model.NewAppError("moveFulfillmentLinesToTargetFulfillment", app.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
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
func (a *ServiceOrder) CreateRefundFulfillment(requester *account.User, ord *order.Order, payMent *payment.Payment, orderLinesToRefund []*order.OrderLineData, manager interface{}, amount *decimal.Decimal, refundShippingCosts bool) (interface{}, *model.AppError) {
	panic("not implt")
}

// populateReplaceOrderFields create new order based on the state of given originalOrder
//
// If original order has shippingAddress/billingAddress, the new order copy these address(es) and change their IDs
func (a *ServiceOrder) populateReplaceOrderFields(transaction *gorp.Transaction, originalOrder *order.Order) (replaceOrder *order.Order, appErr *model.AppError) {
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
		addressesOfOriginalOrder, appErr := a.srv.AccountService().AddressesByOption(&account.AddressFilterOption{
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
			newAddress, appErr := a.srv.AccountService().UpsertAddress(transaction, address)
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

	return a.UpsertOrder(transaction, replaceOrder)
}

// CreateReplaceOrder Create draft order with lines to replace
func (a *ServiceOrder) CreateReplaceOrder(requester *account.User, originalOrder *order.Order, orderLinesToReplace []*order.OrderLineData, fulfillmentLinesToReplace []*order.FulfillmentLineData) (*order.Order, *model.AppError) {
	transaction, err := a.srv.Store.GetMaster().Begin()
	if err != nil {
		return nil, model.NewAppError("CreateReplaceOrder", app.ErrorCreatingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}
	defer a.srv.Store.FinalizeTransaction(transaction)

	replaceOrder, appErr := a.populateReplaceOrderFields(transaction, originalOrder)
	if appErr != nil {
		return nil, appErr
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
		return nil, appErr
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

	_, appErr = a.BulkUpsertOrderLines(transaction, orderLinesToCreate)
	if appErr != nil {
		return nil, appErr
	}

	appErr = a.RecalculateOrder(transaction, replaceOrder, nil)
	if appErr != nil {
		return nil, appErr
	}

	var userID *string
	if requester != nil && model.IsValidId(requester.Id) {
		userID = &requester.Id
	}

	_, appErr = a.CommonCreateOrderEvent(transaction, &order.OrderEventOption{
		OrderID: replaceOrder.Id,
		Type:    order.ORDER_EVENT_TYPE__DRAFT_CREATED_FROM_REPLACE,
		UserID:  userID,
		Parameters: &model.StringInterface{
			"related_order_pk": originalOrder.Id,
			"lines":            linesPerQuantityToLineObjectList(orderLinesToQuantityOrderLine(orderLinesToCreate)),
		},
	})

	// commit transaction
	if err = transaction.Commit(); err != nil {
		return nil, model.NewAppError("CreateReplaceOrder", app.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return replaceOrder, nil
}

func (a *ServiceOrder) moveLinesToReturnFulfillment(
	orderLineDatas []*order.OrderLineData,
	fulfillmentLineDatas []*order.FulfillmentLineData,
	fulfillmentStatus string,
	ord *order.Order,
	totalRefundAmount *decimal.Decimal,
	shippingRefundAmount *decimal.Decimal,

) (*order.Fulfillment, *model.AppError) {

	targetFulfillment, appErr := a.UpsertFulfillment(nil, &order.Fulfillment{
		Status:               fulfillmentStatus,
		OrderID:              ord.Id,
		TotalRefundAmount:    totalRefundAmount,
		ShippingRefundAmount: shippingRefundAmount,
	})
	if appErr != nil {
		return nil, appErr
	}

	LinesInTargetFulfillment, appErr := a.moveOrderLinesToTargetFulfillment(orderLineDatas, targetFulfillment)
	if appErr != nil {
		return nil, appErr
	}

	fulfillmentLinesAlreadyRefunded, appErr := a.FulfillmentLinesByOption(&order.FulfillmentLineFilterOption{
		FulfillmentOrderID: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: ord.Id,
			},
		},
		FulfillmentStatus: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: order.FULFILLMENT_REFUNDED,
			},
		},
	})
	if appErr != nil {
		return nil, appErr
	}
	fulfillmentLinesAlreadyRefundedMap := model.MakeStringMapForModelSlice(
		fulfillmentLinesAlreadyRefunded,
		func(i interface{}) string {
			return i.(*order.FulfillmentLine).Id
		},
		nil,
	)

	var (
		refundedFulfillmentLinesToReturn []*order.FulfillmentLineData
		fulfillmentLinesToReturn         []*order.FulfillmentLineData
	)

	for _, lineData := range fulfillmentLineDatas {
		if item, exist := fulfillmentLinesAlreadyRefundedMap[lineData.Line.Id]; exist && item != nil {
			refundedFulfillmentLinesToReturn = append(refundedFulfillmentLinesToReturn, lineData)
			continue
		}

		fulfillmentLinesToReturn = append(fulfillmentLinesToReturn, lineData)
	}

	appErr = a.moveFulfillmentLinesToTargetFulfillment(fulfillmentLinesToReturn, LinesInTargetFulfillment, targetFulfillment)
	if appErr != nil {
		return nil, appErr
	}

	if len(refundedFulfillmentLinesToReturn) > 0 {
		var refundAndReturnFulfillment *order.Fulfillment
		if fulfillmentStatus == order.FULFILLMENT_REFUNDED_AND_RETURNED {
			refundAndReturnFulfillment = targetFulfillment
		} else {
			refundAndReturnFulfillment, appErr = a.UpsertFulfillment(nil, &order.Fulfillment{
				Status:  order.FULFILLMENT_REFUNDED_AND_RETURNED,
				OrderID: ord.Id,
			})
			if appErr != nil {
				return nil, appErr
			}
		}

		appErr = a.moveFulfillmentLinesToTargetFulfillment(refundedFulfillmentLinesToReturn, []*order.FulfillmentLine{}, refundAndReturnFulfillment)
		if appErr != nil {
			return nil, appErr
		}
	}

	return targetFulfillment, nil
}

func (a *ServiceOrder) moveLinesToReplaceFulfillment(
	orderLinesToReplace []*order.OrderLineData,
	fulfillmentLinesToReplace []*order.FulfillmentLineData,
	ord *order.Order,

) (*order.Fulfillment, *model.AppError) {

	targetFulfillment, appErr := a.UpsertFulfillment(nil, &order.Fulfillment{
		Status:  order.FULFILLMENT_REPLACED,
		OrderID: ord.Id,
	})
	if appErr != nil {
		return nil, appErr
	}

	linesInTargetFulfillment, appErr := a.moveOrderLinesToTargetFulfillment(orderLinesToReplace, targetFulfillment)
	if appErr != nil {
		return nil, appErr
	}

	appErr = a.moveFulfillmentLinesToTargetFulfillment(
		fulfillmentLinesToReplace,
		linesInTargetFulfillment,
		targetFulfillment,
	)

	return targetFulfillment, appErr
}

func (a *ServiceOrder) CreateReturnFulfillment(
	requester *account.User, // can be nil
	ord *order.Order,
	orderLineDatas []*order.OrderLineData,
	fulfillmentLineDatas []*order.FulfillmentLineData,
	totalRefundAmount *decimal.Decimal, // can be nil
	shippingRefundAmount *decimal.Decimal, // can be nil

) (*order.Fulfillment, *model.AppError) {

	// begin transaction
	transaction, err := a.srv.Store.GetMaster().Begin()
	if err != nil {
		return nil, model.NewAppError("CreateReturnFulfillment", app.ErrorCreatingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}
	defer a.srv.Store.FinalizeTransaction(transaction)

	status := order.FULFILLMENT_RETURNED
	if totalRefundAmount != nil {
		status = order.FULFILLMENT_REFUNDED_AND_RETURNED
	}

	returnFulfillment, appErr := a.moveLinesToReturnFulfillment(
		orderLineDatas,
		fulfillmentLineDatas,
		status,
		ord,
		totalRefundAmount,
		shippingRefundAmount,
	)
	if appErr != nil {
		return nil, appErr
	}

	returnedLines := map[string]*order.QuantityOrderLine{}

	orderLineIDs := []string{}
	for _, lineData := range fulfillmentLineDatas {
		orderLineIDs = append(orderLineIDs, lineData.Line.OrderLineID)
	}
	orderLinesByIDs, appErr := a.OrderLinesByOption(&order.OrderLineFilterOption{
		Id: &model.StringFilter{
			StringOption: &model.StringOption{
				In: orderLineIDs,
			},
		},
	})
	if appErr != nil {
		return nil, appErr
	}

	orderLinesByIDsMap := model.MakeStringMapForModelSlice(
		orderLinesByIDs,
		func(i interface{}) string {
			return i.(*order.OrderLine).Id
		},
		nil,
	)

	for _, orderLineData := range orderLineDatas {
		returnedLines[orderLineData.Line.Id] = &order.QuantityOrderLine{
			Quantity:  orderLineData.Quantity,
			OrderLine: &orderLineData.Line,
		}
	}

	for _, fulfillmentLineData := range fulfillmentLineDatas {
		if ifaceType := orderLinesByIDsMap[fulfillmentLineData.Line.OrderLineID]; ifaceType != nil {
			orderLine := ifaceType.(*order.OrderLine)
			returnedLine := returnedLines[orderLine.Id]

			if returnedLine != nil {
				returnedLines[orderLine.Id] = &order.QuantityOrderLine{
					Quantity:  returnedLine.Quantity + fulfillmentLineData.Quantity,
					OrderLine: returnedLine.OrderLine,
				}
			} else {
				returnedLines[orderLine.Id] = &order.QuantityOrderLine{
					Quantity:  fulfillmentLineData.Quantity,
					OrderLine: orderLine,
				}
			}
		}
	}

	sliceOfQuantityOrderLine := []*order.QuantityOrderLine{}
	for _, value := range returnedLines {
		sliceOfQuantityOrderLine = append(sliceOfQuantityOrderLine, value)
	}

	// commit
	if err = transaction.Commit(); err != nil {
		return nil, model.NewAppError("CreateReturnFulfillment", app.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	// NOTE: this is called after transaction commit
	appErr = a.OrderReturned(transaction, ord, requester, sliceOfQuantityOrderLine)

	return returnFulfillment, appErr
}

// ProcessReplace Create replace fulfillment and new draft order.
//
// Move all requested lines to fulfillment with status replaced. Based on original
// order create the draft order with all user details, and requested lines.
func (a *ServiceOrder) ProcessReplace(
	requester *account.User,
	ord *order.Order,
	orderLineDatas []*order.OrderLineData,
	fulfillmentLineDatas []*order.FulfillmentLineData,

) (*order.Fulfillment, *order.Order, *model.AppError) {

	transaction, err := a.srv.Store.GetMaster().Begin()
	if err != nil {
		return nil, nil, model.NewAppError("ProcessReplace", app.ErrorCreatingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}
	defer a.srv.Store.FinalizeTransaction(transaction)

	replaceFulfillment, appErr := a.moveLinesToReplaceFulfillment(orderLineDatas, fulfillmentLineDatas, ord)
	if appErr != nil {
		return nil, nil, appErr
	}

	newOrder, appErr := a.CreateReplaceOrder(requester, ord, orderLineDatas, fulfillmentLineDatas)
	if appErr != nil {
		return nil, nil, appErr
	}

	orderLinesOfOrder, appErr := a.OrderLinesByOption(&order.OrderLineFilterOption{
		OrderID: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: newOrder.Id,
			},
		},
	})
	if appErr != nil {
		return nil, nil, appErr
	}

	replacedLines := []*order.QuantityOrderLine{}
	for _, orderLine := range orderLinesOfOrder {
		replacedLines = append(replacedLines, &order.QuantityOrderLine{
			Quantity:  orderLine.Quantity,
			OrderLine: orderLine,
		})
	}

	var userID *string
	if requester != nil {
		userID = &requester.Id
	}

	_, appErr = a.CommonCreateOrderEvent(transaction, &order.OrderEventOption{
		OrderID: ord.Id,
		Type:    order.ORDER_EVENT_TYPE__FULFILLMENT_REPLACED,
		UserID:  userID,
		Parameters: &model.StringInterface{
			"lines": linesPerQuantityToLineObjectList(replacedLines),
		},
	})
	if appErr != nil {
		return nil, nil, appErr
	}

	_, appErr = a.CommonCreateOrderEvent(transaction, &order.OrderEventOption{
		OrderID: ord.Id,
		Type:    order.ORDER_EVENT_TYPE__ORDER_REPLACEMENT_CREATED,
		UserID:  userID,
		Parameters: &model.StringInterface{
			"related_order_pk": newOrder.Id,
		},
	})
	if appErr != nil {
		return nil, nil, appErr
	}

	if err = transaction.Commit(); err != nil {
		return nil, nil, model.NewAppError("ProcessReplace", app.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return replaceFulfillment, newOrder, nil
}

// Process the request for replacing or returning the products.
//
// Process the refund when the refund is set to True. The amount of refund will be
// calculated for all lines with statuses different from refunded.  The lines which
// are set to replace will not be included in the refund amount.
//
// If the amount is provided, the refund will be used for this amount.
//
// If refund_shipping_costs is True, the calculated refund amount will include
// shipping costs.
//
// All lines with replace set to True will be used to create a new draft order, with
// the same order details as the original order.  These lines will be moved to
// fulfillment with status replaced. The events with relation to new order will be
// created.
//
// All lines with replace set to False will be moved to fulfillment with status
// returned/refunded_and_returned - depends on refund flag and current line status.
// If the fulfillment line has refunded status it will be moved to
// returned_and_refunded
//
// NOTE: `payMent`, `amount` , `requester` are optional.
//
// `refund` and `refundShippingCosts` default to false.
//
func (a *ServiceOrder) CreateFulfillmentsForReturnedProducts(
	requester *account.User,
	ord *order.Order,
	payMent *payment.Payment,
	orderLineDatas []*order.OrderLineData,
	fulfillmentLineDatas []*order.FulfillmentLineData,
	manager interface{},
	refund bool,
	amount *decimal.Decimal,
	refundShippingCosts bool,

) (*order.Fulfillment, *order.Fulfillment, *order.Order, *model.AppError) {

	var (
		returnOrderLines        []*order.OrderLineData
		returnFulfillmentLines  []*order.FulfillmentLineData
		replaceOrderLines       []*order.OrderLineData
		replaceFulfillmentLines []*order.FulfillmentLineData
	)
	for _, lineData := range orderLineDatas {
		if !lineData.Replace {
			returnOrderLines = append(returnOrderLines, lineData)
			continue
		}
		replaceOrderLines = append(replaceOrderLines, lineData)
	}
	for _, lineData := range fulfillmentLineDatas {
		if !lineData.Replace {
			returnFulfillmentLines = append(returnFulfillmentLines, lineData)
			continue
		}
		replaceFulfillmentLines = append(replaceFulfillmentLines, lineData)
	}

	shippingRefundAmount := getShippingRefundAmount(refundShippingCosts, amount, ord.ShippingPriceGrossAmount)

	// create transaction
	transaction, err := a.srv.Store.GetMaster().Begin()
	if err != nil {
		return nil, nil, nil, model.NewAppError("CreateFulfillmentsForReturnedProducts", app.ErrorCreatingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}
	defer a.srv.Store.FinalizeTransaction(transaction)

	var (
		totalRefundAmount *decimal.Decimal
		appErr            *model.AppError
	)
	if refund && payMent != nil {
		totalRefundAmount, appErr = a.processRefund(
			requester,
			ord,
			payMent,
			returnOrderLines,
			returnFulfillmentLines,
			amount,
			refundShippingCosts,
			manager,
		)
		if appErr != nil {
			return nil, nil, nil, appErr
		}
	}

	var (
		replaceFulfillment *order.Fulfillment
		newOrder           *order.Order
	)
	if len(replaceFulfillmentLines) > 0 || len(replaceOrderLines) > 0 {
		replaceFulfillment, newOrder, appErr = a.ProcessReplace(
			requester,
			ord,
			replaceOrderLines,
			replaceFulfillmentLines,
		)
		if appErr != nil {
			return nil, nil, nil, appErr
		}
	}

	returnFulfillment, appErr := a.CreateReturnFulfillment(
		requester,
		ord,
		returnOrderLines,
		returnFulfillmentLines,
		totalRefundAmount,
		shippingRefundAmount,
	)
	if appErr != nil {
		return nil, nil, nil, appErr
	}

	fulfillmentsToDelete, appErr := a.FulfillmentsByOption(transaction, &order.FulfillmentFilterOption{
		OrderID: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: ord.Id,
			},
		},
		Status: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: order.FULFILLMENT_FULFILLED,
			},
		},
		FulfillmentLineID: &model.StringFilter{
			StringOption: &model.StringOption{
				NULL: model.NewBool(true),
			},
		},
	})
	if appErr != nil && appErr.StatusCode == http.StatusInternalServerError { // ignore not found err
		return nil, nil, nil, appErr
	}

	if len(fulfillmentsToDelete) > 0 {
		appErr = a.DeleteFulfillmentsByOption(transaction, &order.FulfillmentFilterOption{
			Id: &model.StringFilter{
				StringOption: &model.StringOption{
					In: fulfillmentsToDelete.IDs(),
				},
			},
		})
		if appErr != nil {
			return nil, nil, nil, appErr
		}
	}

	// commit transaction
	if err = transaction.Commit(); err != nil {
		return nil, nil, nil, model.NewAppError("CreateFulfillmentsForReturnedProducts", app.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	//----------------------------
	panic("not implt")

	return returnFulfillment, replaceFulfillment, newOrder, nil
}

func (a *ServiceOrder) calculateRefundAmount(
	returnOrderLineDatas []*order.OrderLineData,
	returnFulfillmentLineDatas []*order.FulfillmentLineData,
	linesToRefund map[string]*order.QuantityOrderLine,

) (*decimal.Decimal, *model.AppError) {

	refundAmount := decimal.Zero
	for _, lineData := range returnOrderLineDatas {
		if unitPriceGrossAmount := lineData.Line.UnitPriceGrossAmount; unitPriceGrossAmount != nil {
			refundAmount = refundAmount.Add(
				unitPriceGrossAmount.Mul(decimal.NewFromInt32(int32(lineData.Quantity))),
			)
		}
		linesToRefund[lineData.Line.Id] = &order.QuantityOrderLine{
			Quantity:  lineData.Quantity,
			OrderLine: &lineData.Line,
		}
	}

	if len(returnFulfillmentLineDatas) == 0 {
		return &refundAmount, nil
	}

	orderLineIDs := []string{}
	fulfillmentIDs := []string{}

	for _, lineData := range returnFulfillmentLineDatas {
		orderLineIDs = append(orderLineIDs, lineData.Line.OrderLineID)
		fulfillmentIDs = append(fulfillmentIDs, lineData.Line.FulfillmentID)
	}

	orderLines, appErr := a.OrderLinesByOption(&order.OrderLineFilterOption{
		Id: &model.StringFilter{
			StringOption: &model.StringOption{
				In: orderLineIDs,
			},
		},
	})
	if appErr != nil {
		return nil, appErr
	}

	fulfillments, appErr := a.FulfillmentsByOption(nil, &order.FulfillmentFilterOption{
		Id: &model.StringFilter{
			StringOption: &model.StringOption{
				In: fulfillmentIDs,
			},
		},
	})
	if appErr != nil {
		return nil, appErr
	}

	orderLinesMap := model.MakeStringMapForModelSlice(
		orderLines,
		func(i interface{}) string {
			return i.(*order.OrderLine).Id
		},
		nil,
	)
	fulfillmentsMap := model.MakeStringMapForModelSlice(
		fulfillments,
		func(i interface{}) string {
			return i.(*order.Fulfillment).Id
		},
		nil,
	)

	for _, lineData := range returnFulfillmentLineDatas {
		// skip lines which were already refunded
		ifaceType := fulfillmentsMap[lineData.Line.FulfillmentID]
		if ifaceType != nil && ifaceType.(*order.Fulfillment).Status == order.FULFILLMENT_REFUNDED {
			continue
		}

		if ifaceType = orderLinesMap[lineData.Line.OrderLineID]; ifaceType != nil {
			orderLine := ifaceType.(*order.OrderLine)
			if unitPriceGrossAmount := orderLine.UnitPriceGrossAmount; unitPriceGrossAmount != nil {
				refundAmount = refundAmount.Add(
					unitPriceGrossAmount.Mul(decimal.NewFromInt32(int32(lineData.Quantity))),
				)
			}

			dataFromAllRefundedLines := linesToRefund[orderLine.Id]
			if dataFromAllRefundedLines != nil {
				linesToRefund[orderLine.Id] = &order.QuantityOrderLine{
					Quantity:  dataFromAllRefundedLines.Quantity + lineData.Quantity,
					OrderLine: dataFromAllRefundedLines.OrderLine,
				}
			} else {
				linesToRefund[orderLine.Id] = &order.QuantityOrderLine{
					Quantity:  lineData.Quantity,
					OrderLine: orderLine,
				}
			}
		}
	}

	return &refundAmount, nil
}

// `requester` and `amount` can be nil
//
func (a *ServiceOrder) processRefund(
	requester *account.User,
	ord *order.Order,
	payMent *payment.Payment,
	orderLinesToRefund []*order.OrderLineData,
	fulfillmentLinesToRefund []*order.FulfillmentLineData,
	amount *decimal.Decimal,
	refundShippingCosts bool,
	manager interface{},

) (*decimal.Decimal, *model.AppError) {

	// transaction begin
	transaction, err := a.srv.Store.GetMaster().Begin()
	if err != nil {
		return nil, model.NewAppError("processRefund", app.ErrorCreatingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}
	defer a.srv.Store.FinalizeTransaction(transaction)

	linesToRefund := map[string]*order.QuantityOrderLine{}

	refundAmount, appErr := a.calculateRefundAmount(orderLinesToRefund, fulfillmentLinesToRefund, linesToRefund)
	if appErr != nil {
		return nil, appErr
	}

	if amount == nil {
		amount = refundAmount
		// we take into consideration the shipping costs only when amount is not provided.
		if refundShippingCosts && ord.ShippingPriceGrossAmount != nil {
			amount = model.NewDecimal(amount.Add(*ord.ShippingPriceGrossAmount))
		}
	}

	if amount != nil && !amount.Equal(decimal.Zero) {

	}

	// TODO: fix me
	panic("not implt")
}

package warehouse

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/mattermost/gorp"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/exception"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/model/warehouse"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
)

type StockData struct {
	Pk       string // ID of a stock
	Quantity int    // Quantity of the stock
}

// Allocate stocks for given `order_lines` in given country.
//
// Function lock for update all stocks and allocations for variants in
// given country and order by pk. Next, generate the dictionary
// ({"stock_pk": "quantity_allocated"}) with actual allocated quantity for stocks.
// Iterate by stocks and allocate as many items as needed or available in stock
// for order line, until allocated all required quantity for the order line.
// If there is less quantity in stocks then rise InsufficientStock exception.
func (a *ServiceWarehouse) AllocateStocks(orderLineInfos order.OrderLineDatas, countryCode string, channelSlug string, manager interface{}, additionalFilterLookup model.StringInterface) (*exception.InsufficientStock, *model.AppError) {
	transaction, err := a.srv.Store.GetMaster().Begin()
	if err != nil {
		return nil, model.NewAppError("AllocateStocks", app.ErrorCreatingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}
	defer a.srv.Store.FinalizeTransaction(transaction)

	// allocation only applied to order lines with variants with track inventory set to True
	orderLineInfos = a.GetOrderLinesWithTrackInventory(orderLineInfos)
	if len(orderLineInfos) == 0 {
		return nil, nil
	}

	stockFilterOption := &warehouse.StockFilterForCountryAndChannel{
		CountryCode: countryCode,
		ChannelSlug: channelSlug,

		ProductVariantIDFilter: &model.StringFilter{
			StringOption: &model.StringOption{
				In: orderLineInfos.Variants().IDs(),
			},
		},

		LockForUpdate: true,                 // FOR UPDATE
		ForUpdateOf:   store.StockTableName, // FOR UPDATE OF Stocks
	}

	// update lookup options:
	if additionalFilterLookup != nil {
		if warehouseId, ok := additionalFilterLookup["warehouse_id"]; ok && warehouseId != nil {
			if warehouseIdString, canCast := warehouseId.(string); canCast {
				stockFilterOption.WarehouseID = warehouseIdString
			}
		}
	}

	stocks, appErr := a.FilterStocksForCountryAndChannel(transaction, stockFilterOption)
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr // return immediately if error was caused by system
		}
		stocks = []*warehouse.Stock{} // just incase stocks is nil
	}

	quantityAllocationList, appErr := a.AllocationsByOption(nil, &warehouse.AllocationFilterOption{
		StockID: &model.StringFilter{
			StringOption: &model.StringOption{
				In: warehouse.Stocks(stocks).IDs(),
			},
		},
		QuantityAllocated: &model.NumberFilter{
			NumberOption: &model.NumberOption{
				Gt: model.NewFloat64(0),
			},
		},
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr // return immediately if error was caused by system
		}
		quantityAllocationList = []*warehouse.Allocation{} // just in case quantityAllocationList is nil
	}

	quantityAllocationForStocks := map[string]int{} // keys are stock IDs and values are sum of allocatedQuantity of allocations (which belong to a stock)
	for _, allocation := range quantityAllocationList {
		quantityAllocationForStocks[allocation.StockID] += allocation.QuantityAllocated
	}

	// the map below has: keys are IDs of product variants
	variantToStocks := map[string][]*StockData{}
	for _, stock := range stocks {
		variantToStocks[stock.ProductVariantID] = append(
			variantToStocks[stock.ProductVariantID],
			&StockData{
				Pk:       stock.Id,
				Quantity: stock.Quantity,
			},
		)
	}

	var (
		insufficientStock []*exception.InsufficientStockData
		allocations       warehouse.Allocations
		allocationItems   []*warehouse.Allocation
	)

	for _, lineInfo := range orderLineInfos {
		stockAllocations := variantToStocks[lineInfo.Variant.Id]
		insufficientStock, allocationItems = a.createAllocations(
			lineInfo,
			stockAllocations,
			quantityAllocationForStocks,
			insufficientStock,
		)
		allocations = append(allocations, allocationItems...)
	}

	if len(insufficientStock) > 0 {
		return &exception.InsufficientStock{Items: insufficientStock}, nil
	}

	// outOfStocks is a list of stocks that are have no item left
	var outOfStocks []*warehouse.Stock

	if len(allocations) > 0 {
		allocations, appErr = a.BulkUpsertAllocations(transaction, allocations)
		if appErr != nil {
			return nil, appErr
		}

		stockIDsOfAllocations := allocations.StockIDs()

		stocks, appErr := a.StocksByOption(transaction, &warehouse.StockFilterOption{
			Id: &model.StringFilter{
				StringOption: &model.StringOption{
					In: stockIDsOfAllocations,
				},
			},
		})
		if appErr != nil {
			return nil, appErr
		}
		// stockMap has keys are stock ids
		var stockMap = map[string]*warehouse.Stock{}
		for _, stock := range stocks {
			stockMap[stock.Id] = stock
		}

		allocationsOfStocks, appErr := a.AllocationsByOption(transaction, &warehouse.AllocationFilterOption{
			StockID: &model.StringFilter{
				StringOption: &model.StringOption{
					In: stockIDsOfAllocations,
				},
			},
		})
		if appErr != nil {
			return nil, appErr
		}

		// totalQuantityAllocatedOfStocksMap has keys are stock ids.
		// values are total quantity allocated of allocations of each stock
		var totalQuantityAllocatedOfStocksMap = map[string]int{}
		for _, allocation := range allocationsOfStocks {
			totalQuantityAllocatedOfStocksMap[allocation.StockID] += allocation.QuantityAllocated
		}

		for _, allocation := range allocations {
			if stock := stockMap[allocation.StockID]; stock != nil {
				if allocatedStock, ok := totalQuantityAllocatedOfStocksMap[stock.Id]; ok {
					if (stock.Quantity - allocatedStock) <= 0 {
						outOfStocks = append(outOfStocks, stock)
					}
				}
			}
		}
	}

	// commit transaction
	if err = transaction.Commit(); err != nil {
		return nil, model.NewAppError("AllocateStocks", app.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	if len(outOfStocks) > 0 {
		// TODO: fixme
		panic("not implemented")
	}

	return nil, nil
}

func (a *ServiceWarehouse) createAllocations(lineInfo *order.OrderLineData, stocks []*StockData, quantityAllocationForStocks map[string]int, insufficientStock []*exception.InsufficientStockData) ([]*exception.InsufficientStockData, []*warehouse.Allocation) {
	quantity := lineInfo.Quantity
	quantityAllocated := 0
	allocations := []*warehouse.Allocation{}

	for _, stockData := range stocks {
		quantityAllocatedInStock := quantityAllocationForStocks[stockData.Pk]
		quantityAvailableInStock := stockData.Quantity - quantityAllocatedInStock

		quantityToAllocate := util.Min(
			(quantity - quantityAllocated),
			quantityAvailableInStock,
		)

		if quantityToAllocate > 0 {
			allocations = append(allocations, &warehouse.Allocation{
				OrderLineID:       lineInfo.Line.Id,
				StockID:           stockData.Pk,
				QuantityAllocated: quantityToAllocate,
			})

			quantityAllocated += quantityToAllocate
			if quantityAllocated == quantity {
				return insufficientStock, allocations
			}
		}
	}

	if quantityAllocated != quantity {
		insufficientStock = append(insufficientStock, &exception.InsufficientStockData{
			Variant:   *lineInfo.Variant,
			OrderLine: &lineInfo.Line,
		})
	}

	return insufficientStock, []*warehouse.Allocation{}
}

// DeallocateStock Deallocate stocks for given `order_lines`.
//
// Function lock for update stocks and allocations related to given `order_lines`.
// Iterate over allocations sorted by `stock.pk` and deallocate as many items
// as needed of available in stock for order line, until deallocated all required
// quantity for the order line. If there is less quantity in stocks then
// raise an exception.
func (a *ServiceWarehouse) DeallocateStock(orderLineDatas []*order.OrderLineData, manager interface{}) (*warehouse.AllocationError, *model.AppError) {
	transaction, err := a.srv.Store.GetMaster().Begin()
	if err != nil {
		return nil, model.NewAppError("DeallocateStock", app.ErrorCreatingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}
	defer a.srv.Store.FinalizeTransaction(transaction)

	linesAllocations, appErr := a.AllocationsByOption(transaction, &warehouse.AllocationFilterOption{
		OrderLineID: &model.StringFilter{
			StringOption: &model.StringOption{
				In: order.OrderLineDatas(orderLineDatas).OrderLines().IDs(),
			},
		},
		LockForUpdate:        true,
		ForUpdateOf:          store.AllocationTableName + ", " + store.StockTableName,
		SelectedRelatedStock: true,
	})
	if appErr != nil {
		return nil, appErr
	}

	// lineToAllocations has keys are order line ids
	var lineToAllocations = map[string][]*warehouse.Allocation{}
	for _, allocation := range linesAllocations {
		lineToAllocations[allocation.OrderLineID] = append(lineToAllocations[allocation.OrderLineID], allocation)
	}

	var (
		allocationsToUpdate warehouse.Allocations
		notDeallocatedLines order.OrderLines
	)
	for _, lineInfo := range orderLineDatas {
		var (
			orderLine           = lineInfo.Line
			quantity            = lineInfo.Quantity
			allocations         = lineToAllocations[orderLine.Id]
			quantityDeAllocated = 0
		)

		for _, allocation := range allocations {
			quantityToDeallocate := util.Min(
				(quantity - quantityDeAllocated),
				allocation.QuantityAllocated,
			)
			if quantityToDeallocate > 0 {
				allocation.QuantityAllocated = allocation.QuantityAllocated - quantityToDeallocate
				quantityDeAllocated += quantityToDeallocate

				allocationsToUpdate = append(allocationsToUpdate, allocation)
				if quantityDeAllocated == quantity {
					break
				}
			}
		}

		if quantityDeAllocated != quantity {
			notDeallocatedLines = append(notDeallocatedLines, &orderLine)
		}
	}

	if len(notDeallocatedLines) > 0 {
		return &warehouse.AllocationError{OrderLines: notDeallocatedLines}, nil
	}

	allocationsBeforeUpdate, appErr := a.AllocationsByOption(transaction, &warehouse.AllocationFilterOption{
		Id: &model.StringFilter{
			StringOption: &model.StringOption{
				In: allocationsToUpdate.IDs(),
			},
		},

		SelectedRelatedStock:           true, // this tells store to attach `Stock` to each of returning allocations
		AnnotateStockAvailableQuantity: true, // this tells store to populate `StockAvailableQuantity` fields of returning allocations.
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr
		}
		// ignore not found error
		allocationsBeforeUpdate = make([]*warehouse.Allocation, 0)
	}

	_, appErr = a.BulkUpsertAllocations(transaction, allocationsToUpdate)
	if appErr != nil {
		return nil, appErr
	}

	// commit transaction:
	if err := transaction.Commit(); err != nil {
		return nil, model.NewAppError("DeallocateStock", app.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	// stockAndTotalQuantityAllocatedMap has keys are stock ids
	var stockAndTotalQuantityAllocatedMap = map[string]int{}
	for _, allocation := range allocationsBeforeUpdate {
		stockAndTotalQuantityAllocatedMap[allocation.StockID] += allocation.QuantityAllocated
	}

	for _, allocation := range allocationsBeforeUpdate {
		availableStockNow := util.Max(allocation.Stock.Quantity-stockAndTotalQuantityAllocatedMap[allocation.StockID], 0)

		if allocation.StockAvailableQuantity == 0 && availableStockNow > 0 {
			// TODO: fix me
			panic("not implemented")
		}
	}
	return nil, nil
}

// IncreaseStock Increse stock quantity for given `order_line` in a given warehouse.
//
// Function lock for update stock and allocations related to given `order_line`
// in a given warehouse. If the stock exists, increase the stock quantity
// by given value. If not exist create a stock with the given quantity. This function
// can create the allocation for increased quantity in stock by passing True
// to `allocate` argument. If the order line has the allocation in this stock
// function increase `quantity_allocated`. If allocation does not exist function
// create a new allocation for this order line in this stock.
//
// NOTE: allocate is default to false
func (a *ServiceWarehouse) IncreaseStock(orderLine *order.OrderLine, wareHouse *warehouse.WareHouse, quantity int, allocate bool) *model.AppError {
	transaction, err := a.srv.Store.GetMaster().Begin()
	if err != nil {
		return model.NewAppError("IncreaseStock", app.ErrorCreatingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}
	defer a.srv.Store.FinalizeTransaction(transaction)

	var stock *warehouse.Stock

	stocks, appErr := a.StocksByOption(transaction, &warehouse.StockFilterOption{
		WarehouseID: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: wareHouse.Id,
			},
		},
		ProductVariantID: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: *orderLine.VariantID,
			},
		},
		LockForUpdate: true,                 // FOR UPDATE
		ForUpdateOf:   store.StockTableName, // FOR UPDATE Stocks
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return appErr
		}
	} else {
		stock = stocks[0]
	}

	if stock != nil {
		stock.Quantity += quantity
	} else {
		// validate given `orderLine` has VariantID property not nil
		if orderLine.VariantID == nil {
			return model.NewAppError("IncreaseStock", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "orderLine"}, "orderLine must has VariantID property not nil", http.StatusBadRequest)
		}

		stock = &warehouse.Stock{
			WarehouseID:      wareHouse.Id,
			ProductVariantID: *orderLine.VariantID, // validated above
			Quantity:         quantity,
		}
	}
	_, appErr = a.BulkUpsertStocks(transaction, []*warehouse.Stock{stock})
	if appErr != nil {
		return appErr
	}

	if allocate && stock != nil {
		var allocation *warehouse.Allocation

		allocations, appErr := a.AllocationsByOption(transaction, &warehouse.AllocationFilterOption{
			OrderLineID: &model.StringFilter{
				StringOption: &model.StringOption{
					Eq: orderLine.Id,
				},
			},
			StockID: &model.StringFilter{
				StringOption: &model.StringOption{
					Eq: stock.Id,
				},
			},
		})
		if appErr != nil {
			if appErr.StatusCode == http.StatusInternalServerError {
				return appErr
			}
		} else {
			allocation = allocations[0]
		}

		if allocation != nil {
			allocation.QuantityAllocated += quantity
		} else {
			allocation = &warehouse.Allocation{
				OrderLineID:       orderLine.Id,
				StockID:           stock.Id,
				QuantityAllocated: quantity,
			}
		}

		_, appErr = a.BulkUpsertAllocations(transaction, []*warehouse.Allocation{allocation})
		if appErr != nil {
			return appErr
		}
	}

	if err = transaction.Commit(); err != nil {
		return model.NewAppError("IncreaseStock", app.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}

// IncreaseAllocations ncrease allocation for order lines with appropriate quantity
func (a *ServiceWarehouse) IncreaseAllocations(lineInfos []*order.OrderLineData, channelSlug string, manager interface{}) (*exception.InsufficientStock, *model.AppError) {
	// validate lineInfos is not nil nor empty
	if len(lineInfos) == 0 {
		return nil, model.NewAppError("IncreaseAllocations", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "lineInfos"}, "", http.StatusBadRequest)
	}

	// start a transaction
	transaction, err := a.srv.Store.GetMaster().Begin()
	if err != nil {
		return nil, model.NewAppError("IncreaseAllocations", app.ErrorCreatingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}
	defer a.srv.Store.FinalizeTransaction(transaction)

	allocations, appErr := a.AllocationsByOption(transaction, &warehouse.AllocationFilterOption{
		OrderLineID: &model.StringFilter{
			StringOption: &model.StringOption{
				In: order.OrderLineDatas(lineInfos).OrderLines().IDs(),
			},
		},
		LockForUpdate:          true,
		ForUpdateOf:            fmt.Sprintf("%s, %s", store.AllocationTableName, store.StockTableName),
		SelectedRelatedStock:   true,
		SelectRelatedOrderLine: true,
	})
	if appErr != nil {
		return nil, appErr
	}

	// evaluate allocations query to trigger select_for_update lock

	var (
		allocationIDsToDelete = warehouse.Allocations(allocations).IDs()

		// keys are IDs of order lines.
		// Values are lists of allocated quantities of allocations
		allocationQuantityMap = map[string][]int{}
	)

	for _, allocation := range allocations {
		allocationQuantityMap[allocation.OrderLineID] = append(allocationQuantityMap[allocation.OrderLineID], allocation.QuantityAllocated)
	}

	for _, lineInfo := range lineInfos {
		// lineInfo.quantity resembles amount to add, sum it with already allocated.
		lineInfo.Quantity += util.SumOfIntSlice(allocationQuantityMap[lineInfo.Line.Id])
	}

	if len(allocationIDsToDelete) > 0 {
		appErr = a.BulkDeleteAllocations(transaction, allocationIDsToDelete)
		if appErr != nil {
			return nil, appErr
		}
	}

	// find address of order of orderLine
	address, appErr := a.srv.OrderService().AnAddressOfOrder(lineInfos[0].Line.OrderID, order.ShippingAddressID)
	if appErr != nil {
		return nil, appErr
	}
	insufficientErr, appErr := a.AllocateStocks(lineInfos, address.Country, channelSlug, manager, nil)
	if insufficientErr != nil || appErr != nil {
		return insufficientErr, appErr
	}

	if err = transaction.Commit(); err != nil {
		return nil, model.NewAppError("IncreaseAllocations", app.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return nil, nil
}

// DecreaseAllocations Decreate allocations for provided order lines.
func (a *ServiceWarehouse) DecreaseAllocations(lineInfos []*order.OrderLineData, manager interface{}) (*exception.InsufficientStock, *model.AppError) {
	trackedOrderLines := a.GetOrderLinesWithTrackInventory(lineInfos)
	if len(trackedOrderLines) == 0 {
		return nil, nil
	}

	return a.DecreaseStock(lineInfos, manager, false, false)
}

// Decrease stocks quantities for given `order_lines` in given warehouses.
//
// Function deallocate as many quantities as requested if order_line has less quantity
// from requested function deallocate whole quantity. Next function try to find the
// stock in a given warehouse, if stock not exists or have not enough stock,
// the function raise InsufficientStock exception. When the stock has enough quantity
// function decrease it by given value.
// If update_stocks is False, allocations will decrease but stocks quantities
// will stay unmodified (case of unconfirmed order editing).
// If allow_stock_to_be_exceeded flag is True then quantity could be < 0.
//
// updateStocks default to true
func (a *ServiceWarehouse) DecreaseStock(orderLineInfos []*order.OrderLineData, manager interface{}, updateStocks bool, allowStockTobeExceeded bool) (*exception.InsufficientStock, *model.AppError) {
	// validate orderLineInfos is not nil nor empty
	if len(orderLineInfos) == 0 {
		return nil, model.NewAppError("DecreaseStock", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "orderLineInfos"}, "", http.StatusBadRequest)
	}

	transaction, err := a.srv.Store.GetMaster().Begin()
	if err != nil {
		return nil, model.NewAppError("DecreaseStock", app.ErrorCreatingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}
	defer a.srv.Store.FinalizeTransaction(transaction)

	var (
		variantIDs   = order.OrderLineDatas(orderLineInfos).Variants().IDs()
		warehouseIDs = order.OrderLineDatas(orderLineInfos).WarehouseIDs()
	)

	allocationErr, appErr := a.DeallocateStock(orderLineInfos, manager)
	if appErr != nil {
		return nil, appErr
	}
	if allocationErr != nil {
		allocations, appErr := a.AllocationsByOption(transaction, &warehouse.AllocationFilterOption{
			OrderLineID: &model.StringFilter{
				StringOption: &model.StringOption{
					In: allocationErr.OrderLines.IDs(),
				},
			},
		})
		if appErr != nil {
			if appErr.StatusCode == http.StatusInternalServerError {
				return nil, appErr
			}
		} else {
			for _, allocation := range allocations {
				allocation.QuantityAllocated = 0
			}

			_, appErr = a.BulkUpsertAllocations(transaction, allocations)
			if appErr != nil {
				return nil, appErr
			}
		}
	}

	stocks, appErr := a.StocksByOption(nil, &warehouse.StockFilterOption{
		ProductVariantID: &model.StringFilter{
			StringOption: &model.StringOption{
				In: variantIDs,
			},
		},
		WarehouseID: &model.StringFilter{
			StringOption: &model.StringOption{
				In: warehouseIDs,
			},
		},
		SelectRelatedProductVariant: true,
		SelectRelatedWarehouse:      true,
		LockForUpdate:               true,                 // add FOR UPDATE
		ForUpdateOf:                 store.StockTableName, // FOR UPDATE OF Stocks
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr
		}
		stocks = make([]*warehouse.Stock, 0)
	}

	variantAndWarehouseToStock := map[string]map[string]*warehouse.Stock{}
	for _, stock := range stocks {
		variantAndWarehouseToStock[stock.ProductVariantID][stock.WarehouseID] = stock
	}

	quantityAllocationList, appErr := a.AllocationsByOption(nil, &warehouse.AllocationFilterOption{
		StockID: &model.StringFilter{
			StringOption: &model.StringOption{
				In: warehouse.Stocks(stocks).IDs(),
			},
		},
		QuantityAllocated: &model.NumberFilter{
			NumberOption: &model.NumberOption{
				Gt: model.NewFloat64(0),
			},
		},
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr
		}
		quantityAllocationList = []*warehouse.Allocation{}
	}

	quantityAllocationForStocks := map[string]int{}
	for _, allocation := range quantityAllocationList {
		quantityAllocationForStocks[allocation.StockID] += allocation.QuantityAllocated
	}

	if updateStocks {
		insufficientErr, appErr := a.decreaseStocksQuantity(transaction, orderLineInfos, variantAndWarehouseToStock, quantityAllocationForStocks)
		if insufficientErr != nil || appErr != nil {
			return insufficientErr, appErr
		}
	}

	if err = transaction.Commit(); err != nil {
		return nil, model.NewAppError("DecreaseStock", app.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	if updateStocks {
		foundStocks, appErr := a.StocksByOption(nil, &warehouse.StockFilterOption{
			Id: &model.StringFilter{
				StringOption: &model.StringOption{
					In: warehouse.Stocks(stocks).IDs(),
				},
			},

			AnnotateAvailabeQuantity: true, // this tells store to populate AvailableQuantity fields of every returning stocks
		})
		if appErr != nil {
			if appErr.StatusCode == http.StatusInternalServerError {
				return nil, appErr
			}
			foundStocks = make([]*warehouse.Stock, 0)
		}

		for _, stock := range foundStocks {
			if stock.AvailableQuantity <= 0 {
				// TODO: fixme
				panic("not implemented")
			}
		}
	}

	return nil, nil
}

// decreaseStocksQuantity
func (a *ServiceWarehouse) decreaseStocksQuantity(transaction *gorp.Transaction, orderLinesInfo order.OrderLineDatas, variantAndwarehouseToStock map[string]map[string]*warehouse.Stock, quantityAllocationForStocks map[string]int) (*exception.InsufficientStock, *model.AppError) {

	var (
		insufficientStocks []*exception.InsufficientStockData
		stocksToUpdate     []*warehouse.Stock
	)

	for _, lineInfo := range orderLinesInfo {
		variant := lineInfo.Variant
		if variant == nil {
			continue
		}

		var stock *warehouse.Stock
		stockMap, ok := variantAndwarehouseToStock[variant.Id]
		if ok && stockMap != nil {
			if lineInfo.WarehouseID != nil {
				stock = stockMap[*lineInfo.WarehouseID]
			}
		}

		if stock == nil {
			insufficientStocks = append(insufficientStocks, &exception.InsufficientStockData{
				Variant:     *variant, // variant nil case is checked
				OrderLine:   &lineInfo.Line,
				WarehouseID: lineInfo.WarehouseID,
			})
			continue
		}

		quantityAllocated := quantityAllocationForStocks[stock.Id] // stock == nil already continued the loop
		if (stock.Quantity - quantityAllocated) < lineInfo.Quantity {
			insufficientStocks = append(insufficientStocks, &exception.InsufficientStockData{
				Variant:     *variant, // nil case checked
				OrderLine:   &lineInfo.Line,
				WarehouseID: lineInfo.WarehouseID,
			})
			continue
		}

		stock.Quantity -= lineInfo.Quantity
		stocksToUpdate = append(stocksToUpdate, stock)
	}

	if len(insufficientStocks) > 0 {
		return &exception.InsufficientStock{
			Items: insufficientStocks,
		}, nil
	}

	_, appErr := a.BulkUpsertStocks(transaction, stocksToUpdate)

	return nil, appErr
}

// GetOrderLinesWithTrackInventory Return order lines with variants with track inventory set to True
func (a *ServiceWarehouse) GetOrderLinesWithTrackInventory(orderLineInfos []*order.OrderLineData) []*order.OrderLineData {
	var res []*order.OrderLineData

	for _, lineInfo := range orderLineInfos {
		if lineInfo.Variant == nil || !*lineInfo.Variant.TrackInventory {
			res = append(res, lineInfo)
		}
	}

	return res
}

// DeAllocateStockForOrder Remove all allocations for given order
func (a *ServiceWarehouse) DeAllocateStockForOrder(ord *order.Order, manager interface{}) *model.AppError {
	transaction, err := a.srv.Store.GetMaster().Begin()
	if err != nil {
		return model.NewAppError("DeAllocateStockForOrder", app.ErrorCreatingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}
	defer a.srv.Store.FinalizeTransaction(transaction)

	allocations, appErr := a.AllocationsByOption(transaction, &warehouse.AllocationFilterOption{
		QuantityAllocated: &model.NumberFilter{
			NumberOption: &model.NumberOption{
				Gt: model.NewFloat64(0),
			},
		},
		OrderLineOrderID: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: ord.Id,
			},
		},

		AnnotateStockAvailableQuantity: true, // this tells store to populate StockAvailableQuantity fields of returning allocations
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return appErr
		}
		return nil
	}

	for i := range allocations {
		allocations[i].QuantityAllocated = 0
	}

	_, appErr = a.BulkUpsertAllocations(transaction, allocations)
	if appErr != nil {
		return appErr
	}

	if err = transaction.Commit(); err != nil {
		return model.NewAppError("DeAllocateStockForOrder", app.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	for _, allocation := range allocations {
		if allocation.StockAvailableQuantity <= 0 {
			// TODO: fix me
			panic("not implemented")
		}
	}

	return nil
}

// AllocatePreOrders allocates pre-order variant for given `order_lines` in given channel
func (s *ServiceWarehouse) AllocatePreOrders(orderLinesInfo order.OrderLineDatas, channelSlun string) *model.AppError {
	// init transaction
	transaction, err := s.srv.Store.GetMaster().Begin()
	if err != nil {
		return model.NewAppError("AllocatePreOrders", app.ErrorCreatingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}
	defer s.srv.Store.FinalizeTransaction(transaction)

	orderLinesInfoWithPreOrder := s.GetOrderLinesWithPreOrder(orderLinesInfo)
	if len(orderLinesInfoWithPreOrder) == 0 {
		return nil
	}

	variants := orderLinesInfoWithPreOrder.Variants()

	allVariantChannelListings, appErr := s.srv.ProductService().ProductVariantChannelListingsByOption(transaction, &product_and_discount.ProductVariantChannelListingFilterOption{
		VariantID: &model.StringFilter{
			StringOption: &model.StringOption{
				In: variants.IDs(),
			},
		},
		SelectRelatedChannel: true,
		SelectForUpdate:      true,
		SelectForUpdateOf:    store.ProductVariantChannelListingTableName,
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return appErr
		}
	}

	quantityAllocationList, appErr := s.PreOrderAllocationsByOptions(&warehouse.PreorderAllocationFilterOption{
		ProductVariantChannelListingID: &model.StringFilter{
			StringOption: &model.StringOption{
				In: allVariantChannelListings.IDs(),
			},
		},
	})

	var (
		quantityAllocationForChannel = map[string]int{}
	)

}

// GetOrderLinesWithPreOrder returns order lines with variants with preorder flag set to true
func (s *ServiceWarehouse) GetOrderLinesWithPreOrder(orderLinesInfo order.OrderLineDatas) order.OrderLineDatas {
	res := order.OrderLineDatas{}

	for _, lineInfo := range orderLinesInfo {
		if lineInfo.Variant != nil && lineInfo.Variant.IsPreorderActive() {
			res = append(res, lineInfo)
		}
	}

	return res
}

// variantChannelDataType
type variantChannelDataType struct {
	ChannelListingID         string
	ChannelQuantityThreshold *int
}

// createPreorderAllocation
func (s *ServiceWarehouse) createPreorderAllocation(lineInfo *order.OrderLineData, variantChannelData *variantChannelDataType, variantGlobalAllocation int, quantityAllocationForChannel map[string]int) (*warehouse.PreorderAllocation, *exception.InsufficientStockData, *model.AppError) {
	// validate valid arguments are provided:
	var invalidParams []string
	if variantChannelData == nil {
		invalidParams = append(invalidParams, "variantChannelData")
	}
	if lineInfo.Variant == nil {
		invalidParams = append(invalidParams, "lineInfo.Variant")
	}
	if len(invalidParams) > 0 {
		return nil, nil, model.NewAppError("createPreorderAllocation", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": strings.Join(invalidParams, ", ")}, "", http.StatusBadRequest)
	}

	var (
		variant                  = lineInfo.Variant // non-nil
		quantity                 = lineInfo.Quantity
		channelListingID         = variantChannelData.ChannelListingID
		channelQuantityThreshold = variantChannelData.ChannelQuantityThreshold
	)

	if channelQuantityThreshold != nil {
		channelAvailability := *channelQuantityThreshold - quantityAllocationForChannel[channelListingID]
		if quantity > channelAvailability {
			return nil, &exception.InsufficientStockData{
				Variant:           *variant,
				AvailableQuantity: &channelAvailability,
			}, nil
		}
	}

	if variant.PreOrderGlobalThreshold != nil {
		globalAvailability := *variant.PreOrderGlobalThreshold - variantGlobalAllocation
		if quantity > globalAvailability {
			return nil, &exception.InsufficientStockData{
				Variant:           *variant,
				AvailableQuantity: &globalAvailability,
			}, nil
		}
	}

	return &warehouse.PreorderAllocation{
		OrderLineID:                    lineInfo.Line.Id,
		ProductVariantChannelListingID: channelListingID,
		Quantity:                       quantity,
	}, nil, nil
}

// DeactivatePreorderForVariant Complete preorder for product variant.
// All preorder settings should be cleared and all preorder allocations
// should be replaced by regular allocations.
func (s *ServiceWarehouse) DeactivatePreorderForVariant(productVariant *product_and_discount.ProductVariant) (*exception.PreorderAllocationError, *model.AppError) {
	// init transaction:
	transaction, err := s.srv.Store.GetMaster().Begin()
	if err != nil {
		return nil, model.NewAppError("DeactivatePreorderForVariant", app.ErrorCreatingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}
	defer s.srv.Store.FinalizeTransaction(transaction)

	if !productVariant.IsPreOrder {
		return nil, nil
	}

	channelListings, appErr := s.srv.ProductService().ProductVariantChannelListingsByOption(transaction, &product_and_discount.ProductVariantChannelListingFilterOption{
		VariantID: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: productVariant.Id,
			},
		},
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr
		}
	}

	preorderAllocations, appErr := s.srv.WarehouseService().PreOrderAllocationsByOptions(&warehouse.PreorderAllocationFilterOption{
		ProductVariantChannelListingID: &model.StringFilter{
			StringOption: &model.StringOption{
				In: channelListings.IDs(),
			},
		},
		SelectRelated_OrderLine:       true,
		SelectRelated_OrderLine_Order: true,
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr
		}
	}

	var (
		allocationsToCreate []*warehouse.Allocation
		stocksToCreate      []*warehouse.Stock
	)
	for _, preorderAllocation := range preorderAllocations {
		stock, preorderAllocationErr, appErr := s.getStockForPreorderAllocation(preorderAllocation, productVariant)
		if preorderAllocationErr != nil || appErr != nil {
			return preorderAllocationErr, appErr
		}
		if !model.IsValidId(stock.Id) {
			stocksToCreate = append(stocksToCreate, stock)
		}
		allocationsToCreate = append(allocationsToCreate, &warehouse.Allocation{
			OrderLineID:       preorderAllocation.OrderLineID,
			StockID:           stock.Id,
			QuantityAllocated: preorderAllocation.Quantity,
		})
	}

	if len(stocksToCreate) > 0 {
		_, appErr = s.BulkUpsertStocks(transaction, stocksToCreate)
		if appErr != nil {
			return nil, appErr
		}
	}

	if len(allocationsToCreate) > 0 {
		_, appErr = s.BulkUpsertAllocations(transaction, allocationsToCreate)
		if appErr != nil {
			return nil, appErr
		}
	}

	if len(preorderAllocations) > 0 {
		appErr = s.DeletePreorderAllocations(transaction, preorderAllocations.IDs()...)
		if appErr != nil {
			return nil, appErr
		}
	}

	productVariant.PreOrderGlobalThreshold = nil
	productVariant.PreorderEndDate = nil
	productVariant.IsPreOrder = false
	_, appErr = s.srv.ProductService().UpsertProductVariant(transaction, productVariant)
	if appErr != nil {
		return nil, appErr
	}

	// NOTE: call the same query as above
	// the found result may difer the above since some row(s) may have been added during the period prior to this moment.
	productVariantChannelListings, appErr := s.srv.ProductService().ProductVariantChannelListingsByOption(transaction, &product_and_discount.ProductVariantChannelListingFilterOption{
		VariantID: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: productVariant.Id,
			},
		},
	})
	if appErr != nil {
		return nil, appErr
	}

	for _, listing := range productVariantChannelListings {
		listing.PreorderQuantityThreshold = nil
	}

	_, appErr = s.srv.ProductService().BulkUpsertProductVariantChannelListings(transaction, productVariantChannelListings)
	return nil, appErr
}

// getStockForPreorderAllocation Return stock where preordered variant should be allocated.
// By default this function uses any warehouse from the shipping zone that matches
// order's shipping method. If order has no shipping method set, it uses any warehouse
// that matches order's country. Function returns existing stock for selected warehouse
// or creates a new one unsaved `Stock` instance. Function raises an error if there is
// no warehouse assigned to any shipping zone handles order's country.
func (s *ServiceWarehouse) getStockForPreorderAllocation(preorderAllocation *warehouse.PreorderAllocation, productVariant *product_and_discount.ProductVariant) (*warehouse.Stock, *exception.PreorderAllocationError, *model.AppError) {
	panic("not implemented")
}

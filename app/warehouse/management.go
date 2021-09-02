package warehouse

import (
	"fmt"
	"net/http"

	"github.com/mattermost/gorp"
	"github.com/sitename/sitename/app"
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
func (a *ServiceWarehouse) AllocateStocks(orderLineInfos order.OrderLineDatas, countryCode string, channelSlug string) (*warehouse.InsufficientStock, *model.AppError) {
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

	stocks, appErr := a.FilterStocksForCountryAndChannel(transaction, &warehouse.StockFilterForCountryAndChannel{
		CountryCode: countryCode,
		ChannelSlug: channelSlug,

		ProductVariantIDFilter: &model.StringFilter{
			StringOption: &model.StringOption{
				In: product_and_discount.ProductVariants(orderLineInfos.Variants()).IDs(),
			},
		},

		LockForUpdate: true,                 // FOR UPDATE
		ForUpdateOf:   store.StockTableName, // FOR UPDATE OF Stocks
	})
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
		insufficientStock []*warehouse.InsufficientStockData
		allocations       []*warehouse.Allocation
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
		return &warehouse.InsufficientStock{Items: insufficientStock}, nil
	}

	if len(allocations) > 0 {
		_, appErr = a.BulkUpsertAllocations(transaction, allocations)
		if appErr != nil {
			return nil, appErr
		}
	}

	// commit transaction
	if err = transaction.Commit(); err != nil {
		return nil, model.NewAppError("AllocateStocks", app.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return nil, nil
}

func (a *ServiceWarehouse) createAllocations(lineInfo *order.OrderLineData, stocks []*StockData, quantityAllocationForStocks map[string]int, insufficientStock []*warehouse.InsufficientStockData) ([]*warehouse.InsufficientStockData, []*warehouse.Allocation) {
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
		insufficientStock = append(insufficientStock, &warehouse.InsufficientStockData{
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
func (a *ServiceWarehouse) DeallocateStock(orderLineDatas []*order.OrderLineData) (*warehouse.AllocationError, *model.AppError) {
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
		ForUpdateOf:          fmt.Sprintf("%s, %s", store.AllocationTableName, store.StockTableName),
		SelectedRelatedStock: true,
	})
	if appErr != nil {
		return nil, appErr
	}

	lineToAllocations := map[string][]*warehouse.Allocation{}
	for _, allocation := range linesAllocations {
		lineToAllocations[allocation.OrderLineID] = append(lineToAllocations[allocation.OrderLineID], allocation)
	}

	var (
		allocationsToUpdate []*warehouse.Allocation
		notDeallocatedLines []*order.OrderLine
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

	_, appErr = a.BulkUpsertAllocations(transaction, allocationsToUpdate)
	if appErr != nil {
		return nil, appErr
	}

	if err := transaction.Commit(); err != nil {
		return nil, model.NewAppError("DeallocateStock", app.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
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
func (a *ServiceWarehouse) IncreaseAllocations(lineInfos []*order.OrderLineData, channelSlug string) (*warehouse.InsufficientStock, *model.AppError) {
	// validate lineInfos is not nil nor empty
	if lineInfos == nil || len(lineInfos) == 0 {
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
	insufficientErr, appErr := a.AllocateStocks(lineInfos, address.Country, channelSlug)
	if insufficientErr != nil || appErr != nil {
		return insufficientErr, appErr
	}

	if err = transaction.Commit(); err != nil {
		return nil, model.NewAppError("IncreaseAllocations", app.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return nil, nil
}

// DecreaseAllocations Decreate allocations for provided order lines.
func (a *ServiceWarehouse) DecreaseAllocations(lineInfos []*order.OrderLineData) (*warehouse.InsufficientStock, *model.AppError) {
	trackedOrderLines := a.GetOrderLinesWithTrackInventory(lineInfos)
	if len(trackedOrderLines) == 0 {
		return nil, nil
	}

	return a.DecreaseStock(lineInfos, false)
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
//
// updateStocks default to true
func (a *ServiceWarehouse) DecreaseStock(orderLineInfos []*order.OrderLineData, updateStocks bool) (*warehouse.InsufficientStock, *model.AppError) {
	// validate orderLineInfos is not nil nor empty
	if orderLineInfos == nil || len(orderLineInfos) == 0 {
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

	allocationErr, appErr := a.DeallocateStock(orderLineInfos)
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
		stocks = []*warehouse.Stock{}
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

	return nil, nil
}

// decreaseStocksQuantity
func (a *ServiceWarehouse) decreaseStocksQuantity(transaction *gorp.Transaction, orderLinesInfo order.OrderLineDatas, variantAndwarehouseToStock map[string]map[string]*warehouse.Stock, quantityAllocationForStocks map[string]int) (*warehouse.InsufficientStock, *model.AppError) {

	var (
		insufficientStocks []*warehouse.InsufficientStockData
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
			insufficientStocks = append(insufficientStocks, &warehouse.InsufficientStockData{
				Variant:     *variant, // variant nil case is checked
				OrderLine:   &lineInfo.Line,
				WarehouseID: lineInfo.WarehouseID,
			})
			continue
		}

		quantityAllocated := quantityAllocationForStocks[stock.Id] // stock == nil already continued the loop
		if (stock.Quantity - quantityAllocated) < lineInfo.Quantity {
			insufficientStocks = append(insufficientStocks, &warehouse.InsufficientStockData{
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
		return &warehouse.InsufficientStock{
			Items: insufficientStocks,
		}, nil
	}

	_, appErr := a.BulkUpsertStocks(transaction, stocksToUpdate)

	return nil, appErr
}

// GetOrderLinesWithTrackInventory Return order lines with variants with track inventory set to True
func (a *ServiceWarehouse) GetOrderLinesWithTrackInventory(orderLineInfos []*order.OrderLineData) []*order.OrderLineData {
	for i, lineInfo := range orderLineInfos {
		if lineInfo.Variant == nil || !*lineInfo.Variant.TrackInventory {
			orderLineInfos = append(orderLineInfos[:i], orderLineInfos[i:]...)
		}
	}

	return orderLineInfos
}

// DeAllocateStockForOrder Remove all allocations for given order
func (a *ServiceWarehouse) DeAllocateStockForOrder(ord *order.Order) *model.AppError {
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
		LockForUpdate: true,                      // add `FOR UPDATE`
		ForUpdateOf:   store.AllocationTableName, // FOR UPDATE OF Allocations
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

	return nil
}

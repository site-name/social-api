package warehouse

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/plugin/interfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"gorm.io/gorm"
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
func (a *ServiceWarehouse) AllocateStocks(orderLineInfos model.OrderLineDatas, countryCode model.CountryCode, channelSlug string, manager interfaces.PluginManagerInterface, additionalFilterLookup model.StringInterface) (*model.InsufficientStock, *model.AppError) {
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

	stockFilterOption := &model.StockFilterForCountryAndChannel{
		CountryCode:            countryCode,
		ChannelSlug:            channelSlug,
		ProductVariantIDFilter: squirrel.Eq{model.StockTableName + ".ProductVariantID": orderLineInfos.Variants().IDs()},
		LockForUpdate:          true,                 // FOR UPDATE
		ForUpdateOf:            model.StockTableName, // FOR UPDATE OF Stocks
	}

	// update lookup options:
	if additionalFilterLookup != nil {
		if warehouseId, ok := additionalFilterLookup["warehouse_id"]; ok && warehouseId != nil {
			if warehouseIdString, ok := warehouseId.(string); ok {
				stockFilterOption.WarehouseID = warehouseIdString
			}
		}
	}

	stocks, appErr := a.FilterStocksForCountryAndChannel(stockFilterOption)
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr // return immediately if error was caused by system
		}
	}

	quantityAllocationList, appErr := a.AllocationsByOption(&model.AllocationFilterOption{
		StockID:           squirrel.Eq{model.AllocationTableName + ".StockID": model.Stocks(stocks).IDs()},
		QuantityAllocated: squirrel.Gt{model.AllocationTableName + ".QuantityAllocated": 0},
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr // return immediately if error was caused by system
		}
	}

	// quantityAllocationForStocks has keys are stock IDs and values are sum of allocatedQuantity of allocations (which belong to a stock)
	var quantityAllocationForStocks = map[string]int{}
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
		insufficientStock []*model.InsufficientStockData
		allocations       model.Allocations
		allocationItems   model.Allocations
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
		return &model.InsufficientStock{Items: insufficientStock}, nil
	}

	// outOfStocks is a list of stocks that are have no item left
	var outOfStocks []*model.Stock

	if len(allocations) > 0 {
		allocations, appErr = a.BulkUpsertAllocations(transaction, allocations)
		if appErr != nil {
			return nil, appErr
		}

		stockIDsOfAllocations := allocations.StockIDs()

		stocks, appErr := a.StocksByOption(&model.StockFilterOption{
			Conditions: squirrel.Eq{model.StockTableName + ".Id": stockIDsOfAllocations},
		})
		if appErr != nil {
			return nil, appErr
		}
		// stockMap has keys are stock ids
		var stockMap = map[string]*model.Stock{}
		for _, stock := range stocks {
			stockMap[stock.Id] = stock
		}

		allocationsOfStocks, appErr := a.AllocationsByOption(&model.AllocationFilterOption{
			StockID: squirrel.Eq{model.AllocationTableName + ".StockID": stockIDsOfAllocations},
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
		for _, stock := range outOfStocks {
			if appErr := manager.ProductVariantOutOfStock(*stock); appErr != nil {
				return nil, appErr
			}
		}
	}

	return nil, nil
}

func (a *ServiceWarehouse) createAllocations(lineInfo *model.OrderLineData, stocks []*StockData, quantityAllocationForStocks map[string]int, insufficientStock []*model.InsufficientStockData) ([]*model.InsufficientStockData, []*model.Allocation) {
	quantity := lineInfo.Quantity
	quantityAllocated := 0
	allocations := []*model.Allocation{}

	for _, stockData := range stocks {
		quantityAllocatedInStock := quantityAllocationForStocks[stockData.Pk]
		quantityAvailableInStock := stockData.Quantity - quantityAllocatedInStock

		quantityToAllocate := util.GetMinMax(
			(quantity - quantityAllocated),
			quantityAvailableInStock,
		).Min

		if quantityToAllocate > 0 {
			allocations = append(allocations, &model.Allocation{
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
		insufficientStock = append(insufficientStock, &model.InsufficientStockData{
			Variant:   *lineInfo.Variant,
			OrderLine: &lineInfo.Line,
		})
	}

	return insufficientStock, []*model.Allocation{}
}

// DeallocateStock Deallocate stocks for given `order_lines`.
//
// Function lock for update stocks and allocations related to given `order_lines`.
// Iterate over allocations sorted by `stock.pk` and deallocate as many items
// as needed of available in stock for order line, until deallocated all required
// quantity for the order line. If there is less quantity in stocks then
// raise an exception.
func (a *ServiceWarehouse) DeallocateStock(orderLineDatas model.OrderLineDatas, manager interfaces.PluginManagerInterface) (*model.AllocationError, *model.AppError) {

	linesAllocations, appErr := a.AllocationsByOption(&model.AllocationFilterOption{
		OrderLineID:          squirrel.Eq{model.AllocationTableName + ".OrderLineID": orderLineDatas.OrderLines().IDs()},
		LockForUpdate:        true,
		ForUpdateOf:          model.AllocationTableName + ", " + model.StockTableName,
		SelectedRelatedStock: true, //
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr
		}
		return nil, nil
	}

	// lineToAllocations has keys are order line ids
	var lineToAllocations = map[string][]*model.Allocation{}
	for _, allocation := range linesAllocations {
		lineToAllocations[allocation.OrderLineID] = append(lineToAllocations[allocation.OrderLineID], allocation)
	}

	var (
		allocationsToUpdate model.Allocations
		notDeallocatedLines model.OrderLines
	)
	for _, lineInfo := range orderLineDatas {
		var (
			orderLine           = lineInfo.Line
			quantity            = lineInfo.Quantity
			allocations         = lineToAllocations[orderLine.Id]
			quantityDeAllocated = 0
		)

		for _, allocation := range allocations {
			quantityToDeallocate := util.GetMinMax(
				(quantity - quantityDeAllocated),
				allocation.QuantityAllocated,
			).Min
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
		return &model.AllocationError{OrderLines: notDeallocatedLines}, nil
	}

	allocationsBeforeUpdate, appErr := a.AllocationsByOption(&model.AllocationFilterOption{
		Id:                             squirrel.Eq{model.AllocationTableName + ".Id": allocationsToUpdate.IDs()},
		SelectedRelatedStock:           true, // this tells store to attach `Stock` to each of returning allocations
		AnnotateStockAvailableQuantity: true, // this tells store to populate `StockAvailableQuantity` fields of returning allocations.
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr
		}
		// ignore not found error
	}

	transaction, err := a.srv.Store.GetMaster().Begin()
	if err != nil {
		return nil, model.NewAppError("DeallocateStock", app.ErrorCreatingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}
	defer a.srv.Store.FinalizeTransaction(transaction)

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
		availableStockNow := util.GetMinMax(allocation.GetStock().Quantity-stockAndTotalQuantityAllocatedMap[allocation.StockID], 0).Max

		if allocation.GetStockAvailableQuantity() <= 0 && availableStockNow > 0 {
			if appErr := manager.ProductVariantBackInStock(*allocation.GetStock()); appErr != nil {
				return nil, appErr
			}
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
func (a *ServiceWarehouse) IncreaseStock(orderLine *model.OrderLine, wareHouse *model.WareHouse, quantity int, allocate bool) *model.AppError {
	transaction, err := a.srv.Store.GetMaster().Begin()
	if err != nil {
		return model.NewAppError("IncreaseStock", app.ErrorCreatingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}
	defer a.srv.Store.FinalizeTransaction(transaction)

	var stock *model.Stock

	stocks, appErr := a.StocksByOption(&model.StockFilterOption{
		Conditions: squirrel.And{
			squirrel.Eq{model.ProductVariantTableName + ".Id": *orderLine.VariantID},
			squirrel.Eq{model.WarehouseTableName + ".Id": wareHouse.Id},
		},
		LockForUpdate: true,                 // FOR UPDATE
		ForUpdateOf:   model.StockTableName, // FOR UPDATE Stocks
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

		stock = &model.Stock{
			WarehouseID:      wareHouse.Id,
			ProductVariantID: *orderLine.VariantID, // validated above
			Quantity:         quantity,
		}
	}
	_, appErr = a.BulkUpsertStocks(transaction, []*model.Stock{stock})
	if appErr != nil {
		return appErr
	}

	if allocate && stock != nil {
		var allocation *model.Allocation

		allocations, appErr := a.AllocationsByOption(&model.AllocationFilterOption{
			OrderLineID: squirrel.Eq{model.AllocationTableName + ".OrderLineID": orderLine.Id},
			StockID:     squirrel.Eq{model.AllocationTableName + ".StockID": stock.Id},
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
			allocation = &model.Allocation{
				OrderLineID:       orderLine.Id,
				StockID:           stock.Id,
				QuantityAllocated: quantity,
			}
		}

		_, appErr = a.BulkUpsertAllocations(transaction, []*model.Allocation{allocation})
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
func (a *ServiceWarehouse) IncreaseAllocations(lineInfos model.OrderLineDatas, channelSlug string, manager interfaces.PluginManagerInterface) (*model.InsufficientStock, *model.AppError) {
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

	allocations, appErr := a.AllocationsByOption(&model.AllocationFilterOption{
		OrderLineID:            squirrel.Eq{model.AllocationTableName + ".OrderLineID": lineInfos.OrderLines().IDs()},
		LockForUpdate:          true,
		ForUpdateOf:            fmt.Sprintf("%s, %s", model.AllocationTableName, model.StockTableName),
		SelectedRelatedStock:   true,
		SelectRelatedOrderLine: true,
	})
	if appErr != nil {
		return nil, appErr
	}

	// evaluate allocations query to trigger select_for_update lock

	var (
		allocationIDsToDelete = model.Allocations(allocations).IDs()

		// keys are IDs of order lines.
		// Values are lists of allocated quantities of allocations
		allocationQuantityMap = map[string]util.AnyArray[int]{}
	)

	for _, allocation := range allocations {
		allocationQuantityMap[allocation.OrderLineID] = append(allocationQuantityMap[allocation.OrderLineID], allocation.QuantityAllocated)
	}

	for _, lineInfo := range lineInfos {
		// lineInfo.quantity resembles amount to add, sum it with already allocated.
		lineInfo.Quantity += allocationQuantityMap[lineInfo.Line.Id].Sum()
	}

	if len(allocationIDsToDelete) > 0 {
		appErr = a.BulkDeleteAllocations(transaction, allocationIDsToDelete)
		if appErr != nil {
			return nil, appErr
		}
	}

	// find address of order of orderLine
	address, appErr := a.srv.OrderService().AnAddressOfOrder(lineInfos[0].Line.OrderID, model.ShippingAddressID)
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
func (a *ServiceWarehouse) DecreaseAllocations(lineInfos []*model.OrderLineData, manager interfaces.PluginManagerInterface) (*model.InsufficientStock, *model.AppError) {
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
func (a *ServiceWarehouse) DecreaseStock(orderLineInfos model.OrderLineDatas, manager interfaces.PluginManagerInterface, updateStocks bool, allowStockTobeExceeded bool) (*model.InsufficientStock, *model.AppError) {
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
		variantIDs   = orderLineInfos.Variants().IDs()
		warehouseIDs = orderLineInfos.WarehouseIDs()
	)

	allocationErr, appErr := a.DeallocateStock(orderLineInfos, manager)
	if appErr != nil {
		return nil, appErr
	}
	if allocationErr != nil {
		allocations, appErr := a.AllocationsByOption(&model.AllocationFilterOption{
			OrderLineID: squirrel.Eq{model.AllocationTableName + ".OrderLineID": allocationErr.OrderLines.IDs()},
		})
		if appErr != nil {
			if appErr.StatusCode == http.StatusInternalServerError {
				return nil, appErr
			}
			// ignore not found error.
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

	stocks, appErr := a.StocksByOption(&model.StockFilterOption{
		Conditions: squirrel.And{
			squirrel.Eq{model.StockTableName + ".ProductVariantID": variantIDs},
			squirrel.Eq{model.StockTableName + ".WarehouseID": warehouseIDs},
		},
		SelectRelatedProductVariant: true,
		SelectRelatedWarehouse:      true,
		LockForUpdate:               true,                 // add FOR UPDATE
		ForUpdateOf:                 model.StockTableName, // FOR UPDATE OF Stocks
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr
		}
		// ignore not found error.
	}

	// variantAndWarehouseToStock has keys are product variant ids
	// values are map with keys are warehouse ids
	var variantAndWarehouseToStock = map[string]map[string]*model.Stock{}
	for _, stock := range stocks {
		variantAndWarehouseToStock[stock.ProductVariantID][stock.WarehouseID] = stock
	}

	quantityAllocationList, appErr := a.AllocationsByOption(&model.AllocationFilterOption{
		StockID:           squirrel.Eq{model.AllocationTableName + ".StockID": stocks.IDs()},
		QuantityAllocated: squirrel.Gt{model.AllocationTableName + ".QuantityAllocated": 0},
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr
		}
		// ignore not found error
	}

	// quantityAllocationForStocks has keys are stock ids
	var quantityAllocationForStocks = map[string]int{}
	for _, allocation := range quantityAllocationList {
		quantityAllocationForStocks[allocation.StockID] += allocation.QuantityAllocated
	}

	if updateStocks {
		insufficientErr, appErr := a.decreaseStocksQuantity(transaction, orderLineInfos, variantAndWarehouseToStock, quantityAllocationForStocks)
		if insufficientErr != nil || appErr != nil {
			return insufficientErr, appErr
		}
	}

	// commit transaction
	if err = transaction.Commit(); err != nil {
		return nil, model.NewAppError("DecreaseStock", app.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	if updateStocks {
		foundStocks, appErr := a.StocksByOption(&model.StockFilterOption{
			Conditions:               squirrel.Eq{model.StockTableName + ".Id": stocks.IDs()},
			AnnotateAvailabeQuantity: true, // this tells store to populate AvailableQuantity fields of every returning stocks
		})
		if appErr != nil {
			if appErr.StatusCode == http.StatusInternalServerError {
				return nil, appErr
			}
			// ignore not found error
		}

		for _, stock := range foundStocks {
			if stock.AvailableQuantity <= 0 {
				appErr = manager.ProductVariantOutOfStock(*stock)
				if appErr != nil {
					return nil, appErr
				}
			}
		}
	}

	return nil, nil
}

// decreaseStocksQuantity
func (a *ServiceWarehouse) decreaseStocksQuantity(transaction *gorm.DB, orderLinesInfo model.OrderLineDatas, variantAndwarehouseToStock map[string]map[string]*model.Stock, quantityAllocationForStocks map[string]int) (*model.InsufficientStock, *model.AppError) {

	var (
		insufficientStocks []*model.InsufficientStockData
		stocksToUpdate     []*model.Stock
	)

	for _, lineInfo := range orderLinesInfo {
		variant := lineInfo.Variant
		if variant == nil {
			continue
		}

		var stock *model.Stock
		stockMap, ok := variantAndwarehouseToStock[variant.Id]
		if ok && stockMap != nil {
			if lineInfo.WarehouseID != nil {
				stock = stockMap[*lineInfo.WarehouseID]
			}
		}

		if stock == nil {
			insufficientStocks = append(insufficientStocks, &model.InsufficientStockData{
				Variant:     *variant, // variant nil case is checked
				OrderLine:   &lineInfo.Line,
				WarehouseID: lineInfo.WarehouseID,
			})
			continue
		}

		quantityAllocated := quantityAllocationForStocks[stock.Id] // stock == nil already continued the loop
		if (stock.Quantity - quantityAllocated) < lineInfo.Quantity {
			insufficientStocks = append(insufficientStocks, &model.InsufficientStockData{
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
		return &model.InsufficientStock{
			Items: insufficientStocks,
		}, nil
	}

	_, appErr := a.BulkUpsertStocks(transaction, stocksToUpdate)

	return nil, appErr
}

// GetOrderLinesWithTrackInventory Return order lines with variants with track inventory set to True
func (a *ServiceWarehouse) GetOrderLinesWithTrackInventory(orderLineInfos []*model.OrderLineData) []*model.OrderLineData {
	var res []*model.OrderLineData

	for _, lineInfo := range orderLineInfos {
		if lineInfo.Variant == nil || !*lineInfo.Variant.TrackInventory {
			res = append(res, lineInfo)
		}
	}

	return res
}

// DeAllocateStockForOrder Remove all allocations for given order
func (a *ServiceWarehouse) DeAllocateStockForOrder(ord *model.Order, manager interfaces.PluginManagerInterface) *model.AppError {
	transaction, err := a.srv.Store.GetMaster().Begin()
	if err != nil {
		return model.NewAppError("DeAllocateStockForOrder", app.ErrorCreatingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}
	defer a.srv.Store.FinalizeTransaction(transaction)

	allocations, appErr := a.AllocationsByOption(&model.AllocationFilterOption{
		QuantityAllocated:              squirrel.Gt{model.AllocationTableName + ".QuantityAllocated": 0},
		OrderLineOrderID:               squirrel.Eq{model.OrderLineTableName + ".OrderID": ord.Id},
		AnnotateStockAvailableQuantity: true, // this tells store to populate StockAvailableQuantity fields of returning allocations
		SelectedRelatedStock:           true, //
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return appErr
		}
		return nil
	}

	allocationsToHandleAfterCommit := []*model.Allocation{}

	for _, allocation := range allocations {

		allocation.QuantityAllocated = 0

		if allocation.GetStockAvailableQuantity() <= 0 {
			allocationsToHandleAfterCommit = append(allocationsToHandleAfterCommit, allocation)
		}
	}

	_, appErr = a.BulkUpsertAllocations(transaction, allocations)
	if appErr != nil {
		return appErr
	}

	// commit transaction
	if err := transaction.Commit(); err != nil {
		return model.NewAppError("DeAllocateStockForOrder", app.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	for _, allocation := range allocationsToHandleAfterCommit {
		appErr = manager.ProductVariantBackInStock(*allocation.GetStock())
		if appErr != nil {
			return appErr
		}
	}

	return nil
}

// AllocatePreOrders allocates pre-order variant for given `order_lines` in given channel
func (s *ServiceWarehouse) AllocatePreOrders(orderLinesInfo model.OrderLineDatas, channelSlug string) (*model.InsufficientStock, *model.AppError) {
	// init transaction
	transaction, err := s.srv.Store.GetMaster().Begin()
	if err != nil {
		return nil, model.NewAppError("AllocatePreOrders", app.ErrorCreatingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}
	defer s.srv.Store.FinalizeTransaction(transaction)

	orderLinesInfoWithPreOrder := s.GetOrderLinesWithPreOrder(orderLinesInfo)
	if len(orderLinesInfoWithPreOrder) == 0 {
		return nil, nil
	}

	variants := orderLinesInfoWithPreOrder.Variants()

	allVariantChannelListings, appErr := s.srv.ProductService().
		ProductVariantChannelListingsByOption(&model.ProductVariantChannelListingFilterOption{
			VariantID:            squirrel.Eq{model.ProductVariantChannelListingTableName + ".VariantID": variants.IDs()},
			SelectRelatedChannel: true,
			SelectForUpdate:      true,
			SelectForUpdateOf:    model.ProductVariantChannelListingTableName,
		})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr
		}
	}

	quantityAllocationList, appErr := s.PreOrderAllocationsByOptions(&model.PreorderAllocationFilterOption{
		ProductVariantChannelListingID: squirrel.Eq{model.PreOrderAllocationTableName + ".ProductVariantChannelListingID": allVariantChannelListings.IDs()},
		Quantity:                       squirrel.Gt{model.PreOrderAllocationTableName + ".Quantity": 0},
	})
	if appErr != nil {
		return nil, appErr
	}

	var (
		// quantityAllocationForChannel has keys are product variant channel listing ids
		quantityAllocationForChannel = map[string]int{}
		// variantToChannelListings has keys are product variant ids
		variantToChannelListings = map[string]*variantChannelDataType{}
		// variantsGlobalAllocations has keys are product variant ids
		variantsGlobalAllocations = map[string]int{}
	)

	for _, allocation := range quantityAllocationList {
		quantityAllocationForChannel[allocation.ProductVariantChannelListingID] += allocation.Quantity
	}

	for _, channelListing := range allVariantChannelListings {
		if channelListing.GetChannel() != nil && channelListing.GetChannel().Slug == channelSlug {
			variantToChannelListings[channelListing.VariantID] = &variantChannelDataType{
				ChannelListingID:         channelListing.Id,
				ChannelQuantityThreshold: channelListing.PreorderQuantityThreshold,
			}
		}

		variantsGlobalAllocations[channelListing.VariantID] += quantityAllocationForChannel[channelListing.Id]
	}

	var (
		insufficientStocks []*model.InsufficientStockData
		allocations        []*model.PreorderAllocation
	)

	for _, lineInfo := range orderLinesInfo {
		variant := lineInfo.Variant
		if variant != nil {
			allocationItem, insufficientStockData, appErr := s.createPreorderAllocation(
				lineInfo,
				variantToChannelListings[variant.Id],
				variantsGlobalAllocations[variant.Id],
				quantityAllocationForChannel,
			)

			if appErr != nil {
				return nil, appErr // invalid argument app error
			}

			if allocationItem != nil {
				allocations = append(allocations, allocationItem)
			}
			if insufficientStockData != nil {
				insufficientStocks = append(insufficientStocks, insufficientStockData)
			}
		}
	}

	if len(insufficientStocks) > 0 {
		return model.NewInsufficientStock(insufficientStocks), nil
	}

	if len(allocations) > 0 {
		_, appErr = s.BulkCreate(transaction, allocations)
		if appErr != nil {
			return nil, appErr
		}
	}

	// commit transaction.
	if err = transaction.Commit(); err != nil {
		return nil, model.NewAppError("AllocatePreorders", app.ErrorCreatingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return nil, nil
}

// GetOrderLinesWithPreOrder returns order lines with variants with preorder flag set to true
func (s *ServiceWarehouse) GetOrderLinesWithPreOrder(orderLinesInfo model.OrderLineDatas) model.OrderLineDatas {
	res := model.OrderLineDatas{}

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
func (s *ServiceWarehouse) createPreorderAllocation(lineInfo *model.OrderLineData, variantChannelData *variantChannelDataType, variantGlobalAllocation int, quantityAllocationForChannel map[string]int) (*model.PreorderAllocation, *model.InsufficientStockData, *model.AppError) {
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
			return nil, &model.InsufficientStockData{
				Variant:           *variant,
				AvailableQuantity: &channelAvailability,
			}, nil
		}
	}

	if variant.PreOrderGlobalThreshold != nil {
		globalAvailability := *variant.PreOrderGlobalThreshold - variantGlobalAllocation
		if quantity > globalAvailability {
			return nil, &model.InsufficientStockData{
				Variant:           *variant,
				AvailableQuantity: &globalAvailability,
			}, nil
		}
	}

	return &model.PreorderAllocation{
		OrderLineID:                    lineInfo.Line.Id,
		ProductVariantChannelListingID: channelListingID,
		Quantity:                       quantity,
	}, nil, nil
}

// DeactivatePreorderForVariant Complete preorder for product variant.
// All preorder settings should be cleared and all preorder allocations
// should be replaced by regular allocations.
func (s *ServiceWarehouse) DeactivatePreorderForVariant(productVariant *model.ProductVariant) (*model.PreorderAllocationError, *model.AppError) {
	// init transaction:
	transaction, err := s.srv.Store.GetMaster().Begin()
	if err != nil {
		return nil, model.NewAppError("DeactivatePreorderForVariant", app.ErrorCreatingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}
	defer s.srv.Store.FinalizeTransaction(transaction)

	if !productVariant.IsPreOrder {
		return nil, nil
	}

	variantChannelListings, appErr := s.srv.ProductService().
		ProductVariantChannelListingsByOption(&model.ProductVariantChannelListingFilterOption{
			VariantID: squirrel.Eq{model.ProductVariantChannelListingTableName + ".VariantID": productVariant.Id},
		})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr
		}
	}

	preorderAllocations, appErr := s.srv.WarehouseService().PreOrderAllocationsByOptions(&model.PreorderAllocationFilterOption{
		ProductVariantChannelListingID: squirrel.Eq{model.PreOrderAllocationTableName + ".ProductVariantChannelListingID": variantChannelListings.IDs()},
		SelectRelated_OrderLine:        true,
		SelectRelated_OrderLine_Order:  true,
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr
		}
	}

	var (
		allocationsToCreate []*model.Allocation
		stocksToCreate      []*model.Stock
	)
	for _, preorderAllocation := range preorderAllocations {
		stock, preorderAllocationErr, appErr := s.getStockForPreorderAllocation(transaction, preorderAllocation, productVariant)
		if preorderAllocationErr != nil || appErr != nil {
			return preorderAllocationErr, appErr
		}
		if !model.IsValidId(stock.Id) {
			stocksToCreate = append(stocksToCreate, stock)
		}
		allocationsToCreate = append(allocationsToCreate, &model.Allocation{
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
	productVariantChannelListings, appErr := s.srv.ProductService().ProductVariantChannelListingsByOption(&model.ProductVariantChannelListingFilterOption{
		VariantID: squirrel.Eq{model.ProductVariantChannelListingTableName + ".VariantID": productVariant.Id},
	})
	if appErr != nil {
		return nil, appErr
	}

	for _, listing := range productVariantChannelListings {
		listing.PreorderQuantityThreshold = nil
	}

	_, appErr = s.srv.ProductService().BulkUpsertProductVariantChannelListings(transaction, productVariantChannelListings)
	if appErr != nil {
		return nil, appErr
	}

	// commit transaction
	if err = transaction.Commit(); err != nil {
		return nil, model.NewAppError("DeactivatePreorderForVariant", app.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return nil, nil
}

// getStockForPreorderAllocation Return stock where preordered variant should be allocated.
// By default this function uses any warehouse from the shipping zone that matches
// order's shipping method. If order has no shipping method set, it uses any warehouse
// that matches order's country. Function returns existing stock for selected warehouse
// or creates a new one unsaved `Stock` instance. Function raises an error if there is
// no warehouse assigned to any shipping zone handles order's country.
//
// NOTE: `transaction` MUST NOT be nil, otherwise this method'd return error
func (s *ServiceWarehouse) getStockForPreorderAllocation(transaction *gorm.DB, preorderAllocation *model.PreorderAllocation, productVariant *model.ProductVariant) (*model.Stock, *model.PreorderAllocationError, *model.AppError) {
	if transaction == nil {
		return nil, nil, model.NewAppError("getStockForPreorderAllocation", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "transaction"}, "please provide a non-nil transaction", http.StatusBadRequest)
	}

	var orDer *model.Order
	if preorderAllocation != nil &&
		preorderAllocation.GetOrderLine() != nil &&
		preorderAllocation.GetOrderLine().GetOrder() != nil {
		orDer = preorderAllocation.GetOrderLine().GetOrder()

	} else {
		preorderAllocations, appErr := s.srv.
			WarehouseService().
			PreOrderAllocationsByOptions(&model.PreorderAllocationFilterOption{
				SelectRelated_OrderLine:       true,
				SelectRelated_OrderLine_Order: true,
				Id:                            squirrel.Eq{model.PreOrderAllocationTableName + ".Id": preorderAllocation.Id},
			})
		if appErr != nil {
			return nil, nil, appErr
		}
		preorderAllocation = preorderAllocations[0]
		orDer = preorderAllocation.GetOrderLine().GetOrder()
	}

	var wareHouse *model.WareHouse

	if orDer.ShippingMethodID != nil {
		orderShippingMethod, appErr := s.srv.
			ShippingService().
			ShippingMethodByOption(&model.ShippingMethodFilterOption{
				Id: squirrel.Eq{
					model.ShippingMethodTableName + ".Id": *orDer.ShippingMethodID,
				},
			})
		if appErr != nil {
			return nil, nil, appErr
		}

		warehouses, appErr := s.srv.
			WarehouseService().
			WarehousesByOption(&model.WarehouseFilterOption{
				ShippingZonesId: squirrel.Eq{model.ShippingZoneTableName + ".Id": orderShippingMethod.ShippingZoneID},
			})
		if appErr != nil {
			if appErr.StatusCode == http.StatusInternalServerError {
				return nil, nil, appErr
			}
			// ignore not found error
		}
		if len(warehouses) != 0 {
			wareHouse = warehouses[0]
		}
	} else {
		orderCountry, appErr := s.srv.OrderService().GetOrderCountry(orDer)
		if appErr != nil {
			return nil, nil, appErr
		}

		warehouses, appErr := s.srv.
			WarehouseService().
			WarehousesByOption(&model.WarehouseFilterOption{
				ShippingZonesCountries: squirrel.Like{model.ShippingMethodTableName + ".Countries": orderCountry},
			})
		if appErr != nil {
			if appErr.StatusCode == http.StatusInternalServerError {
				return nil, nil, appErr
			}
			// ignore not found error
		}
		if len(warehouses) != 0 {
			wareHouse = warehouses[0]
		}
	}

	if wareHouse == nil {
		return nil, model.NewPreorderAllocationError(preorderAllocation.GetOrderLine()), nil
	}

	stocks, appErr := s.srv.WarehouseService().StocksByOption(&model.StockFilterOption{
		LockForUpdate: true,
		ForUpdateOf:   model.StockTableName,
		Conditions: squirrel.And{
			squirrel.Eq{model.ProductVariantTableName + ".Id": productVariant.Id},
			squirrel.Eq{model.WarehouseTableName + ".Id": wareHouse.Id},
		},
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, nil, appErr
		}
	}

	if len(stocks) != 0 {
		return stocks[0], nil, nil
	}

	return &model.Stock{
		WarehouseID:      wareHouse.Id,
		ProductVariantID: productVariant.Id,
		Quantity:         0,
	}, nil, nil
}

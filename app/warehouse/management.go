package warehouse

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/model/warehouse"
)

// Allocate stocks for given `order_lines` in given country.
//
// Function lock for update all stocks and allocations for variants in
// given country and order by pk. Next, generate the dictionary
// ({"stock_pk": "quantity_allocated"}) with actual allocated quantity for stocks.
// Iterate by stocks and allocate as many items as needed or available in stock
// for order line, until allocated all required quantity for the order line.
// If there is less quantity in stocks then rise InsufficientStock exception.
func (a *AppWarehouse) AllocateStocks(orderLineInfos []*order.OrderLineData, countryCode string, channelSlug string) *model.AppError {
	// allocation only applied to order lines with variants with track inventory set to True
	orderLineInfos = a.GetOrderLinesWithTrackInventory(orderLineInfos)
	if len(orderLineInfos) == 0 {
		return nil
	}

	variants := order.OrderLineDatas(orderLineInfos).Variants()
}

// IncreaseAllocations ncrease allocation for order lines with appropriate quantity
func (a *AppWarehouse) IncreaseAllocations(lineInfos []*order.OrderLineData, channelSlug string) *model.AppError {
	var orderLineIDs []string
	for _, lineInfo := range lineInfos {
		orderLineIDs = append(orderLineIDs, lineInfo.Line.Id)
	}

	panic("not implemented")
}

// DecreaseAllocations Decreate allocations for provided order lines.
func (a *AppWarehouse) DecreaseAllocations(lineInfos []*order.OrderLineData) *model.AppError {
	trackedOrderLines := a.GetOrderLinesWithTrackInventory(lineInfos)
	if len(trackedOrderLines) == 0 {
		return nil
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
func (a *AppWarehouse) DecreaseStock(orderLineInfos []*order.OrderLineData, updateStocks bool) *model.AppError {
	panic("not implemented")
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
func (a *AppWarehouse) IncreaseStock(orderLine *order.OrderLine, warehouse *warehouse.WareHouse, quantity int, allocate bool) *model.AppError {
	panic("not implemented")
}

// GetOrderLinesWithTrackInventory Return order lines with variants with track inventory set to True
func (a *AppWarehouse) GetOrderLinesWithTrackInventory(orderLineInfos []*order.OrderLineData) []*order.OrderLineData {

	for i, lineInfo := range orderLineInfos {
		if lineInfo.Variant == nil || !*lineInfo.Variant.TrackInventory {
			orderLineInfos = append(orderLineInfos[:i], orderLineInfos[i:]...)
		}
	}

	return orderLineInfos
}

// DeallocateStock Deallocate stocks for given `order_lines`.
//
// Function lock for update stocks and allocations related to given `order_lines`.
// Iterate over allocations sorted by `stock.pk` and deallocate as many items
// as needed of available in stock for order line, until deallocated all required
// quantity for the order line. If there is less quantity in stocks then
// raise an exception.
func (a *AppWarehouse) DeallocateStock(orderLineDatas []*order.OrderLineData) (*warehouse.AllocationError, *model.AppError) {
	panic("not implemented")
}

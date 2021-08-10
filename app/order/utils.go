package order

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
)

// FilterOrdersByOptions is common method for filtering orders by given option
func (a *AppOrder) FilterOrdersByOptions(option *order.OrderFilterOption) ([]*order.Order, *model.AppError) {
	orders, err := a.Srv().Store.Order().FilterByOption(option)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("FilterOrdersbyOption", "app.order.error_finding_orders_by_option.app_error", err)
	}

	return orders, nil
}

// UpdateOrderStatus Update order status depending on fulfillments
func (a *AppOrder) UpdateOrderStatus(ord *order.Order) *model.AppError {
	totalQuantity, quantityFulfilled, quantityReturned, appErr := a.calculateQuantityIncludingReturns(ord)
	if appErr != nil {
		return appErr
	}

	var status string
	if totalQuantity == 0 {
		status = ord.Status
	} else if quantityFulfilled <= 0 {
		status = order.UNFULFILLED
	} else if quantityReturned > 0 && quantityReturned < totalQuantity {
		status = order.PARTIALLY_RETURNED
	} else if quantityReturned == totalQuantity {
		status = order.RETURNED
	} else if quantityFulfilled < totalQuantity {
		status = order.PARTIALLY_FULFILLED
	} else {
		status = order.FULFILLED
	}

	if status != ord.Status {
		ord.Status = status
		_, appErr := a.UpsertOrder(ord)
		if appErr != nil {
			return appErr
		}
	}

	return nil
}

func (a *AppOrder) calculateQuantityIncludingReturns(ord *order.Order) (uint, uint, uint, *model.AppError) {
	orderLinesOfOrder, appErr := a.GetAllOrderLinesByOrderId(ord.Id)
	if appErr != nil {
		return 0, 0, 0, appErr
	}

	var (
		totalOrderLinesQuantity uint
		quantityFulfilled       uint
		quantityReturned        uint
		quantityReplaced        uint
	)

	for _, line := range orderLinesOfOrder {
		totalOrderLinesQuantity += line.Quantity
		quantityFulfilled += line.QuantityFulfilled
	}

	fulfillmentsOfOrder, appErr := a.FulfillmentsByOrderID(ord.Id)
	if appErr != nil {
		return 0, 0, 0, appErr
	}

	var (
		hasGoRutines bool
		appError     *model.AppError
	)

	for _, fulfillment := range fulfillmentsOfOrder {
		if status := fulfillment.Status; util.StringInSlice(status, []string{
			order.FULFILLMENT_RETURNED,
			order.FULFILLMENT_REFUNDED_AND_RETURNED,
			order.FULFILLMENT_REPLACED,
		}) {

			a.wg.Add(1)
			hasGoRutines = true

			go func(fulm *order.Fulfillment) {
				fulfillmentLinesOfFulfillment, apErr := a.FulfillmentLinesByFulfillmentID(fulm.Id)

				a.mutex.Lock()
				if appError != nil && appError == nil {
					appError = apErr
				} else {
					for _, line := range fulfillmentLinesOfFulfillment {
						if status == order.FULFILLMENT_RETURNED || status == order.FULFILLMENT_REFUNDED_AND_RETURNED {
							quantityReturned += line.Quantity
						} else {
							quantityReplaced += line.Quantity
						}
					}
				}
				a.mutex.Unlock()
				a.wg.Done()

			}(fulfillment)

		}
	}

	if hasGoRutines {
		a.wg.Wait()
	}

	if appError != nil {
		return 0, 0, 0, appError
	}

	totalOrderLinesQuantity -= quantityReplaced
	quantityFulfilled -= quantityReplaced

	return totalOrderLinesQuantity, quantityFulfilled, quantityReturned, nil
}

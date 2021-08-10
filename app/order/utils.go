package order

import (
	"net/http"

	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app/discount"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/model/shop"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
)

// GetOrderCountry Return country to which order will be shipped
func (a *AppOrder) GetOrderCountry(ord *order.Order) (string, *model.AppError) {
	addressID := ord.BillingAddressID
	orderRequireShipping, appErr := a.OrderShippingIsRequired(ord.Id)
	if appErr != nil {
		return "", appErr
	}
	if orderRequireShipping {
		addressID = ord.ShippingAddressID
	}

	if addressID == nil {
		return model.DEFAULT_COUNTRY, nil
	}

	address, appErr := a.AccountApp().AddressById(*addressID)
	if appErr != nil {
		return "", appErr
	}

	return address.Country, nil
}

// OrderLineNeedsAutomaticFulfillment Check if given line is digital and should be automatically fulfilled.
func (a *AppOrder) OrderLineNeedsAutomaticFulfillment(orderLine *order.OrderLine, shopDigitalSettings *shop.ShopDefaultDigitalContentSettings) (bool, *model.AppError) {
	if orderLine.VariantID == nil {
		return false, nil
	}

	digitalContentOfOrderLineProductVariant, appErr := a.ProductApp().DigitalContentByProductVariantID(*orderLine.VariantID)
	if appErr != nil {
		return false, appErr
	}

	if *digitalContentOfOrderLineProductVariant.UseDefaultSettings && *shopDigitalSettings.AutomaticFulfillmentDigitalProducts {
		return true, nil
	}
	if *digitalContentOfOrderLineProductVariant.AutomaticFulfillment {
		return true, nil
	}

	return false, nil
}

// OrderNeedsAutomaticFulfillment checks if given order has digital products which shoul be automatically fulfilled.
func (a *AppOrder) OrderNeedsAutomaticFulfillment(ord *order.Order) (bool, *model.AppError) {
	// finding shop that hold this order:
	ownerShopOfOrder, appErr := a.ShopApp().ShopById(ord.ShopID)
	if appErr != nil {
		return false, appErr
	}
	shopDefaultDigitalContentSettings := a.ProductApp().GetDefaultDigitalContentSettings(ownerShopOfOrder)

	orderLinesOfOrder, appErr := a.GetAllOrderLinesByOrderId(ord.Id)
	if appErr != nil {
		return false, appErr
	}

	for _, orderLine := range orderLinesOfOrder {
		orderLineNeedsAutomaticFulfillment, appErr := a.OrderLineNeedsAutomaticFulfillment(orderLine, shopDefaultDigitalContentSettings)
		if appErr != nil {
			return false, appErr
		}
		if orderLineNeedsAutomaticFulfillment {
			return true, nil
		}
	}

	return false, nil
}

func (a *AppOrder) GetVoucherDiscountAssignedToOrder(ord *order.Order) (*product_and_discount.OrderDiscount, *model.AppError) {
	orderDiscountsOfOrder, appErr := a.DiscountApp().
		OrderDiscountsByOption(&product_and_discount.OrderDiscountFilterOption{
			Type: &model.StringFilter{
				StringOption: &model.StringOption{
					Eq: product_and_discount.VOUCHER,
				},
			},
		})

	if appErr != nil {
		return nil, appErr
	}

	return orderDiscountsOfOrder[0], nil
}

// Recalculate all order discounts assigned to order.
//
// It returns the list of tuples which contains order discounts where the amount has been changed.
func (a *AppOrder) RecalculateOrderDiscounts(ord *order.Order) ([][2]*product_and_discount.OrderDiscount, *model.AppError) {
	var changedOrderDiscounts [][2]*product_and_discount.OrderDiscount

	orderDiscounts, appErr := a.DiscountApp().
		OrderDiscountsByOption(&product_and_discount.OrderDiscountFilterOption{
			OrderID: &model.StringFilter{
				StringOption: &model.StringOption{
					Eq: ord.Id,
				},
			},
			Type: &model.StringFilter{
				StringOption: &model.StringOption{
					Eq: product_and_discount.MANUAL,
				},
			},
		})

	if appErr != nil {
		return nil, appErr
	}

	for _, orderDiscount := range orderDiscounts {

		previousOrderDiscount := orderDiscount.DeepCopy()
		currentTotal := ord.Total.Gross.Amount

		appErr = a.UpdateOrderDiscountForOrder(ord, orderDiscount, "", "", nil)
		if appErr != nil {
			return nil, appErr
		}

		discountValue := orderDiscount.Value
		amount := orderDiscount.Amount

		if (orderDiscount.ValueType == product_and_discount.PERCENTAGE || currentTotal.LessThan(*discountValue)) &&
			!amount.Amount.Equal(*previousOrderDiscount.Amount.Amount) {
			changedOrderDiscounts = append(changedOrderDiscounts, [2]*product_and_discount.OrderDiscount{
				previousOrderDiscount,
				orderDiscount,
			})
		}
	}

	return changedOrderDiscounts, nil
}

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

// UpdateOrderDiscountForOrder Update the order_discount for an order and recalculate the order's prices
//
// `reason`, `valueType` and `value` can be nil
func (a *AppOrder) UpdateOrderDiscountForOrder(ord *order.Order, orderDiscountToUpdate *product_and_discount.OrderDiscount, reason string, valueType string, value *decimal.Decimal) *model.AppError {
	if value == nil {
		value = orderDiscountToUpdate.Value
	}
	if valueType == "" {
		valueType = orderDiscountToUpdate.ValueType
	}

	if reason != "" {
		orderDiscountToUpdate.Reason = &reason
	}

	netTotal, err := ApplyDiscountToValue(value, valueType, orderDiscountToUpdate.Currency, ord.Total.Net)
	if err != nil {
		return model.NewAppError("UpdateOrderDiscountForOrder", "app.order.error_calculating_net_total_discount.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	grossTotal, err := ApplyDiscountToValue(value, valueType, orderDiscountToUpdate.Currency, ord.Total.Gross)
	if err != nil {
		return model.NewAppError("UpdateOrderDiscountForOrder", "app.order.error_calculating_gross_total_discount.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	sub, _ := ord.Total.Sub(grossTotal)

	orderDiscountToUpdate.Amount = sub.Gross
	orderDiscountToUpdate.Value = value
	orderDiscountToUpdate.ValueType = valueType

	ord.Total, _ = goprices.NewTaxedMoney(netTotal.(*goprices.Money), grossTotal.(*goprices.Money))

	_, appErr := a.DiscountApp().UpsertOrderDiscount(orderDiscountToUpdate)

	return appErr
}

// ApplyDiscountToValue Calculate the price based on the provided values
func ApplyDiscountToValue(value *decimal.Decimal, valueType string, currency string, priceToDiscount interface{}) (interface{}, error) {
	// validate currency
	money, err := goprices.NewMoney(value, currency)
	if err != nil {
		return nil, err
	}

	var discountCalculator discount.DiscountCalculator
	if valueType == product_and_discount.FIXED {
		discountCalculator = discount.Decorator(money)
	} else {
		discountCalculator = discount.Decorator(value)
	}

	return discountCalculator(priceToDiscount)
}

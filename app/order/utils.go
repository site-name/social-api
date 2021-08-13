package order

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/discount"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/giftcard"
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/model/shop"
	"github.com/sitename/sitename/modules/measurement"
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

	digitalOrderLinesOfOrder, appErr := a.AllDigitalOrderLinesOfOrder(ord.Id)
	if appErr != nil {
		return false, appErr
	}

	for _, orderLine := range digitalOrderLinesOfOrder {
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

		if (orderDiscount.ValueType == product_and_discount.PERCENTAGE || currentTotal.LessThan(*discountValue)) && !amount.Amount.Equal(*previousOrderDiscount.Amount.Amount) {
			changedOrderDiscounts = append(changedOrderDiscounts, [2]*product_and_discount.OrderDiscount{
				previousOrderDiscount,
				orderDiscount,
			})
		}
	}

	return changedOrderDiscounts, nil
}

func (a *AppOrder) RecalculateOrderPrices(ord *order.Order, kwargs map[string]interface{}) *model.AppError {
	// TODO: fix me
	panic("not implemented")
}

func (a *AppOrder) RecalculateOrder(ord *order.Order) {
	panic("not implemented")
}

// ReCalculateOrderWeight
func (a *AppOrder) ReCalculateOrderWeight(ord *order.Order) *model.AppError {
	orderLines, appErr := a.GetAllOrderLinesByOrderId(ord.Id)
	if appErr != nil {
		return appErr
	}

	var (
		appError      *model.AppError
		hasGoRoutines bool
		weight        = measurement.ZeroWeight
	)

	setAppError := func(err *model.AppError) {
		a.mutex.Lock()
		if err != nil && appError == nil {
			appError = err
		}
		a.mutex.Unlock()
	}

	for _, orderLine := range orderLines {
		if orderLine.VariantID != nil && model.IsValidId(*orderLine.VariantID) {

			hasGoRoutines = true
			a.wg.Add(1)

			go func(variantID string) {
				productVariantWeight, err := a.Srv().Store.ProductVariant().GetWeight(*orderLine.VariantID)
				if err != nil {
					if _, ok := err.(*store.ErrNotFound); !ok {
						setAppError(model.NewAppError("ReCalculateOrderWeight", app.InternalServerErrorID, nil, err.Error(), http.StatusInternalServerError))
					}
				} else {
					a.mutex.Lock()
					addedWeight, err := weight.Add(productVariantWeight.Mul(float32(orderLine.Quantity)))
					if err != nil {
						setAppError(model.NewAppError("ReCalculateOrderWeight", app.InternalServerErrorID, nil, err.Error(), http.StatusInternalServerError))
					} else {
						weight = addedWeight
					}
					a.mutex.Unlock()
				}

				a.wg.Done()
			}(*orderLine.VariantID)

		}
	}

	if hasGoRoutines {
		a.wg.Wait()
	}

	if appError != nil {
		return appError
	}

	weight, _ = weight.ConvertTo(ord.WeightUnit)
	ord.WeightAmount = *weight.Amount

	_, appError = a.UpsertOrder(ord)
	return appError
}

func (a *AppOrder) UpdateTaxesForOrderLine() {
	panic("not implemented")
}

func (a *AppOrder) UpdateTaxesForOrderLines() {
	panic("not implemented")
}

func (a *AppOrder) UpdateOrderPrices() {
	panic("not implemented")
}

// thereIsAnItem takes a slice and a checker function.
// it iterates through the slice to find out if there is an item that satisfy given checker function
func thereIsAnItem(slice interface{}, checker func(item interface{}) bool) bool {
	valueOf := reflect.ValueOf(slice)
	typeOf := reflect.TypeOf(slice)

	if typeOf.Kind() == reflect.Slice {
		for i := 0; i < valueOf.Len(); i++ {
			valueAtIndex := valueOf.Index(i)
			if checker(valueAtIndex.Interface()) {
				return true
			}
		}
	}

	return false
}

// collectionsIntersection select only common items between two given slices
func collectionsIntersection(
	collectionSlice1 []*product_and_discount.Collection,
	collectionSlice2 []*product_and_discount.Collection,
) []*product_and_discount.Collection {

	var res []*product_and_discount.Collection

	for i := 0; i < len(collectionSlice1); i++ {
		for j := 0; j < len(collectionSlice2); j++ {
			if collectionSlice1[i].Id == collectionSlice2[j].Id {
				res = append(res, collectionSlice1[i])
			}
		}
	}

	return res
}

func (a *AppOrder) GetDiscountedLines(orderLines []*order.OrderLine, voucher *product_and_discount.Voucher) ([]*order.OrderLine, *model.AppError) {
	var (
		discountedProducts    []*product_and_discount.Product
		discountedCategories  []*product_and_discount.Category
		discountedCollections []*product_and_discount.Collection
		firstAppError         *model.AppError
		meetMap               = map[string]bool{}
	)

	setFirstAppErr := func(err *model.AppError) {
		a.mutex.Lock()
		if err != nil {
			firstAppError = err
		}
		a.mutex.Unlock()
	}

	a.wg.Add(3)

	go func() {
		products, appErr := a.ProductApp().ProductsByVoucherID(voucher.Id)
		if appErr != nil {
			setFirstAppErr(appErr)
		} else {
			discountedProducts = products
		}
		a.wg.Done()
	}()

	go func() {
		categories, appErr := a.ProductApp().CategoriesByVoucherID(voucher.Id)
		if appErr != nil {
			setFirstAppErr(appErr)
		} else {
			// remove duplicate categories
			for _, category := range categories {
				if _, met := meetMap[category.Id]; !met {
					discountedCategories = append(discountedCategories, category)
					meetMap[category.Id] = true
				}
			}
		}
		a.wg.Done()
	}()

	go func() {
		collections, appErr := a.ProductApp().CollectionsByVoucherID(voucher.Id)
		if appErr != nil {
			setFirstAppErr(appErr)
		} else {
			// remove duplicate collections
			for _, collection := range collections {
				if _, met := meetMap[collection.Id]; !met {
					discountedCollections = append(discountedCollections, collection)
					meetMap[collection.Id] = true
				}
			}
		}
		a.wg.Done()
	}()

	a.wg.Wait()

	// returns immediately if there is an system error occured
	if firstAppError != nil {
		return nil, firstAppError
	}

	var (
		discountedOrderLines []*order.OrderLine
		appError             *model.AppError
		hasGoRoutines        bool
	)
	setAppError := func(appErr *model.AppError) {
		a.mutex.Lock()
		if appErr != nil && appError == nil {
			appError = appErr
		}
		a.mutex.Unlock()
	}

	if len(discountedProducts) > 0 || len(discountedCategories) > 0 || len(discountedCollections) > 0 {

		for _, orderLine := range orderLines {
			// we can
			if orderLine.VariantID != nil && model.IsValidId(*orderLine.VariantID) {
				hasGoRoutines = true
				a.wg.Add(1)

				go func(anOrderLine *order.OrderLine) {
					orderLineProduct, appErr := a.ProductApp().ProductByProductVariantID(*anOrderLine.VariantID)
					if appErr != nil {
						setAppError(appErr)
					} else {
						orderLineCategory, appErr := a.ProductApp().CategoryByProductID(orderLineProduct.Id)
						if appErr != nil {
							setAppError(appErr)
						} else {
							orderLineCollections, appErr := a.ProductApp().CollectionsByProductID(orderLineProduct.Id)
							if appErr != nil {
								setAppError(appErr)
							} else {
								orderLineProductInDiscountedProducts := thereIsAnItem(discountedProducts, func(i interface{}) bool { return i.(*product_and_discount.Product).Id == orderLineProduct.Id })
								orderLineCategoryInDiscountedCategories := thereIsAnItem(discountedCategories, func(i interface{}) bool { return i.(*product_and_discount.Category).Id == orderLineCategory.Id })
								orderLineCollectionsIntersectDiscountedCollections := collectionsIntersection(orderLineCollections, discountedCollections)

								if orderLineProductInDiscountedProducts || orderLineCategoryInDiscountedCategories || len(orderLineCollectionsIntersectDiscountedCollections) > 0 {
									a.mutex.Lock()
									discountedOrderLines = append(discountedOrderLines, anOrderLine)
									a.mutex.Unlock()
								}
							}
						}
					}

					a.wg.Done()
				}(orderLine)
			}
		}
	} else {
		// If there's no discounted products, collections or categories,
		// it means that all products are discounted
		return orderLines, nil
	}

	if hasGoRoutines {
		a.wg.Wait()
	}

	return discountedOrderLines, nil
}

// Get prices of variants belonging to the discounted specific products.
//
// Specific products are products, collections and categories.
// Product must be assigned directly to the discounted category, assigning
// product to child category won't work
func (a *AppOrder) GetPricesOfDiscountedSpecificProduct(orderLines []*order.OrderLine, voucher *product_and_discount.Voucher) ([]*goprices.Money, *model.AppError) {
	discountedOrderLines, appErr := a.GetDiscountedLines(orderLines, voucher)
	if appErr != nil {
		return nil, appErr
	}

	var orderLinePrices []*goprices.Money
	for _, orderLine := range discountedOrderLines {
		if orderLine.Quantity == 0 {
			continue
		}
		for i := 0; i < int(orderLine.Quantity); i++ {
			orderLinePrices = append(orderLinePrices, orderLine.UnitPriceGross)
		}
	}

	return orderLinePrices, nil
}

// Calculate discount value depending on voucher and discount types.
//
// Raise NotApplicable if voucher of given type cannot be applied.
func (a *AppOrder) GetVoucherDiscountForOrder(ord *order.Order) (interface{}, *model.AppError) {
	ord.PopulateNonDbFields()

	// validate if order has voucher attached to
	if ord.VoucherID == nil || !model.IsValidId(*ord.VoucherID) {
		return &goprices.Money{
			Amount:   &decimal.Zero,
			Currency: ord.Currency,
		}, nil
	}

	appErr := a.DiscountApp().ValidateVoucherInOrder(ord)
	if appErr != nil {
		return nil, appErr
	}

	orderLines, appErr := a.GetAllOrderLinesByOrderId(ord.Id)
	if appErr != nil {
		return nil, appErr
	}

	orderSubTotal, appErr := a.PaymentApp().GetSubTotal(orderLines, ord.Currency)
	if appErr != nil {
		return nil, appErr
	}

	voucherOfDiscount, appErr := a.DiscountApp().VoucherById(*ord.VoucherID)
	if appErr != nil {
		return nil, appErr
	}

	if voucherOfDiscount.Type == product_and_discount.ENTIRE_ORDER {
		return a.DiscountApp().GetDiscountAmountFor(voucherOfDiscount, orderSubTotal.Gross, ord.ChannelID)
	}
	if voucherOfDiscount.Type == product_and_discount.SHIPPING {
		return a.DiscountApp().GetDiscountAmountFor(voucherOfDiscount, ord.ShippingPrice, ord.ChannelID)
	}
	// otherwise: Type is product_and_discount.SPECIFIC_PRODUCT
	prices, appErr := a.GetPricesOfDiscountedSpecificProduct(orderLines, voucherOfDiscount)
	if appErr != nil {
		return nil, appErr
	}
	if len(prices) == 0 {
		return nil, model.NewAppError("GetVoucherDiscountForOrder", "app.order.offer_only_valid_for_selected_items.app_error", nil, "", http.StatusNotAcceptable)
	}

	return a.DiscountApp().GetProductsVoucherDiscount(voucherOfDiscount, prices, ord.ChannelID)
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
			if !hasGoRutines {
				hasGoRutines = true
			}

			go func(fulm *order.Fulfillment) {
				fulfillmentLinesOfFulfillment, apErr := a.FulfillmentLinesByFulfillmentID(fulm.Id)

				a.mutex.Lock()
				if apErr != nil && appError == nil {
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

func (a *AppOrder) AddVariantToOrder() {
	panic("not implemented")
}

// Add gift card to order.
//
// Return a total price left after applying the gift cards.
func (a *AppOrder) AddGiftCardToOrder(ord *order.Order, giftCard *giftcard.GiftCard, totalPriceLeft *goprices.Money) (*goprices.Money, *model.AppError) {
	// validate given arguments's currencies are valid
	_, err := goprices.GetCurrencyPrecision(totalPriceLeft.Currency)
	if err != nil || !strings.EqualFold(giftCard.Currency, totalPriceLeft.Currency) {
		return nil, model.NewAppError("AddGiftCardToOrder", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "totalPriceLeft"}, err.Error(), http.StatusBadRequest)
	}

	giftCard.PopulateNonDbFields()
	// add new order-giftcard relationship
	if totalPriceLeft.Amount.GreaterThan(decimal.Zero) {
		_, appErr := a.GiftcardApp().CreateOrderGiftcardRelation(&giftcard.OrderGiftCard{
			GiftCardID: giftCard.Id,
			OrderID:    ord.Id,
		})
		if appErr != nil {
			return nil, appErr
		}

		if less, err := totalPriceLeft.LessThan(giftCard.CurrentBalance); less && err == nil {
			giftCard.CurrentBalance, _ = giftCard.CurrentBalance.Sub(totalPriceLeft)
			totalPriceLeft, _ = util.ZeroMoney(totalPriceLeft.Currency)
		} else {
			totalPriceLeft, _ = totalPriceLeft.Sub(giftCard.CurrentBalance)
			giftCard.CurrentBalanceAmount = &decimal.Zero
		}

		giftCard.LastUsedOn = model.GetMillis()
		_, appErr = a.GiftcardApp().UpsertGiftcard(giftCard)
		if appErr != nil {
			return nil, appErr
		}
	}

	return totalPriceLeft, nil
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

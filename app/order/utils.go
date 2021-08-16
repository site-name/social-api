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
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/giftcard"
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/model/shipping"
	"github.com/sitename/sitename/model/shop"
	"github.com/sitename/sitename/model/warehouse"
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
//
// NOTE: before calling this, caller can attach related data into `orderLine` so this function does not have to call the database
func (a *AppOrder) OrderLineNeedsAutomaticFulfillment(orderLine *order.OrderLine, shopDigitalSettings *shop.ShopDefaultDigitalContentSettings) (bool, *model.AppError) {
	if orderLine.VariantID == nil || orderLine.ProductVariant == nil {
		return false, nil
	}

	digitalContent := orderLine.ProductVariant.DigitalContent
	var appErr *model.AppError

	if digitalContent == nil {
		digitalContent, appErr = a.ProductApp().DigitalContentByProductVariantID(*orderLine.VariantID)
		if appErr != nil {
			if appErr.StatusCode == http.StatusInternalServerError {
				appErr.Where = "OrderLineNeedsAutomaticFulfillment"
				return false, appErr
			}
			return false, nil
		}
	}

	if *digitalContent.UseDefaultSettings && *shopDigitalSettings.AutomaticFulfillmentDigitalProducts {
		return true, nil
	}
	if *digitalContent.AutomaticFulfillment {
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

	// get first item of the result here. make sure to ordering the query
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

// func (a *AppOrder) RecalculateOrderPrices(ord *order.Order, kwargs map[string]interface{}) *model.AppError {
// 	// TODO: fix me
// 	panic("not implemented")
// }

// Recalculate and assign total price of order.
//
// Total price is a sum of items in order and order shipping price minus
// discount amount.
//
// Voucher discount amount is recalculated by default. To avoid this, pass
// update_voucher_discount argument set to False.
func (a *AppOrder) RecalculateOrder(ord *order.Order, kwargs map[string]interface{}) (appErr *model.AppError) {
	defer func() {
		if appErr != nil {
			appErr.Where = "RecalculateOrder"
		}
	}()

	appErr = a.RecalculateOrderPrices(ord, kwargs)
	if appErr != nil {
		return
	}

	changedOrderDiscounts, appErr := a.RecalculateOrderDiscounts(ord)
	if appErr != nil {
		return
	}

	appErr = a.OrderDiscountsAutomaticallyUpdatedEvent(ord, changedOrderDiscounts)
	if appErr != nil {
		return
	}

	ord, appErr = a.UpsertOrder(ord)
	if appErr != nil {
		return
	}

	return a.ReCalculateOrderWeight(ord)
}

// ReCalculateOrderWeight
func (a *AppOrder) ReCalculateOrderWeight(ord *order.Order) *model.AppError {
	orderLines, appErr := a.OrderLinesByOption(&order.OrderLineFilterOption{
		OrderID: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: ord.Id,
			},
		},
	})
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

			go func(anOrderLine *order.OrderLine) {
				productVariantWeight, err := a.Srv().Store.ProductVariant().GetWeight(*anOrderLine.VariantID)
				if err != nil {
					if _, ok := err.(*store.ErrNotFound); !ok { // set appError if the error is caused by system.
						setAppError(model.NewAppError("ReCalculateOrderWeight", app.InternalServerErrorID, nil, err.Error(), http.StatusInternalServerError))
						// Ignore If not variant found
					}
				} else {
					a.mutex.Lock()
					addedWeight, err := weight.Add(productVariantWeight.Mul(float32(anOrderLine.Quantity)))
					if err != nil {
						setAppError(model.NewAppError("ReCalculateOrderWeight", app.InternalServerErrorID, nil, err.Error(), http.StatusInternalServerError))
					} else {
						weight = addedWeight
					}
					a.mutex.Unlock()
				}

				a.wg.Done()
			}(orderLine)

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

// GetDiscountedLines returns a list of discounted order lines, filterd from given orderLines
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
	if ord.VoucherID == nil {
		return &goprices.Money{
			Amount:   &decimal.Zero,
			Currency: ord.Currency,
		}, nil
	}

	appErr := a.DiscountApp().ValidateVoucherInOrder(ord)
	if appErr != nil {
		return nil, appErr
	}

	orderLines, appErr := a.OrderLinesByOption(&order.OrderLineFilterOption{
		OrderID: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: ord.Id,
			},
		},
	})
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

func (a *AppOrder) calculateQuantityIncludingReturns(ord *order.Order) (int, int, int, *model.AppError) {
	orderLinesOfOrder, appErr := a.OrderLinesByOption(&order.OrderLineFilterOption{
		OrderID: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: ord.Id,
			},
		},
	})
	if appErr != nil {
		return 0, 0, 0, appErr
	}

	var (
		totalOrderLinesQuantity int
		quantityFulfilled       int
		quantityReturned        int
		quantityReplaced        int
	)

	for _, line := range orderLinesOfOrder {
		totalOrderLinesQuantity += line.Quantity
		quantityFulfilled += line.QuantityFulfilled
	}

	fulfillmentsOfOrder, appErr := a.FulfillmentsByOption(&order.FulfillmentFilterOption{
		OrderID: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: ord.Id,
			},
		},
	})
	if appErr != nil {
		return 0, 0, 0, appErr
	}

	// filter all fulfillments that has `Status` is either: "returned", "refunded_and_returned" and "replaced"
	var (
		filteredFulfillmentIDs []string
		fulfillmentMap         = map[string]*order.Fulfillment{}
	)
	for _, fulfillment := range fulfillmentsOfOrder {
		if util.StringInSlice(fulfillment.Status, []string{
			order.FULFILLMENT_RETURNED,
			order.FULFILLMENT_REFUNDED_AND_RETURNED,
			order.FULFILLMENT_REPLACED,
		}) {
			filteredFulfillmentIDs = append(filteredFulfillmentIDs, fulfillment.Id)
			fulfillmentMap[fulfillment.Id] = fulfillment
		}
	}

	// finds all fulfillment lines belong to filtered fulfillments
	fulfillmentLines, appErr := a.FulfillmentLinesByOption(&order.FulfillmentLineFilterOption{
		FulfillmentID: &model.StringFilter{
			StringOption: &model.StringOption{
				In: filteredFulfillmentIDs,
			},
		},
	})
	if appErr != nil {
		return 0, 0, 0, appErr
	}

	for _, fulfillmentLine := range fulfillmentLines {
		parentFulfillmentStatus := fulfillmentMap[fulfillmentLine.FulfillmentID].Status

		if parentFulfillmentStatus == order.FULFILLMENT_RETURNED || parentFulfillmentStatus == order.FULFILLMENT_REFUNDED_AND_RETURNED {
			quantityReturned += fulfillmentLine.Quantity
		} else if parentFulfillmentStatus == order.FULFILLMENT_REPLACED {
			quantityReplaced += fulfillmentLine.Quantity
		}
	}

	totalOrderLinesQuantity -= quantityReplaced
	quantityFulfilled -= quantityReplaced

	return totalOrderLinesQuantity, quantityFulfilled, quantityReturned, nil
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

	// NOTE: must call this before performing any operations on giftcards
	giftCard.PopulateNonDbFields()

	// add new order-giftcard relationship
	if totalPriceLeft.Amount.GreaterThan(decimal.Zero) {
		// create new order-goftcard relation instance
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

		// update giftcard
		giftCard.LastUsedOn = model.GetMillis()
		_, appErr = a.GiftcardApp().UpsertGiftcard(giftCard)
		if appErr != nil {
			return nil, appErr
		}
	}

	return totalPriceLeft, nil
}

func (a *AppOrder) updateAllocationsForLine(lineInfo *order.OrderLineData, oldQuantity int, newQuantity int, channelSlug string) *model.AppError {
	if oldQuantity == newQuantity {
		return nil
	}

	orderLinesWithTrackInventory := a.WarehouseApp().GetOrderLinesWithTrackInventory([]*order.OrderLineData{lineInfo})
	if len(orderLinesWithTrackInventory) == 0 {
		return nil
	}

	if oldQuantity < newQuantity {
		lineInfo.Quantity = newQuantity - oldQuantity
		return a.WarehouseApp().IncreaseAllocations([]*order.OrderLineData{lineInfo}, channelSlug)
	} else {
		lineInfo.Quantity = oldQuantity - newQuantity
		return a.WarehouseApp().DecreaseAllocations([]*order.OrderLineData{lineInfo})
	}
}

// ChangeOrderLineQuantity Change the quantity of ordered items in a order line.
//
// NOTE: userID can be empty
func (a *AppOrder) ChangeOrderLineQuantity(userID string, lineInfo *order.OrderLineData, oldQuantity int, newQuantity int, channelSlug string, sendEvent bool) *model.AppError {

	orderLine := lineInfo.Line
	// NOTE: this must be called
	orderLine.PopulateNonDbFields()

	if newQuantity > 0 {
		order, appErr := a.OrderById(lineInfo.Line.OrderID)
		if appErr != nil {
			appErr.Where = "ChangeOrderLineQuantity"
			return appErr
		}

		if order.IsUnconfirmed() {
			appErr = a.updateAllocationsForLine(lineInfo, oldQuantity, newQuantity, channelSlug)
			if appErr != nil {
				appErr.Where = "ChangeOrderLineQuantity"
				return appErr
			}
		}

		lineInfo.Line.Quantity = newQuantity

		totalPriceNetAmount := orderLine.UnitPriceNetAmount.Mul(decimal.NewFromInt32(int32(orderLine.Quantity)))
		totalPriceGrossAmount := orderLine.UnitPriceGrossAmount.Mul(decimal.NewFromInt32(int32(orderLine.Quantity)))
		orderLine.TotalPriceNetAmount = model.NewDecimal(totalPriceNetAmount.Round(3))
		orderLine.TotalPriceGrossAmount = model.NewDecimal(totalPriceGrossAmount.Round(3))

		unDiscountedTotalPriceNetAmount := orderLine.UnDiscountedUnitPriceNetAmount.Mul(decimal.NewFromInt32(int32(orderLine.Quantity)))
		unDiscountedTotalpriceGrossAmount := orderLine.UnDiscountedUnitPriceGrossAmount.Mul(decimal.NewFromInt32(int32(orderLine.Quantity)))
		orderLine.UnDiscountedTotalPriceNetAmount = model.NewDecimal(unDiscountedTotalPriceNetAmount.Round(3))
		orderLine.UnDiscountedTotalPriceGrossAmount = model.NewDecimal(unDiscountedTotalpriceGrossAmount.Round(3))

		_, appErr = a.UpsertOrderLine(&orderLine)
		if appErr != nil {
			appErr.Where = "ChangeOrderLineQuantity"
			return appErr
		}
	} else { // ------------
		appErr := a.DeleteOrderLine(lineInfo)
		if appErr != nil {
			appErr.Where = "ChangeOrderLineQuantity"
			return appErr
		}
	}

	quantityDiff := int(oldQuantity) - int(newQuantity)

	if sendEvent {
		appErr := a.CreateOrderEvent(&orderLine, userID, quantityDiff)
		if appErr != nil {
			appErr.Where = "ChangeOrderLineQuantity"
			return appErr
		}
	}

	return nil
}

func (a *AppOrder) CreateOrderEvent(orderLine *order.OrderLine, userID string, quantityDiff int) *model.AppError {
	var appErr *model.AppError

	var savingUserID *string
	if userID != "" {
		savingUserID = &userID
	}

	if quantityDiff > 0 {
		_, appErr = a.CommonCreateOrderEvent(&order.OrderEventOption{
			OrderID: orderLine.OrderID,
			UserID:  savingUserID,
			Type:    order.ORDER_EVENT_TYPE__REMOVED_PRODUCTS,
			Parameters: &model.StringInterface{
				"lines": linesPerQuantityToLineObjectList([]*QuantityOrderLine{
					{
						Quantity:  quantityDiff,
						OrderLine: orderLine,
					},
				}),
			},
		})
	} else if quantityDiff < 0 {
		_, appErr = a.CommonCreateOrderEvent(&order.OrderEventOption{
			OrderID: orderLine.OrderID,
			UserID:  savingUserID,
			Type:    order.ORDER_EVENT_TYPE__ADDED_PRODUCTS,
			Parameters: &model.StringInterface{
				"lines": linesPerQuantityToLineObjectList([]*QuantityOrderLine{
					{
						Quantity:  quantityDiff * -1,
						OrderLine: orderLine,
					},
				}),
			},
		})
	}

	return appErr
}

// Delete an order line from an order.
func (a *AppOrder) DeleteOrderLine(lineInfo *order.OrderLineData) *model.AppError {
	ord, appErr := a.OrderById(lineInfo.Line.OrderID)
	if appErr != nil {
		appErr.Where = "DeleteOrderLine"
		return appErr
	}

	if ord.IsUnconfirmed() {
		appErr = a.WarehouseApp().DecreaseAllocations([]*order.OrderLineData{lineInfo})
		if appErr != nil {
			return appErr
		}
	}

	return a.DeleteOrderLines([]string{lineInfo.Line.Id})
}

// RestockOrderLines Return ordered products to corresponding stocks
func (a *AppOrder) RestockOrderLines(ord *order.Order) *model.AppError {
	countryCode, appError := a.GetOrderCountry(ord)
	if appError != nil {
		return appError
	}

	warehouses, appError := a.WarehouseApp().WarehouseByOption(&warehouse.WarehouseFilterOption{
		ShippingZonesCountries: &model.StringFilter{
			StringOption: &model.StringOption{
				Like: countryCode,
			},
		},
	})
	if appError != nil {
		appError.Where = "RestockOrderLines"
		return appError
	}
	defaultWarehouse := warehouses[0]

	orderLinesOfOrder, appError := a.OrderLinesByOption(&order.OrderLineFilterOption{
		OrderID: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: ord.Id,
			},
		},
	})
	if appError != nil {
		appError.Where = "RestockOrderLines"
		return appError
	}

	var (
		dellocatingStockLines []*order.OrderLineData
		hasGoRoutines         bool
	)

	setAppError := func(err *model.AppError) {
		if err != nil {
			a.mutex.Lock()
			if appError == nil {
				appError = err
				appError.Where = "RestockOrderLines"
			}
			a.mutex.Unlock()
		}
	}

	for _, orderLine := range orderLinesOfOrder {
		if orderLine.VariantID != nil {

			hasGoRoutines = true
			a.wg.Add(1)

			go func(anOrderLine *order.OrderLine) {
				productVariant, appErr := a.ProductApp().ProductVariantById(*anOrderLine.VariantID)
				if appErr != nil {
					setAppError(appErr) //
				} else {
					if *productVariant.TrackInventory {
						if anOrderLine.QuantityUnFulfilled() > 0 {

							a.mutex.Lock()
							dellocatingStockLines = append(dellocatingStockLines, &order.OrderLineData{
								Line:     *anOrderLine,
								Quantity: anOrderLine.QuantityUnFulfilled(),
							})
							a.mutex.Unlock()

						}

						if anOrderLine.QuantityFulfilled > 0 {
							allocations, appErr := a.WarehouseApp().AllocationsByOption(&warehouse.AllocationFilterOption{
								OrderLineID: &model.StringFilter{
									StringOption: &model.StringOption{
										Eq: anOrderLine.Id,
									},
								},
							})
							if appErr != nil {
								setAppError(appErr) //
							} else {
								warehouse := defaultWarehouse
								if len(allocations) > 0 {
									warehouseOfOrderLine, appErr := a.WarehouseApp().WarehouseByStockID(allocations[0].StockID)
									if appErr != nil {
										setAppError(appErr) //
									} else {
										warehouse = warehouseOfOrderLine
									}
								}

								appErr = a.WarehouseApp().IncreaseStock(anOrderLine, warehouse, anOrderLine.QuantityFulfilled, false)
								setAppError(appErr) //
							}
						}
					}

					if anOrderLine.QuantityFulfilled > 0 {
						anOrderLine.QuantityFulfilled = 0

						_, appErr = a.UpsertOrderLine(anOrderLine)
						setAppError(appErr) //
					}
				}

				a.wg.Done()
			}(orderLine)
		}
	}

	if hasGoRoutines {
		a.wg.Wait()
	}

	if len(dellocatingStockLines) > 0 {
		_, appError = a.WarehouseApp().DeallocateStock(dellocatingStockLines)
	}

	return appError
}

// RestockFulfillmentLines Return fulfilled products to corresponding stocks.
//
// Return products to stocks and update order lines quantity fulfilled values.
func (a *AppOrder) RestockFulfillmentLines(fulfillment *order.Fulfillment, warehouse *warehouse.WareHouse) (appErr *model.AppError) {
	defer func() {
		if appErr != nil {
			appErr.Where = "RestockFulfillmentLines"
		}
	}()

	fulfillmentLines, appErr := a.FulfillmentLinesByOption(&order.FulfillmentLineFilterOption{
		FulfillmentID: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: fulfillment.Id,
			},
		},
	})
	if appErr != nil {
		return appErr
	}

	orderLinesOfFulfillmentLines, appErr := a.OrderLinesByOption(&order.OrderLineFilterOption{
		Id: &model.StringFilter{
			StringOption: &model.StringOption{
				In: order.FulfillmentLines(fulfillmentLines).OrderLineIDs(),
			},
		},
	})
	if appErr != nil {
		return appErr
	}

	// create map with key: fulfillmentLine.Id, value: *OrderLine
	mapFulfillmentLine_OrderLine := map[string]*order.OrderLine{}
	for _, fulfillmentLine := range fulfillmentLines {
		for _, orderLine := range orderLinesOfFulfillmentLines {
			if fulfillmentLine.OrderLineID == orderLine.Id {
				mapFulfillmentLine_OrderLine[fulfillmentLine.Id] = orderLine
			}
		}
	}

	productVariantsOfOrderLines, appErr := a.ProductApp().ProductVariantsByOption(&product_and_discount.ProductVariantFilterOption{
		Id: &model.StringFilter{
			StringOption: &model.StringOption{
				In: order.OrderLines(orderLinesOfFulfillmentLines).ProductVariantIDs(),
			},
		},
	})
	if appErr != nil {
		return appErr
	}

	// create map with key: orderLine.Id, value: *ProductVariant
	mapOrderLine_productVariant := map[string]*product_and_discount.ProductVariant{}
	for _, orderLine := range orderLinesOfFulfillmentLines {
		if orderLine.VariantID == nil { // since some order line have no product variant attached
			continue
		}
		for _, variant := range productVariantsOfOrderLines {
			if variant.Id == *orderLine.VariantID {
				mapOrderLine_productVariant[orderLine.Id] = variant
			}
		}
	}

	for _, fulfillmentLine := range fulfillmentLines {
		orderLineOfFulfillment := mapFulfillmentLine_OrderLine[fulfillmentLine.Id]   // number of order lines == number of fulfillment lines
		variantOfOrderLine := mapOrderLine_productVariant[orderLineOfFulfillment.Id] // variantOfOrderLine can be nil

		if variantOfOrderLine != nil && *variantOfOrderLine.TrackInventory {
			appErr := a.WarehouseApp().IncreaseStock(orderLineOfFulfillment, warehouse, fulfillmentLine.Quantity, true)
			if appErr != nil {
				return appErr
			}
		}

		orderLineOfFulfillment.QuantityFulfilled -= fulfillmentLine.Quantity
	}

	_, appErr = a.BulkUpsertOrderLines(orderLinesOfFulfillmentLines)
	return appErr
}

func (a *AppOrder) SumOrderTotals(orders []*order.Order, currencyCode string) (*goprices.TaxedMoney, *model.AppError) {
	taxedSum, _ := util.ZeroTaxedMoney(currencyCode)
	if len(orders) == 0 {
		return taxedSum, nil
	}
	// validate given `currencyCode` is valid
	currencyCode = strings.ToUpper(currencyCode)
	if _, err := goprices.GetCurrencyPrecision(currencyCode); err != nil || currencyCode != orders[0].Currency {
		return nil, model.NewAppError("SumOrderTotals", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "currencyCode"}, err.Error(), http.StatusBadRequest)
	}

	for _, order := range orders {
		order.PopulateNonDbFields() //
		added, err := taxedSum.Add(order.Total)
		if err != nil {
			return nil, model.NewAppError("SumOrderTotals", "app.order.error_adding_taxed_money.app_error", nil, err.Error(), http.StatusInternalServerError)
		}
		taxedSum = added
	}

	return taxedSum, nil
}

// GetValidShippingMethodsForOrder returns a list of valid shipping methods for given order
func (a *AppOrder) GetValidShippingMethodsForOrder(ord *order.Order) ([]*shipping.ShippingMethod, *model.AppError) {
	orderRequireShipping, appErr := a.OrderShippingIsRequired(ord.Id)
	if appErr != nil {
		appErr.Where = "GetValidShippingMethodsForOrder"
		return nil, appErr
	}

	if orderRequireShipping {
		return nil, nil
	}

	if ord.ShippingAddressID == nil {
		return nil, nil
	}

	orderSubTotal, appErr := a.OrderSubTotal(ord)
	if appErr != nil {
		appErr.Where = "GetValidShippingMethodsForOrder"
		return nil, appErr
	}

	shippingAddress, appErr := a.AccountApp().AddressById(*ord.ShippingAddressID)
	if appErr != nil {
		appErr.Where = "GetValidShippingMethodsForOrder"
		return nil, appErr
	}

	return a.ShippingApp().ApplicableShippingMethodsForOrder(ord, ord.ChannelID, orderSubTotal.Gross, shippingAddress.Country, nil)
}

// UpdateOrderDiscountForOrder Update the order_discount for an order and recalculate the order's prices
//
// `reason`, `valueType` and `value` can be nil
func (a *AppOrder) UpdateOrderDiscountForOrder(ord *order.Order, orderDiscountToUpdate *product_and_discount.OrderDiscount, reason string, valueType string, value *decimal.Decimal) *model.AppError {
	ord.PopulateNonDbFields() // NOTE: call this first

	if value == nil {
		value = orderDiscountToUpdate.Value
	}
	if valueType == "" {
		valueType = orderDiscountToUpdate.ValueType
	}

	if reason != "" {
		orderDiscountToUpdate.Reason = &reason
	}

	netTotal, err := a.ApplyDiscountToValue(value, valueType, orderDiscountToUpdate.Currency, ord.Total.Net)
	if err != nil {
		return model.NewAppError("UpdateOrderDiscountForOrder", app.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
	}
	grossTotal, err := a.ApplyDiscountToValue(value, valueType, orderDiscountToUpdate.Currency, ord.Total.Gross)
	if err != nil {
		return model.NewAppError("UpdateOrderDiscountForOrder", app.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	newAmount, _ := ord.Total.Sub(grossTotal)

	orderDiscountToUpdate.Amount = newAmount.Gross
	orderDiscountToUpdate.Value = value
	orderDiscountToUpdate.ValueType = valueType

	newOrderTotal, err := goprices.NewTaxedMoney(netTotal.(*goprices.Money), grossTotal.(*goprices.Money))
	if err != nil {
		return model.NewAppError("UpdateOrderDiscountForOrder", app.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
	}
	ord.Total = newOrderTotal

	_, appErr := a.DiscountApp().UpsertOrderDiscount(orderDiscountToUpdate)
	if appErr != nil {
		appErr.Where = "UpdateOrderDiscountForOrder"
		return appErr
	}
	return nil
}

// ApplyDiscountToValue Calculate the price based on the provided values
func (a *AppOrder) ApplyDiscountToValue(value *decimal.Decimal, valueType string, currency string, priceToDiscount interface{}) (interface{}, error) {
	// validate currency
	money, _ := goprices.NewMoney(value, currency)
	// MOTE: we can safely ignore the error here since OrderDiscounts's Currencies were validated before saving into database

	var discountCalculator discount.DiscountCalculator
	if valueType == product_and_discount.FIXED {
		discountCalculator = discount.Decorator(money)
	} else {
		discountCalculator = discount.Decorator(value)
	}

	return discountCalculator(priceToDiscount)
}

// GetProductsVoucherDiscountForOrder Calculate products discount value for a voucher, depending on its type.
func (a *AppOrder) GetProductsVoucherDiscountForOrder(ord *order.Order) (*goprices.Money, *model.AppError) {
	var (
		prices  []*goprices.Money
		voucher *product_and_discount.Voucher
	)

	if ord.VoucherID != nil {
		voucher, appErr := a.DiscountApp().VoucherById(*ord.VoucherID)
		if appErr != nil {
			if appErr.StatusCode == http.StatusInternalServerError {
				appErr.Where = "GetProductsVoucherDiscountForOrder"
				return nil, appErr
			}
			// ignore not found error
		} else {
			if voucher.Type == product_and_discount.SPECIFIC_PRODUCT {
				orderLinesOfOrder, appErr := a.OrderLinesByOption(&order.OrderLineFilterOption{
					OrderID: &model.StringFilter{
						StringOption: &model.StringOption{
							Eq: ord.Id,
						},
					},
				})
				if appErr != nil {
					appErr.Where = "GetProductsVoucherDiscountForOrder"
					return nil, appErr
				}

				discountedPrices, appErr := a.GetPricesOfDiscountedSpecificProduct(orderLinesOfOrder, voucher)
				if appErr != nil {
					appErr.Where = "GetProductsVoucherDiscountForOrder"
					return nil, appErr
				}

				prices = discountedPrices
			}
		}
	}

	if len(prices) == 0 {
		return nil, model.NewAppError("GetProductsVoucherDiscountForOrder", "app.order.offer_only_valid_for_selected_items.app_error", nil, "", http.StatusNotAcceptable)
	}

	return a.DiscountApp().GetProductsVoucherDiscount(voucher, prices, ord.ChannelID)
}

func (a *AppOrder) MatchOrdersWithNewUser(user *account.User) *model.AppError {
	ordersByOption, appErr := a.FilterOrdersByOptions(&order.OrderFilterOption{
		Status: &model.StringFilter{
			StringOption: &model.StringOption{
				NotIn: []string{
					order.DRAFT,
					order.UNCONFIRMED,
				},
			},
		},
		UserEmail: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: user.Email,
			},
		},
		UserID: &model.StringFilter{
			StringOption: &model.StringOption{
				NULL: model.NewBool(true),
			},
		},
	})
	if appErr != nil {
		appErr.Where = "MatchOrdersWithNewUser"
		return appErr
	}

	_, appErr = a.BulkUpsertOrders(ordersByOption)
	if appErr != nil {
		appErr.Where = "MatchOrdersWithNewUser"
		return appErr
	}
	return nil
}

// GetTotalOrderDiscount Return total order discount assigned to the order
func (a *AppOrder) GetTotalOrderDiscount(ord *order.Order) (*goprices.Money, *model.AppError) {
	orderDiscountsOfOrder, appErr := a.DiscountApp().OrderDiscountsByOption(&product_and_discount.OrderDiscountFilterOption{
		OrderID: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: ord.Id,
			},
		},
	})
	if appErr != nil {
		appErr.Where = "GetTotalOrderDiscount"
		return nil, appErr
	}

	totalOrderDiscount, _ := util.ZeroMoney(ord.Currency)
	for _, orderDiscount := range orderDiscountsOfOrder {
		orderDiscount.PopulateNonDbFields()
		addedMoney, err := totalOrderDiscount.Add(orderDiscount.Amount)
		if err != nil { // order's Currency != orderDiscount.Currency
			return nil, model.NewAppError("GetTotalOrderDiscount", app.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
		} else {
			totalOrderDiscount = addedMoney
		}
	}

	if less, err := totalOrderDiscount.LessThan(ord.UnDiscountedTotalGross); less && err == nil {
		return totalOrderDiscount, nil
	}

	return ord.UnDiscountedTotalGross, nil
}

// GetOrderDiscounts Return all discounts applied to the order by staff user
func (a *AppOrder) GetOrderDiscounts(ord *order.Order) ([]*product_and_discount.OrderDiscount, *model.AppError) {
	orderDiscounts, appErr := a.DiscountApp().OrderDiscountsByOption(&product_and_discount.OrderDiscountFilterOption{
		Type: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: product_and_discount.MANUAL,
			},
		},
		OrderID: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: ord.Id,
			},
		},
	})
	if appErr != nil {
		appErr.Where = "GetOrderDiscounts"
		return nil, appErr
	}

	return orderDiscounts, nil
}

// CreateOrderDiscountForOrder Add new order discount and update the prices
func (a *AppOrder) CreateOrderDiscountForOrder(ord *order.Order, reason string, valueType string, value *decimal.Decimal) (*product_and_discount.OrderDiscount, *model.AppError) {
	ord.PopulateNonDbFields()

	netTotal, err := a.ApplyDiscountToValue(value, valueType, ord.Currency, ord.Total.Net)
	if err != nil {
		return nil, model.NewAppError("CreateOrderDiscountForOrder", app.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
	}
	grossTotal, err := a.ApplyDiscountToValue(value, valueType, ord.Currency, ord.Total.Gross)
	if err != nil {
		return nil, model.NewAppError("CreateOrderDiscountForOrder", app.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	sub, _ := ord.Total.Sub(grossTotal.(*goprices.Money))
	newAmount := sub.Gross

	newOrderDiscount, appErr := a.DiscountApp().UpsertOrderDiscount(&product_and_discount.OrderDiscount{
		ValueType: valueType,
		Value:     value,
		Reason:    &reason,
		Amount:    newAmount,
	})
	if appErr != nil {
		appErr.Where = "CreateOrderDiscountForOrder"
		return nil, appErr
	}

	newOrderTotal, err := goprices.NewTaxedMoney(netTotal.(*goprices.Money), grossTotal.(*goprices.Money))
	if err != nil {
		return nil, model.NewAppError("CreateOrderDiscountForOrder", app.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	ord.Total = newOrderTotal
	_, appErr = a.UpsertOrder(ord)
	if appErr != nil {
		appErr.Where = "CreateOrderDiscountForOrder"
		return nil, appErr
	}

	return newOrderDiscount, nil
}

// RemoveOrderDiscountFromOrder Remove the order discount from order and update the prices.
func (a *AppOrder) RemoveOrderDiscountFromOrder(ord *order.Order, orderDiscount *product_and_discount.OrderDiscount) *model.AppError {
	appErr := a.DiscountApp().BulkDeleteOrderDiscounts([]string{orderDiscount.Id})
	if appErr != nil {
		appErr.Where = "RemoveOrderDiscountFromOrder"
		return appErr
	}

	ord.PopulateNonDbFields()
	orderDiscount.PopulateNonDbFields()

	newOrderTotal, err := ord.Total.Add(orderDiscount.Amount)
	if err != nil {
		return model.NewAppError("RemoveOrderDiscountFromOrder", app.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	ord.Total = newOrderTotal
	_, appErr = a.UpsertOrder(ord)
	if appErr != nil {
		appErr.Where = "RemoveOrderDiscountFromOrder"
		return appErr
	}

	return nil
}

// UpdateDiscountForOrderLine Update discount fields for order line. Apply discount to the price
//
// `reason`, `valueType` can be empty. `value` can be nil
func (a *AppOrder) UpdateDiscountForOrderLine(orderLine *order.OrderLine, ord *order.Order, reason string, valueType string, value *decimal.Decimal, manager interface{}, taxIncluded bool) *model.AppError {

	ord.PopulateNonDbFields()
	orderLine.PopulateNonDbFields()

	if reason != "" {
		orderLine.UnitDiscountReason = &reason
	}
	if value == nil {
		value = orderLine.UnitDiscountValue
	}
	if valueType == "" {
		valueType = orderLine.UnitDiscountType
	}

	if orderLine.UnitDiscountValue != value || orderLine.UnitDiscountType != valueType {
		unitPriceWithDiscount, err := a.ApplyDiscountToValue(value, valueType, orderLine.UnDiscountedUnitPrice.Currency, orderLine.UnDiscountedUnitPrice)
		if err != nil {
			return model.NewAppError("UpdateDiscountForOrderLine", app.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
		}

		newOrderLineUnitDiscount, err := orderLine.UnDiscountedUnitPrice.Sub(unitPriceWithDiscount)
		if err != nil {
			return model.NewAppError("UpdateDiscountForOrderLine", app.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
		}

		orderLine.UnitDiscount = newOrderLineUnitDiscount.Gross
		orderLine.UnitPrice = unitPriceWithDiscount.(*goprices.TaxedMoney)
		orderLine.UnitDiscountType = valueType
		orderLine.UnitDiscountValue = value
		orderLine.TotalPrice, _ = orderLine.UnitPrice.Mul(int(orderLine.Quantity))
		orderLine.UnDiscountedUnitPrice, _ = orderLine.UnitPrice.Sub(orderLine.UnitDiscount)
		orderLine.UnDiscountedTotalPrice, _ = orderLine.UnDiscountedUnitPrice.Mul(orderLine.Quantity)

	}

	// Save lines before calculating the taxes as some plugin can fetch all order data
	// from db
	_, appErr := a.UpsertOrderLine(orderLine)
	if appErr != nil {
		appErr.Where = "UpdateDiscountForOrderLine"
		return appErr
	}

	//-------------------------------------- TOTO: fixme
	panic("not implemented")
}

// RemoveDiscountFromOrderLine Drop discount applied to order line. Restore undiscounted price
func (a *AppOrder) RemoveDiscountFromOrderLine(orderLine *order.OrderLine, ord *order.Order, manager interface{}, taxIncluded bool) *model.AppError {
	orderLine.PopulateNonDbFields()

	orderLine.UnitPrice = orderLine.UnDiscountedUnitPrice
	orderLine.UnitDiscountAmount = &decimal.Zero
	orderLine.UnitDiscountValue = &decimal.Zero
	orderLine.UnitDiscountReason = model.NewString("")
	orderLine.TotalPrice, _ = orderLine.UnitPrice.Mul(int(orderLine.Quantity))

	_, appErr := a.UpsertOrderLine(orderLine)
	if appErr != nil {
		appErr.Where = "RemoveDiscountFromOrderLine"
		return appErr
	}

	//-----------------------TODO: fixme
	panic("not implemented")
}

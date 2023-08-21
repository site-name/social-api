package order

import (
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/samber/lo"
	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app/discount/types"
	"github.com/sitename/sitename/app/plugin/interfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/measurement"
	"github.com/sitename/sitename/modules/util"
	"gorm.io/gorm"
)

// GetOrderCountry Return country to which order will be shipped
func (a *ServiceOrder) GetOrderCountry(order *model.Order) (model.CountryCode, *model.AppError) {
	addressID := order.BillingAddressID
	orderRequireShipping, appErr := a.OrderShippingIsRequired(order.Id)
	if appErr != nil {
		return "", appErr
	}
	if orderRequireShipping {
		addressID = order.ShippingAddressID
	}

	if addressID == nil {
		return model.DEFAULT_COUNTRY, nil
	}

	address, appErr := a.srv.AccountService().AddressById(*addressID)
	if appErr != nil {
		return "", appErr
	}

	return address.Country, nil
}

// OrderLineNeedsAutomaticFulfillment Check if given line is digital and should be automatically fulfilled.
//
// NOTE: before calling this, caller can attach related data into `orderLine` so this function does not have to call the database
func (a *ServiceOrder) OrderLineNeedsAutomaticFulfillment(orderLine *model.OrderLine) (bool, *model.AppError) {
	if orderLine.VariantID == nil || orderLine.ProductVariant == nil {
		return false, nil
	}

	digitalContent := orderLine.ProductVariant.DigitalContent

	if digitalContent == nil {
		var appErr *model.AppError
		digitalContent, appErr = a.srv.ProductService().DigitalContentbyOption(&model.DigitalContentFilterOption{
			Conditions: squirrel.Eq{model.DigitalContentTableName + ".ProductVariantID": *orderLine.VariantID},
		})
		if appErr != nil {
			if appErr.StatusCode == http.StatusInternalServerError {
				return false, appErr
			}
			return false, nil
		}
	}

	if *digitalContent.UseDefaultSettings && *a.srv.Config().ShopSettings.AutomaticFulfillmentDigitalProducts {
		return true, nil
	}
	if *digitalContent.AutomaticFulfillment {
		return true, nil
	}

	return false, nil
}

// OrderNeedsAutomaticFulfillment checks if given order has digital products which shoul be automatically fulfilled.
func (a *ServiceOrder) OrderNeedsAutomaticFulfillment(order model.Order) (bool, *model.AppError) {
	digitalOrderLinesOfOrder, appErr := a.AllDigitalOrderLinesOfOrder(order.Id)
	if appErr != nil {
		return false, appErr
	}

	for _, orderLine := range digitalOrderLinesOfOrder {
		orderLineNeedsAutomaticFulfillment, appErr := a.OrderLineNeedsAutomaticFulfillment(orderLine)
		if appErr != nil {
			return false, appErr
		}
		if orderLineNeedsAutomaticFulfillment {
			return true, nil
		}
	}

	return false, nil
}

func (a *ServiceOrder) GetVoucherDiscountAssignedToOrder(order *model.Order) (*model.OrderDiscount, *model.AppError) {
	orderDiscountsOfOrder, appErr := a.srv.DiscountService().
		OrderDiscountsByOption(&model.OrderDiscountFilterOption{
			Conditions: squirrel.Eq{
				model.OrderDiscountTableName + ".Type":    model.VOUCHER,
				model.OrderDiscountTableName + ".OrderID": order.Id,
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
func (a *ServiceOrder) RecalculateOrderDiscounts(transaction *gorm.DB, order *model.Order) ([][2]*model.OrderDiscount, *model.AppError) {
	var changedOrderDiscounts [][2]*model.OrderDiscount

	orderDiscounts, appErr := a.srv.DiscountService().OrderDiscountsByOption(&model.OrderDiscountFilterOption{
		Conditions: squirrel.Eq{
			model.OrderDiscountTableName + ".OrderID": order.Id,
			model.OrderDiscountTableName + ".Type":    model.MANUAL,
		},
	})

	if appErr != nil {
		return nil, appErr
	}

	for _, orderDiscount := range orderDiscounts {

		previousOrderDiscount := orderDiscount.DeepCopy()
		currentTotal := order.Total.Gross.Amount

		appErr = a.UpdateOrderDiscountForOrder(transaction, order, orderDiscount, "", "", nil)
		if appErr != nil {
			return nil, appErr
		}

		discountValue := orderDiscount.Value
		amount := orderDiscount.Amount

		if (orderDiscount.ValueType == model.DISCOUNT_VALUE_TYPE_PERCENTAGE || currentTotal.LessThan(*discountValue)) &&
			!amount.Amount.Equal(previousOrderDiscount.Amount.Amount) {
			changedOrderDiscounts = append(changedOrderDiscounts, [2]*model.OrderDiscount{
				previousOrderDiscount,
				orderDiscount,
			})
		}
	}

	return changedOrderDiscounts, nil
}

// Recalculate and assign total price of order.
//
// Total price is a sum of items in order and order shipping price minus
// discount amount.
//
// Voucher discount amount is recalculated by default. To avoid this, pass
// update_voucher_discount argument set to False.
//
// NOTE: `kwargs` can be nil
func (a *ServiceOrder) RecalculateOrder(transaction *gorm.DB, order *model.Order, kwargs map[string]interface{}) *model.AppError {
	appErr := a.RecalculateOrderPrices(transaction, order, kwargs)
	if appErr != nil {
		return appErr
	}

	changedOrderDiscounts, appErr := a.RecalculateOrderDiscounts(transaction, order)
	if appErr != nil {
		return appErr
	}

	appErr = a.OrderDiscountsAutomaticallyUpdatedEvent(transaction, order, changedOrderDiscounts)
	if appErr != nil {
		return appErr
	}

	order, appErr = a.UpsertOrder(transaction, order)
	if appErr != nil {
		return appErr
	}

	return a.ReCalculateOrderWeight(transaction, order)
}

// ReCalculateOrderWeight
func (a *ServiceOrder) ReCalculateOrderWeight(transaction *gorm.DB, order *model.Order) *model.AppError {
	orderLines, appErr := a.OrderLinesByOption(&model.OrderLineFilterOption{
		Conditions: squirrel.Expr(model.OrderLineTableName+".OrderID = ? AND Orderlines.VariantID IS NOT NULL", order.Id),
	})
	if appErr != nil {
		return appErr
	}

	var (
		weight       = measurement.ZeroWeight
		atomicValue  atomic.Int32
		appErrorChan = make(chan *model.AppError)
		weightChan   = make(chan *measurement.Weight)
	)
	defer func() {
		close(appErrorChan)
		close(weightChan)
	}()
	atomicValue.Add(int32(len(orderLines)))

	for _, orderLine := range orderLines {
		go func(anOrderLine *model.OrderLine) {
			defer atomicValue.Add(-1)

			productVariantWeight, appErr := a.srv.ProductService().ProductVariantGetWeight(*anOrderLine.VariantID)
			if appErr != nil {
				appErrorChan <- appErr
				return
			}

			weightChan <- productVariantWeight.Mul(float32(anOrderLine.Quantity))
		}(orderLine)
	}

	for atomicValue.Load() != 0 {
		select {
		case appErr := <-appErrorChan:
			return appErr
		case wg := <-weightChan:
			addedWeight, err := weight.Add(wg)
			if err != nil {
				return model.NewAppError("ReCalculateOrderWeight", model.ErrorCalculatingMeasurementID, nil, err.Error(), http.StatusInternalServerError)
			}
			weight = addedWeight
		}
	}

	weight, _ = weight.ConvertTo(order.WeightUnit)
	order.WeightAmount = weight.Amount

	_, appErr = a.UpsertOrder(transaction, order)
	return appErr
}

func (a *ServiceOrder) UpdateTaxesForOrderLine(line model.OrderLine, ord model.Order, manager interfaces.PluginManagerInterface, taxIncluded bool) *model.AppError {
	variant := line.ProductVariant
	if variant == nil {
		var appErr *model.AppError
		variant, appErr = a.srv.ProductService().ProductVariantById(*line.ProductVariantID)
		if appErr != nil {
			return appErr
		}
	}

	product, appErr := a.srv.ProductService().ProductById(variant.ProductID)
	if appErr != nil {
		return appErr
	}

	line.PopulateNonDbFields() // this is needed

	linePrice := line.UnitPrice.Gross
	if !taxIncluded {
		linePrice = line.UnitPrice.Net
	}

	line.UnitPrice = &goprices.TaxedMoney{
		Net:      linePrice,
		Gross:    linePrice,
		Currency: line.Currency,
	}

	unitPrice, appErr := manager.CalculateOrderLineUnit(ord, line, *variant, *product)
	if appErr != nil {
		return appErr
	}

	totalPrice, appErr := manager.CalculateOrderlineTotal(ord, line, *variant, *product)
	if appErr != nil {
		return appErr
	}

	line.UnitPrice = unitPrice
	line.TotalPrice = totalPrice

	line.UnDiscountedUnitPrice, _ = line.UnitPrice.Add(line.UnitDiscount)
	line.UnDiscountedTotalPrice = totalPrice
	if line.UnitDiscount != nil && !line.UnitDiscount.Amount.Equal(decimal.Zero) {
		line.UnDiscountedTotalPrice = line.UnDiscountedUnitPrice.Mul(float64(line.Quantity))
	}

	unitPriceTax := unitPrice.Tax()
	if !unitPriceTax.Amount.Equal(decimal.Zero) && !unitPrice.Net.Amount.Equal(decimal.Zero) {
		line.TaxRate, appErr = manager.GetOrderLineTaxRate(ord, *product, *variant, nil, *unitPrice)
		if appErr != nil {
			return appErr
		}
	}

	return nil
}

func (a *ServiceOrder) UpdateTaxesForOrderLines(lines model.OrderLines, ord model.Order, manager interfaces.PluginManagerInterface, taxIncludeed bool) *model.AppError {
	for _, line := range lines.FilterNils() {
		appErr := a.UpdateTaxesForOrderLine(*line, ord, manager, taxIncludeed)
		if appErr != nil {
			return appErr
		}
	}

	_, appErr := a.BulkUpsertOrderLines(nil, lines)
	return appErr
}

// UpdateOrderPrices Update prices in order with given discounts and proper taxes.
func (a *ServiceOrder) UpdateOrderPrices(ord model.Order, manager interfaces.PluginManagerInterface, taxIncluded bool) *model.AppError {
	lines, appErr := a.OrderLinesByOption(&model.OrderLineFilterOption{
		Conditions: squirrel.Eq{model.OrderLineTableName + ".OrderID": ord.Id},
	})
	if appErr != nil {
		return appErr
	}

	appErr = a.UpdateTaxesForOrderLines(lines, ord, manager, taxIncluded)
	if appErr != nil {
		return appErr
	}

	if ord.ShippingMethodID != nil && model.IsValidId(*ord.ShippingMethodID) {
		shippingPrice, appErr := manager.CalculateOrderShipping(ord)
		if appErr != nil {
			return appErr
		}

		ord.ShippingPrice = shippingPrice
		ord.ShippingTaxRate, appErr = manager.GetOrderShippingTaxRate(ord, *shippingPrice)
		if appErr != nil {
			return appErr
		}

		_, appErr = a.UpsertOrder(nil, &ord)
		if appErr != nil {
			return appErr
		}
	}

	return a.RecalculateOrder(nil, &ord, nil)
}

func (s *ServiceOrder) GetValidCollectionPointsForOrder(lines model.OrderLines, addressCountryCode model.CountryCode) (model.Warehouses, *model.AppError) {
	// check shipping required:
	if !lo.SomeBy(lines, func(l *model.OrderLine) bool { return l.IsShippingRequired }) {
		return model.Warehouses{}, nil
	}
	if !addressCountryCode.IsValid() {
		return model.Warehouses{}, nil
	}

	warehouses, err := s.srv.Store.Warehouse().ApplicableForClickAndCollectOrderLines(lines, addressCountryCode)
	if err != nil {
		return nil, model.NewAppError("GetValidCollectionPointsForOrder", "app.order.valid_collection_points_for_order.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return warehouses, nil
}

// GetDiscountedLines returns a list of discounted order lines, filterd from given orderLines
func (a *ServiceOrder) GetDiscountedLines(orderLines model.OrderLines, voucher *model.Voucher) ([]*model.OrderLine, *model.AppError) {
	if len(orderLines) == 0 {
		return orderLines, nil
	}

	var (
		discountedProducts    model.Products
		discountedCategories  model.Categories
		discountedCollections model.Collections

		atomicValue atomic.Int32
		appErrChan  = make(chan *model.AppError)
	)
	defer func() {
		close(appErrChan)
	}()
	atomicValue.Add(3)

	go func() {
		defer atomicValue.Add(-1)

		products, appErr := a.srv.ProductService().ProductsByVoucherID(voucher.Id)
		if appErr != nil {
			appErrChan <- appErr
			return
		}

		discountedProducts = products
	}()

	go func() {
		defer atomicValue.Add(-1)

		categories, appErr := a.srv.ProductService().CategoriesByOption(&model.CategoryFilterOption{
			VoucherID: squirrel.Eq{model.VoucherCategoryTableName + ".VoucherID": voucher.Id},
		})
		if appErr != nil {
			appErrChan <- appErr
			return
		}

		discountedCategories = categories
	}()

	go func() {
		defer atomicValue.Add(-1)

		collections, appErr := a.srv.ProductService().CollectionsByVoucherID(voucher.Id)
		if appErr != nil {
			appErrChan <- appErr
			return
		}

		discountedCollections = collections
	}()

	for atomicValue.Load() != 0 {
		select {
		case appErr := <-appErrChan:
			return nil, appErr
		default:
		}
	}

	// try prefetching related product variants, products, collections related to given orderlines
	if orderLines[0].ProductVariant == nil {
		var appErr *model.AppError
		orderLines, appErr = a.srv.OrderService().OrderLinesByOption(&model.OrderLineFilterOption{
			Conditions: squirrel.Expr(model.OrderLineTableName+".Id IN ?", orderLines.IDs()),
			Preload:    []string{"ProductVariant.Product.Collections"}, // TODO: check if this works
		})
		if appErr != nil {
			return nil, appErr
		}
	}

	// filter duplicates
	discountedCategories = lo.UniqBy(discountedCategories, func(c *model.Category) string { return c.Id })
	discountedCollections = lo.UniqBy(discountedCollections, func(c *model.Collection) string { return c.Id })

	if len(discountedProducts) > 0 ||
		len(discountedCategories) > 0 ||
		len(discountedCollections) > 0 {

		var discountedOrderLines []*model.OrderLine

		for _, orderLine := range orderLines {
			if orderLine.ProductVariant != nil {
				var (
					lineProduct     = orderLine.ProductVariant.Product
					lineCollections model.Collections
					lineCategory    *model.Category
				)
				if lineProduct != nil {
					lineCollections = lineProduct.Collections

					if lineProduct.CategoryID != nil {
						categories, appErr := a.srv.ProductService().CategoryByIds([]string{*lineProduct.CategoryID}, true)
						if appErr != nil {
							return nil, appErr
						}
						lineCategory = categories[0]
					}
				}

				if (lineProduct != nil && discountedProducts.Contains(lineProduct)) ||
					(lineCategory != nil && discountedCategories.Contains(lineCategory)) ||
					lineCollections.IDs().InterSection(discountedCollections.IDs()).Len() > 0 {

					discountedOrderLines = append(discountedOrderLines, orderLine)
				}
			}
		}

		return discountedOrderLines, nil
	}

	// If there's no discounted products, collections or categories,
	// it means that all products are discounted
	return orderLines, nil
}

// Get prices of variants belonging to the discounted specific products.
//
// Specific products are products, collections and categories.
// Product must be assigned directly to the discounted category, assigning
// product to child category won't work
func (a *ServiceOrder) GetPricesOfDiscountedSpecificProduct(orderLines []*model.OrderLine, voucher *model.Voucher) ([]*goprices.Money, *model.AppError) {
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
func (a *ServiceOrder) GetVoucherDiscountForOrder(ord *model.Order) (result interface{}, notApplicableErr *model.NotApplicable, appErr *model.AppError) {

	ord.PopulateNonDbFields() // NOTE: must call this method before performing money, weight calculations

	// validate if order has voucher attached to
	if ord.VoucherID == nil {
		result = &goprices.Money{
			Amount:   decimal.Zero,
			Currency: ord.Currency,
		}
		return
	}

	notApplicableErr, appErr = a.srv.DiscountService().ValidateVoucherInOrder(ord)
	if appErr != nil || notApplicableErr != nil {
		return
	}

	orderLines, appErr := a.OrderLinesByOption(&model.OrderLineFilterOption{
		Conditions: squirrel.Eq{model.OrderLineTableName + ".OrderID": ord.Id},
	})
	if appErr != nil {
		return
	}

	orderSubTotal, appErr := a.srv.PaymentService().GetSubTotal(orderLines, ord.Currency)
	if appErr != nil {
		return
	}

	voucherOfDiscount, appErr := a.srv.DiscountService().VoucherById(*ord.VoucherID)
	if appErr != nil {
		return
	}

	if voucherOfDiscount.Type == model.VOUCHER_TYPE_ENTIRE_ORDER {
		result, appErr = a.srv.DiscountService().GetDiscountAmountFor(voucherOfDiscount, orderSubTotal.Gross, ord.ChannelID)
		return
	}
	if voucherOfDiscount.Type == model.VOUCHER_TYPE_SHIPPING {
		result, appErr = a.srv.DiscountService().GetDiscountAmountFor(voucherOfDiscount, ord.ShippingPrice, ord.ChannelID)
		return
	}
	// otherwise: Type is model.SPECIFIC_PRODUCT
	prices, appErr := a.GetPricesOfDiscountedSpecificProduct(orderLines, voucherOfDiscount)
	if appErr != nil {
		return
	}
	if len(prices) == 0 {
		appErr = model.NewAppError("GetVoucherDiscountForOrder", "app.order.offer_only_valid_for_selected_items.app_error", nil, "", http.StatusNotAcceptable)
		return
	}

	result, appErr = a.srv.DiscountService().GetProductsVoucherDiscount(voucherOfDiscount, prices, ord.ChannelID)
	return
}

func (a *ServiceOrder) calculateQuantityIncludingReturns(ord model.Order) (int, int, int, *model.AppError) {
	orderLinesOfOrder, appErr := a.OrderLinesByOption(&model.OrderLineFilterOption{
		Conditions: squirrel.Eq{model.OrderLineTableName + ".OrderID": ord.Id},
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

	fulfillmentsOfOrder, appErr := a.FulfillmentsByOption(&model.FulfillmentFilterOption{
		Conditions: squirrel.Eq{model.FulfillmentTableName + ".OrderID": ord.Id},
	})
	if appErr != nil {
		return 0, 0, 0, appErr
	}

	// filter all fulfillments that has `Status` is either: "returned", "refunded_and_returned" and "replaced"
	var (
		filteredFulfillmentIDs []string
		fulfillmentMap         = map[string]*model.Fulfillment{}
	)
	for _, fulfillment := range fulfillmentsOfOrder {
		if fulfillment.Status == model.FULFILLMENT_RETURNED ||
			fulfillment.Status == model.FULFILLMENT_REFUNDED_AND_RETURNED ||
			fulfillment.Status == model.FULFILLMENT_REPLACED {

			filteredFulfillmentIDs = append(filteredFulfillmentIDs, fulfillment.Id)
			fulfillmentMap[fulfillment.Id] = fulfillment
		}
	}

	// finds all fulfillment lines belong to filtered fulfillments
	fulfillmentLines, appErr := a.FulfillmentLinesByOption(&model.FulfillmentLineFilterOption{
		Conditions: squirrel.Eq{model.FulfillmentLineTableName + ".FulfillmentID": filteredFulfillmentIDs},
	})
	if appErr != nil {
		return 0, 0, 0, appErr
	}

	for _, fulfillmentLine := range fulfillmentLines {
		parentFulfillmentStatus := fulfillmentMap[fulfillmentLine.FulfillmentID].Status

		if parentFulfillmentStatus == model.FULFILLMENT_RETURNED || parentFulfillmentStatus == model.FULFILLMENT_REFUNDED_AND_RETURNED {
			quantityReturned += fulfillmentLine.Quantity
		} else if parentFulfillmentStatus == model.FULFILLMENT_REPLACED {
			quantityReplaced += fulfillmentLine.Quantity
		}
	}

	totalOrderLinesQuantity -= quantityReplaced
	quantityFulfilled -= quantityReplaced

	return totalOrderLinesQuantity, quantityFulfilled, quantityReturned, nil
}

// UpdateOrderStatus Update order status depending on fulfillments
func (a *ServiceOrder) UpdateOrderStatus(transaction *gorm.DB, ord model.Order) *model.AppError {

	totalQuantity, quantityFulfilled, quantityReturned, appErr := a.calculateQuantityIncludingReturns(ord)
	if appErr != nil {
		return appErr
	}

	var status model.OrderStatus
	if totalQuantity == 0 {
		status = ord.Status
	} else if quantityFulfilled <= 0 {
		status = model.ORDER_STATUS_UNFULFILLED
	} else if quantityReturned > 0 && quantityReturned < totalQuantity {
		status = model.ORDER_STATUS_PARTIALLY_RETURNED
	} else if quantityReturned == totalQuantity {
		status = model.ORDER_STATUS_RETURNED
	} else if quantityFulfilled < totalQuantity {
		status = model.ORDER_STATUS_PARTIALLY_FULFILLED
	} else {
		status = model.ORDER_STATUS_FULFILLED
	}

	if status != ord.Status {
		ord.Status = status
		_, appErr := a.UpsertOrder(transaction, &ord)
		if appErr != nil {
			return appErr
		}
	}

	return nil
}

// AddVariantToOrder Add total_quantity of variant to order.
//
// Returns an order line the variant was added to.
func (s *ServiceOrder) AddVariantToOrder(order model.Order, variant model.ProductVariant, quantity int, user *model.User, _ interface{}, manager interfaces.PluginManagerInterface, discounts []*model.DiscountInfo, allocateStock bool) (*model.OrderLine, *model.InsufficientStock, *model.AppError) {
	transaction := s.srv.Store.GetMaster().Begin()
	if transaction.Error != nil {
		return nil, nil, model.NewAppError("AddVariantToOrder", model.ErrorCreatingTransactionErrorID, nil, transaction.Error.Error(), http.StatusInternalServerError)
	}
	defer transaction.Rollback()

	chanNel, appErr := s.srv.ChannelService().ChannelByOption(&model.ChannelFilterOption{
		Conditions: squirrel.Eq{model.ChannelTableName + ".Id": order.ChannelID},
	})
	if appErr != nil {
		return nil, nil, appErr
	}

	orderLinesOfOrder, appErr := s.OrderLinesByOption(&model.OrderLineFilterOption{
		Conditions: squirrel.Eq{
			model.OrderLineTableName + ".OrderID":   order.Id,
			model.OrderLineTableName + ".VariantID": variant.Id,
		},
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, nil, appErr
		}
	}

	// order line
	var orderLine *model.OrderLine

	if len(orderLinesOfOrder) > 0 {
		orderLine = orderLinesOfOrder[0]
		oldQuantity := orderLine.Quantity
		newQuantity := oldQuantity + quantity

		lineInfo := &model.OrderLineData{
			Line:     *orderLine,
			Quantity: oldQuantity,
		}
		insufficientStock, appErr := s.ChangeOrderLineQuantity(transaction, user.Id, nil, lineInfo, oldQuantity, newQuantity, chanNel.Slug, manager, false)
		if insufficientStock != nil || appErr != nil {
			return nil, insufficientStock, appErr
		}
	} else {
		// in case no order line found
		product, appErr := s.srv.ProductService().ProductById(variant.ProductID)
		if appErr != nil {
			return nil, nil, appErr
		}

		collections, appErr := s.srv.ProductService().CollectionsByProductID(product.Id)
		if appErr != nil {
			return nil, nil, appErr
		}

		variantChannelListings, appErr := s.srv.ProductService().ProductVariantChannelListingsByOption(&model.ProductVariantChannelListingFilterOption{
			Conditions: squirrel.Eq{
				model.ProductVariantChannelListingTableName + ".VariantID": variant.Id,
				model.ProductVariantChannelListingTableName + ".ChannelID": chanNel.Id,
			},
		})
		if appErr != nil {
			return nil, nil, appErr // NOTE: does not care what type of error, just return
		}

		price, appErr := s.srv.ProductService().ProductVariantGetPrice(&variant, *product, collections, *chanNel, variantChannelListings[0], discounts)
		if appErr != nil {
			return nil, nil, appErr
		}

		taxedUnitPrice := &goprices.TaxedMoney{
			Net:      price,
			Gross:    price,
			Currency: price.Currency,
		}

		totalPrice := taxedUnitPrice.Mul(float64(quantity))
		productName := product.String()
		variantName := variant.String()

		var translatedProductName string
		productTranslations, appErr := s.srv.ProductService().
			ProductTranslationsByOption(&model.ProductTranslationFilterOption{
				Conditions: squirrel.Eq{
					model.ProductTranslationTableName + ".LanguageCode": user.Locale,
					model.ProductTranslationTableName + ".ProductID":    product.Id,
				},
			})
		if appErr != nil {
			if appErr.StatusCode == http.StatusInternalServerError {
				return nil, nil, appErr
			}
		} else {
			translatedProductName = productTranslations[0].Name
		}

		var translatedVariantName string
		variantTranslations, appErr := s.srv.ProductService().ProductVariantTranslationsByOption(&model.ProductVariantTranslationFilterOption{
			Conditions: squirrel.Eq{
				model.ProductVariantTranslationTableName + ".LanguageCode":     user.Locale,
				model.ProductVariantTranslationTableName + ".ProductVariantID": variant.Id,
			},
		})
		if appErr != nil {
			if appErr.StatusCode == http.StatusInternalServerError {
				return nil, nil, appErr
			}
		} else {
			translatedVariantName = variantTranslations[0].Name
		}

		if translatedProductName == productName {
			translatedProductName = ""
		}
		if translatedVariantName == variantName {
			translatedVariantName = ""
		}

		variantRequiresShipping, appErr := s.srv.ProductService().ProductsRequireShipping([]string{variant.ProductID})
		if appErr != nil {
			return nil, nil, appErr
		}
		productType, appErr := s.srv.ProductService().ProductTypeByOption(&model.ProductTypeFilterOption{
			Conditions: squirrel.Eq{model.ProductTypeTableName + ".Id": product.ProductTypeID},
		})
		if appErr != nil {
			return nil, nil, appErr
		}

		orderLine, appErr = s.UpsertOrderLine(transaction, &model.OrderLine{
			ProductName:           productName,
			VariantName:           variantName,
			TranslatedProductName: translatedProductName,
			TranslatedVariantName: translatedVariantName,
			ProductSku:            &variant.Sku,
			IsShippingRequired:    variantRequiresShipping,
			IsGiftcard:            productType.IsGiftcard(),
			Quantity:              quantity,
			UnitPrice:             taxedUnitPrice,
			TotalPrice:            totalPrice,
			VariantID:             &variant.Id,
		})
		if appErr != nil {
			return nil, nil, appErr
		}

		unitPrice, appErr := manager.CalculateOrderLineUnit(order, *orderLine, variant, *product)
		if appErr != nil {
			return nil, nil, appErr
		}

		totalPrice, appErr = manager.CalculateOrderlineTotal(order, *orderLine, variant, *product)
		if appErr != nil {
			return nil, nil, appErr
		}

		orderLine.UnitPrice = unitPrice
		orderLine.TotalPrice = totalPrice
		orderLine.UnDiscountedUnitPrice = unitPrice
		orderLine.UnDiscountedTotalPrice = totalPrice
		orderLine.TaxRate, appErr = manager.GetOrderLineTaxRate(order, *product, variant, nil, *unitPrice)
		if appErr != nil {
			return nil, nil, appErr
		}

		_, appErr = s.UpsertOrderLine(transaction, orderLine)
		if appErr != nil {
			return nil, nil, appErr
		}
	}

	if allocateStock {
		insufficientStockErr, appErr := s.srv.WarehouseService().IncreaseAllocations(
			[]*model.OrderLineData{
				{
					Line:        *orderLine,
					Quantity:    quantity,
					Variant:     &variant,
					WarehouseID: nil,
				},
			},
			chanNel.Slug,
			manager,
		)
		if insufficientStockErr != nil || appErr != nil {
			return nil, insufficientStockErr, appErr
		}
	}

	// commit transaction
	if err := transaction.Commit().Error; err != nil {
		return nil, nil, model.NewAppError("AddVariantToOrder", model.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return orderLine, nil, nil
}

// AddGiftcardsToOrder
func (s *ServiceOrder) AddGiftcardsToOrder(transaction *gorm.DB, checkoutInfo model.CheckoutInfo, order *model.Order, totalPriceLeft *goprices.Money, user *model.User, _ interface{}) *model.AppError {
	var (
		balanceData       = model.BalanceData{}
		usedByUser        = checkoutInfo.User
		usedByEmail       = checkoutInfo.GetCustomerEmail()
		giftcardsToUpdate = []*model.GiftCard{}

		giftcardsToAddToOrder []*model.GiftCard
	)

	giftcards, appErr := s.srv.GiftcardService().GiftcardsByOption(&model.GiftCardFilterOption{
		SelectForUpdate: true,
		CheckoutToken:   squirrel.Eq{model.GiftcardCheckoutTableName + ".CheckoutID": checkoutInfo.Checkout.Token},
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return appErr
		}
	}

	for _, giftCard := range giftcards {
		if totalPriceLeft.Amount.GreaterThan(decimal.Zero) {
			giftcardsToAddToOrder = append(giftcardsToAddToOrder, giftCard)

			balanceData = append(balanceData, s.UpdateGiftcardBalance(giftCard, totalPriceLeft))

			if giftCard.UsedByEmail == nil {
				if usedByUser != nil {
					giftCard.UsedByID = &usedByUser.Id
				}
				giftCard.UsedByEmail = &usedByEmail
			}

			giftCard.LastUsedOn = model.NewPrimitive(model.GetMillis())
			giftcardsToUpdate = append(giftcardsToUpdate, giftCard)
		}
	}

	appErr = s.srv.GiftcardService().AddGiftcardRelations(transaction, giftcardsToAddToOrder, model.Orders{order})
	if appErr != nil {
		return appErr
	}

	if len(giftcardsToUpdate) > 0 {
		_, appErr = s.srv.GiftcardService().UpsertGiftcards(transaction, giftcardsToUpdate...)
		if appErr != nil {
			return appErr
		}
	}

	_, appErr = s.srv.GiftcardService().GiftcardsUsedInOrderEvent(transaction, balanceData, order.Id, user, nil)
	return appErr
}

func (s *ServiceOrder) UpdateGiftcardBalance(giftCard *model.GiftCard, totalPriceLeft *goprices.Money) model.BalanceObject {
	giftCard.PopulateNonDbFields() // NOTE: this call is important

	previousBalance := giftCard.CurrentBalance
	if totalPriceLeft.LessThan(giftCard.CurrentBalance) {
		giftCard.CurrentBalance, _ = giftCard.CurrentBalance.Sub(totalPriceLeft)
		totalPriceLeft, _ = util.ZeroMoney(totalPriceLeft.Currency)
	} else {
		totalPriceLeft, _ = totalPriceLeft.Sub(giftCard.CurrentBalance)
		giftCard.CurrentBalanceAmount = &decimal.Zero
	}

	return model.BalanceObject{
		Giftcard:        *giftCard,
		PreviousBalance: &previousBalance.Amount,
	}
}

func (a *ServiceOrder) updateAllocationsForLine(lineInfo *model.OrderLineData, oldQuantity int, newQuantity int, channelSlug string, manager interfaces.PluginManagerInterface) (*model.InsufficientStock, *model.AppError) {
	if oldQuantity == newQuantity {
		return nil, nil
	}

	orderLinesWithTrackInventory := a.srv.WarehouseService().GetOrderLinesWithTrackInventory([]*model.OrderLineData{lineInfo})
	if len(orderLinesWithTrackInventory) == 0 {
		return nil, nil
	}

	if oldQuantity < newQuantity {
		lineInfo.Quantity = newQuantity - oldQuantity
		return a.srv.WarehouseService().IncreaseAllocations([]*model.OrderLineData{lineInfo}, channelSlug, manager)
	} else {
		lineInfo.Quantity = oldQuantity - newQuantity
		return a.srv.WarehouseService().DecreaseAllocations([]*model.OrderLineData{lineInfo}, manager)
	}
}

// ChangeOrderLineQuantity Change the quantity of ordered items in a order line.
//
// NOTE: userID can be empty
func (a *ServiceOrder) ChangeOrderLineQuantity(transaction *gorm.DB, userID string, _ interface{}, lineInfo *model.OrderLineData, oldQuantity int, newQuantity int, channelSlug string, manager interfaces.PluginManagerInterface, sendEvent bool) (*model.InsufficientStock, *model.AppError) {
	orderLine := lineInfo.Line
	// NOTE: this must be called
	orderLine.PopulateNonDbFields()

	if newQuantity > 0 {
		order, appErr := a.OrderById(lineInfo.Line.OrderID)
		if appErr != nil {
			return nil, appErr
		}

		if order.IsUnconfirmed() {
			insufficientStock, appErr := a.updateAllocationsForLine(lineInfo, oldQuantity, newQuantity, channelSlug, manager)
			if appErr != nil || insufficientStock != nil {
				return insufficientStock, appErr
			}
		}

		lineInfo.Line.Quantity = newQuantity

		totalPriceNetAmount := orderLine.UnitPriceNetAmount.Mul(decimal.NewFromInt32(int32(orderLine.Quantity)))
		totalPriceGrossAmount := orderLine.UnitPriceGrossAmount.Mul(decimal.NewFromInt32(int32(orderLine.Quantity)))
		orderLine.TotalPriceNetAmount = model.NewPrimitive(totalPriceNetAmount.Round(3))
		orderLine.TotalPriceGrossAmount = model.NewPrimitive(totalPriceGrossAmount.Round(3))

		unDiscountedTotalPriceNetAmount := orderLine.UnDiscountedUnitPriceNetAmount.Mul(decimal.NewFromInt32(int32(orderLine.Quantity)))
		unDiscountedTotalpriceGrossAmount := orderLine.UnDiscountedUnitPriceGrossAmount.Mul(decimal.NewFromInt32(int32(orderLine.Quantity)))
		orderLine.UnDiscountedTotalPriceNetAmount = model.NewPrimitive(unDiscountedTotalPriceNetAmount.Round(3))
		orderLine.UnDiscountedTotalPriceGrossAmount = model.NewPrimitive(unDiscountedTotalpriceGrossAmount.Round(3))

		_, appErr = a.UpsertOrderLine(nil, &orderLine)
		if appErr != nil {
			return nil, appErr
		}
	} else {
		insufficientErr, appErr := a.DeleteOrderLine(lineInfo, manager)
		if appErr != nil || insufficientErr != nil {
			return insufficientErr, appErr
		}
	}

	quantityDiff := int(oldQuantity) - int(newQuantity)

	if sendEvent {
		appErr := a.CreateOrderEvent(transaction, &orderLine, userID, quantityDiff)
		if appErr != nil {
			return nil, appErr
		}
	}

	return nil, nil
}

func (a *ServiceOrder) CreateOrderEvent(transaction *gorm.DB, orderLine *model.OrderLine, userID string, quantityDiff int) *model.AppError {
	var appErr *model.AppError

	var savingUserID *string
	if userID != "" {
		savingUserID = &userID
	}

	if quantityDiff > 0 {
		_, appErr = a.CommonCreateOrderEvent(transaction, &model.OrderEventOption{
			OrderID: orderLine.OrderID,
			UserID:  savingUserID,
			Type:    model.ORDER_EVENT_TYPE_REMOVED_PRODUCTS,
			Parameters: model.StringInterface{
				"lines": a.LinesPerQuantityToLineObjectList([]*model.QuantityOrderLine{
					{
						Quantity:  quantityDiff,
						OrderLine: orderLine,
					},
				}),
			},
		})
	} else if quantityDiff < 0 {
		_, appErr = a.CommonCreateOrderEvent(transaction, &model.OrderEventOption{
			OrderID: orderLine.OrderID,
			UserID:  savingUserID,
			Type:    model.ORDER_EVENT_TYPE_ADDED_PRODUCTS,
			Parameters: model.StringInterface{
				"lines": a.LinesPerQuantityToLineObjectList([]*model.QuantityOrderLine{
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

// DeleteOrderLine Delete an order line from an order.
func (a *ServiceOrder) DeleteOrderLine(lineInfo *model.OrderLineData, manager interfaces.PluginManagerInterface) (*model.InsufficientStock, *model.AppError) {
	ord, appErr := a.OrderById(lineInfo.Line.OrderID)
	if appErr != nil {
		return nil, appErr
	}

	if ord.IsUnconfirmed() {
		insufficientErr, appErr := a.srv.WarehouseService().DecreaseAllocations([]*model.OrderLineData{lineInfo}, manager)
		if appErr != nil || insufficientErr != nil {
			return insufficientErr, appErr
		}
	}

	return nil, a.DeleteOrderLines([]string{lineInfo.Line.Id})
}

// RestockOrderLines Return ordered products to corresponding stocks
func (a *ServiceOrder) RestockOrderLines(ord *model.Order, manager interfaces.PluginManagerInterface) *model.AppError {
	countryCode, appError := a.GetOrderCountry(ord)
	if appError != nil {
		return appError
	}

	warehouses, appError := a.srv.WarehouseService().WarehousesByOption(&model.WarehouseFilterOption{
		ShippingZonesCountries: squirrel.Like{model.ShippingZoneTableName + ".Countries": "%" + countryCode + "%"},
	})
	if appError != nil {
		return appError
	}
	defaultWarehouse := warehouses[0]

	orderLinesOfOrder, appError := a.OrderLinesByOption(&model.OrderLineFilterOption{
		Conditions: squirrel.Eq{model.OrderLineTableName + ".OrderID": ord.Id},
	})
	if appError != nil {
		return appError
	}

	var (
		dellocatingStockLines []*model.OrderLineData
		mut                   sync.Mutex
		wg                    sync.WaitGroup
	)

	setAppError := func(err *model.AppError) {
		if err != nil {
			mut.Lock()
			if appError == nil && err != nil {
				appError = err
			}
			mut.Unlock()
		}
	}

	for _, orderLine := range orderLinesOfOrder {
		if orderLine.VariantID != nil {

			wg.Add(1)

			go func(anOrderLine *model.OrderLine) {
				productVariant, appErr := a.srv.ProductService().ProductVariantById(*anOrderLine.VariantID)
				if appErr != nil {
					setAppError(appErr)
					return
				}

				if *productVariant.TrackInventory {
					if anOrderLine.QuantityUnFulfilled() > 0 {

						mut.Lock()
						dellocatingStockLines = append(dellocatingStockLines, &model.OrderLineData{
							Line:     *anOrderLine,
							Quantity: anOrderLine.QuantityUnFulfilled(),
						})
						mut.Unlock()

					}

					if anOrderLine.QuantityFulfilled > 0 {
						allocations, appErr := a.srv.WarehouseService().AllocationsByOption(&model.AllocationFilterOption{
							Conditions: squirrel.Eq{model.AllocationTableName + ".OrderLineID": anOrderLine.Id},
						})
						if appErr != nil {
							setAppError(appErr)
						} else {
							warehouse := defaultWarehouse
							if len(allocations) > 0 {
								warehouseOfOrderLine, appErr := a.srv.WarehouseService().WarehouseByStockID(allocations[0].StockID)
								if appErr != nil {
									setAppError(appErr)
								} else {
									warehouse = warehouseOfOrderLine
								}
							}

							appErr = a.srv.WarehouseService().IncreaseStock(anOrderLine, warehouse, anOrderLine.QuantityFulfilled, false)
							setAppError(appErr)
						}
					}
				}

				if anOrderLine.QuantityFulfilled > 0 {
					anOrderLine.QuantityFulfilled = 0

					_, appErr = a.UpsertOrderLine(nil, anOrderLine)
					setAppError(appErr)
				}

				wg.Done()
			}(orderLine)
		}
	}

	wg.Wait()

	if len(dellocatingStockLines) > 0 {
		_, appError = a.srv.WarehouseService().DeallocateStock(dellocatingStockLines, manager)
	}

	return appError
}

// RestockFulfillmentLines Return fulfilled products to corresponding stocks.
//
// Return products to stocks and update order lines quantity fulfilled values.
func (a *ServiceOrder) RestockFulfillmentLines(transaction *gorm.DB, fulfillment *model.Fulfillment, warehouse *model.WareHouse) (appErr *model.AppError) {
	fulfillmentLines, appErr := a.FulfillmentLinesByOption(&model.FulfillmentLineFilterOption{
		Conditions: squirrel.Eq{model.FulfillmentLineTableName + ".FulfillmentID": fulfillment.Id},
	})
	if appErr != nil {
		return appErr
	}

	orderLinesOfFulfillmentLines, appErr := a.OrderLinesByOption(&model.OrderLineFilterOption{
		Conditions: squirrel.Eq{model.OrderLineTableName + ".Id": fulfillmentLines.OrderLineIDs()},
	})
	if appErr != nil {
		return appErr
	}

	// create map with key: fulfillmentLine.Id, value: *OrderLine
	mapFulfillmentLine_OrderLine := map[string]*model.OrderLine{}
	for _, fulfillmentLine := range fulfillmentLines {
		for _, orderLine := range orderLinesOfFulfillmentLines {
			if fulfillmentLine.OrderLineID == orderLine.Id {
				mapFulfillmentLine_OrderLine[fulfillmentLine.Id] = orderLine
			}
		}
	}

	productVariantsOfOrderLines, appErr := a.srv.ProductService().ProductVariantsByOption(&model.ProductVariantFilterOption{
		Conditions: squirrel.Eq{model.ProductVariantTableName + ".Id": orderLinesOfFulfillmentLines.ProductVariantIDs()},
	})
	if appErr != nil {
		return appErr
	}

	// create map with key: orderLine.Id, value: *ProductVariant
	mapOrderLine_productVariant := map[string]*model.ProductVariant{}
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
			appErr := a.srv.WarehouseService().IncreaseStock(orderLineOfFulfillment, warehouse, fulfillmentLine.Quantity, true)
			if appErr != nil {
				return appErr
			}
		}

		orderLineOfFulfillment.QuantityFulfilled -= fulfillmentLine.Quantity
	}

	_, appErr = a.BulkUpsertOrderLines(transaction, orderLinesOfFulfillmentLines)
	return appErr
}

func (a *ServiceOrder) SumOrderTotals(orders []*model.Order, currencyCode string) (*goprices.TaxedMoney, *model.AppError) {
	taxedSum, _ := util.ZeroTaxedMoney(currencyCode)
	if len(orders) == 0 {
		return taxedSum, nil
	}
	// validate given `currencyCode` is valid
	currencyCode = strings.ToUpper(currencyCode)
	if _, err := goprices.GetCurrencyPrecision(currencyCode); err != nil || currencyCode != orders[0].Currency {
		return nil, model.NewAppError("SumOrderTotals", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "currencyCode"}, err.Error(), http.StatusBadRequest)
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
func (a *ServiceOrder) GetValidShippingMethodsForOrder(ord *model.Order) ([]*model.ShippingMethod, *model.AppError) {
	orderRequireShipping, appErr := a.OrderShippingIsRequired(ord.Id)
	if appErr != nil {
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
		return nil, appErr
	}

	shippingAddress, appErr := a.srv.AccountService().AddressById(*ord.ShippingAddressID)
	if appErr != nil {
		return nil, appErr
	}

	return a.srv.ShippingService().ApplicableShippingMethodsForOrder(ord, ord.ChannelID, orderSubTotal.Gross, shippingAddress.Country, nil)
}

// UpdateOrderDiscountForOrder Update the order_discount for an order and recalculate the order's prices
//
// `reason`, `valueType` and `value` can be nil
func (a *ServiceOrder) UpdateOrderDiscountForOrder(transaction *gorm.DB, ord *model.Order, orderDiscountToUpdate *model.OrderDiscount, reason string, valueType model.DiscountValueType, value *decimal.Decimal) *model.AppError {
	ord.PopulateNonDbFields() // NOTE: call this first

	if value == nil {
		value = orderDiscountToUpdate.Value
	}
	if !valueType.IsValid() {
		valueType = orderDiscountToUpdate.ValueType
	}

	if reason != "" {
		orderDiscountToUpdate.Reason = &reason
	}

	netTotal, err := a.ApplyDiscountToValue(value, valueType, orderDiscountToUpdate.Currency, ord.Total.Net)
	if err != nil {
		return model.NewAppError("UpdateOrderDiscountForOrder", model.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
	}
	grossTotal, err := a.ApplyDiscountToValue(value, valueType, orderDiscountToUpdate.Currency, ord.Total.Gross)
	if err != nil {
		return model.NewAppError("UpdateOrderDiscountForOrder", model.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	newAmount, _ := ord.Total.Sub(grossTotal)

	orderDiscountToUpdate.Amount = newAmount.Gross
	orderDiscountToUpdate.Value = value
	orderDiscountToUpdate.ValueType = valueType

	newOrderTotal, err := goprices.NewTaxedMoney(netTotal.(*goprices.Money), grossTotal.(*goprices.Money))
	if err != nil {
		return model.NewAppError("UpdateOrderDiscountForOrder", model.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
	}
	ord.Total = newOrderTotal

	_, appErr := a.srv.DiscountService().UpsertOrderDiscount(transaction, orderDiscountToUpdate)
	if appErr != nil {
		return appErr
	}
	return nil
}

// ApplyDiscountToValue Calculate the price based on the provided values
func (a *ServiceOrder) ApplyDiscountToValue(value *decimal.Decimal, valueType model.DiscountValueType, currency string, priceToDiscount interface{}) (interface{}, error) {
	// validate currency
	money := &goprices.Money{
		Amount:   *value,
		Currency: currency,
	}
	// MOTE: we can safely ignore the error here since OrderDiscounts's Currencies were validated before saving into database

	var discountCalculator types.DiscountCalculator
	if valueType == model.DISCOUNT_VALUE_TYPE_FIXED {
		discountCalculator = a.srv.DiscountService().Decorator(money)
	} else {
		discountCalculator = a.srv.DiscountService().Decorator(value)
	}

	return discountCalculator(priceToDiscount, nil)
}

// GetProductsVoucherDiscountForOrder Calculate products discount value for a voucher, depending on its type.
func (a *ServiceOrder) GetProductsVoucherDiscountForOrder(ord *model.Order) (*goprices.Money, *model.AppError) {
	var (
		prices  []*goprices.Money
		voucher *model.Voucher
	)

	if ord.VoucherID != nil {
		voucher, appErr := a.srv.DiscountService().VoucherById(*ord.VoucherID)
		if appErr != nil {
			if appErr.StatusCode == http.StatusInternalServerError {
				return nil, appErr
			}
			// ignore not found error
		} else {
			if voucher.Type == model.VOUCHER_TYPE_SPECIFIC_PRODUCT {
				orderLinesOfOrder, appErr := a.OrderLinesByOption(&model.OrderLineFilterOption{
					Conditions: squirrel.Eq{model.OrderLineTableName + ".OrderID": ord.Id},
				})
				if appErr != nil {
					return nil, appErr
				}

				discountedPrices, appErr := a.GetPricesOfDiscountedSpecificProduct(orderLinesOfOrder, voucher)
				if appErr != nil {
					return nil, appErr
				}

				prices = discountedPrices
			}
		}
	}

	if len(prices) == 0 {
		return nil, model.NewAppError("GetProductsVoucherDiscountForOrder", "app.order.offer_only_valid_for_selected_items.app_error", nil, "", http.StatusNotAcceptable)
	}

	return a.srv.DiscountService().GetProductsVoucherDiscount(voucher, prices, ord.ChannelID)
}

func (a *ServiceOrder) MatchOrdersWithNewUser(user *model.User) *model.AppError {
	ordersByOption, appErr := a.FilterOrdersByOptions(&model.OrderFilterOption{
		Conditions: squirrel.And{
			squirrel.Eq{model.OrderTableName + ".UserID": nil},
			squirrel.Eq{model.OrderTableName + ".UserEmail": user.Email},
			squirrel.NotEq{model.OrderTableName + ".Status": []model.OrderStatus{(model.ORDER_STATUS_DRAFT), (model.ORDER_STATUS_UNCONFIRMED)}},
		},
	})
	if appErr != nil {
		return appErr
	}

	_, appErr = a.BulkUpsertOrders(nil, ordersByOption)
	if appErr != nil {
		return appErr
	}
	return nil
}

// GetTotalOrderDiscount Return total order discount assigned to the order
func (a *ServiceOrder) GetTotalOrderDiscount(ord *model.Order) (*goprices.Money, *model.AppError) {
	orderDiscountsOfOrder, appErr := a.srv.DiscountService().OrderDiscountsByOption(&model.OrderDiscountFilterOption{
		Conditions: squirrel.Eq{model.OrderDiscountTableName + ".OrderID": ord.Id},
	})
	if appErr != nil {
		return nil, appErr
	}

	totalOrderDiscount, _ := util.ZeroMoney(ord.Currency)
	for _, orderDiscount := range orderDiscountsOfOrder {
		orderDiscount.PopulateNonDbFields()
		addedMoney, err := totalOrderDiscount.Add(orderDiscount.Amount)
		if err != nil { // order's Currency != orderDiscount.Currency
			return nil, model.NewAppError("GetTotalOrderDiscount", model.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
		} else {
			totalOrderDiscount = addedMoney
		}
	}

	if totalOrderDiscount.LessThan(ord.UnDiscountedTotalGross) {
		return totalOrderDiscount, nil
	}

	return ord.UnDiscountedTotalGross, nil
}

// GetOrderDiscounts Return all discounts applied to the order by staff user
func (a *ServiceOrder) GetOrderDiscounts(ord *model.Order) ([]*model.OrderDiscount, *model.AppError) {
	orderDiscounts, appErr := a.srv.DiscountService().OrderDiscountsByOption(&model.OrderDiscountFilterOption{
		Conditions: squirrel.Eq{
			model.OrderDiscountTableName + ".Type":    model.MANUAL,
			model.OrderDiscountTableName + ".OrderID": ord.Id,
		},
	})
	if appErr != nil {
		return nil, appErr
	}

	return orderDiscounts, nil
}

// CreateOrderDiscountForOrder Add new order discount and update the prices
func (a *ServiceOrder) CreateOrderDiscountForOrder(transaction *gorm.DB, ord *model.Order, reason string, valueType model.DiscountValueType, value *decimal.Decimal) (*model.OrderDiscount, *model.AppError) {
	ord.PopulateNonDbFields()

	netTotal, err := a.ApplyDiscountToValue(value, valueType, ord.Currency, ord.Total.Net)
	if err != nil {
		return nil, model.NewAppError("CreateOrderDiscountForOrder", model.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
	}
	grossTotal, err := a.ApplyDiscountToValue(value, valueType, ord.Currency, ord.Total.Gross)
	if err != nil {
		return nil, model.NewAppError("CreateOrderDiscountForOrder", model.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	sub, _ := ord.Total.Sub(grossTotal.(*goprices.Money))

	newOrderDiscount, appErr := a.srv.DiscountService().UpsertOrderDiscount(transaction, &model.OrderDiscount{
		ValueType: valueType,
		Value:     value,
		Reason:    &reason,
		Amount:    sub.Gross,
	})
	if appErr != nil {
		return nil, appErr
	}

	newOrderTotal, err := goprices.NewTaxedMoney(netTotal.(*goprices.Money), grossTotal.(*goprices.Money))
	if err != nil {
		return nil, model.NewAppError("CreateOrderDiscountForOrder", model.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	ord.Total = newOrderTotal
	_, appErr = a.UpsertOrder(transaction, ord)
	if appErr != nil {
		return nil, appErr
	}

	return newOrderDiscount, nil
}

// RemoveOrderDiscountFromOrder Remove the order discount from order and update the prices.
func (a *ServiceOrder) RemoveOrderDiscountFromOrder(transaction *gorm.DB, ord *model.Order, orderDiscount *model.OrderDiscount) *model.AppError {
	appErr := a.srv.DiscountService().BulkDeleteOrderDiscounts([]string{orderDiscount.Id})
	if appErr != nil {
		return appErr
	}

	ord.PopulateNonDbFields()
	orderDiscount.PopulateNonDbFields()

	newOrderTotal, err := ord.Total.Add(orderDiscount.Amount)
	if err != nil {
		return model.NewAppError("RemoveOrderDiscountFromOrder", model.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	ord.Total = newOrderTotal
	_, appErr = a.UpsertOrder(transaction, ord)
	if appErr != nil {
		return appErr
	}

	return nil
}

// UpdateDiscountForOrderLine Update discount fields for order line. Apply discount to the price
//
// `reason`, `valueType` can be empty. `value` can be nil
func (a *ServiceOrder) UpdateDiscountForOrderLine(orderLine model.OrderLine, ord model.Order, reason string, valueType model.DiscountValueType, value *decimal.Decimal, manager interfaces.PluginManagerInterface, taxIncluded bool) *model.AppError {
	ord.PopulateNonDbFields()
	orderLine.PopulateNonDbFields()

	if reason != "" {
		orderLine.UnitDiscountReason = &reason
	}
	if value == nil {
		value = orderLine.UnitDiscountValue
	}
	if !valueType.IsValid() {
		valueType = orderLine.UnitDiscountType
	}

	if orderLine.UnitDiscountValue != value || orderLine.UnitDiscountType != valueType {
		unitPriceWithDiscount, err := a.ApplyDiscountToValue(value, valueType, orderLine.UnDiscountedUnitPrice.Currency, orderLine.UnDiscountedUnitPrice)
		if err != nil {
			return model.NewAppError("UpdateDiscountForOrderLine", model.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
		}

		newOrderLineUnitDiscount, err := orderLine.UnDiscountedUnitPrice.Sub(unitPriceWithDiscount)
		if err != nil {
			return model.NewAppError("UpdateDiscountForOrderLine", model.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
		}

		orderLine.UnitDiscount = newOrderLineUnitDiscount.Gross
		orderLine.UnitPrice = unitPriceWithDiscount.(*goprices.TaxedMoney)
		orderLine.UnitDiscountType = valueType
		orderLine.UnitDiscountValue = value
		orderLine.TotalPrice = orderLine.UnitPrice.Mul(float64(orderLine.Quantity))
		orderLine.UnDiscountedUnitPrice, _ = orderLine.UnitPrice.Sub(orderLine.UnitDiscount)
		orderLine.UnDiscountedTotalPrice = orderLine.UnDiscountedUnitPrice.Mul(float64(orderLine.Quantity))

	}

	// Save lines before calculating the taxes as some plugin can fetch all order data from db
	_, appErr := a.UpsertOrderLine(nil, &orderLine)
	if appErr != nil {
		return appErr
	}

	appErr = a.UpdateTaxesForOrderLine(orderLine, ord, manager, taxIncluded)
	if appErr != nil {
		return appErr
	}

	_, appErr = a.UpsertOrderLine(nil, &orderLine)
	return appErr
}

// RemoveDiscountFromOrderLine Drop discount applied to order line. Restore undiscounted price
func (a *ServiceOrder) RemoveDiscountFromOrderLine(orderLine model.OrderLine, ord model.Order, manager interfaces.PluginManagerInterface, taxIncluded bool) *model.AppError {
	orderLine.PopulateNonDbFields()

	orderLine.UnitPrice = orderLine.UnDiscountedUnitPrice
	orderLine.UnitDiscountAmount = &decimal.Zero
	orderLine.UnitDiscountValue = &decimal.Zero
	orderLine.UnitDiscountReason = model.NewPrimitive("")
	orderLine.TotalPrice = orderLine.UnitPrice.Mul(float64(orderLine.Quantity))

	_, appErr := a.UpsertOrderLine(nil, &orderLine)
	if appErr != nil {
		return appErr
	}

	appErr = a.UpdateTaxesForOrderLine(orderLine, ord, manager, taxIncluded)
	if appErr != nil {
		return appErr
	}

	_, appErr = a.UpsertOrderLine(nil, &orderLine)
	return appErr
}

// ValidateDraftOrder checks if the given order contains the proper data.
//
//	// Has proper customer data,
//	// Shipping address and method are set up,
//	// Product variants for order lines still exists in database.
//	// Product variants are available in requested quantity.
//	// Product variants are published.
func (s *ServiceOrder) ValidateDraftOrder(order *model.Order) *model.AppError {
	if order.Status != model.ORDER_STATUS_DRAFT {
		return nil
	}

	// validate billing address
	if order.BillingAddressID == nil {
		return model.NewAppError("app.order.ValidateDraftOrder", "app.order.billing_address_not_set.app_error", nil, "Can't finalize draft order without billing address.", http.StatusNotAcceptable)
	}

	orderRequiresShipping, appErr := s.OrderShippingIsRequired(order.Id)
	if appErr != nil {
		return appErr
	}
	if orderRequiresShipping {
		// validate shipping address
		if order.ShippingAddressID == nil {
			return model.NewAppError("app.order.ValidateDraftOrder", "app.order.shipping_address_not_set.app_error", nil, "Can't finalize draft order without shipping address.", http.StatusNotAcceptable)
		}

		// validate shipping method
		if order.ShippingMethodID == nil {
			return model.NewAppError("app.order.ValidateDraftOrder", "app.order.shipping_method_required.app_error", nil, "shipping method is required", http.StatusNotAcceptable)
		}
	}

	// validate order lines
	orderLinesOfOrder, appErr := s.OrderLinesByOption(&model.OrderLineFilterOption{
		Conditions: squirrel.Eq{model.OrderLineTableName + ".OrderID": order.Id},
		Preload:    []string{"ProductVariant"},
	})
	if appErr != nil {
		return appErr
	}
	countryCode, appErr := s.GetOrderCountry(order)
	if appErr != nil {
		return appErr
	}

	// validate channel is active
	channel, appErr := s.srv.ChannelService().ChannelByOption(&model.ChannelFilterOption{
		Conditions: squirrel.Eq{model.ChannelTableName + ".Id": order.ChannelID},
	})
	if appErr != nil {
		return appErr
	}
	if !channel.IsActive {
		return model.NewAppError("app.Order.ValidateDraftOrder", "app.channel.in_active.app_error", nil, "channel is inactive", http.StatusNotAcceptable)
	}

	for _, line := range orderLinesOfOrder {
		if line.VariantID == nil {
			return model.NewAppError("app.order.ValidateDraftOrder", "app.order.variant_required_for_order_line.app_error", nil, "can not create orders without product", http.StatusNotAcceptable)

		} else if *line.ProductVariant.TrackInventory {
			insufficientStockErr, appErr := s.srv.
				WarehouseService().
				CheckStockAndPreorderQuantity(line.ProductVariant, countryCode, channel.Slug, line.Quantity)

			if appErr != nil {
				return appErr
			}
			if insufficientStockErr != nil {
				return model.NewAppError("app.order.ValidateDraftOrder", "app.order.stock_insufficient.app_error", nil, "stock insufficient", http.StatusNotAcceptable)
			}
		}
	}

	// check if there is a selected product variant that belongs to an unpublished product
	productIDsMap := map[string]bool{} // keys are product ids
	totalQuantity := 0
	productVariantIDs := util.AnyArray[string]{}

	for _, line := range orderLinesOfOrder {
		variant := line.ProductVariant
		totalQuantity += line.Quantity
		if variant != nil {
			productIDsMap[variant.ProductID] = true
			productVariantIDs = append(productVariantIDs, variant.Id)
		}
	}

	// validate total quantity must > 0
	if totalQuantity == 0 {
		return model.NewAppError("app.order.ValidateDraftOrder", "app.order.total_quantity_zero.app_error", nil, "cannot create order without product", http.StatusNotAcceptable)
	}

	notPublishedProducts, err := s.srv.Store.Product().NotPublishedProducts(channel.Id)
	if err != nil {
		return model.NewAppError("app.order.ValidateDraftOrder", "app.product.error_finding_not_published_products.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	if lo.SomeBy(notPublishedProducts, func(item *struct {
		model.Product
		IsPublished     bool
		PublicationDate *time.Time
	}) bool {
		return productIDsMap[item.Id]
	}) {
		return model.NewAppError("app.order.ValidateDraftOrder", "app.order.selected_un_published_product.app_error", nil, "you can't buy unpublished products", http.StatusNotAcceptable)
	}

	// validate product available for purchase
	productChannelListings, appErr := s.srv.
		ProductService().
		ProductChannelListingsByOption(&model.ProductChannelListingFilterOption{
			Conditions: squirrel.Eq{
				model.ProductChannelListingTableName + ".ChannelID": order.ChannelID,
				model.ProductChannelListingTableName + ".ProductID": lo.Keys(productIDsMap),
			},
		})
	if appErr != nil {
		return appErr
	}

	availableForPurchase := lo.SomeBy(productChannelListings, func(item *model.ProductChannelListing) bool { return item != nil && item.IsAvailableForPurchase() })
	if !availableForPurchase {
		return model.NewAppError("app.order.ValidateDraftOrder", "app.order.product_not_avaiable_for_purchase.app_error", nil, "product is not available for purchase", http.StatusNotAcceptable)
	}

	// validate variants is available:
	variantChannelListings, appErr := s.srv.ProductService().ProductVariantChannelListingsByOption(&model.ProductVariantChannelListingFilterOption{
		Conditions: squirrel.And{
			squirrel.Eq{
				model.ProductVariantChannelListingTableName + ".VariantID": productVariantIDs,
				model.ProductVariantChannelListingTableName + ".ChannelID": order.ChannelID,
			},
			squirrel.Expr(model.ProductVariantChannelListingTableName + ".PriceAmount IS NOT NULL"),
		},
	})
	if appErr != nil {
		return appErr
	}

	if productVariantIDs.Dedup().Len() > variantChannelListings.VariantIDs().Dedup().Len() {
		return model.NewAppError("app.order.ValidateDraftOrder", "app.order.variant_not_available.app_error", nil, "product variant not available for purchase", http.StatusNotAcceptable)
	}

	return nil
}

// ValidateProductIsPublishedInChannel checks if some of given variants belong to unpublished products
func (s *ServiceOrder) ValidateProductIsPublishedInChannel(variants model.ProductVariants, channelID string) *model.AppError {
	var unPublishedProductsWithData, err = s.srv.Store.Product().NotPublishedProducts(channelID)
	if err != nil {
		return model.NewAppError("ValidateProductIsPublishedInChannel", "app.product.error_finding_not_published_products.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	var unPublisgedProductIdMap = lo.SliceToMap(unPublishedProductsWithData, func(item *struct {
		model.Product
		IsPublished     bool
		PublicationDate *time.Time
	}) (string, bool) {
		return item.Id, true
	})

	variantIdsThatAreUnPublished := []string{}
	for _, item := range variants {
		if unPublisgedProductIdMap[item.ProductID] {
			variantIdsThatAreUnPublished = append(variantIdsThatAreUnPublished, item.Id)
		}
	}

	if len(variantIdsThatAreUnPublished) > 0 {
		return model.NewAppError("ValidateProductIsPublishedInChannel", "app.order.add_unpublished_variants_to_order.app_error", map[string]interface{}{"Variants": strings.Join(variantIdsThatAreUnPublished, ", ")}, "cannot add unpublished variants to order", http.StatusNotAcceptable)
	}

	return nil
}

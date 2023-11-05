package checkout

import (
	"encoding/json"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/Masterminds/squirrel"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app/plugin/interfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"gorm.io/gorm"
)

// getVoucherDataForOrder Fetch, process and return voucher/discount data from checkout.
// Careful! It should be called inside a transaction.
// :raises NotApplicable: When the voucher is not applicable in the current checkout.
func (s *ServiceCheckout) getVoucherDataForOrder(checkoutInfo model.CheckoutInfo) (map[string]*model.Voucher, *model.NotApplicable, *model.AppError) {
	checkout := checkoutInfo.Checkout
	voucher, appErr := s.GetVoucherForCheckout(checkoutInfo, nil, true)
	if appErr != nil {
		return nil, nil, appErr
	}

	if checkout.VoucherCode != nil && voucher == nil {
		return nil, model.NewNotApplicable("getVoucherDataForOrder", "Voucher expired in meantime. Order placement aborted", nil, 0), nil
	}

	if voucher == nil {
		return map[string]*model.Voucher{}, nil, nil
	}

	if voucher.UsageLimit != nil && *voucher.UsageLimit != 0 {
		appErr = s.srv.DiscountService().IncreaseVoucherUsage(voucher)
		if appErr != nil {
			return nil, nil, appErr
		}
	}

	if voucher.ApplyOncePerCustomer {
		notApplicable, appErr := s.srv.DiscountService().AddVoucherUsageByCustomer(voucher, checkoutInfo.GetCustomerEmail())
		if notApplicable != nil || appErr != nil {
			return nil, notApplicable, appErr
		}
	}

	return map[string]*model.Voucher{"voucher": voucher}, nil, nil
}

// processShippingDataForOrder Fetch, process and return shipping data from checkout.
func (s *ServiceCheckout) processShippingDataForOrder(checkoutInfo model.CheckoutInfo, shippingPrice *goprices.TaxedMoney, manager interfaces.PluginManagerInterface, lines []*model.CheckoutLineInfo) (map[string]interface{}, *model.AppError) {
	var (
		deliveryMethodInfo  = checkoutInfo.DeliveryMethodInfo
		shippingAddress     = deliveryMethodInfo.GetShippingAddress()
		copyShippingAddress *model.Address
		appErr              *model.AppError
	)

	if checkoutInfo.User != nil && shippingAddress != nil {
		appErr = s.srv.AccountService().StoreUserAddress(checkoutInfo.User, *shippingAddress, model.ADDRESS_TYPE_SHIPPING, manager)
		if appErr != nil {
			return nil, appErr
		}

		addressesOfUser, appErr := s.srv.AccountService().AddressesByOption(&model.AddressFilterOption{
			Conditions: squirrel.Eq{model.AddressTableName + ".Id": shippingAddress.Id},
			UserID:     squirrel.Eq{model.UserAddressTableName + ".user_id": checkoutInfo.User.Id},
		})
		if appErr != nil {
			if appErr.StatusCode != http.StatusNotFound {
				return nil, appErr
			}
		}

		if len(addressesOfUser) > 0 {
			copyShippingAddress, appErr = s.srv.AccountService().CopyAddress(shippingAddress)
			if appErr != nil {
				return nil, appErr
			}
		}
	}

	checkoutTotalWeight, appErr := s.srv.CheckoutService().CheckoutTotalWeight(lines)
	if appErr != nil {
		return nil, appErr
	}

	result := map[string]interface{}{
		deliveryMethodInfo.GetOrderKey(): deliveryMethodInfo.GetDeliveryMethod(),
	}

	if copyShippingAddress != nil {
		result["shipping_address"] = copyShippingAddress
	} else {
		result["shipping_address"] = shippingAddress
	}

	result["shipping_price"] = shippingPrice
	result["weight"] = checkoutTotalWeight

	for key, value := range deliveryMethodInfo.DeliveryMethodName() {
		result[key] = value
	}

	return result, nil
}

// processUserDataForOrder Fetch, process and return shipping data from checkout.
func (s *ServiceCheckout) processUserDataForOrder(checkoutInfo model.CheckoutInfo, manager interfaces.PluginManagerInterface) (map[string]interface{}, *model.AppError) {
	var (
		billingAddress     = checkoutInfo.BillingAddress
		copyBillingAddress *model.Address
		appErr             *model.AppError
	)

	if checkoutInfo.User != nil && billingAddress != nil {
		appErr = s.srv.AccountService().StoreUserAddress(checkoutInfo.User, *billingAddress, model.ADDRESS_TYPE_BILLING, manager)
		if appErr != nil {
			return nil, appErr
		}

		billingAddressOfUser, appErr := s.srv.AccountService().AddressesByOption(&model.AddressFilterOption{
			UserID:     squirrel.Eq{model.UserAddressTableName + ".user_id": checkoutInfo.User.Id},
			Conditions: squirrel.Eq{model.AddressTableName + ".Id": billingAddress.Id},
		})
		if appErr != nil && appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr
		}

		if len(billingAddressOfUser) > 0 {
			copyBillingAddress, appErr = s.srv.AccountService().CopyAddress(billingAddress)
			if appErr != nil {
				return nil, appErr
			}
		}
	}

	if copyBillingAddress == nil {
		copyBillingAddress = billingAddress
	}

	return map[string]interface{}{
		"user":            checkoutInfo.User,
		"user_email":      checkoutInfo.GetCustomerEmail(),
		"billing_address": copyBillingAddress,
		"customer_note":   checkoutInfo.Checkout.Note,
	}, nil
}

// validateGiftcards Check if all gift cards assigned to checkout are available.
func (s *ServiceCheckout) validateGiftcards(checkout model.Checkout) (*model.NotApplicable, *model.AppError) {
	var (
		totalGiftcardsOfCheckout       int
		totalActiveGiftcardsOfCheckout int
		startOfToday                   = util.StartOfDay(time.Now().UTC())
	)

	_, allGiftcards, appErr := s.srv.GiftcardService().GiftcardsByOption(&model.GiftCardFilterOption{
		CheckoutToken: squirrel.Eq{model.GiftcardCheckoutTableName + ".CheckoutID": checkout.Token},
		Distinct:      true,
	})
	if appErr != nil {
		return nil, appErr
	}

	totalGiftcardsOfCheckout = len(allGiftcards)

	// find active giftcards
	// NOTE: active giftcards are active and has (ExpiryDate IS NULL || ExpiryDate >= beginning of Today)
	var expiryDateOfGiftcard *time.Time
	for _, giftcard := range allGiftcards {
		expiryDateOfGiftcard = giftcard.ExpiryDate
		if (expiryDateOfGiftcard == nil || util.StartOfDay(*expiryDateOfGiftcard).Equal(startOfToday) || util.StartOfDay(*expiryDateOfGiftcard).After(startOfToday)) && *giftcard.IsActive {
			totalActiveGiftcardsOfCheckout++
		}
	}

	if totalActiveGiftcardsOfCheckout != totalGiftcardsOfCheckout {
		return model.NewNotApplicable("validateGiftcards", "Gift card has expired. Order placement cancelled.", nil, 0), nil
	}

	return nil, nil
}

// createLineForOrder Create a line for the given order.
// :raises InsufficientStock: when there is not enough items in stock for this variant.
func (s *ServiceCheckout) createLineForOrder(
	manager interfaces.PluginManagerInterface,
	checkoutInfo model.CheckoutInfo,
	lines []*model.CheckoutLineInfo,
	checkoutLineInfo model.CheckoutLineInfo,
	discounts []*model.DiscountInfo,
	productsTranslation map[string]string,
	variantsTranslation map[string]string,

) (*model.OrderLineData, *model.AppError) {

	var (
		checkoutLine          = checkoutLineInfo.Line
		quantity              = checkoutLine.Quantity
		variant               = checkoutLineInfo.Variant
		product               = checkoutLineInfo.Product
		address               = checkoutInfo.ShippingAddress
		productName           = product.String()
		variantName           = variant.String()
		translatedProductName = productsTranslation[product.Id]
		translatedVariantName = variantsTranslation[variant.Id]
	)
	if address == nil {
		address = checkoutInfo.BillingAddress
	}
	if translatedProductName == productName {
		translatedProductName = ""
	}
	if translatedVariantName == variantName {
		translatedVariantName = ""
	}

	totalLinePrice, appErr := manager.CalculateCheckoutLineTotal(checkoutInfo, lines, checkoutLineInfo, address, discounts)
	if appErr != nil {
		return nil, appErr
	}

	unitPrice, appErr := manager.CalculateCheckoutLineUnitPrice(*totalLinePrice, quantity, checkoutInfo, lines, checkoutLineInfo, address, discounts)
	if appErr != nil {
		return nil, appErr
	}

	taxRate, appErr := manager.GetCheckoutLineTaxRate(checkoutInfo, lines, checkoutLineInfo, address, discounts, *unitPrice)
	if appErr != nil {
		return nil, appErr
	}

	productVariantRequireShipping, appErr := s.srv.ProductService().ProductsRequireShipping([]string{variant.ProductID})
	if appErr != nil {
		return nil, appErr
	}

	orderLine := model.OrderLine{
		ProductName:           productName,
		VariantName:           variantName,
		TranslatedProductName: translatedProductName,
		TranslatedVariantName: translatedVariantName,
		ProductSku:            &variant.Sku,
		IsShippingRequired:    productVariantRequireShipping,
		Quantity:              quantity,
		VariantID:             &variant.Id,
		UnitPrice:             unitPrice,
		TotalPrice:            totalLinePrice,
		TaxRate:               taxRate,
	}

	return &model.OrderLineData{
		Line:        orderLine,
		Quantity:    quantity,
		Variant:     &variant,
		WarehouseID: model.GetPointerOfValue(checkoutInfo.DeliveryMethodInfo.WarehousePK()),
	}, nil
}

// createLinesForOrder Create a lines for the given order.
// :raises InsufficientStock: when there is not enough items in stock for this variant.
func (s *ServiceCheckout) createLinesForOrder(manager interfaces.PluginManagerInterface, checkoutInfo model.CheckoutInfo, lines model.CheckoutLineInfos, discounts []*model.DiscountInfo) ([]*model.OrderLineData, *model.InsufficientStock, *model.AppError) {
	var (
		translationLanguageCode = checkoutInfo.Checkout.LanguageCode
		countryCode             = checkoutInfo.GetCountry()
		variants                model.ProductVariants
		quantities              []int
		products                model.Products
	)

	lines = lines.FilterNils()

	for _, lineInfo := range lines {
		variants = append(variants, &lineInfo.Variant)
		quantities = append(quantities, lineInfo.Line.Quantity)
		products = append(products, &lineInfo.Product)
	}

	productTranslations, appErr := s.srv.ProductService().ProductTranslationsByOption(&model.ProductTranslationFilterOption{
		Conditions: squirrel.Eq{
			model.ProductTranslationTableName + ".ProductID":    products.IDs(),
			model.ProductTranslationTableName + ".LanguageCode": translationLanguageCode,
		},
	})
	if appErr != nil && appErr.StatusCode != http.StatusNotFound {
		return nil, nil, appErr
		// ignore not found error here
	}

	// productTranslationMap has keys are product ids, values are translated product name
	var productTranslationMap = map[string]string{}
	if len(productTranslations) > 0 {
		for _, item := range productTranslations {
			productTranslationMap[item.ProductID] = item.Name
		}
	}

	variantTranslations, appErr := s.srv.ProductService().ProductVariantTranslationsByOption(&model.ProductVariantTranslationFilterOption{
		Conditions: squirrel.Eq{
			model.ProductVariantTranslationTableName + ".ProductVariantID": variants.IDs(),
			model.ProductVariantTranslationTableName + ".LanguageCode":     translationLanguageCode,
		},
	})
	if appErr != nil && appErr.StatusCode != http.StatusNotFound {
		return nil, nil, appErr
		// ignore not found error here
	}

	// productVariantTranslationMap has keys are product variant ids, values are translated product variant names
	var productVariantTranslationMap = map[string]string{}
	if len(variantTranslations) > 0 {
		for _, item := range variantTranslations {
			productVariantTranslationMap[item.ProductVariantID] = item.Name
		}
	}

	additionalWarehouseLookup := checkoutInfo.DeliveryMethodInfo.GetWarehouseFilterLookup()
	insufficientStockErr, appErr := s.srv.WarehouseService().CheckStockAndPreorderQuantityBulk(
		variants,
		countryCode,
		quantities,
		checkoutInfo.Channel.Slug,
		additionalWarehouseLookup,
		nil,
		false,
	)
	if insufficientStockErr != nil || appErr != nil {
		return nil, insufficientStockErr, appErr
	}

	var (
		orderLineDatas []*model.OrderLineData
		appErrorChan   = make(chan *model.AppError)
		dataChan       = make(chan *model.OrderLineData)
		atomicValue    atomic.Int32
	)
	defer func() {
		close(appErrorChan)
		close(dataChan)
	}()

	atomicValue.Add(int32(len(lines))) // specify number of go-routines to wait

	for _, item := range lines {
		go func(lineInfo *model.CheckoutLineInfo) {
			defer atomicValue.Add(-1)

			orderLineData, appErr := s.createLineForOrder(manager, checkoutInfo, lines, *lineInfo, discounts, productTranslationMap, productVariantTranslationMap)
			if appErr != nil {
				appErrorChan <- appErr
				return
			}

			dataChan <- orderLineData
		}(item)
	}

	for atomicValue.Load() != 0 {
		select {
		case err := <-appErrorChan:
			return nil, nil, err
		case data := <-dataChan:
			orderLineDatas = append(orderLineDatas, data)
		default:
		}
	}

	return orderLineDatas, nil, nil
}

// prepareOrderData Run checks and return all the data from a given checkout to create an order.
// :raises NotApplicable InsufficientStock:
func (s *ServiceCheckout) prepareOrderData(manager interfaces.PluginManagerInterface, checkoutInfo model.CheckoutInfo, lines model.CheckoutLineInfos, discounts []*model.DiscountInfo) (map[string]interface{}, *model.InsufficientStock, *model.NotApplicable, *model.TaxError, *model.AppError) {
	checkout := checkoutInfo.Checkout
	checkout.PopulateNonDbFields() // this call is important

	orderData := model.StringInterface{}

	address := checkoutInfo.ShippingAddress
	if address == nil {
		address = checkoutInfo.BillingAddress
	}

	taxedTotal, appErr := s.CheckoutTotal(manager, checkoutInfo, lines, address, discounts)
	if appErr != nil {
		return nil, nil, nil, nil, appErr
	}

	cardsTotal, appErr := s.CheckoutTotalGiftCardsBalance(&checkout)
	if appErr != nil {
		return nil, nil, nil, nil, appErr
	}

	newTaxedTotalGross, err1 := taxedTotal.Gross.Sub(cardsTotal)
	newTaxedTotalNet, err2 := taxedTotal.Net.Sub(cardsTotal)
	if err1 != nil || err2 != nil {
		var errMsg string
		if err1 != nil {
			errMsg = err1.Error()
		} else {
			errMsg = err2.Error()
		}
		return nil, nil, nil, nil, model.NewAppError("prepareOrderData", model.ErrorCalculatingMoneyErrorID, nil, errMsg, http.StatusInternalServerError)
	}

	taxedTotal.Gross = newTaxedTotalGross
	taxedTotal.Net = newTaxedTotalNet

	zeroTaxedMoney, _ := util.ZeroTaxedMoney(checkout.Currency)
	if taxedTotal.LessThan(zeroTaxedMoney) {
		taxedTotal = zeroTaxedMoney
	}

	undiscountedTotal, _ := taxedTotal.Add(checkout.Discount)

	shippingTotal, appErr := manager.CalculateCheckoutShipping(checkoutInfo, lines, address, discounts)
	if appErr != nil {
		return nil, nil, nil, nil, appErr
	}

	shippingTaxRate, appErr := manager.GetCheckoutShippingTaxRate(checkoutInfo, lines, address, discounts, *shippingTotal)
	if appErr != nil {
		return nil, nil, nil, nil, appErr
	}

	data, appErr := s.processShippingDataForOrder(checkoutInfo, shippingTotal, manager, lines)
	if appErr != nil {
		return nil, nil, nil, nil, appErr
	}

	orderData.Merge(data)

	data, appErr = s.processUserDataForOrder(checkoutInfo, manager)
	if appErr != nil {
		return nil, nil, nil, nil, appErr
	}

	orderData.Merge(data)

	var trackingCode string
	if checkout.TrackingCode != nil {
		trackingCode = *checkout.TrackingCode
	}
	orderData.Merge(model.StringInterface{
		"language_code":      checkout.LanguageCode,
		"tracking_client_id": trackingCode,
		"total":              taxedTotal,
		"undiscounted_total": undiscountedTotal,
		"shipping_tax_rate":  shippingTaxRate,
	})

	orderLinesData, insufficient, appErr := s.createLinesForOrder(manager, checkoutInfo, lines, discounts)
	if insufficient != nil || appErr != nil {
		return nil, insufficient, nil, nil, appErr
	}

	orderData["lines"] = orderLinesData

	// validate checkout gift cards
	notApplicable, appErr := s.validateGiftcards(checkout)
	if notApplicable != nil || appErr != nil {
		return nil, nil, notApplicable, nil, appErr
	}
	// Get voucher data (last) as they require a transaction
	voucherMap, notApplicable, appErr := s.getVoucherDataForOrder(checkoutInfo)
	if notApplicable != nil || appErr != nil {
		return nil, nil, notApplicable, nil, appErr
	}

	for key, value := range voucherMap {
		orderData[key] = value
	}

	taxedMoney, appErr := manager.CalculateCheckoutSubTotal(checkoutInfo, lines, address, discounts)
	if appErr != nil {
		return nil, nil, nil, nil, appErr
	}

	taxedMoney, _ = taxedMoney.Add(shippingTotal)
	taxedMoney, _ = taxedMoney.Sub(checkout.Discount)

	orderData["total_price_left"] = taxedMoney.Gross

	manager.PreprocessOrderCreation(checkoutInfo, discounts, lines)

	return orderData, nil, nil, nil, nil
}

// createOrder Create an order from the checkout.
// Each order will get a private copy of both the billing and the shipping
// address (if shipping).
// If any of the addresses is new and the user is logged in the address
// will also get saved to that user's address book.
// Current user's language is saved in the order so we can later determine
// which language to use when sending email.
//
// NOTE: the unused underscore param originally is `app`, but we are not gonna present the feature in early versions.
func (s *ServiceCheckout) createOrder(transaction *gorm.DB, checkoutInfo model.CheckoutInfo, orderData model.StringInterface, user *model.User, _ interface{}, manager interfaces.PluginManagerInterface, siteSettings model.ShopSettings) (*model.Order, *model.InsufficientStock, *model.AppError) {
	checkout := checkoutInfo.Checkout

	_, orders, appErr := s.srv.OrderService().FilterOrdersByOptions(&model.OrderFilterOption{
		Conditions: squirrel.Eq{model.OrderTableName + ".CheckoutToken": checkout.Token},
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, nil, appErr
		}
	}

	if len(orders) > 0 {
		return orders[0], nil, nil
	}

	totalPriceLeft := orderData.Pop("total_price_left").(*goprices.Money)
	orderLinesInfo := orderData.Pop("lines").([]*model.OrderLineData)

	status := model.ORDER_STATUS_UNCONFIRMED
	if *s.srv.Config().ShopSettings.AutomaticallyConfirmAllNewOrders {
		status = model.ORDER_STATUS_UNFULFILLED
	}

	serializedOrderData, err := json.Marshal(orderData)
	if err != nil {
		return nil, nil, model.NewAppError("createOrder", model.ErrorCalculatingMeasurementID, nil, err.Error(), http.StatusInternalServerError)
	}

	// define new order to create
	var newOrder model.Order
	err = json.Unmarshal(serializedOrderData, &newOrder)
	if err != nil {
		return nil, nil, model.NewAppError("createOrder", model.ErrorUnMarshallingDataID, nil, err.Error(), http.StatusInternalServerError)
	}

	newOrder.Id = ""
	newOrder.CheckoutToken = checkout.Token
	newOrder.Status = status
	newOrder.Origin = model.ORDER_ORIGIN_CHECKOUT
	newOrder.ChannelID = checkoutInfo.Channel.Id

	createdNewOrder, appErr := s.srv.OrderService().UpsertOrder(transaction, &newOrder)
	if appErr != nil {
		return nil, nil, appErr
	}

	if checkout.DiscountAmount != nil {
		// store voucher as a fixed value as it this the simplest solution for now.
		// This will be solved when we refactor the voucher logic to use .discounts
		// relations
		_, appErr := s.srv.DiscountService().UpsertOrderDiscount(transaction, &model.OrderDiscount{
			Type:           model.VOUCHER,
			ValueType:      model.DISCOUNT_VALUE_TYPE_FIXED,
			Value:          checkout.DiscountAmount,
			Name:           checkout.DiscountName,
			TranslatedName: checkout.TranslatedDiscountName,
			Currency:       checkout.Currency,
			AmountValue:    checkout.DiscountAmount,
			OrderID:        &createdNewOrder.Id,
		})
		if appErr != nil {
			return nil, nil, appErr
		}
	}

	var orderLines []*model.OrderLine
	for _, lineInfo := range orderLinesInfo {
		line := lineInfo.Line
		line.OrderID = createdNewOrder.Id

		orderLines = append(orderLines, &line)
	}

	_, appErr = s.srv.OrderService().BulkUpsertOrderLines(transaction, orderLines)
	if appErr != nil {
		return nil, nil, appErr
	}

	var (
		countryCode               = checkoutInfo.GetCountry()
		additionalWarehouseLookup = checkoutInfo.DeliveryMethodInfo.GetWarehouseFilterLookup()
	)

	insufficientStockErr, appErr := s.srv.WarehouseService().AllocateStocks(orderLinesInfo, countryCode, checkoutInfo.Channel.Slug, manager, additionalWarehouseLookup)
	if insufficientStockErr != nil || appErr != nil {
		return nil, insufficientStockErr, appErr
	}

	insufficientStockErr, appErr = s.srv.WarehouseService().AllocatePreOrders(orderLinesInfo, checkoutInfo.Channel.Slug)
	if insufficientStockErr != nil || appErr != nil {
		return nil, insufficientStockErr, appErr
	}

	appErr = s.srv.OrderService().AddGiftcardsToOrder(transaction, checkoutInfo, createdNewOrder, totalPriceLeft, user, nil)
	if appErr != nil {
		return nil, nil, appErr
	}

	// assign checkout payments to other order
	appErr = s.srv.PaymentService().UpdatePaymentsOfCheckout(transaction, checkout.Token, &model.PaymentPatch{OrderID: createdNewOrder.Id})
	if appErr != nil {
		return nil, nil, appErr
	}

	// copy metadata from the checkout into the new order
	createdNewOrder.Metadata = checkout.Metadata.DeepCopy()
	createdNewOrder.RedirectUrl = checkout.RedirectURL
	createdNewOrder.PrivateMetadata = checkout.PrivateMetadata.DeepCopy()

	appErr = s.srv.OrderService().UpdateOrderTotalPaid(transaction, createdNewOrder)
	if appErr != nil {
		return nil, nil, appErr
	}

	_, appErr = s.srv.OrderService().UpsertOrder(transaction, createdNewOrder)
	if appErr != nil {
		return nil, nil, appErr
	}

	if *siteSettings.AutomaticallyFulfillNonShippableGiftcard {
		_, insufficientStockErr, appErr = s.srv.GiftcardService().FulfillNonShippableGiftcards(createdNewOrder, orderLines, siteSettings, user, nil, manager)
		if insufficientStockErr != nil || appErr != nil {
			return nil, nil, appErr
		}
	}

	insufficientStock, appErr := s.srv.OrderService().OrderCreated(transaction, *createdNewOrder, user, nil, manager, false)
	if insufficientStock != nil || appErr != nil {
		return nil, insufficientStock, appErr
	}

	// Send the order confirmation email
	var redirectURL string
	if checkout.RedirectURL != nil {
		redirectURL = *checkout.RedirectURL
	}
	appErr = s.srv.OrderService().SendOrderConfirmation(createdNewOrder, redirectURL, manager)
	if appErr != nil {
		return nil, nil, appErr
	}

	return createdNewOrder, nil, nil
}

// prepareCheckout Prepare checkout object to complete the checkout process.
func (s *ServiceCheckout) prepareCheckout(transaction *gorm.DB, manager interfaces.PluginManagerInterface, checkoutInfo model.CheckoutInfo, lines []*model.CheckoutLineInfo, discoutns []*model.DiscountInfo, trackingCode string, redirectURL string, payMent *model.Payment) (*model.PaymentError, *model.AppError) {
	checkout := checkoutInfo.Checkout

	appErr := s.CleanCheckoutShipping(checkoutInfo, lines)
	if appErr != nil {
		return nil, appErr
	}

	paymentErr, appErr := s.CleanCheckoutPayment(transaction, manager, checkoutInfo, lines, discoutns, payMent)
	if paymentErr != nil || appErr != nil {
		return paymentErr, appErr
	}

	if !checkoutInfo.Channel.IsActive {
		return nil, model.NewAppError("prepareCheckout", "app.checkout.channel_inactive.app_error", nil, "", http.StatusNotAcceptable)
	}
	if redirectURL != "" {
		appErr = model.ValidateStoreFrontUrl(s.srv.Config(), redirectURL)
		if appErr != nil {
			return nil, appErr
		}
	}

	var needUpdate bool
	if redirectURL != "" && (checkout.RedirectURL == nil || *checkout.RedirectURL != redirectURL) {
		checkout.RedirectURL = &redirectURL
		needUpdate = true
	}
	if trackingCode != "" && (checkout.TrackingCode == nil || *checkout.TrackingCode != trackingCode) {
		checkout.TrackingCode = &trackingCode
		needUpdate = true
	}

	if needUpdate {
		_, appErr = s.UpsertCheckouts(transaction, []*model.Checkout{&checkout})
		if appErr != nil {
			return nil, appErr
		}
	}

	return nil, nil
}

// ReleaseVoucherUsage
func (s *ServiceCheckout) ReleaseVoucherUsage(orderData map[string]interface{}) *model.AppError {
	if iface, ok := orderData["voucher"]; ok && iface != nil {
		voucher := iface.(*model.Voucher)

		if voucher.UsageLimit != nil && *voucher.UsageLimit != 0 {
			appErr := s.srv.DiscountService().DecreaseVoucherUsage(voucher)
			if appErr != nil {
				return appErr
			}

			if userEmail, ok := orderData["user_email"]; ok {
				appErr = s.srv.DiscountService().RemoveVoucherUsageByCustomer(voucher, userEmail.(string))
				if appErr != nil {
					return appErr
				}
			}
		}
	}

	return nil
}

func (s *ServiceCheckout) getOrderData(manager interfaces.PluginManagerInterface, checkoutInfo model.CheckoutInfo, lines []*model.CheckoutLineInfo, discoutns []*model.DiscountInfo) (map[string]interface{}, *model.AppError) {
	orderData, insufficientStockErr, notApplicableErr, taxError, appErr := s.prepareOrderData(manager, checkoutInfo, lines, discoutns)
	if appErr != nil {
		return nil, appErr
	}

	if insufficientStockErr != nil {
		return nil, s.PrepareInsufficientStockCheckoutValidationAppError("getOrderData", insufficientStockErr)
	}
	if notApplicableErr != nil {
		return nil, model.NewAppError("getOrderData", "app.checkout.voucher_not_applicable.app_error", map[string]interface{}{"code": model.VOUCHER_NOT_APPLICABLE}, notApplicableErr.Error(), 0)
	}
	if taxError != nil {
		return nil, model.NewAppError("getOrderData", "app.checkout.unable_to_calculate_taxes", map[string]interface{}{"code": model.TAX_ERROR}, taxError.Message, 0)
	}
	return orderData, nil
}

// processPayment Process the payment assigned to checkout
func (s *ServiceCheckout) processPayment(dbTransaction *gorm.DB, payMent *model.Payment, customerID *string, storeSource bool, paymentData map[string]interface{}, orderData map[string]interface{}, manager interfaces.PluginManagerInterface, channelSlug string) (*model.PaymentTransaction, *model.PaymentError, *model.AppError) {
	var (
		transaction *model.PaymentTransaction
		paymentErr  *model.PaymentError
		appErr      *model.AppError
		paymentID   = payMent.Id
	)

	if payMent.ToConfirm {
		transaction, paymentErr, appErr = s.srv.PaymentService().Confirm(
			dbTransaction,
			*payMent,
			manager,
			channelSlug,
			paymentData,
		)
	} else {
		transaction, paymentErr, appErr = s.srv.PaymentService().ProcessPayment(
			dbTransaction,
			*payMent,
			payMent.Token,
			manager,
			channelSlug,
			customerID,
			storeSource,
			paymentData,
		)
	}

	if appErr != nil {
		return nil, nil, appErr
	}
	if paymentErr != nil {
		appErr = s.ReleaseVoucherUsage(orderData)
		if appErr != nil {
			return nil, nil, appErr
		}
		return nil, nil, model.NewAppError("processPayment", "app.checkout.payment_error.app_error", nil, paymentErr.Error(), 0)
	}

	_, appErr = s.srv.PaymentService().PaymentByID(nil, paymentID, false)
	if appErr != nil {
		return nil, nil, appErr
	}

	if !transaction.IsSuccess {
		var paymentErrorMessage string
		if transaction.Error != nil {
			paymentErrorMessage = *transaction.Error
		}
		return nil, &model.PaymentError{
			Where:   "processPayment",
			Message: paymentErrorMessage,
		}, nil
	}

	return transaction, nil, nil
}

// Logic required to finalize the checkout and convert it to order.
// Should be used with transaction_with_commit_on_errors, as there is a possibility
// for thread race.
// :raises ValidationError
//
// NOTE: Make sure user is authenticated before calling this method.
func (s *ServiceCheckout) CompleteCheckout(
	dbTransaction *gorm.DB,
	manager interfaces.PluginManagerInterface,
	checkoutInfo model.CheckoutInfo,
	lines []*model.CheckoutLineInfo,
	paymentData map[string]interface{},
	storeSource bool,
	discounts []*model.DiscountInfo,
	user *model.User,
	_ interface{}, // this param originally is `app`, but we not gonna integrate app feature in the early versions
	siteSettings model.ShopSettings,
	trackingCode string,
	redirectURL string,
) (*model.Order, bool, model.StringInterface, *model.PaymentError, *model.AppError) {

	var (
		checkout    = checkoutInfo.Checkout
		channelSlug = checkoutInfo.Channel.Slug
	)

	lastActivePaymentOfCheckout, appErr := s.CheckoutLastActivePayment(&checkout) // NOTE: returned payment still can be nil even when appErr is nil
	if appErr != nil {
		return nil, false, nil, nil, appErr
	}

	paymentErr, appErr := s.prepareCheckout(
		dbTransaction,
		manager,
		checkoutInfo,
		lines,
		discounts,
		trackingCode,
		redirectURL,
		lastActivePaymentOfCheckout,
	)
	if paymentErr != nil || appErr != nil {
		return nil, false, nil, paymentErr, appErr
	}

	orderData, appErr := s.getOrderData(manager, checkoutInfo, lines, discounts)
	if appErr != nil {
		paymentErr, apErr := s.srv.PaymentService().PaymentRefundOrVoid(dbTransaction, lastActivePaymentOfCheckout, manager, channelSlug)
		if paymentErr != nil || apErr != nil {
			return nil, false, nil, paymentErr, apErr
		}

		return nil, false, nil, nil, appErr
	}

	var customerID *string
	if lastActivePaymentOfCheckout != nil && user != nil {
		uuid, appErr := s.srv.PaymentService().FetchCustomerId(user, lastActivePaymentOfCheckout.GateWay)
		if appErr != nil {
			return nil, false, nil, nil, appErr
		}
		if model.IsValidId(uuid) {
			customerID = &uuid
		}
	}

	transaction, paymentErr, appErr := s.processPayment(
		dbTransaction,
		lastActivePaymentOfCheckout,
		customerID,
		storeSource,
		paymentData,
		orderData,
		manager,
		channelSlug,
	)
	if paymentErr != nil || appErr != nil {
		return nil, false, nil, paymentErr, appErr
	}

	if transaction.CustomerID != nil && user != nil {
		appErr = s.srv.PaymentService().StoreCustomerId(user.Id, lastActivePaymentOfCheckout.GateWay, *transaction.CustomerID)
		if appErr != nil {
			return nil, false, nil, nil, appErr
		}
	}

	actionData := transaction.ActionRequiredData
	if !transaction.ActionRequired {
		actionData = make(model.StringInterface)
	}

	var (
		orDer                *model.Order
		insufficientStockErr *model.InsufficientStock
	)
	if !transaction.ActionRequired {
		orDer, insufficientStockErr, appErr = s.createOrder(dbTransaction, checkoutInfo, orderData, user, nil, manager, siteSettings)
		if appErr != nil {
			return nil, false, nil, nil, appErr
		}

		if insufficientStockErr != nil {
			appErr = s.ReleaseVoucherUsage(orderData)
			if appErr != nil {
				return nil, false, nil, nil, appErr
			}

			paymentErr, appErr = s.srv.PaymentService().PaymentRefundOrVoid(dbTransaction, lastActivePaymentOfCheckout, manager, channelSlug)
			if appErr != nil || paymentErr != nil {
				return nil, false, nil, paymentErr, appErr
			}

			return nil, false, nil, nil, s.PrepareInsufficientStockCheckoutValidationAppError("CompleteCheckout", insufficientStockErr)
		}

		// if not appError nor insufficient stock error, remove checkout after order is successfully created:
		appErr = s.DeleteCheckoutsByOption(dbTransaction, &model.CheckoutFilterOption{
			Conditions: squirrel.Eq{model.CheckoutTableName + ".Token": checkout.Token},
		})
		if appErr != nil {
			return nil, false, nil, nil, appErr
		}
	}

	return orDer, transaction.ActionRequired, actionData, nil, nil
}

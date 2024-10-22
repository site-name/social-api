package checkout

import (
	"encoding/json"
	"net/http"
	"sync/atomic"
	"time"

	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app/plugin/interfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/model_types"
	"github.com/sitename/sitename/modules/util"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// getVoucherDataForOrder Fetch, process and return voucher/discount data from checkout.
// Careful! It should be called inside a transaction.
// :raises NotApplicable: When the voucher is not applicable in the current checkout.
func (s *ServiceCheckout) getVoucherDataForOrder(checkoutInfo model_helper.CheckoutInfo) (map[string]*model.Voucher, *model_helper.NotApplicable, *model_helper.AppError) {
	checkout := checkoutInfo.Checkout
	voucher, appErr := s.GetVoucherForCheckout(checkoutInfo, nil, true)
	if appErr != nil {
		return nil, nil, appErr
	}

	if !checkout.VoucherCode.IsNil() && voucher == nil {
		return nil, model_helper.NewNotApplicable("getVoucherDataForOrder", "Voucher expired in meantime. Order placement aborted", nil, 0), nil
	}
	if voucher == nil {
		return map[string]*model.Voucher{}, nil, nil
	}

	if !voucher.UsageLimit.IsNil() {
		voucher, appErr = s.srv.Discount.AlterVoucherUsage(*voucher, 1)
		if appErr != nil {
			return nil, nil, appErr
		}
	}

	if voucher.ApplyOncePerCustomer {
		notApplicable, appErr := s.srv.Discount.AddVoucherUsageByCustomer(*voucher, checkoutInfo.GetCustomerEmail())
		if notApplicable != nil || appErr != nil {
			return nil, notApplicable, appErr
		}
	}

	return map[string]*model.Voucher{"voucher": voucher}, nil, nil
}

// processShippingDataForOrder Fetch, process and return shipping data from checkout.
func (s *ServiceCheckout) processShippingDataForOrder(checkoutInfo model_helper.CheckoutInfo, shippingPrice goprices.TaxedMoney, manager interfaces.PluginManagerInterface, lines model_helper.CheckoutLineInfos) (map[string]any, *model_helper.AppError) {
	var (
		deliveryMethodInfo  = checkoutInfo.DeliveryMethodInfo
		shippingAddress     = deliveryMethodInfo.GetShippingAddress()
		copyShippingAddress *model.Address
		appErr              *model_helper.AppError
	)

	if checkoutInfo.User != nil && shippingAddress != nil {
		appErr = s.srv.Account.StoreUserAddress(*checkoutInfo.User, *shippingAddress, model_helper.ADDRESS_TYPE_SHIPPING, manager)
		if appErr != nil {
			return nil, appErr
		}

		addressesOfUser, appErr := s.srv.Account.AddressesByOption(model_helper.AddressFilterOptions{
			CommonQueryOptions: model_helper.NewCommonQueryOptions(
				model.AddressWhere.ID.EQ(shippingAddress.ID),
				model.AddressWhere.UserID.EQ(checkoutInfo.User.ID),
			),
		})
		if appErr != nil {
			if appErr.StatusCode != http.StatusNotFound {
				return nil, appErr
			}
		}

		if len(addressesOfUser) > 0 {
			copyShippingAddress, appErr = s.srv.Account.CopyAddress(shippingAddress)
			if appErr != nil {
				return nil, appErr
			}
		}
	}

	checkoutTotalWeight, appErr := s.srv.Checkout.CheckoutTotalWeight(lines)
	if appErr != nil {
		return nil, appErr
	}

	result := map[string]any{
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
func (s *ServiceCheckout) processUserDataForOrder(checkoutInfo model_helper.CheckoutInfo, manager interfaces.PluginManagerInterface) (map[string]any, *model_helper.AppError) {
	var (
		billingAddress     = checkoutInfo.BillingAddress
		copyBillingAddress *model.Address
		appErr             *model_helper.AppError
	)

	if checkoutInfo.User != nil && billingAddress != nil {
		appErr = s.srv.Account.StoreUserAddress(*checkoutInfo.User, *billingAddress, model_helper.ADDRESS_TYPE_BILLING, manager)
		if appErr != nil {
			return nil, appErr
		}

		billingAddressOfUser, appErr := s.srv.Account.AddressesByOption(model_helper.AddressFilterOptions{
			CommonQueryOptions: model_helper.NewCommonQueryOptions(
				model.AddressWhere.ID.EQ(billingAddress.ID),
			),
		})
		if appErr != nil && appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr
		}

		if len(billingAddressOfUser) > 0 {
			copyBillingAddress, appErr = s.srv.Account.CopyAddress(billingAddress)
			if appErr != nil {
				return nil, appErr
			}
		}
	}

	if copyBillingAddress == nil {
		copyBillingAddress = billingAddress
	}

	return map[string]any{
		"user":            checkoutInfo.User,
		"user_email":      checkoutInfo.GetCustomerEmail(),
		"billing_address": copyBillingAddress,
		"customer_note":   checkoutInfo.Checkout.Note,
	}, nil
}

// validateGiftcards Check if all gift cards assigned to checkout are available.
func (s *ServiceCheckout) validateGiftcards(checkout model.Checkout) (*model_helper.NotApplicable, *model_helper.AppError) {
	var (
		totalGiftcardsOfCheckout       int
		totalActiveGiftcardsOfCheckout int
		startOfToday                   = util.StartOfDay(time.Now().UTC())
	)

	_, allGiftcards, appErr := s.srv.Giftcard.GiftcardsByOption(model_helper.GiftcardFilterOption{
		CheckoutToken: model.GiftcardCheckoutWhere.CheckoutID.EQ(checkout.Token),
		CommonQueryOptions: model_helper.NewCommonQueryOptions(
			qm.Distinct(model.GiftcardTableColumns.ID),
		),
	})
	if appErr != nil {
		return nil, appErr
	}

	totalGiftcardsOfCheckout = len(allGiftcards)

	// find active giftcards
	// NOTE: active giftcards are active and has (ExpiryDate IS NULL || ExpiryDate >= beginning of Today)
	var expiryDateOfGiftcard *time.Time
	for _, giftcard := range allGiftcards {
		expiryDateOfGiftcard = giftcard.ExpiryDate.Time
		if (expiryDateOfGiftcard == nil || util.StartOfDay(*expiryDateOfGiftcard).Equal(startOfToday) || util.StartOfDay(*expiryDateOfGiftcard).After(startOfToday)) && !giftcard.IsActive.IsNil() && *giftcard.IsActive.Bool {
			totalActiveGiftcardsOfCheckout++
		}
	}

	if totalActiveGiftcardsOfCheckout != totalGiftcardsOfCheckout {
		return model_helper.NewNotApplicable("validateGiftcards", "Gift card has expired. Order placement cancelled.", nil, 0), nil
	}

	return nil, nil
}

// createLineForOrder Create a line for the given order.
// :raises InsufficientStock: when there is not enough items in stock for this variant.
func (s *ServiceCheckout) createLineForOrder(
	manager interfaces.PluginManagerInterface,
	checkoutInfo model_helper.CheckoutInfo,
	lines model_helper.CheckoutLineInfos,
	checkoutLineInfo model_helper.CheckoutLineInfo,
	discounts []*model_helper.DiscountInfo,
	productsTranslation map[string]string,
	variantsTranslation map[string]string,
) (*model_helper.OrderLineData, *model_helper.AppError) {
	var (
		checkoutLine          = checkoutLineInfo.Line
		quantity              = checkoutLine.Quantity
		variant               = checkoutLineInfo.Variant
		product               = checkoutLineInfo.Product
		address               = checkoutInfo.ShippingAddress
		productName           = product.Name
		variantName           = model_helper.ProductVariantString(variant)
		translatedProductName = productsTranslation[product.ID]
		translatedVariantName = variantsTranslation[variant.ID]
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

	productVariantRequireShipping, appErr := s.srv.Product.ProductsRequireShipping([]string{variant.ProductID})
	if appErr != nil {
		return nil, appErr
	}

	orderLine := model.OrderLine{
		ProductName:           productName,
		VariantName:           variantName,
		TranslatedProductName: translatedProductName,
		TranslatedVariantName: translatedVariantName,
		ProductSku:            model_types.NewNullString(variant.Sku),
		IsShippingRequired:    productVariantRequireShipping,
		Quantity:              quantity,
		VariantID:             model_types.NewNullString(variant.ID),
		TaxRate:               model_types.NullDecimal{Decimal: taxRate},
	}

	model_helper.OrderLineSetTotalPrice(&orderLine, *totalLinePrice)
	model_helper.OrderLineSetUnitPrice(&orderLine, *unitPrice)

	return &model_helper.OrderLineData{
		Line:        orderLine,
		Quantity:    quantity,
		Variant:     &variant,
		WarehouseID: model_helper.GetPointerOfValue(checkoutInfo.DeliveryMethodInfo.WarehousePK()),
	}, nil
}

// createLinesForOrder Create a lines for the given order.
// :raises InsufficientStock: when there is not enough items in stock for this variant.
func (s *ServiceCheckout) createLinesForOrder(manager interfaces.PluginManagerInterface, checkoutInfo model_helper.CheckoutInfo, lines model_helper.CheckoutLineInfos, discounts []*model_helper.DiscountInfo) ([]*model_helper.OrderLineData, *model_helper.InsufficientStock, *model_helper.AppError) {
	lines = lines.FilterNils()
	length := len(lines)

	var (
		translationLanguageCode = checkoutInfo.Checkout.LanguageCode
		countryCode             = checkoutInfo.GetCountry()
		variants                = make(model.ProductVariantSlice, length)
		quantities              = make([]int, length)
		productIDs              = make([]string, length)
		variantIDs              = make([]string, length)
	)

	for idx, lineInfo := range lines {
		quantities[idx] = lineInfo.Line.Quantity
		productIDs[idx] = lineInfo.Product.ID
		variants[idx] = &lineInfo.Variant
		variantIDs[idx] = lineInfo.Variant.ID
	}

	productTranslations, appErr := s.srv.Product.ProductTranslationsByOption(model_helper.ProductTranslationFilterOption{
		CommonQueryOptions: model_helper.NewCommonQueryOptions(
			model.ProductTranslationWhere.ProductID.IN(productIDs),
			model.ProductTranslationWhere.LanguageCode.EQ(translationLanguageCode),
		),
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

	variantTranslations, appErr := s.srv.Product.ProductVariantTranslationsByOption(model_helper.ProductVariantTranslationFilterOption{
		CommonQueryOptions: model_helper.NewCommonQueryOptions(
			model.ProductVariantTranslationWhere.ProductVariantID.IN(variantIDs),
			model.ProductVariantTranslationWhere.LanguageCode.EQ(translationLanguageCode),
		),
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
	insufficientStockErr, appErr := s.srv.Warehouse.CheckStockAndPreorderQuantityBulk(
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
		orderLineDatas []*model_helper.OrderLineData
		appErrorChan   = make(chan *model_helper.AppError)
		dataChan       = make(chan *model_helper.OrderLineData)
		atomicValue    atomic.Int32
	)
	defer func() {
		close(appErrorChan)
		close(dataChan)
	}()

	atomicValue.Add(int32(len(lines))) // specify number of go-routines to wait

	for _, item := range lines {
		go func(lineInfo *model_helper.CheckoutLineInfo) {
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

func (s *ServiceCheckout) prepareOrderData(manager interfaces.PluginManagerInterface, checkoutInfo model_helper.CheckoutInfo, lines model_helper.CheckoutLineInfos, discounts []*model_helper.DiscountInfo) (map[string]any, *model_helper.InsufficientStock, *model_helper.NotApplicable, *model_helper.TaxError, *model_helper.AppError) {
	checkout := checkoutInfo.Checkout

	orderData := model_types.JSONString{}

	address := checkoutInfo.ShippingAddress
	if address == nil {
		address = checkoutInfo.BillingAddress
	}

	taxedTotal, appErr := s.CheckoutTotal(manager, checkoutInfo, lines, address, discounts)
	if appErr != nil {
		return nil, nil, nil, nil, appErr
	}

	cardsTotal, appErr := s.CheckoutTotalGiftCardsBalance(checkout)
	if appErr != nil {
		return nil, nil, nil, nil, appErr
	}

	newTaxedTotalGross, err1 := taxedTotal.Gross.Sub(*cardsTotal)
	newTaxedTotalNet, err2 := taxedTotal.Net.Sub(*cardsTotal)
	if err1 != nil || err2 != nil {
		var errMsg string
		if err1 != nil {
			errMsg = err1.Error()
		} else {
			errMsg = err2.Error()
		}
		return nil, nil, nil, nil, model_helper.NewAppError("prepareOrderData", model_helper.ErrorCalculatingMoneyErrorID, nil, errMsg, http.StatusInternalServerError)
	}

	taxedTotal.Gross = *newTaxedTotalGross
	taxedTotal.Net = *newTaxedTotalNet

	zeroTaxedMoney, _ := util.ZeroTaxedMoney(checkout.Currency.String())
	if taxedTotal.LessThan(*zeroTaxedMoney) {
		taxedTotal = zeroTaxedMoney
	}

	undiscountedTotal, _ := taxedTotal.Add(model_helper.CheckoutGetDiscountMoney(checkout))

	shippingTotal, appErr := manager.CalculateCheckoutShipping(checkoutInfo, lines, address, discounts)
	if appErr != nil {
		return nil, nil, nil, nil, appErr
	}

	shippingTaxRate, appErr := manager.GetCheckoutShippingTaxRate(checkoutInfo, lines, address, discounts, *shippingTotal)
	if appErr != nil {
		return nil, nil, nil, nil, appErr
	}

	data, appErr := s.processShippingDataForOrder(checkoutInfo, *shippingTotal, manager, lines)
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
	if !checkout.TrackingCode.IsNil() {
		trackingCode = *checkout.TrackingCode.String
	}
	orderData.Merge(model_types.JSONString{
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
	taxedMoney, _ = taxedMoney.Sub(model_helper.CheckoutGetDiscountMoney(checkout))

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
func (s *ServiceCheckout) createOrder(transaction boil.ContextTransactor, checkoutInfo model_helper.CheckoutInfo, orderData model_types.JSONString, user *model.User, _ any, manager interfaces.PluginManagerInterface, siteSettings model_helper.ShopSettings) (*model.Order, *model_helper.InsufficientStock, *model_helper.AppError) {
	checkout := checkoutInfo.Checkout

	_, orders, appErr := s.srv.Order.FilterOrdersByOptions(model_helper.OrderFilterOption{
		CommonQueryOptions: model_helper.NewCommonQueryOptions(
			model.OrderWhere.CheckoutToken.EQ(checkout.Token),
		),
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
	orderLinesInfo := orderData.Pop("lines").([]*model_helper.OrderLineData)

	status := model.OrderStatusUnconfirmed
	if *s.srv.Config().ShopSettings.AutomaticallyConfirmAllNewOrders {
		status = model.OrderStatusUnfulfilled
	}

	serializedOrderData, err := json.Marshal(orderData)
	if err != nil {
		return nil, nil, model_helper.NewAppError("createOrder", model_helper.ErrorCalculatingMeasurementID, nil, err.Error(), http.StatusInternalServerError)
	}

	// define new order to create
	var newOrder model.Order
	err = json.Unmarshal(serializedOrderData, &newOrder)
	if err != nil {
		return nil, nil, model_helper.NewAppError("createOrder", model_helper.ErrorUnMarshallingDataID, nil, err.Error(), http.StatusInternalServerError)
	}

	newOrder.ID = ""
	newOrder.CheckoutToken = checkout.Token
	newOrder.Status = status
	newOrder.Origin = model.NullOrderOrigin{
		Valid: true,
		Val:   model.OrderOriginCheckout,
	}
	newOrder.ChannelID = checkoutInfo.Channel.ID

	createdNewOrder, appErr := s.srv.Order.UpsertOrder(transaction, &newOrder)
	if appErr != nil {
		return nil, nil, appErr
	}

	// store voucher as a fixed value as it this the simplest solution for now.
	// This will be solved when we refactor the voucher logic to use .discounts
	// relations
	_, appErr = s.srv.Discount.UpsertOrderDiscount(transaction, model.OrderDiscount{
		Type:           model.OrderDiscountTypeVoucher,
		ValueType:      model.DiscountValueTypeFixed,
		Value:          checkout.DiscountAmount,
		Name:           checkout.DiscountName,
		TranslatedName: checkout.TranslatedDiscountName,
		Currency:       checkout.Currency,
		AmountValue:    checkout.DiscountAmount,
		OrderID:        model_types.NewNullString(createdNewOrder.ID),
	})
	if appErr != nil {
		return nil, nil, appErr
	}

	var orderLines model.OrderLineSlice
	for _, lineInfo := range orderLinesInfo {
		line := lineInfo.Line
		line.OrderID = createdNewOrder.ID

		orderLines = append(orderLines, &line)
	}

	_, appErr = s.srv.Order.BulkUpsertOrderLines(transaction, orderLines)
	if appErr != nil {
		return nil, nil, appErr
	}

	var (
		countryCode               = checkoutInfo.GetCountry()
		additionalWarehouseLookup = checkoutInfo.DeliveryMethodInfo.GetWarehouseFilterLookup()
	)

	insufficientStockErr, appErr := s.srv.Warehouse.AllocateStocks(orderLinesInfo, countryCode, checkoutInfo.Channel.Slug, manager, additionalWarehouseLookup)
	if insufficientStockErr != nil || appErr != nil {
		return nil, insufficientStockErr, appErr
	}

	insufficientStockErr, appErr = s.srv.Warehouse.AllocatePreOrders(orderLinesInfo, checkoutInfo.Channel.Slug)
	if insufficientStockErr != nil || appErr != nil {
		return nil, insufficientStockErr, appErr
	}

	appErr = s.srv.Order.AddGiftcardsToOrder(transaction, checkoutInfo, createdNewOrder, totalPriceLeft, user, nil)
	if appErr != nil {
		return nil, nil, appErr
	}

	// assign checkout payments to other order
	appErr = s.srv.Payment.UpdatePaymentsOfCheckout(transaction, checkout.Token, model_helper.PaymentPatch{OrderID: createdNewOrder.ID})
	if appErr != nil {
		return nil, nil, appErr
	}

	// copy metadata from the checkout into the new order
	createdNewOrder.Metadata = checkout.Metadata.DeepCopy()
	createdNewOrder.RedirectURL = checkout.RedirectURL
	createdNewOrder.PrivateMetadata = checkout.PrivateMetadata.DeepCopy()

	appErr = s.srv.Order.UpdateOrderTotalPaid(transaction, createdNewOrder)
	if appErr != nil {
		return nil, nil, appErr
	}

	_, appErr = s.srv.Order.UpsertOrder(transaction, createdNewOrder)
	if appErr != nil {
		return nil, nil, appErr
	}

	if *siteSettings.AutomaticallyFulfillNonShippableGiftcard {
		_, insufficientStockErr, appErr = s.srv.Giftcard.FulfillNonShippableGiftcards(createdNewOrder, orderLines, siteSettings, user, nil, manager)
		if insufficientStockErr != nil || appErr != nil {
			return nil, nil, appErr
		}
	}

	insufficientStock, appErr := s.srv.Order.OrderCreated(transaction, *createdNewOrder, user, nil, manager, false)
	if insufficientStock != nil || appErr != nil {
		return nil, insufficientStock, appErr
	}

	// Send the order confirmation email
	var redirectURL string
	if !checkout.RedirectURL.IsNil() {
		redirectURL = *checkout.RedirectURL.String
	}
	appErr = s.srv.Order.SendOrderConfirmation(createdNewOrder, redirectURL, manager)
	if appErr != nil {
		return nil, nil, appErr
	}

	return createdNewOrder, nil, nil
}

// prepareCheckout Prepare checkout object to complete the checkout process.
func (s *ServiceCheckout) prepareCheckout(transaction boil.ContextTransactor, manager interfaces.PluginManagerInterface, checkoutInfo model_helper.CheckoutInfo, lines model_helper.CheckoutLineInfos, discoutns []*model_helper.DiscountInfo, trackingCode string, redirectURL string, payment *model.Payment) (*model_helper.PaymentError, *model_helper.AppError) {
	checkout := checkoutInfo.Checkout

	appErr := s.CleanCheckoutShipping(checkoutInfo, lines)
	if appErr != nil {
		return nil, appErr
	}

	paymentErr, appErr := s.CleanCheckoutPayment(transaction, manager, checkoutInfo, lines, discoutns, payment)
	if paymentErr != nil || appErr != nil {
		return paymentErr, appErr
	}

	if !checkoutInfo.Channel.IsActive {
		return nil, model_helper.NewAppError("prepareCheckout", "app.checkout.channel_inactive.app_error", nil, "", http.StatusNotAcceptable)
	}
	if redirectURL != "" {
		appErr = model_helper.ValidateStoreFrontUrl(s.srv.Config(), redirectURL)
		if appErr != nil {
			return nil, appErr
		}
	}

	var needUpdate bool
	if redirectURL != "" && (checkout.RedirectURL.IsNil() || *checkout.RedirectURL.String != redirectURL) {
		checkout.RedirectURL = model_types.NewNullString(redirectURL)
		needUpdate = true
	}
	if trackingCode != "" && (checkout.TrackingCode.IsNil() || *checkout.TrackingCode.String != trackingCode) {
		checkout.TrackingCode = model_types.NewNullString(trackingCode)
		needUpdate = true
	}

	if needUpdate {
		_, appErr = s.UpsertCheckouts(transaction, model.CheckoutSlice{&checkout})
		if appErr != nil {
			return nil, appErr
		}
	}

	return nil, nil
}

// ReleaseVoucherUsage
func (s *ServiceCheckout) ReleaseVoucherUsage(orderData map[string]any) *model_helper.AppError {
	if iface, ok := orderData["voucher"]; ok && iface != nil {
		var voucher model.Voucher

		switch t := iface.(type) {
		case *model.Voucher:
			if t == nil {
				return nil
			}
			voucher = *t
		case model.Voucher:
			voucher = t
		}

		if !voucher.UsageLimit.IsNil() && *voucher.UsageLimit.Int != 0 {
			savedVoucher, appErr := s.srv.Discount.AlterVoucherUsage(voucher, -1)
			if appErr != nil {
				return appErr
			}

			if userEmail, ok := orderData["user_email"]; ok {
				appErr = s.srv.Discount.RemoveVoucherUsageByCustomer(*savedVoucher, userEmail.(string))
				if appErr != nil {
					return appErr
				}
			}
		}
	}

	return nil
}

func (s *ServiceCheckout) getOrderData(manager interfaces.PluginManagerInterface, checkoutInfo model_helper.CheckoutInfo, lines model_helper.CheckoutLineInfos, discoutns []*model_helper.DiscountInfo) (map[string]any, *model_helper.AppError) {
	orderData, insufficientStockErr, notApplicableErr, taxError, appErr := s.prepareOrderData(manager, checkoutInfo, lines, discoutns)
	if appErr != nil {
		return nil, appErr
	}

	if insufficientStockErr != nil {
		return nil, s.PrepareInsufficientStockCheckoutValidationAppError("getOrderData", *insufficientStockErr)
	}
	if notApplicableErr != nil {
		return nil, model_helper.NewAppError("getOrderData", "app.checkout.voucher_not_applicable.app_error", map[string]any{"code": model_helper.VOUCHER_NOT_APPLICABLE}, notApplicableErr.Error(), 0)
	}
	if taxError != nil {
		return nil, model_helper.NewAppError("getOrderData", "app.checkout.unable_to_calculate_taxes", map[string]any{"code": model_helper.TAX_ERROR}, taxError.Message, 0)
	}
	return orderData, nil
}

// processPayment Process the payment assigned to checkout
func (s *ServiceCheckout) processPayment(dbTransaction boil.ContextTransactor, payment model.Payment, customerID *string, storeSource bool, paymentData map[string]any, orderData map[string]any, manager interfaces.PluginManagerInterface, channelSlug string) (*model.PaymentTransaction, *model_helper.PaymentError, *model_helper.AppError) {
	var (
		transaction *model.PaymentTransaction
		paymentErr  *model_helper.PaymentError
		appErr      *model_helper.AppError
		paymentID   = payment.ID
	)

	if payment.ToConfirm {
		transaction, paymentErr, appErr = s.srv.Payment.Confirm(
			dbTransaction,
			payment,
			manager,
			channelSlug,
			paymentData,
		)
	} else {
		transaction, paymentErr, appErr = s.srv.Payment.ProcessPayment(
			dbTransaction,
			payment,
			payment.Token,
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
		return nil, nil, model_helper.NewAppError("processPayment", "app.checkout.payment_error.app_error", nil, paymentErr.Error(), 0)
	}

	_, appErr = s.srv.Payment.PaymentByID(nil, paymentID, false)
	if appErr != nil {
		return nil, nil, appErr
	}

	if !transaction.IsSuccess {
		var paymentErrorMessage string
		if !transaction.Error.IsNil() {
			paymentErrorMessage = *transaction.Error.String
		}
		return nil, &model_helper.PaymentError{
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
	dbTransaction boil.ContextTransactor,
	manager interfaces.PluginManagerInterface,
	checkoutInfo model_helper.CheckoutInfo,
	lines model_helper.CheckoutLineInfos,
	paymentData map[string]any,
	storeSource bool,
	discounts []*model_helper.DiscountInfo,
	user *model.User,
	_ any, // this param originally is `app`, but we not gonna integrate app feature in the early versions
	siteSettings model_helper.ShopSettings,
	trackingCode string,
	redirectURL string,
) (*model.Order, bool, model_types.JSONString, *model_helper.PaymentError, *model_helper.AppError) {
	var (
		checkout    = checkoutInfo.Checkout
		channelSlug = checkoutInfo.Channel.Slug
	)

	lastActivePaymentOfCheckout, appErr := s.CheckoutLastActivePayment(checkout) // NOTE: returned payment still can be nil even when appErr is nil
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
		paymentErr, apErr := s.srv.Payment.PaymentRefundOrVoid(dbTransaction, *lastActivePaymentOfCheckout, manager, channelSlug)
		if paymentErr != nil || apErr != nil {
			return nil, false, nil, paymentErr, apErr
		}

		return nil, false, nil, nil, appErr
	}

	var customerID *string
	if lastActivePaymentOfCheckout != nil && user != nil {
		uuid, appErr := s.srv.Payment.FetchCustomerId(*user, lastActivePaymentOfCheckout.Gateway)
		if appErr != nil {
			return nil, false, nil, nil, appErr
		}
		if model_helper.IsValidId(uuid) {
			customerID = &uuid
		}
	}

	transaction, paymentErr, appErr := s.processPayment(
		dbTransaction,
		*lastActivePaymentOfCheckout,
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

	if !transaction.CustomerID.IsNil() && user != nil {
		appErr = s.srv.Payment.StoreCustomerId(user.ID, lastActivePaymentOfCheckout.Gateway, *transaction.CustomerID.String)
		if appErr != nil {
			return nil, false, nil, nil, appErr
		}
	}

	actionData := transaction.ActionRequiredData
	if !transaction.ActionRequired {
		actionData = make(model_types.JSONString)
	}

	var (
		orDer                *model.Order
		insufficientStockErr *model_helper.InsufficientStock
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

			paymentErr, appErr = s.srv.Payment.PaymentRefundOrVoid(dbTransaction, *lastActivePaymentOfCheckout, manager, channelSlug)
			if appErr != nil || paymentErr != nil {
				return nil, false, nil, paymentErr, appErr
			}

			return nil, false, nil, nil, s.PrepareInsufficientStockCheckoutValidationAppError("CompleteCheckout", *insufficientStockErr)
		}

		// if not appError nor insufficient stock error, remove checkout after order is successfully created:
		appErr = s.DeleteCheckoutsByOption(dbTransaction, model_helper.CheckoutFilterOptions{
			// Conditions: squirrel.Eq{model.CheckoutTableName + ".Token": checkout.Token},
			CommonQueryOptions: model_helper.NewCommonQueryOptions(
				model.CheckoutWhere.Token.EQ(checkout.Token),
			),
		})
		if appErr != nil {
			return nil, false, nil, nil, appErr
		}
	}

	return orDer, transaction.ActionRequired, actionData, nil, nil
}

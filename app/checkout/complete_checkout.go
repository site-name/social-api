package checkout

import (
	"net/http"
	"sync"
	"time"

	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/exception"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/checkout"
	"github.com/sitename/sitename/model/giftcard"
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/model/payment"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/modules/json"
	"github.com/sitename/sitename/modules/util"
)

// getVoucherDataForOrder Fetch, process and return voucher/discount data from checkout.
// Careful! It should be called inside a transaction.
// :raises NotApplicable: When the voucher is not applicable in the current checkout.
func (s *ServiceCheckout) getVoucherDataForOrder(checkoutInfo *checkout.CheckoutInfo) (map[string]*product_and_discount.Voucher, *product_and_discount.NotApplicable, *model.AppError) {
	checkOut := checkoutInfo.Checkout
	voucher, appErr := s.GetVoucherForCheckout(checkoutInfo, true)
	if appErr != nil {
		return nil, nil, appErr
	}

	if checkOut.VoucherCode != nil && voucher == nil {
		return nil, product_and_discount.NewNotApplicable("getVoucherDataForOrder", "Voucher expired in meantime. Order placement aborted", nil, 0), nil
	}

	if voucher == nil {
		return map[string]*product_and_discount.Voucher{}, nil, nil
	}

	appErr = s.srv.DiscountService().IncreaseVoucherUsage(voucher)
	if appErr != nil {
		return nil, nil, appErr
	}

	if voucher.ApplyOncePerCustomer {
		notApplicable, appErr := s.srv.DiscountService().AddVoucherUsageByCustomer(voucher, checkoutInfo.GetCustomerEmail())
		if notApplicable != nil || appErr != nil {
			return nil, notApplicable, appErr
		}
	}

	return map[string]*product_and_discount.Voucher{"voucher": voucher}, nil, nil
}

// processShippingDataForOrder Fetch, process and return shipping data from checkout.
func (s *ServiceCheckout) processShippingDataForOrder(checkoutInfo *checkout.CheckoutInfo, shippingPrice *goprices.TaxedMoney, manager interface{}, lines []*checkout.CheckoutLineInfo) (map[string]interface{}, *model.AppError) {
	var (
		deliveryMethodInfo  = checkoutInfo.DeliveryMethodInfo
		shippingAddress     = deliveryMethodInfo.GetShippingAddress()
		copyShippingAddress *account.Address
		appErr              *model.AppError
	)

	if checkoutInfo.User != nil && shippingAddress != nil {
		appErr = s.srv.AccountService().StoreUserAddress(checkoutInfo.User, shippingAddress, account.ADDRESS_TYPE_SHIPPING, manager)
		if appErr != nil {
			return nil, appErr
		}

		anAddressOfUser, appErr := s.srv.AccountService().AddressesByOption(&account.AddressFilterOption{
			Id: &model.StringFilter{
				StringOption: &model.StringOption{
					Eq: shippingAddress.Id,
				},
			},
			UserID: &model.StringFilter{
				StringOption: &model.StringOption{
					Eq: checkoutInfo.User.Id,
				},
			},
		})
		if appErr != nil {
			if appErr.StatusCode != http.StatusNotFound {
				return nil, appErr
			}
		}

		if len(anAddressOfUser) > 0 {
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
func (s *ServiceCheckout) processUserDataForOrder(checkoutInfo *checkout.CheckoutInfo, manager interface{}) (map[string]interface{}, *model.AppError) {
	var (
		billingAddress     = checkoutInfo.BillingAddress
		copyBillingAddress *account.Address
		appErr             *model.AppError
	)

	if checkoutInfo.User != nil && billingAddress != nil {
		appErr = s.srv.AccountService().StoreUserAddress(checkoutInfo.User, billingAddress, account.ADDRESS_TYPE_BILLING, manager)
		if appErr != nil {
			return nil, appErr
		}

		billingAddressOfUser, appErr := s.srv.AccountService().AddressesByOption(&account.AddressFilterOption{
			UserID: &model.StringFilter{
				StringOption: &model.StringOption{
					Eq: checkoutInfo.User.Id,
				},
			},
			Id: &model.StringFilter{
				StringOption: &model.StringOption{
					Eq: billingAddress.Id,
				},
			},
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
func (s *ServiceCheckout) validateGiftcards(checkOut *checkout.Checkout) (*product_and_discount.NotApplicable, *model.AppError) {

	var (
		totalGiftcardsOfCheckout       int
		totalActiveGiftcardsOfCheckout int
		startOfToday                   = util.StartOfDay(time.Now().UTC())
	)

	allGiftcards, appErr := s.srv.GiftcardService().GiftcardsByOption(nil, &giftcard.GiftCardFilterOption{
		CheckoutToken: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: checkOut.Token,
			},
		},
		Distinct: true,
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr
		}
		// ignore not found error
	}

	if allGiftcards != nil {
		totalGiftcardsOfCheckout = len(allGiftcards)
	}

	// find active giftcards
	// NOTE: active giftcards are active and has (ExpiryDate == NULL || ExpiryDate >= beginning of Today)
	var expiryDateOfGiftcard *time.Time
	for _, giftcard := range allGiftcards {
		expiryDateOfGiftcard = giftcard.ExpiryDate
		if (expiryDateOfGiftcard == nil || util.StartOfDay(*expiryDateOfGiftcard).Equal(startOfToday) || util.StartOfDay(*expiryDateOfGiftcard).After(startOfToday)) && *giftcard.IsActive {
			totalActiveGiftcardsOfCheckout++
		}
	}

	if totalActiveGiftcardsOfCheckout != totalGiftcardsOfCheckout {
		return product_and_discount.NewNotApplicable("validateGiftcards", "Gift card has expired. Order placement cancelled.", nil, 0), nil
	}

	return nil, nil
}

// createLineForOrder Create a line for the given order.
// :raises InsufficientStock: when there is not enough items in stock for this variant.
func (s *ServiceCheckout) createLineForOrder(
	manager interface{},
	checkoutInfo *checkout.CheckoutInfo,
	lines []*checkout.CheckoutLineInfo,
	checkoutLineInfo *checkout.CheckoutLineInfo,
	discounts []*product_and_discount.DiscountInfo,
	productsTranslation map[string]string,
	variantsTranslation map[string]string,

) (*order.OrderLineData, *model.AppError) {

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

	// TODO: fixme. This part requires a few works with `manager`
	panic("not implemented")

	productVariantRequireShipping, appErr := s.srv.ProductService().ProductsRequireShipping([]string{variant.ProductID})
	if appErr != nil {
		return nil, appErr
	}

	orderLine := order.OrderLine{
		ProductName:           productName,
		VariantName:           variantName,
		TranslatedProductName: translatedProductName,
		TranslatedVariantName: translatedVariantName,
		ProductSku:            variant.Sku,
		IsShippingRequired:    productVariantRequireShipping,
		Quantity:              quantity,
		VariantID:             &variant.Id,
		UnitPrice:             nil, // TODO: add me
		TotalPrice:            nil, // TODO: add me
		TaxRate:               nil, // TODO: add me
	}

	return &order.OrderLineData{
		Line:        orderLine,
		Quantity:    quantity,
		Variant:     &variant,
		WarehouseID: model.NewString(checkoutInfo.DeliveryMethodInfo.WarehousePK()),
	}, nil
}

// createLinesForOrder Create a lines for the given order.
// :raises InsufficientStock: when there is not enough items in stock for this variant.
func (s *ServiceCheckout) createLinesForOrder(manager interface{}, checkoutInfo *checkout.CheckoutInfo, lines checkout.CheckoutLineInfos, discounts []*product_and_discount.DiscountInfo) ([]*order.OrderLineData, *exception.InsufficientStock, *model.AppError) {
	var (
		translationLanguageCode = checkoutInfo.Checkout.LanguageCode
		countryCode             = checkoutInfo.GetCountry()
		variants                product_and_discount.ProductVariants
		quantities              []int
		products                product_and_discount.Products
		wg                      sync.WaitGroup
		mutex                   sync.Mutex
	)

	lines = lines.FilterNils()

	for _, lineInfo := range lines {
		variants = append(variants, &lineInfo.Variant)
		quantities = append(quantities, lineInfo.Line.Quantity)
		products = append(products, &lineInfo.Product)
	}

	productTranslations, appErr := s.srv.ProductService().ProductTranslationsByOption(&product_and_discount.ProductTranslationFilterOption{
		ProductID: &model.StringFilter{
			StringOption: &model.StringOption{
				In: products.IDs(),
			},
		},
		LanguageCode: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: translationLanguageCode,
			},
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

	variantTranslations, appErr := s.srv.ProductService().ProductVariantTranslationsByOption(&product_and_discount.ProductVariantTranslationFilterOption{
		ProductVariantID: &model.StringFilter{
			StringOption: &model.StringOption{
				In: variants.IDs(),
			},
		},
		LanguageCode: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: translationLanguageCode,
			},
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
	insufficientStockErr, appErr := s.srv.WarehouseService().CheckStockQuantityBulk(
		variants,
		countryCode,
		quantities,
		checkoutInfo.Channel.Slug,
		additionalWarehouseLookup,
		nil,
	)
	if insufficientStockErr != nil || appErr != nil {
		return nil, insufficientStockErr, appErr
	}

	var (
		orderLineDatas []*order.OrderLineData
		appError       *model.AppError
	)

	for _, item := range lines {
		wg.Add(1)

		go func(lineInfo *checkout.CheckoutLineInfo) {
			mutex.Lock()
			defer mutex.Unlock()

			orderLineData, appErr := s.createLineForOrder(manager, checkoutInfo, lines, lineInfo, discounts, productTranslationMap, productVariantTranslationMap)
			if appErr != nil {
				if appErr.StatusCode == http.StatusInternalServerError && appError == nil {
					appError = appErr
				}
			} else if orderLineData != nil {
				orderLineDatas = append(orderLineDatas, orderLineData)
			}

			wg.Done()
		}(item)
	}

	if len(lines) > 0 {
		wg.Wait()
	}

	if appError != nil {
		return nil, nil, appErr
	}

	return orderLineDatas, nil, nil
}

// prepareOrderData Run checks and return all the data from a given checkout to create an order.
// :raises NotApplicable InsufficientStock:
func (s *ServiceCheckout) prepareOrderData(manager interface{}, checkoutInfo *checkout.CheckoutInfo, lines checkout.CheckoutLineInfos, discounts []*product_and_discount.DiscountInfo) (map[string]interface{}, *exception.InsufficientStock, *product_and_discount.NotApplicable, *exception.TaxError, *model.AppError) {
	checkOut := checkoutInfo.Checkout
	checkOut.PopulateNonDbFields() // this call is important

	// orderData = OrderData{}
	address := checkoutInfo.ShippingAddress
	if address == nil {
		address = checkoutInfo.BillingAddress
	}

	taxedTotal, appErr := s.CheckoutTotal(manager, checkoutInfo, lines, address, discounts)
	if appErr != nil {
		return nil, nil, nil, nil, appErr
	}

	cardsTotal, appErr := s.CheckoutTotalGiftCardsBalance(&checkOut)
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
		return nil, nil, nil, nil, model.NewAppError("prepareOrderData", app.ErrorCalculatingMoneyErrorID, nil, errMsg, http.StatusInternalServerError)
	}

	taxedTotal.Gross = newTaxedTotalGross
	taxedTotal.Net = newTaxedTotalNet

	zeroTaxedMoney, _ := util.ZeroTaxedMoney(checkOut.Currency)
	if less, err := taxedTotal.LessThan(zeroTaxedMoney); less && err == nil {
		taxedTotal = zeroTaxedMoney
	}

	// TODO: implement a few works with plugin manager here.
	panic("not implemented")

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
func (s *ServiceCheckout) createOrder(checkoutInfo *checkout.CheckoutInfo, orderData map[string]interface{}, user *account.User, _ interface{}, manager interface{}, siteSettings interface{}) (*order.Order, *exception.InsufficientStock, *model.AppError) {
	// create transaction
	transaction, err := s.srv.Store.GetMaster().Begin()
	if err != nil {
		return nil, nil, model.NewAppError("createOrder", app.ErrorCreatingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}
	defer s.srv.Store.FinalizeTransaction(transaction)

	checkOut := checkoutInfo.Checkout
	// checkOut.PopulateNonDbFields() // this call is important

	orders, appErr := s.srv.OrderService().FilterOrdersByOptions(&order.OrderFilterOption{
		CheckoutToken: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: checkOut.Token,
			},
		},
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, nil, appErr
		}
	}

	if len(orders) > 0 {
		return orders[0], nil, nil
	}

	totalPriceLeft := orderData["total_price_left"].(*goprices.Money)
	delete(orderData, "total_price_left")
	orderLinesInfo := orderData["lines"].([]*order.OrderLineData)
	delete(orderData, "lines")

	status := order.UNCONFIRMED
	if siteSettings == nil {
		shop, appErr := s.srv.ShopService().ShopById(checkOut.ShopID)
		if appErr != nil {
			return nil, nil, appErr
		}
		if *shop.AutomaticallyConfirmAllNewOrders {
			status = order.UNFULFILLED
		}
	}

	/*
		NOTE: we can easily convert a map[string]interface{} to Order{} since:
		the map's keys are exactly match json tags of order struct's fields
	*/
	serializedOrderData, err := json.JSON.Marshal(orderData)
	if err != nil {
		return nil, nil, model.NewAppError("createOrder", app.ErrorMarshallingDataID, nil, err.Error(), http.StatusInternalServerError)
	}

	// define new order to create
	var newOrder order.Order
	err = json.JSON.Unmarshal(serializedOrderData, &newOrder)
	if err != nil {
		return nil, nil, model.NewAppError("createOrder", app.ErrorUnMarshallingDataID, nil, err.Error(), http.StatusInternalServerError)
	}

	newOrder.Id = ""
	newOrder.CheckoutToken = checkOut.Token
	newOrder.Status = status
	newOrder.Origin = order.CHECKOUT
	newOrder.ChannelID = checkoutInfo.Channel.Id

	createdNewOrder, appErr := s.srv.OrderService().UpsertOrder(transaction, &newOrder)
	if appErr != nil {
		return nil, nil, appErr
	}

	if checkOut.DiscountAmount != nil {
		// store voucher as a fixed value as it this the simplest solution for now.
		// This will be solved when we refactor the voucher logic to use .discounts
		// relations
		_, appErr := s.srv.DiscountService().UpsertOrderDiscount(transaction, &product_and_discount.OrderDiscount{
			Type:           product_and_discount.VOUCHER,
			ValueType:      product_and_discount.FIXED,
			Value:          checkOut.DiscountAmount,
			Name:           checkOut.DiscountName,
			TranslatedName: checkOut.TranslatedDiscountName,
			Currency:       checkOut.Currency,
			AmountValue:    checkOut.DiscountAmount,
			OrderID:        &createdNewOrder.Id,
		})
		if appErr != nil {
			return nil, nil, appErr
		}
	}

	var orderLines []*order.OrderLine
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

	appErr = s.srv.OrderService().AddGiftcardsToOrder(transaction, checkoutInfo, createdNewOrder, totalPriceLeft, user, nil)
	if appErr != nil {
		return nil, nil, appErr
	}

	// assign checkout payments to other order
	appErr = s.srv.PaymentService().UpdatePaymentsOfCheckout(transaction, checkOut.Token, &payment.PaymentPatch{OrderID: createdNewOrder.Id})
	if appErr != nil {
		return nil, nil, appErr
	}

	// copy metadata from the checkout into the new order
	createdNewOrder.Metadata = model.CopyStringMap(checkOut.Metadata)
	createdNewOrder.RedirectUrl = checkOut.RedirectURL
	createdNewOrder.PrivateMetadata = model.CopyStringMap(checkOut.PrivateMetadata)

	appErr = s.srv.OrderService().UpdateOrderTotalPaid(transaction, createdNewOrder)
	if appErr != nil {
		return nil, nil, appErr
	}

	_, appErr = s.srv.OrderService().UpsertOrder(transaction, createdNewOrder)
	if appErr != nil {
		return nil, nil, appErr
	}

	// commit transaction
	if err = transaction.Commit(); err != nil {
		return nil, nil, model.NewAppError("createOrder", app.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	appErr = s.srv.OrderService().OrderCreated(createdNewOrder, user, nil, manager, false)
	if appErr != nil {
		return nil, nil, appErr
	}

	// Send the order confirmation email
	// TODO: fixme
	panic("not implemented")

	return createdNewOrder, nil, nil
}

// prepareCheckout Prepare checkout object to complete the checkout process.
func (s *ServiceCheckout) prepareCheckout(manager interface{}, checkoutInfo *checkout.CheckoutInfo, lines []*checkout.CheckoutLineInfo, discoutns []*product_and_discount.DiscountInfo, trackingCode string, redirectURL string, payMent *payment.Payment) (*payment.PaymentError, *model.AppError) {
	checkOut := checkoutInfo.Checkout

	appErr := s.CleanCheckoutShipping(checkoutInfo, lines)
	if appErr != nil {
		return nil, appErr
	}

	paymentErr, appErr := s.CleanCheckoutPayment(manager, checkoutInfo, lines, discoutns, payMent)
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
	if redirectURL != "" && (checkOut.RedirectURL == nil || *checkOut.RedirectURL != redirectURL) {
		checkOut.RedirectURL = &redirectURL
		needUpdate = true
	}
	if trackingCode != "" && (checkOut.TrackingCode == nil || *checkOut.TrackingCode != trackingCode) {
		checkOut.TrackingCode = &trackingCode
		needUpdate = true
	}

	if needUpdate {
		_, appErr = s.UpsertCheckout(&checkOut)
		if appErr != nil {
			return nil, appErr
		}
	}

	return nil, nil
}

// ReleaseVoucherUsage
func (s *ServiceCheckout) ReleaseVoucherUsage(orderData map[string]interface{}) *model.AppError {
	if iface, ok := orderData["voucher"]; ok && iface != nil {
		voucher := iface.(*product_and_discount.Voucher)

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

	return nil
}

func (s *ServiceCheckout) getOrderData(manager interface{}, checkoutInfo *checkout.CheckoutInfo, lines []*checkout.CheckoutLineInfo, discoutns []*product_and_discount.DiscountInfo) (map[string]interface{}, *model.AppError) {
	orderData, insufficientStockErr, notApplicableErr, taxError, appErr := s.prepareOrderData(manager, checkoutInfo, lines, discoutns)
	if appErr != nil {
		return nil, appErr
	}

	if insufficientStockErr != nil {
		return nil, checkout.PrepareInsufficientStockCheckoutValidationAppError("getOrderData", insufficientStockErr)
	}
	if notApplicableErr != nil {
		return nil, model.NewAppError("getOrderData", "app.checkout.voucher_not_applicable.app_error", map[string]interface{}{"code": exception.VOUCHER_NOT_APPLICABLE}, notApplicableErr.Error(), 0)
	}
	if taxError != nil {
		return nil, model.NewAppError("getOrderData", "app.checkout.unable_to_calculate_taxes", map[string]interface{}{"code": exception.TAX_ERROR}, taxError.Message, 0)
	}
	return orderData, nil
}

// processPayment Process the payment assigned to checkout
func (s *ServiceCheckout) processPayment(payMent *payment.Payment, customerID *string, storeSource bool, paymentData map[string]interface{}, orderData map[string]interface{}, manager interface{}, channelSlug string) (*payment.PaymentTransaction, *payment.PaymentError, *model.AppError) {
	var (
		transaction *payment.PaymentTransaction
		paymentErr  *payment.PaymentError
		appErr      *model.AppError
		paymentID   string = payMent.Id
	)

	if payMent.ToConfirm {
		transaction, paymentErr, appErr = s.srv.PaymentService().Confirm(
			payMent,
			manager,
			channelSlug,
			paymentData,
		)
	} else {
		transaction, paymentErr, appErr = s.srv.PaymentService().ProcessPayment(
			payMent,
			payMent.Token,
			manager,
			channelSlug,
			customerID,
			storeSource,
			paymentData,
		)
	}

	// catch errors
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

	// re fetching payment from db since the payment may was modified in two calls above
	payMent, appErr = s.srv.PaymentService().PaymentByID(nil, paymentID, false)
	if appErr != nil {
		return nil, nil, appErr
	}

	if !transaction.IsSuccess {
		var paymentErrorMessage string
		if transaction.Error != nil {
			paymentErrorMessage = *transaction.Error
		}
		return nil, &payment.PaymentError{
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
	manager interface{},
	checkoutInfo *checkout.CheckoutInfo,
	lines []*checkout.CheckoutLineInfo,
	paymentData map[string]interface{},
	storeSource bool,
	discounts []*product_and_discount.DiscountInfo,
	user *account.User, // must be authenticated before this
	_ interface{}, // this param originally is `app`, but we not gonna integrate app feature in the early versions
	siteSettings interface{},
	trackingCode string,
	redirectURL string,

) (*order.Order, bool, model.StringMap, *payment.PaymentError, *model.AppError) {

	var (
		checkOut    = checkoutInfo.Checkout
		channelSlug = checkoutInfo.Channel.Slug
	)

	lastActivePaymentOfCheckout, appErr := s.CheckoutLastActivePayment(&checkOut) // NOTE: returned payment still can be nil even when appErr is nil
	if appErr != nil {
		return nil, false, nil, nil, appErr
	}

	paymentErr, appErr := s.prepareCheckout(
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
		paymentErr, apErr := s.srv.PaymentService().PaymentRefundOrVoid(lastActivePaymentOfCheckout, manager, channelSlug)
		if paymentErr != nil || apErr != nil {
			return nil, false, nil, paymentErr, apErr
		}

		return nil, false, nil, nil, appErr
	}

	var customerID *string
	if lastActivePaymentOfCheckout != nil && user != nil { // NOTE: user must be authenticated before calling this method.
		uuid, appErr := s.srv.PaymentService().FetchCustomerId(user, lastActivePaymentOfCheckout.GateWay)
		if appErr != nil {
			return nil, false, nil, nil, appErr
		}
		if model.IsValidId(uuid) {
			customerID = &uuid
		}
	}

	transaction, paymentErr, appErr := s.processPayment(
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

	if transaction.CustomerID != nil && user != nil && model.IsValidId(user.Id) {
		appErr = s.srv.PaymentService().StoreCustomerId(user.Id, lastActivePaymentOfCheckout.GateWay, *transaction.CustomerID)
		if appErr != nil {
			return nil, false, nil, nil, appErr
		}
	}

	actionData := transaction.ActionRequiredData
	if !transaction.ActionRequired {
		actionData = make(model.StringMap)
	}

	var (
		orDer                *order.Order
		insufficientStockErr *exception.InsufficientStock
	)
	if !transaction.ActionRequired {
		orDer, insufficientStockErr, appErr = s.createOrder(checkoutInfo, orderData, user, nil, manager, siteSettings)
		if appErr != nil {
			return nil, false, nil, nil, appErr
		}

		if insufficientStockErr != nil {
			appErr = s.ReleaseVoucherUsage(orderData)
			if appErr != nil {
				return nil, false, nil, nil, appErr
			}

			paymentErr, appErr = s.srv.PaymentService().PaymentRefundOrVoid(lastActivePaymentOfCheckout, manager, channelSlug)
			if appErr != nil || paymentErr != nil {
				return nil, false, nil, paymentErr, appErr
			}

			return nil, false, nil, nil, checkout.PrepareInsufficientStockCheckoutValidationAppError("", insufficientStockErr)
		}

		// if not appError nor insufficient stock error, remove checkout after order is successfully created:
		appErr = s.DeleteCheckoutsByOption(nil, &checkout.CheckoutFilterOption{
			Token: &model.StringFilter{
				StringOption: &model.StringOption{
					Eq: checkOut.Token,
				},
			},
		})
		if appErr != nil {
			return nil, false, nil, nil, appErr
		}
	}

	return orDer, transaction.ActionRequired, actionData, nil, nil
}
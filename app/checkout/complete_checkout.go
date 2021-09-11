package checkout

import (
	"net/http"
	"sync"
	"time"

	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/checkout"
	"github.com/sitename/sitename/model/giftcard"
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/model/payment"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/model/warehouse"
	"github.com/sitename/sitename/modules/util"
)

// getVoucherDataForOrder Fetch, process and return voucher/discount data from checkout.
// Careful! It should be called inside a transaction.
// :raises NotApplicable: When the voucher is not applicable in the current checkout.
func (s *ServiceCheckout) getVoucherDataForOrder(checkoutInfo *checkout.CheckoutInfo) (interface{}, *model.NotApplicable, *model.AppError) {
	checkOut := checkoutInfo.Checkout
	voucher, appErr := s.GetVoucherForCheckout(checkoutInfo, true)
	if appErr != nil {
		return nil, nil, appErr
	}

	if checkOut.VoucherCode != nil && voucher == nil {
		return nil, model.NewNotApplicable("getVoucherDataForOrder", "Voucher expired in meantime. Order placement aborted", nil, 0), nil
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

	deliveryMethodDict := map[string]interface{}{
		deliveryMethodInfo.GetOrderKey(): deliveryMethodInfo.GetDeliveryMethod(),
	}

	if checkoutInfo.User != nil && shippingAddress != nil {
		appErr = s.srv.AccountService().StoreUserAddress(checkoutInfo.User, shippingAddress, account.ADDRESS_TYPE_SHIPPING, manager)
		if appErr != nil {
			return nil, appErr
		}

		billingAddressOfUser, appErr := s.srv.AccountService().AddressesByOption(&account.AddressFilterOption{
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
		if appErr != nil && appErr.StatusCode != http.StatusNotFound {
			return nil, appErr
		}

		if len(billingAddressOfUser) > 0 {
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

	if copyShippingAddress != nil {
		deliveryMethodDict["shipping_address"] = copyShippingAddress
	} else {
		deliveryMethodDict["shipping_address"] = shippingAddress
	}

	deliveryMethodDict["shipping_price"] = shippingPrice
	deliveryMethodDict["weight"] = checkoutTotalWeight

	return deliveryMethodDict, nil
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
func (s *ServiceCheckout) validateGiftcards(checkOut *checkout.Checkout) (*model.NotApplicable, *model.AppError) {
	startOfToday := util.StartOfDay(time.Now().UTC())

	var (
		TotalGiftcardsOfCheckout       int
		TotalActiveGiftcardsOfCheckout int
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
		TotalGiftcardsOfCheckout = len(allGiftcards)
	}

	// find active giftcards
	// NOTE: active giftcards are active and has (ExpiryDate == NULL || ExpiryDate >= beginning of Today)
	for _, item := range allGiftcards {
		expiryDateOfItem := item.ExpiryDate
		if (expiryDateOfItem == nil || util.StartOfDay(*expiryDateOfItem).Equal(startOfToday) || util.StartOfDay(*expiryDateOfItem).After(startOfToday)) && *item.IsActive {
			TotalActiveGiftcardsOfCheckout++
		}
	}

	if TotalActiveGiftcardsOfCheckout != TotalGiftcardsOfCheckout {
		return model.NewNotApplicable("validateGiftcards", "Gift card has expired. Order placement cancelled.", nil, 0), nil
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
		_                     = checkoutLine.Quantity
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

	panic("not implemented")
}

// createLinesForOrder Create a lines for the given order.
// :raises InsufficientStock: when there is not enough items in stock for this variant.
func (s *ServiceCheckout) createLinesForOrder(manager interface{}, checkoutInfo *checkout.CheckoutInfo, lines checkout.CheckoutLineInfos, discounts []*product_and_discount.DiscountInfo) ([]*order.OrderLineData, *warehouse.InsufficientStock, *model.AppError) {
	var (
		translationLanguageCode = checkoutInfo.Checkout.LanguageCode
		countryCode             = checkoutInfo.GetCountry()
		variants                = product_and_discount.ProductVariants{}
		quantities              = []int{}
		products                = product_and_discount.Products{}
		wg                      sync.WaitGroup
		mutex                   sync.Mutex
	)

	lines = lines.FilterNils()

	for _, lineInfo := range lines {
		if lineInfo.Variant != nil {
			variants = append(variants, lineInfo.Variant)
		}
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
	}

	productTranslationMap := map[string]string{}
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
	}

	productVariantTranslationMap := map[string]string{}
	if len(variantTranslations) > 0 {
		for _, item := range variantTranslations {
			productVariantTranslationMap[item.ProductVariantID] = item.Name
		}
	}

	additionalWarehouseLookup := checkoutInfo.DeliveryMethodInfo.GetWarehouseFilterLookup()
	insufficientStockErr, appErr := s.srv.WarehouseService().CheckStockQuantityBulk(variants, countryCode, quantities, checkoutInfo.Channel.Slug, additionalWarehouseLookup, nil)
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

type OrderData struct {
	order.Order

	TotalPriceLeft *goprices.Money
	Lines          []*order.OrderLineData
}

// prepareOrderData Run checks and return all the data from a given checkout to create an order.
// :raises NotApplicable InsufficientStock:
func (s *ServiceCheckout) prepareOrderData(manager interface{}, checkoutInfo *checkout.CheckoutInfo, lines checkout.CheckoutLineInfos, discounts []*product_and_discount.DiscountInfo) (interface{}, *model.AppError) {
	var (
		checkOut = checkoutInfo.Checkout
		// orderData = OrderData{}
		address = checkoutInfo.ShippingAddress
	)
	checkOut.PopulateNonDbFields() // this call is important

	if address == nil {
		address = checkoutInfo.BillingAddress
	}

	taxedTotal, appErr := s.CheckoutTotal(manager, checkoutInfo, lines, address, discounts)
	if appErr != nil {
		return nil, appErr
	}

	cardsTotal, appErr := s.CheckoutTotalGiftCardsBalance(&checkOut)
	if appErr != nil {
		return nil, appErr
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
		return nil, model.NewAppError("prepareOrderData", app.ErrorCalculatingMoneyErrorID, nil, errMsg, http.StatusInternalServerError)
	}

	taxedTotal.Gross = newTaxedTotalGross
	taxedTotal.Net = newTaxedTotalNet

	zeroTaxedMoney, _ := util.ZeroTaxedMoney(checkOut.Currency)
	if less, err := taxedTotal.LessThan(zeroTaxedMoney); less && err == nil {
		taxedTotal = zeroTaxedMoney
	}

	panic("not implemented")
}

// createOrder Create an order from the checkout.
// Each order will get a private copy of both the billing and the shipping
// address (if shipping).
// If any of the addresses is new and the user is logged in the address
// will also get saved to that user's address book.
// Current user's language is saved in the order so we can later determine
// which language to use when sending email.
func (s *ServiceCheckout) createOrder(checkoutInfo *checkout.CheckoutInfo, orderData OrderData, user *account.User, _ interface{}, manager interface{}, siteSettings interface{}) (*order.Order, *warehouse.InsufficientStock, *model.AppError) {
	// create transaction
	transaction, err := s.srv.Store.GetMaster().Begin()
	if err != nil {
		return nil, nil, model.NewAppError("createOrder", app.ErrorCreatingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}
	defer s.srv.Store.FinalizeTransaction(transaction)

	checkOut := checkoutInfo.Checkout
	checkOut.PopulateNonDbFields() // this call is important

	orders, appErr := s.srv.OrderService().FilterOrdersByOptions(&order.OrderFilterOption{
		CheckoutToken: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: checkOut.Token,
			},
		},
	})
	if appErr != nil {
		return nil, nil, appErr
	}

	if len(orders) > 0 {
		return orders[0], nil, nil
	}

	totalPriceLeft := orderData.TotalPriceLeft
	orderLinesInfo := orderData.Lines

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

	newOrder := orderData.Order
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

	// add giftcards to other order
	giftcardsOfCheckout, appErr := s.srv.GiftcardService().GiftcardsByOption(transaction, &giftcard.GiftCardFilterOption{
		CheckoutToken: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: checkOut.Token,
			},
		},
		SelectForUpdate: true,
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, nil, appErr
		}
		giftcardsOfCheckout = make([]*giftcard.GiftCard, 0)
	}

	var newTotalPriceLeft *goprices.Money
	for _, giftcard := range giftcardsOfCheckout {
		newTotalPriceLeft, appErr = s.srv.OrderService().AddGiftCardToOrder(createdNewOrder, giftcard, totalPriceLeft)
		if appErr != nil {
			return nil, nil, appErr
		}
		totalPriceLeft = newTotalPriceLeft
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

	if err = transaction.Commit(); err != nil {
		return nil, nil, model.NewAppError("createOrder", app.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	appErr = s.srv.OrderService().OrderCreated(createdNewOrder, user, manager, false)
	if appErr != nil {
		return nil, nil, appErr
	}

	// Send the order confirmation email
	// TODO: fixme
	panic("not implemented")

	return createdNewOrder, nil, nil
}

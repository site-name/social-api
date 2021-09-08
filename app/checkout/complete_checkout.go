package checkout

import (
	"net/http"
	"time"

	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/checkout"
	"github.com/sitename/sitename/model/giftcard"
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/model/product_and_discount"
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
		shippingAddress     = deliveryMethodInfo.ShippingAddress
		copyShippingAddress *account.Address
		appErr              *model.AppError
	)

	deliveryMethodDict := map[string]interface{}{
		deliveryMethodInfo.OrderKey: deliveryMethodInfo.DeliveryMethod,
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

	allGiftcards, appErr := s.srv.GiftcardService().GiftcardsByOption(&giftcard.GiftCardFilterOption{
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
	productsTranslation map[string]*string,
	variantsTranslation map[string]*string,

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
	if translatedProductName != nil && *translatedProductName == productName {
		translatedProductName = model.NewString("")
	}
	if translatedVariantName != nil && *translatedVariantName == variantName {
		translatedVariantName = model.NewString("")
	}

	panic("not implemented")
}

// createLinesForOrder Create a lines for the given order.
// :raises InsufficientStock: when there is not enough items in stock for this variant.
func (s *ServiceCheckout) createLinesForOrder(manager interface{}, checkoutInfo *checkout.CheckoutInfo, lines []*checkout.CheckoutLineInfo, discounts []*product_and_discount.DiscountInfo) ([]*order.OrderLineData, *model.AppError) {
	var (
		translationLanguageCode = checkoutInfo.Checkout.LanguageCode
		countryCode             = checkoutInfo.GetCountry()
		variants                = []*product_and_discount.ProductVariant{}
		quantities              = []int{}
		products                = product_and_discount.Products{}
	)

	for _, lineInfo := range lines {
		if lineInfo.Variant != nil {
			variants = append(variants, lineInfo.Variant)
		}
		quantities = append(quantities, lineInfo.Line.Quantity)
		products = append(products, &lineInfo.Product)
	}

}

package checkout

import (
	"net/http"

	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/checkout"
	"github.com/sitename/sitename/model/product_and_discount"
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

}

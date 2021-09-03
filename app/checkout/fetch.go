package checkout

import (
	"context"
	"net/http"

	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/channel"
	"github.com/sitename/sitename/model/checkout"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/model/shipping"
)

// FetchCheckoutLines Fetch checkout lines as CheckoutLineInfo objects.
// It prefetch some related value also
func (a *ServiceCheckout) FetchCheckoutLines(checkOut *checkout.Checkout) ([]*checkout.CheckoutLineInfo, *model.AppError) {
	checkoutLineInfos, err := a.srv.Store.Checkout().FetchCheckoutLinesAndPrefetchRelatedValue(checkOut)
	if err != nil {
		return nil, model.NewAppError("FetchCheckoutLines", "app.checkout.error_collecting_checkout_line_infos.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return checkoutLineInfos, nil
}

// FetCheckoutInfo Fetch checkout as CheckoutInfo object
func (a *ServiceCheckout) FetCheckoutInfo(checkOut *checkout.Checkout, lines []*checkout.CheckoutLineInfo, discounts []*product_and_discount.DiscountInfo, manager interface{}) (*checkout.CheckoutInfo, *model.AppError) {
	// validate arguments:
	if checkOut == nil {
		return nil, model.NewAppError("FetchCheckoutInfo", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "checkOut"}, "", http.StatusBadRequest)
	}

	var (
		appError        *model.AppError
		chanNel         = make(chan *channel.Channel)
		shippingAddress = make(chan *account.Address)
		shippingMethod  = make(chan *shipping.ShippingMethod)
		checkoutUser    = make(chan *account.User)
		billingAddress  = make(chan *account.Address)
	)

	syncSetAppError := func(err *model.AppError) {
		a.mutex.Lock()
		defer a.mutex.Unlock()

		if err != nil && appError == nil {
			appError = err
		}
	}

	defer func() {
		close(chanNel)
		close(shippingAddress)
		close(shippingMethod)
		close(checkoutUser)
		close(billingAddress)
	}()

	go func() {
		chNel, appErr := a.srv.ChannelService().ChannelByOption(&channel.ChannelFilterOption{
			Id: &model.StringFilter{
				StringOption: &model.StringOption{
					Eq: checkOut.ChannelID,
				},
			},
		})
		if appErr != nil {
			syncSetAppError(appErr)
		}
		chanNel <- chNel
	}()

	go func() {
		if checkOut.ShippingAddressID != nil || checkOut.BillingAddressID != nil {

			addressIDs := []string{}
			if checkOut.ShippingAddressID != nil {
				addressIDs = append(addressIDs, *checkOut.ShippingAddressID)
			}
			if checkOut.BillingAddressID != nil {
				addressIDs = append(addressIDs, *checkOut.BillingAddressID)
			}

			addresses, appErr := a.srv.AccountService().AddressesByOption(&account.AddressFilterOption{
				Id: &model.StringFilter{
					StringOption: &model.StringOption{
						In: addressIDs,
					},
				},
			})
			if appErr != nil {
				if appErr.StatusCode == http.StatusInternalServerError {
					syncSetAppError(appErr)
				}

				// ignore other errors
				shippingAddress <- nil
				billingAddress <- nil

			} else {
				switch len(addresses) {
				case 1:
					if checkOut.ShippingAddressID != nil {
						shippingAddress <- addresses[0]
						billingAddress <- nil
					} else if checkOut.BillingAddressID != nil {
						billingAddress <- addresses[0]
						shippingAddress <- nil
					}
				case 2:
					if addresses[0].Id == *checkOut.ShippingAddressID {
						shippingAddress <- addresses[0]
						billingAddress <- addresses[1]
					} else {
						shippingAddress <- addresses[1]
						billingAddress <- addresses[0]
					}
				}
			}

		} else {
			shippingAddress <- nil
			billingAddress <- nil
		}
	}()

	go func() {
		if checkOut.ShippingMethodID != nil {
			shippingMt, appErr := a.srv.ShippingService().ShippingMethodByOption(&shipping.ShippingMethodFilterOption{
				Id: &model.StringFilter{
					StringOption: &model.StringOption{
						Eq: *checkOut.ShippingMethodID,
					},
				},
			})
			if appErr != nil {
				syncSetAppError(appErr)
			}
			shippingMethod <- shippingMt

		} else {
			shippingMethod <- nil
		}
	}()

	go func() {
		if checkOut.UserID != nil {
			user, appErr := a.srv.AccountService().UserById(context.Background(), *checkOut.UserID)
			if appErr != nil {
				if appErr.StatusCode == http.StatusInternalServerError {
					syncSetAppError(appErr)
				}
				// ignore non-system-related errors
				checkoutUser <- user
			} else {
				checkoutUser <- user
			}

		} else {
			checkoutUser <- nil
		}
	}()

	resultChannel := <-chanNel
	resultShippingMethod := <-shippingMethod
	resultShippingAddress := <-shippingAddress
	resultBillingAddress := <-billingAddress
	resultUser := <-checkoutUser

	if appError != nil {
		return nil, appError
	}

	// declare filter option for shipping method channel listings
	shippingMethodChannelListingFilterOption := new(shipping.ShippingMethodChannelListingFilterOption)

	// add filter options
	if resultShippingMethod != nil {
		shippingMethodChannelListingFilterOption.ShippingMethodID = &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: resultShippingMethod.Id,
			},
		}
	}
	if resultChannel != nil {
		shippingMethodChannelListingFilterOption.ChannelID = &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: resultChannel.Id,
			},
		}
	}

	// find the first shipping method channel listing
	var shippingMethodChannelListing *shipping.ShippingMethodChannelListing
	shippingMethodChannelListings, appErr := a.srv.ShippingService().ShippingMethodChannelListingsByOption(shippingMethodChannelListingFilterOption)
	if appErr != nil && appErr.StatusCode == http.StatusInternalServerError {
		return nil, appErr
	}
	if shippingMethodChannelListings != nil && len(shippingMethodChannelListings) > 0 {
		shippingMethodChannelListing = shippingMethodChannelListings[0]
	}

	checkoutInfo := &checkout.CheckoutInfo{
		Checkout:                      *checkOut,
		User:                          resultUser,
		Channel:                       *resultChannel,
		BillingAddress:                resultBillingAddress,
		ShippingAddress:               resultShippingAddress,
		ShippingMethod:                resultShippingMethod,
		ShippingMethodChannelListings: shippingMethodChannelListing,
		ValidShippingMethods:          []*shipping.ShippingMethod{}, // empty
	}

	validShippingMethods, appErr := a.GetValidShippingMethodListForCheckoutInfo(
		checkoutInfo,
		resultShippingAddress,
		lines,
		discounts,
		manager,
	)
	if appErr != nil {
		return nil, appErr
	}
	checkoutInfo.ValidShippingMethods = validShippingMethods

	return checkoutInfo, nil
}

// UpdateCheckoutInfoShippingAddress updates given `checkoutInfo` by setting given `address` as its ShippingAddress.
// then updates its ValidShippingMethods
func (a *ServiceCheckout) UpdateCheckoutInfoShippingAddress(checkoutInfo *checkout.CheckoutInfo, address *account.Address, lines []*checkout.CheckoutLineInfo, discounts []*product_and_discount.DiscountInfo, manager interface{}) *model.AppError {
	checkoutInfo.ShippingAddress = address
	validMethods, appErr := a.GetValidShippingMethodListForCheckoutInfo(checkoutInfo, address, lines, discounts, manager)
	if appErr != nil {
		return appErr
	}

	checkoutInfo.ValidShippingMethods = validMethods
	return nil
}

func (a *ServiceCheckout) GetValidShippingMethodListForCheckoutInfo(checkoutInfo *checkout.CheckoutInfo, shippingAddress *account.Address, lines []*checkout.CheckoutLineInfo, discounts []*product_and_discount.DiscountInfo, manager interface{}) ([]*shipping.ShippingMethod, *model.AppError) {
	panic("not implt")
}

// UpdateCheckoutInfoShippingMethod set CheckoutInfo's ShippingMethod to given shippingMethod
// and set new value for checkoutInfo's ShippingMethodChannelListings
func (a *ServiceCheckout) UpdateCheckoutInfoShippingMethod(checkoutInfo *checkout.CheckoutInfo, shippingMethod *shipping.ShippingMethod) *model.AppError {
	checkoutInfo.ShippingMethod = shippingMethod

	checkoutInfo.ShippingMethodChannelListings = nil
	if shippingMethod != nil {
		listings, appErr := a.srv.ShippingService().ShippingMethodChannelListingsByOption(&shipping.ShippingMethodChannelListingFilterOption{
			ShippingMethodID: &model.StringFilter{
				StringOption: &model.StringOption{
					Eq: shippingMethod.Id,
				},
			},
			ChannelID: &model.StringFilter{
				StringOption: &model.StringOption{
					Eq: checkoutInfo.Channel.Id,
				},
			},
		})

		if appErr != nil {
			return appErr
		}
		checkoutInfo.ShippingMethodChannelListings = listings[0]
	}

	return nil
}

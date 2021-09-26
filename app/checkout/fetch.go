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
	"github.com/sitename/sitename/model/warehouse"
)

// GetDeliveryMethodInfo takes `deliveryMethod` is either *ShippingMethod or *Warehouse
func (s *ServiceCheckout) GetDeliveryMethodInfo(deliveryMethod interface{}, address *account.Address) (checkout.DeliveryMethodBaseInterface, *model.AppError) {
	if deliveryMethod == nil {
		return &checkout.DeliveryMethodBase{}, nil
	}

	switch t := deliveryMethod.(type) {
	case *shipping.ShippingMethod:
		return &checkout.ShippingMethodInfo{
			DeliveryMethod:  *t,
			ShippingAddress: address,
		}, nil

	case *warehouse.WareHouse:
		var (
			addr   = t.Address
			appErr *model.AppError
		)
		if addr == nil && t.AddressID != nil {
			addr, appErr = s.srv.AccountService().AddressById(*t.AddressID)
			if appErr != nil {
				return nil, appErr
			}
		}
		return &checkout.CollectionPointInfo{
			DeliveryMethod:  *t,
			ShippingAddress: addr,
		}, nil

	default:
		return nil, model.NewAppError("GetDeliveryMethodInfo", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "deliveryMethod"}, "", http.StatusBadRequest)
	}
}

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

	chanNel, appErr := a.srv.ChannelService().ChannelByOption(&channel.ChannelFilterOption{
		Id: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: checkOut.ChannelID,
			},
		},
	})
	if appErr != nil {
		return nil, appErr
	}

	var shippingAddress *account.Address
	if checkOut.ShippingAddressID != nil {
		shippingAddress, appErr = a.srv.AccountService().AddressById(*checkOut.ShippingAddressID)
		if appErr != nil {
			if appErr.StatusCode == http.StatusInternalServerError {
				return nil, appErr
			}
			// ignore not found error
		}
	}

	var shippingMethod *shipping.ShippingMethod
	if checkOut.ShippingMethodID != nil {
		shippingMethod, appErr = a.srv.ShippingService().ShippingMethodByOption(&shipping.ShippingMethodFilterOption{
			Id: &model.StringFilter{
				StringOption: &model.StringOption{
					Eq: *checkOut.ShippingMethodID,
				},
			},
		})
		if appErr != nil {
			return nil, appErr
		}
	}
	// build shipping method channel listing filter option:
	shippingMethodChannelListingFilterOption := new(shipping.ShippingMethodChannelListingFilterOption)
	if chanNel != nil {
		shippingMethodChannelListingFilterOption.ChannelID = &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: chanNel.Id,
			},
		}
	}
	if shippingMethod != nil {
		shippingMethodChannelListingFilterOption.ShippingMethodID = &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: shippingMethod.Id,
			},
		}
	}
	var shippingMethodChannelListing *shipping.ShippingMethodChannelListing
	shippingMethodChannelListings, appErr := a.srv.ShippingService().ShippingMethodChannelListingsByOption(shippingMethodChannelListingFilterOption)
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr
		}
	} else {
		shippingMethodChannelListing = shippingMethodChannelListings[0]
	}

	var collectionPoint *warehouse.WareHouse
	if checkOut.CollectionPointID != nil {
		collectionPoint, appErr = a.srv.WarehouseService().WarehouseByOption(&warehouse.WarehouseFilterOption{
			Id: &model.StringFilter{
				StringOption: &model.StringOption{
					Eq: *checkOut.CollectionPointID,
				},
			},
		})
		if appErr != nil {
			if appErr.StatusCode == http.StatusInternalServerError {
				return nil, appErr
			}
			// ignore not found error
		}
	}

	checkout.CheckoutInfo{
		Checkout: *checkOut,
		User: ,
	}
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

// UpdateCheckoutInfoDeliveryMethod set CheckoutInfo's ShippingMethod to given shippingMethod
// and set new value for checkoutInfo's ShippingMethodChannelListings
// deliveryMethod must be either *ShippingMethod or *Warehouse or nil
func (a *ServiceCheckout) UpdateCheckoutInfoDeliveryMethod(checkoutInfo *checkout.CheckoutInfo, deliveryMethod interface{}) *model.AppError {
	// validate `deliveryMethod` is valid:
	if deliveryMethod != nil {
		switch deliveryMethod.(type) {
		case *warehouse.WareHouse, *shipping.ShippingMethod:
		default:
			return model.NewAppError("UpdateCheckoutInfoDeliveryMethod", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "deliveryMethod"}, "", http.StatusBadRequest)
		}
	}

	deliveryMethodIface, appErr := a.GetDeliveryMethodInfo(deliveryMethod, checkoutInfo.ShippingAddress)
	if appErr != nil {
		return appErr
	}

	checkoutInfo.DeliveryMethodInfo = deliveryMethodIface

	err := checkoutInfo.DeliveryMethodInfo.UpdateChannelListings(checkoutInfo)
	// if error is non-nil, this means we need another method that can access database Store
	if err != nil && err == checkout.ErrorNotUsable {
		appErr = a.updateChannelListings(checkoutInfo.DeliveryMethodInfo, checkoutInfo)
		if appErr != nil {
			return appErr
		}
	}

	return nil
}

func (s *ServiceCheckout) updateChannelListings(methodInfo checkout.DeliveryMethodBaseInterface, checkoutInfo *checkout.CheckoutInfo) *model.AppError {
	shippingMethodChannelListings, appErr := s.srv.ShippingService().ShippingMethodChannelListingsByOption(&shipping.ShippingMethodChannelListingFilterOption{
		ShippingMethodID: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: methodInfo.GetDeliveryMethod().(*shipping.ShippingMethod).Id,
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

	if len(shippingMethodChannelListings) > 0 {
		checkoutInfo.ShippingMethodChannelListings = shippingMethodChannelListings[0]
	}

	return nil
}

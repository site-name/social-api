package checkout

import (
	"context"
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/plugin/interfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

// GetDeliveryMethodInfo takes `deliveryMethod` is either *model.ShippingMethod or *model.Warehouse
func (s *ServiceCheckout) GetDeliveryMethodInfo(deliveryMethod interface{}, address *model.Address) (model.DeliveryMethodBaseInterface, *model.AppError) {
	if deliveryMethod == nil {
		return &model.DeliveryMethodBase{}, nil
	}

	switch t := deliveryMethod.(type) {
	case *model.ShippingMethod:
		return &model.ShippingMethodInfo{
			DeliveryMethod:  *t,
			ShippingAddress: address,
		}, nil

	case *model.WareHouse:
		address := t.GetAddress()
		if address == nil && t.AddressID != nil {
			var appErr *model.AppError
			address, appErr = s.srv.AccountService().AddressById(*t.AddressID)
			if appErr != nil {
				return nil, appErr
			}
		}
		return &model.CollectionPointInfo{
			DeliveryMethod:  *t,
			ShippingAddress: address,
		}, nil

	default:
		return nil, model.NewAppError("GetDeliveryMethodInfo", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "deliveryMethod"}, "", http.StatusBadRequest)
	}
}

// FetchCheckoutLines Fetch checkout lines as CheckoutLineInfo objects.
// It prefetch some related value also
func (a *ServiceCheckout) FetchCheckoutLines(checkOut *model.Checkout) ([]*model.CheckoutLineInfo, *model.AppError) {
	checkoutLineInfos, err := a.srv.Store.Checkout().FetchCheckoutLinesAndPrefetchRelatedValue(checkOut)
	if err != nil {
		return nil, model.NewAppError("FetchCheckoutLines", "app.checkout.error_collecting_checkout_line_infos.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return checkoutLineInfos, nil
}

// FetchCheckoutInfo Fetch checkout as CheckoutInfo object
func (a *ServiceCheckout) FetchCheckoutInfo(checkOut *model.Checkout, lines []*model.CheckoutLineInfo, discounts []*model.DiscountInfo, manager interfaces.PluginManagerInterface) (*model.CheckoutInfo, *model.AppError) {
	// validate arguments:
	if checkOut == nil {
		return nil, model.NewAppError("FetchCheckoutInfo", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "checkOut"}, "", http.StatusBadRequest)
	}

	chanNel, appErr := a.srv.ChannelService().ChannelByOption(&model.ChannelFilterOption{
		Id: squirrel.Eq{store.ChannelTableName + ".Id": checkOut.ChannelID},
	})
	if appErr != nil {
		return nil, appErr
	}

	// check if given checkout has both shipping address id and billing address id
	// then perform db lookup both of them in single query.
	var checkoutAddressIDs []string
	if checkOut.ShippingAddressID != nil {
		checkoutAddressIDs = append(checkoutAddressIDs, *checkOut.ShippingAddressID)
	}
	if checkOut.BillingAddressID != nil {
		checkoutAddressIDs = append(checkoutAddressIDs, *checkOut.BillingAddressID)
	}

	var (
		billingAddress  *model.Address
		shippingAddress *model.Address
	)
	if len(checkoutAddressIDs) > 0 {
		addresses, appErr := a.srv.AccountService().AddressesByOption(&model.AddressFilterOption{
			Id: squirrel.Eq{store.AddressTableName + ".Id": checkoutAddressIDs},
		})
		if appErr != nil {
			if appErr.StatusCode == http.StatusInternalServerError {
				return nil, appErr
			}
			// ignore not found error
		}
		if len(addresses) == 1 {
			if checkOut.ShippingAddressID != nil {
				shippingAddress = addresses[0]
			} else {
				billingAddress = addresses[0]
			}
		} else if len(addresses) == 2 {
			if *checkOut.ShippingAddressID == addresses[0].Id {
				shippingAddress = addresses[0]
				billingAddress = addresses[1]
			} else {
				shippingAddress = addresses[1]
				billingAddress = addresses[0]
			}
		}
	}

	var shippingMethod *model.ShippingMethod

	if checkOut.ShippingMethodID != nil {
		shippingMethod, appErr = a.srv.ShippingService().ShippingMethodByOption(&model.ShippingMethodFilterOption{
			Id: squirrel.Eq{store.ShippingMethodTableName + ".Id": *checkOut.ShippingMethodID},
		})
		if appErr != nil {
			if appErr.StatusCode == http.StatusInternalServerError {
				return nil, appErr
			}
			// ignore not found error
		}
	}

	var (
		shippingMethodChannelListing *model.ShippingMethodChannelListing
	)

	// build shipping method channel listings filter option:
	var shippingMethodChannelListingsFilterOption = new(model.ShippingMethodChannelListingFilterOption)

	if shippingMethod != nil {
		shippingMethodChannelListingsFilterOption.ShippingMethodID = squirrel.Eq{
			store.ShippingMethodChannelListingTableName + ".ShippingMethodID": shippingMethod.Id,
		}
	}

	if chanNel != nil {
		shippingMethodChannelListingsFilterOption.ChannelID = squirrel.Eq{
			store.ShippingMethodChannelListingTableName + ".ChannelID": chanNel.Id,
		}
	}

	shippingMethodChannelListings, appErr := a.srv.ShippingService().ShippingMethodChannelListingsByOption(shippingMethodChannelListingsFilterOption)
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr
		}
		// ignore not found error
	}
	shippingMethodChannelListing = shippingMethodChannelListings[0]

	var collectionPoint *model.WareHouse
	if checkOut.CollectionPointID != nil {
		collectionPoint, appErr = a.srv.WarehouseService().WarehouseByOption(&model.WarehouseFilterOption{
			Id:                   squirrel.Eq{store.WarehouseTableName + ".Id": *checkOut.CollectionPointID},
			SelectRelatedAddress: true,
		})
		if appErr != nil {
			if appErr.StatusCode == http.StatusInternalServerError {
				return nil, appErr
			}
			// ignore not found error
		}
	}
	var deliveryMethod interface{} = collectionPoint
	if deliveryMethod == nil {
		deliveryMethod = shippingMethod
	}

	deliveryMethodInfo, appErr := a.GetDeliveryMethodInfo(deliveryMethod, shippingAddress)
	if appErr != nil {
		return nil, appErr
	}

	var user *model.User
	if checkOut.UserID != nil {
		user, appErr = a.srv.AccountService().UserById(context.Background(), *checkOut.UserID)
		if appErr != nil {
			if appErr.StatusCode == http.StatusInternalServerError {
				return nil, appErr
			}
		}
	}

	checkoutInfo := model.CheckoutInfo{
		Checkout:                      *checkOut,
		User:                          user,
		Channel:                       *chanNel,
		BillingAddress:                billingAddress,
		ShippingAddress:               shippingAddress,
		DeliveryMethodInfo:            deliveryMethodInfo,
		ShippingMethodChannelListings: shippingMethodChannelListing,
	}

	validShippingMethods, appErr := a.GetValidShippingMethodListForCheckoutInfo(checkoutInfo, shippingAddress, lines, discounts, manager)
	if appErr != nil {
		return nil, appErr
	}

	validPickupPoints, appErr := a.GetValidCollectionPointsForCheckoutInfo(shippingAddress, lines, &checkoutInfo)
	if appErr != nil {
		return nil, appErr
	}

	checkoutInfo.ValidShippingMethods = validShippingMethods
	checkoutInfo.ValidPickupPoints = validPickupPoints
	checkoutInfo.DeliveryMethodInfo = deliveryMethodInfo

	return &checkoutInfo, nil
}

// UpdateCheckoutInfoShippingAddress updates given `checkoutInfo` by setting given `address` as its ShippingAddress.
// then updates its ValidShippingMethods
func (a *ServiceCheckout) UpdateCheckoutInfoShippingAddress(checkoutInfo model.CheckoutInfo, address *model.Address, lines []*model.CheckoutLineInfo, discounts []*model.DiscountInfo, manager interfaces.PluginManagerInterface) *model.AppError {
	checkoutInfo.ShippingAddress = address
	validMethods, appErr := a.GetValidShippingMethodListForCheckoutInfo(checkoutInfo, address, lines, discounts, manager)
	if appErr != nil {
		return appErr
	}

	checkoutInfo.ValidShippingMethods = validMethods
	deliveryMethod := checkoutInfo.DeliveryMethodInfo.GetDeliveryMethod()
	checkoutInfo.DeliveryMethodInfo, appErr = a.GetDeliveryMethodInfo(deliveryMethod, address)
	return appErr
}

// GetValidShippingMethodListForCheckoutInfo
func (a *ServiceCheckout) GetValidShippingMethodListForCheckoutInfo(checkoutInfo model.CheckoutInfo, shippingAddress *model.Address, lines []*model.CheckoutLineInfo, discounts []*model.DiscountInfo, manager interfaces.PluginManagerInterface) ([]*model.ShippingMethod, *model.AppError) {
	var countryCode model.CountryCode
	if shippingAddress != nil {
		countryCode = shippingAddress.Country
	}

	subTotal, appErr := manager.CalculateCheckoutSubTotal(checkoutInfo, lines, checkoutInfo.ShippingAddress, discounts)
	if appErr != nil {
		return nil, appErr
	}

	checkoutInfo.Checkout.PopulateNonDbFields() // this is important
	subTotal, err := subTotal.Sub(checkoutInfo.Checkout.Discount)
	if err != nil {
		return nil, model.NewAppError("GetValidShippingMethodListForCheckoutInfo", app.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return a.GetValidShippingMethodsForCheckout(checkoutInfo, lines, subTotal, countryCode)
}

func (s *ServiceCheckout) GetValidCollectionPointsForCheckoutInfo(shippingAddress *model.Address, lines []*model.CheckoutLineInfo, checkoutInfo *model.CheckoutInfo) ([]*model.WareHouse, *model.AppError) {
	var countryCode model.CountryCode

	if shippingAddress != nil {
		countryCode = shippingAddress.Country
	} else {
		countryCode = checkoutInfo.Channel.DefaultCountry
	}

	return s.GetValidCollectionPointsForCheckout(lines, countryCode, false)
}

// UpdateCheckoutInfoDeliveryMethod set CheckoutInfo's ShippingMethod to given shippingMethod
// and set new value for checkoutInfo's ShippingMethodChannelListings
// deliveryMethod must be either *ShippingMethod or *Warehouse or nil
func (a *ServiceCheckout) UpdateCheckoutInfoDeliveryMethod(checkoutInfo model.CheckoutInfo, deliveryMethod interface{}) *model.AppError {
	// validate `deliveryMethod` is valid:
	if deliveryMethod != nil {
		switch deliveryMethod.(type) {
		case *model.WareHouse, *model.ShippingMethod:
		default:
			return model.NewAppError("UpdateCheckoutInfoDeliveryMethod", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "deliveryMethod"}, "", http.StatusBadRequest)
		}
	}

	deliveryMethodIface, appErr := a.GetDeliveryMethodInfo(deliveryMethod, checkoutInfo.ShippingAddress)
	if appErr != nil {
		return appErr
	}

	checkoutInfo.DeliveryMethodInfo = deliveryMethodIface

	err := checkoutInfo.DeliveryMethodInfo.UpdateChannelListings(&checkoutInfo)
	// if error is non-nil, this means we need another method that can access database Store
	if err != nil && err == model.ErrorNotUsable {
		appErr = a.updateChannelListings(checkoutInfo.DeliveryMethodInfo, checkoutInfo)
		if appErr != nil {
			return appErr
		}
	}

	return nil
}

func (s *ServiceCheckout) updateChannelListings(methodInfo model.DeliveryMethodBaseInterface, checkoutInfo model.CheckoutInfo) *model.AppError {
	shippingMethodChannelListings, appErr := s.srv.ShippingService().ShippingMethodChannelListingsByOption(&model.ShippingMethodChannelListingFilterOption{
		ShippingMethodID: squirrel.Eq{store.ShippingMethodChannelListingTableName + ".ShippingMethodID": methodInfo.GetDeliveryMethod().(*model.ShippingMethod).Id},
		ChannelID:        squirrel.Eq{store.ShippingMethodChannelListingTableName + ".ChannelID": checkoutInfo.Channel.Id},
	})
	if appErr != nil {
		return appErr
	}

	if len(shippingMethodChannelListings) > 0 {
		checkoutInfo.ShippingMethodChannelListings = shippingMethodChannelListings[0]
	}

	return nil
}

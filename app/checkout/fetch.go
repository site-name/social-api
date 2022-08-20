package checkout

import (
	"context"
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/plugin/interfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/channel"
	"github.com/sitename/sitename/model/checkout"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/model/shipping"
	"github.com/sitename/sitename/model/warehouse"
	"github.com/sitename/sitename/store"
)

// GetDeliveryMethodInfo takes `deliveryMethod` is either *ShippingMethodData or *Warehouse
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

// FetchCheckoutInfo Fetch checkout as CheckoutInfo object
func (a *ServiceCheckout) FetchCheckoutInfo(checkOut *checkout.Checkout, lines []*checkout.CheckoutLineInfo, discounts []*product_and_discount.DiscountInfo, manager interfaces.PluginManagerInterface) (*checkout.CheckoutInfo, *model.AppError) {
	// validate arguments:
	if checkOut == nil {
		return nil, model.NewAppError("FetchCheckoutInfo", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "checkOut"}, "", http.StatusBadRequest)
	}

	chanNel, appErr := a.srv.ChannelService().ChannelByOption(&channel.ChannelFilterOption{
		Id: squirrel.Eq{a.srv.Store.Channel().TableName("Id"): checkOut.ChannelID},
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
		billingAddress  *account.Address
		shippingAddress *account.Address
	)
	if len(checkoutAddressIDs) > 0 {
		addresses, appErr := a.srv.AccountService().AddressesByOption(&account.AddressFilterOption{
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

	var shippingMethod *shipping.ShippingMethod

	if checkOut.ShippingMethodID != nil {
		shippingMethod, appErr = a.srv.ShippingService().ShippingMethodByOption(&shipping.ShippingMethodFilterOption{
			Id: squirrel.Eq{a.srv.Store.ShippingMethod().TableName("Id"): *checkOut.ShippingMethodID},
		})
		if appErr != nil {
			if appErr.StatusCode == http.StatusInternalServerError {
				return nil, appErr
			}
			// ignore not found error
		}
	}

	var (
		shippingMethodChannelListing *shipping.ShippingMethodChannelListing
	)

	// build shipping method channel listings filter option:
	var shippingMethodChannelListingsFilterOption = new(shipping.ShippingMethodChannelListingFilterOption)

	if shippingMethod != nil {
		shippingMethodChannelListingsFilterOption.ShippingMethodID = squirrel.Eq{
			a.srv.Store.ShippingMethodChannelListing().TableName("ShippingMethodID"): shippingMethod.Id,
		}
	}

	if chanNel != nil {
		shippingMethodChannelListingsFilterOption.ChannelID = squirrel.Eq{
			a.srv.Store.ShippingMethodChannelListing().TableName("ChannelID"): chanNel.Id,
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

	var collectionPoint *warehouse.WareHouse
	if checkOut.CollectionPointID != nil {
		collectionPoint, appErr = a.srv.WarehouseService().WarehouseByOption(&warehouse.WarehouseFilterOption{
			Id:                   squirrel.Eq{a.srv.Store.Warehouse().TableName("Id"): *checkOut.CollectionPointID},
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

	var user *account.User
	if checkOut.UserID != nil {
		user, appErr = a.srv.AccountService().UserById(context.Background(), *checkOut.UserID)
		if appErr != nil {
			if appErr.StatusCode == http.StatusInternalServerError {
				return nil, appErr
			}
		}
	}

	checkoutInfo := checkout.CheckoutInfo{
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
func (a *ServiceCheckout) UpdateCheckoutInfoShippingAddress(checkoutInfo checkout.CheckoutInfo, address *account.Address, lines []*checkout.CheckoutLineInfo, discounts []*product_and_discount.DiscountInfo, manager interfaces.PluginManagerInterface) *model.AppError {
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
func (a *ServiceCheckout) GetValidShippingMethodListForCheckoutInfo(checkoutInfo checkout.CheckoutInfo, shippingAddress *account.Address, lines []*checkout.CheckoutLineInfo, discounts []*product_and_discount.DiscountInfo, manager interfaces.PluginManagerInterface) ([]*shipping.ShippingMethod, *model.AppError) {
	var countryCode string
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

func (s *ServiceCheckout) GetValidCollectionPointsForCheckoutInfo(shippingAddress *account.Address, lines []*checkout.CheckoutLineInfo, checkoutInfo *checkout.CheckoutInfo) ([]*warehouse.WareHouse, *model.AppError) {
	var countryCode string

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
func (a *ServiceCheckout) UpdateCheckoutInfoDeliveryMethod(checkoutInfo checkout.CheckoutInfo, deliveryMethod interface{}) *model.AppError {
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

	err := checkoutInfo.DeliveryMethodInfo.UpdateChannelListings(&checkoutInfo)
	// if error is non-nil, this means we need another method that can access database Store
	if err != nil && err == checkout.ErrorNotUsable {
		appErr = a.updateChannelListings(checkoutInfo.DeliveryMethodInfo, checkoutInfo)
		if appErr != nil {
			return appErr
		}
	}

	return nil
}

func (s *ServiceCheckout) updateChannelListings(methodInfo checkout.DeliveryMethodBaseInterface, checkoutInfo checkout.CheckoutInfo) *model.AppError {
	shippingMethodChannelListings, appErr := s.srv.ShippingService().ShippingMethodChannelListingsByOption(&shipping.ShippingMethodChannelListingFilterOption{
		ShippingMethodID: squirrel.Eq{s.srv.Store.ShippingMethodChannelListing().TableName("ShippingMethodID"): methodInfo.GetDeliveryMethod().(*shipping.ShippingMethod).Id},
		ChannelID:        squirrel.Eq{s.srv.Store.ShippingMethodChannelListing().TableName("ChannelID"): checkoutInfo.Channel.Id},
	})
	if appErr != nil {
		return appErr
	}

	if len(shippingMethodChannelListings) > 0 {
		checkoutInfo.ShippingMethodChannelListings = shippingMethodChannelListings[0]
	}

	return nil
}

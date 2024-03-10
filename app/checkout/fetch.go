package checkout

import (
	"context"
	"net/http"

	"github.com/sitename/sitename/app/plugin/interfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// GetDeliveryMethodInfo takes `deliveryMethod` is either model.ShippingMethod or model.Warehouse
func (s *ServiceCheckout) GetDeliveryMethodInfo(deliveryMethod any, address *model.Address) (model_helper.DeliveryMethodBaseInterface, *model_helper.AppError) {
	if deliveryMethod == nil {
		return &model_helper.DeliveryMethodBase{}, nil
	}

	switch t := deliveryMethod.(type) {
	case model.ShippingMethod:
		return &model_helper.ShippingMethodInfo{
			DeliveryMethod:  t,
			ShippingAddress: address,
		}, nil

	case model.Warehouse:
		var address *model.Address
		if t.R != nil && t.R.Address != nil {
			address = t.R.Address
		}

		if address == nil && !t.AddressID.IsNil() {
			var appErr *model_helper.AppError
			address, appErr = s.srv.Account.AddressById(*t.AddressID.String)
			if appErr != nil {
				return nil, appErr
			}
		}
		return &model_helper.CollectionPointInfo{
			DeliveryMethod:  t,
			ShippingAddress: address,
		}, nil

	default:
		return nil, model_helper.NewAppError("GetDeliveryMethodInfo", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "deliveryMethod"}, "", http.StatusBadRequest)
	}
}

// FetchCheckoutLines Fetch checkout lines as CheckoutLineInfo objects.
// It prefetch some related value also
func (a *ServiceCheckout) FetchCheckoutLines(checkOut model.Checkout) (model_helper.CheckoutLineInfos, *model_helper.AppError) {
	checkoutLineInfos, err := a.srv.Store.Checkout().FetchCheckoutLinesAndPrefetchRelatedValue(checkOut)
	if err != nil {
		return nil, model_helper.NewAppError("FetchCheckoutLines", "app.checkout.error_collecting_checkout_line_infos.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return checkoutLineInfos, nil
}

func (a *ServiceCheckout) FetchCheckoutInfo(checkout model.Checkout, lines model_helper.CheckoutLineInfos, discounts []*model_helper.DiscountInfo, manager interfaces.PluginManagerInterface) (*model_helper.CheckoutInfo, *model_helper.AppError) {
	channel, appErr := a.srv.Channel.ChannelByOption(model_helper.ChannelFilterOptions{
		CommonQueryOptions: model_helper.NewCommonQueryOptions(model.ChannelWhere.ID.EQ(checkout.ChannelID)),
	})
	if appErr != nil {
		return nil, appErr
	}

	// check if given checkout has both shipping address id and billing address id
	// then perform db lookup both of them in single query.
	var checkoutAddressIDs []string
	if !checkout.ShippingAddressID.IsNil() {
		checkoutAddressIDs = append(checkoutAddressIDs, *checkout.ShippingAddressID.String)
	}
	if !checkout.BillingAddressID.IsNil() {
		checkoutAddressIDs = append(checkoutAddressIDs, *checkout.BillingAddressID.String)
	}

	var (
		billingAddress  *model.Address
		shippingAddress *model.Address
	)
	if len(checkoutAddressIDs) > 0 {
		addresses, appErr := a.srv.Account.AddressesByOption(model_helper.AddressFilterOptions{
			CommonQueryOptions: model_helper.NewCommonQueryOptions(model.AddressWhere.ID.IN(checkoutAddressIDs)),
		})
		if appErr != nil {
			if appErr.StatusCode == http.StatusInternalServerError {
				return nil, appErr
			}
			// ignore not found error
		}
		if len(addresses) == 1 {
			if !checkout.ShippingAddressID.IsNil() {
				shippingAddress = addresses[0]
			} else {
				billingAddress = addresses[0]
			}
		} else if len(addresses) == 2 {
			if !checkout.ShippingAddressID.IsNil() && *checkout.ShippingAddressID.String == addresses[0].ID {
				shippingAddress = addresses[0]
				billingAddress = addresses[1]
			} else {
				shippingAddress = addresses[1]
				billingAddress = addresses[0]
			}
		}
	}

	var shippingMethod *model.ShippingMethod

	if !checkout.ShippingMethodID.IsNil() {
		shippingMethod, appErr = a.srv.Shipping.ShippingMethodByOption(model_helper.ShippingMethodFilterOption{
			CommonQueryOptions: model_helper.NewCommonQueryOptions(model.ShippingMethodWhere.ID.EQ(*checkout.ShippingMethodID.String)),
		})
		if appErr != nil {
			if appErr.StatusCode == http.StatusInternalServerError {
				return nil, appErr
			}
			// ignore not found error
		}
	}

	andConds := []qm.QueryMod{}
	if shippingMethod != nil {
		andConds = append(andConds, model.ShippingMethodChannelListingWhere.ShippingMethodID.EQ(shippingMethod.ID))
	}
	if channel != nil {
		andConds = append(andConds, model.ShippingMethodChannelListingWhere.ChannelID.EQ(channel.ID))
	}

	shippingMethodChannelListings, appErr := a.srv.Shipping.ShippingMethodChannelListingsByOption(model_helper.ShippingMethodChannelListingFilterOption{
		CommonQueryOptions: model_helper.NewCommonQueryOptions(andConds...),
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr
		}
		// ignore not found error
	}

	var shippingMethodChannelListing *model.ShippingMethodChannelListing
	if len(shippingMethodChannelListings) > 0 {
		shippingMethodChannelListing = shippingMethodChannelListings[0]
	}

	var collectionPoint *model.Warehouse
	if !checkout.CollectionPointID.IsNil() {
		collectionPoint, appErr = a.srv.Warehouse.WarehouseByOption(model_helper.WarehouseFilterOption{
			CommonQueryOptions: model_helper.NewCommonQueryOptions(model.WarehouseWhere.ID.EQ(*checkout.CollectionPointID.String)),
			Preloads:           []string{model.WarehouseRels.Address},
		})
		if appErr != nil {
			if appErr.StatusCode == http.StatusInternalServerError {
				return nil, appErr
			}
			// ignore not found error
		}
	}
	var deliveryMethod any = collectionPoint
	if deliveryMethod.(*model.Warehouse) == nil {
		deliveryMethod = shippingMethod
	}

	deliveryMethodInfo, appErr := a.GetDeliveryMethodInfo(deliveryMethod, shippingAddress)
	if appErr != nil {
		return nil, appErr
	}

	var user *model.User
	if !checkout.UserID.IsNil() {
		user, appErr = a.srv.Account.UserById(context.Background(), *checkout.UserID.String)
		if appErr != nil {
			if appErr.StatusCode == http.StatusInternalServerError {
				return nil, appErr
			}
		}
	}

	checkoutInfo := model_helper.CheckoutInfo{
		Checkout:                      checkout,
		User:                          user,
		Channel:                       *channel,
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

// updateCheckoutInfoShippingAddress updates given `checkoutInfo` by setting given `address` as its ShippingAddress.
// then updates its ValidShippingMethods
func (a *ServiceCheckout) updateCheckoutInfoShippingAddress(checkoutInfo model_helper.CheckoutInfo, address *model.Address, lines model_helper.CheckoutLineInfos, discounts []*model_helper.DiscountInfo, manager interfaces.PluginManagerInterface) *model_helper.AppError {
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

func (a *ServiceCheckout) GetValidShippingMethodListForCheckoutInfo(checkoutInfo model_helper.CheckoutInfo, shippingAddress *model.Address, lines model_helper.CheckoutLineInfos, discounts []*model_helper.DiscountInfo, manager interfaces.PluginManagerInterface) (model.ShippingMethodSlice, *model_helper.AppError) {
	var countryCode model.CountryCode
	if shippingAddress != nil {
		countryCode = shippingAddress.Country
	}

	subTotal, appErr := manager.CalculateCheckoutSubTotal(checkoutInfo, lines, checkoutInfo.ShippingAddress, discounts)
	if appErr != nil {
		return nil, appErr
	}

	subTotal, err := subTotal.Sub(model_helper.CheckoutGetDiscountMoney(checkoutInfo.Checkout))
	if err != nil {
		return nil, model_helper.NewAppError("GetValidShippingMethodListForCheckoutInfo", model_helper.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return a.GetValidShippingMethodsForCheckout(checkoutInfo, lines, subTotal, countryCode)
}

func (s *ServiceCheckout) GetValidCollectionPointsForCheckoutInfo(shippingAddress *model.Address, lines model_helper.CheckoutLineInfos, checkoutInfo *model_helper.CheckoutInfo) (model.WarehouseSlice, *model_helper.AppError) {
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
func (a *ServiceCheckout) UpdateCheckoutInfoDeliveryMethod(checkoutInfo model_helper.CheckoutInfo, deliveryMethod any) *model_helper.AppError {
	// validate `deliveryMethod` is valid:
	if deliveryMethod != nil {
		switch deliveryMethod.(type) {
		case *model.Warehouse, *model.ShippingMethod:
		default:
			return model_helper.NewAppError("UpdateCheckoutInfoDeliveryMethod", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "deliveryMethod"}, "", http.StatusBadRequest)
		}
	}

	deliveryMethodIface, appErr := a.GetDeliveryMethodInfo(deliveryMethod, checkoutInfo.ShippingAddress)
	if appErr != nil {
		return appErr
	}

	checkoutInfo.DeliveryMethodInfo = deliveryMethodIface

	err := checkoutInfo.DeliveryMethodInfo.UpdateChannelListings(&checkoutInfo)
	// if error is non-nil, this means we need another method that can access database Store
	if err != nil && err == model_helper.ErrorNotUsable {
		appErr = a.updateChannelListings(checkoutInfo.DeliveryMethodInfo, checkoutInfo)
		if appErr != nil {
			return appErr
		}
	}

	return nil
}

func (s *ServiceCheckout) updateChannelListings(methodInfo model_helper.DeliveryMethodBaseInterface, checkoutInfo model_helper.CheckoutInfo) *model_helper.AppError {
	shippingMethodChannelListings, appErr := s.srv.Shipping.ShippingMethodChannelListingsByOption(model_helper.ShippingMethodChannelListingFilterOption{
		CommonQueryOptions: model_helper.NewCommonQueryOptions(
			model.ShippingMethodChannelListingWhere.ShippingMethodID.EQ(methodInfo.GetDeliveryMethod().(*model.ShippingMethod).ID),
			model.ShippingMethodChannelListingWhere.ChannelID.EQ(checkoutInfo.Channel.ID),
		),
	})
	if appErr != nil {
		return appErr
	}

	if len(shippingMethodChannelListings) > 0 {
		checkoutInfo.ShippingMethodChannelListings = shippingMethodChannelListings[0]
	}

	return nil
}

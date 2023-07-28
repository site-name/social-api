package shipping

import (
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/samber/lo"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/measurement"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

// ApplicableShippingMethodsForCheckout finds all applicable shipping methods for given checkout, based on given additional arguments
func (a *ServiceShipping) ApplicableShippingMethodsForCheckout(checkout *model.Checkout, channelID string, price *goprices.Money, countryCode model.CountryCode, lines []*model.CheckoutLineInfo) ([]*model.ShippingMethod, *model.AppError) {
	if checkout.ShippingAddressID == nil {
		return nil, nil
	}

	var appErr *model.AppError
	if countryCode == "" {
		countryCode, appErr = a.srv.CheckoutService().CheckoutCountry(checkout)
		if appErr != nil {
			return nil, appErr
		}
	}

	var checkoutProductIDs []string

	// check if checkout line infos are provided
	if len(lines) == 0 {
		_, _, products, err := a.srv.Store.CheckoutLine().CheckoutLinesByCheckoutWithPrefetch(checkout.Token)
		if err != nil {
			return nil, model.NewAppError("ApplicableShippingMethodsForCheckout", "app.shipping.get_applicable_shipping_methods_for_checkout.app_error", nil, err.Error(), http.StatusInternalServerError)
		}

		// if product(s) was/were found
		for _, product := range products {
			checkoutProductIDs = append(checkoutProductIDs, product.Id)
		}
	} else {
		for _, info := range lines {
			checkoutProductIDs = append(checkoutProductIDs, info.Product.Id)
		}
	}

	// calculate total weight of this checkout:
	var checkoutLineIDs []string
	for _, lineInfo := range lines {
		checkoutLineIDs = append(checkoutLineIDs, lineInfo.Line.Id)
	}
	totalWeight, err := a.srv.Store.CheckoutLine().TotalWeightForCheckoutLines(checkoutLineIDs)
	if err != nil {
		if _, ok := err.(*store.ErrNotFound); !ok {
			return nil, model.NewAppError("ApplicableShippingMethodsForCheckout", "app.shipping.get_total_weight_for_checkout.app_error", nil, err.Error(), http.StatusInternalServerError)
		}
		totalWeight = measurement.ZeroWeight
	}

	shippingMethods, err := a.srv.Store.ShippingMethod().ApplicableShippingMethods(
		price,
		channelID,
		totalWeight,
		countryCode,
		checkoutLineIDs,
	)

	if err != nil {
		return nil, model.NewAppError("ApplicableShippingMethodsForCheckout", "app.shipping.shipping_methods_for_checkout.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	checkoutShippingAddress, appErr := a.srv.AccountService().AddressById(*checkout.ShippingAddressID)
	if appErr != nil {
		return nil, appErr
	}
	return a.FilterShippingMethodsByPostalCodeRules(shippingMethods, checkoutShippingAddress), nil
}

// ApplicableShippingMethodsForOrder finds all applicable shippingmethods for given order, based on other arguments passed in
func (a *ServiceShipping) ApplicableShippingMethodsForOrder(order *model.Order, channelID string, price *goprices.Money, countryCode model.CountryCode, lines []*model.CheckoutLineInfo) ([]*model.ShippingMethod, *model.AppError) {
	if order.ShippingAddressID == nil {
		return nil, nil
	}

	if !countryCode.IsValid() {
		address, appErr := a.srv.AccountService().AddressById(*order.ShippingAddressID)
		if appErr != nil {
			return nil, appErr
		}
		countryCode = address.Country
	}

	var orderProductIDs []string
	if len(lines) == 0 {
		orderLines, appErr := a.srv.OrderService().OrderLinesByOption(&model.OrderLineFilterOption{
			Conditions: squirrel.Eq{model.OrderLineTableName + ".OrderID": order.Id},
			PrefetchRelated: model.OrderLinePrefetchRelated{
				VariantProduct: true, // this tells store to prefetch related product variants, products too
			},
		})
		if appErr != nil {
			return nil, appErr
		}

		for _, orderLine := range orderLines {
			if variant := orderLine.GetProductVariant(); variant != nil && variant.GetProduct() != nil {
				orderProductIDs = append(orderProductIDs, variant.GetProduct().Id)
			}
		}
	} else {
		for _, line := range lines {
			orderProductIDs = append(orderProductIDs, line.Product.Id)
		}
	}

	applicableShippingMethods, err := a.srv.Store.ShippingMethod().ApplicableShippingMethods(
		price,
		channelID,
		order.GetTotalWeight(),
		countryCode,
		orderProductIDs,
	)
	if err != nil {
		return nil, model.NewAppError("ApplicableShippingMethodsForOrder", "app.shipping.shipping_methods_for_order.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	orderShippingAddress, appErr := a.srv.AccountService().AddressById(*order.ShippingAddressID)
	if appErr != nil {
		return nil, appErr
	}
	return a.FilterShippingMethodsByPostalCodeRules(applicableShippingMethods, orderShippingAddress), nil
}

// ShippingMethodByOption returns a shipping method with given options
func (s *ServiceShipping) ShippingMethodByOption(option *model.ShippingMethodFilterOption) (*model.ShippingMethod, *model.AppError) {
	method, err := s.srv.Store.ShippingMethod().GetbyOption(option)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}
		return nil, model.NewAppError("ServiceShipping.ShippingMethodByOption", "app.shipping.error_finding_shipping_method_by_option.app_error", nil, err.Error(), statusCode)
	}

	return method, nil
}

// ShippingMethodsByOptions finds and returns all shipping methods that satisfy given fiter options
func (s *ServiceShipping) ShippingMethodsByOptions(options *model.ShippingMethodFilterOption) ([]*model.ShippingMethod, *model.AppError) {
	methods, err := s.srv.Store.ShippingMethod().FilterByOptions(options)
	if err != nil {
		return nil, model.NewAppError("ShippingMethodsByOptions", "app.shipping.error_finding_shipping_methods.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return methods, nil
}

func (s *ServiceShipping) DropInvalidShippingMethodsRelationsForGivenChannels(transaction *gorm.DB, shippingMethodIds, channelIds []string) *model.AppError {
	// unlink shipping methods from order and checkout instances
	// when method is no longer available in given channels
	checkouts, appErr := s.srv.CheckoutService().CheckoutsByOption(&model.CheckoutFilterOption{
		Conditions: squirrel.Eq{
			model.CheckoutTableName + ".ShippingMethodID": shippingMethodIds,
			model.CheckoutTableName + ".ChannelID":        channelIds,
		},
	})
	if appErr != nil {
		return appErr
	}

	lo.ForEach(checkouts, func(checkout *model.Checkout, _ int) { checkout.ShippingMethodID = nil })
	_, appErr = s.srv.CheckoutService().UpsertCheckouts(transaction, checkouts)
	if appErr != nil {
		return appErr
	}

	orders, appErr := s.srv.OrderService().FilterOrdersByOptions(&model.OrderFilterOption{
		Conditions: squirrel.Eq{
			model.OrderTableName + ".Status": []model.OrderStatus{
				model.ORDER_STATUS_UNCONFIRMED,
				model.ORDER_STATUS_DRAFT,
			},
			model.OrderTableName + ".ShippingMethodID": shippingMethodIds,
			model.OrderTableName + ".ChannelID":        channelIds,
		},
	})
	if appErr != nil {
		return appErr
	}

	lo.ForEach(orders, func(order *model.Order, _ int) { order.ShippingMethodID = nil })
	_, appErr = s.srv.OrderService().BulkUpsertOrders(transaction, orders)
	if appErr != nil {
		return appErr
	}

	return nil
}

func (s *ServiceShipping) UpsertShippingMethod(transaction *gorm.DB, method *model.ShippingMethod) (*model.ShippingMethod, *model.AppError) {
	method, err := s.srv.Store.ShippingMethod().Upsert(transaction, method)
	if err != nil {
		return nil, model.NewAppError("UpsertShippingMethod", "app.shipping.error_upserting_shipping_method.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return method, nil
}

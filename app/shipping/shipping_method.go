package shipping

import (
	"net/http"

	"github.com/mattermost/squirrel"
	"github.com/samber/lo"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/measurement"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

// ApplicableShippingMethodsForCheckout finds all applicable shipping methods for given checkout, based on given additional arguments
func (a *ServiceShipping) ApplicableShippingMethodsForCheckout(checkout model.Checkout, channelID string, price goprices.Money, countryCode model.CountryCode, lines model_helper.CheckoutLineInfos) (model.ShippingMethodSlice, *model_helper.AppError) {
	if checkout.ShippingAddressID.IsNil() {
		return nil, nil
	}

	var appErr *model_helper.AppError
	if countryCode.IsValid() != nil {
		countryCode, appErr = a.srv.Checkout.CheckoutCountry(checkout)
		if appErr != nil {
			return nil, appErr
		}
	}

	var checkoutProductIDs []string

	// check if checkout line infos are provided
	if len(lines) == 0 {
		_, _, products, err := a.srv.Store.CheckoutLine().CheckoutLinesByCheckoutWithPrefetch(checkout.Token)
		if err != nil {
			return nil, model_helper.NewAppError("ApplicableShippingMethodsForCheckout", "app.shipping.get_applicable_shipping_methods_for_checkout.app_error", nil, err.Error(), http.StatusInternalServerError)
		}

		// if product(s) was/were found
		for _, product := range products {
			checkoutProductIDs = append(checkoutProductIDs, product.Id)
		}
	} else {
		for _, info := range lines {
			checkoutProductIDs = append(checkoutProductIDs, info.Product.ID)
		}
	}

	// calculate total weight of this checkout:
	var checkoutLineIDs []string
	for _, lineInfo := range lines {
		checkoutLineIDs = append(checkoutLineIDs, lineInfo.Line.ID)
	}
	totalWeight, err := a.srv.Store.CheckoutLine().TotalWeightForCheckoutLines(checkoutLineIDs)
	if err != nil {
		if _, ok := err.(*store.ErrNotFound); !ok {
			return nil, model_helper.NewAppError("ApplicableShippingMethodsForCheckout", "app.shipping.get_total_weight_for_checkout.app_error", nil, err.Error(), http.StatusInternalServerError)
		}
		totalWeight = &measurement.Weight{
			Amount: 0,
			Unit:   measurement.KG,
		}
	}

	shippingMethods, err := a.srv.Store.ShippingMethod().ApplicableShippingMethods(
		price,
		channelID,
		*totalWeight,
		countryCode,
		checkoutLineIDs,
	)

	if err != nil {
		return nil, model_helper.NewAppError("ApplicableShippingMethodsForCheckout", "app.shipping.shipping_methods_for_checkout.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	checkoutShippingAddress, appErr := a.srv.Account.AddressById(*checkout.ShippingAddressID.String)
	if appErr != nil {
		return nil, appErr
	}
	return a.FilterShippingMethodsByPostalCodeRules(shippingMethods, *checkoutShippingAddress), nil
}

// ApplicableShippingMethodsForOrder finds all applicable shippingmethods for given order, based on other arguments passed in
func (a *ServiceShipping) ApplicableShippingMethodsForOrder(order model.Order, channelID string, price goprices.Money, countryCode model.CountryCode, lines model_helper.CheckoutLineInfos) (model.ShippingMethodSlice, *model_helper.AppError) {
	if order.ShippingAddressID.IsNil() {
		return nil, nil
	}

	orderShippingAddress, appErr := a.srv.Account.AddressById(*order.ShippingAddressID.String)
	if appErr != nil {
		return nil, appErr
	}

	if countryCode.IsValid() != nil {
		countryCode = orderShippingAddress.Country
	}

	var orderProductIDs []string
	if len(lines) == 0 {
		orderLines, appErr := a.srv.Order.OrderLinesByOption(model_helper.OrderLineFilterOptions{
			CommonQueryOptions: model_helper.NewCommonQueryOptions(
				model.OrderLineWhere.OrderID.EQ(order.ID),
			),
			Preload: []string{model.OrderLineRels.Variant},
		})
		if appErr != nil {
			return nil, appErr
		}

		orderProductIDs = lo.Map(orderLines, func(o *model.OrderLine, _ int) string { return o.R.Variant.ProductID })
	} else {
		orderProductIDs = lo.Map(lines, func(item *model_helper.CheckoutLineInfo, _ int) string { return item.Product.ID })
	}

	applicableShippingMethods, err := a.srv.Store.ShippingMethod().ApplicableShippingMethods(
		price,
		channelID,
		model_helper.OrderGetTotalWeight(order),
		countryCode,
		orderProductIDs,
	)
	if err != nil {
		return nil, model_helper.NewAppError("ApplicableShippingMethodsForOrder", "app.shipping.shipping_methods_for_order.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return a.FilterShippingMethodsByPostalCodeRules(applicableShippingMethods, *orderShippingAddress), nil
}

func (s *ServiceShipping) ShippingMethodByOption(option model_helper.ShippingMethodFilterOption) (*model.ShippingMethod, *model_helper.AppError) {
	method, err := s.srv.Store.ShippingMethod().GetbyOption(option)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}
		return nil, model_helper.NewAppError("ServiceShipping.ShippingMethodByOption", "app.shipping.error_finding_shipping_method_by_option.app_error", nil, err.Error(), statusCode)
	}

	return method, nil
}

func (s *ServiceShipping) ShippingMethodsByOptions(options model_helper.ShippingMethodFilterOption) (model.ShippingMethodSlice, *model_helper.AppError) {
	methods, err := s.srv.Store.ShippingMethod().FilterByOptions(options)
	if err != nil {
		return nil, model_helper.NewAppError("ShippingMethodsByOptions", "app.shipping.error_finding_shipping_methods.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return methods, nil
}

func (s *ServiceShipping) DropInvalidShippingMethodsRelationsForGivenChannels(transaction boil.ContextTransactor, shippingMethodIds, channelIds []string) *model_helper.AppError {
	// unlink shipping methods from order and checkout instances
	// when method is no longer available in given channels
	_, checkouts, appErr := s.srv.Checkout.CheckoutsByOption(&model.CheckoutFilterOption{
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

	_, orders, appErr := s.srv.OrderService().FilterOrdersByOptions(&model.OrderFilterOption{
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

func (s *ServiceShipping) UpsertShippingMethod(transaction boil.ContextTransactor, method *model.ShippingMethod) (*model.ShippingMethod, *model_helper.AppError) {
	method, err := s.srv.Store.ShippingMethod().Upsert(transaction, method)
	if err != nil {
		return nil, model_helper.NewAppError("UpsertShippingMethod", "app.shipping.error_upserting_shipping_method.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return method, nil
}

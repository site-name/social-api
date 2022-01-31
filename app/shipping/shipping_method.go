package shipping

import (
	"net/http"

	"github.com/Masterminds/squirrel"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/checkout"
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/model/shipping"
	"github.com/sitename/sitename/modules/measurement"
	"github.com/sitename/sitename/store"
)

// ApplicableShippingMethodsForCheckout finds all applicable shipping methods for given checkout, based on given additional arguments
func (a *ServiceShipping) ApplicableShippingMethodsForCheckout(ckout *checkout.Checkout, channelID string, price *goprices.Money, countryCode string, lines []*checkout.CheckoutLineInfo) ([]*shipping.ShippingMethod, *model.AppError) {
	if ckout.ShippingAddressID == nil || !model.IsValidId(*ckout.ShippingAddressID) {
		return nil, nil
	}

	var appErr *model.AppError
	if countryCode == "" {
		countryCode, appErr = a.srv.CheckoutService().CheckoutCountry(ckout)
		if appErr != nil {
			return nil, appErr
		}
	}

	var checkoutProductIDs []string

	// check if checkout line infos are provided
	if len(lines) == 0 {
		_, _, products, err := a.srv.Store.CheckoutLine().CheckoutLinesByCheckoutWithPrefetch(ckout.Token)
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

	return a.FilterShippingMethodsByPostalCodeRules(shippingMethods, *ckout.ShippingAddressID) // already checked ShippingAddressID == nil
}

// ApplicableShippingMethodsForOrder finds all applicable shippingmethods for given order, based on other arguments passed in
func (a *ServiceShipping) ApplicableShippingMethodsForOrder(oder *order.Order, channelID string, price *goprices.Money, countryCode string, lines []*checkout.CheckoutLineInfo) ([]*shipping.ShippingMethod, *model.AppError) {
	if oder.ShippingAddressID == nil || !model.IsValidId(*oder.ShippingAddressID) {
		return nil, nil
	}

	if countryCode == "" {
		address, appErr := a.srv.AccountService().AddressById(*oder.ShippingAddressID)
		if appErr != nil {
			return nil, appErr
		}
		countryCode = address.Country
	}

	var orderProductIDs []string
	if len(lines) == 0 {
		orderLines, appErr := a.srv.OrderService().OrderLinesByOption(&order.OrderLineFilterOption{
			OrderID: squirrel.Eq{a.srv.Store.OrderLine().TableName("OrderID"): oder.Id},
			PrefetchRelated: order.OrderLinePrefetchRelated{
				VariantProduct: true, // this tells store to prefetch related product variants, products too
			},
		})
		if appErr != nil {
			return nil, appErr
		}

		for _, orderLine := range orderLines {
			if orderLine.ProductVariant != nil && orderLine.ProductVariant.Product != nil {
				orderProductIDs = append(orderProductIDs, orderLine.ProductVariant.Product.Id)
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
		oder.GetTotalWeight(),
		countryCode,
		orderProductIDs,
	)

	if err != nil {
		return nil, model.NewAppError("ApplicableShippingMethodsForOrder", "app.shipping.shipping_methods_for_order.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return a.FilterShippingMethodsByPostalCodeRules(applicableShippingMethods, *oder.ShippingAddressID) // already checked ShippingAddressID == nil
}

// ShippingMethodByOption returns a shipping method with given options
func (s *ServiceShipping) ShippingMethodByOption(option *shipping.ShippingMethodFilterOption) (*shipping.ShippingMethod, *model.AppError) {
	method, err := s.srv.Store.ShippingMethod().GetbyOption(option)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("ShippingMethodByOption", "app.shipping.error_finding_shipping_method_by_option.app_error", err)
	}

	return method, nil
}

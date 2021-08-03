package shipping

import (
	"net/http"

	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/checkout"
	"github.com/sitename/sitename/modules/measurement"
	"github.com/sitename/sitename/store"
)

func (a *AppShipping) ApplicableShippingMethodsForCheckout(ckout *checkout.Checkout, channelID string, price *goprices.Money, countryCode string, lines []*checkout.CheckoutLineInfo) (interface{}, *model.AppError) {
	if ckout.ShippingAddressID == nil {
		return nil, nil
	}

	var appErr *model.AppError
	if countryCode == "" {
		countryCode, appErr = a.CheckoutApp().CheckoutCountry(ckout)
		if appErr != nil {
			return nil, appErr
		}
	}

	var checkoutProductIDs []string

	// check if checkout line infos are provided
	if len(lines) == 0 {
		_, _, products, err := a.Srv().Store.CheckoutLine().CheckoutLinesByCheckoutWithPrefetch(ckout.Token)
		if err != nil {
			if _, ok := err.(*store.ErrNotFound); !ok {
				// returns if error is caused by system
				return nil, model.NewAppError("ApplicableShippingMethodsForCheckout", "app.shipping.get_applicable_shipping_methods_for_checkout.app_error", nil, err.Error(), http.StatusInternalServerError)
			}
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
	totalWeight, err := a.Srv().Store.CheckoutLine().TotalWeightForCheckoutLines(checkoutLineIDs)
	if err != nil {
		if _, ok := err.(*store.ErrNotFound); !ok {
			return nil, model.NewAppError("ApplicableShippingMethodsForCheckout", "app.shipping.get_total_weight_for_checkout.app_error", nil, err.Error(), http.StatusInternalServerError)
		}
		totalWeight = measurement.ZeroWeight
	}

	a.Srv().Store.ShippingMethod().ApplicableShippingMethods(
		price,
		channelID,
		totalWeight,
		countryCode,
		checkoutLineIDs,
	)

	return nil, nil
}

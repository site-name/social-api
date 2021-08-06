package checkout

import (
	"net/http"
	"strings"

	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/checkout"
	"github.com/sitename/sitename/model/shipping"
	"github.com/sitename/sitename/modules/util"
)

// BaseCalculationShippingPrice Return checkout shipping price.
func (a *AppCheckout) BaseCalculationShippingPrice(checkoutInfo *checkout.CheckoutInfo, lineInfos []*checkout.CheckoutLineInfo) (*goprices.TaxedMoney, *model.AppError) {
	var (
		shippingRequired bool
		appErr           *model.AppError
	)

	if len(lineInfos) > 0 {
		productIDs := []string{}
		for _, info := range lineInfos {
			productIDs = append(productIDs, info.Product.Id)
		}

		shippingRequired, appErr = a.app.ProductApp().ProductsRequireShipping(productIDs)
	} else {
		shippingRequired, appErr = a.CheckoutShippingRequired(checkoutInfo.Checkout.Token)
	}

	if appErr != nil {
		return nil, appErr
	}

	if checkoutInfo.ShippingMethod == nil || !shippingRequired {
		// ignore error here since checkouts were validated before saving into database
		taxedMoney, _ := util.ZeroTaxedMoney(checkoutInfo.Checkout.Currency)
		return taxedMoney, nil
	}

	shippingMethodChannelListings, appErr := a.app.ShippingApp().ShippingMethodChannelListingsByOption(&shipping.ShippingMethodChannelListingFilterOption{
		ShippingMethodID: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: checkoutInfo.ShippingMethod.Id,
			},
		},
		ChannelID: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: checkoutInfo.Checkout.ChannelID,
			},
		},
	})
	if appErr != nil {
		return nil, appErr
	}

	shippingMethodChannelListings[0].GetTotal().Quantize()
}

// BaseCheckoutTotal returns the total cost of the checkout
func (a *AppCheckout) BaseCheckoutTotal(subTotal *goprices.TaxedMoney, shippingPrice *goprices.TaxedMoney, discount *goprices.TaxedMoney, currency string) (*goprices.TaxedMoney, *model.AppError) {
	currencyMap := map[string]bool{}
	currencyMap[subTotal.Currency] = true
	currencyMap[shippingPrice.Currency] = true
	currencyMap[discount.Currency] = true
	currencyMap[strings.ToLower(currency)] = true

	if len(currencyMap) > 1 {
		return nil, model.NewAppError("BaseCheckoutTotal", "app.checkout.invalid_currencies.app_error", nil, "Please pass in the same currency values", http.StatusBadRequest)
	}
	errorID := "app.checkout.error_calculating_taxed_money.app_error"

	total, err := subTotal.Add(shippingPrice)
	if err != nil {
		return nil, model.NewAppError("BaseCheckoutTotal", errorID, nil, err.Error(), http.StatusInternalServerError)
	}

	total, err = total.Sub(discount)
	if err != nil {
		return nil, model.NewAppError("BaseCheckoutTotal", errorID, nil, err.Error(), http.StatusInternalServerError)
	}

	zeroTaxedMoney, err := util.ZeroTaxedMoney(currency)
	if err != nil {
		return nil, model.NewAppError("BaseCheckoutTotal", "app.checkout.error_creating_taxed_money.app_error", nil, err.Error(), http.StatusBadRequest)
	}

	if lessThanOrEqual, _ := zeroTaxedMoney.LessThanOrEqual(total); lessThanOrEqual {
		return total, nil
	}

	return zeroTaxedMoney, nil
}

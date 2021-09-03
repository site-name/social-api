package checkout

import (
	"net/http"
	"strings"

	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/channel"
	"github.com/sitename/sitename/model/checkout"
	"github.com/sitename/sitename/model/order"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/model/shipping"
	"github.com/sitename/sitename/modules/util"
)

// BaseCalculationShippingPrice Return checkout shipping price.
func (a *ServiceCheckout) BaseCalculationShippingPrice(checkoutInfo *checkout.CheckoutInfo, lineInfos []*checkout.CheckoutLineInfo) (*goprices.TaxedMoney, *model.AppError) {
	var (
		shippingRequired bool
		appErr           *model.AppError
	)

	if len(lineInfos) > 0 {
		productIDs := []string{}
		for _, info := range lineInfos {
			productIDs = append(productIDs, info.Product.Id)
		}

		shippingRequired, appErr = a.srv.ProductService().ProductsRequireShipping(productIDs)
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

	shippingMethodChannelListings, appErr := a.srv.ShippingService().
		ShippingMethodChannelListingsByOption(&shipping.ShippingMethodChannelListingFilterOption{
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

	shippingPrice := shippingMethodChannelListings[0].GetTotal()
	res, _ := (&goprices.TaxedMoney{
		Net:      shippingPrice,
		Gross:    shippingPrice,
		Currency: shippingPrice.Currency,
	}).Quantize()

	return res, nil
}

// BaseCheckoutTotal returns the total cost of the checkout
func (a *ServiceCheckout) BaseCheckoutTotal(subTotal *goprices.TaxedMoney, shippingPrice *goprices.TaxedMoney, discount *goprices.TaxedMoney, currency string) (*goprices.TaxedMoney, *model.AppError) {
	// this method reqires all values's currencies are uppoer-cased and supported by system
	currency = strings.ToUpper(currency)
	currencyMap := map[string]bool{}
	currencyMap[subTotal.Currency] = true
	currencyMap[shippingPrice.Currency] = true
	currencyMap[discount.Currency] = true
	currencyMap[currency] = true

	if _, err := goprices.GetCurrencyPrecision(currency); err != nil || len(currencyMap) > 1 {
		return nil, model.NewAppError("BaseCheckoutTotal", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "money fields"}, "Please pass in the same currency values", http.StatusBadRequest)
	}

	total, _ := subTotal.Add(shippingPrice)
	total, _ = total.Sub(discount)

	zeroTaxedMoney, _ := util.ZeroTaxedMoney(currency)
	if lessThanOrEqual, _ := zeroTaxedMoney.LessThanOrEqual(total); lessThanOrEqual {
		return total, nil
	}

	return zeroTaxedMoney, nil
}

// BaseCheckoutLineTotal Return the total price of this line
//
// `discounts` can be nil
func (a *ServiceCheckout) BaseCheckoutLineTotal(checkoutLineInfo *checkout.CheckoutLineInfo, channel *channel.Channel, discounts []*product_and_discount.DiscountInfo) (*goprices.TaxedMoney, *model.AppError) {
	if discounts == nil {
		discounts = []*product_and_discount.DiscountInfo{}
	}

	variantPrice, appErr := a.srv.ProductService().ProductVariantGetPrice(
		&checkoutLineInfo.Product,
		checkoutLineInfo.Collections,
		channel,
		checkoutLineInfo.ChannelListing,
		discounts,
	)
	if appErr != nil {
		return nil, appErr
	}

	amount, _ := variantPrice.Mul(int(checkoutLineInfo.Line.Quantity))
	amount, _ = amount.Quantize()

	return &goprices.TaxedMoney{
		Net:      amount,
		Gross:    amount,
		Currency: amount.Currency,
	}, nil
}

func (a *ServiceCheckout) BaseOrderLineTotal(orderLine *order.OrderLine) (*goprices.TaxedMoney, *model.AppError) {
	orderLine.PopulateNonDbFields()
	if orderLine.UnitPrice != nil {
		unitPrice, _ := orderLine.UnitPrice.Mul(int(orderLine.Quantity))
		unitPrice, _ = unitPrice.Quantize()

		return unitPrice, nil
	}

	return nil, model.NewAppError("BaseOrderLineTotal", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "orderLine"}, "", http.StatusBadRequest)
}

func (a *ServiceCheckout) BaseTaxRate(price *goprices.TaxedMoney) (*decimal.Decimal, *model.AppError) {
	taxRate := &decimal.Zero
	if price != nil && price.Gross != nil && !price.Gross.Amount.Equal(decimal.Zero) {
		tax, _ := price.Tax()
		div, _ := tax.TrueDiv(price.Net)
		taxRate = div.Amount
	}

	return taxRate, nil
}

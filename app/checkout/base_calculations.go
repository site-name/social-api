package checkout

import (
	"net/http"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/util"
)

// BaseCheckoutShippingPrice
func (s *ServiceCheckout) BaseCheckoutShippingPrice(checkoutInfo *model_helper.CheckoutInfo, lines model.CheckoutLineInfos) (*goprices.TaxedMoney, *model_helper.AppError) {
	deliveryMethodInfo := checkoutInfo.DeliveryMethodInfo.Self()
	if shippingMethodInfo, ok := deliveryMethodInfo.(*model.ShippingMethodInfo); ok {
		return s.CalculatePriceForShippingMethod(checkoutInfo, shippingMethodInfo, lines)
	}

	zeroTaxed, _ := util.ZeroTaxedMoney(checkoutInfo.Checkout.Currency.String())
	return zeroTaxed, nil
}

// CalculatePriceForShippingMethod Return checkout shipping price
func (s *ServiceCheckout) CalculatePriceForShippingMethod(checkoutInfo *model_helper.CheckoutInfo, shippingMethodInfo *model.ShippingMethodInfo, lines model.CheckoutLineInfos) (*goprices.TaxedMoney, *model_helper.AppError) {
	// validate input arguments
	if checkoutInfo == nil {
		return nil, model_helper.NewAppError("CalculatePriceForShippingMethod", model_helper.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "checkoutInfo"}, "", http.StatusBadRequest)
	}

	var (
		shippingMethod   = shippingMethodInfo.DeliveryMethod
		shippingRequired bool
		appErr           *model_helper.AppError
	)

	if lines != nil {
		shippingRequired, appErr = s.srv.ProductService().ProductsRequireShipping(lines.Products().IDs())
	} else {
		shippingRequired, appErr = s.srv.CheckoutService().CheckoutShippingRequired(checkoutInfo.Checkout.Token)
	}

	if appErr != nil {
		return nil, appErr
	}

	if !model_helper.IsValidId(shippingMethod.Id) || !shippingRequired {
		zeroTaxedMoney, _ := util.ZeroTaxedMoney(checkoutInfo.Checkout.Currency)
		return zeroTaxedMoney, nil
	}

	shippingMethodChannelListingsOfShippingMethod, appErr := s.srv.ShippingService().
		ShippingMethodChannelListingsByOption(&model.ShippingMethodChannelListingFilterOption{
			Conditions: squirrel.Eq{
				model.ShippingMethodChannelListingTableName + ".ShippingMethodID": shippingMethod.Id,
				model.ShippingMethodChannelListingTableName + ".ChannelID":        checkoutInfo.Checkout.ChannelID,
			},
		})
	if appErr != nil {
		return nil, appErr
	}

	shippingPrice := shippingMethodChannelListingsOfShippingMethod[0].GetTotal()
	taxedMoney, _ := goprices.NewTaxedMoney(shippingPrice, shippingPrice)

	quantizedPrice, _ := taxedMoney.Quantize(goprices.Up, -1)
	return quantizedPrice, nil
}

// BaseCheckoutTotal returns the total cost of the checkout
//
// NOTE: discount must be either Money, TaxedMoney, *Money, *TaxedMoney
func (a *ServiceCheckout) BaseCheckoutTotal(subTotal *goprices.TaxedMoney, shippingPrice *goprices.TaxedMoney, discount interface{}, currency string) (*goprices.TaxedMoney, *model_helper.AppError) {
	// valudate input
	switch discount.(type) {
	case *goprices.Money, *goprices.TaxedMoney, goprices.Money, goprices.TaxedMoney:
	default:
		return nil, model_helper.NewAppError("BaseCheckoutTotal", model_helper.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "discount"}, "discount must be either Money or TaxedMoney", http.StatusBadRequest)
	}

	// this method reqires all values's currencies are uppoer-cased and supported by system
	currency = strings.ToUpper(currency)
	currencyMap := map[string]bool{}
	currencyMap[subTotal.GetCurrency()] = true
	currencyMap[shippingPrice.GetCurrency()] = true
	currencyMap[discount.(goprices.Currencier).GetCurrency()] = true // validated in the beginning
	currencyMap[currency] = true

	if _, err := goprices.GetCurrencyPrecision(currency); err != nil || len(currencyMap) > 1 {
		return nil, model_helper.NewAppError("BaseCheckoutTotal", model_helper.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "money fields"}, "Please pass in the same currency values", http.StatusBadRequest)
	}

	total, _ := subTotal.Add(shippingPrice)
	total, _ = total.Sub(discount)

	zeroTaxedMoney, _ := util.ZeroTaxedMoney(currency)
	if zeroTaxedMoney.LessThanOrEqual(total) {
		return total, nil
	}

	return zeroTaxedMoney, nil
}

// BaseCheckoutLineTotal Return the total price of this line
//
// `discounts` can be nil
func (a *ServiceCheckout) BaseCheckoutLineTotal(checkoutLineInfo *model.CheckoutLineInfo, channel *model.Channel, discounts []*model.DiscountInfo) (*goprices.TaxedMoney, *model_helper.AppError) {
	if discounts == nil {
		discounts = []*model.DiscountInfo{}
	}

	variantPrice, appErr := a.srv.ProductService().ProductVariantGetPrice(
		&checkoutLineInfo.Variant,
		checkoutLineInfo.Product,
		checkoutLineInfo.Collections,
		*channel,
		&checkoutLineInfo.ChannelListing,
		discounts,
	)
	if appErr != nil {
		return nil, appErr
	}

	amount := variantPrice.Mul(float64(checkoutLineInfo.Line.Quantity))
	amount, _ = amount.Quantize(goprices.Up, -1)

	return &goprices.TaxedMoney{
		Net:      amount,
		Gross:    amount,
		Currency: amount.Currency,
	}, nil
}

func (a *ServiceCheckout) BaseOrderLineTotal(orderLine *model.OrderLine) (*goprices.TaxedMoney, *model_helper.AppError) {
	orderLine.PopulateNonDbFields()
	if orderLine.UnitPrice != nil {
		unitPrice := orderLine.UnitPrice.Mul(float64(orderLine.Quantity))
		unitPrice, _ = unitPrice.Quantize(goprices.Up, -1)

		return unitPrice, nil
	}

	return nil, model_helper.NewAppError("BaseOrderLineTotal", model_helper.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "orderLine"}, "", http.StatusBadRequest)
}

func (a *ServiceCheckout) BaseTaxRate(price goprices.TaxedMoney) (*decimal.Decimal, *model_helper.AppError) {
	taxRate := decimal.Zero
	if !price.Gross.Amount.IsZero() {
		tax := price.Tax()
		div := tax.TrueDiv(price.Net.Amount.InexactFloat64())
		taxRate = div.Amount
	}

	return &taxRate, nil
}

// BaseCheckoutLineUnitPrice divide given totalLinePrice to given quantity and returns the result
func (a *ServiceCheckout) BaseCheckoutLineUnitPrice(totalLinePrice *goprices.TaxedMoney, quantity int) *goprices.TaxedMoney {
	return totalLinePrice.TrueDiv(float64(quantity))
}

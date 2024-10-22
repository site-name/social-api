package checkout

import (
	"net/http"

	"github.com/samber/lo"
	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/util"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (s *ServiceCheckout) BaseCheckoutShippingPrice(checkoutInfo model_helper.CheckoutInfo, lines model_helper.CheckoutLineInfos) (*goprices.TaxedMoney, *model_helper.AppError) {
	deliveryMethodInfo := checkoutInfo.DeliveryMethodInfo.Self()

	if shippingMethodInfo, ok := deliveryMethodInfo.(model_helper.ShippingMethodInfo); ok {
		return s.CalculatePriceForShippingMethod(checkoutInfo, shippingMethodInfo, lines)
	}

	zeroTaxed, _ := util.ZeroTaxedMoney(checkoutInfo.Checkout.Currency.String())
	return zeroTaxed, nil
}

func (s *ServiceCheckout) CalculatePriceForShippingMethod(checkoutInfo model_helper.CheckoutInfo, shippingMethodInfo model_helper.ShippingMethodInfo, lines model_helper.CheckoutLineInfos) (*goprices.TaxedMoney, *model_helper.AppError) {
	var (
		shippingMethod   = shippingMethodInfo.DeliveryMethod
		shippingRequired bool
		appErr           *model_helper.AppError
	)

	if len(lines) > 0 {
		productIDs := lo.Map(lines.Products(), func(item *model.Product, _ int) string { return item.ID })
		shippingRequired, appErr = s.srv.Product.ProductsRequireShipping(productIDs)
	} else {
		shippingRequired, appErr = s.srv.Checkout.CheckoutShippingRequired(checkoutInfo.Checkout.Token)
	}
	if appErr != nil {
		return nil, appErr
	}

	if !model_helper.IsValidId(shippingMethod.ID) || !shippingRequired {
		zeroTaxedMoney, _ := util.ZeroTaxedMoney(checkoutInfo.Checkout.Currency.String())
		return zeroTaxedMoney, nil
	}

	shippingMethodChannelListingsOfShippingMethod, appErr := s.srv.Shipping.
		ShippingMethodChannelListingsByOption(
			model_helper.ShippingMethodChannelListingFilterOption{
				CommonQueryOptions: model_helper.NewCommonQueryOptions(
					model.ShippingMethodChannelListingWhere.ShippingMethodID.EQ(shippingMethod.ID),
					model.ShippingMethodChannelListingWhere.ChannelID.EQ(checkoutInfo.Checkout.ChannelID),
					qm.Limit(1),
				),
			},
		)
	if appErr != nil {
		return nil, appErr
	}

	shippingPrice := model_helper.ShippingMethodChannelListingGetTotal(shippingMethodChannelListingsOfShippingMethod[0])
	taxedMoney, _ := goprices.NewTaxedMoney(shippingPrice, shippingPrice)

	quantizedPrice, _ := taxedMoney.Quantize(goprices.Up, -1)
	return quantizedPrice, nil
}

// BaseCheckoutTotal returns the total cost of the checkout
//
// NOTE: discount must be either Money, TaxedMoney, *Money, *TaxedMoney
func (a *ServiceCheckout) BaseCheckoutTotal(subTotal goprices.TaxedMoney, shippingPrice goprices.TaxedMoney, discount any, currency model.Currency) (*goprices.TaxedMoney, *model_helper.AppError) {
	switch discount.(type) {
	case *goprices.Money, *goprices.TaxedMoney, goprices.Money, goprices.TaxedMoney:
	default:
		return nil, model_helper.NewAppError("BaseCheckoutTotal", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "discount"}, "discount must be either Money or TaxedMoney", http.StatusBadRequest)
	}

	// this method reqires all values's currencies are upper-cased and supported by system
	currencyMap := map[string]bool{}
	currencyMap[subTotal.GetCurrency()] = true
	currencyMap[shippingPrice.GetCurrency()] = true
	currencyMap[discount.(goprices.Currencier).GetCurrency()] = true // validated in the beginning
	currencyMap[currency.String()] = true

	if _, err := goprices.GetCurrencyPrecision(currency.String()); err != nil || len(currencyMap) > 1 {
		return nil, model_helper.NewAppError("BaseCheckoutTotal", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "money fields"}, "Please pass in the same currency values", http.StatusBadRequest)
	}

	total, _ := subTotal.Add(shippingPrice)
	total, _ = total.Sub(discount)

	zeroTaxedMoney, _ := util.ZeroTaxedMoney(currency.String())
	if zeroTaxedMoney.LessThanOrEqual(*total) {
		return total, nil
	}

	return zeroTaxedMoney, nil
}

// BaseCheckoutLineTotal Return the total price of this line
//
// `discounts` can be nil
func (a *ServiceCheckout) BaseCheckoutLineTotal(checkoutLineInfo model_helper.CheckoutLineInfo, channel model.Channel, discounts []*model_helper.DiscountInfo) (*goprices.TaxedMoney, *model_helper.AppError) {
	variantPrice, appErr := a.srv.Product.ProductVariantGetPrice(
		&checkoutLineInfo.Variant,
		checkoutLineInfo.Product,
		checkoutLineInfo.Collections,
		channel,
		&checkoutLineInfo.ChannelListing,
		discounts,
	)
	if appErr != nil {
		return nil, appErr
	}

	amount := variantPrice.Mul(float64(checkoutLineInfo.Line.Quantity))
	quantizedAmount, _ := amount.Quantize(goprices.Up, -1)

	return &goprices.TaxedMoney{
		Net:   *quantizedAmount,
		Gross: *quantizedAmount,
	}, nil
}

func (a *ServiceCheckout) BaseOrderLineTotal(orderLine model.OrderLine) (*goprices.TaxedMoney, *model_helper.AppError) {
	orderLineUnitPrice := model_helper.OrderLineGetUnitPrice(orderLine)

	unitPrice := orderLineUnitPrice.Mul(float64(orderLine.Quantity))
	quantizedUnitPrice, _ := unitPrice.Quantize(goprices.Up, -1)
	return quantizedUnitPrice, nil
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
func (a *ServiceCheckout) BaseCheckoutLineUnitPrice(totalLinePrice goprices.TaxedMoney, quantity int) *goprices.TaxedMoney {
	res := totalLinePrice.TrueDiv(float64(quantity))
	return &res
}

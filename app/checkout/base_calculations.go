package checkout

import (
	"net/http"
	"strings"

	"github.com/Masterminds/squirrel"
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

// BaseCheckoutShippingPrice
func (s *ServiceCheckout) BaseCheckoutShippingPrice(checkoutInfo *checkout.CheckoutInfo, lines checkout.CheckoutLineInfos) (*goprices.TaxedMoney, *model.AppError) {
	deliveryMethodInfo := checkoutInfo.DeliveryMethodInfo.Self()
	if shippingMethodInfo, ok := deliveryMethodInfo.(*checkout.ShippingMethodInfo); ok {
		return s.CalculatePriceForShippingMethod(checkoutInfo, shippingMethodInfo, lines)
	}

	zeroTaxed, _ := util.ZeroTaxedMoney(checkoutInfo.Checkout.Currency)
	return zeroTaxed, nil
}

// CalculatePriceForShippingMethod Return checkout shipping price
func (s *ServiceCheckout) CalculatePriceForShippingMethod(checkoutInfo *checkout.CheckoutInfo, shippingMethodInfo *checkout.ShippingMethodInfo, lines checkout.CheckoutLineInfos) (*goprices.TaxedMoney, *model.AppError) {
	// validate input arguments
	if checkoutInfo == nil {
		return nil, model.NewAppError("CalculatePriceForShippingMethod", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "checkoutInfo"}, "", http.StatusBadRequest)
	}

	var (
		shippingMethod   = shippingMethodInfo.DeliveryMethod
		shippingRequired bool
		appErr           *model.AppError
	)

	if lines != nil {
		shippingRequired, appErr = s.srv.ProductService().ProductsRequireShipping(lines.Products().IDs())
	} else {
		shippingRequired, appErr = s.srv.CheckoutService().CheckoutShippingRequired(checkoutInfo.Checkout.Token)
	}

	if appErr != nil {
		return nil, appErr
	}

	if !model.IsValidId(shippingMethod.Id) || !shippingRequired {
		zeroTaxedMoney, _ := util.ZeroTaxedMoney(checkoutInfo.Checkout.Currency)
		return zeroTaxedMoney, nil
	}

	shippingMethodChannelListingsOfShippingMethod, appErr := s.srv.ShippingService().ShippingMethodChannelListingsByOption(&shipping.ShippingMethodChannelListingFilterOption{
		ShippingMethodID: squirrel.Eq{s.srv.Store.ShippingMethodChannelListing().TableName("ShippingMethodID"): shippingMethod.Id},
		ChannelID:        squirrel.Eq{s.srv.Store.ShippingMethodChannelListing().TableName("ChannelID"): checkoutInfo.Checkout.ChannelID},
	})
	if appErr != nil {
		return nil, appErr
	}

	shippingPrice := shippingMethodChannelListingsOfShippingMethod[0].GetTotal()
	taxedMoney, _ := goprices.NewTaxedMoney(shippingPrice, shippingPrice)

	quantizedPrice, _ := taxedMoney.Quantize(nil, goprices.Up)
	return quantizedPrice, nil
}

// BaseCheckoutTotal returns the total cost of the checkout
//
// NOTE: discount must be either Money, TaxedMoney, *Money, *TaxedMoney
func (a *ServiceCheckout) BaseCheckoutTotal(subTotal *goprices.TaxedMoney, shippingPrice *goprices.TaxedMoney, discount interface{}, currency string) (*goprices.TaxedMoney, *model.AppError) {
	// valudate input
	switch discount.(type) {
	case *goprices.Money, *goprices.TaxedMoney, goprices.Money, goprices.TaxedMoney:
	default:
		return nil, model.NewAppError("BaseCheckoutTotal", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "discount"}, "discount must be either Money or TaxedMoney", http.StatusBadRequest)
	}

	// this method reqires all values's currencies are uppoer-cased and supported by system
	currency = strings.ToUpper(currency)
	currencyMap := map[string]bool{}
	currencyMap[subTotal.Currency] = true
	currencyMap[shippingPrice.Currency] = true
	currencyMap[discount.(goprices.Currencyable).MyCurrency()] = true // validated in the beginning
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
		&checkoutLineInfo.Variant,
		&checkoutLineInfo.Product,
		checkoutLineInfo.Collections,
		channel,
		&checkoutLineInfo.ChannelListing,
		discounts,
	)
	if appErr != nil {
		return nil, appErr
	}

	amount, _ := variantPrice.Mul(int(checkoutLineInfo.Line.Quantity))
	amount, _ = amount.Quantize(nil, goprices.Up)

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
		unitPrice, _ = unitPrice.Quantize(nil, goprices.Up)

		return unitPrice, nil
	}

	return nil, model.NewAppError("BaseOrderLineTotal", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "orderLine"}, "", http.StatusBadRequest)
}

func (a *ServiceCheckout) BaseTaxRate(price *goprices.TaxedMoney) (*decimal.Decimal, *model.AppError) {
	taxRate := decimal.Zero
	if price != nil && price.Gross != nil && !price.Gross.Amount.IsZero() {
		tax, _ := price.Tax()
		div, _ := tax.TrueDiv(price.Net)
		taxRate = div.Amount
	}

	return &taxRate, nil
}

// BaseCheckoutLineUnitPrice divide given totalLinePrice to given quantity and returns the result
func (a *ServiceCheckout) BaseCheckoutLineUnitPrice(totalLinePrice *goprices.TaxedMoney, quantity int) (*goprices.TaxedMoney, *model.AppError) {
	div, err := totalLinePrice.TrueDiv(quantity)
	if err != nil {
		return nil, model.NewAppError("BaseCheckoutLineUnitPrice", app.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return div, nil
}

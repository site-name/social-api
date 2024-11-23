package checkout

import (
	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/util"
)

func (s *ServiceCheckout) BaseCheckoutShippingPrice(checkoutInfo model_helper.CheckoutInfo, lines model_helper.CheckoutLineInfos) (*goprices.TaxedMoney, *model_helper.AppError) {
	deliveryMethodInfo := checkoutInfo.DeliveryMethodInfo.Self()

	if shippingMethodInfo, ok := deliveryMethodInfo.(model_helper.ShippingMethodInfo); ok {
		return s.CalculatePriceForShippingMethod(checkoutInfo, shippingMethodInfo, lines)
	}

	zeroTaxed, _ := util.ZeroTaxedMoney(checkoutInfo.Checkout.Currency.String())
	return zeroTaxed, nil
}

func (s *ServiceCheckout) CalculateBasePriceForShippingMethod(checkoutInfo model_helper.CheckoutInfo, shippingMethodInfo model_helper.ShippingMethodInfo) (*goprices.Money, *model_helper.AppError) {
	shippingRequired, appErr := s.CheckoutShippingRequired(checkoutInfo.Checkout.Token)
	if appErr != nil {
		return nil, appErr
	}

	if !shippingRequired {
		money, _ := util.ZeroMoney(checkoutInfo.Checkout.Currency)
		return money, nil
	}

	result, err := goprices.QuantizePrice(&shippingMethodInfo.DeliveryMethod.Price, goprices.Up)
	if err != nil {
		return nil, model_helper.NewAppError("CalculateBasePriceForShippingMethod", model_helper.ErrorCalculatingMoneyErrorID, nil, err.Error(), 0)
	}

	return result, nil
}

// BaseCheckoutTotal returns the total cost of the checkout
//
// The price includes catalogue promotions, shipping, specific product
// and applied once per order voucher discounts.
// The price does not include order promotions and the entire order vouchers.
func (a *ServiceCheckout) BaseCheckoutTotal(checkoutInfo model_helper.CheckoutInfo, lines model_helper.CheckoutLineInfos) (*goprices.Money, *model_helper.AppError) {
	subTotal, appErr := a.BaseCheckoutSubTotal(lines, checkoutInfo.Channel, checkoutInfo.Checkout.Currency, true)
	if appErr != nil {
		// return nil, appErr
		// lines[0].
	}
}

func (s *ServiceCheckout) BaseCheckoutDeliveryPrice(checkoutInfo model_helper.CheckoutInfo, lines model_helper.CheckoutLineInfos, includeVoucher bool) (*goprices.Money, *model_helper.AppError) {
	shippingPrice, appErr := s.BaseCheckoutUndiscountedDeliveryPrice(checkoutInfo, lines)
	if appErr != nil {
		return nil, appErr
	}

	isShippingVoucher := checkoutInfo.Voucher
}

func (s *ServiceCheckout) BaseCheckoutUndiscountedDeliveryPrice(checkoutInfo model_helper.CheckoutInfo, lines model_helper.CheckoutLineInfos) (*goprices.Money, *model_helper.AppError) {
	switch checkoutInfo.DeliveryMethodInfo.(type) {
	case model_helper.ShippingMethodInfo:
		money, _ := util.ZeroMoney(checkoutInfo.Checkout.Currency)
		return money, nil

	default:
		return s.CalculateBasePriceForShippingMethod(checkoutInfo, checkoutInfo.DeliveryMethodInfo, lines)
	}
}

// Return the checkout subtotal value.
//
// The price includes catalogue promotions, specific product and applied once per order
// voucher discounts.
// The price does not include order promotions and the entire order vouchers.
func (c *ServiceCheckout) BaseCheckoutSubTotal(checkoutLines model_helper.CheckoutLineInfos, _ model.Channel, currency model.Currency, includeVoucher bool) (*goprices.Money, *model_helper.AppError) {
	var result, _ = util.ZeroMoney(currency)

	for _, line := range checkoutLines {
		if line == nil {
			continue
		}
		money, appErr := c.CalculateBaseLineTotalPrice(*line, includeVoucher)
		if appErr != nil {
			return nil, appErr
		}
		result.Add(*money)
	}

	return result, nil
}

func (c *ServiceCheckout) CalculateBaseLineTotalPrice(lineInfo model_helper.CheckoutLineInfo, includeVoucher bool) (*goprices.Money, *model_helper.AppError) {
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

package checkout

import (
	"net/http"

	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/plugin/interfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
)

// CheckoutShippingPrice Return checkout shipping price.
//
// It takes in account all plugins.
func (s *ServiceCheckout) CheckoutShippingPrice(manager interfaces.PluginManagerInterface, checkoutInfo model.CheckoutInfo, lines []*model.CheckoutLineInfo, address *model.Address, discounts []*model.DiscountInfo) (*goprices.TaxedMoney, *model.AppError) {
	if discounts == nil {
		discounts = []*model.DiscountInfo{}
	}
	calculatedCheckoutShipping, appErr := manager.CalculateCheckoutShipping(checkoutInfo, lines, address, discounts)
	if appErr != nil {
		return nil, appErr
	}

	calculatedCheckoutShipping, err := calculatedCheckoutShipping.Quantize(nil, goprices.Up)
	if err != nil {
		return nil, model.NewAppError("CheckoutShippingPrice", app.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return calculatedCheckoutShipping, nil
}

// CheckoutSubTotal Return the total cost of all the checkout lines, taxes included.
func (s *ServiceCheckout) CheckoutSubTotal(manager interfaces.PluginManagerInterface, checkoutInfo model.CheckoutInfo, lines []*model.CheckoutLineInfo, address *model.Address, discounts []*model.DiscountInfo) (*goprices.TaxedMoney, *model.AppError) {
	if discounts == nil {
		discounts = []*model.DiscountInfo{}
	}
	calculatedCheckoutSubTotal, appErr := manager.CalculateCheckoutSubTotal(checkoutInfo, lines, address, discounts)
	if appErr != nil {
		return nil, appErr
	}

	calculatedCheckoutSubTotal, err := calculatedCheckoutSubTotal.Quantize(nil, goprices.Up)
	if err != nil {
		return nil, model.NewAppError("CheckoutSubTotal", app.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return calculatedCheckoutSubTotal, nil
}

// CalculateCheckoutTotalWithGiftcards
func (s *ServiceCheckout) CalculateCheckoutTotalWithGiftcards(manager interfaces.PluginManagerInterface, checkoutInfo model.CheckoutInfo, lines []*model.CheckoutLineInfo, address *model.Address, discounts []*model.DiscountInfo) (*goprices.TaxedMoney, *model.AppError) {
	checkoutTotal, appErr := s.CheckoutTotal(manager, checkoutInfo, lines, address, discounts)
	if appErr != nil {
		return nil, appErr
	}

	checkoutTotalGiftcardBalance, appErr := s.CheckoutTotalGiftCardsBalance(&checkoutInfo.Checkout)
	if appErr != nil {
		return nil, appErr
	}

	total, err := checkoutTotal.Sub(checkoutTotalGiftcardBalance)
	if err != nil {
		return nil, model.NewAppError("CalculateCheckoutTotalWithGiftcards", app.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	zeroTaxedMoney, _ := util.ZeroTaxedMoney(total.Currency)
	if zeroTaxedMoney.LessThan(total) {
		return total, nil
	}

	return zeroTaxedMoney, nil
}

// CheckoutTotal Return the total cost of the checkout.
//
// Total is a cost of all lines and shipping fees, minus checkout discounts,
// taxes included.
//
// It takes in account all plugins.
func (s *ServiceCheckout) CheckoutTotal(manager interfaces.PluginManagerInterface, checkoutInfo model.CheckoutInfo, lines []*model.CheckoutLineInfo, address *model.Address, discounts []*model.DiscountInfo) (*goprices.TaxedMoney, *model.AppError) {
	if discounts == nil {
		discounts = []*model.DiscountInfo{}
	}
	calculatedCheckoutTotal, appErr := manager.CalculateCheckoutTotal(checkoutInfo, lines, address, discounts)
	if appErr != nil {
		return nil, appErr
	}

	calculatedCheckoutTotal, err := calculatedCheckoutTotal.Quantize(nil, goprices.Up)
	if err != nil {
		return nil, model.NewAppError("CheckoutTotal", app.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return calculatedCheckoutTotal, nil
}

// CheckoutLineTotal Return the total price of provided line, taxes included.
//
// It takes in account all plugins.
func (s *ServiceCheckout) CheckoutLineTotal(manager interfaces.PluginManagerInterface, checkoutInfo model.CheckoutInfo, lines []*model.CheckoutLineInfo, checkoutLineInfo *model.CheckoutLineInfo, discounts []*model.DiscountInfo) (*goprices.TaxedMoney, *model.AppError) {
	address := checkoutInfo.ShippingAddress
	if address == nil {
		address = checkoutInfo.BillingAddress
	}

	if discounts == nil {
		discounts = []*model.DiscountInfo{}
	}

	calculatedLineTotal, appErr := manager.CalculateCheckoutLineTotal(checkoutInfo, lines, *checkoutLineInfo, address, discounts)
	if appErr != nil {
		return nil, appErr
	}

	calculatedLineTotal, err := calculatedLineTotal.Quantize(nil, goprices.Up)
	if err != nil {
		return nil, model.NewAppError("CheckoutLineTotal", app.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return calculatedLineTotal, nil
}

package checkout

import (
	"net/http"

	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app/plugin/interfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/util"
)

func (s *ServiceCheckout) CheckoutShippingPrice(manager interfaces.PluginManagerInterface, checkoutInfo model_helper.CheckoutInfo, lines model_helper.CheckoutLineInfos, address *model.Address, discounts []*model_helper.DiscountInfo) (*goprices.TaxedMoney, *model_helper.AppError) {
	calculatedCheckoutShipping, appErr := manager.CalculateCheckoutShipping(checkoutInfo, lines, address, discounts)
	if appErr != nil {
		return nil, appErr
	}

	calculatedCheckoutShipping, err := calculatedCheckoutShipping.Quantize(goprices.Up, -1)
	if err != nil {
		return nil, model_helper.NewAppError("CheckoutShippingPrice", model_helper.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return calculatedCheckoutShipping, nil
}

func (s *ServiceCheckout) CheckoutSubTotal(manager interfaces.PluginManagerInterface, checkoutInfo model_helper.CheckoutInfo, lines model_helper.CheckoutLineInfos, address *model.Address, discounts []*model_helper.DiscountInfo) (*goprices.TaxedMoney, *model_helper.AppError) {
	calculatedCheckoutSubTotal, appErr := manager.CalculateCheckoutSubTotal(checkoutInfo, lines, address, discounts)
	if appErr != nil {
		return nil, appErr
	}

	calculatedCheckoutSubTotal, err := calculatedCheckoutSubTotal.Quantize(goprices.Up, -1)
	if err != nil {
		return nil, model_helper.NewAppError("CheckoutSubTotal", model_helper.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return calculatedCheckoutSubTotal, nil
}

func (s *ServiceCheckout) CalculateCheckoutTotalWithGiftcards(manager interfaces.PluginManagerInterface, checkoutInfo model_helper.CheckoutInfo, lines model_helper.CheckoutLineInfos, address *model.Address, discounts []*model_helper.DiscountInfo) (*goprices.TaxedMoney, *model_helper.AppError) {
	checkoutTotal, appErr := s.CheckoutTotal(manager, checkoutInfo, lines, address, discounts)
	if appErr != nil {
		return nil, appErr
	}

	checkoutTotalGiftcardBalance, appErr := s.CheckoutTotalGiftCardsBalance(checkoutInfo.Checkout)
	if appErr != nil {
		return nil, appErr
	}

	total, err := checkoutTotal.Sub(checkoutTotalGiftcardBalance)
	if err != nil {
		return nil, model_helper.NewAppError("CalculateCheckoutTotalWithGiftcards", model_helper.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	zeroTaxedMoney, _ := util.ZeroTaxedMoney(total.GetCurrency())
	if zeroTaxedMoney.LessThan(*total) {
		return total, nil
	}

	return zeroTaxedMoney, nil
}

func (s *ServiceCheckout) CheckoutTotal(manager interfaces.PluginManagerInterface, checkoutInfo model_helper.CheckoutInfo, lines model_helper.CheckoutLineInfos, address *model.Address, discounts []*model_helper.DiscountInfo) (*goprices.TaxedMoney, *model_helper.AppError) {
	if discounts == nil {
		discounts = []*model_helper.DiscountInfo{}
	}
	calculatedCheckoutTotal, appErr := manager.CalculateCheckoutTotal(checkoutInfo, lines, address, discounts)
	if appErr != nil {
		return nil, appErr
	}

	calculatedCheckoutTotal, err := calculatedCheckoutTotal.Quantize(goprices.Up, -1)
	if err != nil {
		return nil, model_helper.NewAppError("CheckoutTotal", model_helper.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return calculatedCheckoutTotal, nil
}

func (s *ServiceCheckout) CheckoutLineTotal(manager interfaces.PluginManagerInterface, checkoutInfo model_helper.CheckoutInfo, lines model_helper.CheckoutLineInfos, checkoutLineInfo model_helper.CheckoutLineInfo, discounts []*model_helper.DiscountInfo) (*goprices.TaxedMoney, *model_helper.AppError) {
	address := checkoutInfo.ShippingAddress
	if address == nil {
		address = checkoutInfo.BillingAddress
	}

	if discounts == nil {
		discounts = []*model_helper.DiscountInfo{}
	}

	calculatedLineTotal, appErr := manager.CalculateCheckoutLineTotal(checkoutInfo, lines, checkoutLineInfo, address, discounts)
	if appErr != nil {
		return nil, appErr
	}

	calculatedLineTotal, err := calculatedLineTotal.Quantize(goprices.Up, -1)
	if err != nil {
		return nil, model_helper.NewAppError("CheckoutLineTotal", model_helper.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return calculatedLineTotal, nil
}

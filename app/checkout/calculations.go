package checkout

import (
	"net/http"

	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/checkout"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/modules/util"
)

// CheckoutShippingPrice Return checkout shipping price.
//
// It takes in account all plugins.
func (s *ServiceCheckout) CheckoutShippingPrice(manager interface{}, checkoutInfo *checkout.CheckoutInfo, lines []*checkout.CheckoutLineInfo, address *account.Address, discounts []*product_and_discount.DiscountInfo) (*goprices.TaxedMoney, *model.AppError) {
	panic("not implemented")
}

// CheckoutSubTotal Return the total cost of all the checkout lines, taxes included.
//
// It takes in account all plugins.
func (s *ServiceCheckout) CheckoutSubTotal(manager interface{}, checkoutInfo *checkout.CheckoutInfo, lines []*checkout.CheckoutLineInfo, address *account.Address, discounts []*product_and_discount.DiscountInfo) (*goprices.TaxedMoney, *model.AppError) {
	panic("not implemented")
}

// CalculateCheckoutTotalWithGiftcards
func (s *ServiceCheckout) CalculateCheckoutTotalWithGiftcards(manager interface{}, checkoutInfo *checkout.CheckoutInfo, lines []*checkout.CheckoutLineInfo, address *account.Address, discounts []*product_and_discount.DiscountInfo) (*goprices.TaxedMoney, *model.AppError) {
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
	if less, err := zeroTaxedMoney.LessThan(total); less && err == nil {
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
func (s *ServiceCheckout) CheckoutTotal(manager interface{}, checkoutInfo *checkout.CheckoutInfo, lines []*checkout.CheckoutLineInfo, address *account.Address, discounts []*product_and_discount.DiscountInfo) (*goprices.TaxedMoney, *model.AppError) {
	panic("not implemented")
}

// CheckoutLineTotal Return the total price of provided line, taxes included.
//
// It takes in account all plugins.
func (s *ServiceCheckout) CheckoutLineTotal(manager interface{}, checkoutInfo *checkout.CheckoutInfo, lines []*checkout.CheckoutLineInfo, checkoutLineInfo *checkout.CheckoutLineInfo, discounts []*product_and_discount.DiscountInfo) (*goprices.TaxedMoney, *model.AppError) {
	panic("not implemented")
}

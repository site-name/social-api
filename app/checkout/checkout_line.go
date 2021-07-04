package checkout

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/checkout"
	"github.com/sitename/sitename/store"
)

func (a *AppCheckout) CheckoutLineShippingRequired(checkoutLine *checkout.CheckoutLine) (bool, *model.AppError) {
	panic("not implt")
}

func (a *AppCheckout) CheckoutLinesByCheckoutID(checkoutID string) ([]*checkout.CheckoutLine, *model.AppError) {
	lines, err := a.Srv().Store.CheckoutLine().CheckoutLinesByCheckoutID(checkoutID)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("CheckoutLinesByCheckoutID", "app.checkout.checkout_lines_by_checkout.app_error", err)
	}

	return lines, nil
}

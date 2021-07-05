package checkout

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/checkout"
	"github.com/sitename/sitename/store"
)

// func (a *AppCheckout) CheckoutLineShippingRequired(checkoutLine *checkout.CheckoutLine) (bool, *model.AppError) {
// 	productVariant, appErr := a.ProductApp().ProductVariantById(checkoutLine.VariantID)
// 	if appErr != nil {
// 		return false, appErr
// 	}

// }

func (a *AppCheckout) CheckoutLinesByCheckoutID(checkoutToken string) ([]*checkout.CheckoutLine, *model.AppError) {
	lines, err := a.Srv().Store.CheckoutLine().CheckoutLinesByCheckoutID(checkoutToken)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("CheckoutLinesByCheckoutID", "app.checkout.checkout_lines_by_checkout.app_error", err)
	}

	return lines, nil
}

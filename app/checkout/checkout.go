package checkout

import (
	"context"
	"net/http"

	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/checkout"
	"github.com/sitename/sitename/store"
)

type AppCheckout struct {
	app.AppIface
}

func init() {
	app.RegisterCheckoutApp(func(a app.AppIface) sub_app_iface.CheckoutApp {
		return &AppCheckout{a}
	})
}

func (a *AppCheckout) CheckoutbyToken(checkoutToken string) (*checkout.Checkout, *model.AppError) {
	checkout, err := a.Srv().Store.Checkout().Get(checkoutToken)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("CheckoutById", "app.checkout.missing_checkout.app_error", err)
	}

	return checkout, nil
}

func (a *AppCheckout) GetCustomerEmail(ckout *checkout.Checkout) (string, *model.AppError) {
	if ckout.UserID != nil {
		user, appErr := a.AccountApp().UserById(context.Background(), *ckout.UserID)
		if appErr != nil {
			return "", appErr
		}
		return user.Email, nil
	}
	return ckout.Email, nil
}

func (a *AppCheckout) CheckoutShippingRequired(checkoutToken string) (bool, *model.AppError) {
	/*
					checkout
					|      |
		...<--|		   |--> checkoutLine <-- productVariant <-- product <-- productType
																							|												     |
													 ...checkoutLine <--|              ...product <--|
	*/

	productTypes, appErr := a.ProductApp().ProductTypesByCheckoutToken(checkoutToken)
	if appErr != nil {
		// if product types not found for checkout:
		if appErr.StatusCode == http.StatusNotFound {
			return false, nil
		}
		return false, appErr
	}

	for _, prdType := range productTypes {
		if prdType.IsShippingRequired != nil && *prdType.IsShippingRequired {
			return true, nil
		}
	}

	return false, nil
}

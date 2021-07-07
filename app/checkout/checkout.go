package checkout

import (
	"context"
	"net/http"
	"strings"

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

func (a *AppCheckout) CheckoutsByUser(userID string) (*checkout.Checkout, *model.AppError) {
	panic("not implt")
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

func (a *AppCheckout) CheckoutSetCountry(ckout *checkout.Checkout, newCountryCode string) *model.AppError {
	// no need to validate country code here, since checkout.IsValid() does that
	countryCode := strings.ToUpper(strings.TrimSpace(newCountryCode))
	ckout.Country = countryCode
	_, appErr := a.UpdateCheckout(ckout)
	return appErr
}

func (a *AppCheckout) CheckoutCountry(ckout *checkout.Checkout) (string, *model.AppError) {
	addressID := ckout.ShippingAddressID
	if ckout.ShippingAddressID == nil {
		addressID = ckout.BillingAddressID
	}

	if addressID == nil {
		return ckout.Country, nil
	}

	address, appErr := a.AccountApp().AddressById(*addressID)
	// ignore this error even when the lookup fail
	if appErr != nil || address == nil || strings.TrimSpace(address.Country) == "" {
		return ckout.Country, nil
	}

	countryCode := strings.TrimSpace(address.Country)
	if countryCode != ckout.Country {
		// set new country code for checkout:
		appErr := a.CheckoutSetCountry(ckout, countryCode)
		if appErr != nil {
			return "", appErr
		}
	}

	return countryCode, nil
}

func (a *AppCheckout) UpdateCheckout(ckout *checkout.Checkout) (*checkout.Checkout, *model.AppError) {
	newCkout, err := a.Srv().Store.Checkout().Update(ckout)
	if err != nil {
		if appErr, ok := err.(*model.AppError); ok {
			return nil, appErr
		}
		return nil, model.NewAppError(
			"UpdateCheckout",
			"app.checkout.checkout_update_failed.app_error",
			nil, err.Error(),
			http.StatusInternalServerError,
		)
	}

	return newCkout, nil
}

package checkout

import (
	"context"
	"net/http"
	"sort"
	"strings"

	"github.com/shopspring/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/checkout"
	"github.com/sitename/sitename/model/payment"
	"github.com/sitename/sitename/store"
)

type AppCheckout struct {
	app app.AppIface
}

const (
	CheckoutMissingAppErrorId = "app.checkout.missing_checkout.app_error"
)

func init() {
	app.RegisterCheckoutApp(func(a app.AppIface) sub_app_iface.CheckoutApp {
		return &AppCheckout{a}
	})
}

func (a *AppCheckout) CheckoutsByUser(userID string, channelActive bool) ([]*checkout.Checkout, *model.AppError) {
	checkouts, err := a.app.Srv().Store.Checkout().CheckoutsByUserID(userID, channelActive)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("CheckoutsByUser", "app.checkout.checkout_by_user_missing.app_error", err)
	}
	return checkouts, nil
}

func (a *AppCheckout) GetUserCheckout(userID string) (*checkout.Checkout, *model.AppError) {
	checkouts, appErr := a.CheckoutsByUser(userID, true)
	if appErr != nil {
		return nil, appErr
	}
	// sort checkouts by creation time ascending
	sort.Slice(checkouts, func(i, j int) bool {
		return checkouts[i].CreateAt < checkouts[j].CreateAt
	})
	// last checkout is the most recent made by user
	return checkouts[len(checkouts)-1], nil
}

func (a *AppCheckout) CheckoutbyToken(checkoutToken string) (*checkout.Checkout, *model.AppError) {
	checkout, err := a.app.Srv().Store.Checkout().Get(checkoutToken)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("CheckoutById", CheckoutMissingAppErrorId, err)
	}

	return checkout, nil
}

func (a *AppCheckout) GetCustomerEmail(ckout *checkout.Checkout) (string, *model.AppError) {
	if ckout.UserID != nil {
		user, appErr := a.app.AccountApp().UserById(context.Background(), *ckout.UserID)
		if appErr != nil {
			return "", appErr
		}
		return user.Email, nil
	}
	return ckout.Email, nil
}

// CheckoutShippingRequired checks if given checkout require shipping
func (a *AppCheckout) CheckoutShippingRequired(checkoutToken string) (bool, *model.AppError) {
	/*
					checkout
					|      |
		...<--|		   |--> checkoutLine <-- productVariant <-- product <-- productType
																							|												     |
													 ...checkoutLine <--|              ...product <--|
	*/

	productTypes, appErr := a.app.ProductApp().ProductTypesByCheckoutToken(checkoutToken)
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
	ckout.Country = strings.ToUpper(strings.TrimSpace(newCountryCode))
	_, appErr := a.UpsertCheckout(ckout)
	return appErr
}

// UpsertCheckout saves/updates given checkout
func (a *AppCheckout) UpsertCheckout(ckout *checkout.Checkout) (*checkout.Checkout, *model.AppError) {
	ckout, err := a.app.Srv().Store.Checkout().Upsert(ckout)
	if err != nil {
		if appErr, ok := err.(*model.AppError); ok {
			return nil, appErr
		}

		statusCode := http.StatusInternalServerError
		var errID string = "app.checkout.checkout_upsert_error.app_error"

		if _, ok := err.(*store.ErrNotFound); ok { // this error caused by Update
			errID = CheckoutMissingAppErrorId
			statusCode = http.StatusNotFound
		}
		return nil, model.NewAppError("UpsertCheckout", errID, nil, err.Error(), statusCode)
	}

	return ckout, nil
}

func (a *AppCheckout) CheckoutCountry(ckout *checkout.Checkout) (string, *model.AppError) {
	addressID := ckout.ShippingAddressID
	if ckout.ShippingAddressID == nil {
		addressID = ckout.BillingAddressID
	}

	if addressID == nil {
		return ckout.Country, nil
	}

	address, appErr := a.app.AccountApp().AddressById(*addressID)
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

func (a *AppCheckout) CheckoutTotalGiftCardsBalance(checkout *checkout.Checkout) (*goprices.Money, *model.AppError) {
	gcs, appErr := a.app.GiftcardApp().GiftcardsByCheckout(checkout.Token)
	if appErr != nil {
		return nil, appErr
	}

	balanceAmount := decimal.Zero
	for _, gc := range gcs {
		if gc.CurrentBalanceAmount != nil {
			balanceAmount = balanceAmount.Add(*gc.CurrentBalanceAmount)
		}
	}

	return &goprices.Money{
		Amount:   &balanceAmount,
		Currency: checkout.Currency,
	}, nil
}

func (a *AppCheckout) CheckoutLineWithVariant(checkout *checkout.Checkout, productVariantID string) (*checkout.CheckoutLine, *model.AppError) {
	checkoutLines, appErr := a.CheckoutLinesByCheckoutToken(checkout.Token)
	if appErr != nil {
		// in case checkout has no checkout lines:
		if appErr.StatusCode == http.StatusNotFound {
			return nil, nil
		}
		return nil, appErr
	}

	for _, line := range checkoutLines {
		if line.VariantID == productVariantID {
			return line, nil
		}
	}

	return nil, nil
}

func (a *AppCheckout) CheckoutLastActivePayment(checkout *checkout.Checkout) (*payment.Payment, *model.AppError) {
	payments, appErr := a.app.PaymentApp().GetAllPaymentsByCheckout(checkout.Token)
	if appErr != nil {
		if appErr.StatusCode == http.StatusNotFound {
			return nil, nil
		}
		return nil, appErr
	}

	// find latest payment by comparing their creation time
	var latestPayment *payment.Payment
	for _, pm := range payments {
		if pm.IsActive && (latestPayment == nil || latestPayment.CreateAt < pm.CreateAt) {
			latestPayment = pm
		}
	}

	return latestPayment, nil
}

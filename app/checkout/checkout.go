/*
	NOTE: This package is initialized during server startup (modules/imports does that)
	so the init() function get the chance to register a function to create `ServiceAccount`
*/
package checkout

import (
	"context"
	"net/http"
	"strings"
	"sync"

	"github.com/mattermost/gorp"
	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/checkout"
	"github.com/sitename/sitename/model/payment"
	"github.com/sitename/sitename/modules/measurement"
	"github.com/sitename/sitename/store"
)

type ServiceCheckout struct {
	srv   *app.Server
	wg    sync.WaitGroup
	mutex sync.Mutex
}

func init() {
	app.RegisterCheckoutService(func(s *app.Server) (sub_app_iface.CheckoutService, error) {
		return &ServiceCheckout{
			srv: s,
		}, nil
	})
}

const (
	CheckoutMissingAppErrorId = "app.checkout.missing_checkout.app_error"
)

// CheckoutByOption returns a checkout filtered by given option
func (a *ServiceCheckout) CheckoutByOption(option *checkout.CheckoutFilterOption) (*checkout.Checkout, *model.AppError) {
	chekout, err := a.srv.Store.Checkout().GetByOption(option)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("CheckoutbyOption", "app.checkout.error_finding_checkout_by_option.app_error", err)
	}

	return chekout, nil
}

// CheckoutsByOption returns a list of checkouts, filtered by given option
func (a *ServiceCheckout) CheckoutsByOption(option *checkout.CheckoutFilterOption) ([]*checkout.Checkout, *model.AppError) {
	checkouts, err := a.srv.Store.Checkout().FilterByOption(option)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("CheckoutsByOption", "app.checkout_error_finding_checkouts_by_option.app_error", err)
	}
	return checkouts, nil
}

// GetCustomerEmail returns checkout's user's email
func (a *ServiceCheckout) GetCustomerEmail(ckout *checkout.Checkout) (string, *model.AppError) {
	if ckout.UserID != nil {
		user, appErr := a.srv.AccountService().UserById(context.Background(), *ckout.UserID)
		if appErr != nil {
			if appErr.StatusCode == http.StatusNotFound {
				return ckout.Email, nil
			}
			return "", appErr // returns system caused error
		}
		return user.Email, nil
	}
	return ckout.Email, nil
}

// CheckoutShippingRequired checks if given checkout require shipping
func (a *ServiceCheckout) CheckoutShippingRequired(checkoutToken string) (bool, *model.AppError) {
	/*
					checkout
					|      |
		...<--|		   |--> checkoutLine <-- productVariant <-- product <-- productType
																							|												     |
													 ...checkoutLine <--|              ...product <--|
	*/

	productTypes, appErr := a.srv.ProductService().ProductTypesByCheckoutToken(checkoutToken)
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

func (a *ServiceCheckout) CheckoutSetCountry(ckout *checkout.Checkout, newCountryCode string) *model.AppError {
	// no need to validate country code here, since checkout.IsValid() does that
	ckout.Country = strings.ToUpper(strings.TrimSpace(newCountryCode))
	_, appErr := a.UpsertCheckout(ckout)
	return appErr
}

// UpsertCheckout saves/updates given checkout
func (a *ServiceCheckout) UpsertCheckout(ckout *checkout.Checkout) (*checkout.Checkout, *model.AppError) {
	ckout, err := a.srv.Store.Checkout().Upsert(ckout)
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

func (a *ServiceCheckout) CheckoutCountry(ckout *checkout.Checkout) (string, *model.AppError) {
	addressID := ckout.ShippingAddressID
	if addressID == nil {
		addressID = ckout.BillingAddressID
	}

	if addressID == nil {
		return ckout.Country, nil
	}

	address, appErr := a.srv.AccountService().AddressById(*addressID)
	if appErr != nil {
		// return immediately if the error is caused by system
		if appErr.StatusCode == http.StatusInternalServerError {
			return "", appErr
		}
		if address == nil || strings.TrimSpace(address.Country) == "" {
			return ckout.Country, nil
		}
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

// CheckoutTotalGiftCardsBalance Return the total balance of the gift cards assigned to the checkout
func (a *ServiceCheckout) CheckoutTotalGiftCardsBalance(checkout *checkout.Checkout) (*goprices.Money, *model.AppError) {
	giftcards, appErr := a.srv.GiftcardService().GiftcardsByCheckout(checkout.Token)
	if appErr != nil {
		return nil, appErr
	}

	balanceAmount := decimal.Zero
	for _, giftcard := range giftcards {
		if giftcard.CurrentBalanceAmount != nil {
			balanceAmount = balanceAmount.Add(*giftcard.CurrentBalanceAmount)
		}
	}

	return &goprices.Money{
		Amount:   &balanceAmount,
		Currency: checkout.Currency,
	}, nil
}

func (a *ServiceCheckout) CheckoutLineWithVariant(checkout *checkout.Checkout, productVariantID string) (*checkout.CheckoutLine, *model.AppError) {
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

// CheckoutLastActivePayment returns the most recent payment made for given checkout
func (a *ServiceCheckout) CheckoutLastActivePayment(checkout *checkout.Checkout) (*payment.Payment, *model.AppError) {
	payments, appErr := a.srv.PaymentService().PaymentsByOption(&payment.PaymentFilterOption{
		CheckoutToken: checkout.Token,
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusNotFound {
			return nil, nil
		}
		return nil, appErr
	}

	// find latest payment by comparing their creation time
	var latestPayment payment.Payment
	for _, payMent := range payments {
		if *payMent.IsActive && (latestPayment.Id == "" || latestPayment.CreateAt < payMent.CreateAt) {
			latestPayment = *payMent
		}
	}

	return &latestPayment, nil
}

// CheckoutTotalWeight calculate total weight for given checkout lines (these lines belong to a single checkout)
func (a *ServiceCheckout) CheckoutTotalWeight(checkoutLineInfos []*checkout.CheckoutLineInfo) (*measurement.Weight, *model.AppError) {
	checkoutLineIDs := []string{}
	for _, lineInfo := range checkoutLineInfos {
		if !model.IsValidId(lineInfo.Line.Id) {
			checkoutLineIDs = append(checkoutLineIDs, lineInfo.Line.Id)
		}
	}

	totalWeight, err := a.srv.Store.CheckoutLine().TotalWeightForCheckoutLines(checkoutLineIDs)
	if err != nil {
		status := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			status = http.StatusNotFound
		}
		return nil, model.NewAppError("CheckoutTotalWeight", "app.checkout.checkout_total_weight.app_error", nil, err.Error(), status)
	}

	return totalWeight, nil
}

// DeleteCheckoutsByOption tells store to delete checkout(s) rows, filtered using given option
func (s *ServiceCheckout) DeleteCheckoutsByOption(transaction *gorp.Transaction, option *checkout.CheckoutFilterOption) *model.AppError {
	err := s.srv.SqlStore.Checkout().DeleteCheckoutsByOption(transaction, option)
	if err != nil {
		return model.NewAppError("DeleteCheckoutsByOption", "app.checkout.error_deleting_checkouts_by_option.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}

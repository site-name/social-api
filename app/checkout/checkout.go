/*
NOTE: This package is initialized during server startup (modules/imports does that)
so the init() function get the chance to register a function to create `ServiceAccount`
*/
package checkout

import (
	"context"
	"net/http"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/measurement"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type ServiceCheckout struct {
	srv *app.Server
}

func init() {
	app.RegisterService(func(s *app.Server) error {
		s.Checkout = &ServiceCheckout{s}
		return nil
	})
}

// CheckoutByOption returns a checkout filtered by given option
func (a *ServiceCheckout) CheckoutByOption(option *model.CheckoutFilterOption) (*model.Checkout, *model.AppError) {
	chekout, err := a.srv.Store.Checkout().GetByOption(option)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}
		return nil, model.NewAppError("CheckoutbyOption", "app.checkout.error_finding_checkout_by_option.app_error", nil, err.Error(), statusCode)
	}

	return chekout, nil
}

// CheckoutsByOption returns a list of checkouts, filtered by given option
func (a *ServiceCheckout) CheckoutsByOption(option *model.CheckoutFilterOption) ([]*model.Checkout, *model.AppError) {
	checkouts, err := a.srv.Store.Checkout().FilterByOption(option)
	var (
		statusCode int
		errMsg     string
	)
	if err != nil {
		statusCode = http.StatusInternalServerError
		errMsg = err.Error()
	} else if len(checkouts) == 0 {
		statusCode = http.StatusNotFound
	}

	if statusCode != 0 {
		return nil, model.NewAppError("CheckoutsByOption", "app.checkout.error_finding_checkouts.app_error", nil, errMsg, statusCode)
	}

	return checkouts, nil
}

// GetCustomerEmail returns checkout's user's email
func (a *ServiceCheckout) GetCustomerEmail(ckout *model.Checkout) (string, *model.AppError) {
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

func (a *ServiceCheckout) CheckoutSetCountry(ckout *model.Checkout, newCountryCode model.CountryCode) *model.AppError {
	// no need to validate country code here, since checkout.IsValid() does that
	ckout.Country = newCountryCode
	_, appErr := a.UpsertCheckouts(nil, []*model.Checkout{ckout})
	return appErr
}

// UpsertCheckout saves/updates given checkout
func (a *ServiceCheckout) UpsertCheckouts(transaction *gorm.DB, checkouts []*model.Checkout) ([]*model.Checkout, *model.AppError) {
	checkouts, err := a.srv.Store.Checkout().Upsert(transaction, checkouts)
	if err != nil {
		if appErr, ok := err.(*model.AppError); ok {
			return nil, appErr
		}

		statusCode := http.StatusInternalServerError
		var errID string = "app.checkout.checkout_upsert_error.app_error"

		if _, ok := err.(*store.ErrNotFound); ok { // this error caused by Update
			errID = "app.checkout.missing_checkout.app_error"
			statusCode = http.StatusNotFound
		}
		return nil, model.NewAppError("UpsertCheckout", errID, nil, err.Error(), statusCode)
	}

	return checkouts, nil
}

func (a *ServiceCheckout) CheckoutCountry(ckout *model.Checkout) (model.CountryCode, *model.AppError) {
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
		if address == nil || address.Country == "" {
			return ckout.Country, nil
		}
	}

	countryCode := address.Country
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
func (a *ServiceCheckout) CheckoutTotalGiftCardsBalance(checkOut *model.Checkout) (*goprices.Money, *model.AppError) {
	giftcards, appErr := a.srv.GiftcardService().GiftcardsByOption(&model.GiftCardFilterOption{
		CheckoutToken: squirrel.Eq{model.GiftcardCheckoutTableName + ".CheckoutID": checkOut.Token},
		Conditions: squirrel.And{
			squirrel.Or{
				squirrel.Eq{model.GiftcardTableName + ".ExpiryDate": nil},
				squirrel.GtOrEq{model.GiftcardTableName + ".ExpiryDate": util.StartOfDay(time.Now().UTC())},
			},
			squirrel.Eq{model.GiftcardTableName + ".IsActive": true},
		},
	})
	if appErr != nil {
		return nil, appErr
	}

	balanceAmount := decimal.Zero
	for _, giftcard := range giftcards {
		if giftcard != nil && giftcard.CurrentBalanceAmount != nil {
			balanceAmount = balanceAmount.Add(*giftcard.CurrentBalanceAmount)
		}
	}

	return &goprices.Money{
		Amount:   balanceAmount,
		Currency: checkOut.Currency,
	}, nil
}

func (a *ServiceCheckout) CheckoutLineWithVariant(checkout *model.Checkout, productVariantID string) (*model.CheckoutLine, *model.AppError) {
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
func (a *ServiceCheckout) CheckoutLastActivePayment(checkout *model.Checkout) (*model.Payment, *model.AppError) {
	payments, appErr := a.srv.PaymentService().PaymentsByOption(&model.PaymentFilterOption{
		Conditions: squirrel.Eq{model.PaymentTableName + ".CheckoutID": checkout.Token},
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusNotFound {
			return nil, nil
		}
		return nil, appErr
	}

	// find latest payment by comparing their creation time
	var latestPayment model.Payment
	for _, payMent := range payments {
		if *payMent.IsActive && (latestPayment.Id == "" || latestPayment.CreateAt < payMent.CreateAt) {
			latestPayment = *payMent
		}
	}

	return &latestPayment, nil
}

// CheckoutTotalWeight calculate total weight for given checkout lines (these lines belong to a single checkout)
func (a *ServiceCheckout) CheckoutTotalWeight(checkoutLineInfos []*model.CheckoutLineInfo) (*measurement.Weight, *model.AppError) {
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
func (s *ServiceCheckout) DeleteCheckoutsByOption(transaction *gorm.DB, option *model.CheckoutFilterOption) *model.AppError {
	err := s.srv.Store.Checkout().DeleteCheckoutsByOption(transaction, option)
	if err != nil {
		return model.NewAppError("DeleteCheckoutsByOption", "app.checkout.error_deleting_checkouts_by_option.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}

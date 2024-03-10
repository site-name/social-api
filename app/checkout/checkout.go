/*
NOTE: This package is initialized during server startup (modules/imports does that)
so the init() function get the chance to register a function to create `ServiceAccount`
*/
package checkout

import (
	"context"
	"net/http"

	"github.com/mattermost/squirrel"
	"github.com/samber/lo"
	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/measurement"
	"github.com/sitename/sitename/modules/model_types"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
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

func (a *ServiceCheckout) CheckoutByOption(option model_helper.CheckoutFilterOptions) (*model.Checkout, *model_helper.AppError) {
	chekout, err := a.srv.Store.Checkout().GetByOption(option)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}
		return nil, model_helper.NewAppError("CheckoutbyOption", "app.checkout.error_finding_checkout_by_option.app_error", nil, err.Error(), statusCode)
	}

	return chekout, nil
}

func (a *ServiceCheckout) CheckoutsByOption(option model_helper.CheckoutFilterOptions) (model.CheckoutSlice, *model_helper.AppError) {
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
		return nil, model_helper.NewAppError("CheckoutsByOption", "app.checkout.error_finding_checkouts.app_error", nil, errMsg, statusCode)
	}

	return checkouts, nil
}

func (a *ServiceCheckout) GetCustomerEmail(checkout model.Checkout) (string, *model_helper.AppError) {
	if !checkout.UserID.IsNil() {
		user, appErr := a.srv.Account.UserById(context.Background(), *checkout.UserID.String)
		if appErr != nil {
			if appErr.StatusCode == http.StatusNotFound {
				return checkout.Email, nil
			}
			return "", appErr // returns system caused error
		}
		return user.Email, nil
	}
	return checkout.Email, nil
}

// TODO: check if we need this method. Since we don't use product_type anymore
func (a *ServiceCheckout) CheckoutShippingRequired(checkoutToken string) (bool, *model_helper.AppError) {
	productTypes, appErr := a.srv.Product.ProductTypesByCheckoutToken(checkoutToken)
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

func (a *ServiceCheckout) CheckoutSetCountry(checkout model.Checkout, newCountryCode model.CountryCode) *model_helper.AppError {
	checkout.Country = newCountryCode
	_, appErr := a.UpsertCheckouts(nil, model.CheckoutSlice{&checkout})
	return appErr
}

func (a *ServiceCheckout) UpsertCheckouts(transaction boil.ContextTransactor, checkouts model.CheckoutSlice) (model.CheckoutSlice, *model_helper.AppError) {
	checkouts, err := a.srv.Store.Checkout().Upsert(transaction, checkouts)
	if err != nil {
		if appErr, ok := err.(*model_helper.AppError); ok {
			return nil, appErr
		}

		statusCode := http.StatusInternalServerError
		var errID string = "app.checkout.checkout_upsert_error.app_error"

		if _, ok := err.(*store.ErrNotFound); ok { // this error caused by Update
			errID = "app.checkout.missing_checkout.app_error"
			statusCode = http.StatusNotFound
		}
		return nil, model_helper.NewAppError("UpsertCheckout", errID, nil, err.Error(), statusCode)
	}

	return checkouts, nil
}

func (a *ServiceCheckout) CheckoutCountry(checkout model.Checkout) (model.CountryCode, *model_helper.AppError) {
	addressID := checkout.ShippingAddressID
	if addressID.IsNil() {
		addressID = checkout.BillingAddressID
	}

	if addressID.IsNil() {
		return checkout.Country, nil
	}

	address, appErr := a.srv.Account.AddressById(*addressID.String)
	if appErr != nil {
		// return immediately if the error is caused by system
		if appErr.StatusCode == http.StatusInternalServerError {
			return "", appErr
		}
		if address == nil || address.Country == "" {
			return checkout.Country, nil
		}
	}

	countryCode := address.Country
	if countryCode != checkout.Country {
		// set new country code for checkout:
		appErr := a.CheckoutSetCountry(checkout, countryCode)
		if appErr != nil {
			return "", appErr
		}
	}

	return countryCode, nil
}

// CheckoutTotalGiftCardsBalance Return the total balance of the gift cards assigned to the checkout
func (a *ServiceCheckout) CheckoutTotalGiftCardsBalance(checkout model.Checkout) (*goprices.Money, *model_helper.AppError) {
	_, giftcards, appErr := a.srv.Giftcard.GiftcardsByOption(model_helper.GiftcardFilterOption{
		CheckoutToken: model.GiftcardCheckoutWhere.CheckoutID.EQ(checkout.Token),
		CommonQueryOptions: model_helper.NewCommonQueryOptions(
			model_helper.And{
				squirrel.Or{
					squirrel.Eq{model.GiftcardTableColumns.ExpiryDate: nil},
					squirrel.GtOrEq{model.GiftcardTableColumns.ExpiryDate: util.StartOfDay(model_helper.GetTimeUTCNow())},
				},
				squirrel.Eq{model.GiftcardTableColumns.IsActive: true},
			},
		),
	})
	if appErr != nil {
		return nil, appErr
	}

	balanceAmount := decimal.Zero
	for _, giftcard := range giftcards {
		if giftcard != nil && !giftcard.CurrentBalanceAmount.IsNil() {
			balanceAmount = balanceAmount.Add(*giftcard.CurrentBalanceAmount.Decimal)
		}
	}

	return &goprices.Money{
		Amount:   balanceAmount,
		Currency: checkout.Currency.String(),
	}, nil
}

func (a *ServiceCheckout) CheckoutLineWithVariant(checkout model.Checkout, productVariantID string) (*model.CheckoutLine, *model_helper.AppError) {
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

func (a *ServiceCheckout) CheckoutLastActivePayment(checkout model.Checkout) (*model.Payment, *model_helper.AppError) {
	payments, appErr := a.srv.Payment.PaymentsByOption(model_helper.PaymentFilterOptions{
		CommonQueryOptions: model_helper.NewCommonQueryOptions(
			model.PaymentWhere.CheckoutID.EQ(model_types.NewNullString(checkout.Token)),
			qm.OrderBy(model.PaymentColumns.CreatedAt+" "+string(model_helper.DESC)),
			qm.Limit(1),
		),
	})
	if appErr != nil {
		return nil, appErr
	}
	if len(payments) == 0 {
		return nil, nil
	}

	return payments[0], nil
}

// CheckoutTotalWeight calculate total weight for given checkout lines (these lines belong to a single checkout)
func (a *ServiceCheckout) CheckoutTotalWeight(checkoutLineInfos model_helper.CheckoutLineInfos) (*measurement.Weight, *model_helper.AppError) {
	checkoutLineIDs := []string{}
	for _, lineInfo := range checkoutLineInfos {
		if !model_helper.IsValidId(lineInfo.Line.ID) {
			checkoutLineIDs = append(checkoutLineIDs, lineInfo.Line.ID)
		}
	}

	totalWeight, err := a.srv.Store.CheckoutLine().TotalWeightForCheckoutLines(checkoutLineIDs)
	if err != nil {
		status := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			status = http.StatusNotFound
		}
		return nil, model_helper.NewAppError("CheckoutTotalWeight", "app.checkout.checkout_total_weight.app_error", nil, err.Error(), status)
	}

	return totalWeight, nil
}

func (s *ServiceCheckout) DeleteCheckoutsByOption(transaction boil.ContextTransactor, option model_helper.CheckoutFilterOptions) *model_helper.AppError {
	checkouts, appErr := s.CheckoutsByOption(option)
	if appErr != nil {
		return appErr
	}

	checkoutIDs := lo.Map(checkouts, func(item *model.Checkout, _ int) string { return item.Token })
	err := s.srv.Store.Checkout().Delete(transaction, checkoutIDs)
	if err != nil {
		return model_helper.NewAppError("DeleteCheckoutsByOption", "app.checkout.error_deleting_checkouts_by_option.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}

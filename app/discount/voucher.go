package discount

import (
	"net/http"

	"github.com/shopspring/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
)

// WrapperFunc number of `args` must be 1 or 2
//
//  if len(args) == 1 {
//		args[0].(type) == *Money || *MoneyRange || *TaxedMoney || *TaxedMoneyRange
//  }
//  if len(args) == 2 {
//		(args[0].(type) == *Money || *MoneyRange || *TaxedMoney || *TaxedMoneyRange) && args[0].(type) == bool
//  }
type WrapperFunc func(args ...interface{}) (interface{}, error)

func decorator(preValue interface{}) WrapperFunc {
	return func(args ...interface{}) (interface{}, error) {
		// validating number of args
		if l := len(args); l < 1 || l > 2 {
			return nil, model.NewAppError("app.Discount.decorator", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "args"}, "you must provide either 1 or 2 arguments", http.StatusBadRequest)
		}

		if len(args) == 1 { // fixed discount
			discount := preValue.(*goprices.Money)
			return goprices.FixedDiscount(args[0], discount)
		}
		return goprices.PercentageDiscount(args[0], preValue, args[1].(bool))
	}
}

func (a *AppDiscount) GetVoucherDiscount(voucher *product_and_discount.Voucher, channelID string) (WrapperFunc, *model.AppError) {
	voucherChannelListings, appErr := a.VoucherChannelListingsByVoucherAndChannel(voucher.Id, channelID)
	if appErr != nil {
		return nil, appErr
	}

	firstListing := voucherChannelListings[0]
	if firstListing == nil {
		return nil, model.NewAppError("VoucherChannelListingsByVoucherAndChannel", "app.discount.voucher_not_assigned_to_channel.app_error", nil, "", http.StatusNotAcceptable)
	}

	if voucher.DiscountValueType == product_and_discount.FIXED {
		discountAmount, err := goprices.NewMoney(firstListing.DiscountValue, firstListing.Currency)
		if err != nil {
			return nil, model.NewAppError("VoucherChannelListingsByVoucherAndChannel", app.NewMoneyCreationAppErrorID, nil, err.Error(), http.StatusInternalServerError)
		}
		return decorator(discountAmount), nil
	}
	return decorator(firstListing.DiscountValue), nil
}

func (a *AppDiscount) GetDiscountAmountFor(voucher *product_and_discount.Voucher, price *goprices.Money, channelID string) (*goprices.Money, *model.AppError) {
	discountCalcuFunc, appErr := a.GetVoucherDiscount(voucher, channelID)
	if appErr != nil {
		return nil, appErr
	}

	afterDiscount, err := discountCalcuFunc(price)
	if err != nil {
		// this error maybe caused by user. But we tomporarily set status code to 500
		return nil, model.NewAppError("GetDiscountAmountFor", "app.discount.error_calculating_discount.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	if afterDiscount.(*goprices.Money).Amount.LessThan(decimal.Zero) {
		return price, nil
	}
	sub, err := price.Sub(afterDiscount.(*goprices.Money))
	if err != nil {
		return nil, model.NewAppError("GetDiscountAmountFor", "app.discount.error_subtract_money.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return sub, nil
}

// ValidateMinSpent validates if the order cost at least a specific amount of money
func (a *AppDiscount) ValidateMinSpent(voucher *product_and_discount.Voucher, value *goprices.TaxedMoney, channelID string) *model.AppError {
	ownerShopOfVoucher, appErr := a.ShopApp().ShopById(voucher.ShopID)
	if appErr != nil {
		return appErr
	}

	money := value.Net
	if ownerShopOfVoucher.DisplayGrossPrices != nil && *ownerShopOfVoucher.DisplayGrossPrices {
		money = value.Gross
	}

	voucherChannelListings, appErr := a.VoucherChannelListingsByVoucherAndChannel(voucher.Id, channelID)
	if appErr != nil {
		return appErr
	}

	firstVoucherChannelListing := voucherChannelListings[0]
	if firstVoucherChannelListing.MinSpent != nil {
		if less, err := money.LessThan(firstVoucherChannelListing.MinSpent); less && err == nil {
			return model.NewAppError("ValidateMinSpent", "app.discount.voucher_not_applicable_for_cost_below.app_error", map[string]interface{}{"MinMoney": nil}, "", http.StatusNotAcceptable)
		}
	}

	return nil
}

// ValidateOncePerCustomer checks to make sure each customer has ONLY 1 time usage with 1 voucher
func (a *AppDiscount) ValidateOncePerCustomer(voucher *product_and_discount.Voucher, customerEmail string) *model.AppError {
	_, appErr := a.VoucherCustomerByCustomerEmailAndVoucherID(voucher.Id, customerEmail)
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return appErr
		}
	}

	return nil
}

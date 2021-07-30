package discount

import (
	"net/http"

	"github.com/shopspring/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

func (a *AppDiscount) GetVoucherDiscount(voucher *product_and_discount.Voucher, channelID string) (DiscountCalculator, *model.AppError) {
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
		if appErr.StatusCode == http.StatusInternalServerError { // must returns here since it's system error
			return appErr
		}
	}

	return nil
}

// ValidateVoucherOnlyForStaff validate if voucher is only for staff
func (a *AppDiscount) ValidateVoucherOnlyForStaff(voucher *product_and_discount.Voucher, customer *account.User) *model.AppError {
	if voucher.OnlyForStaff != nil && !*voucher.OnlyForStaff {
		return nil
	}

	var violatVoucherOnlyForStaff bool
	if customer == nil {
		violatVoucherOnlyForStaff = true
	}

	_, appErr := a.ShopApp().ShopStaffRelationByShopIDAndStaffID(voucher.ShopID, customer.Id)
	if appErr != nil {
		if appErr.StatusCode == http.StatusNotFound {
			violatVoucherOnlyForStaff = true
		} else {
			return appErr
		}
	}

	if violatVoucherOnlyForStaff {
		return model.NewAppError("ValidateVoucherOnlyForStaff", "app.shop.voucher_for_staff_only.app_error", nil, "", http.StatusNotAcceptable)
	}

	return nil
}

// VouchersByOption finds all vouchers with given option then returns them
func (a *AppDiscount) VouchersByOption(option *product_and_discount.VoucherFilterOption) ([]*product_and_discount.Voucher, *model.AppError) {
	vouchers, err := a.Srv().Store.DiscountVoucher().FilterVouchersByOption(option)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("VouchersByOption", "app.discount.vouchers_by_option_error.app_error", err)
	}

	return vouchers, nil
}

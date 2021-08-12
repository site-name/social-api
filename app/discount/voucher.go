package discount

import (
	"fmt"
	"net/http"

	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
)

// UpsertVoucher update or insert given voucher
func (a *AppDiscount) UpsertVoucher(voucher *product_and_discount.Voucher) (*product_and_discount.Voucher, *model.AppError) {
	voucher, err := a.Srv().Store.DiscountVoucher().Upsert(voucher)
	if err != nil {
		if appErr, ok := err.(*model.AppError); ok {
			return nil, appErr
		}
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}
		return nil, model.NewAppError("UpsertVoucher", "app.discount.upsert_voucher_error.app_error", nil, err.Error(), statusCode)
	}

	return voucher, nil
}

// VoucherById finds and returns a voucher with given id
func (a *AppDiscount) VoucherById(voucherID string) (*product_and_discount.Voucher, *model.AppError) {
	voucher, err := a.Srv().Store.DiscountVoucher().Get(voucherID)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("VoucherById", "app.discount.voucher_missing.app_error", err)
	}
	return voucher, nil
}

func (a *AppDiscount) GetVoucherDiscount(voucher *product_and_discount.Voucher, channelID string) (DiscountCalculator, *model.AppError) {
	voucherChannelListings, appErr := a.VoucherChannelListingsByOption(&product_and_discount.VoucherChannelListingFilterOption{
		VoucherID: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: voucher.Id,
			},
		},
		ChannelID: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: channelID,
			},
		},
	})
	if appErr != nil {
		return nil, appErr
	}

	// chose the first listing since these result is already sorted during database look up
	firstListing := voucherChannelListings[0]

	if voucher.DiscountValueType == product_and_discount.FIXED {
		return Decorator(&goprices.Money{
			Amount:   firstListing.DiscountValue,
			Currency: firstListing.Currency,
		}), nil
	}

	// otherwise DiscountValueType is 'percentage'
	return Decorator(firstListing.DiscountValue), nil
}

// GetDiscountAmountFor checks given voucher's `DiscountValueType` and returns according discount calculator function
//
//  price.(type) == *Money || *MoneyRange || *TaxedMoney || *TaxedMoneyRange
//
// NOTE: the returning interface's type should be identical to given price's type
func (a *AppDiscount) GetDiscountAmountFor(voucher *product_and_discount.Voucher, price interface{}, channelID string) (interface{}, *model.AppError) {
	// validate given price has valid type
	switch priceType := price.(type) {
	case *goprices.Money,
		*goprices.MoneyRange,
		*goprices.TaxedMoney,
		*goprices.TaxedMoneyRange:

	default:
		return nil, model.NewAppError("GetDiscountAmountFor", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "price"}, fmt.Sprintf("price's type is unexpected: %T", priceType), http.StatusBadRequest)
	}

	discountCalculator, appErr := a.GetVoucherDiscount(voucher, channelID)
	if appErr != nil {
		return nil, appErr
	}

	afterDiscount, err := discountCalculator(price) // pass in 1 argument here mean calling fixed discount calculator
	if err != nil {
		// this error maybe caused by user. But we tomporarily set status code to 500
		return nil, model.NewAppError("GetDiscountAmountFor", "app.discount.error_calculating_discount.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	switch priceType := price.(type) {
	case *goprices.Money:
		if afterDiscount.(*goprices.Money).Amount.LessThan(decimal.Zero) {
			return priceType, nil
		}
		sub, _ := priceType.Sub(afterDiscount.(*goprices.Money))
		return sub, nil

	case *goprices.MoneyRange:
		zeroMoneyRange, _ := util.ZeroMoneyRange(priceType.Currency)
		if less, err := afterDiscount.(*goprices.MoneyRange).LessThan(zeroMoneyRange); less && err == nil {
			return priceType, nil
		}

		sub, _ := priceType.Sub(afterDiscount)
		return sub, nil

	case *goprices.TaxedMoney:
		zeroTaxedMoney, _ := util.ZeroTaxedMoney(priceType.Currency)
		if less, err := afterDiscount.(*goprices.TaxedMoney).LessThan(zeroTaxedMoney); less && err == nil {
			return priceType, nil
		}

		sub, _ := priceType.Sub(afterDiscount)
		return sub, nil

	case *goprices.TaxedMoneyRange:
		zeroTaxedMoneyRange, _ := util.ZeroTaxedMoneyRange(priceType.Currency)
		if less, err := afterDiscount.(*goprices.TaxedMoneyRange).LessThan(zeroTaxedMoneyRange); less && err == nil {
			return priceType, nil
		}

		sub, _ := priceType.Sub(afterDiscount)
		return sub, nil

	default:
		return nil, nil // this code is not reached since we've already validated price's type
	}
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

	voucherChannelListings, appErr := a.VoucherChannelListingsByOption(&product_and_discount.VoucherChannelListingFilterOption{
		VoucherID: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: voucher.Id,
			},
		},
		ChannelID: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: channelID,
			},
		},
	})
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
	voucherCustomers, appErr := a.VoucherCustomerByCustomerEmailAndVoucherID(voucher.Id, customerEmail)
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError { // must returns here since it's system error
			return appErr
		}
	}
	if len(voucherCustomers) >= 1 {
		return model.NewAppError("ValidateOncePerCustomer", "app.discount.offer_only_apply_once_per_customer.app_error", nil, "", http.StatusNotAcceptable)
	}

	return nil
}

// ValidateVoucherOnlyForStaff validate if voucher is only for staff
func (a *AppDiscount) ValidateVoucherOnlyForStaff(voucher *product_and_discount.Voucher, customerID string) *model.AppError {
	if !*voucher.OnlyForStaff {
		return nil
	}

	if !model.IsValidId(customerID) {
		return model.NewAppError("ValidateVoucherOnlyForStaff", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "customerID"}, "", http.StatusBadRequest)
	}

	// try checking if there is a relationship between the shop(owner of this voucher) and the customer
	// if no reation found, it means this customer cannot have this voucher
	relation, appErr := a.ShopApp().ShopStaffRelationByShopIDAndStaffID(voucher.ShopID, customerID)
	if appErr != nil {
		if appErr.StatusCode == http.StatusNotFound || relation == nil {
			return model.NewAppError("ValidateVoucherOnlyForStaff", "app.shop.voucher_for_staff_only.app_error", nil, "", http.StatusNotAcceptable)
		}
		// error caused by server, returns immediately
		return appErr
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

// PromoCodeIsVoucher checks if given code is belong to a voucher
func (a *AppDiscount) PromoCodeIsVoucher(code string) (bool, *model.AppError) {
	vouchers, err := a.Srv().Store.DiscountVoucher().FilterVouchersByOption(&product_and_discount.VoucherFilterOption{
		Code: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: code,
			},
		},
	})
	if err != nil {
		if _, ok := err.(*store.ErrNotFound); ok {
			return false, nil
		}
		return false, store.AppErrorFromDatabaseLookupError("PromoCodeIsVoucher", "app.discount.error_finding_vouchers_with_promo_code.app_error", err)
	}

	return len(vouchers) != 0, nil
}

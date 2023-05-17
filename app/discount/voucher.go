package discount

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/discount/types"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
)

// UpsertVoucher update or insert given voucher
func (a *ServiceDiscount) UpsertVoucher(voucher *model.Voucher) (*model.Voucher, *model.AppError) {
	voucher, err := a.srv.Store.DiscountVoucher().Upsert(voucher)
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
func (a *ServiceDiscount) VoucherById(voucherID string) (*model.Voucher, *model.AppError) {
	voucher, err := a.srv.Store.DiscountVoucher().Get(voucherID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}
		return nil, model.NewAppError("VoucherById", "app.discount.voucher_missing.app_error", nil, err.Error(), statusCode)
	}
	return voucher, nil
}

// GetVoucherDiscount
func (a *ServiceDiscount) GetVoucherDiscount(voucher *model.Voucher, channelID string) (types.DiscountCalculator, *model.AppError) {
	voucherChannelListings, appErr := a.VoucherChannelListingsByOption(&model.VoucherChannelListingFilterOption{
		VoucherID: squirrel.Eq{store.VoucherChannelListingTableName + ".VoucherID": voucher.Id},
		ChannelID: squirrel.Eq{store.VoucherChannelListingTableName + ".ChannelID": channelID},
	})
	if appErr != nil {
		return nil, appErr
	}

	// chose the first listing since these result is already sorted during database look up
	firstListing := voucherChannelListings[0]

	if voucher.DiscountValueType == model.FIXED {
		return a.Decorator(&goprices.Money{
			Amount:   *firstListing.DiscountValue,
			Currency: firstListing.Currency,
		}), nil
	}

	// otherwise DiscountValueType is 'percentage'
	return a.Decorator(firstListing.DiscountValue), nil
}

// GetDiscountAmountFor checks given voucher's `DiscountValueType` and returns according discount calculator function
//
//	price.(type) == *Money || *MoneyRange || *TaxedMoney || *TaxedMoneyRange
//
// NOTE: the returning interface's type should be identical to given price's type
func (a *ServiceDiscount) GetDiscountAmountFor(voucher *model.Voucher, price interface{}, channelID string) (interface{}, *model.AppError) {
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

	afterDiscount, err := discountCalculator(price, nil) // pass in 1 argument here mean calling fixed discount calculator
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
		if afterDiscount.(*goprices.MoneyRange).LessThan(zeroMoneyRange) {
			return priceType, nil
		}

		sub, _ := priceType.Sub(afterDiscount)
		return sub, nil

	case *goprices.TaxedMoney:
		zeroTaxedMoney, _ := util.ZeroTaxedMoney(priceType.Currency)
		if afterDiscount.(*goprices.TaxedMoney).LessThan(zeroTaxedMoney) {
			return priceType, nil
		}

		sub, _ := priceType.Sub(afterDiscount)
		return sub, nil

	case *goprices.TaxedMoneyRange:
		zeroTaxedMoneyRange, _ := util.ZeroTaxedMoneyRange(priceType.Currency)
		if afterDiscount.(*goprices.TaxedMoneyRange).LessThan(zeroTaxedMoneyRange) {
			return priceType, nil
		}

		sub, _ := priceType.Sub(afterDiscount)
		return sub, nil

	default:
		return nil, nil // this code is not reached since we've already validated price's type
	}
}

// ValidateMinSpent validates if the order cost at least a specific amount of money
func (a *ServiceDiscount) ValidateMinSpent(voucher *model.Voucher, value *goprices.TaxedMoney, channelID string) (notApplicableErr *model.NotApplicable, appErr *model.AppError) {
	ownerShopOfVoucher, appErr := a.srv.ShopService().ShopById(voucher.ShopID)
	if appErr != nil {
		return
	}

	money := value.Net
	if *ownerShopOfVoucher.DisplayGrossPrices {
		money = value.Gross
	}

	voucherChannelListings, appErr := a.VoucherChannelListingsByOption(&model.VoucherChannelListingFilterOption{
		VoucherID: squirrel.Eq{store.VoucherChannelListingTableName + ".VoucherID": voucher.Id},
		ChannelID: squirrel.Eq{store.VoucherChannelListingTableName + ".ChannelID": channelID},
	})
	if appErr != nil {
		return
	}

	if len(voucherChannelListings) == 0 {
		notApplicableErr = &model.NotApplicable{
			Message: "This voucher is not assigned to this channel",
		}
		return
	}

	voucherChannelListing := voucherChannelListings[0]
	voucherChannelListing.PopulateNonDbFields()

	if voucherChannelListing.MinSpent != nil {
		if money.LessThan(voucherChannelListing.MinSpent) {
			notApplicableErr = &model.NotApplicable{
				Message: "This offer is only valid for orders over " + voucherChannelListing.MinSpent.Amount.String(),
			}
			return
		}
	}

	return
}

// ValidateOncePerCustomer checks to make sure each customer has ONLY 1 time usage with 1 voucher
func (a *ServiceDiscount) ValidateOncePerCustomer(voucher *model.Voucher, customerEmail string) (notApplicableErr *model.NotApplicable, appErr *model.AppError) {
	voucherCustomers, appErr := a.VoucherCustomersByOption(&model.VoucherCustomerFilterOption{
		VoucherID:     squirrel.Eq{store.VoucherCustomerTableName + ".VoucherID": voucher.Id},
		CustomerEmail: squirrel.Eq{store.VoucherCustomerTableName + ".CustomerEmail": customerEmail},
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError { // must returns here since it's system error
			return
		}
	}
	if len(voucherCustomers) >= 1 {
		return &model.NotApplicable{
			Message: "This offer is valid only once per customer.",
		}, nil
	}

	return
}

// ValidateVoucherOnlyForStaff validate if voucher is only for staff
func (a *ServiceDiscount) ValidateVoucherOnlyForStaff(voucher *model.Voucher, customerID string) (*model.NotApplicable, *model.AppError) {
	if !*voucher.OnlyForStaff {
		return nil, nil
	}

	// try checking if there is a relationship between the shop(owner of this voucher) and the customer
	// if no reation found, it means this customer cannot have this voucher
	relation, appErr := a.srv.ShopService().ShopStaffByOptions(&model.ShopStaffFilterOptions{
		StaffID: squirrel.Eq{store.ShopStaffTableName + ".StaffID": customerID},
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusNotFound || relation == nil {
			return &model.NotApplicable{
				Message: "This offer is valid only for staff customers",
			}, nil
		}
		return nil, appErr
	}

	return nil, nil
}

// VouchersByOption finds all vouchers with given option then returns them
func (a *ServiceDiscount) VouchersByOption(option *model.VoucherFilterOption) ([]*model.Voucher, *model.AppError) {
	vouchers, err := a.srv.Store.DiscountVoucher().FilterVouchersByOption(option)
	if err != nil {
		return nil, model.NewAppError("VouchersByOption", "app.discount.error_finding_vouchers_by_option_error.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return vouchers, nil
}

// VoucherByOption returns 1 voucher filtered using given options
func (s *ServiceDiscount) VoucherByOption(options *model.VoucherFilterOption) (*model.Voucher, *model.AppError) {
	voucher, err := s.srv.Store.DiscountVoucher().GetByOptions(options)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}
		return nil, model.NewAppError("VoucherByOption", "app.discount.error_finding_voucher_by_option.app_error", nil, err.Error(), statusCode)
	}
	return voucher, nil
}

// PromoCodeIsVoucher checks if given code is belong to a voucher
func (a *ServiceDiscount) PromoCodeIsVoucher(code string) (bool, *model.AppError) {
	vouchers, appErr := a.VouchersByOption(&model.VoucherFilterOption{
		Code: squirrel.Eq{store.VoucherTableName + ".Code": code},
	})
	if appErr != nil {
		return false, appErr
	}

	return len(vouchers) != 0, nil
}

// FilterActiveVouchers returns a list of vouchers that are active.
//
// `channelSlug` is optional (can be empty). pass this argument if you want to find active vouchers in specific channel
func (s *ServiceDiscount) FilterActiveVouchers(date time.Time, channelSlug string) ([]*model.Voucher, *model.AppError) {
	filterOptions := &model.VoucherFilterOption{
		UsageLimit: squirrel.Or{
			squirrel.Eq{store.VoucherTableName + ".UsageLimit": nil},
			squirrel.Gt{store.VoucherTableName + ".UsageLimit": store.VoucherTableName + ".Used"},
		},
		EndDate: squirrel.Or{
			squirrel.Eq{store.VoucherTableName + ".EndDate": nil},
			squirrel.GtOrEq{store.VoucherTableName + ".EndDate": util.StartOfDay(date)},
		},
		StartDate: squirrel.LtOrEq{store.VoucherTableName + ".StartDate": util.StartOfDay(date)},
	}

	if channelSlug != "" {
		filterOptions.ChannelListingSlug = squirrel.Eq{store.ChannelTableName + ".Slug": channelSlug}
		filterOptions.ChannelListingActive = model.NewPrimitive(true)
	}

	return s.VouchersByOption(filterOptions)
}

// ExpiredVouchers returns vouchers that are expired before given date (beginning of the day). If date is nil, use today instead
func (s *ServiceDiscount) ExpiredVouchers(date *time.Time) ([]*model.Voucher, *model.AppError) {
	expiredVouchers, err := s.srv.Store.DiscountVoucher().ExpiredVouchers(date)
	if err != nil {
		return nil, model.NewAppError("ExpiredVouchers", "app.discount.error_finding_expired_vouchers.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return expiredVouchers, nil
}

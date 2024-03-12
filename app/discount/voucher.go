package discount

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/mattermost/squirrel"
	"github.com/samber/lo"
	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app/discount/types"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func (a *ServiceDiscount) UpsertVoucher(voucher model.Voucher) (*model.Voucher, *model_helper.AppError) {
	upserdVoucher, err := a.srv.Store.DiscountVoucher().Upsert(voucher)
	if err != nil {
		if appErr, ok := err.(*model_helper.AppError); ok {
			return nil, appErr
		}
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrInvalidInput); ok {
			statusCode = http.StatusBadRequest
		}
		return nil, model_helper.NewAppError("UpsertVoucher", "app.discount.upsert_voucher_error.app_error", nil, err.Error(), statusCode)
	}

	return upserdVoucher, nil
}

func (a *ServiceDiscount) VoucherById(voucherID string) (*model.Voucher, *model_helper.AppError) {
	voucher, err := a.srv.Store.DiscountVoucher().Get(voucherID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}
		return nil, model_helper.NewAppError("VoucherById", "app.discount.voucher_missing.app_error", nil, err.Error(), statusCode)
	}
	return voucher, nil
}

func (a *ServiceDiscount) GetVoucherDiscount(voucher model.Voucher, channelID string) (types.DiscountCalculator, *model_helper.AppError) {
	voucherChannelListings, appErr := a.VoucherChannelListingsByOption(model_helper.VoucherChannelListingFilterOption{
		CommonQueryOptions: model_helper.NewCommonQueryOptions(
			model.VoucherChannelListingWhere.VoucherID.EQ(voucher.ID),
			model.VoucherChannelListingWhere.ChannelID.EQ(channelID),
		),
	})
	if appErr != nil {
		return nil, appErr
	}

	// chose the first listing since these result is already sorted during database look up
	firstListing := voucherChannelListings[0]

	if voucher.DiscountValueType == model.DiscountValueTypeFixed {
		return a.Decorator(goprices.Money{
			Amount:   firstListing.DiscountValue,
			Currency: firstListing.Currency.String(),
		}), nil
	}

	// otherwise DiscountValueType is 'percentage'
	return a.Decorator(firstListing.DiscountValue), nil
}

// GetDiscountAmountFor checks given voucher's `DiscountValueType` and returns according discount calculator function
//
//	price.(type) == Money || MoneyRange || TaxedMoney || TaxedMoneyRange
//
// NOTE: the returning interface's type should be identical to given price's type
func (a *ServiceDiscount) GetDiscountAmountFor(voucher model.Voucher, price any, channelID string) (any, *model_helper.AppError) {
	// validate given price has valid type
	switch priceType := price.(type) {
	case goprices.Money,
		goprices.MoneyRange,
		goprices.TaxedMoney,
		goprices.TaxedMoneyRange:

	default:
		return nil, model_helper.NewAppError("GetDiscountAmountFor", model_helper.InvalidArgumentAppErrorID, map[string]any{"Fields": "price"}, fmt.Sprintf("price's type is unexpected: %T", priceType), http.StatusBadRequest)
	}

	discountCalculator, appErr := a.GetVoucherDiscount(voucher, channelID)
	if appErr != nil {
		return nil, appErr
	}

	afterDiscount, err := discountCalculator(price, nil) // pass in 1 argument here mean calling fixed discount calculator
	if err != nil {
		// this error maybe caused by user. But we tomporarily set status code to 500
		return nil, model_helper.NewAppError("GetDiscountAmountFor", "app.discount.error_calculating_discount.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	switch priceType := price.(type) {
	case goprices.Money:
		if afterDiscount.(goprices.Money).Amount.LessThan(decimal.Zero) {
			return priceType, nil
		}
		sub, _ := priceType.Sub(afterDiscount.(goprices.Money))
		return sub, nil

	case goprices.MoneyRange:
		zeroMoneyRange, _ := util.ZeroMoneyRange(priceType.GetCurrency())
		afterDiscountMoneyRange := afterDiscount.(goprices.MoneyRange)
		if afterDiscountMoneyRange.LessThan(*zeroMoneyRange) {
			return priceType, nil
		}

		sub, _ := priceType.Sub(afterDiscount)
		return sub, nil

	case goprices.TaxedMoney:
		zeroTaxedMoney, _ := util.ZeroTaxedMoney(priceType.GetCurrency())
		afterDiscountTaxedMoney := afterDiscount.(goprices.TaxedMoney)
		if afterDiscountTaxedMoney.LessThan(*zeroTaxedMoney) {
			return priceType, nil
		}

		sub, _ := priceType.Sub(afterDiscount)
		return sub, nil

	case goprices.TaxedMoneyRange:
		zeroTaxedMoneyRange, _ := util.ZeroTaxedMoneyRange(priceType.GetCurrency())
		afterDiscountTaxedMoneyRange := afterDiscount.(goprices.TaxedMoneyRange)
		if afterDiscountTaxedMoneyRange.LessThan(*zeroTaxedMoneyRange) {
			return priceType, nil
		}

		sub, _ := priceType.Sub(afterDiscount)
		return sub, nil

	default:
		return nil, nil // this code is not reached since we've already validated price's type
	}
}

// ValidateMinSpent validates if the order cost at least a specific amount of money
func (a *ServiceDiscount) ValidateMinSpent(voucher model.Voucher, value goprices.TaxedMoney, channelID string) (notApplicableErr *model_helper.NotApplicable, appErr *model_helper.AppError) {
	money := value.Net
	if *a.srv.Config().ShopSettings.DisplayGrossPrices {
		money = value.Gross
	}

	voucherChannelListings, appErr := a.VoucherChannelListingsByOption(model_helper.VoucherChannelListingFilterOption{
		CommonQueryOptions: model_helper.NewCommonQueryOptions(
			model.VoucherChannelListingWhere.VoucherID.EQ(voucher.ID),
			model.VoucherChannelListingWhere.ChannelID.EQ(channelID),
		),
	})
	if appErr != nil {
		return
	}

	if len(voucherChannelListings) == 0 {
		notApplicableErr = &model_helper.NotApplicable{
			Message: "This voucher is not assigned to this channel",
		}
		return
	}

	minSpent := model_helper.VoucherChannelListingGetMinSpent(*voucherChannelListings[0])
	if money.LessThan(minSpent) {
		notApplicableErr = &model_helper.NotApplicable{
			Message: "This offer is only valid for orders over " + minSpent.Amount.String(),
		}
		return
	}

	return
}

// ValidateOncePerCustomer checks to make sure each customer has ONLY 1 time usage with 1 voucher
func (a *ServiceDiscount) ValidateOncePerCustomer(voucher model.Voucher, customerEmail string) (notApplicableErr *model_helper.NotApplicable, appErr *model_helper.AppError) {
	voucherCustomers, appErr := a.VoucherCustomersByOption(model_helper.VoucherCustomerFilterOption{
		CommonQueryOptions: model_helper.NewCommonQueryOptions(
			model.VoucherCustomerWhere.VoucherID.EQ(voucher.ID),
			model.VoucherCustomerWhere.CustomerEmail.EQ(customerEmail),
		),
	})
	if appErr != nil {
		return nil, appErr
	}
	if len(voucherCustomers) >= 1 {
		return &model_helper.NotApplicable{
			Message: "This offer is valid only once per customer.",
		}, nil
	}

	return
}

func (a *ServiceDiscount) ValidateOnlyForStaff(voucher model.Voucher, customerID string) (*model_helper.NotApplicable, *model_helper.AppError) {
	if !voucher.OnlyForStaff.IsNil() && !*voucher.OnlyForStaff.Bool {
		return nil, nil
	}

	customer, appErr := a.srv.Account.UserById(context.Background(), customerID)
	if appErr != nil {
		return nil, appErr
	}

	// NOTE: shop admin also has staff role
	if !lo.Contains(strings.Fields(customer.Roles), model_helper.ShopStaffRoleId) {
		return model_helper.NewNotApplicable("ValidateOnlyForStaff", "this offer is for shop staffs only", nil, 0), nil
	}

	return nil, nil
}

func (a *ServiceDiscount) VouchersByOption(option model_helper.VoucherFilterOption) (model.VoucherSlice, *model_helper.AppError) {
	vouchers, err := a.srv.Store.DiscountVoucher().FilterVouchersByOption(option)
	if err != nil {
		var statusCode = http.StatusInternalServerError
		if _, ok := err.(*store.ErrInvalidInput); ok {
			statusCode = http.StatusBadRequest
		}
		return nil, model_helper.NewAppError("VouchersByOption", "app.discount.error_finding_vouchers_by_option_error.app_error", nil, err.Error(), statusCode)
	}

	return vouchers, nil
}

func (s *ServiceDiscount) VoucherByOption(options model_helper.VoucherFilterOption) (*model.Voucher, *model_helper.AppError) {
	vouchers, appErr := s.VouchersByOption(options)
	if appErr != nil {
		return nil, appErr
	}
	return vouchers[0], nil
}

func (a *ServiceDiscount) PromoCodeIsVoucher(code string) (bool, *model_helper.AppError) {
	vouchers, appErr := a.VouchersByOption(model_helper.VoucherFilterOption{
		CommonQueryOptions: model_helper.NewCommonQueryOptions(
			model.VoucherWhere.Code.EQ(code),
		),
	})
	if appErr != nil {
		return false, appErr
	}

	return len(vouchers) != 0, nil
}

// FilterActiveVouchers returns a list of vouchers that are active.
//
// `channelSlug` is optional (can be empty). pass this argument if you want to find active vouchers in specific channel
func (s *ServiceDiscount) FilterActiveVouchers(date time.Time, channelSlug string) (model.VoucherSlice, *model_helper.AppError) {
	startOfDay := util.StartOfDay(date)
	filterOptions := &model.VoucherFilterOption{
		Conditions: squirrel.And{
			squirrel.Expr(model.VoucherTableName + ".UsageLimit IS NULL OR Vouchers.UsageLimit > Vouchers.Used"),
			squirrel.Expr(model.VoucherTableName+".EndDate IS NULL OR Vouchers.EndDate >= ?", startOfDay),
			squirrel.Expr(model.VoucherTableName+".StartDate <= ?", startOfDay),
		},
	}

	if channelSlug != "" {
		filterOptions.VoucherChannelListing_ChannelSlug = squirrel.Expr(model.ChannelTableName+".Slug = ?", channelSlug)
		filterOptions.VoucherChannelListing_ChannelIsActive = squirrel.Expr(model.ChannelTableName + ".IsActive")
	}

	_, vouchers, appErr := s.VouchersByOption(filterOptions)
	return vouchers, appErr
}

// ExpiredVouchers returns vouchers that are expired before given date (beginning of the day). If date is nil, use today instead
func (s *ServiceDiscount) ExpiredVouchers(date *time.Time) (model.VoucherSlice, *model_helper.AppError) {
	expiredVouchers, err := s.srv.Store.DiscountVoucher().ExpiredVouchers(date)
	if err != nil {
		return nil, model_helper.NewAppError("ExpiredVouchers", "app.discount.error_finding_expired_vouchers.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return expiredVouchers, nil
}

func (s *ServiceDiscount) ToggleVoucherRelations(transaction boil.ContextTransactor, vouchers model.Vouchers, productIDs, variantIDs, categoryIDs, collectionIDs []string, isDelete bool) *model_helper.AppError {
	err := s.srv.Store.DiscountVoucher().ToggleVoucherRelations(transaction, vouchers, collectionIDs, productIDs, variantIDs, categoryIDs, isDelete)
	if err != nil {
		return model_helper.NewAppError("ToggleVoucherRelations", "app.discount.insert_voucher_relations.app_error", nil, "failed to insert voucher relations", http.StatusInternalServerError)
	}

	return nil
}

// VoucherChannelListingsByOption finds voucher channel listings based on given options
func (a *ServiceDiscount) VoucherChannelListingsByOption(option model_helper.VoucherChannelListingFilterOption) (model.VoucherChannelListingSlice, *model_helper.AppError) {
	listings, err := a.srv.Store.VoucherChannelListing().FilterbyOption(option)
	if err != nil {
		return nil, model_helper.NewAppError("VoucherChannelListingsByOption", "app.discount.error_finding_voucher_channel_listings_by_option.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return listings, nil
}

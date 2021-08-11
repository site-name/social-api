package discount

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

// VoucherChannelListingsByOption finds voucher channel listings based on given options
func (a *AppDiscount) VoucherChannelListingsByOption(option *product_and_discount.VoucherChannelListingFilterOption) ([]*product_and_discount.VoucherChannelListing, *model.AppError) {
	listings, err := a.Srv().Store.VoucherChannelListing().FilterbyOption(option)

	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("VoucherChannelListingsByOption", "app.discount.error_finding_voucher_channel_listings_by_option.app_error", err)
	}

	return listings, nil
}

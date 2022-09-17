package discount

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

// VoucherChannelListingsByOption finds voucher channel listings based on given options
func (a *ServiceDiscount) VoucherChannelListingsByOption(option *model.VoucherChannelListingFilterOption) ([]*model.VoucherChannelListing, *model.AppError) {
	listings, err := a.srv.Store.VoucherChannelListing().FilterbyOption(option)

	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("VoucherChannelListingsByOption", "app.discount.error_finding_voucher_channel_listings_by_option.app_error", err)
	}

	return listings, nil
}

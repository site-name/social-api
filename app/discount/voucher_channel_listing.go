package discount

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

func (a *AppDiscount) VoucherChannelListingsByVoucherAndChannel(voucherID string, channelID string) ([]*product_and_discount.VoucherChannelListing, *model.AppError) {
	listings, err := a.Srv().Store.VoucherChannelListing().FilterByVoucherAndChannel(voucherID, channelID)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("VoucherChannelListingsByVoucherAndChannel", "app.discount.error_finding_voucher_channel_listings_by_channel_and_voucher.app_error", err)
	}

	return listings, nil
}

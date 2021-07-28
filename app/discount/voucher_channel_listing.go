package discount

import (
	"net/http"

	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

func (a *AppDiscount) VoucherChannelListingsByVoucherAndChannel(voucherID string, channelID string) ([]*product_and_discount.VoucherChannelListing, *model.AppError) {
	var invalidArg bool
	if !model.IsValidId(voucherID) {
		invalidArg = true
	}
	if !model.IsValidId(channelID) {
		invalidArg = true
	}
	if invalidArg {
		return nil, model.NewAppError("VoucherChannelListingsByVoucherAndChannel", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "voucherID, channelID"}, "", http.StatusBadRequest)
	}

	listings, err := a.Srv().Store.VoucherChannelListing().FilterByVoucherAndChannel(voucherID, channelID)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("VoucherChannelListingsByVoucherAndChannel", "app.discount.error_finding_voucher_channel_listings_by_channel_and_voucher.app_error", err)
	}

	return listings, nil
}

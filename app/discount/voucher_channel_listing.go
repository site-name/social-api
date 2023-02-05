package discount

import (
	"net/http"

	"github.com/sitename/sitename/model"
)

// VoucherChannelListingsByOption finds voucher channel listings based on given options
func (a *ServiceDiscount) VoucherChannelListingsByOption(option *model.VoucherChannelListingFilterOption) ([]*model.VoucherChannelListing, *model.AppError) {
	listings, err := a.srv.Store.VoucherChannelListing().FilterbyOption(option)
	if err != nil {
		return nil, model.NewAppError("VoucherChannelListingsByOption", "app.discount.error_finding_voucher_channel_listings_by_option.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return listings, nil
}

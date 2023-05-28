package shipping

import (
	"net/http"

	"github.com/sitename/sitename/model"
)

// ShippingMethodChannelListingsByOption returns a list of shipping method channel listings by given option
func (a *ServiceShipping) ShippingMethodChannelListingsByOption(option *model.ShippingMethodChannelListingFilterOption) (model.ShippingMethodChannelListings, *model.AppError) {
	listings, err := a.srv.Store.ShippingMethodChannelListing().FilterByOption(option)
	if err != nil {
		return nil, model.NewAppError("ShippingMethodChannelListingsByOption", "app.shipping.error_finding_shipping_method_channel_listings_by_option.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return listings, nil
}

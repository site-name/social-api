package shipping

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/shipping"
	"github.com/sitename/sitename/store"
)

// ShippingMethodChannelListingsByOption returns a list of shipping method channel listings by given option
func (a *ServiceShipping) ShippingMethodChannelListingsByOption(option *shipping.ShippingMethodChannelListingFilterOption) ([]*shipping.ShippingMethodChannelListing, *model.AppError) {
	listings, err := a.srv.Store.ShippingMethodChannelListing().FilterByOption(option)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("ShippingMethodChannelListingsByOption", "app.shipping.shipping_method_channel_listings_by_option.app_error", err)
	}

	return listings, nil
}

package shipping

import (
	"net/http"

	"github.com/sitename/sitename/model"
)

// ShippingMethodChannelListingsByOption returns a list of shipping method channel listings by given option
func (a *ServiceShipping) ShippingMethodChannelListingsByOption(option *model.ShippingMethodChannelListingFilterOption) ([]*model.ShippingMethodChannelListing, *model.AppError) {
	listings, err := a.srv.Store.ShippingMethodChannelListing().FilterByOption(option)
	var (
		statusCode int
		errMessage string
	)
	if err != nil {
		statusCode = http.StatusInternalServerError
		errMessage = err.Error()
	} else if len(listings) == 0 {
		statusCode = http.StatusNotFound
	}

	if statusCode != 0 {
		return nil, model.NewAppError("ShippingMethodChannelListingsByOption", "app.shipping.error_finding_shipping_method_channel_listings_by_option.app_error", nil, errMessage, statusCode)
	}

	return listings, nil
}

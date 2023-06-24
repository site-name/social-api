package shipping

import (
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/samber/lo"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

// ShippingMethodChannelListingsByOption returns a list of shipping method channel listings by given option
func (a *ServiceShipping) ShippingMethodChannelListingsByOption(option *model.ShippingMethodChannelListingFilterOption) (model.ShippingMethodChannelListings, *model.AppError) {
	listings, err := a.srv.Store.ShippingMethodChannelListing().FilterByOption(option)
	if err != nil {
		return nil, model.NewAppError("ShippingMethodChannelListingsByOption", "app.shipping.error_finding_shipping_method_channel_listings_by_option.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return listings, nil
}

// Prepare mapping shipping method to price from channel listings
func (a *ServiceShipping) GetShippingMethodToShippingPriceMapping(shippingMethods model.ShippingMethods, channelSlug string) (map[string]*goprices.Money, *model.AppError) {
	listings, appErr := a.ShippingMethodChannelListingsByOption(&model.ShippingMethodChannelListingFilterOption{
		ShippingMethodID: squirrel.Eq{store.ShippingMethodChannelListingTableName + ".ShippingMethodID": shippingMethods.IDs()},
		ChannelSlug:      squirrel.Eq{store.ChannelTableName + ".Slug": channelSlug},
	})
	if appErr != nil {
		return nil, appErr
	}

	return lo.SliceToMap(listings, func(lst *model.ShippingMethodChannelListing) (string, *goprices.Money) {
		lst.PopulateNonDbFields() // this call is required
		return lst.ShippingMethodID, lst.Price
	}), nil
}

package shipping

import (
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/samber/lo"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/store_iface"
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

func (s *ServiceShipping) UpsertShippingMethodChannelListings(transaction store_iface.SqlxTxExecutor, listings model.ShippingMethodChannelListings) (model.ShippingMethodChannelListings, *model.AppError) {
	listings, err := s.srv.Store.ShippingMethodChannelListing().Upsert(transaction, listings)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrInvalidInput); ok {
			statusCode = http.StatusBadRequest
		}

		return nil, model.NewAppError("UpsertShippingMethodChannelListings", "app.shipping.upsert_shipping_method_channel_listings.app_error", nil, err.Error(), statusCode)
	}

	return listings, nil
}

func (s *ServiceShipping) DeleteShippingMethodChannelListings(transaction store_iface.SqlxTxExecutor, options *model.ShippingMethodChannelListingFilterOption) *model.AppError {
	err := s.srv.Store.ShippingMethodChannelListing().BulkDelete(transaction, options)
	if err != nil {
		model.NewAppError("ShippingMethodChannelListingUpdate", "app.shipping.delete_shipping_method_channel_listings.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

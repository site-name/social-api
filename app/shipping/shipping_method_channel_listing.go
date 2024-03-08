package shipping

import (
	"net/http"

	"github.com/mattermost/squirrel"
	"github.com/samber/lo"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

func (a *ServiceShipping) ShippingMethodChannelListingsByOption(option model_helper.ShippingMethodChannelListingFilterOption) (model.ShippingMethodChannelListingSlice, *model_helper.AppError) {
	listings, err := a.srv.Store.ShippingMethodChannelListing().FilterByOption(option)
	if err != nil {
		return nil, model_helper.NewAppError("ShippingMethodChannelListingsByOption", "app.shipping.error_finding_shipping_method_channel_listings_by_option.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return listings, nil
}

// Prepare mapping shipping method to price from channel listings
func (a *ServiceShipping) GetShippingMethodToShippingPriceMapping(shippingMethods model.ShippingMethods, channelSlug string) (map[string]*goprices.Money, *model_helper.AppError) {
	listings, appErr := a.ShippingMethodChannelListingsByOption(&model.ShippingMethodChannelListingFilterOption{
		Conditions:  squirrel.Eq{model.ShippingMethodChannelListingTableName + ".ShippingMethodID": shippingMethods.IDs()},
		ChannelSlug: squirrel.Eq{model.ChannelTableName + ".Slug": channelSlug},
	})
	if appErr != nil {
		return nil, appErr
	}

	return lo.SliceToMap(listings, func(lst *model.ShippingMethodChannelListing) (string, *goprices.Money) {
		lst.PopulateNonDbFields() // this call is required
		return lst.ShippingMethodID, lst.Price
	}), nil
}

func (s *ServiceShipping) UpsertShippingMethodChannelListings(transaction *gorm.DB, listings model.ShippingMethodChannelListings) (model.ShippingMethodChannelListings, *model_helper.AppError) {
	listings, err := s.srv.Store.ShippingMethodChannelListing().Upsert(transaction, listings)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrInvalidInput); ok {
			statusCode = http.StatusBadRequest
		}

		return nil, model_helper.NewAppError("UpsertShippingMethodChannelListings", "app.shipping.upsert_shipping_method_channel_listings.app_error", nil, err.Error(), statusCode)
	}

	return listings, nil
}

func (s *ServiceShipping) DeleteShippingMethodChannelListings(transaction *gorm.DB, options *model.ShippingMethodChannelListingFilterOption) *model_helper.AppError {
	err := s.srv.Store.ShippingMethodChannelListing().BulkDelete(transaction, options)
	if err != nil {
		model_helper.NewAppError("ShippingMethodChannelListingUpdate", "app.shipping.delete_shipping_method_channel_listings.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

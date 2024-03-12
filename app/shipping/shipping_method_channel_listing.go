package shipping

import (
	"net/http"

	"github.com/samber/lo"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func (a *ServiceShipping) ShippingMethodChannelListingsByOption(option model_helper.ShippingMethodChannelListingFilterOption) (model.ShippingMethodChannelListingSlice, *model_helper.AppError) {
	listings, err := a.srv.Store.ShippingMethodChannelListing().FilterByOption(option)
	if err != nil {
		return nil, model_helper.NewAppError("ShippingMethodChannelListingsByOption", "app.shipping.error_finding_shipping_method_channel_listings_by_option.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return listings, nil
}

func (a *ServiceShipping) GetShippingMethodToShippingPriceMapping(shippingMethods model.ShippingMethodSlice, channelSlug string) (map[string]*goprices.Money, *model_helper.AppError) {
	shippinMethodIDs := lo.Map(shippingMethods, func(item *model.ShippingMethod, _ int) string { return item.ID })

	listings, appErr := a.ShippingMethodChannelListingsByOption(model_helper.ShippingMethodChannelListingFilterOption{
		CommonQueryOptions: model_helper.NewCommonQueryOptions(
			model.ShippingMethodChannelListingWhere.ShippingMethodID.IN(shippinMethodIDs),
		),
		ChannelSlug: model.ChannelWhere.Slug.EQ(channelSlug),
	})
	if appErr != nil {
		return nil, appErr
	}

	return lo.SliceToMap(listings, func(lst *model.ShippingMethodChannelListing) (string, *goprices.Money) {
		lst.PopulateNonDbFields() // this call is required
		return lst.ShippingMethodID, lst.Price
	}), nil
}

func (s *ServiceShipping) UpsertShippingMethodChannelListings(transaction boil.ContextTransactor, listings model.ShippingMethodChannelListingSlice) (model.ShippingMethodChannelListingSlice, *model_helper.AppError) {
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

func (s *ServiceShipping) DeleteShippingMethodChannelListings(transaction boil.ContextTransactor, ids []string) *model_helper.AppError {
	err := s.srv.Store.ShippingMethodChannelListing().Delete(transaction, ids)
	if err != nil {
		model_helper.NewAppError("ShippingMethodChannelListingUpdate", "app.shipping.delete_shipping_method_channel_listings.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

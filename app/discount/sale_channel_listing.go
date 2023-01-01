package discount

import (
	"net/http"

	"github.com/sitename/sitename/model"
)

func (s *ServiceDiscount) SaleChannelListingsByOptions(options *model.SaleChannelListingFilterOption) ([]*model.SaleChannelListing, *model.AppError) {
	listings, err := s.srv.Store.DiscountSaleChannelListing().SaleChannelListingsWithOption(options)
	if err != nil {
		return nil, model.NewAppError("SaleChannelListingsByOptions", "app.discount.error_finding_sale_channel_listings_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return listings, nil
}

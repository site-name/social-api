package product

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

func (a *AppProduct) ProductChannelListingsByOption(option *product_and_discount.ProductChannelListingFilterOption) ([]*product_and_discount.ProductChannelListing, *model.AppError) {
	listings, err := a.Srv().Store.ProductChannelListing().FilterByOption(option)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("ProductChannelListingsByOption", "app.product.product_channel_listings_by_option_missing.app_error", err)
	}

	return listings, nil
}

package product

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

// ProductChannelListingsByOption returns a list of product channel listings filtered using given option
func (a *ServiceProduct) ProductChannelListingsByOption(option *model.ProductChannelListingFilterOption) ([]*model.ProductChannelListing, *model.AppError) {
	listings, err := a.srv.Store.ProductChannelListing().FilterByOption(option)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("ProductChannelListingsByOption", "app.product.product_channel_listings_by_option_missing.app_error", err)
	}

	return listings, nil
}

// BulkUpsertProductChannelListings bulk update/inserts given product channel listings and returns them
func (a *ServiceProduct) BulkUpsertProductChannelListings(listings []*model.ProductChannelListing) ([]*model.ProductChannelListing, *model.AppError) {
	listings, err := a.srv.Store.ProductChannelListing().BulkUpsert(listings)
	if err != nil {
		if appErr, ok := err.(*model.AppError); ok {
			return nil, appErr
		}
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}
		if _, ok := err.(*store.ErrInvalidInput); ok {
			statusCode = http.StatusBadRequest
		}

		return nil, model.NewAppError("BulkUpsertProductChannelListings", "app.product.error_bulk_upserting_product_channel_listings.app_error", nil, err.Error(), statusCode)
	}

	return listings, nil
}

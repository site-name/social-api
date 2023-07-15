package product

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/store_iface"
)

// ProductChannelListingsByOption returns a list of product channel listings filtered using given option
func (a *ServiceProduct) ProductChannelListingsByOption(option *model.ProductChannelListingFilterOption) ([]*model.ProductChannelListing, *model.AppError) {
	listings, err := a.srv.Store.ProductChannelListing().FilterByOption(option)
	if err != nil {
		return nil, model.NewAppError("ProductChannelListingsByOption", "app.product.product_channel_listings_by_option_missing.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return listings, nil
}

// BulkUpsertProductChannelListings bulk update/inserts given product channel listings and returns them
func (a *ServiceProduct) BulkUpsertProductChannelListings(transaction store_iface.SqlxExecutor, listings []*model.ProductChannelListing) ([]*model.ProductChannelListing, *model.AppError) {
	listings, err := a.srv.Store.ProductChannelListing().BulkUpsert(transaction, listings)
	if err != nil {
		if appErr, ok := err.(*model.AppError); ok {
			return nil, appErr
		}
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		} else if _, ok := err.(*store.ErrInvalidInput); ok {
			statusCode = http.StatusBadRequest
		}

		return nil, model.NewAppError("BulkUpsertProductChannelListings", "app.product.error_bulk_upserting_product_channel_listings.app_error", nil, err.Error(), statusCode)
	}

	return listings, nil
}

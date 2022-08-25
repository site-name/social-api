package product

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/store_iface"
)

// ProductVariantChannelListingsByOption returns a slice of product variant channel listings by given option
func (a *ServiceProduct) ProductVariantChannelListingsByOption(transaction store_iface.SqlxTxExecutor, option *product_and_discount.ProductVariantChannelListingFilterOption) (product_and_discount.ProductVariantChannelListings, *model.AppError) {
	listings, err := a.srv.Store.ProductVariantChannelListing().FilterbyOption(transaction, option)
	var (
		statusCode   int
		errorMessage string
	)
	if err != nil {
		statusCode = http.StatusInternalServerError
		errorMessage = err.Error()
	} else if len(listings) == 0 {
		statusCode = http.StatusNotFound
	}

	if statusCode != 0 {
		return nil, model.NewAppError("ProductVariantChannelListingsByOption", "app.product_error_finding_product_variant_channel_listings_by_option.app_error", nil, errorMessage, statusCode)
	}

	return listings, nil
}

// BulkUpsertProductVariantChannelListings tells store to bulk upserts given product variant channel listings
func (s *ServiceProduct) BulkUpsertProductVariantChannelListings(transaction store_iface.SqlxTxExecutor, listings []*product_and_discount.ProductVariantChannelListing) ([]*product_and_discount.ProductVariantChannelListing, *model.AppError) {
	variantChannelListings, err := s.srv.Store.ProductVariantChannelListing().BulkUpsert(transaction, listings)
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

		return nil, model.NewAppError("BulkUpsertProductVariantChannelListings", "app.product.error_bulk_upserting_product_variant_channel_listings.app_error", nil, err.Error(), statusCode)
	}

	return variantChannelListings, nil
}

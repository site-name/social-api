package product

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

// ProductVariantChannelListingsByOption returns a slice of product variant channel listings by given option
func (a *ServiceProduct) ProductVariantChannelListingsByOption(options *model.ProductVariantChannelListingFilterOption) (model.ProductVariantChannelListings, *model.AppError) {
	listings, err := a.srv.Store.ProductVariantChannelListing().FilterbyOption(options)
	if err != nil {
		return nil, model.NewAppError("ProductVariantChannelListingsByOption", "app.product_error_finding_product_variant_channel_listings_by_option.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return listings, nil
}

// BulkUpsertProductVariantChannelListings tells store to bulk upserts given product variant channel listings
func (s *ServiceProduct) BulkUpsertProductVariantChannelListings(transaction *gorm.DB, listings []*model.ProductVariantChannelListing) ([]*model.ProductVariantChannelListing, *model.AppError) {
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

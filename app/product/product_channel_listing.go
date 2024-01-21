package product

import (
	"net/http"

	"github.com/samber/lo"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

// ProductChannelListingsByOption returns a list of product channel listings filtered using given option
func (a *ServiceProduct) ProductChannelListingsByOption(option *model.ProductChannelListingFilterOption) ([]*model.ProductChannelListing, *model_helper.AppError) {
	listings, err := a.srv.Store.ProductChannelListing().FilterByOption(option)
	if err != nil {
		return nil, model_helper.NewAppError("ProductChannelListingsByOption", "app.product.product_channel_listings_by_option_missing.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return listings, nil
}

// BulkUpsertProductChannelListings bulk update/inserts given product channel listings and returns them
func (a *ServiceProduct) BulkUpsertProductChannelListings(transaction *gorm.DB, listings []*model.ProductChannelListing) ([]*model.ProductChannelListing, *model_helper.AppError) {
	listings, err := a.srv.Store.ProductChannelListing().BulkUpsert(transaction, listings)
	if err != nil {
		if appErr, ok := err.(*model_helper.AppError); ok {
			return nil, appErr
		}
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		} else if _, ok := err.(*store.ErrInvalidInput); ok {
			statusCode = http.StatusBadRequest
		}

		return nil, model_helper.NewAppError("BulkUpsertProductChannelListings", "app.product.error_bulk_upserting_product_channel_listings.app_error", nil, err.Error(), statusCode)
	}

	return listings, nil
}

func (s *ServiceProduct) ValidateVariantsAvailableForPurchase(variantIds []string, channelID string) *model_helper.AppError {
	variants, err := s.srv.Store.ProductVariant().FindVariantsAvailableForPurchase(variantIds, channelID)
	if err != nil {
		return model_helper.NewAppError("ValidateVariantsAvailableForPurchase", "app.product.finding_available_for_purchase_variants.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	notAvailableVariants, _ := lo.Difference(variantIds, variants.IDs())
	if len(notAvailableVariants) > 0 {
		return model_helper.NewAppError("ValidateVariantsAvailableForPurchase", "app.product.add_unavailable_variants_to_checkout_line.app_error", nil, "cannot add lines of unavailable for purchase variants", http.StatusNotAcceptable)
	}

	return nil
}

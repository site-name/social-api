package product

import (
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/samber/lo"
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

func (s *ServiceProduct) ValidateVariantsAvailableInChannel(variantIds []string, channelId string) *model.AppError {
	variantChannelListings, appErr := s.ProductVariantChannelListingsByOption(&model.ProductVariantChannelListingFilterOption{
		Conditions: squirrel.And{
			squirrel.Eq{
				model.ProductVariantChannelListingTableName + ".VariantID": variantIds,
				model.ProductVariantChannelListingTableName + ".ChannelID": channelId,
			},
			squirrel.Expr(model.ProductVariantChannelListingTableName + ".PriceAmount IS NOT NULL"),
		},
	})
	if appErr != nil {
		return appErr
	}

	variantsNotAvailable, _ := lo.Difference(variantIds, variantChannelListings.VariantIDs())
	if len(variantsNotAvailable) > 0 {
		return model.NewAppError("ValidateVariantsAvailableInChannel", "app.product.add_not_available_variants_in_channel_to_lines.app_error", nil, "cannot add lines with unavailable variants", http.StatusNotAcceptable)
	}

	return nil
}

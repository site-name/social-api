package product

import (
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/samber/lo"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

// ProductVariantChannelListingsByOption returns a slice of product variant channel listings by given option
func (a *ServiceProduct) ProductVariantChannelListingsByOption(options *model.ProductVariantChannelListingFilterOption) (model.ProductVariantChannelListings, *model_helper.AppError) {
	listings, err := a.srv.Store.ProductVariantChannelListing().FilterbyOption(options)
	if err != nil {
		return nil, model_helper.NewAppError("ProductVariantChannelListingsByOption", "app.product_error_finding_product_variant_channel_listings_by_option.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return listings, nil
}

// BulkUpsertProductVariantChannelListings tells store to bulk upserts given product variant channel listings
func (s *ServiceProduct) BulkUpsertProductVariantChannelListings(transaction *gorm.DB, listings []*model.ProductVariantChannelListing) ([]*model.ProductVariantChannelListing, *model_helper.AppError) {
	variantChannelListings, err := s.srv.Store.ProductVariantChannelListing().BulkUpsert(transaction, listings)
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

		return nil, model_helper.NewAppError("BulkUpsertProductVariantChannelListings", "app.product.error_bulk_upserting_product_variant_channel_listings.app_error", nil, err.Error(), statusCode)
	}

	return variantChannelListings, nil
}

func (s *ServiceProduct) ValidateVariantsAvailableInChannel(variantIds []string, channelId string) *model_helper.AppError {
	variantChannelListings, appErr := s.ProductVariantChannelListingsByOption(&model.ProductVariantChannelListingFilterOption{
		Conditions: squirrel.And{
			squirrel.Eq{
				model.ProductVariantChannelListingTableName + "." + model.ProductVariantChannelListingColumnVariantID: variantIds,
				model.ProductVariantChannelListingTableName + "." + model.ProductVariantChannelListingColumnChannelID: channelId,
			},
			squirrel.Expr(model.ProductVariantChannelListingTableName + ".PriceAmount IS NOT NULL"),
		},
	})
	if appErr != nil {
		return appErr
	}

	variantsNotAvailable, _ := lo.Difference(variantIds, variantChannelListings.VariantIDs())
	if len(variantsNotAvailable) > 0 {
		return model_helper.NewAppError("ValidateVariantsAvailableInChannel", "app.product.add_not_available_variants_in_channel_to_lines.app_error", nil, "cannot add lines with unavailable variants", http.StatusNotAcceptable)
	}

	return nil
}

func (s *ServiceProduct) UpdateOrCreateProductVariantChannelListings(variantID string, inputList []model.ProductVariantChannelListingAddInput) *model_helper.AppError {
	tx := s.srv.Store.GetMaster().Begin()
	if tx.Error != nil {
		return model_helper.NewAppError("UpdateOrCreateProductVariantChannelListings", model.ErrorCreatingTransactionErrorID, nil, tx.Error.Error(), http.StatusInternalServerError)
	}

	relationsToUpsert := make(model.ProductVariantChannelListings, len(inputList))

	for idx, input := range inputList {
		existingRelations, appErr := s.ProductVariantChannelListingsByOption(&model.ProductVariantChannelListingFilterOption{
			Conditions: squirrel.Eq{
				model.ProductVariantChannelListingTableName + "." + model.ProductVariantChannelListingColumnChannelID: input.ChannelID,
				model.ProductVariantChannelListingTableName + "." + model.ProductVariantChannelListingColumnVariantID: variantID,
			},
		})
		if appErr != nil {
			return appErr
		}

		if len(existingRelations) > 0 {
			relation := existingRelations[0]
			relation.Patch(input)

			relationsToUpsert[idx] = relation
			continue
		}

		relationsToUpsert[idx] = &model.ProductVariantChannelListing{
			VariantID:                 variantID,
			ChannelID:                 input.ChannelID,
			PriceAmount:               &input.Price,
			CostPriceAmount:           input.CostPrice,
			PreorderQuantityThreshold: input.PreorderThreshold,
		}
	}

	_, appErr := s.BulkUpsertProductVariantChannelListings(tx, relationsToUpsert)
	if appErr != nil {
		return appErr
	}

	// update product discounted price
	s.srv.Go(func() {
		defer s.srv.Store.FinalizeTransaction(tx)

		product, appErr := s.ProductByOption(&model.ProductFilterOption{
			ProductVariantID: squirrel.Eq{model.ProductVariantTableName + "." + model.ProductVariantColumnId: variantID},
		})
		if appErr != nil {
			slog.Error("failed to find parent product of given variant", slog.Err(appErr))
			return
		}

		appErr = s.UpdateProductDiscountedPrice(tx, *product, []*model.DiscountInfo{})
		if appErr != nil {
			slog.Error("failed to update discounted price for parent product of given channel", slog.Err(appErr))
			return
		}

		err := tx.Commit().Error
		if err != nil {
			slog.Error("failed to commit transaction after updating discounted price for parent product of given variant", slog.Err(err))
		}
	})

	productVariant, appErr := s.ProductVariantById(variantID)
	if appErr != nil {
		return appErr
	}

	pluginMng := s.srv.PluginService().GetPluginManager()
	_, appErr = pluginMng.ProductVariantUpdated(*productVariant)
	return appErr
}

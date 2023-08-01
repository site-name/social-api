package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/samber/lo"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/web"
)

// NOTE: directives checked. Refer to ./schemas/collections.graphqls for details.
func (r *Resolver) CollectionAddProducts(ctx context.Context, args struct {
	CollectionID string
	Products     []string
}) (*CollectionAddProducts, error) {
	// validate arguments
	if !model.IsValidId(args.CollectionID) {
		return nil, model.NewAppError("CollectionAddProducts", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "collectionID"}, fmt.Sprintf("%s is invalid collection id", args.CollectionID), http.StatusBadRequest)
	}
	if !lo.EveryBy(args.Products, model.IsValidId) {
		return nil, model.NewAppError("CollectionAddProducts", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "products"}, "please provide valid product ids", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	/* check if there is at least 1 product that has no variant */
	productVariants, appErr := embedCtx.App.Srv().ProductService().ProductVariantsByOption(&model.ProductVariantFilterOption{
		Conditions: squirrel.Eq{model.ProductVariantTableName + ".ProductID": args.Products},
	})
	if appErr != nil {
		return nil, appErr
	}
	productsWithVariantsMap := lo.SliceToMap(productVariants, func(v *model.ProductVariant) (string, bool) { return v.ProductID, true })

	if !lo.SomeBy(args.Products, func(pid string) bool { return !productsWithVariantsMap[pid] }) {
		// meaning one of given products has no related productvariants
		return nil, model.NewAppError("CollectionAddProducts", "api.collection.cannot_add_products_without_variants.app_error", nil, "Cannot manage products without variants.", http.StatusBadRequest)
	}

	transaction := embedCtx.App.Srv().Store.GetMaster().Begin()
	if transaction.Error != nil {
		return nil, model.NewAppError("CollectionAddProducts", model.ErrorCreatingTransactionErrorID, nil, transaction.Error.Error(), http.StatusInternalServerError)
	}

	// add collection-product relations:
	collectionProductRels := lo.Map(args.Products, func(pid string, _ int) *model.CollectionProduct {
		return &model.CollectionProduct{CollectionID: args.CollectionID, ProductID: pid}
	})
	_, appErr = embedCtx.App.Srv().ProductService().CreateCollectionProductRelations(transaction, collectionProductRels)
	if appErr != nil {
		return nil, appErr
	}

	// check if the collection has some sale
	collectionHasSales := embedCtx.App.Srv().Store.GetReplica().
		Model(&model.Collection{Id: args.CollectionID}).
		Association("Sales").
		Count() > 0

	if collectionHasSales {
		appErr = embedCtx.App.Srv().ProductService().UpdateProductsDiscountedPricesOfCatalogues(transaction, args.Products, nil, nil, nil)
		if appErr != nil {
			return nil, appErr
		}
	}

	// commit transaction
	err := transaction.Commit().Error
	if err != nil {
		return nil, model.NewAppError("CollectionAddProducts", model.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}
	transaction.Rollback()

	// TODO: Determine if we need call plugins' product updated methods

	return &CollectionAddProducts{
		Collection: &Collection{ID: args.CollectionID},
	}, nil
}

// NOTE: directives checked. Refer to ./schemas/collections.graphqls for details.
func (r *Resolver) CollectionCreate(ctx context.Context, args struct {
	Input CollectionCreateInput
}) (*CollectionCreate, error) {
	// validate params
	if !lo.EveryBy(args.Input.Products, model.IsValidId) {
		return nil, model.NewAppError("CollectionCreate", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "products"}, "please provide valid product ids", http.StatusBadRequest)
	}

	panic(fmt.Errorf("not implemented"))
}

// NOTE: Refer to ./schemas/collection.graphqls for details on directive used.
func (r *Resolver) CollectionDelete(ctx context.Context, args struct{ Id string }) (*CollectionDelete, error) {
	// validate params
	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("CollectionDelete", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "please provide valid collection id", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	err := embedCtx.App.Srv().Store.Collection().Delete(args.Id)
	if err != nil {
		return nil, model.NewAppError("CollectionDelete", "app.product.error_deleting_collections.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	// TODO: check if we need create webhook events for products changes

	return &CollectionDelete{
		Collection: &Collection{ID: args.Id},
	}, nil
}

// NOTE: Refer to ./schemas/collection.graphqls for details on directive used.
func (r *Resolver) CollectionReorderProducts(ctx context.Context, args struct {
	CollectionID string
	Moves        []*MoveProductInput
}) (*CollectionReorderProducts, error) {
	// validate params
	if !model.IsValidId(args.CollectionID) {
		return nil, model.NewAppError("CollectionReorderProducts", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "CollectionID"}, "please provide valid collection id", http.StatusBadRequest)
	}

	panic("not implemented")
}

// NOTE: Refer to ./schemas/collection.graphqls for details on directive used.
func (r *Resolver) CollectionBulkDelete(ctx context.Context, args struct{ Ids []string }) (*CollectionBulkDelete, error) {
	// validate params
	if !lo.EveryBy(args.Ids, model.IsValidId) {
		return nil, model.NewAppError("CollectionBulkDelete", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "ids"}, "please provide valid collection ids", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	err := embedCtx.App.Srv().Store.Collection().Delete(args.Ids...)
	if err != nil {
		return nil, model.NewAppError("CollectionBulkDelete", "app.product.error_deleting_collections.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	// TODO: check if we need create webhook events for products changes
	return &CollectionBulkDelete{
		Count: int32(len(args.Ids)),
	}, nil
}

// NOTE: Refer to ./schemas/collection.graphqls for details on directive used.
func (r *Resolver) CollectionRemoveProducts(ctx context.Context, args struct {
	CollectionID string
	Products     []string
}) (*CollectionRemoveProducts, error) {
	// validate arguments
	if !model.IsValidId(args.CollectionID) {
		return nil, model.NewAppError("CollectionRemoveProducts", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "collection id"}, "please provide valid collection id", http.StatusBadRequest)
	}
	if !lo.EveryBy(args.Products, model.IsValidId) {
		return nil, model.NewAppError("CollectionRemoveProducts", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "product ids"}, "please provide valid product ids", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	err := embedCtx.App.Srv().Store.CollectionProduct().Delete(nil, &model.CollectionProductFilterOptions{
		Conditions: squirrel.Eq{
			model.CollectionProductRelationTableName + ".CollectionID": args.CollectionID,
			model.CollectionProductRelationTableName + ".ProductID":    args.Products,
		},
	})
	if err != nil {
		return nil, model.NewAppError("CollectionRemoveProducts", "app.product.remove_collection_products.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	// TODO: check if we need to call plugins' productUpdated methods

	collectionHasSales := embedCtx.App.Srv().Store.
		GetReplica().
		Model(&model.Collection{Id: args.CollectionID}).
		Association("Sales").
		Count() > 0

	if collectionHasSales {
		// Updated the db entries, recalculating discounts of affected products
		embedCtx.App.Srv().Go(func() {
			appErr := embedCtx.App.Srv().ProductService().UpdateProductsDiscountedPricesOfCatalogues(nil, args.Products, nil, nil, nil)
			if appErr != nil {
				slog.Error("failed to update product discounted prices of catalogues", slog.Err(appErr))
			}
		})
	}

	collections, appErr := embedCtx.App.Srv().ProductService().CollectionsByOption(&model.CollectionFilterOption{
		Conditions: squirrel.Eq{model.CollectionTableName + ".Id": args.CollectionID},
	})
	if appErr != nil {
		return nil, appErr
	}

	return &CollectionRemoveProducts{
		Collection: systemCollectionToGraphqlCollection(collections[0]),
	}, nil
}

// NOTE: Refer to ./schemas/collection.graphqls for details on directive used.
func (r *Resolver) CollectionUpdate(ctx context.Context, args struct {
	Id    string
	Input CollectionInput
}) (*CollectionUpdate, error) {
	// validate arguments
	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("CollectionUpdate", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "please provide valid collection id", http.StatusBadRequest)
	}

	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CollectionTranslate(ctx context.Context, args struct {
	Id           string
	Input        TranslationInput
	LanguageCode LanguageCodeEnum
}) (*CollectionTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

// NOTE: Refer to ./schemas/collection.graphqls for details on directive used.
func (r *Resolver) CollectionChannelListingUpdate(ctx context.Context, args struct {
	Id    string
	Input CollectionChannelListingUpdateInput
}) (*CollectionChannelListingUpdate, error) {
	// validate arguments
	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("CollectionChannelListingUpdate", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "please provide valid collection id", http.StatusBadRequest)
	}

	// validate input
	var addChannelIds util.AnyArray[string] = lo.Map(args.Input.AddChannels, func(item *PublishableChannelListingInput, _ int) string { return item.ChannelID })
	var removeChannelIds util.AnyArray[string] = args.Input.RemoveChannels

	if addChannelIds.InterSection(removeChannelIds).Len() > 0 {
		return nil, model.NewAppError("CollectionChannelListingUpdate", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "input"}, "some channels are both being added and removed", http.StatusBadRequest)
	}
	if addChannelIds.HasDuplicates() || !lo.EveryBy(addChannelIds, model.IsValidId) {
		return nil, model.NewAppError("CollectionChannelListingUpdate", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "add channels"}, "please provide valid channel ids and avoid duplicating", http.StatusBadRequest)
	}
	if removeChannelIds.HasDuplicates() || !lo.EveryBy(removeChannelIds, model.IsValidId) {
		return nil, model.NewAppError("CollectionChannelListingUpdate", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "remove channels"}, "please provide valid channel ids and avoid duplicating", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	// begin transaction:
	transaction := embedCtx.App.Srv().Store.GetMaster().Begin()
	if transaction.Error != nil {
		return nil, model.NewAppError("CollectionChannelListingUpdate", model.ErrorCreatingTransactionErrorID, nil, transaction.Error.Error(), http.StatusInternalServerError)
	}
	defer transaction.Rollback()

	// delete collection-channel listings
	err := embedCtx.App.Srv().Store.CollectionChannelListing().Delete(transaction, &model.CollectionChannelListingFilterOptions{
		Conditions: squirrel.Eq{
			model.CollectionChannelListingTableName + ".CollectionID": args.Id,
			model.CollectionChannelListingTableName + ".ChannelID":    removeChannelIds,
		},
	})
	if err != nil {
		return nil, model.NewAppError("CollectionChannelListingUpdate", "app.product.error_deleting_collection_channel_listings.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	// add collection-channel listings
	now := time.Now()
	collectionChannelListingsToAdd := lo.Map(args.Input.AddChannels, func(item *PublishableChannelListingInput, _ int) *model.CollectionChannelListing {
		relation := model.CollectionChannelListing{
			CollectionID: args.Id,
			ChannelID:    item.ChannelID,
		}

		if item.IsPublished != nil && *item.IsPublished {
			relation.IsPublished = true

			if item.PublicationDate == nil {
				date := util.StartOfDay(now)
				relation.PublicationDate = &date
			} else {
				relation.PublicationDate = &item.PublicationDate.Time
			}
		}

		return &relation
	})

	_, err = embedCtx.App.Srv().Store.CollectionChannelListing().Upsert(transaction, collectionChannelListingsToAdd...)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrInvalidInput); ok {
			statusCode = http.StatusBadRequest
		}
		return nil, model.NewAppError("CollectionChannelListingUpdate", "app.product.upsert_collection_channel_listings.app_error", nil, err.Error(), statusCode)
	}

	// commit transaction
	err = transaction.Commit().Error
	if err != nil {
		return nil, model.NewAppError("CollectionChannelListingUpdate", model.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	collections, appErr := embedCtx.App.Srv().ProductService().CollectionsByOption(&model.CollectionFilterOption{
		Conditions: squirrel.Eq{model.CollectionTableName + ".Id": args.Id},
	})
	if appErr != nil {
		return nil, appErr
	}

	return &CollectionChannelListingUpdate{
		Collection: systemCollectionToGraphqlCollection(collections[0]),
	}, nil
}

func (r *Resolver) Collection(ctx context.Context, args struct {
	Id      *string
	Slug    *string
	Channel *string
}) (*Collection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Collections(ctx context.Context, args struct {
	Filter  *CollectionFilterInput
	SortBy  *CollectionSortingInput
	Channel *string
	GraphqlParams
}) (*CollectionCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

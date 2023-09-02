package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"unsafe"

	"github.com/Masterminds/squirrel"
	"github.com/gosimple/slug"
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

	// validate all given products have variants
	productsWithNoVariants, appErr := embedCtx.App.Srv().ProductService().ProductsByOption(&model.ProductFilterOption{
		Conditions:           squirrel.Eq{model.ProductTableName + ".Id": args.Products},
		HasNoProductVariants: true,
	})
	if appErr != nil {
		return nil, appErr
	}
	if len(productsWithNoVariants) > 0 {
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
	appErr := args.Input.validate("CollectionCreate")
	if appErr != nil {
		return nil, appErr
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
	for _, move := range args.Moves {
		if !model.IsValidId(move.ProductID) {
			return nil, model.NewAppError("CollectionReorderProducts", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "moves"}, "please provide valid product ids for moving", http.StatusBadRequest)
		}
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	collectionProducts, appErr := embedCtx.App.Srv().ProductService().CollectionProductRelationsByOptions(&model.CollectionProductFilterOptions{
		Conditions: squirrel.Expr(model.CollectionProductRelationTableName+".CollectionID = ?", args.CollectionID),
	})
	if appErr != nil {
		return nil, appErr
	}

	// keys are product ids
	var collectionProductsMap = lo.SliceToMap(collectionProducts, func(cp *model.CollectionProduct) (string, *model.CollectionProduct) { return cp.ProductID, cp })
	var operations = map[string]*int32{}

	for _, move := range args.Moves {
		relation, found := collectionProductsMap[move.ProductID]
		if !found {
			return nil, model.NewAppError("CollectionReorderProducts", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "moves"}, "some products provided does not relate to the collection", http.StatusBadRequest)
		}

		operations[relation.Id] = move.SortOrder
	}

	// begin transaction
	tran := embedCtx.App.Srv().Store.GetMaster().Begin()
	if tran.Error != nil {
		return nil, model.NewAppError("CollectionReorderProducts", model.ErrorCreatingTransactionErrorID, nil, tran.Error.Error(), http.StatusInternalServerError)
	}
	defer tran.Rollback()

	panic("not implemented")

	return nil, nil
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

	_, collections, appErr := embedCtx.App.Srv().ProductService().CollectionsByOption(&model.CollectionFilterOption{
		Conditions: squirrel.Expr(model.CollectionTableName+".Id = ?", args.CollectionID),
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
	appErr := args.Input.validate("CollectionUpdate")
	if appErr != nil {
		return nil, appErr
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

type CollectionChannelListingUpdateArgs struct {
	Id    string
	Input CollectionChannelListingUpdateInput
}

func (args *CollectionChannelListingUpdateArgs) validate() *model.AppError {
	if !model.IsValidId(args.Id) {
		return model.NewAppError("CollectionChannelListingUpdate", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "please provide valid collection id", http.StatusBadRequest)
	}

	var addChannelIds util.AnyArray[string] = lo.Map(args.Input.AddChannels, func(item *PublishableChannelListingInput, _ int) string { return item.ChannelID })
	var removeChannelIds util.AnyArray[string] = args.Input.RemoveChannels

	if addChannelIds.InterSection(removeChannelIds).Len() > 0 {
		return model.NewAppError("CollectionChannelListingUpdate", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "input"}, "some channels are both being added and removed", http.StatusBadRequest)
	}
	if addChannelIds.HasDuplicates() || !lo.EveryBy(addChannelIds, model.IsValidId) {
		return model.NewAppError("CollectionChannelListingUpdate", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "add channels"}, "please provide valid channel ids and avoid duplicating", http.StatusBadRequest)
	}
	if removeChannelIds.HasDuplicates() || !lo.EveryBy(removeChannelIds, model.IsValidId) {
		return model.NewAppError("CollectionChannelListingUpdate", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "remove channels"}, "please provide valid channel ids and avoid duplicating", http.StatusBadRequest)
	}

	return nil
}

// NOTE: Refer to ./schemas/collection.graphqls for details on directive used.
func (r *Resolver) CollectionChannelListingUpdate(ctx context.Context, args CollectionChannelListingUpdateArgs) (*CollectionChannelListingUpdate, error) {
	// validate arguments
	appErr := args.validate()
	if appErr != nil {
		return nil, appErr
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
			model.CollectionChannelListingTableName + ".ChannelID":    args.Input.RemoveChannels,
		},
	})
	if err != nil {
		return nil, model.NewAppError("CollectionChannelListingUpdate", "app.product.error_deleting_collection_channel_listings.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	// add collection-channel listings
	today := util.StartOfDay(time.Now())

	collectionChannelListingsToAdd := lo.Map(args.Input.AddChannels, func(item *PublishableChannelListingInput, _ int) *model.CollectionChannelListing {
		relation := model.CollectionChannelListing{
			CollectionID: args.Id,
			ChannelID:    item.ChannelID,
		}

		if item.IsPublished != nil && *item.IsPublished {
			relation.IsPublished = true

			if item.PublicationDate == nil {
				relation.PublicationDate = &today
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

	_, collections, appErr := embedCtx.App.Srv().ProductService().CollectionsByOption(&model.CollectionFilterOption{
		Conditions: squirrel.Eq{model.CollectionTableName + ".Id": args.Id},
	})
	if appErr != nil {
		return nil, appErr
	}

	return &CollectionChannelListingUpdate{
		Collection: systemCollectionToGraphqlCollection(collections[0]),
	}, nil
}

type CollectionArgs struct {
	Id *string // -------|
	//                   OR
	Slug    *string // --|
	Channel *string // this is channel slug
}

func (c *CollectionArgs) validate() *model.AppError {
	if (c.Id == nil && c.Slug == nil) || (c.Id != nil && c.Slug != nil) {
		return model.NewAppError("Collection", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id, slug"}, "please provide either id or slug", http.StatusBadRequest)
	}
	if c.Id != nil && !model.IsValidId(*c.Id) {
		return model.NewAppError("Collection", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "please provide valid collection id", http.StatusBadRequest)
	}
	if c.Slug != nil && !slug.IsSlug(*c.Slug) {
		return model.NewAppError("Collection", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "slug"}, "please provide valid collection slug", http.StatusBadRequest)
	}
	if c.Channel != nil && !slug.IsSlug(*c.Channel) {
		return model.NewAppError("Collection", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "channel"}, "please provide valid channel slug", http.StatusBadRequest)
	}

	return nil
}

func (r *Resolver) Collection(ctx context.Context, args CollectionArgs) (*Collection, error) {
	// validate params
	appErr := args.validate()
	if appErr != nil {
		return nil, appErr
	}

	var (
		collectionID   string
		collectionSlug string
		channelSlug    string
	)
	if args.Id != nil {
		collectionID = *args.Id
	} else if args.Slug != nil {
		collectionSlug = *args.Slug
	}

	if args.Channel != nil {
		channelSlug = *args.Channel
	}

	// check if requester can see all collections or not
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	embedCtx.CheckAuthenticatedAndHasRoleAny("Collection", model.ShopAdminRoleId, model.ShopStaffRoleId)
	userIsShopStaff := embedCtx.Err == nil

	collections, appErr := embedCtx.App.Srv().ProductService().VisibleCollectionsToUser(channelSlug, userIsShopStaff)
	if appErr != nil {
		return nil, appErr
	}

	for _, collection := range collections {
		if collection.Id == collectionID || collection.Slug == collectionSlug {
			return systemCollectionToGraphqlCollection(collection), nil
		}
	}

	return nil, nil
}

type CollectionsArgs struct {
	Filter  *CollectionFilterInput
	SortBy  *CollectionSortingInput
	Channel *string // channel slug
	GraphqlParams
}

func (c *CollectionsArgs) parse(embedCtx *web.Context) (*model.CollectionFilterOption, *model.AppError) {
	// validate params
	if c.Filter != nil {
		appErr := c.Filter.validate("Collections")
		if appErr != nil {
			return nil, appErr
		}
	}
	if c.SortBy != nil && !c.SortBy.Field.IsValid() {
		return nil, model.NewAppError("Collections", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "SortField"}, "please provide valid sort field", http.StatusBadRequest)
	}
	if c.Channel != nil && !slug.IsSlug(*c.Channel) {
		return nil, model.NewAppError("Collections", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Channel"}, "please provide valid channel slug", http.StatusBadRequest)
	}
	appErr := c.GraphqlParams.validate("Collections")
	if appErr != nil {
		return nil, appErr
	}

	// parse pagination
	paginationValues, _ := c.GraphqlParams.Parse("Collections")

	var res = &model.CollectionFilterOption{
		CountTotal:              true,
		GraphqlPaginationValues: *paginationValues,
	}

	conditions := squirrel.And{}

	if c.Filter != nil {
		if len(c.Filter.Ids) > 0 {
			conditions = append(conditions, squirrel.Eq{model.CollectionTableName + ".Id": c.Filter.Ids})
		}

		if c.Filter.Published != nil {
			switch *c.Filter.Published {
			case CollectionPublishedPublished:
				res.ChannelListingIsPublished = squirrel.Expr(model.CollectionChannelListingTableName + ".IsPublished")
			case CollectionPublishedHidden:
				res.ChannelListingIsPublished = squirrel.Expr(model.CollectionChannelListingTableName + ".IsPublished = false")
			}

			if c.Channel != nil {
				res.ChannelListingChannelSlug = squirrel.Expr(model.ChannelTableName+".Slug = ?", *c.Channel)
			}
		}

		// search
		if c.Filter.Search != nil {
			expr := "%" + *c.Filter.Search + "%"
			conditions = append(conditions,
				squirrel.Or{
					squirrel.Expr(model.CollectionTableName+".Name ILIKE ?", expr),
					squirrel.Expr(model.CollectionTableName+".Slug ILIKE ?", expr),
				})
		}

		// meta data
		if c.Filter.Metadata != nil {
			for _, metaItem := range c.Filter.Metadata {
				if metaItem != nil && metaItem.Key != "" {
					if metaItem.Value == "" {
						expr := fmt.Sprintf(model.SaleTableName+".Metadata::jsonb ? '%s'", metaItem.Key)
						conditions = append(conditions, squirrel.Expr(expr))
					} else {
						expr := fmt.Sprintf(model.SaleTableName+".Metadata::jsonb @> '{%q:%q}'", metaItem.Key, metaItem.Value)
						conditions = append(conditions, squirrel.Expr(expr))
					}
				}
			}
		}
	}

	res.Conditions = conditions

	if res.GraphqlPaginationValues.OrderBy == "" {
		sortfields := collectionSortFieldMap[CollectionSortFieldName].fields

		// parse sort
		if c.SortBy != nil {
			sortfields = collectionSortFieldMap[c.SortBy.Field].fields

			switch c.SortBy.Field {
			case CollectionSortFieldAvailability, CollectionSortFieldPublicationDate:
				if c.Channel != nil {
					res.ChannelSlugForIsPublishedAndPublicationDateAnnotation = *c.Channel
				} else {
					// we need a channel slug to be able to annotate collection's Availability and publicationDate
					defaultChannel, appErr := embedCtx.App.Srv().ChannelService().GetDefaultChannel()
					if appErr != nil {
						// this is because there is no active channel in the system
						return nil, appErr
					}
					res.ChannelSlugForIsPublishedAndPublicationDateAnnotation = defaultChannel.Slug
				}
				if c.SortBy.Field == CollectionSortFieldAvailability {
					res.AnnotateIsPublished = true
				} else {
					res.AnnotatePublicationDate = true
				}

			case CollectionSortFieldProductCount:
				res.AnnotateProductCount = true
			}
		}

		orderDirection := c.GraphqlParams.orderDirection()
		res.GraphqlPaginationValues.OrderBy = sortfields.
			Map(func(_ int, item string) string { return item + " " + orderDirection }).
			Join(", ")
	}

	return res, nil
}

func (r *Resolver) Collections(ctx context.Context, args CollectionsArgs) (*CollectionCountableConnection, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	collectionFilterOpts, appErr := args.parse(embedCtx)
	if appErr != nil {
		return nil, appErr
	}

	panic("not done")
	// add filter when requester is shop staff or outer customer

	totalCount, collections, appErr := embedCtx.App.Srv().ProductService().CollectionsByOption(collectionFilterOpts)
	hasNextPage, hasPrevPage := args.checkNextPageAndPreviousPage(len(collections))
	keyFunc := collectionSortFieldMap[CollectionSortFieldName].keyFunc

	if args.SortBy != nil {
		keyFunc = collectionSortFieldMap[args.SortBy.Field].keyFunc
	}
	res := constructCountableConnection(collections, totalCount, hasNextPage, hasPrevPage, keyFunc, systemCollectionToGraphqlCollection)
	return (*CollectionCountableConnection)(unsafe.Pointer(res)), nil
}

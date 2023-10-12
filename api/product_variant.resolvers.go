package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"net/http"
	"unsafe"

	"github.com/Masterminds/squirrel"
	"github.com/samber/lo"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/web"
)

// NOTE: Refer to ./graphql/schemas/product_variant.graphqls for details on directives used.
func (r *Resolver) VariantMediaUnassign(ctx context.Context, args struct {
	MediaID   UUID
	VariantID UUID
}) (*VariantMediaUnassign, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	// create tx:
	tx := embedCtx.App.Srv().Store.GetMaster().Begin()
	if tx.Error != nil {
		return nil, model.NewAppError("VariantMediaUnassign", model.ErrorCreatingTransactionErrorID, nil, tx.Error.Error(), http.StatusInternalServerError)
	}
	defer tx.Rollback()

	err := embedCtx.App.Srv().Store.
		ProductVariant().
		ToggleProductVariantRelations(
			tx,
			model.ProductVariants{{Id: args.VariantID.String()}},
			model.ProductMedias{{Id: args.MediaID.String()}},
			nil,
			nil,
			nil,
			true,
		)
	if err != nil {
		return nil, model.NewAppError("VariantMediaUnassign", "app.product.delete_variant_media_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, model.NewAppError("VariantMediaUnassign", model.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	// TODO: check if this logic is needed, since the system doesn't have support for wekhook yet

	// pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()
	// _, appErr = pluginMng.ProductVariantUpdated(*productVariant)
	// if appErr != nil {
	// 	return nil, appErr
	// }

	return &VariantMediaUnassign{
		ProductVariant: &ProductVariant{ID: args.VariantID.String()},
		Media:          &ProductMedia{ID: args.MediaID.String()},
	}, nil
}

// NOTE: Refer to ./graphql/schemas/product_variant.graphqls for details on directives used.
func (r *Resolver) VariantMediaAssign(ctx context.Context, args struct {
	MediaID   UUID
	VariantID UUID
}) (*VariantMediaAssign, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	productVariant, appErr := embedCtx.App.Srv().ProductService().ProductVariantById(args.VariantID.String())
	if appErr != nil {
		return nil, appErr
	}

	productMedias, appErr := embedCtx.App.Srv().ProductService().ProductMediasByOption(&model.ProductMediaFilterOption{
		Conditions: squirrel.Eq{model.ProductMediaTableName + ".Id": args.MediaID},
	})
	if appErr != nil {
		return nil, appErr
	}
	if len(productMedias) == 0 {
		return nil, model.NewAppError("VariantMediaAssign", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "mediaID"}, "please provide valid product media id", http.StatusBadRequest)
	}
	media := productMedias[0]

	// check if the given image and variant can be matched together
	if media.ProductID != productVariant.ProductID {
		return nil, model.NewAppError("VariantMediaAssign", "app.product.product_does_not_own_media.app_error", nil, "This media doesn't belong to that product.", http.StatusNotAcceptable)
	}

	// create tx:
	tx := embedCtx.App.Srv().Store.GetMaster().Begin()
	if tx.Error != nil {
		return nil, model.NewAppError("VariantMediaAssign", model.ErrorCreatingTransactionErrorID, nil, tx.Error.Error(), http.StatusInternalServerError)
	}
	defer tx.Rollback()

	err := embedCtx.App.Srv().Store.
		ProductVariant().
		ToggleProductVariantRelations(
			tx,
			model.ProductVariants{{Id: args.VariantID.String()}},
			model.ProductMedias{{Id: args.MediaID.String()}},
			nil,
			nil,
			nil,
			false,
		)
	if err != nil {
		return nil, model.NewAppError("VariantMediaAssign", "app.product.upsert_variant_media.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	// commit tx
	err = tx.Commit().Error
	if err != nil {
		return nil, model.NewAppError("VariantMediaAssign", model.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()
	_, appErr = pluginMng.ProductVariantUpdated(*productVariant)
	if appErr != nil {
		return nil, appErr
	}

	return &VariantMediaAssign{
		ProductVariant: SystemProductVariantToGraphqlProductVariant(productVariant),
		Media:          systemProductMediaToGraphqlProductMedia(media),
	}, nil
}

// NOTE: Refer to ./graphql/schemas/product_variant.graphqls for details on directives used.
func (r *Resolver) ProductVariantReorder(ctx context.Context, args struct {
	Moves     []*ReorderInput
	ProductID UUID
}) (*ProductVariantReorder, error) {
	panic(fmt.Errorf("not implemented"))
}

// NOTE: Refer to ./graphql/schemas/product_variant.graphqls for details on directives used.
func (r *Resolver) ProductVariantCreate(ctx context.Context, args struct {
	Input ProductVariantCreateInput
}) (*ProductVariantCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

// NOTE: Refer to ./graphql/schemas/product_variant.graphqls for details on directives used.
func (r *Resolver) ProductVariantUpdate(ctx context.Context, args struct {
	Id    UUID
	Input ProductVariantInput
}) (*ProductVariantUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

// NOTE: Refer to ./graphql/schemas/product_variant.graphqls for details on directives used.
func (r *Resolver) ProductVariantDelete(ctx context.Context, args struct{ Id UUID }) (*ProductVariantDelete, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	// delete variants, delete related draft order lines, create order events
	_, appErr := embedCtx.App.Srv().ProductService().DeleteProductVariants([]string{args.Id.String()}, embedCtx.AppContext.Session().UserId)
	if appErr != nil {
		return nil, appErr
	}

	return &ProductVariantDelete{
		ProductVariant: &ProductVariant{ID: args.Id.String()},
	}, nil
}

// NOTE: Refer to ./graphql/schemas/product_variant.graphqls for details on directives used.
func (r *Resolver) ProductVariantBulkCreate(ctx context.Context, args struct {
	Product  UUID
	Variants []*ProductVariantBulkCreateInput
}) (*ProductVariantBulkCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

// NOTE: Refer to ./graphql/schemas/product_variant.graphqls for details on directives used.
func (r *Resolver) ProductVariantBulkDelete(ctx context.Context, args struct{ Ids []UUID }) (*ProductVariantBulkDelete, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	// delete variants, delete related draft order lines, create order events
	stringIds := *(*[]string)(unsafe.Pointer(&args.Ids))
	numDeleted, appErr := embedCtx.App.Srv().ProductService().DeleteProductVariants(stringIds, embedCtx.AppContext.Session().UserId)
	if appErr != nil {
		return nil, appErr
	}

	return &ProductVariantBulkDelete{
		Count: int32(numDeleted),
	}, nil
}

// NOTE: Refer to ./graphql/schemas/product_variant.graphqls for details on directives used.
func (r *Resolver) ProductVariantStocksDelete(ctx context.Context, args struct {
	VariantID    UUID
	WarehouseIds []UUID
}) (*ProductVariantStocksDelete, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	strWarehouseIDs := *(*[]string)(unsafe.Pointer(&args.WarehouseIds))

	_, appErr := embedCtx.App.Srv().WarehouseService().DeleteStocks(&model.StockFilterOption{
		Conditions: squirrel.Eq{
			model.StockTableName + "." + model.StockColumnProductVariantID: args.VariantID,
			model.StockTableName + "." + model.StockColumnWarehouseID:      strWarehouseIDs,
		},
	})
	if appErr != nil {
		return nil, appErr
	}

	return &ProductVariantStocksDelete{
		ProductVariant: &ProductVariant{ID: args.VariantID.String()},
	}, nil
}

// NOTE: Refer to ./graphql/schemas/product_variant.graphqls for details on directives used.
func (r *Resolver) ProductVariantStocksCreate(ctx context.Context, args struct {
	Stocks    []StockInput
	VariantID UUID
}) (*ProductVariantStocksCreate, error) {
	if len(args.Stocks) == 0 {
		return nil, model.NewAppError("ProductVariantStocksCreate", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Stocks"}, "please provide stock information to create", http.StatusBadRequest)
	}

	// validate duplication:
	warehouseIdsForStocksMap := map[string]bool{} // keys are warehouse ids
	stocksToCreate := make(model.Stocks, len(args.Stocks))
	strVariantID := args.VariantID.String()

	for idx, stockInput := range args.Stocks {
		strWarehouseID := stockInput.Warehouse.String()
		if warehouseIdsForStocksMap[strWarehouseID] {
			return nil, model.NewAppError("ProductVariantStocksCreate", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Stocks"}, "please provide unique warehouse ids", http.StatusBadRequest)
		}

		warehouseIdsForStocksMap[strWarehouseID] = true
		stocksToCreate[idx] = &model.Stock{
			WarehouseID:      strWarehouseID,
			ProductVariantID: strVariantID,
			Quantity:         int(stockInput.Quantity),
		}
	}

	// check if given variant does already have stocks that belong to any of given warehouses
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	_, stocks, appErr := embedCtx.App.Srv().WarehouseService().StocksByOption(&model.StockFilterOption{
		Conditions: squirrel.Eq{
			model.StockTableName + "." + model.StockColumnWarehouseID:      lo.Keys(warehouseIdsForStocksMap),
			model.StockTableName + "." + model.StockColumnProductVariantID: args.VariantID,
		},
	})
	if appErr != nil {
		return nil, appErr
	}

	existedStock, found := lo.Find(stocks, func(item *model.Stock) bool { return item != nil && warehouseIdsForStocksMap[item.WarehouseID] })
	if found {
		return nil, model.NewAppError("ProductVariantStocksCreate", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "stocks"}, fmt.Sprintf("there is already a stock in warehouse with id=%s existed for given variant", existedStock.WarehouseID), http.StatusBadRequest)
	}

	// begin tx
	tx := embedCtx.App.Srv().Store.GetMaster().Begin()
	if tx.Error != nil {
		return nil, model.NewAppError("ProductVariantStocksCreate", model.ErrorCreatingTransactionErrorID, nil, tx.Error.Error(), http.StatusInternalServerError)
	}
	defer embedCtx.App.Srv().Store.FinalizeTransaction(tx)

	stocks, appErr = embedCtx.App.Srv().WarehouseService().BulkUpsertStocks(tx, stocksToCreate)
	if appErr != nil {
		return nil, appErr
	}

	// commit tx
	if err := tx.Commit().Error; err != nil {
		return nil, model.NewAppError("ProductVariantStocksCreate", model.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()
	for _, stock := range stocks {
		appErr = pluginMng.ProductVariantBackInStock(*stock)
		if appErr != nil {
			return nil, appErr
		}
	}

	return &ProductVariantStocksCreate{
		ProductVariant: &ProductVariant{ID: args.VariantID.String()},
	}, nil
}

// NOTE: Refer to ./graphql/schemas/product_variant.graphqls for details on directives used.
func (r *Resolver) ProductVariantStocksUpdate(ctx context.Context, args struct {
	Stocks    []StockInput
	VariantID UUID
}) (*ProductVariantStocksUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

// NOTE: Refer to ./graphql/schemas/product_variant.graphqls for details on directives used.
func (r *Resolver) ProductVariantSetDefault(ctx context.Context, args struct {
	ProductID UUID
	VariantID UUID
}) (*ProductVariantSetDefault, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	updatedProduct, appErr := embedCtx.App.Srv().ProductService().SetDefaultProductVariantForProduct(args.ProductID.String(), args.VariantID.String())
	if appErr != nil {
		return nil, appErr
	}

	return &ProductVariantSetDefault{
		Product: SystemProductToGraphqlProduct(updatedProduct),
	}, nil
}

// NOTE: Refer to ./graphql/schemas/product_variant.graphqls for details on directives used.
func (r *Resolver) ProductVariantTranslate(ctx context.Context, args struct {
	Id           UUID
	Input        NameTranslationInput
	LanguageCode LanguageCodeEnum
}) (*ProductVariantTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

// NOTE: Refer to ./graphql/schemas/product_variant.graphqls for details on directives used.
func (r *Resolver) ProductVariantChannelListingUpdate(ctx context.Context, args struct {
	Id    UUID
	Input []ProductVariantChannelListingAddInput
}) (*ProductVariantChannelListingUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

// NOTE: Refer to ./graphql/schemas/product_variant.graphqls for details on directives used.
func (r *Resolver) ProductVariantReorderAttributeValues(ctx context.Context, args struct {
	AttributeID UUID
	Moves       []*ReorderInput
	VariantID   UUID
}) (*ProductVariantReorderAttributeValues, error) {
	panic(fmt.Errorf("not implemented"))
}

// NOTE: Refer to ./graphql/schemas/product_variant.graphqls for details on directives used.
func (r *Resolver) ProductVariant(ctx context.Context, args struct {
	Id      *UUID
	Sku     *string
	Channel *string
}) (*ProductVariant, error) {
	panic(fmt.Errorf("not implemented"))
}

// NOTE: Refer to ./graphql/schemas/product_variant.graphqls for details on directives used.
func (r *Resolver) ProductVariants(ctx context.Context, args struct {
	Ids     []UUID
	Channel *string
	Filter  *ProductVariantFilterInput
	GraphqlParams
}) (*ProductVariantCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

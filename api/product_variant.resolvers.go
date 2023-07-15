package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/web"
)

// NOTE: Refer to ./schemas/product_variant.graphqls for details on directives used.
func (r *Resolver) VariantMediaUnassign(ctx context.Context, args struct {
	MediaID   string
	VariantID string
}) (*VariantMediaUnassign, error) {
	// validate params
	if !model.IsValidId(args.MediaID) {
		return nil, model.NewAppError("VariantMediaUnassign", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "MediaID"}, "please provide valid media id", http.StatusBadRequest)
	}
	if !model.IsValidId(args.VariantID) {
		return nil, model.NewAppError("VariantMediaUnassign", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "VariantID"}, "please provide valid variant id", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	// create transaction:
	transaction, err := embedCtx.App.Srv().Store.GetMasterX().Beginx()
	if err != nil {
		return nil, model.NewAppError("VariantMediaUnassign", app.ErrorCreatingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}
	defer store.FinalizeTransaction(transaction)

	// NOTE: delete does not return error on wrong values provided.
	err = embedCtx.App.Srv().Store.VariantMedia().Delete(transaction, &model.VariantMediaFilterOptions{
		Conditions: squirrel.And{
			squirrel.Eq{model.ProductVariantMediaTableName + ".VariantID": args.VariantID},
			squirrel.Eq{model.ProductVariantMediaTableName + ".MediaID": args.MediaID},
		},
	})
	if err != nil {
		return nil, model.NewAppError("VariantMediaUnassign", "app.product.delete_variant_media_by_options.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	// TODO: check if this logic is needed, since the system doesn't have support for wekhook yet

	// pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()
	// _, appErr = pluginMng.ProductVariantUpdated(*productVariant)
	// if appErr != nil {
	// 	return nil, appErr
	// }

	return &VariantMediaUnassign{
		ProductVariant: &ProductVariant{ID: args.VariantID},
		Media:          &ProductMedia{ID: args.MediaID},
	}, nil
}

// NOTE: Refer to ./schemas/product_variant.graphqls for details on directives used.
func (r *Resolver) VariantMediaAssign(ctx context.Context, args struct {
	MediaID   string
	VariantID string
}) (*VariantMediaAssign, error) {
	// validate params
	if !model.IsValidId(args.MediaID) {
		return nil, model.NewAppError("VariantMediaAssign", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "MediaID"}, "please provide valid media id", http.StatusBadRequest)
	}
	if !model.IsValidId(args.VariantID) {
		return nil, model.NewAppError("VariantMediaAssign", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "VariantID"}, "please provide valid variant id", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	productVariant, appErr := embedCtx.App.Srv().ProductService().ProductVariantById(args.VariantID)
	if appErr != nil {
		return nil, appErr
	}

	productMedias, appErr := embedCtx.App.Srv().ProductService().ProductMediasByOption(&model.ProductMediaFilterOption{
		Conditions: squirrel.Eq{model.ProductMediaTableName + ".Id": args.MediaID},
	})
	if appErr != nil {
		// NOTE: This appError covers 404 code also so no need to worry if productMedias is empty
		return nil, appErr
	}
	media := productMedias[0]

	// create transaction:
	transaction, err := embedCtx.App.Srv().Store.GetMasterX().Beginx()
	if err != nil {
		return nil, model.NewAppError("VariantMediaAssign", app.ErrorCreatingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}
	defer store.FinalizeTransaction(transaction)

	if media != nil && productVariant != nil {
		// check if the given image and variant can be matched together
		if media.ProductID == productVariant.ProductID {
			_, err := embedCtx.App.Srv().Store.VariantMedia().Upsert(transaction, &model.VariantMedia{VariantID: args.VariantID, MediaID: args.MediaID})
			if err != nil {
				// NOTE: SQL logic for handle uniqueness violation is added in upsert process.
				// So this error is internal system error
				return nil, model.NewAppError("VariantMediaAssign", "app.product.upsert_variant_media.app_error", nil, err.Error(), http.StatusInternalServerError)
			}
		} else {
			return nil, model.NewAppError("VariantMediaAssign", "app.product.product_does_not_own_media.app_error", nil, "This media doesn't belong to that product.", http.StatusNotAcceptable)
		}
	}

	// commit transaction
	err = transaction.Commit()
	if err != nil {
		return nil, model.NewAppError("VariantMediaAssign", app.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	// TODO: check if this logic is needed, since the system doesn't have support for wekhook yet

	// pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()
	// _, appErr = pluginMng.ProductVariantUpdated(*productVariant)
	// if appErr != nil {
	// 	return nil, appErr
	// }

	return &VariantMediaAssign{
		ProductVariant: SystemProductVariantToGraphqlProductVariant(productVariant),
		Media:          systemProductMediaToGraphqlProductMedia(media),
	}, nil
}

func (r *Resolver) ProductVariantReorder(ctx context.Context, args struct {
	Moves     []*ReorderInput
	ProductID string
}) (*ProductVariantReorder, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductVariantCreate(ctx context.Context, args struct {
	Input ProductVariantCreateInput
}) (*ProductVariantCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductVariantDelete(ctx context.Context, args struct{ Id string }) (*ProductVariantDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductVariantBulkCreate(ctx context.Context, args struct {
	Product  string
	Variants []*ProductVariantBulkCreateInput
}) (*ProductVariantBulkCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductVariantBulkDelete(ctx context.Context, args struct{ Ids []string }) (*ProductVariantBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductVariantStocksCreate(ctx context.Context, args struct {
	Stocks    []StockInput
	VariantID string
}) (*ProductVariantStocksCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductVariantStocksDelete(ctx context.Context, args struct {
	VariantID    string
	WarehouseIds []string
}) (*ProductVariantStocksDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductVariantStocksUpdate(ctx context.Context, args struct {
	Stocks    []StockInput
	VariantID string
}) (*ProductVariantStocksUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductVariantUpdate(ctx context.Context, args struct {
	Id    string
	Input ProductVariantInput
}) (*ProductVariantUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductVariantSetDefault(ctx context.Context, args struct {
	ProductID string
	VariantID string
}) (*ProductVariantSetDefault, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductVariantTranslate(ctx context.Context, args struct {
	Id           string
	Input        NameTranslationInput
	LanguageCode LanguageCodeEnum
}) (*ProductVariantTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductVariantChannelListingUpdate(ctx context.Context, args struct {
	Id    string
	Input []ProductVariantChannelListingAddInput
}) (*ProductVariantChannelListingUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductVariantReorderAttributeValues(ctx context.Context, args struct {
	AttributeID string
	Moves       []*ReorderInput
	VariantID   string
}) (*ProductVariantReorderAttributeValues, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductVariant(ctx context.Context, args struct {
	Id      *string
	Sku     *string
	Channel *string
}) (*ProductVariant, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductVariants(ctx context.Context, args struct {
	Ids     []string
	Channel *string
	Filter  *ProductVariantFilterInput
	GraphqlParams
}) (*ProductVariantCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

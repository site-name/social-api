package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"net/http"
	"unsafe"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/web"
)

// NOTE: Please refer to ./graphql/schemas/product_media.graphqls for details on directives used
func (r *Resolver) ProductMediaCreate(ctx context.Context, args struct {
	Input ProductMediaCreateInput
}) (*ProductMediaCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

// NOTE: Please refer to ./graphql/schemas/product_media.graphqls for details on directives used
func (r *Resolver) ProductMediaDelete(ctx context.Context, args struct{ Id UUID }) (*ProductMediaDelete, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	productMedias, appErr := embedCtx.App.Srv().ProductService().ProductMediasByOption(&model.ProductMediaFilterOption{
		Conditions: squirrel.Expr(model.ProductMediaTableName+".Id = ?", args.Id),
		Preloads:   []string{"Product"},
	})
	if appErr != nil {
		return nil, appErr
	}
	if len(productMedias) == 0 {
		return nil, model_helper.NewAppError("ProductMediaDelete", model_helper.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Id"}, "please provide valid product media id", http.StatusBadRequest)
	}

	_, appErr = embedCtx.App.Srv().ProductService().DeleteProductMedias(nil, []string{args.Id.String()})
	if appErr != nil {
		return nil, appErr
	}

	media := productMedias[0]
	product := media.Product

	if product != nil {
		pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()
		_, appErr = pluginMng.ProductUpdated(*product)
		if appErr != nil {
			return nil, appErr
		}
	}

	return &ProductMediaDelete{
		Product: SystemProductToGraphqlProduct(product),
		Media:   systemProductMediaToGraphqlProductMedia(media),
	}, nil
}

// NOTE: Please refer to ./graphql/schemas/product_media.graphqls for details on directives used
func (r *Resolver) ProductMediaBulkDelete(ctx context.Context, args struct{ Ids []UUID }) (*ProductMediaBulkDelete, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	ids := *(*[]string)(unsafe.Pointer(&args.Ids))

	numDeleted, appErr := embedCtx.App.Srv().ProductService().DeleteProductMedias(nil, ids)
	if appErr != nil {
		return nil, appErr
	}

	return &ProductMediaBulkDelete{
		Count: *(*int32)(unsafe.Pointer(&numDeleted)),
	}, nil
}

// NOTE: Please refer to ./graphql/schemas/product_media.graphqls for details on directives used
func (r *Resolver) ProductMediaReorder(ctx context.Context, args struct {
	ProductID UUID
}) (*ProductMediaReorder, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	medias, appErr := embedCtx.App.Srv().ProductService().ProductMediasByOption(&model.ProductMediaFilterOption{
		Conditions: squirrel.Expr(model.ProductMediaTableName+".ProductID = ?", args.ProductID),
	})
	if appErr != nil {
		return nil, appErr
	}

	if len(medias) == 0 {
		return nil, model_helper.NewAppError("ProductMediaReorder", model_helper.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "ProductID"}, "given product has no related product medias", http.StatusBadRequest)
	}

	for idx, media := range medias {
		media.SortOrder = &idx
	}

	medias, appErr = embedCtx.App.Srv().ProductService().UpsertProductMedias(nil, medias)
	if appErr != nil {
		return nil, appErr
	}

	product, appErr := embedCtx.App.Srv().ProductService().ProductById(args.ProductID.String())
	if appErr != nil {
		return nil, appErr
	}

	pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()
	_, appErr = pluginMng.ProductUpdated(*product)
	if appErr != nil {
		return nil, appErr
	}

	return &ProductMediaReorder{
		Product: SystemProductToGraphqlProduct(product),
		Media:   systemRecordsToGraphql(medias, systemProductMediaToGraphqlProductMedia),
	}, nil
}

// NOTE: Please refer to ./graphql/schemas/product_media.graphqls for details on directives used
func (r *Resolver) ProductMediaUpdate(ctx context.Context, args struct {
	Id    UUID
	Input ProductMediaUpdateInput
}) (*ProductMediaUpdate, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	productMedias, appErr := embedCtx.App.Srv().ProductService().ProductMediasByOption(&model.ProductMediaFilterOption{
		Conditions: squirrel.Expr(model.ProductMediaTableName+".Id = ?", args.Id),
		Preloads:   []string{"Product"},
	})
	if appErr != nil {
		return nil, appErr
	}
	if len(productMedias) == 0 {
		return nil, model_helper.NewAppError("ProductMediaUpdate", model_helper.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Id"}, "please provide valid product media id", http.StatusBadRequest)
	}

	media := productMedias[0]
	media.Alt = args.Input.Alt

	updatedMedias, appErr := embedCtx.App.Srv().ProductService().UpsertProductMedias(nil, model.ProductMedias{media})
	if appErr != nil {
		return nil, appErr
	}

	return &ProductMediaUpdate{
		Product: SystemProductToGraphqlProduct(media.Product),
		Media:   systemProductMediaToGraphqlProductMedia(updatedMedias[0]),
	}, nil
}

package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"net/http"
	"unsafe"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/web"
)

// NOTE: Refer to ./schemas/digital_content.graphqls for details on directive used
func (r *Resolver) DigitalContentCreate(ctx context.Context, args struct {
	Input     DigitalContentUploadInput
	VariantID string
}) (*DigitalContentCreate, error) {
	// validate params
	if !model_helper.IsValidId(args.VariantID) {
		return nil, model_helper.NewAppError("DigitalContentCreate", model_helper.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "VariantID"}, "please provide valid variant id", http.StatusBadRequest)
	}

	panic("not implemented")
}

// NOTE: Refer to ./schemas/digital_content.graphqls for details on directive used
func (r *Resolver) DigitalContentDelete(ctx context.Context, args struct{ VariantID string }) (*DigitalContentDelete, error) {
	// validate params
	if !model_helper.IsValidId(args.VariantID) {
		return nil, model_helper.NewAppError("DigitalContentDelete", model_helper.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "VariantID"}, "please provide valid variant id", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	err := embedCtx.App.Srv().Store.DigitalContent().Delete(nil, &model.DigitalContentFilterOption{
		Conditions: squirrel.Expr(model.DigitalContentTableName+".ProductVariantID = ?", args.VariantID),
	})
	if err != nil {
		return nil, model_helper.NewAppError("DigitalContentDelete", "app.product.error_delete_digital_content.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return &DigitalContentDelete{
		Variant: &ProductVariant{ID: args.VariantID},
	}, nil
}

// NOTE: Refer to ./schemas/digital_content.graphqls for details on directive used
func (r *Resolver) DigitalContentUpdate(ctx context.Context, args struct {
	Input     DigitalContentInput
	VariantID string
}) (*DigitalContentUpdate, error) {
	// validate params
	if !model_helper.IsValidId(args.VariantID) {
		return nil, model_helper.NewAppError("DigitalContentUpdate", model_helper.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "VariantID"}, "please provide valid variant id", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	content, appErr := embedCtx.App.Srv().ProductService().DigitalContentbyOption(&model.DigitalContentFilterOption{
		Conditions: squirrel.Expr(model.DigitalContentTableName+".ProductVariantID = ?", args.VariantID),
	})
	if appErr != nil {
		return nil, appErr
	}

	content.UseDefaultSettings = &args.Input.UseDefaultSettings

	// clean input
	switch {
	case args.Input.MaxDownloads == nil:
		return nil, model_helper.NewAppError("DigitalContentUpdate", model_helper.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "MaxDownloads"}, "please provide MaxDownloads", http.StatusBadRequest)
	case args.Input.URLValidDays == nil:
		return nil, model_helper.NewAppError("DigitalContentUpdate", model_helper.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "URLValidDays"}, "please provide URLValidDays", http.StatusBadRequest)
	case args.Input.AutomaticFulfillment == nil:
		return nil, model_helper.NewAppError("DigitalContentUpdate", model_helper.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "AutomaticFulfillment"}, "please provide AutomaticFulfillment", http.StatusBadRequest)
	}

	content.MaxDownloads = (*int)(unsafe.Pointer(args.Input.MaxDownloads))
	content.UrlValidDays = (*int)(unsafe.Pointer(args.Input.URLValidDays))
	content.AutomaticFulfillment = args.Input.AutomaticFulfillment

	content, appErr = embedCtx.App.Srv().ProductService().UpsertDigitalContent(content)
	if appErr != nil {
		return nil, appErr
	}

	return &DigitalContentUpdate{
		Variant: &ProductVariant{ID: content.ProductVariantID},
		Content: systemDigitalContentToGraphqlDigitalContent(content),
	}, nil
}

// NOTE: Refer to ./schemas/digital_content.graphqls for details on directive used
func (r *Resolver) DigitalContentURLCreate(ctx context.Context, args struct {
	Input DigitalContentURLCreateInput
}) (*DigitalContentURLCreate, error) {
	if !model_helper.IsValidId(args.Input.Content) {
		return nil, model_helper.NewAppError("DigitalContentURLCreate", model_helper.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Content"}, "please provide valid digital content id", http.StatusBadRequest)
	}

	contentUrl := &model.DigitalContentUrl{
		ContentID: args.Input.Content,
	}
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	contentUrl, appErr := embedCtx.App.Srv().ProductService().UpsertDigitalContentURL(contentUrl)
	if appErr != nil {
		return nil, appErr
	}

	return &DigitalContentURLCreate{
		DigitalContentURL: systemDigitalContentURLToGraphqlDigitalContentURL(contentUrl),
	}, nil
}

// NOTE: Refer to ./schemas/digital_content.graphqls for details on directive used
// TODO: check if we need permissions to see this.
func (r *Resolver) DigitalContent(ctx context.Context, args struct{ Id string }) (*DigitalContent, error) {
	if !model_helper.IsValidId(args.Id) {
		return nil, model_helper.NewAppError("DigitalContentURLCreate", model_helper.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Content"}, "please provide valid digital content id", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	content, appErr := embedCtx.App.Srv().ProductService().DigitalContentbyOption(&model.DigitalContentFilterOption{
		Conditions: squirrel.Expr(model.DigitalContentTableName+".Id = ?", args.Id),
	})
	if appErr != nil {
		return nil, appErr
	}

	return systemDigitalContentToGraphqlDigitalContent(content), nil
}

// NOTE: Refer to ./schemas/digital_content.graphqls for details on directive used
// TODO: check if we need permissions to see this.
// NOTE: Digital contents are sort by `ContentType` ASC, `ContentFile` ASC
func (r *Resolver) DigitalContents(ctx context.Context, args GraphqlParams) (*DigitalContentCountableConnection, error) {
	graphqlPagin, appErr := args.Parse("DigitalContents")
	if appErr != nil {
		return nil, appErr
	}

	// check if this is initial query, then no order by is passed, we need to provide it here
	if graphqlPagin.OrderBy == "" {
		orderDirection := args.orderDirection().String()
		graphqlPagin.OrderBy = util.AnyArray[string]{model.DigitalContentTableName + ".ContentType", model.DigitalContentTableName + ".ContentFile"}.
			Map(func(_ int, item string) string { return item + " " + orderDirection }).
			Join(", ")
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	totalCount, contents, appErr := embedCtx.App.Srv().ProductService().DigitalContentsbyOptions(&model.DigitalContentFilterOption{
		PaginationValues: *graphqlPagin,
		CountTotal:       true,
	})

	keyFunc := func(d *model.DigitalContent) []any {
		return []any{
			model.DigitalContentTableName + ".ContentType", d.ContentType,
			model.DigitalContentTableName + ".ContentFile", d.ContentFile,
		}
	}
	res := constructCountableConnection(contents, totalCount, args, keyFunc, systemDigitalContentToGraphqlDigitalContent)
	return (*DigitalContentCountableConnection)(unsafe.Pointer(res)), nil
}

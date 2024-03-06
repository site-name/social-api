package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"net/http"
	"unsafe"

	"github.com/mattermost/squirrel"
	"github.com/samber/lo"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/web"
)

// NOTE: Please refer to ./graphql/schemas/product_media.graphqls for details on directives used
func (r *Resolver) ProductTypeCreate(ctx context.Context, args struct{ Input ProductTypeInput }) (*ProductTypeCreate, error) {
	appErr := args.Input.validate("ProductTypeCreate")
	if appErr != nil {
		return nil, appErr
	}

	// construct new product type:
	var productType model.ProductType
	args.Input.patch(&productType)

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	if args.Input.TaxCode != nil {
		pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()
		_, appErr = pluginMng.AssignTaxCodeToObjectMeta(productType, *args.Input.TaxCode)
		if appErr != nil {
			return nil, appErr
		}
	}

	// NOTE: product attributes go first [0], then variant attributes [1]
	var attributes = [2]model.Attributes{}

	for idx, ids := range [2][]UUID{
		args.Input.ProductAttributes,
		args.Input.VariantAttributes,
	} {
		if len(ids) > 0 {
			var appErr *model_helper.AppError
			attributes[idx], appErr = embedCtx.App.Srv().AttributeService().AttributesByOption(&model.AttributeFilterOption{
				Conditions: squirrel.Eq{model.AttributeTableName + ".Id": ids},
			})
			if appErr != nil {
				return nil, appErr
			}

			// check if there are some attribute(s) that is not product type
			if attributes[idx] != nil && lo.SomeBy(attributes[idx], func(item *model.Attribute) bool { return item.Type != model.PRODUCT_TYPE }) {
				return nil, model_helper.NewAppError("ProductTypeCreate", model_helper.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "attributes"}, "please provide attributes with types are product type", http.StatusBadRequest)
			}
		}
	}

	// save product type
	savedProductType, appErr := embedCtx.App.Srv().ProductService().UpsertProductType(nil, &productType)
	if appErr != nil {
		return nil, appErr
	}

	// add many to many attributes
	appErr = embedCtx.App.Srv().ProductService().ToggleProductTypeAttributeRelations(nil, savedProductType.Id, attributes[1], attributes[0], false)
	if appErr != nil {
		return nil, appErr
	}

	return &ProductTypeCreate{
		ProductType: SystemProductTypeToGraphqlProductType(savedProductType),
	}, nil
}

// NOTE: Please refer to ./graphql/schemas/product_media.graphqls for details on directives used
func (r *Resolver) ProductTypeDelete(ctx context.Context, args struct{ Id UUID }) (*ProductTypeDelete, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	// begin tx
	tx := embedCtx.App.Srv().Store.GetMaster().Begin()
	if tx.Error != nil {
		return nil, model_helper.NewAppError("ProductTypeDelete", model.ErrorCreatingTransactionErrorID, nil, tx.Error.Error(), http.StatusInternalServerError)
	}
	defer embedCtx.App.Srv().Store.FinalizeTransaction(tx)

	_, appErr := embedCtx.App.Srv().ProductService().DeleteProductTypes(tx, []string{args.Id.String()})
	if appErr != nil {
		return nil, appErr
	}

	// commit
	if err := tx.Commit().Error; err != nil {
		return nil, model_helper.NewAppError("ProductTypeDelete", model.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return &ProductTypeDelete{
		ProductType: &ProductType{
			ID: args.Id.String(),
		},
	}, nil
}

// NOTE: Please refer to ./graphql/schemas/product_media.graphqls for details on directives used
func (r *Resolver) ProductTypeBulkDelete(ctx context.Context, args struct{ Ids []UUID }) (*ProductTypeBulkDelete, error) {
	if len(args.Ids) == 0 {
		return &ProductTypeBulkDelete{
			Count: 0,
		}, nil
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	// begin tx
	tx := embedCtx.App.Srv().Store.GetMaster().Begin()
	if tx.Error != nil {
		return nil, model_helper.NewAppError("ProductTypeBulkDelete", model.ErrorCreatingTransactionErrorID, nil, tx.Error.Error(), http.StatusInternalServerError)
	}
	defer embedCtx.App.Srv().Store.FinalizeTransaction(tx)

	ids := *(*[]string)(unsafe.Pointer(&args.Ids))
	delCount, appErr := embedCtx.App.Srv().ProductService().DeleteProductTypes(tx, ids)
	if appErr != nil {
		return nil, appErr
	}

	// commit
	if err := tx.Commit().Error; err != nil {
		return nil, model_helper.NewAppError("ProductTypeBulkDelete", model.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return &ProductTypeBulkDelete{
		Count: int32(delCount),
	}, nil
}

// NOTE: Please refer to ./graphql/schemas/product_media.graphqls for details on directives used
func (r *Resolver) ProductTypeUpdate(ctx context.Context, args struct {
	Id    UUID
	Input ProductTypeInput
}) (*ProductTypeUpdate, error) {
	appErr := args.Input.validate("ProductTypeUpdate")
	if appErr != nil {
		return nil, appErr
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	productType, appErr := embedCtx.App.Srv().ProductService().ProductTypeByOption(&model.ProductTypeFilterOption{
		Conditions: squirrel.Eq{model.ProductTypeTableName + ".Id": args.Id},
	})
	if appErr != nil {
		return nil, appErr
	}

	args.Input.patch(productType)
	updatedProductType, appErr := embedCtx.App.Srv().ProductService().UpsertProductType(nil, productType)
	if appErr != nil {
		return nil, appErr
	}

	// NOTE: original code does rename related product variants
	// but this code doesn't do that
	// we should consider whether to add that

	// NOTE: product attributes go first [0], then variant attributes [1]
	var attributes = [2]model.Attributes{}

	for idx, ids := range [2][]UUID{
		args.Input.ProductAttributes,
		args.Input.VariantAttributes,
	} {
		if len(ids) > 0 {
			var appErr *model_helper.AppError
			attributes[idx], appErr = embedCtx.App.Srv().AttributeService().AttributesByOption(&model.AttributeFilterOption{
				Conditions: squirrel.Eq{model.AttributeTableName + ".Id": ids},
			})
			if appErr != nil {
				return nil, appErr
			}

			// check if there are some attribute(s) that is not product type
			if attributes[idx] != nil && lo.SomeBy(attributes[idx], func(item *model.Attribute) bool { return item.Type != model.PRODUCT_TYPE }) {
				return nil, model_helper.NewAppError("ProductTypeUpdate", model_helper.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "attributes"}, "please provide attributes with types are product type", http.StatusBadRequest)
			}
		}
	}

	// add many to many attributes
	appErr = embedCtx.App.Srv().ProductService().ToggleProductTypeAttributeRelations(nil, updatedProductType.Id, attributes[1], attributes[0], false)
	if appErr != nil {
		return nil, appErr
	}

	return &ProductTypeUpdate{
		ProductType: SystemProductTypeToGraphqlProductType(updatedProductType),
	}, nil
}

// NOTE: Please refer to ./graphql/schemas/product_media.graphqls for details on directives used
func (r *Resolver) ProductTypeReorderAttributes(ctx context.Context, args struct {
	Moves         []*ReorderInput
	ProductTypeID UUID
	Type          ProductAttributeType
}) (*ProductTypeReorderAttributes, error) {
	// validate params
	// if !args.Type.IsValid() {
	// 	return nil, model_helper.NewAppError("ProductTypeReorderAttributes", model_helper.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Type"}, "please provide valid product attribute type", http.StatusBadRequest)
	// }

	// embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	// embedCtx.App.Srv().ProductService().ProductTypeByOption(&model.ProductTypeFilterOption{
	// 	Conditions: squirrel.Expr(model.ProductTypeTableName + ".Id = ?", args.ProductTypeID),
	// })

	panic("not implemented")
}

// NOTE: Please refer to ./graphql/schemas/product_media.graphqls for details on directives used
func (r *Resolver) ProductType(ctx context.Context, args struct{ Id UUID }) (*ProductType, error) {
	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	productType, appErr := embedCtx.App.Srv().ProductService().ProductTypeByOption(&model.ProductTypeFilterOption{
		Conditions: squirrel.Eq{model.ProductTypeTableName + ".Id": args.Id},
	})
	if appErr != nil {
		return nil, appErr
	}

	return SystemProductTypeToGraphqlProductType(productType), nil
}

// NOTE: Please refer to ./graphql/schemas/product_media.graphqls for details on directives used
func (r *Resolver) ProductTypes(ctx context.Context, args struct {
	Filter *ProductTypeFilterInput
	SortBy *ProductTypeSortingInput
	GraphqlParams
}) (*ProductTypeCountableConnection, error) {
	paginValues, appErr := args.GraphqlParams.Parse("ProductTypes")
	if appErr != nil {
		return nil, appErr
	}

	filterOpts := args.Filter.parse("ProductTypes")
	filterOpts.GraphqlPaginationValues = *paginValues
	filterOpts.CountTotal = true

	if filterOpts.GraphqlPaginationValues.OrderBy == "" {
		orderFields := productTypeSortFieldsMap[ProductTypeSortFieldName].fields

		if args.SortBy != nil && args.SortBy.Field.IsValid() {
			orderFields = productTypeSortFieldsMap[args.SortBy.Field].fields
		}

		ordering := args.GraphqlParams.orderDirection().String()
		filterOpts.GraphqlPaginationValues.OrderBy = orderFields.
			Map(func(_ int, item string) string { return item + " " + ordering }).
			Join(",")
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	totalCount, productTypes, appErr := embedCtx.App.Srv().ProductService().ProductTypesByOptions(filterOpts)
	if appErr != nil {
		return nil, appErr
	}

	keyFunc := productTypeSortFieldsMap[ProductTypeSortFieldName].keyFunc
	if args.SortBy != nil {
		keyFunc = productTypeSortFieldsMap[args.SortBy.Field].keyFunc
	}
	res := constructCountableConnection(productTypes, totalCount, args.GraphqlParams, keyFunc, SystemProductTypeToGraphqlProductType)
	return (*ProductTypeCountableConnection)(unsafe.Pointer(res)), nil
}

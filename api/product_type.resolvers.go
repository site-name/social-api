package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/samber/lo"
	"github.com/sitename/sitename/model"
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
			var appErr *model.AppError
			attributes[idx], appErr = embedCtx.App.Srv().AttributeService().AttributesByOption(&model.AttributeFilterOption{
				Conditions: squirrel.Eq{model.AttributeTableName + ".Id": ids},
			})
			if appErr != nil {
				return nil, appErr
			}

			// check if there are some attribute(s) that is not product type
			if attributes[idx] != nil && lo.SomeBy(attributes[idx], func(item *model.Attribute) bool { return item.Type != model.PRODUCT_TYPE }) {
				return nil, model.NewAppError("ProductTypeCreate", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "attributes"}, "please provide attributes with types are product type", http.StatusBadRequest)
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
	// embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	// embedCtx.App.Srv().ProductService().DeleteProductTypes(nil, []string{args.Id.String()})
	panic("not implemented")
}

// NOTE: Please refer to ./graphql/schemas/product_media.graphqls for details on directives used
func (r *Resolver) ProductTypeBulkDelete(ctx context.Context, args struct{ Ids []string }) (*ProductTypeBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

// NOTE: Please refer to ./graphql/schemas/product_media.graphqls for details on directives used
func (r *Resolver) ProductTypeUpdate(ctx context.Context, args struct {
	Id    UUID
	Input ProductTypeInput
}) (*ProductTypeUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

// NOTE: Please refer to ./graphql/schemas/product_media.graphqls for details on directives used
func (r *Resolver) ProductTypeReorderAttributes(ctx context.Context, args struct {
	Moves         []*ReorderInput
	ProductTypeID UUID
	Type          ProductAttributeType
}) (*ProductTypeReorderAttributes, error) {
	panic(fmt.Errorf("not implemented"))
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

func (r *Resolver) ProductTypes(ctx context.Context, args struct {
	Filter *ProductTypeFilterInput
	SortBy *ProductTypeSortingInput
	GraphqlParams
}) (*ProductTypeCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

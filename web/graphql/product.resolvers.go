package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/web/graphql/gqlmodel"
)

func (r *mutationResolver) ProductAttributeAssign(ctx context.Context, operations []*gqlmodel.ProductAttributeAssignInput, productTypeID string) (*gqlmodel.ProductAttributeAssign, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductAttributeUnassign(ctx context.Context, attributeIds []*string, productTypeID string) (*gqlmodel.ProductAttributeUnassign, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductCreate(ctx context.Context, input gqlmodel.ProductCreateInput) (*gqlmodel.ProductCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductDelete(ctx context.Context, id string) (*gqlmodel.ProductDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductBulkDelete(ctx context.Context, ids []*string) (*gqlmodel.ProductBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductUpdate(ctx context.Context, id string, input gqlmodel.ProductInput) (*gqlmodel.ProductUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductTranslate(ctx context.Context, id string, input gqlmodel.TranslationInput, languageCode gqlmodel.LanguageCodeEnum) (*gqlmodel.ProductTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductChannelListingUpdate(ctx context.Context, id string, input gqlmodel.ProductChannelListingUpdateInput) (*gqlmodel.ProductChannelListingUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductReorderAttributeValues(ctx context.Context, attributeID string, moves []*gqlmodel.ReorderInput, productID string) (*gqlmodel.ProductReorderAttributeValues, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Product(ctx context.Context, id *string, slug *string, channel *string) (*gqlmodel.Product, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Products(ctx context.Context, filter *gqlmodel.ProductFilterInput, sortBy *gqlmodel.ProductOrder, channel *string, before *string, after *string, first *int, last *int) (*gqlmodel.ProductCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

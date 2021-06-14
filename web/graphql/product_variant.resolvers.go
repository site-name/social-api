package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/web/graphql/gqlmodel"
)

func (r *mutationResolver) ProductVariantReorder(ctx context.Context, moves []*gqlmodel.ReorderInput, productID string) (*gqlmodel.ProductVariantReorder, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductVariantCreate(ctx context.Context, input gqlmodel.ProductVariantCreateInput) (*gqlmodel.ProductVariantCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductVariantDelete(ctx context.Context, id string) (*gqlmodel.ProductVariantDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductVariantBulkCreate(ctx context.Context, product string, variants []*gqlmodel.ProductVariantBulkCreateInput) (*gqlmodel.ProductVariantBulkCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductVariantBulkDelete(ctx context.Context, ids []*string) (*gqlmodel.ProductVariantBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductVariantStocksCreate(ctx context.Context, stocks []*gqlmodel.StockInput, variantID string) (*gqlmodel.ProductVariantStocksCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductVariantStocksDelete(ctx context.Context, variantID string, warehouseIds []string) (*gqlmodel.ProductVariantStocksDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductVariantStocksUpdate(ctx context.Context, stocks []*gqlmodel.StockInput, variantID string) (*gqlmodel.ProductVariantStocksUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductVariantUpdate(ctx context.Context, id string, input gqlmodel.ProductVariantInput) (*gqlmodel.ProductVariantUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductVariantSetDefault(ctx context.Context, productID string, variantID string) (*gqlmodel.ProductVariantSetDefault, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductVariantTranslate(ctx context.Context, id string, input gqlmodel.NameTranslationInput, languageCode gqlmodel.LanguageCodeEnum) (*gqlmodel.ProductVariantTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductVariantChannelListingUpdate(ctx context.Context, id string, input []*gqlmodel.ProductVariantChannelListingAddInput) (*gqlmodel.ProductVariantChannelListingUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductVariantReorderAttributeValues(ctx context.Context, attributeID string, moves []*gqlmodel.ReorderInput, variantID string) (*gqlmodel.ProductVariantReorderAttributeValues, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) ProductVariant(ctx context.Context, id *string, sku *string, channel *string) (*gqlmodel.ProductVariant, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) ProductVariants(ctx context.Context, ids []*string, channel *string, filter *gqlmodel.ProductVariantFilterInput, before *string, after *string, first *int, last *int) (*gqlmodel.ProductVariantCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

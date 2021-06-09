package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
)

func (r *mutationResolver) ProductVariantReorder(ctx context.Context, moves []*ReorderInput, productID string) (*ProductVariantReorder, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductVariantCreate(ctx context.Context, input ProductVariantCreateInput) (*ProductVariantCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductVariantDelete(ctx context.Context, id string) (*ProductVariantDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductVariantBulkCreate(ctx context.Context, product string, variants []*ProductVariantBulkCreateInput) (*ProductVariantBulkCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductVariantBulkDelete(ctx context.Context, ids []*string) (*ProductVariantBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductVariantStocksCreate(ctx context.Context, stocks []StockInput, variantID string) (*ProductVariantStocksCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductVariantStocksDelete(ctx context.Context, variantID string, warehouseIds []string) (*ProductVariantStocksDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductVariantStocksUpdate(ctx context.Context, stocks []StockInput, variantID string) (*ProductVariantStocksUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductVariantUpdate(ctx context.Context, id string, input ProductVariantInput) (*ProductVariantUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductVariantSetDefault(ctx context.Context, productID string, variantID string) (*ProductVariantSetDefault, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductVariantTranslate(ctx context.Context, id string, input NameTranslationInput, languageCode LanguageCodeEnum) (*ProductVariantTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductVariantChannelListingUpdate(ctx context.Context, id string, input []ProductVariantChannelListingAddInput) (*ProductVariantChannelListingUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProductVariantReorderAttributeValues(ctx context.Context, attributeID string, moves []*ReorderInput, variantID string) (*ProductVariantReorderAttributeValues, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) ProductVariant(ctx context.Context, id *string, sku *string, channel *string) (*ProductVariant, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) ProductVariants(ctx context.Context, ids []*string, channel *string, filter *ProductVariantFilterInput, before *string, after *string, first *int, last *int) (*ProductVariantCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

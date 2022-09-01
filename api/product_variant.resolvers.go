package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/api/gqlmodel"
)

func (r *Resolver) ProductVariantReorder(ctx context.Context, args struct {
	moves     []*gqlmodel.ReorderInput
	productID string
}) (*gqlmodel.ProductVariantReorder, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductVariantCreate(ctx context.Context, args struct {
	input gqlmodel.ProductVariantCreateInput
}) (*gqlmodel.ProductVariantCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductVariantDelete(ctx context.Context, args struct{ id string }) (*gqlmodel.ProductVariantDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductVariantBulkCreate(ctx context.Context, args struct {
	product  string
	variants []*gqlmodel.ProductVariantBulkCreateInput
}) (*gqlmodel.ProductVariantBulkCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductVariantBulkDelete(ctx context.Context, args struct{ ids []*string }) (*gqlmodel.ProductVariantBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductVariantStocksCreate(ctx context.Context, args struct {
	stocks    []gqlmodel.StockInput
	variantID string
}) (*gqlmodel.ProductVariantStocksCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductVariantStocksDelete(ctx context.Context, args struct {
	variantID    string
	warehouseIds []string
}) (*gqlmodel.ProductVariantStocksDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductVariantStocksUpdate(ctx context.Context, args struct {
	stocks    []gqlmodel.StockInput
	variantID string
}) (*gqlmodel.ProductVariantStocksUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductVariantUpdate(ctx context.Context, args struct {
	id    string
	input gqlmodel.ProductVariantInput
}) (*gqlmodel.ProductVariantUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductVariantSetDefault(ctx context.Context, args struct {
	productID string
	variantID string
}) (*gqlmodel.ProductVariantSetDefault, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductVariantTranslate(ctx context.Context, args struct {
	id           string
	input        gqlmodel.NameTranslationInput
	languageCode gqlmodel.LanguageCodeEnum
}) (*gqlmodel.ProductVariantTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductVariantChannelListingUpdate(ctx context.Context, args struct {
	id    string
	input []gqlmodel.ProductVariantChannelListingAddInput
}) (*gqlmodel.ProductVariantChannelListingUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductVariantReorderAttributeValues(ctx context.Context, args struct {
	attributeID string
	moves       []*gqlmodel.ReorderInput
	variantID   string
}) (*gqlmodel.ProductVariantReorderAttributeValues, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductVariant(ctx context.Context, args struct {
	id      *string
	sku     *string
	channel *string
}) (*gqlmodel.ProductVariant, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductVariants(ctx context.Context, args struct {
	ids     []*string
	channel *string
	filter  *gqlmodel.ProductVariantFilterInput
	before  *string
	after   *string
	first   *int
	last    *int
}) (*gqlmodel.ProductVariantCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

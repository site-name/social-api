package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/api/gqlmodel"
)

func (r *Resolver) ProductVariantReorder(ctx context.Context, args struct {
	Moves     []*gqlmodel.ReorderInput
	ProductID string
}) (*gqlmodel.ProductVariantReorder, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductVariantCreate(ctx context.Context, args struct {
	Input gqlmodel.ProductVariantCreateInput
}) (*gqlmodel.ProductVariantCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductVariantDelete(ctx context.Context, args struct{ Id string }) (*gqlmodel.ProductVariantDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductVariantBulkCreate(ctx context.Context, args struct {
	Product  string
	Variants []*gqlmodel.ProductVariantBulkCreateInput
}) (*gqlmodel.ProductVariantBulkCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductVariantBulkDelete(ctx context.Context, args struct{ Ids []*string }) (*gqlmodel.ProductVariantBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductVariantStocksCreate(ctx context.Context, args struct {
	Stocks    []gqlmodel.StockInput
	VariantID string
}) (*gqlmodel.ProductVariantStocksCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductVariantStocksDelete(ctx context.Context, args struct {
	VariantID    string
	WarehouseIds []string
}) (*gqlmodel.ProductVariantStocksDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductVariantStocksUpdate(ctx context.Context, args struct {
	Stocks    []gqlmodel.StockInput
	VariantID string
}) (*gqlmodel.ProductVariantStocksUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductVariantUpdate(ctx context.Context, args struct {
	Id    string
	Input gqlmodel.ProductVariantInput
}) (*gqlmodel.ProductVariantUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductVariantSetDefault(ctx context.Context, args struct {
	ProductID string
	VariantID string
}) (*gqlmodel.ProductVariantSetDefault, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductVariantTranslate(ctx context.Context, args struct {
	Id           string
	Input        gqlmodel.NameTranslationInput
	LanguageCode gqlmodel.LanguageCodeEnum
}) (*gqlmodel.ProductVariantTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductVariantChannelListingUpdate(ctx context.Context, args struct {
	Id    string
	Input []gqlmodel.ProductVariantChannelListingAddInput
}) (*gqlmodel.ProductVariantChannelListingUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductVariantReorderAttributeValues(ctx context.Context, args struct {
	AttributeID string
	Moves       []*gqlmodel.ReorderInput
	VariantID   string
}) (*gqlmodel.ProductVariantReorderAttributeValues, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductVariant(ctx context.Context, args struct {
	Id      *string
	Sku     *string
	Channel *string
}) (*gqlmodel.ProductVariant, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductVariants(ctx context.Context, args struct {
	Ids     []*string
	Channel *string
	Filter  *gqlmodel.ProductVariantFilterInput
	Before  *string
	After   *string
	First   *int
	Last    *int
}) (*gqlmodel.ProductVariantCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

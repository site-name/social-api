package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
)

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
	Before  *string
	After   *string
	First   *int
	Last    *int
}) (*ProductVariantCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

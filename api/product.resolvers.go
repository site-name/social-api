package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/api/gqlmodel"
)

func (r *Resolver) ProductAttributeAssign(ctx context.Context, args struct {
	operations    []*gqlmodel.ProductAttributeAssignInput
	productTypeID string
}) (*gqlmodel.ProductAttributeAssign, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductAttributeUnassign(ctx context.Context, args struct {
	attributeIds  []*string
	productTypeID string
}) (*gqlmodel.ProductAttributeUnassign, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductCreate(ctx context.Context, args struct{ input gqlmodel.ProductCreateInput }) (*gqlmodel.ProductCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductDelete(ctx context.Context, args struct{ id string }) (*gqlmodel.ProductDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductBulkDelete(ctx context.Context, args struct{ ids []*string }) (*gqlmodel.ProductBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductUpdate(ctx context.Context, args struct {
	id    string
	input gqlmodel.ProductInput
}) (*gqlmodel.ProductUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductTranslate(ctx context.Context, id string, input gqlmodel.TranslationInput, languageCode gqlmodel.LanguageCodeEnum) (*gqlmodel.ProductTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductChannelListingUpdate(ctx context.Context, id string, input gqlmodel.ProductChannelListingUpdateInput) (*gqlmodel.ProductChannelListingUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductReorderAttributeValues(ctx context.Context, attributeID string, moves []*gqlmodel.ReorderInput, productID string) (*gqlmodel.ProductReorderAttributeValues, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Product(ctx context.Context, id *string, slug *string, channel *string) (*gqlmodel.Product, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Products(ctx context.Context, filter *gqlmodel.ProductFilterInput, sortBy *gqlmodel.ProductOrder, channel *string, before *string, after *string, first *int, last *int) (*gqlmodel.ProductCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

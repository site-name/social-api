package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/api/gqlmodel"
)

func (r *Resolver) ProductAttributeAssign(ctx context.Context, args struct {
	Operations    []*gqlmodel.ProductAttributeAssignInput
	ProductTypeID string
}) (*gqlmodel.ProductAttributeAssign, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductAttributeUnassign(ctx context.Context, args struct {
	AttributeIds  []*string
	ProductTypeID string
}) (*gqlmodel.ProductAttributeUnassign, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductCreate(ctx context.Context, args struct{ Input gqlmodel.ProductCreateInput }) (*gqlmodel.ProductCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductDelete(ctx context.Context, args struct{ Id string }) (*gqlmodel.ProductDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductBulkDelete(ctx context.Context, args struct{ Ids []*string }) (*gqlmodel.ProductBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductUpdate(ctx context.Context, args struct {
	Id    string
	Input gqlmodel.ProductInput
}) (*gqlmodel.ProductUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductTranslate(ctx context.Context, args struct {
	Id           string
	Input        gqlmodel.TranslationInput
	LanguageCode gqlmodel.LanguageCodeEnum
}) (*gqlmodel.ProductTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductChannelListingUpdate(ctx context.Context, args struct {
	Id    string
	Input gqlmodel.ProductChannelListingUpdateInput
}) (*gqlmodel.ProductChannelListingUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductReorderAttributeValues(ctx context.Context, args struct {
	AttributeID string
	Moves       []*gqlmodel.ReorderInput
	ProductID   string
}) (*gqlmodel.ProductReorderAttributeValues, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Product(ctx context.Context, args struct {
	Id      *string
	Slug    *string
	Channel *string
}) (*gqlmodel.Product, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Products(ctx context.Context, args struct {
	Filter  *gqlmodel.ProductFilterInput
	SortBy  *gqlmodel.ProductOrder
	Channel *string
	Before  *string
	After   *string
	First   *int
	Last    *int
}) (*gqlmodel.ProductCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
)

func (r *Resolver) ProductAttributeAssign(ctx context.Context, args struct {
	Operations    []*ProductAttributeAssignInput
	ProductTypeID string
}) (*ProductAttributeAssign, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductAttributeUnassign(ctx context.Context, args struct {
	AttributeIds  []*string
	ProductTypeID string
}) (*ProductAttributeUnassign, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductCreate(ctx context.Context, args struct{ Input ProductCreateInput }) (*ProductCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductDelete(ctx context.Context, args struct{ Id string }) (*ProductDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductBulkDelete(ctx context.Context, args struct{ Ids []*string }) (*ProductBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductUpdate(ctx context.Context, args struct {
	Id    string
	Input ProductInput
}) (*ProductUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductTranslate(ctx context.Context, args struct {
	Id           string
	Input        TranslationInput
	LanguageCode LanguageCodeEnum
}) (*ProductTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductChannelListingUpdate(ctx context.Context, args struct {
	Id    string
	Input ProductChannelListingUpdateInput
}) (*ProductChannelListingUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ProductReorderAttributeValues(ctx context.Context, args struct {
	AttributeID string
	Moves       []*ReorderInput
	ProductID   string
}) (*ProductReorderAttributeValues, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Product(ctx context.Context, args struct {
	Id      *string
	Slug    *string
	Channel *string
}) (*Product, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Products(ctx context.Context, args struct {
	Filter  *ProductFilterInput
	SortBy  *ProductOrder
	Channel *string
	Before  *string
	After   *string
	First   *int
	Last    *int
}) (*ProductCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

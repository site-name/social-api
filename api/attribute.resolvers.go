package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/api/gqlmodel"
)

func (r *Resolver) AttributeCreate(ctx context.Context, args struct{ Input gqlmodel.AttributeCreateInput }) (*gqlmodel.AttributeCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AttributeDelete(ctx context.Context, args struct{ Id string }) (*gqlmodel.AttributeDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AttributeUpdate(ctx context.Context, args struct {
	Id    string
	Input gqlmodel.AttributeUpdateInput
}) (*gqlmodel.AttributeUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AttributeTranslate(ctx context.Context, args struct {
	Id           string
	Input        gqlmodel.NameTranslationInput
	LanguageCode gqlmodel.LanguageCodeEnum
}) (*gqlmodel.AttributeTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AttributeBulkDelete(ctx context.Context, args struct{ Ids []*string }) (*gqlmodel.AttributeBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AttributeValueBulkDelete(ctx context.Context, args struct{ Ids []*string }) (*gqlmodel.AttributeValueBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AttributeValueCreate(ctx context.Context, args struct {
	AttributeID string
	Input       gqlmodel.AttributeValueCreateInput
}) (*gqlmodel.AttributeValueCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AttributeValueDelete(ctx context.Context, args struct{ Id string }) (*gqlmodel.AttributeValueDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AttributeValueUpdate(ctx context.Context, args struct {
	Id    string
	Input gqlmodel.AttributeValueUpdateInput
}) (*gqlmodel.AttributeValueUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AttributeValueTranslate(ctx context.Context, args struct {
	Id           string
	Input        gqlmodel.AttributeValueTranslationInput
	LanguageCode gqlmodel.LanguageCodeEnum
}) (*gqlmodel.AttributeValueTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AttributeReorderValues(ctx context.Context, args struct {
	AttributeID string
	Moves       []*gqlmodel.ReorderInput
}) (*gqlmodel.AttributeReorderValues, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Attributes(ctx context.Context, args struct {
	Filter     *gqlmodel.AttributeFilterInput
	SortBy     *gqlmodel.AttributeSortingInput
	ChanelSlug *string
	Before     *string
	After      *string
	First      *int
	Last       *int
}) (*gqlmodel.AttributeCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Attribute(ctx context.Context, args struct {
	Id   *string
	Slug *string
}) (*gqlmodel.Attribute, error) {
	panic(fmt.Errorf("not implemented"))
}

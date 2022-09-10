package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
)

func (r *Resolver) AttributeCreate(ctx context.Context, args struct{ Input AttributeCreateInput }) (*AttributeCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AttributeDelete(ctx context.Context, args struct{ Id string }) (*AttributeDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AttributeUpdate(ctx context.Context, args struct {
	Id    string
	Input AttributeUpdateInput
}) (*AttributeUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AttributeTranslate(ctx context.Context, args struct {
	Id           string
	Input        NameTranslationInput
	LanguageCode LanguageCodeEnum
}) (*AttributeTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AttributeBulkDelete(ctx context.Context, args struct{ Ids []*string }) (*AttributeBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AttributeValueBulkDelete(ctx context.Context, args struct{ Ids []*string }) (*AttributeValueBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AttributeValueCreate(ctx context.Context, args struct {
	AttributeID string
	Input       AttributeValueCreateInput
}) (*AttributeValueCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AttributeValueDelete(ctx context.Context, args struct{ Id string }) (*AttributeValueDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AttributeValueUpdate(ctx context.Context, args struct {
	Id    string
	Input AttributeValueUpdateInput
}) (*AttributeValueUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AttributeValueTranslate(ctx context.Context, args struct {
	Id           string
	Input        AttributeValueTranslationInput
	LanguageCode LanguageCodeEnum
}) (*AttributeValueTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) AttributeReorderValues(ctx context.Context, args struct {
	AttributeID string
	Moves       []*ReorderInput
}) (*AttributeReorderValues, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Attributes(ctx context.Context, args struct {
	Filter     *AttributeFilterInput
	SortBy     *AttributeSortingInput
	ChanelSlug *string
	Before     *string
	After      *string
	First      *int
	Last       *int
}) (*AttributeCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Attribute(ctx context.Context, args struct {
	Id   *string
	Slug *string
}) (*Attribute, error) {
	panic(fmt.Errorf("not implemented"))
}

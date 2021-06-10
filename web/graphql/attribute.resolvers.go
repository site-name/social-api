package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
)

func (r *mutationResolver) AttributeCreate(ctx context.Context, input AttributeCreateInput) (*AttributeCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AttributeDelete(ctx context.Context, id string) (*AttributeDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AttributeUpdate(ctx context.Context, id string, input AttributeUpdateInput) (*AttributeUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AttributeTranslate(ctx context.Context, id string, input NameTranslationInput, languageCode LanguageCodeEnum) (*AttributeTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AttributeBulkDelete(ctx context.Context, ids []*string) (*AttributeBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AttributeValueBulkDelete(ctx context.Context, ids []*string) (*AttributeValueBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AttributeValueCreate(ctx context.Context, attribute string, input AttributeValueCreateInput) (*AttributeValueCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AttributeValueDelete(ctx context.Context, id string) (*AttributeValueDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AttributeValueUpdate(ctx context.Context, id string, input AttributeValueCreateInput) (*AttributeValueUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AttributeValueTranslate(ctx context.Context, id string, input AttributeValueTranslationInput, languageCode LanguageCodeEnum) (*AttributeValueTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AttributeReorderValues(ctx context.Context, attributeID string, moves []*ReorderInput) (*AttributeReorderValues, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Attributes(ctx context.Context, filter *AttributeFilterInput, sortBy *AttributeSortingInput, before *string, after *string, first *int, last *int) (*AttributeCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Attribute(ctx context.Context, id *string, slug *string) (*Attribute, error) {
	panic(fmt.Errorf("not implemented"))
}

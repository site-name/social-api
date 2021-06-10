package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
)

func (r *mutationResolver) PageCreate(ctx context.Context, input PageCreateInput) (*PageCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PageDelete(ctx context.Context, id string) (*PageDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PageBulkDelete(ctx context.Context, ids []*string) (*PageBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PageBulkPublish(ctx context.Context, ids []*string, isPublished bool) (*PageBulkPublish, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PageUpdate(ctx context.Context, id string, input PageInput) (*PageUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PageTranslate(ctx context.Context, id string, input PageTranslationInput, languageCode LanguageCodeEnum) (*PageTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PageAttributeAssign(ctx context.Context, attributeIds []string, pageTypeID string) (*PageAttributeAssign, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PageAttributeUnassign(ctx context.Context, attributeIds []string, pageTypeID string) (*PageAttributeUnassign, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PageReorderAttributeValues(ctx context.Context, attributeID string, moves []*ReorderInput, pageID string) (*PageReorderAttributeValues, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Page(ctx context.Context, id *string, slug *string) (*Page, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Pages(ctx context.Context, sortBy *PageSortingInput, filter *PageFilterInput, before *string, after *string, first *int, last *int) (*PageCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

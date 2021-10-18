package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/graphql/gqlmodel"
)

func (r *mutationResolver) PageCreate(ctx context.Context, input gqlmodel.PageCreateInput) (*gqlmodel.PageCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PageDelete(ctx context.Context, id string) (*gqlmodel.PageDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PageBulkDelete(ctx context.Context, ids []*string) (*gqlmodel.PageBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PageBulkPublish(ctx context.Context, ids []*string, isPublished bool) (*gqlmodel.PageBulkPublish, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PageUpdate(ctx context.Context, id string, input gqlmodel.PageInput) (*gqlmodel.PageUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PageTranslate(ctx context.Context, id string, input gqlmodel.PageTranslationInput, languageCode gqlmodel.LanguageCodeEnum) (*gqlmodel.PageTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PageAttributeAssign(ctx context.Context, attributeIds []string, pageTypeID string) (*gqlmodel.PageAttributeAssign, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PageAttributeUnassign(ctx context.Context, attributeIds []string, pageTypeID string) (*gqlmodel.PageAttributeUnassign, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PageReorderAttributeValues(ctx context.Context, attributeID string, moves []*gqlmodel.ReorderInput, pageID string) (*gqlmodel.PageReorderAttributeValues, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Page(ctx context.Context, id *string, slug *string) (*gqlmodel.Page, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Pages(ctx context.Context, sortBy *gqlmodel.PageSortingInput, filter *gqlmodel.PageFilterInput, before *string, after *string, first *int, last *int) (*gqlmodel.PageCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

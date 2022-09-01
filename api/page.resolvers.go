package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/api/gqlmodel"
)

func (r *Resolver) PageCreate(ctx context.Context, args struct{ input gqlmodel.PageCreateInput }) (*gqlmodel.PageCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) PageDelete(ctx context.Context, args struct{ id string }) (*gqlmodel.PageDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) PageBulkDelete(ctx context.Context, args struct{ ids []*string }) (*gqlmodel.PageBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) PageBulkPublish(ctx context.Context, args struct {
	ids         []*string
	isPublished bool
}) (*gqlmodel.PageBulkPublish, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) PageUpdate(ctx context.Context, args struct {
	id    string
	input gqlmodel.PageInput
}) (*gqlmodel.PageUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) PageTranslate(ctx context.Context, args struct {
	id           string
	input        gqlmodel.PageTranslationInput
	languageCode gqlmodel.LanguageCodeEnum
}) (*gqlmodel.PageTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) PageAttributeAssign(ctx context.Context, args struct {
	attributeIds []string
	pageTypeID   string
}) (*gqlmodel.PageAttributeAssign, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) PageAttributeUnassign(ctx context.Context, args struct {
	attributeIds []string
	pageTypeID   string
}) (*gqlmodel.PageAttributeUnassign, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) PageReorderAttributeValues(ctx context.Context, attributeID string, moves []*gqlmodel.ReorderInput, pageID string) (*gqlmodel.PageReorderAttributeValues, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Page(ctx context.Context, id *string, slug *string) (*gqlmodel.Page, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Pages(ctx context.Context, sortBy *gqlmodel.PageSortingInput, filter *gqlmodel.PageFilterInput, before *string, after *string, first *int, last *int) (*gqlmodel.PageCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

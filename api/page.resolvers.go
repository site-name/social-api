package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/api/gqlmodel"
)

func (r *Resolver) PageCreate(ctx context.Context, args struct{ Input gqlmodel.PageCreateInput }) (*gqlmodel.PageCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) PageDelete(ctx context.Context, args struct{ Id string }) (*gqlmodel.PageDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) PageBulkDelete(ctx context.Context, args struct{ Ids []*string }) (*gqlmodel.PageBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) PageBulkPublish(ctx context.Context, args struct {
	Ids         []*string
	IsPublished bool
}) (*gqlmodel.PageBulkPublish, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) PageUpdate(ctx context.Context, args struct {
	Id    string
	Input gqlmodel.PageInput
}) (*gqlmodel.PageUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) PageTranslate(ctx context.Context, args struct {
	Id           string
	Input        gqlmodel.PageTranslationInput
	LanguageCode gqlmodel.LanguageCodeEnum
}) (*gqlmodel.PageTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) PageAttributeAssign(ctx context.Context, args struct {
	AttributeIds []string
	PageTypeID   string
}) (*gqlmodel.PageAttributeAssign, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) PageAttributeUnassign(ctx context.Context, args struct {
	AttributeIds []string
	PageTypeID   string
}) (*gqlmodel.PageAttributeUnassign, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) PageReorderAttributeValues(ctx context.Context, args struct {
	AttributeID string
	Moves       []*gqlmodel.ReorderInput
	PageID      string
}) (*gqlmodel.PageReorderAttributeValues, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Page(ctx context.Context, id *string, slug *string) (*gqlmodel.Page, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Pages(ctx context.Context, args struct {
	SortBy *gqlmodel.PageSortingInput
	Filter *gqlmodel.PageFilterInput
	Before *string
	After  *string
	First  *int
	Last   *int
}) (*gqlmodel.PageCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

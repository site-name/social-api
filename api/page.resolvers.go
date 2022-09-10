package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
)

func (r *Resolver) PageCreate(ctx context.Context, args struct{ Input PageCreateInput }) (*PageCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) PageDelete(ctx context.Context, args struct{ Id string }) (*PageDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) PageBulkDelete(ctx context.Context, args struct{ Ids []*string }) (*PageBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) PageBulkPublish(ctx context.Context, args struct {
	Ids         []*string
	IsPublished bool
}) (*PageBulkPublish, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) PageUpdate(ctx context.Context, args struct {
	Id    string
	Input PageInput
}) (*PageUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) PageTranslate(ctx context.Context, args struct {
	Id           string
	Input        PageTranslationInput
	LanguageCode LanguageCodeEnum
}) (*PageTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) PageAttributeAssign(ctx context.Context, args struct {
	AttributeIds []string
	PageTypeID   string
}) (*PageAttributeAssign, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) PageAttributeUnassign(ctx context.Context, args struct {
	AttributeIds []string
	PageTypeID   string
}) (*PageAttributeUnassign, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) PageReorderAttributeValues(ctx context.Context, args struct {
	AttributeID string
	Moves       []*ReorderInput
	PageID      string
}) (*PageReorderAttributeValues, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Page(ctx context.Context, id *string, slug *string) (*Page, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Pages(ctx context.Context, args struct {
	SortBy *PageSortingInput
	Filter *PageFilterInput
	Before *string
	After  *string
	First  *int
	Last   *int
}) (*PageCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

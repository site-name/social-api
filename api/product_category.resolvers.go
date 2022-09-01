package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/api/gqlmodel"
)

func (r *Resolver) CategoryCreate(ctx context.Context, args struct {
	input  gqlmodel.CategoryInput
	parent *string
}) (*gqlmodel.CategoryCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CategoryDelete(ctx context.Context, args struct{ id string }) (*gqlmodel.CategoryDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CategoryBulkDelete(ctx context.Context, args struct{ ids []*string }) (*gqlmodel.CategoryBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CategoryUpdate(ctx context.Context, args struct {
	id    string
	input gqlmodel.CategoryInput
}) (*gqlmodel.CategoryUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CategoryTranslate(ctx context.Context, args struct {
	id           string
	input        gqlmodel.TranslationInput
	languageCode gqlmodel.LanguageCodeEnum
}) (*gqlmodel.CategoryTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Categories(ctx context.Context, args struct {
	filter *gqlmodel.CategoryFilterInput
	sortBy *gqlmodel.CategorySortingInput
	level  *int
	before *string
	after  *string
	first  *int
	last   *int
}) (*gqlmodel.CategoryCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Category(ctx context.Context, args struct {
	id   *string
	slug *string
}) (*gqlmodel.Category, error) {
	panic(fmt.Errorf("not implemented"))
}

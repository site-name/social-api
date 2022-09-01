package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/api/gqlmodel"
)

func (r *Resolver) CategoryCreate(ctx context.Context, args struct {
	Input  gqlmodel.CategoryInput
	Parent *string
}) (*gqlmodel.CategoryCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CategoryDelete(ctx context.Context, args struct{ Id string }) (*gqlmodel.CategoryDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CategoryBulkDelete(ctx context.Context, args struct{ Ids []*string }) (*gqlmodel.CategoryBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CategoryUpdate(ctx context.Context, args struct {
	Id    string
	Input gqlmodel.CategoryInput
}) (*gqlmodel.CategoryUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CategoryTranslate(ctx context.Context, args struct {
	Id           string
	Input        gqlmodel.TranslationInput
	LanguageCode gqlmodel.LanguageCodeEnum
}) (*gqlmodel.CategoryTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Categories(ctx context.Context, args struct {
	Filter *gqlmodel.CategoryFilterInput
	SortBy *gqlmodel.CategorySortingInput
	Level  *int
	Before *string
	After  *string
	First  *int
	Last   *int
}) (*gqlmodel.CategoryCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Category(ctx context.Context, args struct {
	Id   *string
	Slug *string
}) (*gqlmodel.Category, error) {
	panic(fmt.Errorf("not implemented"))
}

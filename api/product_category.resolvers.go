package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
)

func (r *Resolver) CategoryCreate(ctx context.Context, args struct {
	Input  CategoryInput
	Parent *string
}) (*CategoryCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CategoryDelete(ctx context.Context, args struct{ Id string }) (*CategoryDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CategoryBulkDelete(ctx context.Context, args struct{ Ids []*string }) (*CategoryBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CategoryUpdate(ctx context.Context, args struct {
	Id    string
	Input CategoryInput
}) (*CategoryUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CategoryTranslate(ctx context.Context, args struct {
	Id           string
	Input        TranslationInput
	LanguageCode LanguageCodeEnum
}) (*CategoryTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Categories(ctx context.Context, args struct {
	Filter *CategoryFilterInput
	SortBy *CategorySortingInput
	Level  *int
	Before *string
	After  *string
	First  *int
	Last   *int
}) (*CategoryCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Category(ctx context.Context, args struct {
	Id   *string
	Slug *string
}) (*Category, error) {
	panic(fmt.Errorf("not implemented"))
}

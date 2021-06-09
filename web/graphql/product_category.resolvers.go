package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
)

func (r *mutationResolver) CategoryCreate(ctx context.Context, input CategoryInput, parent *string) (*CategoryCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CategoryDelete(ctx context.Context, id string) (*CategoryDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CategoryBulkDelete(ctx context.Context, ids []*string) (*CategoryBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CategoryUpdate(ctx context.Context, id string, input CategoryInput) (*CategoryUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CategoryTranslate(ctx context.Context, id string, input TranslationInput, languageCode LanguageCodeEnum) (*CategoryTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Categories(ctx context.Context, filter *CategoryFilterInput, sortBy *CategorySortingInput, level *int, before *string, after *string, first *int, last *int) (*CategoryCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Category(ctx context.Context, id *string, slug *string) (*Category, error) {
	panic(fmt.Errorf("not implemented"))
}

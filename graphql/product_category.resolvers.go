package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	graphql1 "github.com/sitename/sitename/graphql/generated"
	"github.com/sitename/sitename/graphql/gqlmodel"
	"github.com/sitename/sitename/model"
)

func (r *categoryResolver) Description(ctx context.Context, obj *gqlmodel.Category) (model.StringInterface, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *categoryResolver) Parent(ctx context.Context, obj *gqlmodel.Category) (*gqlmodel.Category, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *categoryResolver) Translation(ctx context.Context, obj *gqlmodel.Category, languageCode gqlmodel.LanguageCodeEnum) (*gqlmodel.CategoryTranslation, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CategoryCreate(ctx context.Context, input gqlmodel.CategoryInput, parent *string) (*gqlmodel.CategoryCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CategoryDelete(ctx context.Context, id string) (*gqlmodel.CategoryDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CategoryBulkDelete(ctx context.Context, ids []*string) (*gqlmodel.CategoryBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CategoryUpdate(ctx context.Context, id string, input gqlmodel.CategoryInput) (*gqlmodel.CategoryUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CategoryTranslate(ctx context.Context, id string, input gqlmodel.TranslationInput, languageCode gqlmodel.LanguageCodeEnum) (*gqlmodel.CategoryTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Categories(ctx context.Context, filter *gqlmodel.CategoryFilterInput, sortBy *gqlmodel.CategorySortingInput, level *int, before *string, after *string, first *int, last *int) (*gqlmodel.CategoryCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Category(ctx context.Context, id *string, slug *string) (*gqlmodel.Category, error) {
	panic(fmt.Errorf("not implemented"))
}

// Category returns graphql1.CategoryResolver implementation.
func (r *Resolver) Category() graphql1.CategoryResolver { return &categoryResolver{r} }

type categoryResolver struct{ *Resolver }

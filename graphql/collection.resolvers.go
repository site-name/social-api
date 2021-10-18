package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/graphql/gqlmodel"
)

func (r *mutationResolver) CollectionAddProducts(ctx context.Context, collectionID string, products []*string) (*gqlmodel.CollectionAddProducts, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CollectionCreate(ctx context.Context, input gqlmodel.CollectionCreateInput) (*gqlmodel.CollectionCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CollectionDelete(ctx context.Context, id string) (*gqlmodel.CollectionDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CollectionReorderProducts(ctx context.Context, collectionID string, moves []*gqlmodel.MoveProductInput) (*gqlmodel.CollectionReorderProducts, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CollectionBulkDelete(ctx context.Context, ids []*string) (*gqlmodel.CollectionBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CollectionRemoveProducts(ctx context.Context, collectionID string, products []*string) (*gqlmodel.CollectionRemoveProducts, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CollectionUpdate(ctx context.Context, id string, input gqlmodel.CollectionInput) (*gqlmodel.CollectionUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CollectionTranslate(ctx context.Context, id string, input gqlmodel.TranslationInput, languageCode gqlmodel.LanguageCodeEnum) (*gqlmodel.CollectionTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CollectionChannelListingUpdate(ctx context.Context, id string, input gqlmodel.CollectionChannelListingUpdateInput) (*gqlmodel.CollectionChannelListingUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Collection(ctx context.Context, id *string, slug *string, channel *string) (*gqlmodel.Collection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Collections(ctx context.Context, filter *gqlmodel.CollectionFilterInput, sortBy *gqlmodel.CollectionSortingInput, channel *string, before *string, after *string, first *int, last *int) (*gqlmodel.CollectionCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

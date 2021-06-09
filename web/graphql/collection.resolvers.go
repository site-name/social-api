package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
)

func (r *mutationResolver) CollectionAddProducts(ctx context.Context, collectionID string, products []*string) (*CollectionAddProducts, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CollectionCreate(ctx context.Context, input CollectionCreateInput) (*CollectionCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CollectionDelete(ctx context.Context, id string) (*CollectionDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CollectionReorderProducts(ctx context.Context, collectionID string, moves []*MoveProductInput) (*CollectionReorderProducts, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CollectionBulkDelete(ctx context.Context, ids []*string) (*CollectionBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CollectionRemoveProducts(ctx context.Context, collectionID string, products []*string) (*CollectionRemoveProducts, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CollectionUpdate(ctx context.Context, id string, input CollectionInput) (*CollectionUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CollectionTranslate(ctx context.Context, id string, input TranslationInput, languageCode LanguageCodeEnum) (*CollectionTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CollectionChannelListingUpdate(ctx context.Context, id string, input CollectionChannelListingUpdateInput) (*CollectionChannelListingUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Collection(ctx context.Context, id *string, slug *string, channel *string) (*Collection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Collections(ctx context.Context, filter *CollectionFilterInput, sortBy *CollectionSortingInput, channel *string, before *string, after *string, first *int, last *int) (*CollectionCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

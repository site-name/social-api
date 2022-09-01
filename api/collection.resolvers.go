package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/api/gqlmodel"
)

func (r *Resolver) CollectionAddProducts(ctx context.Context, args struct {
	collectionID string
	products     []*string
}) (*gqlmodel.CollectionAddProducts, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CollectionCreate(ctx context.Context, args struct {
	input gqlmodel.CollectionCreateInput
}) (*gqlmodel.CollectionCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CollectionDelete(ctx context.Context, args struct{ id string }) (*gqlmodel.CollectionDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CollectionReorderProducts(ctx context.Context, args struct {
	collectionID string
	moves        []*gqlmodel.MoveProductInput
}) (*gqlmodel.CollectionReorderProducts, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CollectionBulkDelete(ctx context.Context, args struct{ ids []*string }) (*gqlmodel.CollectionBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CollectionRemoveProducts(ctx context.Context, args struct {
	collectionID string
	products     []*string
}) (*gqlmodel.CollectionRemoveProducts, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CollectionUpdate(ctx context.Context, args struct {
	id    string
	input gqlmodel.CollectionInput
}) (*gqlmodel.CollectionUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CollectionTranslate(ctx context.Context, args struct {
	id           string
	input        gqlmodel.TranslationInput
	languageCode gqlmodel.LanguageCodeEnum
}) (*gqlmodel.CollectionTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CollectionChannelListingUpdate(ctx context.Context, args struct {
	id    string
	input gqlmodel.CollectionChannelListingUpdateInput
}) (*gqlmodel.CollectionChannelListingUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Collection(ctx context.Context, args struct {
	id      *string
	slug    *string
	channel *string
}) (*gqlmodel.Collection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Collections(ctx context.Context, args struct {
	filter  *gqlmodel.CollectionFilterInput
	sortBy  *gqlmodel.CollectionSortingInput
	channel *string
	before  *string
	after   *string
	first   *int
	last    *int
}) (*gqlmodel.CollectionCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

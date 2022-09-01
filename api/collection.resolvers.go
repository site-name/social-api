package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/api/gqlmodel"
)

func (r *Resolver) CollectionAddProducts(ctx context.Context, args struct {
	CollectionID string
	Products     []string
}) (*gqlmodel.CollectionAddProducts, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CollectionCreate(ctx context.Context, args struct {
	Input gqlmodel.CollectionCreateInput
}) (*gqlmodel.CollectionCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CollectionDelete(ctx context.Context, args struct{ Id string }) (*gqlmodel.CollectionDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CollectionReorderProducts(ctx context.Context, args struct {
	CollectionID string
	Moves        []*gqlmodel.MoveProductInput
}) (*gqlmodel.CollectionReorderProducts, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CollectionBulkDelete(ctx context.Context, args struct{ Ids []string }) (*gqlmodel.CollectionBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CollectionRemoveProducts(ctx context.Context, args struct {
	CollectionID string
	Products     []string
}) (*gqlmodel.CollectionRemoveProducts, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CollectionUpdate(ctx context.Context, args struct {
	Id    string
	Input gqlmodel.CollectionInput
}) (*gqlmodel.CollectionUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CollectionTranslate(ctx context.Context, args struct {
	Id           string
	Input        gqlmodel.TranslationInput
	LanguageCode gqlmodel.LanguageCodeEnum
}) (*gqlmodel.CollectionTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CollectionChannelListingUpdate(ctx context.Context, args struct {
	Id    string
	Input gqlmodel.CollectionChannelListingUpdateInput
}) (*gqlmodel.CollectionChannelListingUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Collection(ctx context.Context, args struct {
	Id      *string
	Slug    *string
	Channel *string
}) (*gqlmodel.Collection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Collections(ctx context.Context, args struct {
	Filter  *gqlmodel.CollectionFilterInput
	SortBy  *gqlmodel.CollectionSortingInput
	Channel *string
	Before  *string
	After   *string
	First   *int
	Last    *int
}) (*gqlmodel.CollectionCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

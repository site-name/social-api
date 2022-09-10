package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
)

func (r *Resolver) CollectionAddProducts(ctx context.Context, args struct {
	CollectionID string
	Products     []string
}) (*CollectionAddProducts, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CollectionCreate(ctx context.Context, args struct {
	Input CollectionCreateInput
}) (*CollectionCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CollectionDelete(ctx context.Context, args struct{ Id string }) (*CollectionDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CollectionReorderProducts(ctx context.Context, args struct {
	CollectionID string
	Moves        []*MoveProductInput
}) (*CollectionReorderProducts, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CollectionBulkDelete(ctx context.Context, args struct{ Ids []string }) (*CollectionBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CollectionRemoveProducts(ctx context.Context, args struct {
	CollectionID string
	Products     []string
}) (*CollectionRemoveProducts, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CollectionUpdate(ctx context.Context, args struct {
	Id    string
	Input CollectionInput
}) (*CollectionUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CollectionTranslate(ctx context.Context, args struct {
	Id           string
	Input        TranslationInput
	LanguageCode LanguageCodeEnum
}) (*CollectionTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) CollectionChannelListingUpdate(ctx context.Context, args struct {
	Id    string
	Input CollectionChannelListingUpdateInput
}) (*CollectionChannelListingUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Collection(ctx context.Context, args struct {
	Id      *string
	Slug    *string
	Channel *string
}) (*Collection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Collections(ctx context.Context, args struct {
	Filter  *CollectionFilterInput
	SortBy  *CollectionSortingInput
	Channel *string
	Before  *string
	After   *string
	First   *int
	Last    *int
}) (*CollectionCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

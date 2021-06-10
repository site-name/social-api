package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
)

func (r *mutationResolver) SaleCreate(ctx context.Context, input SaleInput) (*SaleCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SaleDelete(ctx context.Context, id string) (*SaleDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SaleBulkDelete(ctx context.Context, ids []*string) (*SaleBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SaleUpdate(ctx context.Context, id string, input SaleInput) (*SaleUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SaleCataloguesAdd(ctx context.Context, id string, input CatalogueInput) (*SaleAddCatalogues, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SaleCataloguesRemove(ctx context.Context, id string, input CatalogueInput) (*SaleRemoveCatalogues, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SaleTranslate(ctx context.Context, id string, input NameTranslationInput, languageCode LanguageCodeEnum) (*SaleTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SaleChannelListingUpdate(ctx context.Context, id string, input SaleChannelListingInput) (*SaleChannelListingUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Sale(ctx context.Context, id string, channel *string) (*Sale, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Sales(ctx context.Context, filter *SaleFilterInput, sortBy *SaleSortingInput, query *string, channel *string, before *string, after *string, first *int, last *int) (*SaleCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

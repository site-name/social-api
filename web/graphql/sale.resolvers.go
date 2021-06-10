package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/web/graphql/gqlmodel"
)

func (r *mutationResolver) SaleCreate(ctx context.Context, input gqlmodel.SaleInput) (*gqlmodel.SaleCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SaleDelete(ctx context.Context, id string) (*gqlmodel.SaleDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SaleBulkDelete(ctx context.Context, ids []*string) (*gqlmodel.SaleBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SaleUpdate(ctx context.Context, id string, input gqlmodel.SaleInput) (*gqlmodel.SaleUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SaleCataloguesAdd(ctx context.Context, id string, input gqlmodel.CatalogueInput) (*gqlmodel.SaleAddCatalogues, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SaleCataloguesRemove(ctx context.Context, id string, input gqlmodel.CatalogueInput) (*gqlmodel.SaleRemoveCatalogues, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SaleTranslate(ctx context.Context, id string, input gqlmodel.NameTranslationInput, languageCode gqlmodel.LanguageCodeEnum) (*gqlmodel.SaleTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SaleChannelListingUpdate(ctx context.Context, id string, input gqlmodel.SaleChannelListingInput) (*gqlmodel.SaleChannelListingUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Sale(ctx context.Context, id string, channel *string) (*gqlmodel.Sale, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Sales(ctx context.Context, filter *gqlmodel.SaleFilterInput, sortBy *gqlmodel.SaleSortingInput, query *string, channel *string, before *string, after *string, first *int, last *int) (*gqlmodel.SaleCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

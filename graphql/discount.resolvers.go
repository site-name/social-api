package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/graphql/gqlmodel"
)

func (r *mutationResolver) VoucherCreate(ctx context.Context, input gqlmodel.VoucherInput) (*gqlmodel.VoucherCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) VoucherDelete(ctx context.Context, id string) (*gqlmodel.VoucherDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) VoucherBulkDelete(ctx context.Context, ids []*string) (*gqlmodel.VoucherBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) VoucherUpdate(ctx context.Context, id string, input gqlmodel.VoucherInput) (*gqlmodel.VoucherUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) VoucherCataloguesAdd(ctx context.Context, id string, input gqlmodel.CatalogueInput) (*gqlmodel.VoucherAddCatalogues, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) VoucherCataloguesRemove(ctx context.Context, id string, input gqlmodel.CatalogueInput) (*gqlmodel.VoucherRemoveCatalogues, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) VoucherTranslate(ctx context.Context, id string, input gqlmodel.NameTranslationInput, languageCode gqlmodel.LanguageCodeEnum) (*gqlmodel.VoucherTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) VoucherChannelListingUpdate(ctx context.Context, id string, input gqlmodel.VoucherChannelListingInput) (*gqlmodel.VoucherChannelListingUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Voucher(ctx context.Context, id string, channel *string) (*gqlmodel.Voucher, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Vouchers(ctx context.Context, filter *gqlmodel.VoucherFilterInput, sortBy *gqlmodel.VoucherSortingInput, query *string, channel *string, before *string, after *string, first *int, last *int) (*gqlmodel.VoucherCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}
package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
)

func (r *mutationResolver) VoucherCreate(ctx context.Context, input VoucherInput) (*VoucherCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) VoucherDelete(ctx context.Context, id string) (*VoucherDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) VoucherBulkDelete(ctx context.Context, ids []*string) (*VoucherBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) VoucherUpdate(ctx context.Context, id string, input VoucherInput) (*VoucherUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) VoucherCataloguesAdd(ctx context.Context, id string, input CatalogueInput) (*VoucherAddCatalogues, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) VoucherCataloguesRemove(ctx context.Context, id string, input CatalogueInput) (*VoucherRemoveCatalogues, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) VoucherTranslate(ctx context.Context, id string, input NameTranslationInput, languageCode LanguageCodeEnum) (*VoucherTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) VoucherChannelListingUpdate(ctx context.Context, id string, input VoucherChannelListingInput) (*VoucherChannelListingUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Voucher(ctx context.Context, id string, channel *string) (*Voucher, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Vouchers(ctx context.Context, filter *VoucherFilterInput, sortBy *VoucherSortingInput, query *string, channel *string, before *string, after *string, first *int, last *int) (*VoucherCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

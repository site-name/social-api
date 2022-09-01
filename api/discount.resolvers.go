package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/api/gqlmodel"
)

func (r *Resolver) VoucherCreate(ctx context.Context, args struct{ Input gqlmodel.VoucherInput }) (*gqlmodel.VoucherCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) VoucherDelete(ctx context.Context, args struct{ Id string }) (*gqlmodel.VoucherDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) VoucherBulkDelete(ctx context.Context, args struct{ Ids []*string }) (*gqlmodel.VoucherBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) VoucherUpdate(ctx context.Context, args struct {
	Id    string
	Input gqlmodel.VoucherInput
}) (*gqlmodel.VoucherUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) VoucherCataloguesAdd(ctx context.Context, args struct {
	Id    string
	Input gqlmodel.CatalogueInput
}) (*gqlmodel.VoucherAddCatalogues, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) VoucherCataloguesRemove(ctx context.Context, args struct {
	Id    string
	Input gqlmodel.CatalogueInput
}) (*gqlmodel.VoucherRemoveCatalogues, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) VoucherTranslate(ctx context.Context, args struct {
	Id           string
	Input        gqlmodel.NameTranslationInput
	languageCode gqlmodel.LanguageCodeEnum
}) (*gqlmodel.VoucherTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) VoucherChannelListingUpdate(ctx context.Context, args struct {
	Id    string
	Input gqlmodel.VoucherChannelListingInput
}) (*gqlmodel.VoucherChannelListingUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Voucher(ctx context.Context, args struct {
	Id      string
	Channel *string
}) (*gqlmodel.Voucher, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Vouchers(ctx context.Context, args struct {
	Filter  *gqlmodel.VoucherFilterInput
	SortBy  *gqlmodel.VoucherSortingInput
	Query   *string
	Channel *string
	Before  *string
	After   *string
	First   *int
	Last    *int
}) (*gqlmodel.VoucherCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
)

func (r *Resolver) VoucherCreate(ctx context.Context, args struct{ Input VoucherInput }) (*VoucherCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) VoucherDelete(ctx context.Context, args struct{ Id string }) (*VoucherDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) VoucherBulkDelete(ctx context.Context, args struct{ Ids []string }) (*VoucherBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) VoucherUpdate(ctx context.Context, args struct {
	Id    string
	Input VoucherInput
}) (*VoucherUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) VoucherCataloguesAdd(ctx context.Context, args struct {
	Id    string
	Input CatalogueInput
}) (*VoucherAddCatalogues, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) VoucherCataloguesRemove(ctx context.Context, args struct {
	Id    string
	Input CatalogueInput
}) (*VoucherRemoveCatalogues, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) VoucherTranslate(ctx context.Context, args struct {
	Id           string
	Input        NameTranslationInput
	LanguageCode LanguageCodeEnum
}) (*VoucherTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) VoucherChannelListingUpdate(ctx context.Context, args struct {
	Id    string
	Input VoucherChannelListingInput
}) (*VoucherChannelListingUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Voucher(ctx context.Context, args struct {
	Id      string
	Channel *string
}) (*Voucher, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Vouchers(ctx context.Context, args struct {
	Filter  *VoucherFilterInput
	SortBy  *VoucherSortingInput
	Query   *string
	Channel *string
	GraphqlParams
}) (*VoucherCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

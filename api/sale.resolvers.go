package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
)

func (r *Resolver) SaleCreate(ctx context.Context, args struct{ Input SaleInput }) (*SaleCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) SaleDelete(ctx context.Context, args struct{ Id string }) (*SaleDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) SaleBulkDelete(ctx context.Context, args struct{ Ids []*string }) (*SaleBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) SaleUpdate(ctx context.Context, args struct {
	Id    string
	Input SaleInput
}) (*SaleUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) SaleCataloguesAdd(ctx context.Context, args struct {
	Id    string
	Input CatalogueInput
}) (*SaleAddCatalogues, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) SaleCataloguesRemove(ctx context.Context, args struct {
	Id    string
	Input CatalogueInput
}) (*SaleRemoveCatalogues, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) SaleTranslate(ctx context.Context, args struct {
	Id           string
	Input        NameTranslationInput
	LanguageCode LanguageCodeEnum
}) (*SaleTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) SaleChannelListingUpdate(ctx context.Context, args struct {
	Id    string
	Input SaleChannelListingInput
}) (*SaleChannelListingUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Sale(ctx context.Context, args struct {
	Id      string
	Channel *string
}) (*Sale, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Sales(ctx context.Context, args struct {
	Filter  *SaleFilterInput
	SortBy  *SaleSortingInput
	Query   *string
	Channel *string
	Before  *string
	After   *string
	First   *int
	Last    *int
}) (*SaleCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

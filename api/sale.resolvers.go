package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sitename/sitename/api/gqlmodel"
)

func (r *Resolver) SaleCreate(ctx context.Context, args struct{ Input gqlmodel.SaleInput }) (*gqlmodel.SaleCreate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) SaleDelete(ctx context.Context, args struct{ Id string }) (*gqlmodel.SaleDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) SaleBulkDelete(ctx context.Context, args struct{ Ids []*string }) (*gqlmodel.SaleBulkDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) SaleUpdate(ctx context.Context, args struct {
	Id    string
	Input gqlmodel.SaleInput
}) (*gqlmodel.SaleUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) SaleCataloguesAdd(ctx context.Context, args struct {
	Id    string
	Input gqlmodel.CatalogueInput
}) (*gqlmodel.SaleAddCatalogues, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) SaleCataloguesRemove(ctx context.Context, args struct {
	Id    string
	Input gqlmodel.CatalogueInput
}) (*gqlmodel.SaleRemoveCatalogues, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) SaleTranslate(ctx context.Context, args struct {
	Id           string
	Input        gqlmodel.NameTranslationInput
	LanguageCode gqlmodel.LanguageCodeEnum
}) (*gqlmodel.SaleTranslate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) SaleChannelListingUpdate(ctx context.Context, args struct {
	Id    string
	Input gqlmodel.SaleChannelListingInput
}) (*gqlmodel.SaleChannelListingUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Sale(ctx context.Context, args struct {
	Id      string
	Channel *string
}) (*gqlmodel.Sale, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Sales(ctx context.Context, args struct {
	Filter  *gqlmodel.SaleFilterInput
	SortBy  *gqlmodel.SaleSortingInput
	Query   *string
	Channel *string
	Before  *string
	After   *string
	First   *int
	Last    *int
}) (*gqlmodel.SaleCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

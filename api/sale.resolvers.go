package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/web"
)

// NOTE: Refer to ./schemas/sale.graphqls for details on directives used.
func (r *Resolver) SaleCreate(ctx context.Context, args struct{ Input SaleInput }) (*SaleCreate, error) {
	// validate params
	appErr := args.Input.Validate()
	if appErr != nil {
		return nil, appErr
	}

	sale := new(model.Sale)
	args.Input.PatchSale(sale)

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	// begin transaction
	transaction := embedCtx.App.Srv().Store.GetMaster().Begin()
	defer transaction.Rollback()

	// insert sale
	sale, err := embedCtx.App.Srv().Store.DiscountSale().Upsert(transaction, sale)
	if err != nil {
		return nil, model.NewAppError("SaleCreate", "app.discount.create_sale.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	// save m2m relations
	appErr = embedCtx.App.Srv().DiscountService().AddSaleRelations(transaction, sale.Id, args.Input.Products, args.Input.Variants, args.Input.Categories, args.Input.Collections)
	if appErr != nil {
		return nil, appErr
	}

	err = transaction.Commit().Error
	if err != nil {
		return nil, model.NewAppError("SaleCreate", app.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	embedCtx.App.Srv().DiscountService().FetchCatalogueInfo(*sale)

}

func (r *Resolver) SaleDelete(ctx context.Context, args struct{ Id string }) (*SaleDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) SaleBulkDelete(ctx context.Context, args struct{ Ids []string }) (*SaleBulkDelete, error) {
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
	GraphqlParams
}) (*SaleCountableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

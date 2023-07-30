package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/samber/lo"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/web"
)

// NOTE: Refer to ./schemas/sale.graphqls for details on directives used.
func (r *Resolver) SaleCreate(ctx context.Context, args struct{ Input SaleInput }) (*SaleCreate, error) {
	// validate params
	appErr := args.Input.Validate("SaleCreate")
	if appErr != nil {
		return nil, appErr
	}

	sale := new(model.Sale)
	args.Input.PatchSale(sale)

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	// begin transaction
	transaction := embedCtx.App.Srv().Store.GetMaster().Begin()
	if transaction.Error != nil {
		return nil, model.NewAppError("Resolver.SaleCreate", app.ErrorCreatingTransactionErrorID, nil, transaction.Error.Error(), http.StatusInternalServerError)
	}
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

	appErr = embedCtx.App.Srv().ProductService().UpdateProductsDiscountedPricesOfDiscount(transaction, sale)
	if appErr != nil {
		return nil, appErr
	}

	// commit transaction
	err = transaction.Commit().Error
	if err != nil {
		return nil, model.NewAppError("SaleCreate", app.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	catalogInfo, appErr := embedCtx.App.Srv().DiscountService().FetchCatalogueInfo(*sale)
	if appErr != nil {
		return nil, appErr
	}

	pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()
	_, appErr = pluginMng.SaleCreated(*sale, catalogInfo)
	if appErr != nil {
		return nil, appErr
	}

	return &SaleCreate{
		Sale: systemSaleToGraphqlSale(sale),
	}, nil
}

// NOTE: Refer to ./schemas/sale.graphqls for details on directives used.
func (r *Resolver) SaleDelete(ctx context.Context, args struct{ Id string }) (*SaleDelete, error) {
	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("SaleDelete", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "please provide valid sale id", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	// begin transaction
	transaction := embedCtx.App.Srv().Store.GetMaster().Begin()
	if transaction.Error != nil {
		return nil, model.NewAppError("Resolver.SaleDelete", app.ErrorCreatingTransactionErrorID, nil, transaction.Error.Error(), http.StatusInternalServerError)
	}
	defer transaction.Rollback()

	_, err := embedCtx.App.Srv().Store.DiscountSale().Delete(transaction, &model.SaleFilterOption{
		Conditions: squirrel.Eq{model.SaleTableName + ".Id": args.Id},
	})
	if err != nil {
		return nil, model.NewAppError("SaleDelete", "app.discount.delete_sale_by_id.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	sale := &model.Sale{Id: args.Id}

	appErr := embedCtx.App.Srv().ProductService().UpdateProductsDiscountedPricesOfDiscount(transaction, sale)
	if appErr != nil {
		return nil, appErr
	}

	// commit transaction
	err = transaction.Commit().Error
	if err != nil {
		return nil, model.NewAppError("SaleCreate", app.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	catalogInfo, appErr := embedCtx.App.Srv().DiscountService().FetchCatalogueInfo(*sale)
	if appErr != nil {
		return nil, appErr
	}

	pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()
	_, appErr = pluginMng.SaleDeleted(*sale, catalogInfo)
	if appErr != nil {
		return nil, appErr
	}

	return &SaleDelete{
		Sale: systemSaleToGraphqlSale(sale),
	}, nil
}

// NOTE: Refer to ./schemas/sale.graphqls for details on directives used.
func (r *Resolver) SaleBulkDelete(ctx context.Context, args struct{ Ids []string }) (*SaleBulkDelete, error) {
	// validate params
	if !lo.EveryBy(args.Ids, model.IsValidId) {
		return nil, model.NewAppError("SaleBulkDelete", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "ids"}, "please provide valid sale ids", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	numDeleted, err := embedCtx.App.Srv().Store.DiscountSale().Delete(nil, &model.SaleFilterOption{
		Conditions: squirrel.Eq{model.SaleTableName + ".Id": args.Ids},
	})
	if err != nil {
		return nil, model.NewAppError("SaleBulkDelete", "app.discount.delete_sale_by_id.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return &SaleBulkDelete{
		Count: int32(numDeleted),
	}, nil
}

// NOTE: Refer to ./schemas/sale.graphqls for details on directives used.
func (r *Resolver) SaleUpdate(ctx context.Context, args struct {
	Id    string
	Input SaleInput
}) (*SaleUpdate, error) {
	// validate params
	appErr := args.Input.Validate("SaleUpdate")
	if appErr != nil {
		return nil, appErr
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	// begin transaction
	transaction := embedCtx.App.Srv().Store.GetMaster().Begin()
	if transaction.Error != nil {
		return nil, model.NewAppError("Resolver.SaleDelete", app.ErrorCreatingTransactionErrorID, nil, transaction.Error.Error(), http.StatusInternalServerError)
	}
	defer transaction.Rollback()

	sales, appErr := embedCtx.App.Srv().DiscountService().FilterSalesByOption(&model.SaleFilterOption{
		Conditions: squirrel.Expr(model.SaleTableName+".Id = ?", args.Id),
	})
	if appErr != nil {
		return nil, appErr
	}

	// really update sale
	sale := sales[0]
	args.Input.PatchSale(sale)

	// fetch catalogue before updating sale
	catalogBeforeUpdate, appErr := embedCtx.App.Srv().DiscountService().FetchCatalogueInfo(*sale)
	if appErr != nil {
		return nil, appErr
	}

	sale, appErr = embedCtx.App.Srv().DiscountService().UpsertSale(transaction, sale)
	if appErr != nil {
		return nil, appErr
	}

	appErr = embedCtx.App.Srv().ProductService().UpdateProductsDiscountedPricesOfDiscount(transaction, sale)
	if appErr != nil {
		return nil, appErr
	}

	// commit transaction
	err := transaction.Commit().Error
	if err != nil {
		return nil, model.NewAppError("SaleCreate", app.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	// fetch catalogue after updating sale
	catalogueAfterUpdate, appErr := embedCtx.App.Srv().DiscountService().FetchCatalogueInfo(*sale)
	if appErr != nil {
		return nil, appErr
	}

	pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()
	_, appErr = pluginMng.SaleUpdated(*sale, catalogBeforeUpdate, catalogueAfterUpdate)
	if appErr != nil {
		return nil, appErr
	}

	return &SaleUpdate{
		Sale: systemSaleToGraphqlSale(sale),
	}, nil
}

// NOTE: Refer to ./schemas/sale.graphqls for details on directives used.
func (r *Resolver) SaleCataloguesAdd(ctx context.Context, args struct {
	Id    string
	Input CatalogueInput
}) (*SaleAddCatalogues, error) {
	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("SaleCataloguesAdd", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Id"}, "please provide valid sale id", http.StatusBadRequest)
	}
	appErr := args.Input.Validate("SaleCataloguesAdd")
	if appErr != nil {
		return nil, appErr
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	// NOTE: only products that have variants can be added to sale.
	// So we have to verify that every given products have variant(s)
	if len(args.Input.Products) > 0 {
		productsWithNovariants, appErr := embedCtx.App.Srv().ProductService().ProductsByOption(&model.ProductFilterOption{
			Conditions:           squirrel.Eq{model.ProductTableName + ".Id": args.Input.Products},
			HasNoProductVariants: true,
		})
		if appErr != nil {
			return nil, appErr
		}
		if len(productsWithNovariants) > 0 {
			return nil, model.NewAppError("SaleCataloguesAdd", "app.discount.add_products_with_no_variants_to_sale.app_error", nil, "cant add products that have no variants to sale", http.StatusNotAcceptable)
		}
	}

	// begin transaction
	transaction := embedCtx.App.Srv().Store.GetMaster().Begin()
	if transaction.Error != nil {
		return nil, model.NewAppError("Resolver.SaleCataloguesAdd", app.ErrorCreatingTransactionErrorID, nil, transaction.Error.Error(), http.StatusInternalServerError)
	}
	defer transaction.Rollback()

	sale := &model.Sale{Id: args.Id}

	// fetch catalogue before updating sale
	catalogBeforeUpdate, appErr := embedCtx.App.Srv().DiscountService().FetchCatalogueInfo(*sale)
	if appErr != nil {
		return nil, appErr
	}

	appErr = embedCtx.App.Srv().DiscountService().AddSaleRelations(transaction, args.Id, args.Input.Products, nil, args.Input.Categories, args.Input.Collections)
	if appErr != nil {
		return nil, appErr
	}

	appErr = embedCtx.App.Srv().ProductService().UpdateProductsDiscountedPricesOfCatalogues(transaction, args.Input.Products, args.Input.Categories, args.Input.Collections, nil)
	if appErr != nil {
		return nil, appErr
	}

	// commit transaction
	err := transaction.Commit().Error
	if err != nil {
		return nil, model.NewAppError("SaleCataloguesAdd", app.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	// fetch catalogue after updating sale
	catalogAfterUpdate, appErr := embedCtx.App.Srv().DiscountService().FetchCatalogueInfo(*sale)
	if appErr != nil {
		return nil, appErr
	}

	pluginMng := embedCtx.App.Srv().PluginService().GetPluginManager()
	_, appErr = pluginMng.SaleUpdated(*sale, catalogBeforeUpdate, catalogAfterUpdate)
	if appErr != nil {
		return nil, appErr
	}

	return &SaleAddCatalogues{
		Sale: systemSaleToGraphqlSale(sale),
	}, nil
}

// NOTE: Refer to ./schemas/sale.graphqls for details on directives used.
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

// NOTE: Refer to ./schemas/sale.graphqls for details on directives used.
func (r *Resolver) SaleChannelListingUpdate(ctx context.Context, args struct {
	Id    string
	Input SaleChannelListingInput
}) (*SaleChannelListingUpdate, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Sale(ctx context.Context, args struct {
	Id      string
	Channel *string // TODO: Consider removing this field
}) (*Sale, error) {
	// validate params
	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("Resolve.Sale", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Id"}, "please provide valid sale id", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	sales, appErr := embedCtx.App.Srv().DiscountService().FilterSalesByOption(&model.SaleFilterOption{
		Conditions: squirrel.Eq{model.SaleTableName + ".Id": args.Id},
	})
	if appErr != nil {
		return nil, appErr
	}

	return systemSaleToGraphqlSale(sales[0]), nil
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

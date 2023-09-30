package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"net/http"
	"unsafe"

	"github.com/Masterminds/squirrel"
	"github.com/gosimple/slug"
	"github.com/samber/lo"
	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
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
		return nil, model.NewAppError("Resolver.SaleCreate", model.ErrorCreatingTransactionErrorID, nil, transaction.Error.Error(), http.StatusInternalServerError)
	}
	defer transaction.Rollback()

	// insert sale
	sale, err := embedCtx.App.Srv().Store.DiscountSale().Upsert(transaction, sale)
	if err != nil {
		return nil, model.NewAppError("SaleCreate", "app.discount.create_sale.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	// save m2m relations
	appErr = embedCtx.App.Srv().DiscountService().ToggleSaleRelations(transaction, sale.Id, args.Input.Products, args.Input.Variants, args.Input.Categories, args.Input.Collections, false)
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
		return nil, model.NewAppError("SaleCreate", model.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
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
		return nil, model.NewAppError("SaleDelete", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "please provide valid sale id", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	// begin transaction
	transaction := embedCtx.App.Srv().Store.GetMaster().Begin()
	if transaction.Error != nil {
		return nil, model.NewAppError("Resolver.SaleDelete", model.ErrorCreatingTransactionErrorID, nil, transaction.Error.Error(), http.StatusInternalServerError)
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
		return nil, model.NewAppError("SaleCreate", model.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
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
		return nil, model.NewAppError("SaleBulkDelete", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "ids"}, "please provide valid sale ids", http.StatusBadRequest)
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

	// really update sale
	sale := &model.Sale{Id: args.Id}
	args.Input.PatchSale(sale)

	// fetch catalogue before updating sale
	catalogBeforeUpdate, appErr := embedCtx.App.Srv().DiscountService().FetchCatalogueInfo(*sale)
	if appErr != nil {
		return nil, appErr
	}

	// begin transaction
	transaction := embedCtx.App.Srv().Store.GetMaster().Begin()
	if transaction.Error != nil {
		return nil, model.NewAppError("Resolver.SaleDelete", model.ErrorCreatingTransactionErrorID, nil, transaction.Error.Error(), http.StatusInternalServerError)
	}
	defer transaction.Rollback()

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
		return nil, model.NewAppError("SaleCreate", model.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
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
		return nil, model.NewAppError("SaleCataloguesAdd", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Id"}, "please provide valid sale id", http.StatusBadRequest)
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

	sale := &model.Sale{Id: args.Id}

	// fetch catalogue before updating sale
	catalogBeforeUpdate, appErr := embedCtx.App.Srv().DiscountService().FetchCatalogueInfo(*sale)
	if appErr != nil {
		return nil, appErr
	}

	// begin transaction
	transaction := embedCtx.App.Srv().Store.GetMaster().Begin()
	if transaction.Error != nil {
		return nil, model.NewAppError("Resolver.SaleCataloguesAdd", model.ErrorCreatingTransactionErrorID, nil, transaction.Error.Error(), http.StatusInternalServerError)
	}
	defer transaction.Rollback()

	// add sale relations
	appErr = embedCtx.App.Srv().DiscountService().ToggleSaleRelations(transaction, args.Id, args.Input.Products, nil, args.Input.Categories, args.Input.Collections, false)
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
		return nil, model.NewAppError("SaleCataloguesAdd", model.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
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
		Sale: &Sale{ID: args.Id},
	}, nil
}

// NOTE: Refer to ./schemas/sale.graphqls for details on directives used.
func (r *Resolver) SaleCataloguesRemove(ctx context.Context, args struct {
	Id    string
	Input CatalogueInput
}) (*SaleRemoveCatalogues, error) {
	// validate params
	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("SaleCataloguesRemove", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "please provide valid sale id", http.StatusBadRequest)
	}
	appErr := args.Input.Validate("SaleCataloguesRemove")
	if appErr != nil {
		return nil, appErr
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	sale := &model.Sale{Id: args.Id}

	// fetch catalogue before updating sale
	catalogBeforeUpdate, appErr := embedCtx.App.Srv().DiscountService().FetchCatalogueInfo(*sale)
	if appErr != nil {
		return nil, appErr
	}

	// create transaction
	transaction := embedCtx.App.Srv().Store.GetMaster().Begin()
	if transaction.Error != nil {
		return nil, model.NewAppError("SaleCataloguesRemove", model.ErrorCreatingTransactionErrorID, nil, transaction.Error.Error(), http.StatusInternalServerError)
	}
	defer transaction.Rollback()

	// remove sale relations
	appErr = embedCtx.App.Srv().DiscountService().ToggleSaleRelations(transaction, args.Id, args.Input.Products, nil, args.Input.Categories, args.Input.Collections, true)
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
		return nil, model.NewAppError("SaleCataloguesRemove", model.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
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

	return &SaleRemoveCatalogues{
		Sale: &Sale{ID: args.Id},
	}, nil
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
	// validate params
	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("SaleChannelListingUpdate", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "id"}, "please provide valid sale id", http.StatusBadRequest)
	}
	appErr := args.Input.Validate()
	if appErr != nil {
		return nil, appErr
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)

	// get sale by given id
	_, sales, appErr := embedCtx.App.Srv().DiscountService().FilterSalesByOption(&model.SaleFilterOption{
		Conditions: squirrel.Expr(model.SaleTableName+".Id = ?", args.Id),
	})
	if appErr != nil {
		return nil, appErr
	}
	sale := sales[0]

	// clean discount values
	allChannels, appErr := embedCtx.App.Srv().ChannelService().ChannelsByOption(&model.ChannelFilterOption{})
	if appErr != nil {
		return nil, appErr
	}
	// keys are channel ids, values are channels' according currency units
	channelCurrenciesMap := lo.SliceToMap(allChannels, func(c *model.Channel) (string, string) { return c.Id, c.Currency })

	for _, addChannelObj := range args.Input.AddChannels {
		if addChannelObj == nil || decimal.Decimal(addChannelObj.DiscountValue).IsZero() {
			continue // TODO: check if this is right
		}

		switch sale.Type {
		case model.DISCOUNT_VALUE_TYPE_FIXED:
			currencyPrecision, _ := goprices.GetCurrencyPrecision(channelCurrenciesMap[addChannelObj.ChannelID])
			roundedDiscountValue := decimal.Decimal(addChannelObj.DiscountValue).Round(int32(currencyPrecision))
			addChannelObj.DiscountValue = PositiveDecimal(roundedDiscountValue)

		case model.DISCOUNT_VALUE_TYPE_PERCENTAGE:
			if decimal.Decimal(addChannelObj.DiscountValue).GreaterThan(decimal.NewFromInt(100)) {
				// discount can't > 100%
				return nil, model.NewAppError("SaleChannelListingUpdate", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "add channels"}, "discount cannot be greater than 100%", http.StatusBadRequest)
			}
		}
	}

	// begin transaction
	tran := embedCtx.App.Srv().Store.GetMaster().Begin()
	if tran.Error != nil {
		return nil, model.NewAppError("SaleChannelListingUpdate", model.ErrorCreatingTransactionErrorID, nil, tran.Error.Error(), http.StatusInternalServerError)
	}
	defer tran.Rollback()

	// perform insert/delete in database:
	err := embedCtx.App.Srv().Store.DiscountSaleChannelListing().Delete(tran, &model.SaleChannelListingFilterOption{
		Conditions: squirrel.Eq{
			model.SaleChannelListingTableName + ".SaleID":    args.Id,
			model.SaleChannelListingTableName + ".ChannelID": args.Input.RemoveChannels,
		},
	})
	if err != nil {
		return nil, model.NewAppError("SaleChannelListingUpdate", "app.sale.error_delete_sale_channel_listings.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	listingsToAdd := lo.Map(args.Input.AddChannels, func(ac *SaleChannelListingAddInput, _ int) *model.SaleChannelListing {
		return &model.SaleChannelListing{
			SaleID:        args.Id,
			ChannelID:     ac.ChannelID,
			DiscountValue: (*decimal.Decimal)(unsafe.Pointer(&ac.DiscountValue)),
			Currency:      channelCurrenciesMap[ac.ChannelID],
		}
	})

	_, err = embedCtx.App.Srv().Store.DiscountSaleChannelListing().Upsert(tran, listingsToAdd)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrInvalidInput); ok {
			statusCode = http.StatusBadRequest
		}
		return nil, model.NewAppError("SaleChannelListingUpdate", "app.sale.add_sale_channel_listings.app_error", nil, err.Error(), statusCode)
	}

	appErr = embedCtx.App.Srv().ProductService().UpdateProductsDiscountedPricesOfDiscount(tran, sale)
	if appErr != nil {
		return nil, appErr
	}

	// commit transaction:
	err = tran.Commit().Error
	if err != nil {
		return nil, model.NewAppError("SaleChannelListingUpdate", model.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	return &SaleChannelListingUpdate{
		Sale: systemSaleToGraphqlSale(sale),
	}, nil
}

// TODO: Check if we need any role or permission to see this
func (r *Resolver) Sale(ctx context.Context, args struct {
	Id      string
	Channel *string // TODO: Consider removing this field
}) (*Sale, error) {
	// validate params
	if !model.IsValidId(args.Id) {
		return nil, model.NewAppError("Resolve.Sale", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Id"}, "please provide valid sale id", http.StatusBadRequest)
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	_, sales, appErr := embedCtx.App.Srv().DiscountService().FilterSalesByOption(&model.SaleFilterOption{
		Conditions: squirrel.Expr(model.SaleTableName+".Id = ?", args.Id),
	})
	if appErr != nil {
		return nil, appErr
	}

	return systemSaleToGraphqlSale(sales[0]), nil
}

type SalesArgs struct {
	Filter  *SaleFilterInput
	SortBy  *SaleSortingInput
	Channel *string // This is channel slug
	GraphqlParams
}

func (a *SalesArgs) parse() (*model.SaleFilterOption, *model.AppError) {
	var conditions squirrel.Sqlizer

	if a.Filter != nil {
		var appErr *model.AppError
		conditions, appErr = a.Filter.parse()
		if appErr != nil {
			return nil, appErr
		}
	}
	if a.SortBy != nil && !a.SortBy.Field.IsValid() {
		return nil, model.NewAppError("SalesArgs.SortBy", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "Field"}, "please provide valid field to sort", http.StatusBadRequest)
	}
	if a.Channel != nil && !slug.IsSlug(*a.Channel) {
		return nil, model.NewAppError("SalesArgs.Channel", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "channel"}, *a.Channel+" is not a valid channel slug", http.StatusBadRequest)
	}

	paginationValues, appErr := a.GraphqlParams.Parse("SalesArgs")
	if appErr != nil {
		return nil, appErr
	}

	res := &model.SaleFilterOption{
		Conditions:              conditions,
		GraphqlPaginationValues: *paginationValues,
		CountTotal:              true,
	}

	// in case no sort field is provided
	if res.GraphqlPaginationValues.OrderBy == "" {
		saleSortFields := saleSortFieldsMap[SaleSortFieldName].fields

		if a.SortBy != nil {
			saleSortFields = saleSortFieldsMap[a.SortBy.Field].fields

			// check if users want to sort sales by Values
			if a.SortBy.Field == SaleSortFieldValue {
				res.Annotate_Value = true
			}
		}

		saleOrder := a.GraphqlParams.orderDirection()
		res.GraphqlPaginationValues.OrderBy = saleSortFields.Map(func(_ int, item string) string { return item + " " + saleOrder }).Join(",")
	}

	return res, nil
}

// TODO: Check if we need any role or permission to see this
func (r *Resolver) Sales(ctx context.Context, args SalesArgs) (*SaleCountableConnection, error) {
	saleFilterOpts, appErr := args.parse()
	if appErr != nil {
		return nil, appErr
	}

	embedCtx := GetContextValue[*web.Context](ctx, WebCtx)
	totalSales, sales, appErr := embedCtx.App.Srv().DiscountService().FilterSalesByOption(saleFilterOpts)
	if appErr != nil {
		return nil, appErr
	}

	saleKeyFunc := saleSortFieldsMap[SaleSortFieldName].keyFunc
	if args.SortBy != nil && args.SortBy.Field != SaleSortFieldName {
		saleKeyFunc = saleSortFieldsMap[args.SortBy.Field].keyFunc
	}

	res := constructCountableConnection(sales, totalSales, args.GraphqlParams, saleKeyFunc, systemSaleToGraphqlSale)

	return (*SaleCountableConnection)(unsafe.Pointer(res)), nil
}

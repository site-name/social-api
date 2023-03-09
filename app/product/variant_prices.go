package product

import (
	"net/http"
	"strings"
	"sync"

	"github.com/Masterminds/squirrel"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
)

// getVariantPricesInChannelsDict
func (a *ServiceProduct) getVariantPricesInChannelsDict(product model.Product) (map[string][]*goprices.Money, *model.AppError) {
	variantChannelListings, appErr := a.ProductVariantChannelListingsByOption(nil, &model.ProductVariantChannelListingFilterOption{
		VariantProductID: squirrel.Eq{store.ProductVariantTableName + ".ProductID": product.Id},
		PriceAmount:      squirrel.NotEq{store.ProductVariantChannelListingTableName + ".PriceAmount": nil},
	})
	if appErr != nil {
		return nil, appErr
	}

	pricesDict := map[string][]*goprices.Money{}
	for _, listing := range variantChannelListings {
		listing.PopulateNonDbFields() // must run this first
		pricesDict[listing.ChannelID] = append(pricesDict[listing.ChannelID], listing.Price)
	}

	return pricesDict, nil
}

func (a *ServiceProduct) getProductDiscountedPrice(
	variantPrices []*goprices.Money,
	product model.Product,
	collections []*model.Collection,
	discounts []*model.DiscountInfo,
	chanNel model.Channel,

) (*goprices.Money, *model.AppError) {

	// validate variantPrices have same currencies
	var (
		standardCurrency        string
		discountedVariantPrices []*goprices.Money
	)

	for i, item := range variantPrices {
		if i == 0 {
			standardCurrency = item.Currency
		} else if !strings.EqualFold(standardCurrency, item.Currency) {
			return nil, model.NewAppError("getProductDiscountedPrice", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "variantPrices"}, "", http.StatusBadRequest)
		}

		discoutnedvariantPrice, appErr := a.srv.DiscountService().CalculateDiscountedPrice(
			product,
			item,
			collections,
			discounts,
			chanNel,
			"",
		)
		if appErr != nil {
			return nil, appErr
		}

		discountedVariantPrices = append(discountedVariantPrices, discoutnedvariantPrice)
	}

	min, _ := util.MinMaxMoneyInMoneySlice(discountedVariantPrices)
	return min, nil
}

// UpdateProductDiscountedPrice
//
// NOTE: `discounts` can be nil
func (a *ServiceProduct) UpdateProductDiscountedPrice(product model.Product, discounts []*model.DiscountInfo) *model.AppError {

	var functionAppError *model.AppError
	if len(discounts) == 0 {
		discounts, functionAppError = a.srv.DiscountService().FetchActiveDiscounts()
		if functionAppError != nil {
			return functionAppError
		}
	}

	var (
		collectionsContainProduct   []*model.Collection
		variantPricesInChannelsDict map[string][]*goprices.Money
		productChannelListings      []*model.ProductChannelListing
		wg                          sync.WaitGroup
		mut                         sync.Mutex
	)

	syncSetAppErr := func(err *model.AppError) {
		mut.Lock()
		defer mut.Unlock()
		if err != nil && functionAppError == nil {
			functionAppError = err
		}
	}

	wg.Add(3)

	go func() {
		mut.Lock()
		defer wg.Done()
		defer mut.Unlock()

		res, appErr := a.CollectionsByProductID(product.Id)
		if appErr != nil {
			syncSetAppErr(appErr)
		} else {
			collectionsContainProduct = res
		}
	}()

	go func() {
		mut.Lock()
		defer wg.Done()
		defer mut.Unlock()

		res, appErr := a.getVariantPricesInChannelsDict(product)
		if appErr != nil {
			syncSetAppErr(appErr)
		} else {
			variantPricesInChannelsDict = res
		}
	}()

	go func() {
		mut.Lock()
		defer wg.Done()
		defer mut.Unlock()

		res, appErr := a.ProductChannelListingsByOption(&model.ProductChannelListingFilterOption{
			ProductID:       squirrel.Eq{store.ProductChannelListingTableName + ".ProductID": product.Id},
			PrefetchChannel: true, // this will populate `Channel` fields of every product channel listings
		})
		if appErr != nil {
			syncSetAppErr(appErr)
		} else {
			productChannelListings = res
		}
	}()

	wg.Wait()

	// check appError:
	if functionAppError != nil {
		return functionAppError
	}

	var productChannelListingsToUpdate []*model.ProductChannelListing

	for _, listing := range productChannelListings {
		listing.PopulateNonDbFields() // this call is crutial

		variantPrices := variantPricesInChannelsDict[listing.ChannelID]
		if len(variantPrices) == 0 {
			continue
		}

		if listing.GetChannel() != nil { // check if there is a channel populated
			productDiscountedPrice, appErr := a.getProductDiscountedPrice(
				variantPrices,
				product,
				collectionsContainProduct,
				discounts,
				*listing.GetChannel(),
			)
			if appErr != nil {
				return appErr
			}

			if listing.DiscountedPrice != nil &&
				productDiscountedPrice != nil &&
				// notice below: NOT equal
				!listing.DiscountedPrice.Amount.Equal(productDiscountedPrice.Amount) {

				listing.DiscountedPriceAmount = &productDiscountedPrice.Amount
				productChannelListingsToUpdate = append(productChannelListingsToUpdate, listing)
			}
		}
	}

	if len(productChannelListingsToUpdate) > 0 {
		_, functionAppError = a.BulkUpsertProductChannelListings(productChannelListingsToUpdate)
	}

	return functionAppError
}

// UpdateProductsDiscountedPrices
func (a *ServiceProduct) UpdateProductsDiscountedPrices(products []*model.Product, discounts []*model.DiscountInfo) *model.AppError {
	var (
		appError *model.AppError
		wg       sync.WaitGroup
		mut      sync.Mutex
	)
	if len(discounts) == 0 {
		discounts, appError = a.srv.DiscountService().FetchActiveDiscounts()
		if appError != nil {
			return appError
		}
	}

	wg.Add(len(products))

	syncSetAppError := func(err *model.AppError) {
		mut.Lock()
		defer mut.Unlock()

		if err != nil && appError == nil {
			appError = err
		}
	}

	for _, product := range products {
		go func(prd *model.Product) {
			defer wg.Done()
			syncSetAppError(a.UpdateProductDiscountedPrice(*prd, discounts))
		}(product)
	}

	wg.Wait()

	return appError
}

func (a *ServiceProduct) UpdateProductsDiscountedPricesOfCatalogues(productIDs, categoryIDs, collectionIDs, variantIDs []string) *model.AppError {
	products, err := a.srv.Store.Product().SelectForUpdateDiscountedPricesOfCatalogues(productIDs, categoryIDs, collectionIDs, variantIDs)
	if err != nil {
		return model.NewAppError("UpdateProductsDiscountedPricesOfCatalogues", "app.product.error_finding_products_by_given_id_lists.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return a.UpdateProductsDiscountedPrices(products, nil)
}

// UpdateProductsDiscountedPricesOfDiscount
//
// NOTE: discount must be either *Sale or *Voucher
func (a *ServiceProduct) UpdateProductsDiscountedPricesOfDiscount(discount interface{}) *model.AppError {
	var (
		productFilterOption    model.ProductFilterOption
		categoryFilterOption   model.CategoryFilterOption
		collectionFilterOption model.CollectionFilterOption
		variantFilterOptions   model.ProductVariantFilterOption
	)
	switch t := discount.(type) {
	case *model.Sale:
		productFilterOption.SaleID = squirrel.Eq{store.SaleProductRelationTableName + ".SaleID": t.Id}
		categoryFilterOption.SaleID = squirrel.Eq{store.SaleCategoryRelationTableName + ".SaleID": t.Id}
		collectionFilterOption.SaleID = squirrel.Eq{store.SaleCollectionRelationTableName + ".SaleID": t.Id}
		variantFilterOptions.SaleID = squirrel.Eq{store.SaleProductVariantTableName + ".SaleID": t.Id}
	case *model.Voucher:
		productFilterOption.VoucherID = squirrel.Eq{store.VoucherProductTableName + ".VoucherID": t.Id}
		categoryFilterOption.VoucherID = squirrel.Eq{store.VoucherCategoryTableName + ".VoucherID": t.Id}
		collectionFilterOption.VoucherID = squirrel.Eq{store.VoucherCollectionTableName + ".VoucherID": t.Id}
		variantFilterOptions.VoucherID = squirrel.Eq{store.VoucherProductVariantTableName + ".VoucherID": t.Id}

	default:
		return model.NewAppError("UpdateProductsDiscountedPricesOfDiscount", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "discount"}, "", http.StatusBadRequest)
	}

	var (
		productIDs    []string
		categoryIDs   []string
		collectionIDs []string
		variantIDs    []string
		appError      *model.AppError
		wg            sync.WaitGroup
		mut           sync.Mutex
	)

	syncSetAppError := func(err *model.AppError) {
		mut.Lock()
		defer mut.Unlock()
		if err != nil && appError == nil {
			appError = err
		}
	}

	wg.Add(4)

	go func() {
		defer wg.Done()
		products, appErr := a.ProductsByOption(&productFilterOption)
		syncSetAppError(appErr)
		productIDs = products.IDs()
	}()

	go func() {
		defer wg.Done()
		categories, appErr := a.CategoriesByOption(&categoryFilterOption)
		syncSetAppError(appErr)
		categoryIDs = categories.IDs(false)
	}()

	go func() {
		defer wg.Done()
		collections, appErr := a.CollectionsByOption(&collectionFilterOption)
		syncSetAppError(appErr)
		collectionIDs = collections.IDs()
	}()

	go func() {
		defer wg.Done()
		variants, appErr := a.ProductVariantsByOption(&variantFilterOptions)
		syncSetAppError(appErr)
		variantIDs = variants.IDs()
	}()

	wg.Wait()

	return a.UpdateProductsDiscountedPricesOfCatalogues(productIDs, categoryIDs, collectionIDs, variantIDs)
}

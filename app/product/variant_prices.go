package product

import (
	"net/http"
	"strings"
	"sync/atomic"

	"github.com/Masterminds/squirrel"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

// getVariantPricesInChannelsDict
func (a *ServiceProduct) getVariantPricesInChannelsDict(product model.Product) (map[string][]*goprices.Money, *model_helper.AppError) {
	variantChannelListings, appErr := a.ProductVariantChannelListingsByOption(&model.ProductVariantChannelListingFilterOption{
		VariantProductID: squirrel.Eq{model.ProductVariantTableName + ".ProductID": product.Id},
		Conditions:       squirrel.NotEq{model.ProductVariantChannelListingTableName + ".PriceAmount": nil},
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

) (*goprices.Money, *model_helper.AppError) {
	// validate variantPrices have same currencies
	var (
		standardCurrency        string
		discountedVariantPrices []*goprices.Money
	)

	for i, item := range variantPrices {
		if i == 0 {
			standardCurrency = item.Currency
		} else if !strings.EqualFold(standardCurrency, item.Currency) {
			return nil, model_helper.NewAppError("getProductDiscountedPrice", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "variantPrices"}, "", http.StatusBadRequest)
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
func (a *ServiceProduct) UpdateProductDiscountedPrice(transaction *gorm.DB, product model.Product, discounts []*model.DiscountInfo) *model_helper.AppError {
	if len(discounts) == 0 {
		var appErr *model_helper.AppError
		discounts, appErr = a.srv.DiscountService().FetchActiveDiscounts()
		if appErr != nil {
			return appErr
		}
	}

	var (
		collectionsContainProduct   model.Collections
		variantPricesInChannelsDict map[string][]*goprices.Money
		productChannelListings      model.ProductChannelListings
		atomicValue                 atomic.Int32
		appError                    = make(chan *model_helper.AppError)
	)
	defer close(appError)
	atomicValue.Add(3)

	go func() {
		defer atomicValue.Add(-1)

		collections, appErr := a.CollectionsByProductID(product.Id)
		if appErr != nil {
			appError <- appErr
			return
		}
		collectionsContainProduct = collections
	}()

	go func() {
		defer atomicValue.Add(-1)

		res, appErr := a.getVariantPricesInChannelsDict(product)
		if appErr != nil {
			appError <- appErr
			return
		}
		variantPricesInChannelsDict = res
	}()

	go func() {
		defer atomicValue.Add(-1)

		listings, appErr := a.ProductChannelListingsByOption(&model.ProductChannelListingFilterOption{
			Conditions: squirrel.Eq{model.ProductChannelListingTableName + ".ProductID": product.Id},
			Preloads:   []string{"Channel"}, // this will populate `Channel` fields of every product channel listings
		})
		if appErr != nil {
			appError <- appErr
			return
		}
		productChannelListings = listings
	}()

	for atomicValue.Load() != 0 {
		select {
		case err := <-appError:
			return err
		default:
		}
	}

	var productChannelListingsToUpdate []*model.ProductChannelListing

	for _, listing := range productChannelListings {
		listing.PopulateNonDbFields() // this call is needed

		variantPrices := variantPricesInChannelsDict[listing.ChannelID]
		if len(variantPrices) == 0 {
			continue
		}

		if listing.Channel != nil { // check if there is a channel populated
			productDiscountedPrice, appErr := a.getProductDiscountedPrice(
				variantPrices,
				product,
				collectionsContainProduct,
				discounts,
				*listing.Channel,
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
		_, appErr := a.BulkUpsertProductChannelListings(transaction, productChannelListingsToUpdate)
		if appErr != nil {
			return appErr
		}
	}

	return nil
}

// UpdateProductsDiscountedPrices
func (a *ServiceProduct) UpdateProductsDiscountedPrices(transaction *gorm.DB, products []*model.Product, discounts []*model.DiscountInfo) *model_helper.AppError {
	if len(discounts) == 0 {
		var appErr *model_helper.AppError
		discounts, appErr = a.srv.DiscountService().FetchActiveDiscounts()
		if appErr != nil {
			return appErr
		}
	}

	var (
		appError    = make(chan *model_helper.AppError)
		atomicValue atomic.Int32
	)
	defer close(appError)

	atomicValue.Add(int32(len(products)))

	for _, product := range products {
		go func(product *model.Product) {
			defer atomicValue.Add(-1)

			appErr := a.UpdateProductDiscountedPrice(transaction, *product, discounts)
			if appErr != nil {
				appError <- appErr
			}
		}(product)
	}

	for atomicValue.Load() != 0 {
		select {
		case err := <-appError:
			return err
		default:
		}
	}

	return nil
}

func (a *ServiceProduct) UpdateProductsDiscountedPricesOfCatalogues(transaction *gorm.DB, productIDs, categoryIDs, collectionIDs, variantIDs []string) *model_helper.AppError {
	products, err := a.srv.Store.Product().SelectForUpdateDiscountedPricesOfCatalogues(transaction, productIDs, categoryIDs, collectionIDs, variantIDs)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrInvalidInput); ok {
			statusCode = http.StatusBadRequest
		}
		return model_helper.NewAppError("UpdateProductsDiscountedPricesOfCatalogues", "app.product.error_finding_products_by_given_id_lists.app_error", nil, err.Error(), statusCode)
	}

	return a.UpdateProductsDiscountedPrices(transaction, products, nil)
}

// UpdateProductsDiscountedPricesOfDiscount
//
// NOTE: discount must be either *Sale or *Voucher
func (a *ServiceProduct) UpdateProductsDiscountedPricesOfDiscount(transaction *gorm.DB, discount interface{}) *model_helper.AppError {
	var (
		productFilterOption    model.ProductFilterOption
		categoryFilterOption   model.CategoryFilterOption
		collectionFilterOption model.CollectionFilterOption
		variantFilterOptions   model.ProductVariantFilterOption
	)
	switch t := discount.(type) {
	case *model.Sale:
		productFilterOption.SaleID = squirrel.Eq{model.SaleProductTableName + ".sale_id": t.Id}
		categoryFilterOption.SaleID = squirrel.Eq{model.SaleCategoryTableName + ".sale_id": t.Id}
		collectionFilterOption.SaleID = squirrel.Eq{model.SaleCollectionTableName + ".sale_id": t.Id}
		variantFilterOptions.SaleID = squirrel.Eq{model.SaleProductVariantTableName + ".sale_id": t.Id}
	case *model.Voucher:
		productFilterOption.VoucherID = squirrel.Eq{model.VoucherProductTableName + ".voucher_id": t.Id}
		categoryFilterOption.VoucherID = squirrel.Eq{model.VoucherCategoryTableName + ".voucher_id": t.Id}
		collectionFilterOption.VoucherID = squirrel.Eq{model.VoucherCollectionTableName + ".voucher_id": t.Id}
		variantFilterOptions.VoucherID = squirrel.Eq{model.VoucherProductVariantTableName + ".voucher_id": t.Id}

	default:
		return model_helper.NewAppError("UpdateProductsDiscountedPricesOfDiscount", model.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "discount"}, "", http.StatusBadRequest)
	}

	var (
		atomicInt32   atomic.Int32
		appError      = make(chan *model_helper.AppError)
		value         = make(chan any)
		productIDs    []string
		categoryIDs   []string
		collectionIDs []string
		variantIDs    []string
	)
	defer func() {
		close(appError)
		close(value)
	}()

	atomicInt32.Add(4) // specify there are 4 goroutines to run

	go func() {
		products, appErr := a.ProductsByOption(&productFilterOption)
		if appErr != nil {
			appError <- appErr
			return
		}
		value <- products
	}()

	go func() {
		categories, appErr := a.CategoriesByOption(&categoryFilterOption)
		if appErr != nil {
			appError <- appErr
			return
		}
		value <- categories
	}()

	go func() {
		_, collections, appErr := a.CollectionsByOption(&collectionFilterOption)
		if appErr != nil {
			appError <- appErr
			return
		}
		value <- collections
	}()

	go func() {
		variants, appErr := a.ProductVariantsByOption(&variantFilterOptions)
		if appErr != nil {
			appError <- appErr
			return
		}
		value <- variants
	}()

	for atomicInt32.Load() != 0 {
		select {
		case err := <-appError:
			return err
		case val := <-value:
			atomicInt32.Add(-1)

			switch t := val.(type) {
			case model.Products:
				productIDs = t.IDs()
			case model.Categories:
				categoryIDs = t.IDs(false)
			case model.ProductVariants:
				variantIDs = t.IDs()
			case model.Collections:
				collectionIDs = t.IDs()
			}
		default:
		}
	}

	return a.UpdateProductsDiscountedPricesOfCatalogues(transaction, productIDs, categoryIDs, collectionIDs, variantIDs)
}

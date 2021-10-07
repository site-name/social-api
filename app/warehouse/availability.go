package warehouse

import (
	"net/http"

	"github.com/sitename/sitename/exception"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/checkout"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/model/warehouse"
	"github.com/sitename/sitename/modules/util"
)

// getAvailableQuantity get all stocks quantity (both in stocks and their allocations) not exported
func (a *ServiceWarehouse) getAvailableQuantity(stocks warehouse.Stocks) (int, *model.AppError) {
	if len(stocks) == 0 {
		return 0, nil
	}

	// reference: https://github.com/mirumee/saleor/blob/master/saleor/warehouse/availability.py, function (_get_available_quantity)
	// not sure yet why using SUM(DISTINCT 'quantity') on stocks
	var totalQuantity int
	meetMap := make(map[int]bool)
	stockIDs := stocks.IDs()

	for _, stock := range stocks {
		if _, met := meetMap[stock.Quantity]; !met {
			totalQuantity += stock.Quantity
			meetMap[stock.Quantity] = true // this replicates SUM(DISTINCT ...)
		}
	}

	allocations, appErr := a.AllocationsByOption(nil, &warehouse.AllocationFilterOption{
		StockID: &model.StringFilter{
			StringOption: &model.StringOption{
				In: stockIDs,
			},
		},
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return 0, appErr
		}
		allocations = []*warehouse.Allocation{}
	}

	var quantityAllocated int
	for _, allocation := range allocations {
		quantityAllocated += allocation.QuantityAllocated
	}

	return util.Max(totalQuantity-quantityAllocated, 0), nil
}

// CheckStockAndPreorderQuantity Validate if there is stock/preorder available for given variant.
// :raises InsufficientStock: when there is not enough items in stock for a variant
// or there is not enough available preorder items for a variant.
func (s *ServiceWarehouse) CheckStockAndPreorderQuantity(variant *product_and_discount.ProductVariant, countryCode string, channelSlug string, quantity int) {
	panic("not implemented")
}

// Validate if there is stock available for given variant in given country.
//
// If so - returns None. If there is less stock then required raise InsufficientStock
// exception.
func (a *ServiceWarehouse) CheckStockQuantity(variant *product_and_discount.ProductVariant, countryCode string, channelSlug string, quantity int) (*exception.InsufficientStock, *model.AppError) {
	if *variant.TrackInventory {
		stocks, appErr := a.GetVariantStocksForCountry(nil, countryCode, channelSlug, variant.Id)
		if appErr != nil {
			if appErr.StatusCode == http.StatusInternalServerError {
				return nil, appErr
			}
			stocks = []*warehouse.Stock{} // in case stocks is nil
		}

		if len(stocks) == 0 {
			return &exception.InsufficientStock{
				Items: []*exception.InsufficientStockData{
					{Variant: *variant},
				},
			}, nil
		}

		availableQuantity, appErr := a.getAvailableQuantity(stocks)
		if appErr != nil {
			return nil, appErr
		}
		if quantity > availableQuantity {
			return &exception.InsufficientStock{
				Items: []*exception.InsufficientStockData{
					{Variant: *variant},
				},
			}, nil
		}
	}

	return nil, nil
}

// Validate if there is stock available for given variants in given country.
//
// :raises InsufficientStock: when there is not enough items in stock for a variant
func (a *ServiceWarehouse) CheckStockQuantityBulk(
	variants product_and_discount.ProductVariants,
	countryCode string,
	quantities []int,
	channelSlug string,
	additionalFilterLookup model.StringInterface, // can be nil, if non-nil then it must be map[string]interface{}{"warehouse_id": <an UUID string>}
	existingLines []*checkout.CheckoutLineInfo, // can be nil

) (*exception.InsufficientStock, *model.AppError) {

	variants = variants.FilterNils()

	// build a filter option
	allVariantStockFilterOption := &warehouse.StockFilterForCountryAndChannel{
		CountryCode: countryCode,
		ChannelSlug: channelSlug,

		ProductVariantIDFilter: &model.StringFilter{
			StringOption: &model.StringOption{
				In: variants.IDs(),
			},
		},
		AnnotateAvailabeQuantity: true,
	}

	// check if `additionalFilterLookup` is not nil:
	if additionalFilterLookup != nil && additionalFilterLookup["warehouse_id"] != nil {
		allVariantStockFilterOption.WarehouseID = additionalFilterLookup["warehouse_id"].(string)
	}

	allVariantStocks, appErr := a.FilterStocksForCountryAndChannel(nil, allVariantStockFilterOption)
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr
		}
		// ignore not found error
	}

	// variantStocks has keys of product variant ids
	var variantStocks = map[string][]*warehouse.Stock{}
	for _, stock := range allVariantStocks {
		if stock != nil {
			variantStocks[stock.ProductVariantID] = append(variantStocks[stock.ProductVariantID], stock)
		}
	}

	var (
		// keys are product variant ids, values are order line quantities
		variantsQuantities = map[string]int{}
		insufficientStocks = []*exception.InsufficientStockData{}
	)

	for _, lineInfo := range existingLines {
		if lineInfo != nil && model.IsValidId(lineInfo.Variant.Id) {
			variantsQuantities[lineInfo.Variant.Id] = lineInfo.Line.Quantity
		}
	}

	for i := 0; i < util.Min(len(variants), len(quantities)); i++ {
		quantity := quantities[i]
		variant := variants[i]

		quantity += variantsQuantities[variant.Id]

		stocks, exists := variantStocks[variant.Id]
		if !exists || stocks == nil {
			stocks = []*warehouse.Stock{}
		}

		var availableQuantity int
		for _, stock := range stocks {
			availableQuantity += stock.AvailableQuantity
		}

		if len(stocks) == 0 {
			insufficientStocks = append(insufficientStocks, &exception.InsufficientStockData{
				Variant:           *variant,
				AvailableQuantity: &availableQuantity,
			})
		} else if *variant.TrackInventory {
			if quantities[i] > availableQuantity {
				insufficientStocks = append(insufficientStocks, &exception.InsufficientStockData{
					Variant:           *variant,
					AvailableQuantity: &availableQuantity,
				})
			}
		}
	}

	if len(insufficientStocks) > 0 {
		return &exception.InsufficientStock{
			Items: insufficientStocks,
		}, nil
	}

	return nil, nil
}

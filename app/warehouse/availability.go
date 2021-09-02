package warehouse

import (
	"net/http"

	"github.com/sitename/sitename/model"
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

// Validate if there is stock available for given variant in given country.
//
// If so - returns None. If there is less stock then required raise InsufficientStock
// exception.
func (a *ServiceWarehouse) CheckStockQuantity(variant *product_and_discount.ProductVariant, countryCode string, channelSlug string, quantity int) (*warehouse.InsufficientStock, *model.AppError) {
	if *variant.TrackInventory {
		stocks, appErr := a.GetVariantStocksForCountry(nil, countryCode, channelSlug, variant.Id)
		if appErr != nil {
			if appErr.StatusCode == http.StatusInternalServerError {
				return nil, appErr
			}
			stocks = []*warehouse.Stock{} // in case stocks is nil
		}

		if len(stocks) == 0 {
			return &warehouse.InsufficientStock{
				Items: []*warehouse.InsufficientStockData{
					{Variant: *variant},
				},
			}, nil
		}

		availableQuantity, appErr := a.getAvailableQuantity(stocks)
		if appErr != nil {
			return nil, appErr
		}
		if quantity > availableQuantity {
			return &warehouse.InsufficientStock{
				Items: []*warehouse.InsufficientStockData{
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
func (a *ServiceWarehouse) CheckStockQuantityBulk(variants product_and_discount.ProductVariants, countryCode string, quantities []int, channelSlug string) (*warehouse.InsufficientStock, *model.AppError) {
	allVariantStocks, appErr := a.FilterStocksForCountryAndChannel(nil, &warehouse.StockFilterForCountryAndChannel{
		CountryCode: countryCode,
		ChannelSlug: channelSlug,

		ProductVariantIDFilter: &model.StringFilter{
			StringOption: &model.StringOption{
				In: variants.IDs(),
			},
		},
		AnnotateAvailabeQuantity: true,
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr
		}
		allVariantStocks = []*warehouse.Stock{} // just in case allVariantStocks is nil
	}

	variantStocks := map[string][]*warehouse.Stock{}
	for _, stock := range allVariantStocks {
		variantStocks[stock.ProductVariantID] = append(variantStocks[stock.ProductVariantID], stock)
	}

	insufficientStocks := []*warehouse.InsufficientStockData{}
	for i := 0; i < util.Min(len(variants), len(quantities)); i++ {
		stocks, exists := variantStocks[variants[i].Id]
		if !exists || stocks == nil {
			stocks = []*warehouse.Stock{}
		}

		var availableQuantity int
		for _, stock := range stocks {
			availableQuantity += stock.AvailableQuantity
		}

		if len(stocks) == 0 {
			insufficientStocks = append(insufficientStocks, &warehouse.InsufficientStockData{
				Variant:           *variants[i],
				AvailableQuantity: &availableQuantity,
			})
		} else if *variants[i].TrackInventory {
			if quantities[i] > availableQuantity {
				insufficientStocks = append(insufficientStocks, &warehouse.InsufficientStockData{
					Variant:           *variants[i],
					AvailableQuantity: &availableQuantity,
				})
			}
		}
	}

	if len(insufficientStocks) > 0 {
		return &warehouse.InsufficientStock{
			Items: insufficientStocks,
		}, nil
	}

	return nil, nil
}

// Check if there is any variant of given product available in given country
func (a *ServiceWarehouse) IsProductInStock(productID string, countryCode string, channelSlug string) (bool, *model.AppError) {
	stocks, appErr := a.GetProductStocksForCountryAndChannel(nil, &warehouse.StockFilterForCountryAndChannel{
		CountryCode:              countryCode,
		ChannelSlug:              channelSlug,
		ProductID:                productID,
		AnnotateAvailabeQuantity: true,
	})
	if appErr != nil {
		return false, appErr
	}

	// check at least 1 item in slice > 0
	for _, stocks := range stocks {
		if stocks.AvailableQuantity > 0 {
			return true, nil
		}
	}

	return false, nil
}

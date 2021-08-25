package warehouse

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/model/warehouse"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
)

// getAvailableQuantity get all stocks quantity (both in stocks and their allocations) not exported
func (a *AppWarehouse) getAvailableQuantity(stocks []*warehouse.Stock) (int, *model.AppError) {
	if len(stocks) == 0 {
		return 0, nil
	}

	// reference: https://github.com/mirumee/saleor/blob/master/saleor/warehouse/availability.py, function (_get_available_quantity)
	// not sure yet why using SUM(DISTINCT 'quantity') on stocks
	var totalQuantity int
	meetMap := make(map[int]bool)
	stockIDs := make([]string, len(stocks)) // get all stock ids from `stocks`

	for i, stock := range stocks {
		stockIDs[i] = stock.Id
		if _, ok := meetMap[stock.Quantity]; !ok {
			totalQuantity += stock.Quantity
			meetMap[stock.Quantity] = true
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
		return 0, appErr
	}

	var allocatedQuantity int
	for _, allocation := range allocations {
		allocatedQuantity += allocation.QuantityAllocated
	}

	if sub := totalQuantity - allocatedQuantity; sub > 0 {
		return sub, nil
	}

	return 0, nil
}

// Validate if there is stock available for given variant in given country.
//
// If so - returns None. If there is less stock then required raise InsufficientStock
// exception.
func (a *AppWarehouse) CheckStockQuantity(variant *product_and_discount.ProductVariant, countryCode string, channelSlug string, quantity int) (*warehouse.InsufficientStock, *model.AppError) {
	if *variant.TrackInventory {
		stocks, appErr := a.GetVariantStocksForCountry(countryCode, channelSlug, variant.Id, quantity)
		if appErr != nil {
			return nil, appErr
		}

		availableQuantity, appErr := a.getAvailableQuantity(stocks)
		if appErr != nil {
			// error server lookup
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
func (a *AppWarehouse) CheckStockQuantityBulk(variants []*product_and_discount.ProductVariant, countryCode string, quantities []int, channelSlug string) (*warehouse.InsufficientStock, *model.AppError) {
	stocks, appErr := a.StocksByOption(nil, &warehouse.StockFilterOption{
		ForCountryAndChannel: &warehouse.StockFilterForCountryAndChannel{
			CountryCode: countryCode,
			ChannelSlug: channelSlug,
		},
		ProductVariantID: &model.StringFilter{
			StringOption: &model.StringOption{
				In: product_and_discount.ProductVariants(variants).IDs(),
			},
		},
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return nil, appErr
		}
		stocks = []*warehouse.Stock{} // just in case stocks is nil
	}

	allVariantStocks := []*warehouse.Stock{}
	for _, stock := range stocks {
		for _, variant := range variants {
			if stock.ProductVariantID == variant.Id {
				allVariantStocks = append(allVariantStocks, stock)
			}
		}
	}

	variantStocks := map[string][]*warehouse.Stock{}
	for _, stock := range allVariantStocks {
		if _, ok := variantStocks[stock.ProductVariantID]; !ok {
			variantStocks[stock.ProductVariantID] = []*warehouse.Stock{}
		}
		variantStocks[stock.ProductVariantID] = append(variantStocks[stock.ProductVariantID], stock)
	}

	insufficientStocks := []*warehouse.InsufficientStockData{}
	for i := 0; i < util.Min(len(variants), len(quantities)); i++ {
		stocks_, ok := variantStocks[variants[i].Id]

		availableQuantity, appErr := a.getAvailableQuantity(stocks_)
		if appErr != nil {
			return nil, appErr
		}

		if !ok {
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
func (a *AppWarehouse) IsProductInStock(productID string, countryCode string, channelSlug string) (bool, *model.AppError) {
	stocks, err := a.Srv().Store.Stock().FilterProductStocksForCountryAndChannel(&warehouse.StockFilterForCountryAndChannel{
		CountryCode: countryCode,
		ChannelSlug: channelSlug,
		ProductID:   productID,
	})

	if err != nil {
		return false, store.AppErrorFromDatabaseLookupError("IsProductInStock", "app.warehouse.product_stocks_for_country_and_channel_missing.app_error", err)
	}

	availableQuantity, appErr := a.getAvailableQuantity(stocks)
	if appErr != nil {
		return false, appErr
	}

	return availableQuantity > 0, nil
}

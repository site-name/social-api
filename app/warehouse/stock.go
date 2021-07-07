package warehouse

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/model/warehouse"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
)

func (a *AppWarehouse) CheckStockQuantity(variant *product_and_discount.ProductVariant, countryCode string, channelSlug string, quantity uint) (*warehouse.InsufficientStock, *model.AppError) {
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
					{Variant: variant.Id},
				},
			}, nil
		}
	}

	return nil, nil
}

// GetVariantStocksForCountry validates if stock for given country are valid.
// Not exported.
func (a *AppWarehouse) GetVariantStocksForCountry(countryCode string, channelSlug string, variantID string, quantity uint) ([]*warehouse.Stock, *model.AppError) {
	stocks, _, _, err := a.Srv().Store.Stock().FilterVariantStocksForCountry(
		&warehouse.ForCountryAndChannelFilter{
			CountryCode: countryCode,
			ChannelSlug: channelSlug,
		},
		variantID,
	)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("GetVariantStocksForCountry", "app.warehouse.stock_filter_forcountry_missing.app_error", err)
	}

	return stocks, nil
}

// getAvailableQuantity get all stocks quantity (both in stocks and their allocations)
// not exported
func (a *AppWarehouse) getAvailableQuantity(stocks []*warehouse.Stock) (uint, *model.AppError) {
	if len(stocks) == 0 {
		return 0, nil
	}

	// try looking up all allocations of given stocks:
	stockIDs := make([]string, len(stocks))
	for i := range stocks {
		stockIDs[i] = stocks[i].Id
	}

	// reference: https://github.com/mirumee/saleor/blob/master/saleor/warehouse/availability.py
	var totalQuantity uint
	for _, stock := range stocks {
		// TODO: check whether to use SUM DISTINCT here
		totalQuantity += stock.Quantity
	}

	allocations, err := a.Srv().Store.Allocation().AllocationsByParentIDs(stockIDs, warehouse.ByStock)
	if err != nil {
		// 2 types of errors could happend here: not found OR server error
		appErr := store.AppErrorFromDatabaseLookupError("getAvailableQuantity", "app.warehouse.allocations_by_stocks_missing.app_error", err)
		if appErr.StatusCode == http.StatusNotFound {
			return totalQuantity, nil
		}
		return 0, appErr
	}

	var allocatedQuantity uint
	for _, allocation := range allocations {
		allocatedQuantity += allocation.QuantityAllocated
	}

	if totalQuantity-allocatedQuantity > 0 {
		return totalQuantity - allocatedQuantity, nil
	}

	return 0, nil
}

func (a *AppWarehouse) CheckStockQuantityBulk(variants []*product_and_discount.ProductVariant, countryCode string, quantities []uint, channelSlug string) (*warehouse.InsufficientStock, *model.AppError) {
	stocks, _, _, err := a.Srv().Store.Stock().FilterForCountryAndChannel(&warehouse.ForCountryAndChannelFilter{
		CountryCode: countryCode,
		ChannelSlug: channelSlug,
	})

	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("CheckStockQuantityBulk", "app.warehouse.stocks_filter_for_country_and_channel.app_error", err)
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
				Variant:           variants[i].Id,
				AvailableQuantity: availableQuantity,
			})
		} else if *variants[i].TrackInventory {
			if quantities[i] > availableQuantity {
				insufficientStocks = append(insufficientStocks, &warehouse.InsufficientStockData{
					Variant:           variants[i].Id,
					AvailableQuantity: availableQuantity,
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

func (a *AppWarehouse) IsProductInStock(productID string, countryCode string, channelSlug string) (bool, *model.AppError) {
	stocks, _, _, err := a.Srv().Store.Stock().FilterProductStocksForCountryAndChannel(&warehouse.ForCountryAndChannelFilter{
		CountryCode: countryCode,
		ChannelSlug: channelSlug,
	}, productID)

	if err != nil {
		return false, store.AppErrorFromDatabaseLookupError("IsProductInStock", "app.warehouse.product_stocks_for_country_and_channel_missing.app_error", err)
	}

	availableQuantity, appErr := a.getAvailableQuantity(stocks)
	if appErr != nil {
		return false, appErr
	}

	return availableQuantity > 0, nil
}

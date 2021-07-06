package warehouse

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/model/warehouse"
	"github.com/sitename/sitename/store"
)

func (a *AppWarehouse) CheckStockQuantity(variant *product_and_discount.ProductVariant, countryCode string, channelSlug string, quantity uint) *model.AppError {
	if *variant.TrackInventory {
		appErr := a.GetVariantStocksForCountry(countryCode, channelSlug, variant.Id, quantity)
		return appErr
	}

	return nil
}

// GetVariantStocksForCountry validates if stock for given country are valid.
// Not exported.
func (a *AppWarehouse) GetVariantStocksForCountry(countryCode string, channelSlug string, variantID string, quantity uint) *model.AppError {
	stocks, _, _, err := a.Srv().Store.Stock().FilterVariantStocksForCountry(
		&warehouse.ForCountryAndChannelFilter{
			CountryCode: countryCode,
			ChannelSlug: channelSlug,
		},
		variantID,
	)
	if err != nil {
		return store.AppErrorFromDatabaseLookupError("GetVariantStocksForCountry", "app.warehouse.stock_filter_forcountry_missing.app_error", err)
	}

	commonReturnError := model.NewAppError("GetVariantStocksForCountry", "app.warehouse.stock_insufficient.app_error", map[string]interface{}{"VariantID": variantID}, "", http.StatusInsufficientStorage)
	// if stock is insufficient for country:
	if len(stocks) == 0 {
		return commonReturnError
	}

	availableQuantity, appERr := a.getAvailableQuantity(stocks)
	if appERr != nil {
		// error server lookup
		return appERr
	}
	if quantity > availableQuantity {
		return commonReturnError
	}

	return nil
}

// getAvailableQuantity get all stocks quantity (both in stocks and their allocations)
// not exported
func (a *AppWarehouse) getAvailableQuantity(stocks []*warehouse.Stock) (uint, *model.AppError) {
	// try looking up all allocations of given stocks:
	stockIDs := make([]string, len(stocks))
	for i := range stocks {
		stockIDs[i] = stocks[i].Id
	}

	var totalQuantity uint
	for _, stock := range stocks {
		totalQuantity += stock.Quantity
	}

	allocations, err := a.Srv().Store.Allocation().AllocationsByParentIDs(stockIDs, store.ByStock)
	if err != nil {
		// 2 types of errors could happend here:
		// not found || server error
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

package warehouse

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/warehouse"
	"github.com/sitename/sitename/store"
)

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

// GetStockByOption takes options for filtering 1 stock
func (a *AppWarehouse) GetStockByOption(option *warehouse.StockFilterOption) (*warehouse.Stock, *model.AppError) {
	stock, err := a.Srv().Store.Stock().GetbyOption(option)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("GetStockByOption", "app.warehouse.error_finding_stock_by_option.app_error", err)
	}

	return stock, nil
}

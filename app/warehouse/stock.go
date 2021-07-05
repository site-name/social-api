package warehouse

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/model/warehouse"
	"github.com/sitename/sitename/store"
)

func (a *AppWarehouse) CheckStockQuantity(variant *product_and_discount.ProductVariant, countryCode string, channelSlug string, quantity uint) *model.AppError {
	if *variant.TrackInventory {
		_, appErr := a.GetVariantStocksForCountry(countryCode, channelSlug, variant.Id)
		return appErr
	}

	return nil
}

func (a *AppWarehouse) GetVariantStocksForCountry(countryCode string, channelSlug string, variantID string) (interface{}, *model.AppError) {
	_, _, _, err := a.Srv().Store.Stock().FilterVariantStocksForCountry(
		&warehouse.ForCountryAndChannelFilter{
			CountryCode: countryCode,
			ChannelSlug: channelSlug,
		},
		variantID,
	)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("GetVariantStocksForCountry", "app.warehouse.stock_filter_forcountry_missing.app_error", err)
	}

	return nil, nil
}

package warehouse

import (
	"net/http"

	"github.com/mattermost/gorp"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/warehouse"
	"github.com/sitename/sitename/store"
)

// BulkUpsertStocks updates or insderts given stock based on its Id property
func (a *AppWarehouse) BulkUpsertStocks(transaction *gorp.Transaction, stocks []*warehouse.Stock) ([]*warehouse.Stock, *model.AppError) {
	stocks, err := a.Srv().Store.Stock().BulkUpsert(transaction, stocks)
	if err != nil {
		if appErr, ok := err.(*model.AppError); ok {
			return nil, appErr
		}
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		} else if _, ok := err.(*store.ErrInvalidInput); ok {
			statusCode = http.StatusBadRequest
		}

		return nil, model.NewAppError("UpsertStocks", "app.warehouse.error_upserting_stocks.app_error", nil, err.Error(), statusCode)
	}

	return stocks, nil
}

// StocksByOption returns a list of stocks filtered using given options
func (a *AppWarehouse) StocksByOption(transaction *gorp.Transaction, option *warehouse.StockFilterOption) ([]*warehouse.Stock, *model.AppError) {
	stocks, err := a.Srv().Store.Stock().FilterByOption(transaction, option)
	var (
		statusCode   int
		errorMessage string
	)
	if err != nil {
		statusCode = http.StatusInternalServerError
		errorMessage = err.Error()
	} else if len(stocks) == 0 {
		statusCode = http.StatusNotFound
	}

	if statusCode != 0 {
		return nil, model.NewAppError("StocksByOption", "app.warehouse.error_finding_stocks_by_option.app_error", nil, errorMessage, statusCode)
	}

	return stocks, nil
}

// GetVariantStocksForCountry Return the stock information about the a stock for a given country.
//
// Note it will raise a 'Stock.DoesNotExist' exception if no such stock is found.
func (a *AppWarehouse) GetVariantStocksForCountry(countryCode string, channelSlug string, variantID string, quantity int) ([]*warehouse.Stock, *model.AppError) {
	stocks, err := a.Srv().Store.Stock().FilterVariantStocksForCountry(&warehouse.StockFilterForCountryAndChannel{
		CountryCode:      countryCode,
		ChannelSlug:      channelSlug,
		ProductVariantID: variantID,
	})
	var (
		statusCode   int
		errorMessage string
	)
	if err != nil {
		statusCode = http.StatusInternalServerError
		errorMessage = err.Error()
	} else if len(stocks) == 0 {
		statusCode = http.StatusNotFound
	}

	if statusCode != 0 {
		return nil, model.NewAppError("GetVariantStocksForCountry", "app.warehouse.error_finding_variant_stocks_for_country.app_error", nil, errorMessage, statusCode)
	}

	return stocks, nil
}

// GetProductStocksForCountryAndChannel finds stocks by given options
// func (a *AppWarehouse) GetProductStocksForCountryAndChannel(countryCode string, channelSlug string, productID string) ([]*warehouse.Stock, *model.AppError) {
// 	stocks, err := a.Srv().Store.Stock().FilterProductStocksForCountryAndChannel(&warehouse.StockFilterOption{
// 		CountryCode: countryCode,
// 		ChannelSlug: channelSlug,
// 		ProductID:   productID,
// 	})
// 	var (
// 		statusCode   int
// 		errorMessage string
// 	)
// 	if err != nil {
// 		statusCode = http.StatusInternalServerError
// 		errorMessage = err.Error()
// 	} else if len(stocks) == 0 {
// 		statusCode = http.StatusNotFound
// 	}

// 	if statusCode != 0 {
// 		return nil, model.NewAppError("GetProductStocksForCountryAndChannel", "app.warehouse.error_finding_product_stocks_for_country_and_chanel.app_error", nil, errorMessage, statusCode)
// 	}

// 	return stocks, nil
// }

// GetStocksForCountryAndChannel returns stocks with given arguments, it also attach related data to returning stocks
// func (a *AppWarehouse) GetStocksForCountryAndChannel(countryCode string, channelSlug string, lockForUpdate bool) ([]*warehouse.Stock, *model.AppError) {
// 	stocks, err := a.Srv().Store.Stock().FilterProductStocksForCountryAndChannel(&warehouse.StockFilterOption{
// 		CountryCode:   countryCode,
// 		ChannelSlug:   channelSlug,
// 		LockForUpdate: lockForUpdate,
// 	})
// 	var (
// 		statusCode   int
// 		errorMessage string
// 	)
// 	if err != nil {
// 		statusCode = http.StatusInternalServerError
// 		errorMessage = err.Error()
// 	} else if len(stocks) == 0 {
// 		statusCode = http.StatusNotFound
// 	}

// 	if statusCode != 0 {
// 		return nil, model.NewAppError("GetStocksForCountryAndChannel", "app.warehouse.error_finding_stocks_for_country_and_chanel.app_error", nil, errorMessage, statusCode)
// 	}

// 	return stocks, nil
// }

// GetStockByOption takes options for filtering 1 stock
func (a *AppWarehouse) GetStockByOption(option *warehouse.StockFilterOption) (*warehouse.Stock, *model.AppError) {
	stock, err := a.Srv().Store.Stock().GetbyOption(option)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("GetStockByOption", "app.warehouse.error_finding_stock_by_option.app_error", err)
	}

	return stock, nil
}

// StockIncreaseQuantity Return given quantity of product to a stock.
func (a *AppWarehouse) StockIncreaseQuantity(stockID string, quantity int) *model.AppError {
	err := a.Srv().Store.Stock().ChangeQuantity(stockID, quantity)
	if err != nil {
		return model.NewAppError("StockIncrease", "app.warehouse.error_increasing_stock_quantity", nil, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

// StockDecreaseQuantity Return given quantity of product to a stock.
func (a *AppWarehouse) StockDecreaseQuantity(stockID string, quantity int) *model.AppError {
	err := a.Srv().Store.Stock().ChangeQuantity(stockID, -quantity)
	if err != nil {
		return model.NewAppError("StockIncrease", "app.warehouse.error_increasing_stock_quantity", nil, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

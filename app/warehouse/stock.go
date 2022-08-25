package warehouse

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/warehouse"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/store_iface"
)

// BulkUpsertStocks updates or insderts given stock based on its Id property
func (a *ServiceWarehouse) BulkUpsertStocks(transaction store_iface.SqlxTxExecutor, stocks []*warehouse.Stock) ([]*warehouse.Stock, *model.AppError) {
	stocks, err := a.srv.Store.Stock().BulkUpsert(transaction, stocks)
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
func (a *ServiceWarehouse) StocksByOption(transaction store_iface.SqlxTxExecutor, option *warehouse.StockFilterOption) (warehouse.Stocks, *model.AppError) {
	stocks, err := a.srv.Store.Stock().FilterByOption(transaction, option)
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
func (a *ServiceWarehouse) GetVariantStocksForCountry(transaction store_iface.SqlxTxExecutor, countryCode string, channelSlug string, variantID string) ([]*warehouse.Stock, *model.AppError) {
	stocks, err := a.srv.Store.Stock().FilterVariantStocksForCountry(transaction, &warehouse.StockFilterForCountryAndChannel{
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

// GetProductStocksForCountryAndChannel
func (a *ServiceWarehouse) GetProductStocksForCountryAndChannel(transaction store_iface.SqlxTxExecutor, options *warehouse.StockFilterForCountryAndChannel) ([]*warehouse.Stock, *model.AppError) {
	stocks, err := a.srv.Store.Stock().FilterProductStocksForCountryAndChannel(transaction, options)
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
		return nil, model.NewAppError("GetProductStocksForCountryAndChannel", ErrorFindingStocksId, nil, errorMessage, statusCode)
	}

	return stocks, nil
}

// FilterStocksForCountryAndChannel finds stocks by given options
func (a *ServiceWarehouse) FilterStocksForCountryAndChannel(transaction store_iface.SqlxTxExecutor, options *warehouse.StockFilterForCountryAndChannel) ([]*warehouse.Stock, *model.AppError) {
	stocks, err := a.srv.Store.Stock().FilterForCountryAndChannel(transaction, options)
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
		return nil, model.NewAppError("FilterStocksForCountryAndChannel", ErrorFindingStocksId, nil, errorMessage, statusCode)
	}

	return stocks, nil
}

// GetStockById takes options for filtering 1 stock
func (a *ServiceWarehouse) GetStockById(stockID string) (*warehouse.Stock, *model.AppError) {
	stock, err := a.srv.Store.Stock().Get(stockID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}
		return nil, model.NewAppError("GetStockById", "app.warehouse.error_finding_stock_by_option.app_error", nil, err.Error(), statusCode)
	}

	return stock, nil
}

// FilterStocksForChannel returns a slice of stocks that filtered using given options
func (a *ServiceWarehouse) FilterStocksForChannel(option *warehouse.StockFilterForChannelOption) ([]*warehouse.Stock, *model.AppError) {
	_, stocks, err := a.srv.Store.Stock().FilterForChannel(option)
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
		return nil, model.NewAppError("FilterStocksByChannel", "app.warehouse.error_finding_stocks_for_channel.app_error", nil, errorMessage, statusCode)
	}

	return stocks, nil
}

// StockIncreaseQuantity Return given quantity of product to a stock.
func (a *ServiceWarehouse) StockIncreaseQuantity(stockID string, quantity int) *model.AppError {
	err := a.srv.Store.Stock().ChangeQuantity(stockID, quantity)
	if err != nil {
		return model.NewAppError("StockIncrease", "app.warehouse.error_increasing_stock_quantity", nil, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

// StockDecreaseQuantity Return given quantity of product to a stock.
func (a *ServiceWarehouse) StockDecreaseQuantity(stockID string, quantity int) *model.AppError {
	err := a.srv.Store.Stock().ChangeQuantity(stockID, -quantity)
	if err != nil {
		return model.NewAppError("StockIncrease", "app.warehouse.error_increasing_stock_quantity", nil, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

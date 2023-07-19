package warehouse

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

// BulkUpsertStocks updates or insderts given stock based on its Id property
func (a *ServiceWarehouse) BulkUpsertStocks(transaction *gorm.DB, stocks []*model.Stock) ([]*model.Stock, *model.AppError) {
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
func (a *ServiceWarehouse) StocksByOption(option *model.StockFilterOption) (model.Stocks, *model.AppError) {
	stocks, err := a.srv.Store.Stock().FilterByOption(option)
	if err != nil {
		return nil, model.NewAppError("StocksByOption", "app.warehouse.error_finding_stocks_by_option.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return stocks, nil
}

// GetVariantStocksForCountry Return the stock information about the a stock for a given country.
//
// Note it will raise a 'Stock.DoesNotExist' exception if no such stock is found.
func (a *ServiceWarehouse) GetVariantStocksForCountry(countryCode model.CountryCode, channelSlug string, variantID string) ([]*model.Stock, *model.AppError) {
	stocks, err := a.srv.Store.Stock().FilterVariantStocksForCountry(&model.StockFilterForCountryAndChannel{
		CountryCode:      countryCode,
		ChannelSlug:      channelSlug,
		ProductVariantID: variantID,
	})

	if err != nil {
		return nil, model.NewAppError("GetVariantStocksForCountry", "app.warehouse.error_finding_variant_stocks_for_country.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return stocks, nil
}

// GetProductStocksForCountryAndChannel
func (a *ServiceWarehouse) GetProductStocksForCountryAndChannel(options *model.StockFilterForCountryAndChannel) ([]*model.Stock, *model.AppError) {
	stocks, err := a.srv.Store.Stock().FilterProductStocksForCountryAndChannel(options)
	if err != nil {
		return nil, model.NewAppError("GetProductStocksForCountryAndChannel", ErrorFindingStocksId, nil, err.Error(), http.StatusInternalServerError)
	}

	return stocks, nil
}

// FilterStocksForCountryAndChannel finds stocks by given options
func (a *ServiceWarehouse) FilterStocksForCountryAndChannel(options *model.StockFilterForCountryAndChannel) ([]*model.Stock, *model.AppError) {
	stocks, err := a.srv.Store.Stock().FilterForCountryAndChannel(options)

	if err != nil {
		return nil, model.NewAppError("FilterStocksForCountryAndChannel", ErrorFindingStocksId, nil, err.Error(), http.StatusInternalServerError)
	}

	return stocks, nil
}

// GetStockById takes options for filtering 1 stock
func (a *ServiceWarehouse) GetStockById(stockID string) (*model.Stock, *model.AppError) {
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
func (a *ServiceWarehouse) FilterStocksForChannel(option *model.StockFilterForChannelOption) ([]*model.Stock, *model.AppError) {
	_, stocks, err := a.srv.Store.Stock().FilterForChannel(option)

	if err != nil {
		return nil, model.NewAppError("FilterStocksByChannel", "app.warehouse.error_finding_stocks_for_channel.app_error", nil, err.Error(), http.StatusInternalServerError)
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

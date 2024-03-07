package warehouse

import (
	"net/http"

	"github.com/mattermost/squirrel"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

// BulkUpsertStocks updates or insderts given stock based on its Id property
func (a *ServiceWarehouse) BulkUpsertStocks(transaction *gorm.DB, stocks []*model.Stock) ([]*model.Stock, *model_helper.AppError) {
	stocks, err := a.srv.Store.Stock().BulkUpsert(transaction, stocks)
	if err != nil {
		if appErr, ok := err.(*model_helper.AppError); ok {
			return nil, appErr
		}
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		} else if _, ok := err.(*store.ErrInvalidInput); ok {
			statusCode = http.StatusBadRequest
		}

		return nil, model_helper.NewAppError("UpsertStocks", "app.warehouse.error_upserting_stocks.app_error", nil, err.Error(), statusCode)
	}

	return stocks, nil
}

// StocksByOption returns a list of stocks filtered using given options
func (a *ServiceWarehouse) StocksByOption(option *model.StockFilterOption) (int64, model.Stocks, *model_helper.AppError) {
	total, stocks, err := a.srv.Store.Stock().FilterByOption(option)
	if err != nil {
		return 0, nil, model_helper.NewAppError("StocksByOption", "app.warehouse.error_finding_stocks_by_option.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return total, stocks, nil
}

// GetVariantStocksForCountry Return the stock information about the a stock for a given country.
//
// Note it will raise a 'Stock.DoesNotExist' exception if no such stock is found.
func (a *ServiceWarehouse) GetVariantStocksForCountry(countryCode model.CountryCode, channelSlug string, variantID string) ([]*model.Stock, *model_helper.AppError) {
	stocks, err := a.srv.Store.Stock().FilterVariantStocksForCountry(&model.StockFilterOptionsForCountryAndChannel{
		CountryCode:      countryCode,
		ChannelSlug:      channelSlug,
		ProductVariantID: variantID,
	})

	if err != nil {
		return nil, model_helper.NewAppError("GetVariantStocksForCountry", "app.warehouse.error_finding_variant_stocks_for_country.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return stocks, nil
}

// GetProductStocksForCountryAndChannel
func (a *ServiceWarehouse) GetProductStocksForCountryAndChannel(options *model.StockFilterOptionsForCountryAndChannel) ([]*model.Stock, *model_helper.AppError) {
	stocks, err := a.srv.Store.Stock().FilterProductStocksForCountryAndChannel(options)
	if err != nil {
		return nil, model_helper.NewAppError("GetProductStocksForCountryAndChannel", ErrorFindingStocksId, nil, err.Error(), http.StatusInternalServerError)
	}

	return stocks, nil
}

// FilterStocksForCountryAndChannel finds stocks by given options
func (a *ServiceWarehouse) FilterStocksForCountryAndChannel(options *model.StockFilterOptionsForCountryAndChannel) (model.Stocks, *model_helper.AppError) {
	stocks, err := a.srv.Store.Stock().FilterForCountryAndChannel(options)

	if err != nil {
		return nil, model_helper.NewAppError("FilterStocksForCountryAndChannel", ErrorFindingStocksId, nil, err.Error(), http.StatusInternalServerError)
	}

	return stocks, nil
}

// GetStockById takes options for filtering 1 stock
func (a *ServiceWarehouse) GetStockById(stockID string) (*model.Stock, *model_helper.AppError) {
	stock, err := a.srv.Store.Stock().Get(stockID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}
		return nil, model_helper.NewAppError("GetStockById", "app.warehouse.error_finding_stock_by_option.app_error", nil, err.Error(), statusCode)
	}

	return stock, nil
}

// FilterStocksForChannel returns a slice of stocks that filtered using given options
func (a *ServiceWarehouse) FilterStocksForChannel(option model_helper.StockFilterForChannelOption) ([]*model.Stock, *model_helper.AppError) {
	_, stocks, err := a.srv.Store.Stock().FilterForChannel(option)

	if err != nil {
		return nil, model_helper.NewAppError("FilterStocksByChannel", "app.warehouse.error_finding_stocks_for_channel.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return stocks, nil
}

// StockIncreaseQuantity Return given quantity of product to a stock.
func (a *ServiceWarehouse) StockIncreaseQuantity(stockID string, quantity int) *model_helper.AppError {
	err := a.srv.Store.Stock().ChangeQuantity(stockID, quantity)
	if err != nil {
		return model_helper.NewAppError("StockIncrease", "app.warehouse.error_increasing_stock_quantity", nil, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

// StockDecreaseQuantity Return given quantity of product to a stock.
func (a *ServiceWarehouse) StockDecreaseQuantity(stockID string, quantity int) *model_helper.AppError {
	err := a.srv.Store.Stock().ChangeQuantity(stockID, -quantity)
	if err != nil {
		return model_helper.NewAppError("StockIncrease", "app.warehouse.error_increasing_stock_quantity", nil, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

func (s *ServiceWarehouse) DeleteStocks(options *model.StockFilterOption) (int64, *model_helper.AppError) {
	// begin
	tx := s.srv.Store.GetMaster().Begin()
	if tx.Error != nil {
		return 0, model_helper.NewAppError("DeleteStocks", model.ErrorCreatingTransactionErrorID, nil, tx.Error.Error(), http.StatusInternalServerError)
	}
	defer s.srv.Store.FinalizeTransaction(tx)

	_, stocksToDelete, appErr := s.StocksByOption(options)
	if appErr != nil {
		return 0, appErr
	}

	numDeleted, err := s.srv.Store.Stock().Delete(tx, &model.StockFilterOption{
		Conditions: squirrel.Eq{model.StockTableName + "." + model.StockColumnId: stocksToDelete.IDs()},
	})
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrInvalidInput); ok {
			statusCode = http.StatusBadRequest
		}
		return 0, model_helper.NewAppError("DeleteStocks", "app.warehouse.delete_stocks.app_error", nil, err.Error(), statusCode)
	}

	// commit
	err = tx.Commit().Error
	if err != nil {
		return 0, model_helper.NewAppError("DeleteStocks", model.ErrorCommittingTransactionErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	// perform plugin callbacks
	pluginMng := s.srv.PluginService().GetPluginManager()
	for _, stock := range stocksToDelete {
		appErr = pluginMng.ProductVariantOutOfStock(*stock)
		if appErr != nil {
			return 0, appErr
		}
	}

	return numDeleted, nil
}

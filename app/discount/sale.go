package discount

import (
	"net/http"
	"time"

	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/discount/types"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
)

func (a *ServiceDiscount) GetSaleDiscount(sale *product_and_discount.Sale, saleChannelListing *product_and_discount.SaleChannelListing) (types.DiscountCalculator, *model.AppError) {
	if saleChannelListing == nil {
		return nil, model.NewAppError("GetSaleDiscount", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "saleChannelListing"}, "", http.StatusBadRequest)
	}

	if sale.Type == product_and_discount.FIXED {
		discountAmount := &goprices.Money{ // can use directly here since sale channel listings are validated before saving
			Amount:   *saleChannelListing.DiscountValue,
			Currency: saleChannelListing.Currency,
		}
		return Decorator(discountAmount), nil
	}
	return Decorator(saleChannelListing.DiscountValue), nil
}

// FilterSalesByOption should be used to filter active or expired sales
// refer: saleor/discount/models.SaleQueryset for details
func (a *ServiceDiscount) FilterSalesByOption(option *product_and_discount.SaleFilterOption) ([]*product_and_discount.Sale, *model.AppError) {
	sales, err := a.srv.Store.DiscountSale().FilterSalesByOption(option)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("FilterSalesByOption", "app.discount.filter_sales_by_option.app_error", err)
	}

	return sales, nil
}

// ActiveSales finds active sales by given date. If date is nil then set date to UTC now
//
//  (end_date == NULL || end_date >= date) && start_date <= date
func (a *ServiceDiscount) ActiveSales(date *time.Time) (product_and_discount.Sales, *model.AppError) {
	if date == nil {
		date = util.NewTime(time.Now().UTC())
	}

	activeSalesByDate, err := a.srv.Store.DiscountSale().
		FilterSalesByOption(&product_and_discount.SaleFilterOption{
			EndDate: &model.TimeFilter{
				Or: &model.TimeOption{
					NULL: model.NewBool(true),
					GtE:  date,
				},
			},
			StartDate: &model.TimeFilter{
				TimeOption: &model.TimeOption{
					LtE: date,
				},
			},
		})
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("ActiveSales", "app.discount.active_sales_by_date.app_error", err)
	}

	return activeSalesByDate, nil
}

// ExpiredSales returns sales that are expired by date. If date is nil, default to UTC now
//
//  end_date <= date && start_date <= date
func (a *ServiceDiscount) ExpiredSales(date *time.Time) ([]*product_and_discount.Sale, *model.AppError) {
	if date == nil {
		date = util.NewTime(time.Now().UTC())
	}

	expiredSalesByDate, err := a.srv.Store.DiscountSale().
		FilterSalesByOption(&product_and_discount.SaleFilterOption{
			EndDate: &model.TimeFilter{
				TimeOption: &model.TimeOption{
					Lt: date,
				},
			},
			StartDate: &model.TimeFilter{
				TimeOption: &model.TimeOption{
					Lt: date,
				},
			},
		})

	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("ExpiredSales", "app.discount.expired_sales_by_date.app_error", err)
	}

	return expiredSalesByDate, nil
}

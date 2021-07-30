package discount

import (
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

func (a *AppDiscount) GetSaleDiscount(sale *product_and_discount.Sale, saleChannelListing *product_and_discount.SaleChannelListing) DiscountCalculator {
	if saleChannelListing == nil {
		return nil
	}

	if sale.Type == product_and_discount.FIXED {
		discountAmount := &goprices.Money{ // can use directly here since sale channel listings are validated before saving
			Amount:   saleChannelListing.DiscountValue,
			Currency: saleChannelListing.Currency,
		}
		return decorator(discountAmount)
	}
	return decorator(saleChannelListing.DiscountValue)
}

// FilterSalesByOption should be used to filter active or expired sales
// refer: saleor/discount/models.SaleQueryset for details
func (a *AppDiscount) FilterSalesByOption(option *product_and_discount.SaleFilterOption) ([]*product_and_discount.Sale, *model.AppError) {
	sales, err := a.Srv().Store.DiscountSale().FilterSalesByOption(option)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("FilterSalesByOption", "app.discount.filter_sales_by_option.app_error", err)
	}

	return sales, nil
}

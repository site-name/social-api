package discount

import (
	"net/http"

	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

func (a *AppDiscount) GetSaleDiscount(sale *product_and_discount.Sale, saleChannelListing *product_and_discount.SaleChannelListing) (DiscountCalculator, *model.AppError) {
	if saleChannelListing == nil {
		return nil, model.NewAppError("GetSaleDiscount", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "sale channel listing"}, "", http.StatusBadRequest)
	}

	if sale.Type == product_and_discount.FIXED {
		discountAmount := &goprices.Money{ // can use directly here since sale channel listings are validated before saving
			Amount:   saleChannelListing.DiscountValue,
			Currency: saleChannelListing.Currency,
		}
		return decorator(discountAmount), nil
	}
	return decorator(saleChannelListing.DiscountValue), nil
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

func (a *AppDiscount) SaleCategoriesByOption(option *product_and_discount.SaleCategoryRelationFilterOption) ([]*product_and_discount.SaleCategoryRelation, *model.AppError) {
	saleCategories, err := a.Srv().Store.SaleCategoryRelation().SaleCategoriesByOption(option)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("SaleCategoriesByOption", "app.discount.error_finding_sale_categories_by_option.app_error", err)
	}

	return saleCategories, nil
}

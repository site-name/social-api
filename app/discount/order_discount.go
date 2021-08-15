package discount

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

// OrderDiscountsByOption filters and returns order discounts with given option
func (a *AppDiscount) OrderDiscountsByOption(option *product_and_discount.OrderDiscountFilterOption) ([]*product_and_discount.OrderDiscount, *model.AppError) {
	orderDiscounts, err := a.Srv().Store.OrderDiscount().FilterbyOption(option)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("OrderDiscountsByOption", "app.discount.order_discount_by_option.app_error.app_error", err)
	}

	return orderDiscounts, nil
}

// UpsertOrderDiscount updates or inserts given order discount
func (a *AppDiscount) UpsertOrderDiscount(orderDiscount *product_and_discount.OrderDiscount) (*product_and_discount.OrderDiscount, *model.AppError) {
	orderDiscount, err := a.Srv().Store.OrderDiscount().Upsert(orderDiscount)
	if err != nil {
		if appErr, ok := err.(*model.AppError); ok {
			return nil, appErr
		}
		return nil, model.NewAppError("UpsertOrderDiscount", "app.error_upserting_order_discount.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return orderDiscount, nil
}

// BulkDeleteOrderDiscounts performs bulk delete given order discounts
func (a *AppDiscount) BulkDeleteOrderDiscounts(orderDiscountIDs []string) *model.AppError {
	err := a.Srv().Store.OrderDiscount().BulkDelete(orderDiscountIDs)
	if err != nil {
		return model.NewAppError("BulkDeleteOrderDiscounts", "app.discount.error_bulk_deleting_order_discounts.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}

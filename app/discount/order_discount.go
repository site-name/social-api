package discount

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/store_iface"
)

// OrderDiscountsByOption filters and returns order discounts with given option
func (a *ServiceDiscount) OrderDiscountsByOption(option *product_and_discount.OrderDiscountFilterOption) ([]*product_and_discount.OrderDiscount, *model.AppError) {
	orderDiscounts, err := a.srv.Store.OrderDiscount().FilterbyOption(option)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("OrderDiscountsByOption", "app.discount.order_discount_by_option.app_error.app_error", err)
	}

	return orderDiscounts, nil
}

// UpsertOrderDiscount updates or inserts given order discount
func (a *ServiceDiscount) UpsertOrderDiscount(transaction store_iface.SqlxTxExecutor, orderDiscount *product_and_discount.OrderDiscount) (*product_and_discount.OrderDiscount, *model.AppError) {
	orderDiscount, err := a.srv.Store.OrderDiscount().Upsert(transaction, orderDiscount)
	if err != nil {
		if appErr, ok := err.(*model.AppError); ok {
			return nil, appErr
		}
		return nil, model.NewAppError("UpsertOrderDiscount", "app.error_upserting_order_discount.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return orderDiscount, nil
}

// BulkDeleteOrderDiscounts performs bulk delete given order discounts
func (a *ServiceDiscount) BulkDeleteOrderDiscounts(orderDiscountIDs []string) *model.AppError {
	err := a.srv.Store.OrderDiscount().BulkDelete(orderDiscountIDs)
	if err != nil {
		return model.NewAppError("BulkDeleteOrderDiscounts", "app.discount.error_bulk_deleting_order_discounts.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}

package discount

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"gorm.io/gorm"
)

// OrderDiscountsByOption filters and returns order discounts with given option
func (a *ServiceDiscount) OrderDiscountsByOption(option *model.OrderDiscountFilterOption) ([]*model.OrderDiscount, *model.AppError) {
	orderDiscounts, err := a.srv.Store.OrderDiscount().FilterbyOption(option)
	if err != nil {
		return nil, model.NewAppError("OrderDiscountsByOption", "app.discount.order_discount_by_option.app_error.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return orderDiscounts, nil
}

// UpsertOrderDiscount updates or inserts given order discount
func (a *ServiceDiscount) UpsertOrderDiscount(transaction *gorm.DB, orderDiscount *model.OrderDiscount) (*model.OrderDiscount, *model.AppError) {
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

package discount

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func (a *ServiceDiscount) OrderDiscountsByOption(option model_helper.OrderDiscountFilterOption) (model.OrderDiscountSlice, *model_helper.AppError) {
	orderDiscounts, err := a.srv.Store.OrderDiscount().FilterbyOption(option)
	if err != nil {
		return nil, model_helper.NewAppError("OrderDiscountsByOption", "app.discount.order_discount_by_option.app_error.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return orderDiscounts, nil
}

func (a *ServiceDiscount) UpsertOrderDiscount(transaction boil.ContextTransactor, orderDiscount model.OrderDiscount) (*model.OrderDiscount, *model_helper.AppError) {
	upsertOrderDiscount, err := a.srv.Store.OrderDiscount().Upsert(transaction, orderDiscount)
	if err != nil {
		if appErr, ok := err.(*model_helper.AppError); ok {
			return nil, appErr
		}
		return nil, model_helper.NewAppError("UpsertOrderDiscount", "app.error_upserting_order_discount.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return upsertOrderDiscount, nil
}

func (a *ServiceDiscount) BulkDeleteOrderDiscounts(orderDiscountIDs []string) *model_helper.AppError {
	err := a.srv.Store.OrderDiscount().BulkDelete(orderDiscountIDs)
	if err != nil {
		return model_helper.NewAppError("BulkDeleteOrderDiscounts", "app.discount.error_bulk_deleting_order_discounts.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}

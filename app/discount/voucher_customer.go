package discount

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
)

func (a *ServiceDiscount) CreateNewVoucherCustomer(voucherID string, customerEmail string) (*model.VoucherCustomer, *model_helper.AppError) {
	voucher, err := a.srv.Store.VoucherCustomer().Save(model.VoucherCustomer{
		CustomerEmail: customerEmail,
		VoucherID:     voucherID,
	})
	if err != nil {
		if appErr, ok := err.(*model_helper.AppError); ok {
			return nil, appErr
		}
		return nil, model_helper.NewAppError("CreateNewVoucherCustomer", "app.discount.error_creating_new_customer_voucher.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return voucher, nil
}

func (a *ServiceDiscount) VoucherCustomerByOptions(options model_helper.VoucherCustomerFilterOption) (*model.VoucherCustomer, *model_helper.AppError) {
	voucherCustomer, err := a.srv.Store.VoucherCustomer().GetByOption(options)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := err.(*store.ErrNotFound); ok {
			statusCode = http.StatusNotFound
		}
		return nil, model_helper.NewAppError("VoucherCustomerByOptions", "app.discount.voucher_customer_missing.app_error", nil, err.Error(), statusCode)
	}

	return voucherCustomer, nil
}

func (s *ServiceDiscount) VoucherCustomersByOption(options model_helper.VoucherCustomerFilterOption) (model.VoucherCustomerSlice, *model_helper.AppError) {
	voucherCustomers, err := s.srv.Store.VoucherCustomer().FilterByOptions(options)
	if err != nil {
		return nil, model_helper.NewAppError("VoucherCustomersByOption", "app.discount.error_finding_voucher_customers_by_options", nil, err.Error(), http.StatusInternalServerError)
	}

	return voucherCustomers, nil
}

package discount

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

// CreateNewVoucherCustomer tells store to insert new voucher customer into database, then returns it
func (a *ServiceDiscount) CreateNewVoucherCustomer(voucherID string, customerEmail string) (*model.VoucherCustomer, *model.AppError) {
	voucher, err := a.srv.Store.VoucherCustomer().Save(&model.VoucherCustomer{
		CustomerEmail: customerEmail,
		VoucherID:     voucherID,
	})
	if err != nil {
		if appErr, ok := err.(*model.AppError); ok {
			return nil, appErr
		}
		return nil, model.NewAppError("CreateNewVoucherCustomer", "app.discount.error_creating_new_customer_voucher.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return voucher, nil
}

// VoucherCustomerByOptions finds a voucher customer relation and returns it with an error
func (a *ServiceDiscount) VoucherCustomerByOptions(options *model.VoucherCustomerFilterOption) (*model.VoucherCustomer, *model.AppError) {
	voucherCustomer, err := a.srv.Store.VoucherCustomer().GetByOption(options)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("VoucherCustomerByOptions", "app.discount.voucher_customer_missing.app_error", err)
	}

	return voucherCustomer, nil
}

// VoucherCustomersByOption returns a slice of voucher customers filtered using given options
func (s *ServiceDiscount) VoucherCustomersByOption(options *model.VoucherCustomerFilterOption) ([]*model.VoucherCustomer, *model.AppError) {
	voucherCustomers, err := s.srv.Store.VoucherCustomer().FilterByOptions(options)
	var (
		statusCode int
		errMessage string
	)
	if err != nil {
		statusCode = http.StatusInternalServerError
		errMessage = err.Error()
	} else if len(voucherCustomers) == 0 {
		statusCode = http.StatusNotFound
	}

	if statusCode != 0 {
		return nil, model.NewAppError("VoucherCustomersByOption", "app.discount.error_finding_voucher_customers_by_options", nil, errMessage, statusCode)
	}

	return voucherCustomers, nil
}

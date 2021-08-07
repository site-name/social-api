package discount

import (
	"net/http"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

const (
	VoucherCustomerNotFoundErrId = "app.discount.voucher_customer_missing.app_error"
)

func (a *AppDiscount) CreateNewVoucherCustomer(voucherID string, customerEmail string) (*product_and_discount.VoucherCustomer, *model.AppError) {
	voucher, err := a.Srv().Store.VoucherCustomer().Save(&product_and_discount.VoucherCustomer{
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

// VoucherCustomerById finds a voucher customer relation and returns it with an error
func (a *AppDiscount) VoucherCustomerById(id string) (*product_and_discount.VoucherCustomer, *model.AppError) {
	voucherCustomer, err := a.Srv().Store.VoucherCustomer().Get(id)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("VoucherCustomerById", VoucherCustomerNotFoundErrId, err)
	}

	return voucherCustomer, nil
}

// VoucherCustomerByCustomerEmailAndVoucherID finds voucher customer with given voucherID and customerEmail
func (a *AppDiscount) VoucherCustomerByCustomerEmailAndVoucherID(voucherID string, customerEmail string) ([]*product_and_discount.VoucherCustomer, *model.AppError) {
	res, err := a.Srv().Store.VoucherCustomer().FilterByEmailAndCustomerEmail(voucherID, customerEmail)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("VoucherCustomerByCustomerEmailAndVoucherID", VoucherCustomerNotFoundErrId, err)
	}

	return res, nil
}

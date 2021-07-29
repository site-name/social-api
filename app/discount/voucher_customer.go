package discount

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

const (
	VoucherCustomerNotFoundErrId = "app.discount.voucher_customer_missing.app_error"
)

// VoucherCustomerById finds a voucher customer relation and returns it with an error
func (a *AppDiscount) VoucherCustomerById(id string) (*product_and_discount.VoucherCustomer, *model.AppError) {
	voucherCustomer, err := a.Srv().Store.VoucherCustomer().Get(id)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("VoucherCustomerById", VoucherCustomerNotFoundErrId, err)
	}

	return voucherCustomer, nil
}

// VoucherCustomerByCustomerEmailAndVoucherID finds voucher customer with given voucherID and customerEmail
func (a *AppDiscount) VoucherCustomerByCustomerEmailAndVoucherID(voucherID string, customerEmail string) (*product_and_discount.VoucherCustomer, *model.AppError) {
	res, err := a.Srv().Store.VoucherCustomer().FilterByVoucherAndEmail(voucherID, customerEmail)
	if err != nil {
		return nil, store.AppErrorFromDatabaseLookupError("VoucherCustomerByCustomerEmailAndVoucherID", VoucherCustomerNotFoundErrId, err)
	}

	return res, nil
}

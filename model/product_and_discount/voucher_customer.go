package product_and_discount

import (
	"github.com/sitename/sitename/model"
)

type VoucherCustomer struct {
	Id            string `json:"id"`
	VoucherID     string `json:"voucher_id"`
	CustomerEmail string `json:"customer_email"`
}

func (vc *VoucherCustomer) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.voucher_customer.is_valid.%s.app_error",
		"voucher_customer_id=",
		"VoucherCustomer.IsValid",
	)
	if !model.IsValidId(vc.Id) {
		return outer("id", nil)
	}
	if !model.IsValidId(vc.VoucherID) {
		return outer("voucher_id", &vc.Id)
	}
	if !model.IsValidEmail(vc.CustomerEmail) {
		return outer("customer_email", &vc.Id)
	}

	return nil
}

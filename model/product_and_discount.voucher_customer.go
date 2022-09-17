package model

import (
	"github.com/Masterminds/squirrel"
)

// VoucherCustomer represents m2m relation ship between customers and vouchers
type VoucherCustomer struct {
	Id            string `json:"id"`
	VoucherID     string `json:"voucher_id"`
	CustomerEmail string `json:"customer_email"`
}

// VoucherCustomerFilterOption is used to build squirrel sql queries
type VoucherCustomerFilterOption struct {
	Id            squirrel.Sqlizer
	VoucherID     squirrel.Sqlizer
	CustomerEmail squirrel.Sqlizer
}

func (vc *VoucherCustomer) PreSave() {
	if vc.Id == "" {
		vc.Id = NewId()
	}
}

func (vc *VoucherCustomer) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"voucher_customer.is_valid.%s.app_error",
		"voucher_customer_id=",
		"VoucherCustomer.IsValid",
	)
	if !IsValidId(vc.Id) {
		return outer("id", nil)
	}
	if !IsValidId(vc.VoucherID) {
		return outer("voucher_id", &vc.Id)
	}
	if !IsValidEmail(vc.CustomerEmail) {
		return outer("customer_email", &vc.Id)
	}

	return nil
}

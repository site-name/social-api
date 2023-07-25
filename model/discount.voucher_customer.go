package model

import (
	"net/http"

	"github.com/Masterminds/squirrel"
	"gorm.io/gorm"
)

// VoucherCustomer represents m2m relation ship between customers and vouchers
type VoucherCustomer struct {
	Id            string `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	VoucherID     string `json:"voucher_id" gorm:"type:uuid;column:VoucherID;index:voucherid_customeremail_key"`
	CustomerEmail string `json:"customer_email" gorm:"type:varchar(128);column:CustomerEmail;index:voucherid_customeremail_key"`
}

func (c *VoucherCustomer) BeforeCreate(_ *gorm.DB) error { return c.IsValid() }
func (c *VoucherCustomer) BeforeUpdate(_ *gorm.DB) error { return c.IsValid() }
func (c *VoucherCustomer) TableName() string             { return VoucherCustomerTableName }

// VoucherCustomerFilterOption is used to build squirrel sql queries
type VoucherCustomerFilterOption struct {
	Conditions squirrel.Sqlizer
}

func (vc *VoucherCustomer) IsValid() *AppError {
	if !IsValidId(vc.VoucherID) {
		return NewAppError("VoucherCustomer.IsValid", "model.voucher_customer.is_valid.voucher_id.app_error", nil, "please provide valid voucher id", http.StatusBadRequest)
	}
	if !IsValidEmail(vc.CustomerEmail) {
		return NewAppError("VoucherCustomer.IsValid", "model.voucher_customer.is_valid.customer_email.app_error", nil, "please provide valid customer email", http.StatusBadRequest)
	}

	return nil
}

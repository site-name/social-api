package model

import (
	"github.com/Masterminds/squirrel"
)

type VoucherProductVariant struct {
	Id               string `json:"id"`
	VoucherID        string `json:"voucher_id"`
	ProductVariantID string `json:"product_variant_id"`
	CreateAt         int64  `json:"create_at"`
}

// VoucherProductVariantFilterOption is used to build squirrel sql queries
type VoucherProductVariantFilterOption struct {
	Id               squirrel.Sqlizer
	VoucherID        squirrel.Sqlizer
	ProductVariantID squirrel.Sqlizer
}

func (v *VoucherProductVariant) PreSave() {
	if !IsValidId(v.Id) {
		v.Id = NewId()
	}
	v.CreateAt = GetMillis()
}

func (v *VoucherProductVariant) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"model.voucher_product_variant.is_valid.%s.app_error",
		"voucher_product_variant_id=",
		"VoucherProductVariant.IsValid",
	)

	if !IsValidId(v.Id) {
		return outer("id", nil)
	}
	if !IsValidId(v.VoucherID) {
		return outer("voucher_id", &v.Id)
	}
	if !IsValidId(v.ProductVariantID) {
		return outer("product_variant_id", &v.Id)
	}
	if v.CreateAt <= 0 {
		return outer("create_at", &v.Id)
	}

	return nil
}

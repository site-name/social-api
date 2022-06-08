package product_and_discount

import (
	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/model"
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
	if !model.IsValidId(v.Id) {
		v.Id = model.NewId()
	}
	v.CreateAt = model.GetMillis()
}

func (v *VoucherProductVariant) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.voucher_product_variant.is_valid.%s.app_error",
		"voucher_product_variant_id=",
		"VoucherProductVariant.IsValid",
	)

	if !model.IsValidId(v.Id) {
		return outer("id", nil)
	}
	if !model.IsValidId(v.VoucherID) {
		return outer("voucher_id", &v.Id)
	}
	if !model.IsValidId(v.ProductVariantID) {
		return outer("product_variant_id", &v.Id)
	}
	if v.CreateAt <= 0 {
		return outer("create_at", &v.Id)
	}

	return nil
}

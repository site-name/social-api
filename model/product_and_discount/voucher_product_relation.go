package product_and_discount

import (
	"github.com/sitename/sitename/model"
)

// VoucherProduct represents relationship between vouchers and products
type VoucherProduct struct {
	Id        string `json:"id"`
	VoucherID string `json:"voucher_id"`
	ProductID string `json:"product_id"`
}

func (v *VoucherProduct) PreSave() {
	if v.Id == "" {
		v.Id = model.NewId()
	}
}

func (v *VoucherProduct) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.voucher_product.is_valid.%s.app_error",
		"voucher_product_id=",
		"VoucherProduct.IsValid",
	)

	if !model.IsValidId(v.Id) {
		return outer("id", nil)
	}
	if !model.IsValidId(v.VoucherID) {
		return outer("voucher_id", &v.Id)
	}
	if !model.IsValidId(v.ProductID) {
		return outer("product_id", &v.Id)
	}

	return nil
}

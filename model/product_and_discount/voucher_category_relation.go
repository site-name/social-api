package product_and_discount

import "github.com/sitename/sitename/model"

type VoucherCategory struct {
	Id         string `json:"id"`
	VoucherID  string `json:"voucher_id"`
	CategoryID string `json:"category_id"`
}

func (v *VoucherCategory) PreSave() {
	if v.Id == "" {
		v.Id = model.NewId()
	}
}

func (v *VoucherCategory) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.voucher_category.is_valid.%s.app_error",
		"voucher_category_id=",
		"VoucherCategory.IsValid",
	)
	if !model.IsValidId(v.Id) {
		return outer("id", nil)
	}
	if !model.IsValidId(v.VoucherID) {
		return outer("voucher_id", nil)
	}
	if !model.IsValidId(v.CategoryID) {
		return outer("category_id", nil)
	}

	return nil
}

package model

// VoucherProduct represents relationship between vouchers and products
type VoucherProduct struct {
	Id        string `json:"id"`
	VoucherID string `json:"voucher_id"`
	ProductID string `json:"product_id"`
}

func (v *VoucherProduct) PreSave() {
	if v.Id == "" {
		v.Id = NewId()
	}
}

func (v *VoucherProduct) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"voucher_product.is_valid.%s.app_error",
		"voucher_product_id=",
		"VoucherProduct.IsValid",
	)

	if !IsValidId(v.Id) {
		return outer("id", nil)
	}
	if !IsValidId(v.VoucherID) {
		return outer("voucher_id", &v.Id)
	}
	if !IsValidId(v.ProductID) {
		return outer("product_id", &v.Id)
	}

	return nil
}
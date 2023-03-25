package model

import (
	"github.com/Masterminds/squirrel"
)

type VoucherCategory struct {
	Id         string `json:"id"`
	VoucherID  string `json:"voucher_id"`
	CategoryID string `json:"category_id"`
	CreateAt   int64  `json:"create_at"` // this field is used to ordering
}

// VoucherCategoryFilterOption is used when building sql queries
type VoucherCategoryFilterOption struct {
	Id         squirrel.Sqlizer
	VoucherID  squirrel.Sqlizer
	CategoryID squirrel.Sqlizer
}

func (v *VoucherCategory) PreSave() {
	if v.Id == "" {
		v.Id = NewId()
	}
}

func (v *VoucherCategory) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"model.voucher_category.is_valid.%s.app_error",
		"voucher_category_id=",
		"VoucherCategory.IsValid",
	)
	if !IsValidId(v.Id) {
		return outer("id", nil)
	}
	if !IsValidId(v.VoucherID) {
		return outer("voucher_id", nil)
	}
	if !IsValidId(v.CategoryID) {
		return outer("category_id", nil)
	}

	return nil
}

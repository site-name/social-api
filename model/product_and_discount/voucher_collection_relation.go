package product_and_discount

import (
	"github.com/sitename/sitename/model"
)

// VoucherCollection represents voucher collection relationship
type VoucherCollection struct {
	Id           string `json:"id"`
	VoucherID    string `json:"voucher_id"`
	CollectionID string `json:"collection_id"`
}

func (v *VoucherCollection) PreSave() {
	if v.Id == "" {
		v.Id = model.NewId()
	}
}

func (v *VoucherCollection) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.voucher_collection.is_valid.%s.app_error",
		"voucher_collection_id=",
		"VoucherCollection.IsValid",
	)
	if !model.IsValidId(v.Id) {
		return outer("id", nil)
	}
	if !model.IsValidId(v.VoucherID) {
		return outer("voucher_id", &v.Id)
	}
	if !model.IsValidId(v.CollectionID) {
		return outer("collection_id", &v.Id)
	}

	return nil
}

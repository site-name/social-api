package model

// VoucherCollection represents voucher collection relationship
type VoucherCollection struct {
	Id           string `json:"id"`
	VoucherID    string `json:"voucher_id"`
	CollectionID string `json:"collection_id"`
}

func (v *VoucherCollection) PreSave() {
	if v.Id == "" {
		v.Id = NewId()
	}
}

func (v *VoucherCollection) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"voucher_collection.is_valid.%s.app_error",
		"voucher_collection_id=",
		"VoucherCollection.IsValid",
	)
	if !IsValidId(v.Id) {
		return outer("id", nil)
	}
	if !IsValidId(v.VoucherID) {
		return outer("voucher_id", &v.Id)
	}
	if !IsValidId(v.CollectionID) {
		return outer("collection_id", &v.Id)
	}

	return nil
}

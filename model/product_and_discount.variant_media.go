package model

import "github.com/Masterminds/squirrel"

type VariantMedia struct {
	Id        string `json:"id"`
	VariantID string `json:"variant_id"`
	MediaID   string `json:"media_id"`
}

type VariantMediaFilterOptions struct {
	Conditions squirrel.Sqlizer
}

func (v *VariantMedia) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"model.variant_product.is_valid.%s.app_error",
		"variant_product_id=",
		"VariantProduct.IsValid",
	)

	if !IsValidId(v.Id) {
		return outer("id", nil)
	}
	if !IsValidId(v.VariantID) {
		return outer("variant_id", &v.Id)
	}
	if !IsValidId(v.MediaID) {
		return outer("media_id", &v.Id)
	}

	return nil
}

func (v *VariantMedia) ToJSON() string {
	return ModelToJson(v)
}

func (v *VariantMedia) PreSave() {
	if v.Id == "" {
		v.Id = NewId()
	}
}

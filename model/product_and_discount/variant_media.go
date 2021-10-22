package product_and_discount

import (
	"io"

	"github.com/sitename/sitename/model"
)

type VariantMedia struct {
	Id        string `json:"id"`
	VariantID string `json:"variant_id"`
	MediaID   string `json:"media_id"`
}

func (v *VariantMedia) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.variant_product.is_valid.%s.app_error",
		"variant_product_id=",
		"VariantProduct.IsValid",
	)

	if !model.IsValidId(v.Id) {
		return outer("id", nil)
	}
	if !model.IsValidId(v.VariantID) {
		return outer("variant_id", &v.Id)
	}
	if !model.IsValidId(v.MediaID) {
		return outer("media_id", &v.Id)
	}

	return nil
}

func (v *VariantMedia) ToJSON() string {
	return model.ModelToJson(v)
}

func VariantMediaToJson(data io.Reader) *VariantMedia {
	var v VariantMedia
	model.ModelFromJson(&v, data)
	return &v
}

func (v *VariantMedia) PreSave() {
	if v.Id == "" {
		v.Id = model.NewId()
	}
}

package product_and_discount

import (
	"io"

	"github.com/sitename/sitename/model"
)

type ProductMedia struct {
	Id        string `json:"id"`
	ProductID string `json:"product_id"`
	*model.Sortable
}

// TODO: not done yet
func (p *ProductMedia) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.product_media.is_valid.%s.app_error",
		"product_media_id=",
		"ProductMedia.IsValid")
	if p.Id == "" {
		return outer("id", nil)
	}
	if p.ProductID == "" {
		return outer("product_id", &p.Id)
	}

	return nil
}

func (p *ProductMedia) ToJson() string {
	return model.ModelToJson(p)
}

func ProductMediaFromJson(data io.Reader) *ProductMedia {
	var prd ProductMedia
	model.ModelFromJson(&prd, data)
	return &prd
}

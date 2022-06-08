package product_and_discount

import (
	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/model"
)

type SaleProductVariant struct {
	Id               string `json:"id"`
	SaleID           string `json:"sale_id"`
	ProductVariantID string `json:"product_variant_id"`
	CreateAt         int64  `json:"create_at"`
}

// SaleProductVariantFilterOption is used to build squirrel sql queries
type SaleProductVariantFilterOption struct {
	Id               squirrel.Sqlizer
	SaleID           squirrel.Sqlizer
	ProductVariantID squirrel.Sqlizer
}

func (s *SaleProductVariant) PreSave() {
	if !model.IsValidId(s.Id) {
		s.Id = model.NewId()
	}
	s.CreateAt = model.GetMillis()
}

func (s *SaleProductVariant) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.sale_product_variant.is_valid.%s.app_error",
		"sale_product_variant_id=",
		"SaleProductVariant.IsValid",
	)

	if !model.IsValidId(s.Id) {
		return outer("id", nil)
	}
	if !model.IsValidId(s.SaleID) {
		return outer("sale_id", &s.Id)
	}
	if !model.IsValidId(s.ProductVariantID) {
		return outer("product_variant_id", &s.Id)
	}
	if s.CreateAt <= 0 {
		return outer("create_at", &s.Id)
	}

	return nil
}

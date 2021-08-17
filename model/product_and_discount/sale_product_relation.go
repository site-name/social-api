package product_and_discount

import "github.com/sitename/sitename/model"

type SaleProductRelation struct {
	Id        string `json:"id"`
	SaleID    string `json:"sale_id"`
	ProductID string `json:"product_id"`
	CreateAt  int64  `json:"create_at"`
}

type SaleProductRelationFilterOption struct {
	Id        *model.StringFilter
	SaleID    *model.StringFilter
	ProductID *model.StringFilter
}

func (s *SaleProductRelation) PreSave() {
	if s.Id == "" {
		s.Id = model.NewId()
	}
	s.CreateAt = model.GetMillis()
}

func (s *SaleProductRelation) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.sale_product_relation.is_valid.%s.app_error",
		"sale_product_relation_id=",
		"SaleCategory.IsValid",
	)

	if !model.IsValidId(s.Id) {
		return outer("id", nil)
	}
	if s.CreateAt == 0 {
		return outer("create_at", &s.Id)
	}
	if !model.IsValidId(s.SaleID) {
		return outer("sale_id", &s.Id)
	}
	if !model.IsValidId(s.ProductID) {
		return outer("product_id", &s.Id)
	}

	return nil
}

package product_and_discount

import (
	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/model"
)

type SaleCategoryRelation struct {
	Id         string `json:"id"`
	SaleID     string `json:"sale_id"`
	CategoryID string `json:"category_id"`
	CreateAt   int64  `json:"create_at"`
}

type SaleCategoryRelationFilterOption struct {
	Id         squirrel.Sqlizer
	SaleID     squirrel.Sqlizer
	CategoryID squirrel.Sqlizer
}

func (s *SaleCategoryRelation) PreSave() {
	if s.Id == "" {
		s.Id = model.NewId()
	}
	s.CreateAt = model.GetMillis()
}

func (s *SaleCategoryRelation) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.sale_category_relation.is_valid.%s.app_error",
		"sale_category_relation_id=",
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
	if !model.IsValidId(s.CategoryID) {
		return outer("category_id", &s.Id)
	}

	return nil
}

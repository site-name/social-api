package model

import (
	"github.com/Masterminds/squirrel"
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
		s.Id = NewId()
	}
	s.CreateAt = GetMillis()
}

func (s *SaleCategoryRelation) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"sale_category_relation.is_valid.%s.app_error",
		"sale_category_relation_id=",
		"SaleCategory.IsValid",
	)

	if !IsValidId(s.Id) {
		return outer("id", nil)
	}
	if s.CreateAt == 0 {
		return outer("create_at", &s.Id)
	}
	if !IsValidId(s.SaleID) {
		return outer("sale_id", &s.Id)
	}
	if !IsValidId(s.CategoryID) {
		return outer("category_id", &s.Id)
	}

	return nil
}

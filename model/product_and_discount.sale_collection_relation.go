package model

import (
	"github.com/Masterminds/squirrel"
)

type SaleCollectionRelation struct {
	Id           string `json:"id"`
	SaleID       string `json:"sale_id"`
	CollectionID string `json:"collection_id"`
	CreateAt     int64  `json:"create_at"`
}

// SaleCollectionRelationFilterOption is used to build sql queries
type SaleCollectionRelationFilterOption struct {
	Id           squirrel.Sqlizer
	SaleID       squirrel.Sqlizer
	CollectionID squirrel.Sqlizer
}

func (s *SaleCollectionRelation) PreSave() {
	if s.Id == "" {
		s.Id = NewId()
	}
	s.CreateAt = GetMillis()
}

func (s *SaleCollectionRelation) IsValid() *AppError {
	outer := CreateAppErrorForModel(
		"model.sale_collection_relation.is_valid.%s.app_error",
		"sale_collection_relation_id=",
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
	if !IsValidId(s.CollectionID) {
		return outer("collection_id", &s.Id)
	}

	return nil
}

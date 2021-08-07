package product_and_discount

import "github.com/sitename/sitename/model"

type SaleCollectionRelation struct {
	Id           string `json:"id"`
	SaleID       string `json:"sale_id"`
	CollectionID string `json:"collection_id"`
	CreateAt     int64  `json:"create_at"`
}

func (s *SaleCollectionRelation) PreSave() {
	if s.Id == "" {
		s.Id = model.NewId()
	}
	if s.CreateAt == 0 {
		s.CreateAt = model.GetMillis()
	}
}

func (s *SaleCollectionRelation) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.sale_collection_relation.is_valid.%s.app_error",
		"sale_collection_relation_id=",
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
	if !model.IsValidId(s.CollectionID) {
		return outer("collection_id", &s.Id)
	}

	return nil
}

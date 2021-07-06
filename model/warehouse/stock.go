package warehouse

import (
	"io"

	"github.com/sitename/sitename/model"
)

type Stock struct {
	Id               string `json:"id"`
	WarehouseID      string `json:"warehouse_id"`       // NOT NULL
	ProductVariantID string `json:"product_variant_id"` // NOT NULL
	Quantity         uint   `json:"quantity"`           // DEFAULT 0
}

func (s *Stock) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.stock.is_valid.%s.app_error",
		"stock_id=",
		"Stock.IsValid",
	)
	if !model.IsValidId(s.Id) {
		return outer("id", nil)
	}
	if !model.IsValidId(s.WarehouseID) {
		return outer("warehouse_id", &s.Id)
	}
	if !model.IsValidId(s.ProductVariantID) {
		return outer("product_variant_id", &s.Id)
	}

	return nil
}

func (s *Stock) ToJson() string {
	return model.ModelToJson(s)
}

func StockFromJson(data io.Reader) *Stock {
	var s Stock
	model.ModelFromJson(&s, data)
	return &s
}

func (s *Stock) PreSave() {
	if s.Id == "" {
		s.Id = model.NewId()
	}
}

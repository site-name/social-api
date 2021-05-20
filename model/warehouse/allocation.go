package warehouse

import (
	"io"

	"github.com/sitename/sitename/model"
)

type Allocation struct {
	Id                string `json:"id"`
	OrderLineID       string `json:"order_ldine_id"` // NOT NULL
	StockID           string `json:"stock_id"`       // NOT NULL
	QuantityAllocated uint64 `json:"quantity_allocated"`
}

func (a *Allocation) IsValid() *model.AppError {
	outer := model.CreateAppErrorForModel(
		"model.allocation.is_valid.%s.app_error",
		"allocation_id=",
		"Allocation.isValid",
	)
	if !model.IsValidId(a.Id) {
		return outer("id", nil)
	}
	if !model.IsValidId(a.OrderLineID) {
		return outer("order_line_id", &a.Id)
	}
	if !model.IsValidId(a.StockID) {
		return outer("stock_id", &a.Id)
	}

	return nil
}

func (a *Allocation) ToJson() string {
	return model.ModelToJson(a)
}

func AllocationFromJson(data io.Reader) *Allocation {
	var a Allocation
	model.ModelFromJson(&a, data)
	return &a
}

func (a *Allocation) PreSave() {
	if a.Id == "" {
		a.Id = model.NewId()
	}
}

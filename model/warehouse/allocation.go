package warehouse

import (
	"github.com/sitename/sitename/model"
)

type Allocation struct {
	Id                string `json:"id"`
	CreateAt          int64  `json:"create_at"`
	OrderLineID       string `json:"order_ldine_id"`     // NOT NULL
	StockID           string `json:"stock_id"`           // NOT NULL
	QuantityAllocated uint   `json:"quantity_allocated"` // default 0
}

// AllocationFilterOption is used to build sql queries to filtering warehouse allocations
type AllocationFilterOption struct {
	Id          *model.StringFilter
	OrderLineID *model.StringFilter
	StockID     *model.StringFilter
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
	if a.CreateAt == 0 {
		return outer("create_at", &a.Id)
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

func (a *Allocation) PreSave() {
	if a.Id == "" {
		a.Id = model.NewId()
	}
	if a.CreateAt == 0 {
		a.CreateAt = model.GetMillis()
	}
}
